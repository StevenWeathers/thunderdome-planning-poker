import { test } from '../fixtures/user-sessions'
import { expect } from '@playwright/test'
import { PokerGamesPage } from '../fixtures/poker-games-page'

const pageTitle = 'My Battles'

test.describe('Poker Games page', () => {
    test.describe('Unauthenticated user', () => {
        test('redirects to login', async ({ page }) => {
            const battlesPage = new PokerGamesPage(page)
            await battlesPage.goto()

            const title = battlesPage.page.locator('[data-formtitle="login"]')
            await expect(title).toHaveText('Login')
        })
    })

    test.describe('Guest user', () => {
        test('should load page', async ({ guestPage }) => {
            const battlesPage = new PokerGamesPage(guestPage.page)
            await battlesPage.goto()
            const title = battlesPage.page.locator('h1')
            await expect(title).toHaveText(pageTitle)
        })

        test('should allow creating a game', async ({ guestPage }) => {
            const battleName = 'Test Game'
            const battlesPage = new PokerGamesPage(guestPage.page)
            await battlesPage.goto()

            await battlesPage.createBattle({ name: battleName })

            const battleTitle = battlesPage.page.locator('h2')
            await expect(battleTitle).toHaveText(battleName)
        })

        test('should display games', async ({ guestPage }) => {
            const battleName = 'Test Display Games'

            const battlesPage = new PokerGamesPage(guestPage.page)
            await battlesPage.goto()

            await battlesPage.createBattle({ name: battleName })

            const battleTitle = battlesPage.page.locator('h2')
            await expect(battleTitle).toHaveText(battleName)

            await battlesPage.goto()

            const title = await battlesPage.gameCardName.filter({
                hasText: battleName,
            })
            await expect(title).toBeVisible()
        })
    })

    test.describe('Registered user', () => {
        test('should load page', async ({ registeredPage }) => {
            const battlesPage = new PokerGamesPage(registeredPage.page)
            await battlesPage.goto()
            const title = battlesPage.page.locator('h1')
            await expect(title).toHaveText(pageTitle)
        })

        test('should allow creating a game', async ({ registeredPage }) => {
            const battleName = 'Test Game'
            const battlesPage = new PokerGamesPage(registeredPage.page)
            await battlesPage.goto()

            await battlesPage.createBattle({ name: battleName })

            const battleTitle = battlesPage.page.locator('h2')
            await expect(battleTitle).toHaveText(battleName)
        })

        test('should display games', async ({ registeredPage }) => {
            const battleName = 'Test Display Game'

            const battlesPage = new PokerGamesPage(registeredPage.page)
            await battlesPage.goto()

            await battlesPage.createBattle({ name: battleName })

            const battleTitle = battlesPage.page.locator('h2')
            await expect(battleTitle).toHaveText(battleName)

            await battlesPage.goto()

            const title = await battlesPage.gameCardName.filter({
                hasText: battleName,
            })
            await expect(title).toBeVisible()
        })
    })
})
