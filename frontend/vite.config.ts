import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [sveltekit()],
	server: {
		port: 5173,
		host: '0.0.0.0',
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
