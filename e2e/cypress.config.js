const { defineConfig } = require('cypress')

module.exports = defineConfig({
  blockHosts: '*.google-analytics.com',
  db: {
    name: 'thunderdome',
    user: 'thor',
    pass: 'odinson',
    port: '5432',
    host: 'localhost',
  },
  viewportWidth: 1280,
  viewportHeight: 720,
  e2e: {
    // We've imported your old cypress plugins here.
    // You may want to clean this up later by importing these.
    setupNodeEvents(on, config) {
      return require('./cypress/plugins/index.js')(on, config)
    },
    baseUrl: 'http://localhost:8080',
  },
})
