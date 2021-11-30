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
      cy.task('db:teardown:registeredUser')
      cy.task('db:seed:registeredUser').as('currentUser')
    })

    it('successfully loads for authenticated registered user', function () {
      cy.login(this.currentUser)
      cy.createUserBattle(this.currentUser).then(() => {
        cy.visit(`/battle/${this.currentBattle.id}`)

        cy.get('h2').should('contain', 'Test Battle')
      })
    })

    describe('Plans', () => {
      it('should display existing plans', function () {
        cy.login(this.currentUser)
        cy.createUserBattle(this.currentUser, {
          name: 'Test Battle',
          pointValuesAllowed: ['1', '2', '3', '5', '8', '13', '?'],
          autoFinishVoting: false,
          plans: [
            {
              name: 'Defeat Loki',
              type: 'Story'
            }
          ],
          pointAverageRounding: 'ceil'
        }).then(() => {
          cy.visit(`/battle/${this.currentBattle.id}`)

          cy.getByTestId('plan-name').should('contain', 'Defeat Loki')
        })
      })

      it('should allow adding', function () {
        cy.login(this.currentUser)
        cy.createUserBattle(this.currentUser).then(() => {
          cy.visit(`/battle/${this.currentBattle.id}`)

          cy.getByTestId('plans-add').click()

          cy.get('input[name=planName]').type('Test Plan')

          cy.getByTestId('plan-save').click()

          cy.getByTestId('plan-name').should('contain', 'Test Plan')
        })
      })

      it('should allow editing plans', function () {
        cy.login(this.currentUser)
        cy.createUserBattle(this.currentUser, {
          name: 'Test Battle',
          pointValuesAllowed: ['1', '2', '3', '5', '8', '13', '?'],
          autoFinishVoting: false,
          plans: [
            {
              name: 'Defeat Loki',
              type: 'Epic'
            }
          ],
          pointAverageRounding: 'ceil'
        }).then(() => {
          cy.visit(`/battle/${this.currentBattle.id}`)

          cy.getByTestId('plan-type').should('contain', 'Epic')

          cy.getByTestId('plan-edit').click()

          cy.get('[name=planType]').select('Story')

          cy.getByTestId('plan-save').click()

          cy.getByTestId('plan-type').should('contain', 'Story')
        })
      })

      it('should allow deleting plans', function () {
        cy.login(this.currentUser)
        cy.createUserBattle(this.currentUser, {
          name: 'Test Battle',
          pointValuesAllowed: ['1', '2', '3', '5', '8', '13', '?'],
          autoFinishVoting: false,
          plans: [
            {
              name: 'Defeat Loki',
              type: 'Story'
            }
          ],
          pointAverageRounding: 'ceil'
        }).then(() => {
          cy.visit(`/battle/${this.currentBattle.id}`)

          cy.getByTestId('plan-name').should('contain', 'Defeat Loki')

          cy.getByTestId('plan-delete').click()

          cy.getByTestId('plan-name').should('not.exist')
        })
      })

      it('should allow activating plans', function () {
        cy.login(this.currentUser)
        cy.createUserBattle(this.currentUser, {
          name: 'Test Battle',
          pointValuesAllowed: ['1', '2', '3', '5', '8', '13', '?'],
          autoFinishVoting: false,
          plans: [
            {
              name: 'Defeat Loki',
              type: 'Story'
            }
          ],
          pointAverageRounding: 'ceil'
        }).then(() => {
          cy.visit(`/battle/${this.currentBattle.id}`)

          cy.getByTestId('currentplan-name').should('contain', '[Voting not started]')

          cy.getByTestId('pointCard').invoke('attr', 'data-locked').should('contain', 'true')

          cy.getByTestId('plan-name').should('contain', 'Defeat Loki')

          cy.getByTestId('plan-activate').click()

          cy.getByTestId('currentplan-name').should('contain', 'Defeat Loki')

          cy.getByTestId('pointCard').invoke('attr', 'data-locked').should('contain', 'false')
        })
      })

      it('should allow skipping plan voting', function () {
        cy.login(this.currentUser)
        cy.createUserBattle(this.currentUser, {
          name: 'Test Battle',
          pointValuesAllowed: ['1', '2', '3', '5', '8', '13', '?'],
          autoFinishVoting: false,
          plans: [
            {
              name: 'Defeat Loki',
              type: 'Story'
            }
          ],
          pointAverageRounding: 'ceil'
        }).then(() => {
          cy.visit(`/battle/${this.currentBattle.id}`)

          cy.getByTestId('currentplan-name').should('contain', '[Voting not started]')

          cy.getByTestId('pointCard').invoke('attr', 'data-locked').should('contain', 'true')

          cy.getByTestId('plan-name').should('contain', 'Defeat Loki')

          cy.getByTestId('plan-activate').click()

          cy.getByTestId('currentplan-name').should('contain', 'Defeat Loki')

          cy.getByTestId('pointCard').invoke('attr', 'data-locked').should('contain', 'false')

          cy.getByTestId('voting-skip').click()

          cy.getByTestId('currentplan-name').should('contain', '[Voting not started]')

          cy.getByTestId('pointCard').invoke('attr', 'data-locked').should('contain', 'true')
        })
      })
    })

    describe('Delete Battle', () => {
      it('successfully deletes battle and navigates to my battles page', function () {
        cy.login(this.currentUser)
        cy.createUserBattle(this.currentUser).then(() => {
          cy.visit(`/battle/${this.currentBattle.id}`)

          cy.getByTestId('battle-delete').click()

          // should have delete confirmation button
          cy.getByTestId('confirm-confirm').click()

          // we should be redirected to landing
          cy.location('pathname').should('equal', '/battles')
        })
      })

      it('cancel does not delete battle', function () {
        cy.login(this.currentUser)
        cy.createUserBattle(this.currentUser).then(() => {
          const battleUrl = `/battle/${this.currentBattle.id}`
          cy.visit(battleUrl)

          cy.getByTestId('battle-delete').click()

          // should have confirmation cancel button
          cy.getByTestId('confirm-cancel').click()

          // we should remain on battle
          cy.get('h2').should('contain', 'Test Battle')
          cy.location('pathname').should('equal', battleUrl)
        })
      })
    })
  })
})