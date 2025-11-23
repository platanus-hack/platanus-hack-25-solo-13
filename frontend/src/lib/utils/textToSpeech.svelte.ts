import { auth } from '$lib/stores/auth.svelte';

/**
 * Text-to-Speech utility using ElevenLabs via backend proxy
 * Manages audio playback state and provides controls
 */

export class TextToSpeechPlayer {
	private audio: HTMLAudioElement | null = $state(null);
	private currentText: string = $state('');

	isPlaying = $state(false);
	isLoading = $state(false);
	error = $state<string | null>(null);

	constructor() {}

	/**
	 * Play text as speech
	 * @param text - Text to convert to speech
	 */
	async play(text: string) {
		// If same text is already playing, just pause/resume
		if (this.currentText === text && this.audio) {
			if (this.isPlaying) {
				this.pause();
			} else {
				this.resume();
			}
			return;
		}

		// Stop any currently playing audio
		this.stop();

		this.isLoading = true;
		this.error = null;
		this.currentText = text;

		try {
			// Call backend to generate audio
			const response = await fetch('/api/tts/generate', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					'Authorization': auth.token ? `Bearer ${auth.token}` : ''
				},
				body: JSON.stringify({ text })
			});

			if (!response.ok) {
				throw new Error('Failed to generate audio');
			}

			// Get audio blob from response
			const audioBlob = await response.blob();
			const audioUrl = URL.createObjectURL(audioBlob);

			// Create audio element
			this.audio = new Audio(audioUrl);

			// Setup event listeners
			this.audio.onplay = () => {
				this.isPlaying = true;
			};

			this.audio.onpause = () => {
				this.isPlaying = false;
			};

			this.audio.onended = () => {
				this.isPlaying = false;
				this.currentText = '';
				if (this.audio) {
					URL.revokeObjectURL(this.audio.src);
				}
			};

			this.audio.onerror = () => {
				this.error = 'Error playing audio';
				this.isPlaying = false;
				this.isLoading = false;
			};

			// Start playback
			await this.audio.play();
			this.isLoading = false;
		} catch (err) {
			this.error = err instanceof Error ? err.message : 'Unknown error';
			this.isLoading = false;
			console.error('TTS error:', err);
		}
	}

	/**
	 * Pause current audio
	 */
	pause() {
		if (this.audio && !this.audio.paused) {
			this.audio.pause();
		}
	}

	/**
	 * Resume paused audio
	 */
	resume() {
		if (this.audio && this.audio.paused) {
			this.audio.play();
		}
	}

	/**
	 * Stop and cleanup audio
	 */
	stop() {
		if (this.audio) {
			this.audio.pause();
			this.audio.currentTime = 0;
			URL.revokeObjectURL(this.audio.src);
			this.audio = null;
		}
		this.isPlaying = false;
		this.currentText = '';
	}

	/**
	 * Cleanup when component is destroyed
	 */
	destroy() {
		this.stop();
	}
}

/**
 * Create a new TTS player instance
 */
export function createTTSPlayer() {
	return new TextToSpeechPlayer();
}
