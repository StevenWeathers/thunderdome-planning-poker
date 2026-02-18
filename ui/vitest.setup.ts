/// <reference types="@vitest/browser/matchers" />

import { loadedLocales, loadedFormatters } from './src/i18n/i18n-util';
import { initFormatters } from './src/i18n/formatters';
import { setLocale } from './src/i18n/i18n-svelte';
import en from './src/i18n/en/index';
import type { Translation } from './src/i18n/i18n-types';

// Ensure we're testing in browser mode, not SSR mode
Object.defineProperty(globalThis, '__SSR__', {
  value: false,
  writable: true,
});

// Initialize i18n for tests
loadedLocales['en'] = en as Translation;
loadedFormatters['en'] = initFormatters('en');
setLocale('en');
