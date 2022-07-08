/// <reference types="cypress" />
// ***********************************************************
// This example plugins/index.js can be used to load plugins
//
// You can change the location of this file or turn off loading
// the plugins file with the 'pluginsFile' configuration option.
//
// You can read more here:
// https://on.cypress.io/plugins-guide
// ***********************************************************

// This function is called when a project is opened or re-opened (e.g. due to
// the project's config changing)
const dbSetup = require('../db/setup.js')
const guestUser = require('../db/guestUser.js')
const registeredUser = require('../db/registeredUser.js')
const verifiedUser = require('../db/verifiedUser.js')
const adminUser = require('../db/adminUser.js')
const { db } = require('../../cypress.config.js')
/**
 * @type {Cypress.PluginConfig}
 */
// eslint-disable-next-line no-unused-vars
module.exports = (on, config) => {
  const pool = dbSetup(db)

  on('task', {
    'db:teardown:guestUser': () => {
      return guestUser.teardown(pool)
    },
    'db:teardown:registeredUser': () => {
      return registeredUser.teardown(pool)
    },
    'db:seed:registeredUser': () => {
      const user = registeredUser.seed(pool)

      return user
    },
    'db:teardown:verifiedUser': () => {
      return verifiedUser.teardown(pool)
    },
    'db:seed:verifiedUser': () => {
      const user = verifiedUser.seed(pool)

      return user
    },
    'db:teardown:adminUser': () => {
      return adminUser.teardown(pool)
    },
    'db:seed:adminUser': () => {
      const user = adminUser.seed(pool)

      return user
    },
  })
}
