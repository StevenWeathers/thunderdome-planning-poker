import { expect, test } from '../../fixtures/user-sessions'
import { AlertsPage } from '../../fixtures/admin/alerts-page'

test.describe('The Admin Alerts Page', () => {
    test.describe('Unauthenticated user', () => {
        test('redirects to login', async ({ page }) => {
            const adminPage = new AlertsPage(page)

            await adminPage.goto()

            const title = adminPage.page.locator('[data-formtitle="login"]')
            await expect(title).toHaveText('Login')
        })
    })

    test.describe('Guest user', () => {
        test('redirects to landing', async ({ guestPage }) => {
            const adminPage = new AlertsPage(guestPage.page)

            await adminPage.goto()

            const title = adminPage.page.locator('h1')
            await expect(title).toHaveText(
                'Thunderdome is an Agile Planning Poker app with a fun theme',
            )
        })
    })

    test.describe('Non Admin Registered User', () => {
        test('redirects to landing', async ({ registeredPage }) => {
            const adminPage = new AlertsPage(registeredPage.page)

            await adminPage.goto()

            const title = adminPage.page.locator('h1')
            await expect(title).toHaveText(
                'Thunderdome is an Agile Planning Poker app with a fun theme',
            )
        })
    })

    test.describe('Admin User', () => {
        test('loads alerts page', async ({ adminPage }) => {
            const ap = new AlertsPage(adminPage.page)

            await ap.goto()

            const title = ap.page.locator('h1')
            await expect(title).toHaveText('Alerts')
        })
    })
})
