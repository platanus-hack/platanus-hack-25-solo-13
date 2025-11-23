package handlers

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// TTSRequest es el payload para generar audio
type TTSRequest struct {
	Text string `json:"text"`
}

// GenerateTTSHandler genera audio usando ElevenLabs API con caché
// POST /api/tts/generate
func GenerateTTSHandler(w http.ResponseWriter, r *http.Request) {
	var req TTSRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	if req.Text == "" {
		http.Error(w, `{"error":"text is required"}`, http.StatusBadRequest)
		return
	}

	// Generate hash for cache key
	hash := md5.Sum([]byte(req.Text))
	cacheKey := hex.EncodeToString(hash[:])

	// Setup cache directory
	cacheDir := "./static/tts-cache"
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		fmt.Printf("Warning: failed to create cache directory: %v\n", err)
	}

	cachePath := filepath.Join(cacheDir, cacheKey+".mp3")

	// Check if cached file exists
	if _, err := os.Stat(cachePath); err == nil {
		// Cache hit - serve cached file
		fmt.Printf("TTS Cache HIT: %s\n", cacheKey)

		audioData, err := os.ReadFile(cachePath)
		if err != nil {
			fmt.Printf("Error reading cached file: %v\n", err)
			// Continue to generation if cache read fails
		} else {
			w.Header().Set("Content-Type", "audio/mpeg")
			w.Header().Set("Cache-Control", "public, max-age=86400")
			w.Header().Set("X-Cache", "HIT")
			w.Write(audioData)
			return
		}
	}

	// Cache miss - generate new audio
	fmt.Printf("TTS Cache MISS: %s - generating...\n", cacheKey)

	// Get API key from environment
	apiKey := os.Getenv("ELEVENLABS_API_KEY")
	if apiKey == "" {
		http.Error(w, `{"error":"ElevenLabs API key not configured"}`, http.StatusInternalServerError)
		return
	}

	// Use a fast, multilingual voice
	voiceID := "cLzIVykddLltvgkzos6C"

	// Prepare ElevenLabs API request
	elevenLabsURL := fmt.Sprintf("https://api.elevenlabs.io/v1/text-to-speech/%s", voiceID)

	payload := map[string]interface{}{
		"text": req.Text,
		"model_id": "eleven_multilingual_v2",
		"voice_settings": map[string]interface{}{
			"stability": 0.5,
			"similarity_boost": 0.75,
			"style": 0.0,
			"use_speaker_boost": true,
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, `{"error":"failed to prepare request"}`, http.StatusInternalServerError)
		return
	}

	// Create request to ElevenLabs
	elevenLabsReq, err := http.NewRequest("POST", elevenLabsURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		http.Error(w, `{"error":"failed to create request"}`, http.StatusInternalServerError)
		return
	}

	elevenLabsReq.Header.Set("Content-Type", "application/json")
	elevenLabsReq.Header.Set("xi-api-key", apiKey)

	// Make request to ElevenLabs
	client := &http.Client{}
	resp, err := client.Do(elevenLabsReq)
	if err != nil {
		http.Error(w, `{"error":"failed to call ElevenLabs API"}`, http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Printf("ElevenLabs API error: %s\n", string(bodyBytes))
		http.Error(w, `{"error":"ElevenLabs API error"}`, resp.StatusCode)
		return
	}

	// Read audio data
	audioData, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, `{"error":"failed to read audio data"}`, http.StatusInternalServerError)
		return
	}

	// Save to cache (best effort - don't fail if cache write fails)
	if err := os.WriteFile(cachePath, audioData, 0644); err != nil {
		fmt.Printf("Warning: failed to cache audio: %v\n", err)
	} else {
		fmt.Printf("TTS cached successfully: %s\n", cacheKey)
	}

	// Stream audio response back to client
	w.Header().Set("Content-Type", "audio/mpeg")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Header().Set("X-Cache", "MISS")
	w.Write(audioData)
}
