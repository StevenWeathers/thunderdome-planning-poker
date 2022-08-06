import { test } from '../fixtures/user-sessions'
import { expect } from '@playwright/test'
import { BattlesPage } from '../fixtures/battles-page'

const pageTitle = 'My Battles'

test.describe('Battles page', () => {
    test.describe('Unauthenticated user', () => {
        test('redirects to login', async ({ page }) => {
            const battlesPage = new BattlesPage(page)
            await battlesPage.goto()

            const title = battlesPage.page.locator('[data-formtitle="login"]')
            await expect(title).toHaveText('Login')
        })
    })

    test.describe('Guest user', () => {
        test('should load page', async ({ guestPage }) => {
            const battlesPage = new BattlesPage(guestPage.page)
            await battlesPage.goto()
            const title = battlesPage.page.locator('h1')
            await expect(title).toHaveText(pageTitle)
        })

        test('should allow creating a battle', async ({ guestPage }) => {
            const battleName = 'Test Battle'
            const battlesPage = new BattlesPage(guestPage.page)
            await battlesPage.goto()

            await battlesPage.createBattle({ name: battleName })

            const battleTitle = battlesPage.page.locator('h2')
            await expect(battleTitle).toHaveText(battleName)
        })

        test('should display battles', async ({ guestPage }) => {
            const battleName = 'Test Display Battle'

            const battlesPage = new BattlesPage(guestPage.page)
            await battlesPage.goto()

            await battlesPage.createBattle({ name: battleName })

            const battleTitle = battlesPage.page.locator('h2')
            await expect(battleTitle).toHaveText(battleName)

            await battlesPage.goto()

            const title = await battlesPage.battleCardName.filter({
                hasText: battleName,
            })
            await expect(title).toBeVisible()
        })
    })

    test.describe('Registered user', () => {
        test('should load page', async ({ registeredPage }) => {
            const battlesPage = new BattlesPage(registeredPage.page)
            await battlesPage.goto()
            const title = battlesPage.page.locator('h1')
            await expect(title).toHaveText(pageTitle)
        })

        test('should allow creating a battle', async ({ registeredPage }) => {
            const battleName = 'Test Battle'
            const battlesPage = new BattlesPage(registeredPage.page)
            await battlesPage.goto()

            await battlesPage.createBattle({ name: battleName })

            const battleTitle = battlesPage.page.locator('h2')
            await expect(battleTitle).toHaveText(battleName)
        })

        test('should display battles', async ({ registeredPage }) => {
            const battleName = 'Test Display Battle'

            const battlesPage = new BattlesPage(registeredPage.page)
            await battlesPage.goto()

            await battlesPage.createBattle({ name: battleName })

            const battleTitle = battlesPage.page.locator('h2')
            await expect(battleTitle).toHaveText(battleName)

            await battlesPage.goto()

            const title = await battlesPage.battleCardName.filter({
                hasText: battleName,
            })
            await expect(title).toBeVisible()
        })
    })
})
