import { test } from '../fixtures/user-sessions'
import { expect } from '@playwright/test'
import { RetrosPage } from '../fixtures/retros-page'

const pageTitle = 'My Retros'

test.describe('Retros page', () => {
    test.describe('Unauthenticated user', () => {
        test('redirects to login', async ({ page }) => {
            const retrosPage = new RetrosPage(page)
            await retrosPage.goto()

            const title = retrosPage.page.locator('[data-formtitle="login"]')
            await expect(title).toHaveText('Login')
        })
    })

    test.describe('Guest user', () => {
        test('should load page', async ({ guestPage }) => {
            const retrosPage = new RetrosPage(guestPage.page)
            await retrosPage.goto()
            const title = retrosPage.page.locator('h1')
            await expect(title).toHaveText(pageTitle)
        })

        test('should allow creating a retro', async ({ guestPage }) => {
            const retroName = 'Test Retro'
            const retrosPage = new RetrosPage(guestPage.page)
            await retrosPage.goto()

            await retrosPage.createRetro({ name: retroName })

            const retroTitle = retrosPage.page.locator('h1')
            await expect(retroTitle).toHaveText(retroName)
        })

        test('should display retros', async ({ guestPage }) => {
            const retroName = 'Test Display Retro'

            const retrosPage = new RetrosPage(guestPage.page)
            await retrosPage.goto()

            await retrosPage.createRetro({ name: retroName })

            const retroTitle = retrosPage.page.locator('h1')
            await expect(retroTitle).toHaveText(retroName)

            await retrosPage.goto()

            const title = await retrosPage.retroCardName.filter({
                hasText: retroName,
            })
            await expect(title).toBeVisible()
        })
    })

    test.describe('Registered user', () => {
        test('should load page', async ({ registeredPage }) => {
            const retrosPage = new RetrosPage(registeredPage.page)
            await retrosPage.goto()
            const title = retrosPage.page.locator('h1')
            await expect(title).toHaveText(pageTitle)
        })

        test('should allow creating a retro', async ({ registeredPage }) => {
            const retroName = 'Test Retro'
            const retrosPage = new RetrosPage(registeredPage.page)
            await retrosPage.goto()

            await retrosPage.createRetro({ name: retroName })

            const retroTitle = retrosPage.page.locator('h1')
            await expect(retroTitle).toHaveText(retroName)
        })

        test('should display retros', async ({ registeredPage }) => {
            const retroName = 'Test Display Retro'

            const retrosPage = new RetrosPage(registeredPage.page)
            await retrosPage.goto()

            await retrosPage.createRetro({ name: retroName })

            const retroTitle = retrosPage.page.locator('h1')
            await expect(retroTitle).toHaveText(retroName)

            await retrosPage.goto()

            const title = await retrosPage.retroCardName.filter({
                hasText: retroName,
            })
            await expect(title).toBeVisible()
        })
    })
})
