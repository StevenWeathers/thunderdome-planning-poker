/** @type {import('ts-jest').JestConfigWithTsJest} */
module.exports = {
  testEnvironment: 'jsdom',
  transform: {
    '^.+\\.svelte$': [
      'svelte-jester',
      {
        'preprocess': '../build/svelte.config.js'
      }
    ],
    '^.+\\.ts$': 'ts-jest'
  },
  moduleFileExtensions: [
    'js',
    'ts',
    'svelte'
  ],

  setupFilesAfterEnv: [
    '@testing-library/jest-dom/extend-expect'
  ],
  moduleNameMapper: {
    'typesafe-i18n/svelte': 'typesafe-i18n/svelte/index.cjs',
    'typesafe-i18n/formatters': 'typesafe-i18n/formatters/index.cjs',
    'typesafe-i18n/detectors': 'typesafe-i18n/detectors/index.cjs'
  }
}