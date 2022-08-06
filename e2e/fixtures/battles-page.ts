import { Locator, Page } from '@playwright/test'

export class BattlesPage {
    readonly page: Page
    readonly battleNameFormField: Locator
    readonly battleCardName: Locator

    constructor(page: Page) {
        this.page = page
        this.battleNameFormField = page.locator(
            'form[name="createBattle"] [name="battleName"]',
        )
        this.battleCardName = page.locator('[data-testid="battle-name"]')
    }

    async goto() {
        await this.page.goto('/battles')
    }

    async createBattle({ name }) {
        await this.battleNameFormField.fill(name)
        await this.battleNameFormField.press('Enter')
    }
}
