describe('The User Profile Page', () => {
  describe('Unauthenticated User', () => {
    it('redirects to register', function () {
      cy.visit('/profile')

      cy.location('pathname').should('equal', '/register')
    })
  })

  describe('Guest User', () => {})

  describe('Registered User', () => {
    beforeEach(() => {
      // seed a user in the DB that we can control from our tests
      cy.createUser()
    })

    it('successfully loads', function () {
      cy.login(this.currentUser)

      cy.visit('/profile')

      cy.get('h2').should('contain', 'Your Profile')

      // cleanup our user (for some reason can't access this context in after utility
      cy.logout(this.currentUser)
    })

    describe('Delete Account', function () {
      it('successfully deletes the user', function () {
        cy.login(this.currentUser)

        cy.visit('/profile')

        cy.get('button').contains('Delete Account').click()

        // should have delete confirmation button
        cy.getByTestId('confirm-confirm').click()

        // we should be redirected to landing
        cy.location('pathname').should('equal', '/')

        // our user cookie should not be present
        cy.getCookie('warrior').should('not.exist')

        // UI should reflect this user being logged out
        cy.getByTestId('userprofile-link').should('not.exist')
      })

      it('cancel does not delete the user and remains on profile page', function () {
        cy.login(this.currentUser)

        cy.visit('/profile')

        cy.get('button').contains('Delete Account').click()

        cy.getByTestId('confirm-cancel').click()

        // we should be redirected to landing
        cy.location('pathname').should('equal', '/profile')

        // our user cookie should not be present
        cy.getCookie('warrior').should('exist')

        // UI should reflect this user being logged out
        cy.getByTestId('userprofile-link').should('contain', this.currentUser.name)

        // cleanup our user (for some reason can't access this context in after utility
        cy.logout(this.currentUser)
      })
    })
  })
})