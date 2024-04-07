import postcss from './postcss.config.js';
import { defineConfig } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';

// https://vitejs.dev/config/
export default defineConfig({
  build: {
    assetsDir: 'static',
    // minify: false,
  },
  plugins: [svelte()],
  css: {
    postcss,
  },
  // base: 'STATIC_DIR_BASE', // to be replaced by backend before rendering template
});
