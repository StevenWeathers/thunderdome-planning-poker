import { Locator, Page } from '@playwright/test'

export class LoginPage {
    readonly page: Page
    readonly emailField: Locator
    readonly passwordField: Locator

    constructor(page: Page) {
        this.emailField = page.locator('[name="yourEmail"]')
        this.passwordField = page.locator('[name="yourPassword"]')
        this.page = page
    }

    async goto() {
        await this.page.goto(`/login`)
    }

    async login(email, password) {
        await this.emailField.fill(email)
        await this.passwordField.fill(password)
        await this.passwordField.press('Enter')
    }
}
