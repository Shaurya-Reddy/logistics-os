import { defineConfig } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';

export default defineConfig({
  plugins: [svelte()],
  build: {
    manifest: true,
    sourcemap: false,
    target: 'es2022'
  },
  server: {
    proxy: {
      '/health': 'http://localhost:8080',
      '/ready': 'http://localhost:8080'
    }
  }
});
