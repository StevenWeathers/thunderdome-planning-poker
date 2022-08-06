import { Locator, Page } from '@playwright/test'

export class StoryboardPage {
    readonly page: Page
    readonly storyboardTitle: Locator
    readonly storyboardDeleteBtn: Locator
    readonly storyboardDeleteConfirmBtn: Locator
    readonly storyboardDeleteCancelBtn: Locator

    constructor(page: Page) {
        this.storyboardTitle = page.locator('h1')
        this.storyboardDeleteBtn = page.locator(
            '[data-testid="storyboard-delete"]',
        )
        this.storyboardDeleteConfirmBtn = page.locator(
            'data-testid=confirm-confirm',
        )
        this.storyboardDeleteCancelBtn = page.locator(
            'data-testid=confirm-cancel',
        )

        this.page = page
    }

    async goto(id) {
        await this.page.goto(`/storyboard/${id}`)
    }
}
