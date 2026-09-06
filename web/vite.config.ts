import { defineConfig } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';

export default defineConfig({
  plugins: [svelte()],
  build: { manifest: true, sourcemap: false },
  server: { proxy: { '/health': 'http://127.0.0.1:8080', '/ready': 'http://127.0.0.1:8080' } },
});
