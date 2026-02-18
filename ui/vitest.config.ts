import { defineConfig } from 'vitest/config';
import { svelte } from '@sveltejs/vite-plugin-svelte';
import { playwright } from '@vitest/browser-playwright';
import path from 'path';

export default defineConfig({
  plugins: [
    svelte({
      hot: !process.env.VITEST,
      compilerOptions: {
        dev: true,
      },
    }),
  ],
  server: {
    fs: {
      allow: [
        '.',
        path.resolve(__dirname, './node_modules/vitest-browser-svelte'),
        path.resolve(__dirname, './node_modules/@vitest/browser'),
      ],
    },
  },
  test: {
    globals: true,
    ui: false,
    browser: {
      provider: playwright({
        launchOptions: {
          headless: true,
          args: ['--headless=new', '--disable-gpu', '--no-sandbox'],
        },
      }),
      enabled: true,
      instances: [{ browser: 'chromium' }],
    },
    setupFiles: ['vitest-browser-svelte', './vitest.setup.ts'],
    include: ['src/**/*.test.ts', 'src/**/__tests__/*.test.ts'],
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json', 'html'],
      exclude: ['node_modules/', 'src/setupTests.ts', 'src/i18n/**'],
    },
  },
});
