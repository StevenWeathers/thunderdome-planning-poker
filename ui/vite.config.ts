import { defineConfig } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';

// https://vite.dev/config/
export default defineConfig({
  resolve: {
    dedupe: ['svelte'],
  },
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
              /<link rel="modulepreload" crossorigin href="/g,
              '<link rel="modulepreload" nonce="{{.Nonce}}" href="{{.UIConfig.AppConfig.PathPrefix}}',
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
  optimizeDeps: {
    include: ['lucide-svelte'],
  },
  build: {
    sourcemap: true,
    chunkSizeWarningLimit: 1200,
    rollupOptions: {
      output: {
        manualChunks(id: string) {
          const path = id.replace(/\\/g, '/');
          if (path.includes('node_modules')) {
            // Keep lucide and lucide-svelte together to avoid circular init issues across chunks
            if (path.includes('/lucide-svelte/') || path.includes('/lucide/'))
              return 'vendor_lucide';
            if (path.includes('/svelte/')) return 'vendor_svelte';
            if (path.includes('/quill/')) return 'vendor_quill';
            return 'vendor';
          }
          if (path.includes('/src/pages/poker/')) return 'poker';
          if (path.includes('/src/pages/retro/')) return 'retro';
          if (path.includes('/src/pages/storyboard/')) return 'storyboard';
          if (path.includes('/src/pages/admin/')) return 'admin';
          if (
            path.includes('/src/pages/organization/') ||
            path.includes('/src/pages/department/') ||
            path.includes('/src/pages/team/')
          )
            return 'team';
          if (path.includes('/src/pages/user/')) return 'user';
          if (path.includes('/src/pages/subscription/')) return 'subscription';
          if (path.includes('/src/pages/support/')) return 'support';
          return undefined;
        },
      },
    },
  },
});
