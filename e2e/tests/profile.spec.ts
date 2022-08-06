import { expect, test } from '../fixtures/user-sessions'
import { ProfilePage } from '../fixtures/profile-page'

test.describe('User Profile page', () => {
    test('Unauthenticated user redirects to login', async ({ page }) => {
        const profilePage = new ProfilePage(page)
        await profilePage.goto()

        const title = profilePage.page.locator('[data-formtitle="login"]')
        await expect(title).toHaveText('Login')
    })

    test('Guest user successfully loads', async ({ guestPage }) => {
        const profilePage = new ProfilePage(guestPage.page)
        await profilePage.goto()

        await expect(profilePage.page.locator('h1')).toHaveText('Your Profile')
        await expect(profilePage.page.locator('[name=yourName]')).toHaveValue(
            guestPage.user.name,
        )
        await expect(profilePage.page.locator('[name=yourEmail]')).toHaveValue(
            '',
        )

        await expect(
            profilePage.page.locator('[data-testid="user-verified"]'),
        ).not.toBeVisible()
        await expect(
            profilePage.page.locator('[data-testid="request-verify"]'),
        ).not.toBeVisible()
    })

    test('Guest user cannot create API keys', async ({ guestPage }) => {
        const profilePage = new ProfilePage(guestPage.page)
        await profilePage.goto()

        await expect(profilePage.page.locator('h2')).toHaveText('API Keys')
        await profilePage.page.locator('[data-testid="apikey-create"]').click()
        await profilePage.page
            .locator('[name=keyName]')
            .fill('Create API Key Test')
        await profilePage.page
            .locator('[name=createApiKey] [type=submit]')
            .click()

        await expect(
            profilePage.page.locator('[data-testid="notification-msg"]'),
        ).toContainText('Only verified registered users can create API keys.')
        await expect(profilePage.page.locator('[name=keyName]')).toBeVisible()
    })

    test('Registered user successfully loads', async ({ registeredPage }) => {
        const profilePage = new ProfilePage(registeredPage.page)
        await profilePage.goto()

        await expect(profilePage.page.locator('h1')).toHaveText('Your Profile')
        await expect(profilePage.page.locator('[name=yourName]')).toHaveValue(
            registeredPage.user.name,
        )
        await expect(profilePage.page.locator('[name=yourEmail]')).toHaveValue(
            registeredPage.user.email,
        )

        await expect(
            profilePage.page.locator('[data-testid="user-verified"]'),
        ).not.toBeVisible()
        await expect(
            profilePage.page.locator('[data-testid="request-verify"]'),
        ).toBeVisible()
    })

    test('Registered non verified user cannot create API keys', async ({
        registeredPage,
    }) => {
        const profilePage = new ProfilePage(registeredPage.page)
        await profilePage.goto()

        await expect(profilePage.page.locator('h2')).toHaveText('API Keys')
        await profilePage.page.locator('[data-testid="apikey-create"]').click()
        await profilePage.page
            .locator('[name=keyName]')
            .fill('Create API Key Test')
        await profilePage.page
            .locator('[name=createApiKey] [type=submit]')
            .click()

        await expect(
            profilePage.page.locator('[data-testid="notification-msg"]'),
        ).toContainText('Only verified registered users can create API keys.')
        await expect(profilePage.page.locator('[name=keyName]')).toBeVisible()
    })

    test('Verified user should have verified status next to email field label', async ({
        verifiedPage,
    }) => {
        const profilePage = new ProfilePage(verifiedPage.page)
        await profilePage.goto()

        await expect(profilePage.page.locator('[name=yourEmail]')).toHaveValue(
            verifiedPage.user.email,
        )

        await expect(
            profilePage.page.locator('[data-testid="user-verified"]'),
        ).toBeVisible()
        await expect(
            profilePage.page.locator('[data-testid="request-verify"]'),
        ).not.toBeVisible()
    })

    // it('Verified user should display existing API keys', function () {
    //     cy.login(this.currentUser)
    //
    //     cy.createUserApikey(this.currentUser).then(() => {
    //         cy.visit('/profile')
    //
    //         cy.get('h2').should('contain', 'API Keys')
    //
    //         cy.getByTestId('apikey-name').should('contain', this.currentAPIKey.name)
    //         cy.getByTestId('apikey-prefix').should('contain', this.currentAPIKey.prefix)
    //         cy.getByTestId('apikey-active').invoke('attr', 'data-active').should('eq', 'true')
    //     })
    // })

    test('Verified user can create API keys', async ({ verifiedPage }) => {
        const apiKeyName = 'Create API Key Test'
        const profilePage = new ProfilePage(verifiedPage.page)
        await profilePage.goto()

        await expect(profilePage.page.locator('h2')).toHaveText('API Keys')
        await profilePage.page.locator('[data-testid="apikey-create"]').click()
        await profilePage.page.locator('[name=keyName]').fill(apiKeyName)
        await profilePage.page
            .locator('[name=createApiKey] [type=submit]')
            .click()

        await expect(profilePage.page.locator('[id="apiKey"]')).toBeVisible()
        await profilePage.page.locator('[data-testid="apikey-close"]').click()

        await expect(
            profilePage.page.locator('[data-testid="apikey-name"]', {
                hasText: apiKeyName,
            }),
        ).toBeVisible()
    })

    // it('can toggle api key active status', function () {
    //     cy.login(this.currentUser)
    //
    //     cy.createUserApikey(this.currentUser).then(() => {
    //         cy.visit('/profile')
    //
    //         cy.get('h2').should('contain', 'API Keys')
    //
    //         cy.getByTestId('apikey-active').invoke('attr', 'data-active').should('eq', 'true')
    //
    //         cy.getByTestId('apikey-activetoggle').click()
    //
    //         cy.getByTestId('apikey-active').invoke('attr', 'data-active').should('eq', 'false')
    //     })
    // })
    //
    // it('can delete api key', function () {
    //     cy.login(this.currentUser)
    //
    //     cy.createUserApikey(this.currentUser).then(() => {
    //         cy.visit('/profile')
    //
    //         cy.get('h2').should('contain', 'API Keys')
    //
    //         cy.getByTestId('apikey-name').should('contain', this.currentAPIKey.name)
    //
    //         cy.getByTestId('apikey-delete').click()
    //
    //         cy.getByTestId('apikey-name').should('not.exist')
    //     })
    // })
    //
    // it('can create no more than 5 API keys (default for config)', function () {
    //     cy.login(this.currentUser)
    //
    //     cy.createUserApikey(this.currentUser, {name: 'testkey1'})
    //     cy.createUserApikey(this.currentUser, {name: 'testkey2'})
    //     cy.createUserApikey(this.currentUser, {name: 'testkey3'})
    //     cy.createUserApikey(this.currentUser, {name: 'testkey4'})
    //     cy.createUserApikey(this.currentUser, {name: 'testkey5'})
    //
    //     cy.visit('/profile')
    //
    //     cy.getByTestId('apikey-create').click()
    //     cy.get('[name=keyName]').type('Create API Key Test')
    //     cy.get('[name=createApiKey] [type=submit]').click()
    //
    //     cy.get('[name=keyName]').should('exist')
    //     cy.getByTestId('notification-msg').should('contain', 'You have the max number of API keys allowed.')
    // })

    test('delete account confirmation cancel does not delete account', async ({
        registeredPage,
    }) => {
        const profilePage = new ProfilePage(registeredPage.page)
        await profilePage.goto()

        await expect(profilePage.page.locator('h1')).toHaveText('Your Profile')

        await profilePage.page
            .locator('button', { hasText: 'Delete Account' })
            .click()
        await profilePage.page.locator('data-testid=confirm-cancel').click()

        await expect(profilePage.page.locator('h1')).toHaveText('Your Profile')
        await expect(
            profilePage.page.locator('data-testid=userprofile-link'),
        ).toHaveText(registeredPage.user.name)
    })

    // it('can delete account', function () {
    //     cy.login(this.currentUser)
    //
    //     cy.visit('/profile')
    //
    //     cy.get('button').contains('Delete Account').click()
    //
    //     // should have delete confirmation button
    //     cy.getByTestId('confirm-confirm').click()
    //
    //     // we should be redirected to landing
    //     cy.location('pathname').should('equal', '/')
    //
    //     // our user cookie should not be present
    //     cy.getCookie('warrior').should('not.exist')
    //
    //     // UI should reflect this user being logged out
    //     cy.getByTestId('userprofile-link').should('not.exist')
    // })
})
