import {expect, test} from '../fixtures/user-sessions';
import {BattlePage} from "../fixtures/battle-page";

test.describe('Battle page', () => {
    let battle = {
        id: '',
        name: ''
    }
    let battleWithAutoFinishVoting = {
        id: '',
        name: ''
    }

    test.beforeAll(async ({registeredPage, verifiedPage, adminPage}) => {
        const b = await registeredPage.createBattle({
            name: 'e2e battle page tests',
            pointValuesAllowed: [
                '0', '1,', '2', '3', '5', '8', '13', '?'
            ],
            pointAverageRounding: 'ceil',
            plans: [
                {
                    name: 'Defeat Loki',
                    type: 'Story'
                }
            ],
            autoFinishVoting: false,
            battleLeaders: [
                adminPage.user.email
            ]
        });
        battle = b

        const bWithAutoFinishVoting = await verifiedPage.createBattle({
            name: 'e2e battle page tests',
            pointValuesAllowed: [
                '0', '1,', '2', '3', '5', '8', '13', '?'
            ],
            pointAverageRounding: 'ceil',
            plans: [
                {
                    name: 'Defeat Scarlet Witch',
                    type: 'Story'
                }
            ],
            autoFinishVoting: true,
        });
        battleWithAutoFinishVoting = bWithAutoFinishVoting
    })

    test.describe('Unauthenticated user', () => {
        test('redirects to register', async ({page}) => {
            const bp = new BattlePage(page);
            await bp.goto(battle.id);

            const title = bp.page.locator('h1');
            await expect(title).toHaveText('Enlist to Battle');
        })
    })

    test.describe('Guest User', () => {
        test('successfully loads', async ({guestPage}) => {
            const bp = new BattlePage(guestPage.page);
            await bp.goto(battle.id);

            const title = bp.page.locator('h2');
            await expect(title).toHaveText(battle.name);
        })
    })

    test.describe('Registered User', () => {
        test('successfully loads', async ({registeredPage}) => {
            const bp = new BattlePage(registeredPage.page);
            await bp.goto(battle.id);

            const title = bp.page.locator('h2');
            await expect(title).toHaveText(battle.name);
        })

        test.describe('User', () => {
            test('can become spectator (when autoFinishVoting is true)', async ({registeredPage}) => {
                const bp = new BattlePage(registeredPage.page);
                await bp.goto(battleWithAutoFinishVoting.id);

                const spectatorButton = bp.page.locator('[data-testid="user-togglespectator"]');

                await spectatorButton.click();
                await expect(spectatorButton).toHaveText('Become Participant');
            })

            test('cannot become spectator (when autoFinishVoting is false)', async ({registeredPage}) => {
                const bp = new BattlePage(registeredPage.page);
                await bp.goto(battle.id);

                const spectatorButton = bp.page.locator('[data-testid="user-togglespectator"]');

                await expect(spectatorButton).not.toBeVisible()
            })

            test('can demote leader (when is a leader)', async ({adminPage}) => {
                const bp = new BattlePage(adminPage.page);
                await bp.goto(battle.id);

                const userDemoteBtn = bp.page.locator(`[data-testid="user-card"][data-userid="${adminPage.user.id}"] [data-testid="user-demote"]`);
                const battleDeleteBtn = bp.page.locator('[data-testid="battle-delete"]');
                const addPlansBtn = bp.page.locator('[data-testid="plans-add"]');
                const editPlanBtn = bp.page.locator('[data-testid="plan-edit"]');
                const deletePlanBtn = bp.page.locator('[data-testid="plan-delete"]');
                const activatePlanBtn = bp.page.locator('[data-testid="plan-activate"]');
                const abandonBattleBtn = bp.page.locator('[data-testid="battle-abandon"]');
                const viewPlanBtn = bp.page.locator('[data-testid="plan-view"]');

                await expect(userDemoteBtn).toBeVisible();
                await expect(battleDeleteBtn).toBeVisible();
                await expect(addPlansBtn).toBeVisible();
                await expect(editPlanBtn).toBeVisible();
                await expect(deletePlanBtn).toBeVisible();
                await expect(activatePlanBtn).toBeVisible();
                await expect(viewPlanBtn).toBeVisible();
                await expect(abandonBattleBtn).not.toBeVisible();

                // yes you can demote yourself!
                await userDemoteBtn.click();

                await expect(userDemoteBtn).not.toBeVisible();
                await expect(battleDeleteBtn).not.toBeVisible();
                await expect(addPlansBtn).not.toBeVisible();
                await expect(editPlanBtn).not.toBeVisible();
                await expect(deletePlanBtn).not.toBeVisible();
                await expect(activatePlanBtn).not.toBeVisible();
                await expect(viewPlanBtn).toBeVisible();
                await expect(abandonBattleBtn).toBeVisible();
            });

            test('can abandon battle', async ({registeredPage}) => {
                const bp = new BattlePage(registeredPage.page);
                await bp.goto(battleWithAutoFinishVoting.id);

                await bp.page.click('[data-testid="battle-abandon"]');
                await expect(bp.page.locator('h1')).toHaveText('My Battles');
            })
        })

        test.describe('Plans', () => {
            // test('should display existing plans', async ({} => {
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
            // test('should allow adding', async ({} => {
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
            // test('should allow editing plans', async ({} => {
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
            // test('should allow deleting plans', async ({} => {
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
            // test('should allow activating plans', async ({} => {
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
            // test('should allow skipping plan voting', async ({} => {
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
            // test('should allow finishing plan voting', async ({} => {
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
            // test('should allow saving plan voting final points', async ({} => {
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
            // test('successfully deletes battle and navigates to my battles page', async ({} => {
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
            // test('cancel does not delete battle', async ({} => {
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