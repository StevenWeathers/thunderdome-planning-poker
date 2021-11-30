// ***********************************************
// This example commands.js shows you how to
// create various custom commands and overwrite
// existing commands.
//
// For more comprehensive examples of custom
// commands please read more here:
// https://on.cypress.io/custom-commands
// ***********************************************
//
//
// -- This is a parent command --
// Cypress.Commands.add('login', (email, password) => { ... })
Cypress.Commands.add('login', (user) => {
  const { id, name, email, rank, password } = user

  cy.request('POST', '/api/auth', {
    email,
    password,
  })
  cy.setCookie('warrior', `{%22id%22:%22${id}%22%2C%22name%22:%22${name}%22%2C%22email%22:%22${email}%22%2C%22rank%22:%22${rank}%22%2C%22locale%22:%22%22%2C%22notificationsEnabled%22:true}`)
})

Cypress.Commands.add('guestLogin', (user) => {
  const { id, name, email, rank } = user

  cy.setCookie('warrior', `{%22id%22:%22${id}%22%2C%22name%22:%22${name}%22%2C%22email%22:%22${email}%22%2C%22rank%22:%22${rank}%22%2C%22locale%22:%22%22%2C%22notificationsEnabled%22:true}`)
})

Cypress.Commands.add('createGuestUser', () => {
  cy.request('POST', '/api/auth/guest', {
    name: 'Guest Test User'
  })
    .its('body.data')
    .as('currentUser')
    .then((user) => {
      const { id, name, rank } = user

      cy.setCookie('warrior', `{%22id%22:%22${id}%22%2C%22name%22:%22${name}%22%2C%22email%22:%22null%22%2C%22rank%22:%22${rank}%22%2C%22locale%22:%22%22%2C%22notificationsEnabled%22:true}`)
    })
})

Cypress.Commands.add('logout', (user) => {
  cy.request('DELETE', `/api/users/${user.id}`)
  cy.clearCookie('warrior')
})

Cypress.Commands.add('createUserBattle', (user, battle = {
  name: 'Test Battle',
  pointValuesAllowed: ['1', '2', '3', '5', '8', '13', '?'],
  autoFinishVoting: false,
  plans: [],
  pointAverageRounding: 'ceil'
}) => {
  cy.request('POST', `/api/users/${user.id}/battles`, battle)
    .its('body.data')
    .as('currentBattle')
})

Cypress.Commands.add('createUserTeam', (user, team = {
  name: 'Test Team',
}) => {
  cy.request('POST', `/api/users/${user.id}/teams`, team)
    .its('body.data')
    .as('currentTeam')
})

Cypress.Commands.add('createUserOrganization', (user, org = {
  name: 'Test Organization',
}) => {
  cy.request('POST', `/api/users/${user.id}/organizations`, org)
    .its('body.data')
    .as('currentOrganization')
})

Cypress.Commands.add('createUserDepartment', (user, organizationId, department = {
  name: 'Test Department',
}) => {
  cy.request('POST', `/api/organizations/${organizationId}/departments`, department)
    .its('body.data')
    .as('currentDepartment')
})

Cypress.Commands.add('createUserApikey', (user, apikey = {
  name: 'Test API Key',
}) => {
  cy.request('POST', `/api/users/${user.id}/apikeys`, apikey)
    .its('body.data')
    .as('currentAPIKey')
})

//
//
// -- This is a child command --
// Cypress.Commands.add('drag', { prevSubject: 'element'}, (subject, options) => { ... })
//
//
// -- This is a dual command --
// Cypress.Commands.add('dismiss', { prevSubject: 'optional'}, (subject, options) => { ... })
//
//
// -- This will overwrite an existing command --
// Cypress.Commands.overwrite('visit', (originalFn, url, options) => { ... })

Cypress.Commands.add('getByTestId', (selector, ...args) => {
  return cy.get(`[data-testid=${selector}]`, ...args)
})