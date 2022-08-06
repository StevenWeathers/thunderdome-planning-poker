import { Locator, Page } from '@playwright/test'

export class RegisterPage {
    readonly page: Page
    readonly guestUserNameField: Locator
    readonly registeredUserNameField: Locator
    readonly registeredUserEmailField: Locator
    readonly registeredUserPassword1Field: Locator
    readonly registeredUserPassword2Field: Locator

    constructor(page: Page) {
        this.page = page
        this.guestUserNameField = page.locator('[name="yourName1"]')
        this.registeredUserNameField = page.locator('[name="yourName2"]')
        this.registeredUserEmailField = page.locator('[name="yourEmail"]')
        this.registeredUserPassword1Field = page.locator(
            '[name="yourPassword1"]',
        )
        this.registeredUserPassword2Field = page.locator(
            '[name="yourPassword2"]',
        )
    }

    async goto() {
        await this.page.goto('/register')
    }

    async createGuestUser(name) {
        await this.guestUserNameField.fill(name)
        await this.guestUserNameField.press('Enter')
    }

    async createRegisteredUser(name, email, password1, password2) {
        await this.registeredUserNameField.fill(name)
        await this.registeredUserEmailField.fill(email)
        await this.registeredUserPassword1Field.fill(password1)
        await this.registeredUserPassword2Field.fill(password2)
        await this.registeredUserPassword2Field.press('Enter')
    }
}
