import { Locator, Page } from '@playwright/test'

export class BattlePage {
    readonly page: Page
    readonly battleTitle: Locator
    readonly toggleSpectator: Locator
    readonly userDemoteBtn: Locator
    readonly battleDeleteBtn: Locator
    readonly battleDeleteConfirmBtn: Locator
    readonly battleDeleteCancelBtn: Locator
    readonly addPlansBtn: Locator
    readonly editPlanBtn: Locator
    readonly deletePlanBtn: Locator
    readonly activatePlanBtn: Locator
    readonly abandonBattleBtn: Locator
    readonly viewPlanBtn: Locator
    readonly planName: Locator
    readonly planType: Locator
    readonly planNameField: Locator
    readonly planTypeField: Locator
    readonly savePlanBtn: Locator

    constructor(page: Page) {
        this.battleTitle = page.locator('h2')
        this.toggleSpectator = page.locator(
            '[data-testid="user-togglespectator"]',
        )
        this.userDemoteBtn = page.locator(`[data-testid="user-demote"]`)
        this.battleDeleteBtn = page.locator('[data-testid="battle-delete"]')
        this.battleDeleteConfirmBtn = page.locator(
            'data-testid=confirm-confirm',
        )
        this.battleDeleteCancelBtn = page.locator('data-testid=confirm-cancel')
        this.addPlansBtn = page.locator('[data-testid="plans-add"]')
        this.editPlanBtn = page.locator('[data-testid="plan-edit"]')
        this.deletePlanBtn = page.locator('[data-testid="plan-delete"]')
        this.activatePlanBtn = page.locator('[data-testid="plan-activate"]')
        this.abandonBattleBtn = page.locator('[data-testid="battle-abandon"]')
        this.viewPlanBtn = page.locator('[data-testid="plan-view"]')
        this.planName = page.locator('data-testid=plan-name')
        this.planType = page.locator('data-testid=plan-type')
        this.planNameField = page.locator('input[name=planName]')
        this.planTypeField = page.locator('select[name=planType]')
        this.savePlanBtn = page.locator('data-testid=plan-save')

        this.page = page
    }

    async goto(id) {
        await this.page.goto(`/battle/${id}`)
    }

    async addPlan(name) {
        await this.addPlansBtn.click()
        await this.planNameField.fill(name)
        await this.savePlanBtn.click()
    }
}
