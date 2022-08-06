import { Locator, Page } from '@playwright/test'

export class RetroPage {
    readonly page: Page
    readonly retroTitle: Locator
    readonly retroDeleteBtn: Locator
    readonly retroDeleteConfirmBtn: Locator
    readonly retroDeleteCancelBtn: Locator

    constructor(page: Page) {
        this.retroTitle = page.locator('h1')
        this.retroDeleteBtn = page.locator('[data-testid="retro-delete"]')
        this.retroDeleteConfirmBtn = page.locator('data-testid=confirm-confirm')
        this.retroDeleteCancelBtn = page.locator('data-testid=confirm-cancel')

        this.page = page
    }

    async goto(id) {
        await this.page.goto(`/retro/${id}`)
    }
}
