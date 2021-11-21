describe('The Organizations Page', () => {
  beforeEach(() => {
    // seed a user in the DB that we can control from our tests
    cy.createUser()
  })

  it('redirects to login for unauthenticated user', function () {
    cy.visit('/organizations')

    cy.url().should('include', '/login')

    // cleanup our user (for some reason can't access this context in after utility
    cy.logout(this.currentUser)
  })

  it('successfully loads for authenticated registered user', function () {
    cy.login(this.currentUser)

    cy.visit('/organizations')

    cy.get('h2').should('contain', 'Organizations')
    cy.get('h2').should('contain', 'Teams')

    // cleanup our user (for some reason can't access this context in after utility
    cy.logout(this.currentUser)
  })

  describe('Create Organization', () => {
    it('should successfully submit and navigate to new organization page', function () {
      cy.login(this.currentUser)

      cy.visit('/organizations')

      cy.get('button').contains('Create Organization').click()

      cy.get('form[name="createOrganization"] [name="organizationName"]').type('Test Organization{enter}')

      cy.get('h2').should('contain', 'Departments')
      cy.get('h2').should('contain', 'Teams')
      cy.get('h2').should('contain', 'Users')

      // cleanup our user (for some reason can't access this context in after utility
      cy.logout(this.currentUser)
    })
  })
})