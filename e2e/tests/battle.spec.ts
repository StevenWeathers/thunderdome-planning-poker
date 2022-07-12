import {test} from '@playwright/test';

test.describe('The Battle Page', () => {
    test.describe('Unauthenicated User', () => {
        // it('redirects to register for unauthenticated user', function () {
        //     cy.visit('/battle/bbaf82ef-a2d3-4e9a-b824-5e56a03ac3aa')
        //
        //     cy.url().should('include', '/register')
        // })
    })

    test.describe('Guest User', () => {
        // beforeEach(() => {
        //     cy.task('db:teardown:guestUser')
        //     cy.createGuestUser()
        // })
        //
        // it('successfully loads', function () {
        //     cy.createUserBattle(this.currentUser).then(() => {
        //         cy.visit(`/battle/${this.currentBattle.id}`)
        //
        //         cy.get('h2').should('contain', 'Test Battle')
        //     })
        // })
    })

    test.describe('Registered User', () => {
        // beforeEach(() => {
        //     cy.task('db:teardown:registeredUser')
        //     cy.task('db:seed:registeredUser').as('currentUser')
        // })
        //
        // it('successfully loads for authenticated registered user', function () {
        //     cy.login(this.currentUser)
        //     cy.createUserBattle(this.currentUser).then(() => {
        //         cy.visit(`/battle/${this.currentBattle.id}`)
        //
        //         cy.get('h2').should('contain', 'Test Battle')
        //     })
        // })

        test.describe('User', () => {
            // it('can become spectator (when autoFinishVoting is true)', function () {
            //     cy.login(this.currentUser)
            //     cy.createUserBattle(this.currentUser, {
            //         name: 'Test Battle',
            //         pointValuesAllowed: ['1', '2', '3', '5', '8', '13', '?'],
            //         autoFinishVoting: true,
            //         plans: [],
            //         pointAverageRounding: 'ceil'
            //     }).then(() => {
            //         cy.visit(`/battle/${this.currentBattle.id}`)
            //
            //         cy.getByTestId('user-togglespectator').click()
            //
            //         cy.getByTestId('user-togglespectator').should('contain', 'Become Participant')
            //     })
            // })
            //
            // it('cannot become spectator (when autoFinishVoting is false)', function () {
            //     cy.login(this.currentUser)
            //     cy.createUserBattle(this.currentUser).then(() => {
            //         cy.visit(`/battle/${this.currentBattle.id}`)
            //
            //         cy.getByTestId('user-togglespectator').should('not.exist')
            //     })
            // })
            //
            // it('can demote leader (when is a leader)', function () {
            //     cy.login(this.currentUser)
            //     cy.createUserBattle(this.currentUser, {
            //         name: 'Test Battle',
            //         pointValuesAllowed: ['1', '2', '3', '5', '8', '13', '?'],
            //         autoFinishVoting: false,
            //         plans: [
            //             {
            //                 name: 'Defeat Loki',
            //                 type: 'Story'
            //             }
            //         ],
            //         pointAverageRounding: 'ceil'
            //     }).then(() => {
            //         cy.visit(`/battle/${this.currentBattle.id}`)
            //
            //         // yes you can demote yourself even
            //         cy.getByTestId('user-demote').click()
            //
            //         cy.getByTestId('user-demote').should('not.exist')
            //         cy.getByTestId('battle-delete').should('not.exist')
            //         cy.getByTestId('plans-add').should('not.exist')
            //         cy.getByTestId('plan-edit').should('not.exist')
            //         cy.getByTestId('plan-delete').should('not.exist')
            //         cy.getByTestId('plan-activate').should('not.exist')
            //
            //         cy.getByTestId('battle-abandon').should('exist')
            //         cy.getByTestId('plan-view').should('exist')
            //     })
            // })
            //
            // it('can abandon battle', function () {
            //     cy.login(this.currentUser)
            //     cy.createUserBattle(this.currentUser).then(() => {
            //         cy.visit(`/battle/${this.currentBattle.id}`)
            //
            //         // can't abandon battle if a leader
            //         cy.getByTestId('user-demote').click()
            //
            //         cy.getByTestId('battle-abandon').click()
            //
            //         // we should be redirected to battles page
            //         cy.location('pathname').should('equal', '/battles')
            //
            //         cy.getByTestId('battle-name').should('not.exist')
            //     })
            // })
        })

        test.describe('Plans', () => {
            // it('should display existing plans', function () {
            //     cy.login(this.currentUser)
            //     cy.createUserBattle(this.currentUser, {
            //         name: 'Test Battle',
            //         pointValuesAllowed: ['1', '2', '3', '5', '8', '13', '?'],
            //         autoFinishVoting: false,
            //         plans: [
            //             {
            //                 name: 'Defeat Loki',
            //                 type: 'Story'
            //             }
            //         ],
            //         pointAverageRounding: 'ceil'
            //     }).then(() => {
            //         cy.visit(`/battle/${this.currentBattle.id}`)
            //
            //         cy.getByTestId('plan-name').should('contain', 'Defeat Loki')
            //     })
            // })
            //
            // it('should allow adding', function () {
            //     cy.login(this.currentUser)
            //     cy.createUserBattle(this.currentUser).then(() => {
            //         cy.visit(`/battle/${this.currentBattle.id}`)
            //
            //         cy.getByTestId('plans-add').click()
            //
            //         cy.get('input[name=planName]').type('Test Plan')
            //
            //         cy.getByTestId('plan-save').click()
            //
            //         cy.getByTestId('plan-name').should('contain', 'Test Plan')
            //     })
            // })
            //
            // it('should allow editing plans', function () {
            //     cy.login(this.currentUser)
            //     cy.createUserBattle(this.currentUser, {
            //         name: 'Test Battle',
            //         pointValuesAllowed: ['1', '2', '3', '5', '8', '13', '?'],
            //         autoFinishVoting: false,
            //         plans: [
            //             {
            //                 name: 'Defeat Loki',
            //                 type: 'Epic'
            //             }
            //         ],
            //         pointAverageRounding: 'ceil'
            //     }).then(() => {
            //         cy.visit(`/battle/${this.currentBattle.id}`)
            //
            //         cy.getByTestId('plan-type').should('contain', 'Epic')
            //
            //         cy.getByTestId('plan-edit').click()
            //
            //         cy.get('[name=planType]').select('Story')
            //
            //         cy.getByTestId('plan-save').click()
            //
            //         cy.getByTestId('plan-type').should('contain', 'Story')
            //     })
            // })
            //
            // it('should allow deleting plans', function () {
            //     cy.login(this.currentUser)
            //     cy.createUserBattle(this.currentUser, {
            //         name: 'Test Battle',
            //         pointValuesAllowed: ['1', '2', '3', '5', '8', '13', '?'],
            //         autoFinishVoting: false,
            //         plans: [
            //             {
            //                 name: 'Defeat Loki',
            //                 type: 'Story'
            //             }
            //         ],
            //         pointAverageRounding: 'ceil'
            //     }).then(() => {
            //         cy.visit(`/battle/${this.currentBattle.id}`)
            //
            //         cy.getByTestId('plan-name').should('contain', 'Defeat Loki')
            //
            //         cy.getByTestId('plan-delete').click()
            //
            //         cy.getByTestId('plan-name').should('not.exist')
            //     })
            // })
            //
            // it('should allow activating plans', function () {
            //     cy.login(this.currentUser)
            //     cy.createUserBattle(this.currentUser, {
            //         name: 'Test Battle',
            //         pointValuesAllowed: ['1', '2', '3', '5', '8', '13', '?'],
            //         autoFinishVoting: false,
            //         plans: [
            //             {
            //                 name: 'Defeat Loki',
            //                 type: 'Story'
            //             }
            //         ],
            //         pointAverageRounding: 'ceil'
            //     }).then(() => {
            //         cy.visit(`/battle/${this.currentBattle.id}`)
            //
            //         cy.getByTestId('currentplan-name').should('contain', '[Voting not started]')
            //
            //         cy.getByTestId('pointCard').invoke('attr', 'data-locked').should('contain', 'true')
            //
            //         cy.getByTestId('plan-name').should('contain', 'Defeat Loki')
            //
            //         cy.getByTestId('plan-activate').click()
            //
            //         cy.getByTestId('currentplan-name').should('contain', 'Defeat Loki')
            //
            //         cy.getByTestId('pointCard').invoke('attr', 'data-locked').should('contain', 'false')
            //     })
            // })
            //
            // it('should allow skipping plan voting', function () {
            //     cy.login(this.currentUser)
            //     cy.createUserBattle(this.currentUser, {
            //         name: 'Test Battle',
            //         pointValuesAllowed: ['1', '2', '3', '5', '8', '13', '?'],
            //         autoFinishVoting: false,
            //         plans: [
            //             {
            //                 name: 'Defeat Loki',
            //                 type: 'Story'
            //             }
            //         ],
            //         pointAverageRounding: 'ceil'
            //     }).then(() => {
            //         cy.visit(`/battle/${this.currentBattle.id}`)
            //
            //         cy.getByTestId('currentplan-name').should('contain', '[Voting not started]')
            //
            //         cy.getByTestId('pointCard').invoke('attr', 'data-locked').should('contain', 'true')
            //
            //         cy.getByTestId('plan-name').should('contain', 'Defeat Loki')
            //
            //         cy.getByTestId('plan-activate').click()
            //
            //         cy.getByTestId('currentplan-name').should('contain', 'Defeat Loki')
            //
            //         cy.getByTestId('pointCard').invoke('attr', 'data-locked').should('contain', 'false')
            //
            //         cy.getByTestId('voting-skip').click()
            //
            //         cy.getByTestId('currentplan-name').should('contain', '[Voting not started]')
            //
            //         cy.getByTestId('pointCard').invoke('attr', 'data-locked').should('contain', 'true')
            //     })
            // })
            //
            // it('should allow finishing plan voting', function () {
            //     cy.login(this.currentUser)
            //     cy.createUserBattle(this.currentUser, {
            //         name: 'Test Battle',
            //         pointValuesAllowed: ['1', '2', '3', '5', '8', '13', '?'],
            //         autoFinishVoting: false,
            //         plans: [
            //             {
            //                 name: 'Defeat Loki',
            //                 type: 'Story'
            //             }
            //         ],
            //         pointAverageRounding: 'ceil'
            //     }).then(() => {
            //         cy.visit(`/battle/${this.currentBattle.id}`)
            //
            //         cy.getByTestId('currentplan-name').should('contain', '[Voting not started]')
            //
            //         cy.getByTestId('pointCard').invoke('attr', 'data-locked').should('contain', 'true')
            //
            //         cy.getByTestId('plan-name').should('contain', 'Defeat Loki')
            //
            //         cy.getByTestId('plan-activate').click()
            //
            //         cy.getByTestId('currentplan-name').should('contain', 'Defeat Loki')
            //
            //         cy.getByTestId('pointCard').invoke('attr', 'data-locked').should('contain', 'false')
            //
            //         cy.getByTestId('voting-finish').click()
            //
            //         cy.getByTestId('voteresult-total').should('exist')
            //         cy.getByTestId('voteresult-average').should('exist')
            //         cy.getByTestId('voteresult-high').should('exist')
            //         cy.getByTestId('voteresult-highcount').should('exist')
            //
            //         cy.getByTestId('pointCard').invoke('attr', 'data-locked').should('contain', 'true')
            //
            //         cy.getByTestId('currentplan-name').should('contain', 'Defeat Loki')
            //     })
            // })
            //
            // it('should allow saving plan voting final points', function () {
            //     cy.login(this.currentUser)
            //     cy.createUserBattle(this.currentUser, {
            //         name: 'Test Battle',
            //         pointValuesAllowed: ['1', '2', '3', '5', '8', '13', '?'],
            //         autoFinishVoting: false,
            //         plans: [
            //             {
            //                 name: 'Defeat Loki',
            //                 type: 'Story'
            //             }
            //         ],
            //         pointAverageRounding: 'ceil'
            //     }).then(() => {
            //         cy.visit(`/battle/${this.currentBattle.id}`)
            //
            //         cy.getByTestId('currentplan-name').should('contain', '[Voting not started]')
            //
            //         cy.getByTestId('plans-unpointed').should('contain', 'Unpointed (1)')
            //         cy.getByTestId('plans-pointed').should('contain', 'Pointed (0)')
            //         cy.getByTestId('pointCard').invoke('attr', 'data-locked').should('contain', 'true')
            //         cy.getByTestId('plan-name').should('contain', 'Defeat Loki')
            //         cy.getByTestId('plan-points').should('not.exist')
            //
            //         cy.getByTestId('plan-activate').click()
            //
            //         cy.getByTestId('currentplan-name').should('contain', 'Defeat Loki')
            //
            //         cy.getByTestId('pointCard').invoke('attr', 'data-locked').should('contain', 'false')
            //
            //         cy.getByTestId('voting-finish').click()
            //
            //         cy.getByTestId('voteresult-total').should('exist')
            //
            //         cy.get('[name="planPoints"]').select('1')
            //
            //         cy.getByTestId('voting-save').click()
            //
            //         cy.getByTestId('currentplan-name').should('contain', '[Voting not started]')
            //
            //         cy.getByTestId('plan-name').should('not.exist')
            //
            //         cy.getByTestId('plans-unpointed').should('contain', 'Unpointed (0)')
            //         cy.getByTestId('plans-pointed').should('contain', 'Pointed (1)')
            //         cy.getByTestId('plans-pointed').click()
            //         cy.getByTestId('plan-name').should('contain', 'Defeat Loki')
            //         cy.getByTestId('plan-points').should('contain', '1')
            //     })
            // })
        })

        test.describe('Delete Battle', () => {
            // it('successfully deletes battle and navigates to my battles page', function () {
            //     cy.login(this.currentUser)
            //     cy.createUserBattle(this.currentUser).then(() => {
            //         cy.visit(`/battle/${this.currentBattle.id}`)
            //
            //         cy.getByTestId('battle-delete').click()
            //
            //         // should have delete confirmation button
            //         cy.getByTestId('confirm-confirm').click()
            //
            //         // we should be redirected to battles page
            //         cy.location('pathname').should('equal', '/battles')
            //     })
            // })
            //
            // it('cancel does not delete battle', function () {
            //     cy.login(this.currentUser)
            //     cy.createUserBattle(this.currentUser).then(() => {
            //         const battleUrl = `/battle/${this.currentBattle.id}`
            //         cy.visit(battleUrl)
            //
            //         cy.getByTestId('battle-delete').click()
            //
            //         // should have confirmation cancel button
            //         cy.getByTestId('confirm-cancel').click()
            //
            //         // we should remain on battle
            //         cy.get('h2').should('contain', 'Test Battle')
            //         cy.location('pathname').should('equal', battleUrl)
            //     })
            // })
        })
    })
})