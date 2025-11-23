<script>
	import { onMount } from 'svelte';
	import gsap from 'gsap';

	// Props esperadas desde el backend
	let {
		titulo = '',
		bloques = [],
		onNext = null,
		onPrevious = null,
		showNavigation = true
	} = $props();

	let bloqueRefs = [];

	onMount(() => {
		// Animate blocks on mount with stagger
		gsap.from(bloqueRefs, {
			opacity: 0,
			y: 20,
			duration: 0.5,
			stagger: 0.1,
			ease: 'power2.out'
		});
	});

	// Obtener ícono según el estilo de nota
	function getNotaIcon(estilo) {
		switch (estilo) {
			case 'info':
				return 'ℹ️';
			case 'warning':
				return '⚠️';
			case 'tip':
				return '💡';
			default:
				return 'ℹ️';
		}
	}

	// Obtener color de fondo según el estilo de nota
	function getNotaColor(estilo) {
		switch (estilo) {
			case 'info':
				return 'bg-blue-50 border-blue-300 text-blue-900';
			case 'warning':
				return 'bg-yellow-50 border-yellow-300 text-yellow-900';
			case 'tip':
				return 'bg-green-50 border-green-300 text-green-900';
			default:
				return 'bg-gray-50 border-gray-300 text-gray-900';
		}
	}
</script>

<div class="w-full max-w-4xl mx-auto p-6">
	<!-- Título principal del componente -->
	<h2 class="text-3xl font-bold mb-8 text-gray-900">{titulo}</h2>

	<!-- Renderizar bloques secuencialmente -->
	<div class="space-y-6">
		{#each bloques as bloque, i}
			<div bind:this={bloqueRefs[i]}>
				<!-- BLOQUE TEXTO -->
				{#if bloque.tipo === 'texto'}
					<div class="prose prose-lg max-w-none">
						<p class="text-gray-700 leading-relaxed whitespace-pre-line">
							{bloque.contenido}
						</p>
					</div>

				<!-- BLOQUE EJEMPLO -->
				{:else if bloque.tipo === 'ejemplo'}
					<div class="bg-purple-50 border-l-4 border-purple-500 rounded-r-lg p-6 shadow-sm">
						<h3 class="text-xl font-semibold text-purple-900 mb-3">
							✨ {bloque.titulo}
						</h3>
						<div class="bg-white rounded-lg p-4 mb-3">
							<p class="text-gray-800 whitespace-pre-line">{bloque.contenido}</p>
						</div>
						{#if bloque.analisis}
							<div class="mt-3 text-sm text-purple-700 bg-purple-100 rounded-lg p-3">
								<strong>Análisis:</strong> {bloque.analisis}
							</div>
						{/if}
					</div>

				<!-- BLOQUE DEFINICIÓN -->
				{:else if bloque.tipo === 'definicion'}
					<div class="bg-indigo-50 border-l-4 border-indigo-600 rounded-r-lg p-5">
						<dt class="text-lg font-bold text-indigo-900 mb-2">
							{bloque.termino}
						</dt>
						<dd class="text-gray-700 italic">{bloque.texto}</dd>
					</div>

				<!-- BLOQUE NOTA -->
				{:else if bloque.tipo === 'nota'}
					<div class={`border-l-4 rounded-r-lg p-5 flex gap-3 ${getNotaColor(bloque.estilo || 'info')}`}>
						<span class="text-2xl flex-shrink-0">{getNotaIcon(bloque.estilo || 'info')}</span>
						<p class="whitespace-pre-line">{bloque.texto}</p>
					</div>

				<!-- BLOQUE EJERCICIO -->
				{:else if bloque.tipo === 'ejercicio'}
					<div class="bg-orange-50 border border-orange-300 rounded-lg p-6 shadow-sm">
						<div class="flex items-center gap-2 mb-4">
							<span class="text-2xl">✏️</span>
							<h3 class="text-xl font-semibold text-orange-900">Practica</h3>
						</div>
						<div class="bg-white rounded-lg p-4 border border-orange-200">
							<p class="text-gray-800 font-medium mb-3">{bloque.instruccion}</p>
							{#if bloque.ejemplo}
								<div class="mt-3 pt-3 border-t border-orange-200">
									<p class="text-sm text-orange-700">
										<strong>Ejemplo:</strong> {bloque.ejemplo}
									</p>
								</div>
							{/if}
						</div>
					</div>

				<!-- BLOQUE RESUMEN -->
				{:else if bloque.tipo === 'resumen'}
					<div class="bg-gray-100 rounded-lg p-6 shadow-sm">
						<h3 class="text-xl font-semibold text-gray-900 mb-4 flex items-center gap-2">
							<span>📝</span>
							<span>Resumen</span>
						</h3>
						<ul class="space-y-2">
							{#each bloque.puntos as punto}
								<li class="flex items-start gap-3">
									<span class="text-purple-600 font-bold mt-1">•</span>
									<span class="text-gray-800">{punto}</span>
								</li>
							{/each}
						</ul>
					</div>

				<!-- BLOQUE COMPARACIÓN -->
				{:else if bloque.tipo === 'comparacion'}
					<div class="bg-gray-50 rounded-lg p-6 shadow-sm overflow-x-auto">
						<h3 class="text-xl font-semibold text-gray-900 mb-4 flex items-center gap-2">
							<span>⚖️</span>
							<span>Comparación</span>
						</h3>
						<table class="w-full border-collapse">
							<thead>
								<tr class="bg-gray-200">
									<th class="border border-gray-300 px-4 py-2 text-left font-semibold text-gray-800">
										Aspecto
									</th>
									<th class="border border-gray-300 px-4 py-2 text-left font-semibold text-gray-800">
										Opción 1
									</th>
									<th class="border border-gray-300 px-4 py-2 text-left font-semibold text-gray-800">
										Opción 2
									</th>
								</tr>
							</thead>
							<tbody>
								{#each bloque.items as item}
									<tr class="hover:bg-gray-100 transition-colors">
										<td class="border border-gray-300 px-4 py-3 font-medium text-gray-700">
											{item.aspecto}
										</td>
										<td class="border border-gray-300 px-4 py-3 text-gray-700">
											{item.opcion1}
										</td>
										<td class="border border-gray-300 px-4 py-3 text-gray-700">
											{item.opcion2}
										</td>
									</tr>
								{/each}
							</tbody>
						</table>
					</div>

				<!-- TIPO DE BLOQUE DESCONOCIDO -->
				{:else}
					<div class="bg-red-100 border border-red-400 rounded-lg p-4 text-red-800">
						<p class="font-semibold">Tipo de bloque desconocido: {bloque.tipo}</p>
					</div>
				{/if}
			</div>
		{/each}
	</div>

	<!-- Navegación -->
	{#if showNavigation}
		<div class="flex items-center justify-between pt-6 mt-8 border-t border-gray-300">
			<button
				onclick={onPrevious}
				disabled={!onPrevious}
				class="px-6 py-3 rounded-xl font-semibold bg-gray-200 text-gray-700
				       border border-gray-300 transition-all duration-300
				       hover:bg-gray-300 hover:scale-105
				       disabled:opacity-30 disabled:cursor-not-allowed disabled:hover:scale-100"
			>
				← Anterior
			</button>

			<div class="text-center">
				<p class="text-xs text-gray-500">
					{bloques.length} bloques
				</p>
			</div>

			<button
				onclick={onNext}
				disabled={!onNext}
				class="px-6 py-3 rounded-xl font-semibold
				       bg-gradient-to-r from-blue-500 to-purple-500 text-white
				       transition-all duration-300
				       hover:shadow-lg hover:shadow-blue-500/50 hover:scale-105
				       disabled:opacity-30 disabled:cursor-not-allowed disabled:hover:scale-100"
			>
				Siguiente →
			</button>
		</div>
	{/if}
</div>

<style>
	/* Mejorar legibilidad de texto */
	.prose p {
		line-height: 1.7;
	}

	/* Animaciones suaves en hover */
	.hover\:bg-gray-100 {
		transition: background-color 0.2s ease;
	}

	/* Asegurar que el whitespace-pre-line mantenga los saltos de línea */
	.whitespace-pre-line {
		white-space: pre-line;
	}
</style>
