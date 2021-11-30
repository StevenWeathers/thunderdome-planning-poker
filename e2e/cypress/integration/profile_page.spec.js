describe('The User Profile Page', () => {
  describe('Unauthenticated User', () => {
    it('redirects to register', function () {
      cy.visit('/profile')

      cy.location('pathname').should('equal', '/register')
    })
  })

  describe('Guest User', () => {
    beforeEach(() => {
      cy.task('db:teardown:guestUser')
      cy.createGuestUser()
    })

    it('successfully loads', function () {
      cy.visit('/profile')

      cy.get('h2').should('contain', 'Your Profile')

      cy.get('[name=yourName]').should('have.value', this.currentUser.name)
    })

    describe('API Keys', () => {
      it('can not create API keys', function () {
        cy.visit('/profile')

        cy.get('h2').should('contain', 'API Keys')

        cy.getByTestId('apikey-create').click()

        cy.get('[name=keyName]').type('Create API Key Test')
        cy.get('[name=createApiKey] [type=submit]').click()

        cy.get('[name=keyName]').should('exist')

        cy.getByTestId('notification-msg').should('contain', 'Only verified registered users can create API keys.')
      })
    })
  })

  describe('Registered User', () => {
    beforeEach(() => {
      cy.task('db:teardown:registeredUser')
      cy.task('db:seed:registeredUser').as('currentUser')
    })

    it('successfully loads', function () {
      cy.login(this.currentUser)

      cy.visit('/profile')

      cy.get('h2').should('contain', 'Your Profile')

      cy.get('[name=yourName]').should('have.value', this.currentUser.name)
    })

    describe('API Keys', () => {
      it('can not create API keys', function () {
        cy.login(this.currentUser)

        cy.visit('/profile')

        cy.get('h2').should('contain', 'API Keys')

        cy.getByTestId('apikey-create').click()

        cy.get('[name=keyName]').type('Create API Key Test')
        cy.get('[name=createApiKey] [type=submit]').click()

        cy.get('[name=keyName]').should('exist')

        cy.getByTestId('notification-msg').should('contain', 'Only verified registered users can create API keys.')
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
      })
    })
  })

  describe('Verified Registered User', () => {
    beforeEach(() => {
      cy.task('db:teardown:verifiedUser')
      cy.task('db:seed:verifiedUser').as('currentUser')
    })

    describe('API Keys', () => {
      it('displays users API keys', function () {
        cy.login(this.currentUser)

        cy.createUserApikey(this.currentUser).then(() => {
          cy.visit('/profile')

          cy.get('h2').should('contain', 'API Keys')

          cy.getByTestId('apikey-name').should('contain', this.currentAPIKey.name)
          cy.getByTestId('apikey-prefix').should('contain', this.currentAPIKey.prefix)
          cy.getByTestId('apikey-active').should('contain', 'true')
        })
      })

      it('can create API key', function () {
        cy.login(this.currentUser)

        cy.visit('/profile')

        cy.get('h2').should('contain', 'API Keys')

        cy.getByTestId('apikey-create').click()

        cy.get('[name=keyName]').type('Create API Key Test')
        cy.get('[name=createApiKey] [type=submit]').click()

        cy.get('[id="apiKey"]').should('exist')

        cy.getByTestId('apikey-close').click()

        cy.getByTestId('apikey-name').should('contain', 'Create API Key Test')
      })

      it('can toggle api key active status', function () {
        cy.login(this.currentUser)

        cy.createUserApikey(this.currentUser).then(() => {
          cy.visit('/profile')

          cy.get('h2').should('contain', 'API Keys')

          cy.getByTestId('apikey-active').should('contain', 'true')

          cy.getByTestId('apikey-activetoggle').click()

          cy.getByTestId('apikey-active').should('contain', 'false')
        })
      })

      it('can delete api key', function () {
        cy.login(this.currentUser)

        cy.createUserApikey(this.currentUser).then(() => {
          cy.visit('/profile')

          cy.get('h2').should('contain', 'API Keys')

          cy.getByTestId('apikey-name').should('contain', this.currentAPIKey.name)

          cy.getByTestId('apikey-delete').click()

          cy.getByTestId('apikey-name').should('not.exist')
        })
      })

      it('can create no more than 5 API keys (default for config)', function () {
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
      })
    })
  })
})