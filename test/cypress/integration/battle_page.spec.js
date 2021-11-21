describe('The Battle Page', () => {
  describe('Unauthenicated User', () => {
    it('redirects to register for unauthenticated user', function () {
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

    describe('Delete Battle', () => {
      it('successfully deletes battle and navigates to my battles page', function () {
        cy.login(this.currentUser)
        cy.createUserBattle(this.currentUser).then(() => {
          cy.visit(`/battle/${this.currentBattle.id}`)

          cy.get('button[data-testid="battle-delete"]').click()

          // should have delete confirmation button
          cy.get('[data-testid="confirm-actions"] button').contains('Delete Battle').click()

          // we should be redirected to landing
          cy.location('pathname').should('equal', '/battles')
        })

        // cleanup our user (for some reason can't access this context in after utility
        cy.logout(this.currentUser)
      })

      it('cancel does not delete battle', function () {
        cy.login(this.currentUser)
        cy.createUserBattle(this.currentUser).then(() => {
          const battleUrl = `/battle/${this.currentBattle.id}`
          cy.visit(battleUrl)

          cy.get('button[data-testid="battle-delete"]').click()

          // should have confirmation cancel button
          cy.get('[data-testid="confirm-actions"] button').contains('Cancel').click()

          // we should remain on battle
          cy.get('h2').should('contain', 'Test Battle')
          cy.location('pathname').should('equal', battleUrl)
        })

        // cleanup our user (for some reason can't access this context in after utility
        cy.logout(this.currentUser)
      })
    })
  })
})