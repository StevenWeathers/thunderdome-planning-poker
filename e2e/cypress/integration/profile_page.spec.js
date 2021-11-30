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

    describe('API Keys', () => {
      it.skip('displays users API keys', function () {
        cy.login(this.currentUser)

        cy.createUserApikey(this.currentUser).then(() => {
          cy.visit('/profile')

          cy.get('h2').should('contain', 'API Keys')

          cy.getByTestId('apikey-name').should('contain', this.currentAPIKey.name)
          cy.getByTestId('apikey-prefix').should('contain', this.currentAPIKey.prefix)
          cy.getByTestId('apikey-active').should('contain', 'true')
        })

        // cleanup our user (for some reason can't access this context in after utility
        cy.logout(this.currentUser)
      })

      it.skip('can create API key', function () {
        cy.login(this.currentUser)

        cy.visit('/profile')

        cy.get('h2').should('contain', 'API Keys')

        cy.getByTestId('apikey-create').click()

        cy.get('[name=keyName]').type('Create API Key Test')
        cy.get('[name=createApiKey] [type=submit]').click()

        cy.get('[id="apiKey"]').should('exist')

        cy.getByTestId('apikey-close').click()

        cy.getByTestId('apikey-name').should('contain', 'Create API Key Test')

        // cleanup our user (for some reason can't access this context in after utility
        cy.logout(this.currentUser)
      })

      it.skip('can toggle api key active status', function () {
        cy.login(this.currentUser)

        cy.createUserApikey(this.currentUser).then(() => {
          cy.visit('/profile')

          cy.get('h2').should('contain', 'API Keys')

          cy.getByTestId('apikey-active').should('contain', 'true')

          cy.getByTestId('apikey-activetoggle').click()

          cy.getByTestId('apikey-active').should('contain', 'false')
        })

        // cleanup our user (for some reason can't access this context in after utility
        cy.logout(this.currentUser)
      })

      it.skip('can delete api key', function () {
        cy.login(this.currentUser)

        cy.createUserApikey(this.currentUser).then(() => {
          cy.visit('/profile')

          cy.get('h2').should('contain', 'API Keys')

          cy.getByTestId('apikey-name').should('contain', this.currentAPIKey.name)

          cy.getByTestId('apikey-delete').click()

          cy.getByTestId('apikey-name').should('not.exist')
        })

        // cleanup our user (for some reason can't access this context in after utility
        cy.logout(this.currentUser)
      })

      it.skip('can create no more than 5 API keys (default for config)', function () {
        cy.login(this.currentUser)

        cy.createUserApikey(this.currentUser, { name: 'testkey1' })
        cy.createUserApikey(this.currentUser, { name: 'testkey2' })
        cy.createUserApikey(this.currentUser, { name: 'testkey3' })
        cy.createUserApikey(this.currentUser, { name: 'testkey4' })
        cy.createUserApikey(this.currentUser, { name: 'testkey5' })

        cy.visit('/profile')

        cy.getByTestId('apikey-create').click()
        cy.get('[name=keyName]').type('Create API Key Test')
        cy.get('[name=createApiKey] [type=submit]').click()

        cy.get('[name=keyName]').should('exist')
        cy.getByTestId('notification-msg').should('contain', 'You have the max number of API keys allowed.')

        // cleanup our user (for some reason can't access this context in after utility
        cy.logout(this.currentUser)
      })
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