import {expect, test} from '../fixtures/user-sessions';
import {BattlePage} from "../fixtures/battle-page";

const allowedPointValues = [
    '0', '1,', '2', '3', '5', '8', '13', '?'
];
const pointAverageRounding = 'ceil';
const lokiPlan = {name: 'Defeat Loki', type: 'Story'};
const thanosPlan = {name: 'Defeat Thanos', type: 'Epic'};
const scarletPlan = {name: 'Defeat Scarlet Witch', type: 'Epic'};

test.describe('Battle page', () => {
    let battle = {id: '', name: 'e2e battle page tests'};
    let battleWithPlan = {id: ''};
    let battleAddPlan = {id: ''};
    let battleEditPlan = {id: ''};
    let battleDeletePlan = {id: ''};
    let battleWithoutAutoVoting = {id: ''};
    let battleWithAutoVoting = {id: '', name: ''};
    let battleAbandon = {id: ''};
    let battleCancelDelete = {id: ''}
    let battleDelete = {id: ''};

    test.beforeAll(async ({registeredPage, verifiedPage, adminPage}) => {
        const commonBattle = {
            name: `${battle.name}`,
            pointValuesAllowed: [...allowedPointValues],
            pointAverageRounding: `${pointAverageRounding}`,
            plans: [],
            autoFinishVoting: false,
            battleLeaders: [
                `${adminPage.user.email}`
            ]
        };
        battle = await registeredPage.createBattle({...commonBattle, plans: [lokiPlan]});
        battleWithPlan = await registeredPage.createBattle({...commonBattle, plans: [lokiPlan]});
        battleAddPlan = await registeredPage.createBattle({...commonBattle});
        battleEditPlan = await registeredPage.createBattle({...commonBattle, plans: [thanosPlan]});
        battleDeletePlan = await registeredPage.createBattle({...commonBattle, plans: [scarletPlan]});
        battleWithoutAutoVoting = await registeredPage.createBattle({...commonBattle});
        battleWithAutoVoting = await verifiedPage.createBattle({
            ...commonBattle,
            autoFinishVoting: true,
        });
        battleAbandon = await verifiedPage.createBattle({...commonBattle});
        battleCancelDelete = await registeredPage.createBattle({...commonBattle});
        battleDelete = await registeredPage.createBattle({...commonBattle});
    })

    test('unauthenticated user redirects to register', async ({page}) => {
        const bp = new BattlePage(page);
        await bp.goto(battle.id);

        const title = bp.page.locator('h1');
        await expect(title).toHaveText('Enlist to Battle');
    })

    test('guest user successfully loads', async ({guestPage}) => {
        const bp = new BattlePage(guestPage.page);
        await bp.goto(battle.id);

        await expect(bp.battleTitle).toHaveText(battle.name);
    })

    test('registered user successfully loads', async ({registeredPage}) => {
        const bp = new BattlePage(registeredPage.page);
        await bp.goto(battle.id);

        await expect(bp.battleTitle).toHaveText(battle.name);
    })

    test('user cannot become spectator when autoFinishVoting is false', async ({registeredPage}) => {
        const bp = new BattlePage(registeredPage.page);
        await bp.goto(battleWithoutAutoVoting.id);

        await expect(bp.toggleSpectator).not.toBeVisible();
    });

    test('user can become spectator when autoFinishVoting is true', async ({registeredPage}) => {
        const bp = new BattlePage(registeredPage.page);
        await bp.goto(battleWithAutoVoting.id);

        const spectatorButton = bp.toggleSpectator;

        await spectatorButton.click();
        await expect(spectatorButton).toHaveText('Become Participant');
    });

    test('leader can demote user leader status', async ({adminPage}) => {
        const bp = new BattlePage(adminPage.page);
        await bp.goto(battle.id);

        const userDemoteBtn = bp.page.locator(
            `[data-testid="user-card"][data-userid="${adminPage.user.id}"]`
        ).locator('[data-testid="user-demote"]');

        await expect(userDemoteBtn).toBeVisible();
        await expect(bp.battleDeleteBtn).toBeVisible();
        await expect(bp.addPlansBtn).toBeVisible();
        await expect(bp.editPlanBtn).toBeVisible();
        await expect(bp.deletePlanBtn).toBeVisible();
        await expect(bp.activatePlanBtn).toBeVisible();
        await expect(bp.viewPlanBtn).toBeVisible();
        await expect(bp.abandonBattleBtn).not.toBeVisible();

        // yes you can demote yourself!
        await userDemoteBtn.click();

        await expect(userDemoteBtn).not.toBeVisible();
        await expect(bp.battleDeleteBtn).not.toBeVisible();
        await expect(bp.addPlansBtn).not.toBeVisible();
        await expect(bp.editPlanBtn).not.toBeVisible();
        await expect(bp.deletePlanBtn).not.toBeVisible();
        await expect(bp.activatePlanBtn).not.toBeVisible();
        await expect(bp.viewPlanBtn).toBeVisible();
        await expect(bp.abandonBattleBtn).toBeVisible();
    });

    test('user can abandon battle', async ({registeredPage}) => {
        const bp = new BattlePage(registeredPage.page);
        await bp.goto(battleAbandon.id);

        await bp.page.click('[data-testid="battle-abandon"]');
        await expect(bp.page.locator('h1')).toHaveText('My Battles');
    });

    test('should display existing plans', async ({registeredPage}) => {
        const bp = new BattlePage(registeredPage.page);
        await bp.goto(battleWithPlan.id);

        await expect(bp.planName.filter({hasText: lokiPlan.name})).toBeVisible();
    });

    test('should allow adding plans', async ({registeredPage}) => {
        const newPlanName = 'Defeat Thanos'
        const bp = new BattlePage(registeredPage.page);
        await bp.goto(battleAddPlan.id);

        await bp.addPlan(newPlanName);
        await expect(bp.planName.filter({hasText: newPlanName})).toBeVisible();
    });

    test('should allow editing plans', async ({registeredPage}) => {
        const newType = 'Story';
        const bp = new BattlePage(registeredPage.page);
        await bp.goto(battleEditPlan.id);

        await expect(bp.planType.filter({hasText: thanosPlan.type})).toBeVisible();
        await bp.editPlanBtn.click()
        await bp.planTypeField.selectOption(newType)
        await bp.savePlanBtn.click();
        await expect(bp.planType.filter({hasText: newType})).toBeVisible();
    });

    test('should allow deleting plans', async ({registeredPage}) => {
        const bp = new BattlePage(registeredPage.page);
        await bp.goto(battleDeletePlan.id);

        await expect(bp.planName).toBeVisible();
        await bp.deletePlanBtn.click();
        await expect(bp.planName).not.toBeVisible();
    });

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

    test('delete battle confirmation cancel does not delete battle', async ({registeredPage}) => {
        const bp = new BattlePage(registeredPage.page);
        await bp.goto(battleCancelDelete.id);

        await bp.battleDeleteBtn.click();
        await bp.battleDeleteCancelBtn.click();

        await expect(bp.battleTitle).toHaveText(battle.name);
    });

    test('delete battle confirmation confirm deletes battle and redirects to battles page', async ({registeredPage}) => {
        const bp = new BattlePage(registeredPage.page);
        await bp.goto(battleDelete.id);

        await bp.battleDeleteBtn.click();
        await bp.battleDeleteConfirmBtn.click();

        await expect(bp.page.locator('h1')).toHaveText('My Battles');
    });
});