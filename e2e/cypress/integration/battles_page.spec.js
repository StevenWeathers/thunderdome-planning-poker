describe('The Battles Page', () => {
  describe('Unauthenicated User', () => {
    it('redirects to login for unauthenticated user', function () {
      cy.visit('/battles')

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

      cy.visit('/battles')

      cy.get('h1').should('contain', 'My Battles')
    })

    it('displays users battles', function () {
      cy.login(this.currentUser)
      cy.createUserBattle(this.currentUser)
      cy.visit('/battles')

      // we should be in battle
      cy.getByTestId('battle-name').should('contain', 'Test Battle')
    })

    describe('Create Battle Form', () => {
      it('submitts successfully', function () {
        cy.login(this.currentUser)
        cy.visit('/battles')

        // fill form and submit to go to battle!
        cy.get('form[name="createBattle"] [name="battleName"]').type('Test Battle{enter}')

        // we should be in battle
        cy.get('h2').should('contain', 'Test Battle')
      })
    })
  })
})