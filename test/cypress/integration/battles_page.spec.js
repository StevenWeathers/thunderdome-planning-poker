describe('The Battles Page', () => {
  beforeEach(() => {
    // seed a user in the DB that we can control from our tests
    cy.createUser()
  })

  it('redirects to login for unauthenticated user', function () {
    cy.visit('/battles')

    cy.url().should('include', '/login')

    // cleanup our user (for some reason can't access this context in after utility
    cy.logout(this.currentUser)
  })

  it('successfully loads for authenticated registered user', function () {
    cy.login(this.currentUser)

    cy.visit('/battles')

    cy.get('h1').should('contain', 'My Battles')

    // cleanup our user (for some reason can't access this context in after utility
    cy.logout(this.currentUser)
  })

  it('displays users battles', function () {
    cy.login(this.currentUser)
    cy.createUserBattle(this.currentUser)
    cy.visit('/battles')

    // we should be in battle
    cy.get('[data-testid="battle-name"]').should('contain', 'Test Battle')

    // cleanup our user (for some reason can't access this context in after utility
    cy.logout(this.currentUser)
  })

  describe('Create Battle Form', () => {
    it('submitts successfully', function () {
      cy.login(this.currentUser)
      cy.visit('/battles')

      // fill form and submit to go to battle!
      cy.get('form[name="createBattle"] [name="battleName"]').type('Test Battle{enter}')

      // we should be in battle
      cy.get('h2').should('contain', 'Test Battle')

      // cleanup our user (for some reason can't access this context in after utility
      cy.logout(this.currentUser)
    })
  })
})