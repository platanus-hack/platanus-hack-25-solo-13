// @ts-check
import { defineConfig } from 'astro/config';

import tailwindcss from '@tailwindcss/vite';

// https://astro.build/config
export default defineConfig({
  server: {
    host: '0.0.0.0',
    port: 4321
  },
  vite: {
    plugins: [tailwindcss()],
    server: {
      host: '0.0.0.0',
      allowedHosts: [
        'lumera.cl',
        'www.lumera.cl',
        'lumera.lat',
        'www.lumera.lat',
        'localhost'
      ]
    }
  }
});