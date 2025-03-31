import { defineConfig } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    svelte(),
    {
      name: 'alter-js-css-tags-for-nonce-and-path-prefix',
      apply: 'build',
      transformIndexHtml: {
        order: 'post',
        handler(html: string) {
          return html
            .replace(
              /<script type="module" crossorigin src="/g,
              '<script type="module" nonce="{{.Nonce}}" src="{{.UIConfig.AppConfig.PathPrefix}}',
            )
            .replace(
              /<link rel="stylesheet" crossorigin href="/g,
              '<link rel="stylesheet" nonce="{{.Nonce}}" href="{{.UIConfig.AppConfig.PathPrefix}}',
            );
        },
      },
    },
  ],
  css: {
    postcss: './postcss.config.js',
  },
  build: {
    sourcemap: true,
  },
});
