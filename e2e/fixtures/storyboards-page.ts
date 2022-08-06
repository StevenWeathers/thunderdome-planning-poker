import { Locator, Page } from '@playwright/test'

export class StoryboardsPage {
    readonly page: Page
    readonly storyboardNameFormField: Locator
    readonly storyboardCardName: Locator

    constructor(page: Page) {
        this.page = page
        this.storyboardNameFormField = page.locator(
            'form[name="createStoryboard"] [name="storyboardName"]',
        )
        this.storyboardCardName = page.locator(
            '[data-testid="storyboard-name"]',
        )
    }

    async goto() {
        await this.page.goto('/storyboards')
    }

    async createStoryboard({ name }) {
        await this.storyboardNameFormField.fill(name)
        await this.storyboardNameFormField.press('Enter')
    }
}
