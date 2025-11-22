export default defineConfig({
	plugins: [sveltekit()],
	server: {
		port: 5173,
		host: '0.0.0.0',
		allowedHosts: [           // ← AGREGAR ESTO
			'lumera.cl',
			'www.lumera.cl',
			'app.lumera.cl',
			'lumera.lat',
			'www.lumera.lat',
			'app.lumera.lat',
			'localhost'
		],                         // ← HASTA AQUÍ
		watch: {
			usePolling: true,
			interval: 1000
		},
		// ... resto del config
	}
});