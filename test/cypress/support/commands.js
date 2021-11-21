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
const genericPassword = 'kentRules!'

function generateEmail (length) {
  const characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'
  const charactersLength = characters.length
  let result = ''

  for (let i = 0; i < length; i++) {
    result += characters.charAt(Math.floor(Math.random() * charactersLength))
  }

  return result
}

Cypress.Commands.add('createUser', () => {
  cy.request('POST', '/api/auth/register', {
    name: 'loki',
    email: `${generateEmail(6)}@thunderdome.dev`,
    password1: genericPassword,
    password2: genericPassword
  })
    .its('body.data')
    .as('currentUser')
})
Cypress.Commands.add('login', (user) => {
  const { id, name, email, rank } = user

  cy.request('POST', '/api/auth', {
    email,
    password: genericPassword,
  })
  cy.setCookie('warrior', `{%22id%22:%22${id}%22%2C%22name%22:%22${name}%22%2C%22email%22:%22${email}%22%2C%22rank%22:%22${rank}%22%2C%22locale%22:%22%22%2C%22notificationsEnabled%22:true}`)
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
