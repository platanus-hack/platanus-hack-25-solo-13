<script>
  import { onMount } from 'svelte';
  import gsap from 'gsap';

  // Props
  let {
    // Data del componente
    title = "Crea un mapa conceptual sobre el tema",
    topic = "Fotosíntesis",
    instructions = "Identifica los conceptos clave y sus relaciones",
    requiredConcepts = [], // Conceptos que deben aparecer
    suggestedConnections = [], // { from: "concepto1", to: "concepto2", label: "relación" }
    minConcepts = 4,
    minConnections = 3,

    // Metadata educativa
    bloomLevel = "crear",
    materia = "biología",
    oaId = null,

    // Configuración
    showFeedback = true,
    allowMultipleAttempts = true,
    provideConcepts = false, // Si true, los conceptos están pre-definidos

    // Callbacks
    onAnswer = null,
    onComplete = null
  } = $props();

  // Estados locales
  let nodes = $state(provideConcepts ? requiredConcepts.map((c, i) => ({
    id: i + 1,
    text: c,
    x: 100 + (i % 3) * 200,
    y: 100 + Math.floor(i / 3) * 150
  })) : []);
  let connections = $state([]);
  let hasSubmitted = $state(false);
  let attemptCount = $state(0);
  let containerRef = $state(null);
  let canvasRef = $state(null);

  // Estados de interacción
  let newNodeText = $state("");
  let selectedNodes = $state([]); // Para crear conexiones
  let connectionLabel = $state("");
  let showConnectionDialog = $state(false);
  let draggedNode = $state(null);

  // Colores por nivel de Bloom
  const bloomColors = {
    recordar: { bg: 'bg-red-500/20', border: 'border-red-500', text: 'text-red-400' },
    comprender: { bg: 'bg-orange-500/20', border: 'border-orange-500', text: 'text-orange-400' },
    aplicar: { bg: 'bg-yellow-500/20', border: 'border-yellow-500', text: 'text-yellow-400' },
    analizar: { bg: 'bg-green-500/20', border: 'border-green-500', text: 'text-green-400' },
    evaluar: { bg: 'bg-blue-500/20', border: 'border-blue-500', text: 'text-blue-400' },
    crear: { bg: 'bg-purple-500/20', border: 'border-purple-500', text: 'text-purple-400' }
  };

  const materiaColors = {
    matemáticas: 'cyan',
    lenguaje: 'purple',
    historia: 'amber',
    física: 'blue',
    química: 'green',
    biología: 'emerald',
    default: 'slate'
  };

  const currentMateriaColor = materiaColors[materia] || materiaColors.default;
  const currentBloomStyle = bloomColors[bloomLevel] || bloomColors.crear;

  // Verificar completitud
  const meetsRequirements = $derived(
    nodes.length >= minConcepts && connections.length >= minConnections
  );

  // Calcular score
  const score = $derived.by(() => {
    if (!hasSubmitted) return 0;

    let conceptScore = 0;
    let connectionScore = 0;

    // Verificar conceptos requeridos
    if (requiredConcepts.length > 0) {
      const foundConcepts = requiredConcepts.filter(required =>
        nodes.some(node =>
          node.text.toLowerCase().includes(required.toLowerCase()) ||
          required.toLowerCase().includes(node.text.toLowerCase())
        )
      );
      conceptScore = (foundConcepts.length / requiredConcepts.length) * 50;
    } else {
      // Si no hay conceptos requeridos, dar puntos por cantidad
      conceptScore = Math.min((nodes.length / minConcepts) * 50, 50);
    }

    // Verificar conexiones sugeridas
    if (suggestedConnections.length > 0) {
      const foundConnections = suggestedConnections.filter(suggested => {
        return connections.some(conn => {
          const fromNode = nodes.find(n => n.id === conn.from);
          const toNode = nodes.find(n => n.id === conn.to);

          const matchesFrom = fromNode?.text.toLowerCase().includes(suggested.from.toLowerCase());
          const matchesTo = toNode?.text.toLowerCase().includes(suggested.to.toLowerCase());

          return matchesFrom && matchesTo;
        });
      });
      connectionScore = (foundConnections.length / suggestedConnections.length) * 50;
    } else {
      // Si no hay conexiones sugeridas, dar puntos por cantidad
      connectionScore = Math.min((connections.length / minConnections) * 50, 50);
    }

    return conceptScore + connectionScore;
  });

  const allCorrect = $derived(score >= 70); // 70% o más se considera éxito

  // Funciones de nodos
  function addNode() {
    if (!newNodeText.trim()) return;

    const newNode = {
      id: Date.now(),
      text: newNodeText.trim(),
      x: 150 + (nodes.length % 4) * 150,
      y: 150 + Math.floor(nodes.length / 4) * 120
    };

    nodes = [...nodes, newNode];
    newNodeText = "";

    // Animación
    setTimeout(() => {
      const nodeEl = document.getElementById(`node-${newNode.id}`);
      if (nodeEl) {
        gsap.from(nodeEl, {
          scale: 0,
          duration: 0.3,
          ease: 'back.out(1.7)'
        });
      }
    }, 10);
  }

  function deleteNode(nodeId) {
    nodes = nodes.filter(n => n.id !== nodeId);
    connections = connections.filter(c => c.from !== nodeId && c.to !== nodeId);
    selectedNodes = selectedNodes.filter(id => id !== nodeId);
  }

  function selectNode(nodeId) {
    if (hasSubmitted && !allowMultipleAttempts) return;

    if (selectedNodes.includes(nodeId)) {
      selectedNodes = selectedNodes.filter(id => id !== nodeId);
    } else if (selectedNodes.length < 2) {
      selectedNodes = [...selectedNodes, nodeId];

      // Si se seleccionaron 2 nodos, mostrar diálogo de conexión
      if (selectedNodes.length === 2) {
        showConnectionDialog = true;
      }
    }
  }

  function createConnection() {
    if (selectedNodes.length !== 2 || !connectionLabel.trim()) return;

    const newConnection = {
      id: Date.now(),
      from: selectedNodes[0],
      to: selectedNodes[1],
      label: connectionLabel.trim()
    };

    connections = [...connections, newConnection];
    selectedNodes = [];
    connectionLabel = "";
    showConnectionDialog = false;
  }

  function cancelConnection() {
    selectedNodes = [];
    connectionLabel = "";
    showConnectionDialog = false;
  }

  function deleteConnection(connId) {
    connections = connections.filter(c => c.id !== connId);
  }

  // Drag node (simple version)
  function handleNodeDragStart(event, nodeId) {
    draggedNode = nodeId;
  }

  function handleNodeDrag(event, nodeId) {
    if (draggedNode !== nodeId) return;

    const canvas = canvasRef;
    if (!canvas) return;

    const rect = canvas.getBoundingClientRect();
    const x = event.clientX - rect.left;
    const y = event.clientY - rect.top;

    if (x > 0 && y > 0 && x < rect.width && y < rect.height) {
      nodes = nodes.map(n =>
        n.id === nodeId ? { ...n, x, y } : n
      );
    }
  }

  function handleNodeDragEnd() {
    draggedNode = null;
  }

  // Verificar respuestas
  function handleSubmit() {
    if (!meetsRequirements) {
      alert(`Por favor, crea al menos ${minConcepts} conceptos y ${minConnections} conexiones.`);
      return;
    }

    hasSubmitted = true;
    attemptCount++;

    // Animación
    if (containerRef) {
      gsap.to(containerRef, {
        scale: 0.98,
        duration: 0.1,
        yoyo: true,
        repeat: 1
      });
    }

    // Callback para backend
    if (onAnswer) {
      onAnswer({
        oaId: oaId,
        bloomLevel: bloomLevel,
        materia: materia,
        nodes: nodes,
        connections: connections,
        score: score,
        attemptCount: attemptCount,
        timestamp: new Date().toISOString()
      });
    }

    // Si está bien, llamar onComplete
    if (allCorrect && onComplete) {
      setTimeout(() => {
        onComplete({
          oaId: oaId,
          bloomLevel: bloomLevel,
          score: score,
          attempts: attemptCount
        });
      }, 1500);
    }
  }

  function handleTryAgain() {
    if (allowMultipleAttempts) {
      hasSubmitted = false;
    }
  }

  // Animación de entrada
  onMount(() => {
    if (containerRef) {
      gsap.from(containerRef, {
        opacity: 0,
        y: 20,
        duration: 0.5,
        ease: 'power2.out'
      });
    }
  });
</script>

<div
  bind:this={containerRef}
  class="w-full max-w-7xl mx-auto p-6 bg-canvas-950 rounded-2xl border border-slate-800 shadow-2xl"
>
  <!-- Header -->
  <div class="flex items-center justify-between mb-6">
    <div class="flex items-center gap-3">
      <span class="px-3 py-1 rounded-full text-xs font-semibold uppercase {currentBloomStyle.bg} {currentBloomStyle.border} {currentBloomStyle.text} border">
        {bloomLevel}
      </span>
      <span class="px-3 py-1 rounded-full text-xs font-semibold bg-{currentMateriaColor}-500/20 text-{currentMateriaColor}-400 border border-{currentMateriaColor}-500">
        {materia}
      </span>
    </div>
    {#if attemptCount > 0}
      <span class="text-xs text-slate-500">Intento {attemptCount}</span>
    {/if}
  </div>

  <!-- Título e instrucciones -->
  <div class="mb-6">
    <h3 class="text-xl font-semibold text-white mb-2">
      {title}: <span class="text-purple-400">{topic}</span>
    </h3>
    <p class="text-slate-400 text-sm">
      {instructions}
    </p>
  </div>

  <div class="grid grid-cols-1 lg:grid-cols-4 gap-6">
    <!-- Panel de control -->
    <div class="lg:col-span-1 space-y-4">
      <!-- Agregar nodo -->
      {#if !provideConcepts && (!hasSubmitted || allowMultipleAttempts)}
        <div class="p-4 bg-canvas-900/50 rounded-2xl border border-slate-700">
          <h4 class="text-sm font-semibold text-slate-300 mb-3">Agregar Concepto</h4>
          <input
            type="text"
            bind:value={newNodeText}
            onkeydown={(e) => e.key === 'Enter' && addNode()}
            placeholder="Ej: Clorofila"
            class="w-full px-3 py-2 bg-canvas-950 border border-slate-700 rounded-lg text-white text-sm
              focus:outline-none focus:border-purple-500 mb-2"
          />
          <button
            onclick={addNode}
            disabled={!newNodeText.trim()}
            class="w-full px-3 py-2 bg-purple-500 text-white rounded-lg text-sm font-semibold
              hover:bg-purple-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          >
            + Agregar
          </button>
        </div>
      {/if}

      <!-- Estadísticas -->
      <div class="p-4 bg-canvas-900/50 rounded-2xl border border-slate-700">
        <h4 class="text-sm font-semibold text-slate-300 mb-3">Estadísticas</h4>
        <div class="space-y-2 text-sm">
          <div class="flex justify-between">
            <span class="text-slate-400">Conceptos:</span>
            <span class="text-white font-semibold">{nodes.length}/{minConcepts}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-slate-400">Conexiones:</span>
            <span class="text-white font-semibold">{connections.length}/{minConnections}</span>
          </div>
        </div>

        {#if meetsRequirements && !hasSubmitted}
          <div class="mt-3 p-2 bg-green-500/10 rounded-lg border border-green-500/50">
            <p class="text-green-400 text-xs text-center">✓ Listo para enviar</p>
          </div>
        {/if}
      </div>

      <!-- Instrucciones -->
      <div class="p-4 bg-purple-500/10 rounded-2xl border border-purple-500/30">
        <h4 class="text-sm font-semibold text-purple-400 mb-2">💡 Cómo usar:</h4>
        <ol class="text-xs text-slate-400 space-y-1 list-decimal list-inside">
          <li>Agrega conceptos clave</li>
          <li>Haz clic en 2 conceptos</li>
          <li>Etiqueta su relación</li>
          <li>Arrastra para reorganizar</li>
        </ol>
      </div>

      <!-- Lista de conexiones -->
      {#if connections.length > 0}
        <div class="p-4 bg-canvas-900/50 rounded-2xl border border-slate-700">
          <h4 class="text-sm font-semibold text-slate-300 mb-3">Conexiones</h4>
          <div class="space-y-2 max-h-48 overflow-y-auto">
            {#each connections as conn (conn.id)}
              {@const fromNode = nodes.find(n => n.id === conn.from)}
              {@const toNode = nodes.find(n => n.id === conn.to)}
              <div class="p-2 bg-canvas-950 rounded-lg border border-slate-700 relative group">
                <p class="text-xs text-slate-300 pr-6">
                  <span class="text-purple-400">{fromNode?.text}</span>
                  →
                  <span class="text-focus-400">{toNode?.text}</span>
                </p>
                <p class="text-xs text-slate-500 italic mt-1">"{conn.label}"</p>

                {#if !hasSubmitted || allowMultipleAttempts}
                  <button
                    onclick={() => deleteConnection(conn.id)}
                    class="absolute top-1 right-1 text-slate-600 hover:text-red-400 opacity-0 group-hover:opacity-100 transition-opacity"
                  >
                    ✕
                  </button>
                {/if}
              </div>
            {/each}
          </div>
        </div>
      {/if}
    </div>

    <!-- Canvas del mapa conceptual -->
    <div class="lg:col-span-3">
      <div
        bind:this={canvasRef}
        class="relative w-full h-[600px] bg-canvas-900/30 rounded-2xl border-2 border-dashed border-slate-700 overflow-hidden"
      >
        <!-- SVG para conexiones -->
        <svg class="absolute inset-0 w-full h-full pointer-events-none">
          {#each connections as conn (conn.id)}
            {@const fromNode = nodes.find(n => n.id === conn.from)}
            {@const toNode = nodes.find(n => n.id === conn.to)}
            {#if fromNode && toNode}
              <g>
                <line
                  x1={fromNode.x}
                  y1={fromNode.y}
                  x2={toNode.x}
                  y2={toNode.y}
                  stroke="#6366f1"
                  stroke-width="2"
                  marker-end="url(#arrowhead)"
                />
                <text
                  x={(fromNode.x + toNode.x) / 2}
                  y={(fromNode.y + toNode.y) / 2}
                  fill="#a78bfa"
                  font-size="12"
                  text-anchor="middle"
                  class="pointer-events-none select-none"
                >
                  {conn.label}
                </text>
              </g>
            {/if}
          {/each}

          <!-- Flecha -->
          <defs>
            <marker
              id="arrowhead"
              markerWidth="10"
              markerHeight="10"
              refX="9"
              refY="3"
              orient="auto"
            >
              <polygon points="0 0, 10 3, 0 6" fill="#6366f1" />
            </marker>
          </defs>
        </svg>

        <!-- Nodos -->
        {#each nodes as node (node.id)}
          {@const isSelected = selectedNodes.includes(node.id)}
          <div
            id="node-{node.id}"
            style="left: {node.x}px; top: {node.y}px;"
            class="absolute transform -translate-x-1/2 -translate-y-1/2 cursor-move group"
            draggable={!hasSubmitted || allowMultipleAttempts}
            ondragstart={(e) => handleNodeDragStart(e, node.id)}
            ondrag={(e) => handleNodeDrag(e, node.id)}
            ondragend={handleNodeDragEnd}
            onclick={() => selectNode(node.id)}
          >
            <div class="
              px-4 py-3 rounded-xl border-2 transition-all duration-300 min-w-[120px] text-center
              {isSelected ? 'bg-purple-500 border-purple-400 scale-110 shadow-lg shadow-purple-500/50' : 'bg-canvas-800 border-slate-600 hover:border-purple-500'}
            ">
              <p class="text-sm font-semibold {isSelected ? 'text-white' : 'text-slate-200'}">
                {node.text}
              </p>

              {#if !provideConcepts && (!hasSubmitted || allowMultipleAttempts)}
                <button
                  onclick={(e) => { e.stopPropagation(); deleteNode(node.id); }}
                  class="absolute -top-2 -right-2 w-5 h-5 bg-red-500 text-white rounded-full text-xs
                    opacity-0 group-hover:opacity-100 transition-opacity hover:bg-red-600"
                >
                  ✕
                </button>
              {/if}
            </div>
          </div>
        {/each}

        <!-- Mensaje cuando está vacío -->
        {#if nodes.length === 0}
          <div class="absolute inset-0 flex items-center justify-center">
            <p class="text-slate-500 text-center">
              {provideConcepts ? 'Los conceptos aparecerán aquí' : 'Agrega conceptos para empezar'}
            </p>
          </div>
        {/if}
      </div>
    </div>
  </div>

  <!-- Botones de acción -->
  <div class="flex gap-3 mt-6">
    {#if !hasSubmitted}
      <button
        onclick={handleSubmit}
        disabled={!meetsRequirements}
        class="
          flex-1 px-6 py-3 rounded-xl font-semibold
          bg-gradient-to-r from-purple-500 to-pink-500
          text-white
          transition-all duration-300
          hover:shadow-lg hover:shadow-purple-500/50
          disabled:opacity-50 disabled:cursor-not-allowed
        "
      >
        Enviar Mapa Conceptual
      </button>
    {:else if allowMultipleAttempts && !allCorrect}
      <button
        onclick={handleTryAgain}
        class="
          flex-1 px-6 py-3 rounded-xl font-semibold
          bg-canvas-800 text-white
          border border-slate-700
          transition-all duration-300
          hover:bg-slate-700
        "
      >
        Continuar Editando
      </button>
    {/if}
  </div>

  <!-- Feedback final -->
  {#if hasSubmitted && showFeedback}
    <div class="mt-6 p-4 rounded-2xl {allCorrect ? 'bg-purple-500/10 border border-purple-500/50' : 'bg-yellow-500/10 border border-yellow-500/50'}">
      <div class="flex items-center gap-3">
        <span class="text-3xl">
          {allCorrect ? '🎨' : '💭'}
        </span>
        <div class="flex-1">
          <p class="{allCorrect ? 'text-purple-400' : 'text-yellow-400'} font-semibold mb-1">
            {allCorrect
              ? '¡Excelente mapa conceptual! Has capturado las ideas clave y sus relaciones'
              : 'Buen inicio. Considera agregar más conceptos o conexiones importantes'}
          </p>
          {#if requiredConcepts.length > 0}
            <p class="text-slate-400 text-sm mt-1">
              Conceptos requeridos encontrados: {requiredConcepts.filter(req =>
                nodes.some(n => n.text.toLowerCase().includes(req.toLowerCase()))
              ).length}/{requiredConcepts.length}
            </p>
          {/if}
          <div class="flex items-center gap-3 mt-2">
            <div class="flex-1 bg-canvas-900 rounded-full h-2 overflow-hidden">
              <div
                class="h-full bg-gradient-to-r from-purple-500 to-pink-500 transition-all duration-500"
                style="width: {score}%"
              ></div>
            </div>
            <span class="text-sm font-semibold text-white">
              {score.toFixed(0)}%
            </span>
          </div>
        </div>
      </div>
    </div>
  {/if}
</div>

<!-- Diálogo de conexión -->
{#if showConnectionDialog}
  {@const fromNode = nodes.find(n => n.id === selectedNodes[0])}
  {@const toNode = nodes.find(n => n.id === selectedNodes[1])}
  <div class="fixed inset-0 bg-black/50 flex items-center justify-center z-50" onclick={cancelConnection}>
    <div class="bg-canvas-900 p-6 rounded-2xl border border-slate-700 max-w-md w-full mx-4" onclick={(e) => e.stopPropagation()}>
      <h3 class="text-lg font-semibold text-white mb-4">Etiquetar Conexión</h3>

      <p class="text-sm text-slate-400 mb-4">
        <span class="text-purple-400">{fromNode?.text}</span>
        →
        <span class="text-focus-400">{toNode?.text}</span>
      </p>

      <input
        type="text"
        bind:value={connectionLabel}
        onkeydown={(e) => e.key === 'Enter' && createConnection()}
        placeholder="Ej: produce, requiere, causa..."
        class="w-full px-4 py-3 bg-canvas-950 border border-slate-700 rounded-lg text-white
          focus:outline-none focus:border-purple-500 mb-4"
        autofocus
      />

      <div class="flex gap-3">
        <button
          onclick={cancelConnection}
          class="flex-1 px-4 py-2 bg-canvas-800 text-white rounded-lg hover:bg-slate-700"
        >
          Cancelar
        </button>
        <button
          onclick={createConnection}
          disabled={!connectionLabel.trim()}
          class="flex-1 px-4 py-2 bg-purple-500 text-white rounded-lg hover:bg-purple-600
            disabled:opacity-50 disabled:cursor-not-allowed"
        >
          Crear
        </button>
      </div>
    </div>
  </div>
{/if}

<style>
  /* Estilos adicionales si es necesario */
</style>
