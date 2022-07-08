describe('The Admin Teams Page', () => {
  describe('Unauthenticated User', () => {
    it('redirects to login', function () {
      cy.visit('/admin')

      cy.location('pathname').should('equal', '/login')
    })
  })

  describe('Guest User', () => {
    beforeEach(() => {
      cy.task('db:teardown:guestUser')
      cy.createGuestUser()
    })

    it('redirects to landing', function () {
      cy.visit('/admin/teams')

      cy.location('pathname').should('equal', '/')
    })
  })

  describe('Non Admin Registered User', () => {
    beforeEach(() => {
      cy.task('db:teardown:registeredUser')
      cy.task('db:seed:registeredUser').as('currentUser')
    })

    it('redirects to landing', function () {
      cy.login(this.currentUser)
      cy.visit('/admin/teams')

      cy.location('pathname').should('equal', '/')
    })
  })

  describe('Admin User', () => {
    beforeEach(() => {
      cy.task('db:teardown:adminUser')
      cy.task('db:seed:adminUser').as('currentUser')
    })

    it('successfully loads', function () {
      cy.login(this.currentUser)
      cy.visit('/admin/teams') // change URL to match your dev URL

      cy.get('h1').should('contain', 'Teams')
    })
  })
})