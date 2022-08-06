import { test } from '../fixtures/user-sessions'
import { expect } from '@playwright/test'
import { StoryboardsPage } from '../fixtures/storyboards-page'

const pageTitle = 'My Storyboards'

test.describe('Storyboards page', () => {
    test.describe('Unauthenticated user', () => {
        test('redirects to login', async ({ page }) => {
            const storyboardsPage = new StoryboardsPage(page)
            await storyboardsPage.goto()

            const title = storyboardsPage.page.locator(
                '[data-formtitle="login"]',
            )
            await expect(title).toHaveText('Login')
        })
    })

    test.describe('Guest user', () => {
        test('should load page', async ({ guestPage }) => {
            const storyboardsPage = new StoryboardsPage(guestPage.page)
            await storyboardsPage.goto()
            const title = storyboardsPage.page.locator('h1')
            await expect(title).toHaveText(pageTitle)
        })

        test('should allow creating a storyboard', async ({ guestPage }) => {
            const storyboardName = 'Test Storyboard'
            const storyboardsPage = new StoryboardsPage(guestPage.page)
            await storyboardsPage.goto()

            await storyboardsPage.createStoryboard({ name: storyboardName })

            const storyboardTitle = storyboardsPage.page.locator('h1')
            await expect(storyboardTitle).toHaveText(storyboardName)
        })

        test('should display storyboards', async ({ guestPage }) => {
            const storyboardName = 'Test Display Storyboard'

            const storyboardsPage = new StoryboardsPage(guestPage.page)
            await storyboardsPage.goto()

            await storyboardsPage.createStoryboard({ name: storyboardName })

            const storyboardTitle = storyboardsPage.page.locator('h1')
            await expect(storyboardTitle).toHaveText(storyboardName)

            await storyboardsPage.goto()

            const title = await storyboardsPage.storyboardCardName.filter({
                hasText: storyboardName,
            })
            await expect(title).toBeVisible()
        })
    })

    test.describe('Registered user', () => {
        test('should load page', async ({ registeredPage }) => {
            const storyboardsPage = new StoryboardsPage(registeredPage.page)
            await storyboardsPage.goto()
            const title = storyboardsPage.page.locator('h1')
            await expect(title).toHaveText(pageTitle)
        })

        test('should allow creating a storyboard', async ({
            registeredPage,
        }) => {
            const storyboardName = 'Test Storyboard'
            const storyboardsPage = new StoryboardsPage(registeredPage.page)
            await storyboardsPage.goto()

            await storyboardsPage.createStoryboard({ name: storyboardName })

            const storyboardTitle = storyboardsPage.page.locator('h1')
            await expect(storyboardTitle).toHaveText(storyboardName)
        })

        test('should display storyboards', async ({ registeredPage }) => {
            const storyboardName = 'Test Display Storyboard'

            const storyboardsPage = new StoryboardsPage(registeredPage.page)
            await storyboardsPage.goto()

            await storyboardsPage.createStoryboard({ name: storyboardName })

            const storyboardTitle = storyboardsPage.page.locator('h1')
            await expect(storyboardTitle).toHaveText(storyboardName)

            await storyboardsPage.goto()

            const title = await storyboardsPage.storyboardCardName.filter({
                hasText: storyboardName,
            })
            await expect(title).toBeVisible()
        })
    })
})
