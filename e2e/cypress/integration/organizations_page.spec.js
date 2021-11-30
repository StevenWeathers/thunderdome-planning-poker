describe('The Organizations Page', () => {
  describe('Unauthenicated User', () => {
    it('redirects to login for unauthenticated user', function () {
      cy.visit('/organizations')

      cy.url().should('include', '/login')
    })
  })

  describe('Guest User', () => {})

  describe('Registered User', () => {
    beforeEach(() => {
      cy.task('db:teardown:registeredUser')
      cy.task('db:seed:registeredUser').as('currentUser')
    })

    it('successfully loads for authenticated registered user', function () {
      cy.login(this.currentUser)

      cy.visit('/organizations')

      cy.get('h2').should('contain', 'Organizations')
      cy.get('h2').should('contain', 'Teams')
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
      })
    })
  })
})