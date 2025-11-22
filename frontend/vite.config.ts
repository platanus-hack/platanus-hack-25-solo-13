import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
plugins: [sveltekit()],
server: {
	port: 5173,
	host: '0.0.0.0',
	allowedHosts: [
		'lumera.cl',
		'www.lumera.cl',
		'app.lumera.cl',
		'lumera.lat',
		'www.lumera.lat',
		'app.lumera.lat',
		'localhost'
	],
	watch: {
		usePolling: true,
		interval: 1000
	},
	hmr: {
		host: 'localhost',
		port: 5173
	},
	proxy: {
		'/api': {
			target: 'http://backend:8080',
			changeOrigin: true
		}
	}
}
});