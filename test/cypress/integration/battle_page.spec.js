describe('The Battles Page', () => {
  describe('Unauthenicated User', () => {
    it('redirects to login for unauthenticated user', function () {
      cy.visit('/battle/bbaf82ef-a2d3-4e9a-b824-5e56a03ac3aa')

      cy.url().should('include', '/register')
    })
  })

  describe('Guest User', () => {})

  describe('Registered User', () => {
    beforeEach(() => {
      // seed a user in the DB that we can control from our tests
      cy.createUser()
    })

    it('successfully loads for authenticated registered user', function () {
      cy.login(this.currentUser)
      cy.createUserBattle(this.currentUser).then(() => {
        cy.visit(`/battle/${this.currentBattle.id}`)

        cy.get('h2').should('contain', 'Test Battle')
      })

      // cleanup our user (for some reason can't access this context in after utility
      cy.logout(this.currentUser)
    })
  })
})