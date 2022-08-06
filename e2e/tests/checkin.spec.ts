import { expect, test } from '../fixtures/user-sessions'
import { TeamCheckinPage } from '../fixtures/checkin-page'

test.describe('Team Checkin page', () => {
    test.describe('Unauthenticated user', () => {
        test('redirects to login', async ({ page }) => {
            const teamPage = new TeamCheckinPage(page)
            await teamPage.goto('bbaf82ef-a2d3-4e9a-b824-5e56a03ac3aa')

            const title = teamPage.page.locator('[data-formtitle="login"]')
            await expect(title).toHaveText('Login')
        })
    })

    test.describe('Registered user', () => {
        const sharedTestTeamName = 'E2E SHARED TEST TEAM'
        let sharedTeam = { id: '' }
        test.beforeAll(async ({ registeredPage }) => {
            sharedTeam = await registeredPage.createTeam(sharedTestTeamName)
        })

        test('loads team checkin page', async ({ registeredPage }) => {
            const teamPage = new TeamCheckinPage(registeredPage.page)

            await teamPage.goto(sharedTeam.id)
            await expect(
                teamPage.page.locator('a', { hasText: sharedTestTeamName }),
            ).toBeVisible()
            await expect(teamPage.page.locator('h1')).toHaveText('Check In')
        })

        test('loads organization team checkin page', async ({
            registeredPage,
        }) => {
            const testTeamName = 'E2E TEST TEAM'
            const testOrgName = 'E2E TEST ORGANIZATION'
            const teamPage = new TeamCheckinPage(registeredPage.page)
            const org = await registeredPage.createOrg(testOrgName)
            const team = await registeredPage.createTeam(testTeamName)

            await teamPage.gotoOrg(org.id, team.id)
            await expect(
                teamPage.page.locator('a', { hasText: testOrgName }),
            ).toBeVisible()
            await expect(
                teamPage.page.locator('a', { hasText: testTeamName }),
            ).toBeVisible()
            await expect(teamPage.page.locator('h1')).toHaveText('Check In')
        })

        test('loads department team checkin page', async ({
            registeredPage,
        }) => {
            const testTeamName = 'E2E TEST TEAM'
            const testDepartmentName = 'E2E TEST DEPARTMENT'
            const testOrgName = 'E2E TEST ORGANIZATION'
            const teamPage = new TeamCheckinPage(registeredPage.page)
            const org = await registeredPage.createOrg(testOrgName)
            const dept = await registeredPage.createOrgDepartment(
                org.id,
                testDepartmentName,
            )
            const team = await registeredPage.createDepartmentTeam(
                org.id,
                dept.id,
                testTeamName,
            )

            await teamPage.gotoOrgDept(org.id, dept.id, team.id)
            await expect(
                teamPage.page.locator('a', { hasText: testOrgName }),
            ).toBeVisible()
            await expect(
                teamPage.page.locator('a', { hasText: testDepartmentName }),
            ).toBeVisible()
            await expect(
                teamPage.page.locator('a', { hasText: testTeamName }),
            ).toBeVisible()
            await expect(teamPage.page.locator('h1')).toHaveText('Check In')
        })

        test('can checkin once today', async ({ registeredPage }) => {
            const teamPage = new TeamCheckinPage(registeredPage.page)

            await teamPage.goto(sharedTeam.id)

            await expect(teamPage.page.locator('h1')).toHaveText('Check In')

            // check in
            await teamPage.page.locator('[data-testid="check-in"]').click()
            await teamPage.page
                .locator('#yesterday >> p')
                .fill('Yesterday I fixed bugs')
            await teamPage.page
                .locator('#today >> p')
                .fill('Today I will write e2e tests')
            await teamPage.page.locator('data-testid=save').click()

            await expect(
                teamPage.page.locator('[data-testid="check-in"]'),
            ).toBeDisabled()
            await expect(
                teamPage.page.locator('[data-testid="checkin"]'),
            ).toBeVisible()
            await expect(
                teamPage.page.locator('[data-testid="checkin-yesterday"]'),
            ).toHaveText('Yesterday I fixed bugs')
            await expect(
                teamPage.page.locator('[data-testid="checkin-today"]'),
            ).toHaveText('Today I will write e2e tests')
            await expect(
                teamPage.page.locator('[data-testid="checkin-blockers"]'),
            ).not.toBeVisible()
            await expect(
                teamPage.page.locator('[data-testid="checkin-discuss"]'),
            ).not.toBeVisible()

            // edit checkin - @TODO separate this into its own test
            await teamPage.page.locator('[data-testid="checkin-edit"]').click()
            await teamPage.page
                .locator('#blockers >> p')
                .fill('Blocked by procrastination')
            await teamPage.page.locator('#discuss >> p').fill('Whats next?')
            await teamPage.page.locator('data-testid=save').click()

            await expect(
                teamPage.page.locator('[data-testid="checkin-blockers"]'),
            ).toHaveText('Blocked by procrastination')
            await expect(
                teamPage.page.locator('[data-testid="checkin-discuss"]'),
            ).toHaveText('Whats next?')

            // delete checkin - @TODO separate this into its own test
            await teamPage.page
                .locator('[data-testid="checkin-delete"]')
                .click()
            await expect(
                teamPage.page.locator('[data-testid="check-in"]'),
            ).not.toBeDisabled()
            await expect(
                teamPage.page.locator('[data-testid="checkin"]'),
            ).not.toBeVisible()
        })
    })
})
