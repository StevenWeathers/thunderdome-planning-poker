import { test } from '../fixtures/user-sessions'
import { expect } from '@playwright/test'
import { TeamsPage } from '../fixtures/teams-page'

test.describe('Teams page', () => {
    test.describe('Unauthenticated user', () => {
        test('redirects to login', async ({ page }) => {
            const teamsPage = new TeamsPage(page)
            await teamsPage.goto()

            const title = teamsPage.page.locator('[data-formtitle="login"]')
            await expect(title).toHaveText('Login')
        })
    })

    test.describe('Registered user', () => {
        test('successfully loads page', async ({ registeredPage }) => {
            const teamsPage = new TeamsPage(registeredPage.page)
            await teamsPage.goto()

            await expect(
                teamsPage.page.locator('h2', { hasText: 'Organizations' }),
            ).toBeVisible()
            await expect(
                teamsPage.page.locator('h2', { hasText: 'Teams' }),
            ).toBeVisible()
        })

        test.describe('Create Organization', () => {
            test('should successfully submit and navigate to new organization page', async ({
                registeredPage,
            }) => {
                const teamsPage = new TeamsPage(registeredPage.page)
                await teamsPage.goto()

                await teamsPage.createOrganization({
                    name: 'Test Organization',
                })

                await expect(
                    teamsPage.page.locator('h2', { hasText: 'Departments' }),
                ).toBeVisible()
                await expect(
                    teamsPage.page.locator('h2', { hasText: 'Teams' }),
                ).toBeVisible()
                await expect(
                    teamsPage.page.locator('h2', { hasText: 'Users' }),
                ).toBeVisible()
            })
        })

        test.describe('Create Team', () => {
            test('should successfully submit and navigate to new team page', async ({
                registeredPage,
            }) => {
                const teamsPage = new TeamsPage(registeredPage.page)
                await teamsPage.goto()

                await teamsPage.createTeam({ name: 'Test Team' })

                await expect(
                    teamsPage.page.locator('h2', { hasText: 'Battles' }),
                ).toBeVisible()
                await expect(
                    teamsPage.page.locator('h2', { hasText: 'Retros' }),
                ).toBeVisible()
                await expect(
                    teamsPage.page.locator('h2', { hasText: 'Storyboards' }),
                ).toBeVisible()
                await expect(
                    teamsPage.page.locator('h2', { hasText: 'Users' }),
                ).toBeVisible()
            })
        })
    })
})
