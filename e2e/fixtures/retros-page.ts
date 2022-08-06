import { Locator, Page } from '@playwright/test'

export class RetrosPage {
    readonly page: Page
    readonly retroNameFormField: Locator
    readonly retroCardName: Locator

    constructor(page: Page) {
        this.page = page
        this.retroNameFormField = page.locator(
            'form[name="createRetro"] [name="retroName"]',
        )
        this.retroCardName = page.locator('[data-testid="retro-name"]')
    }

    async goto() {
        await this.page.goto('/retros')
    }

    async createRetro({ name }) {
        await this.retroNameFormField.fill(name)
        await this.retroNameFormField.press('Enter')
    }
}
