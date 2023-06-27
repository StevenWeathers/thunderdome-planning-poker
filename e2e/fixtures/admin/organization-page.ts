import { Page } from '@playwright/test'

export class AdminOrganizationPage {
    readonly page: Page

    constructor(page: Page) {
        this.page = page
    }

    async goto(id) {
        await this.page.goto(`/admin/organizations/${id}`)
    }
}
