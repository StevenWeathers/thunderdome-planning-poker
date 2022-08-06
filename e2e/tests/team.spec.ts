import { expect, test } from '../fixtures/user-sessions'
import { TeamPage } from '../fixtures/team-page'

test.describe('Team page', () => {
    test('Unauthenticated user redirects to login', async ({ page }) => {
        const teamPage = new TeamPage(page)
        await teamPage.goto('bbaf82ef-a2d3-4e9a-b824-5e56a03ac3aa')

        const title = teamPage.page.locator('[data-formtitle="login"]')
        await expect(title).toHaveText('Login')
    })

    test('Registered user loads team page', async ({ registeredPage }) => {
        const testTeamName = 'E2E TEST TEAM'
        const teamPage = new TeamPage(registeredPage.page)
        const team = await registeredPage.createTeam(testTeamName)

        await teamPage.goto(team.id)
        await expect(teamPage.page.locator('h1')).toContainText(testTeamName)
    })

    test('Registered user loads organization team page', async ({
        registeredPage,
    }) => {
        const testTeamName = 'E2E TEST TEAM'
        const testOrgName = 'E2E TEST ORGANIZATION'
        const teamPage = new TeamPage(registeredPage.page)
        const org = await registeredPage.createOrg(testOrgName)
        const team = await registeredPage.createTeam(testTeamName)

        await teamPage.gotoOrgTeam(org.id, team.id)
        await expect(
            teamPage.page.locator('a', { hasText: testOrgName }),
        ).toBeVisible()
        await expect(teamPage.page.locator('h1')).toContainText(testTeamName)
    })

    test('Registered user loads department team page', async ({
        registeredPage,
    }) => {
        const testTeamName = 'E2E TEST TEAM'
        const testDepartmentName = 'E2E TEST DEPARTMENT'
        const testOrgName = 'E2E TEST ORGANIZATION'
        const teamPage = new TeamPage(registeredPage.page)
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

        await teamPage.gotoOrgDeptTeam(org.id, dept.id, team.id)
        await expect(
            teamPage.page.locator('a', { hasText: testOrgName }),
        ).toBeVisible()
        await expect(
            teamPage.page.locator('a', { hasText: testDepartmentName }),
        ).toBeVisible()
        await expect(teamPage.page.locator('h1')).toContainText(testTeamName)
    })

    test('can add user to team', async ({ registeredPage }) => {
        const verifiedEmail = 'e2everified@thunderdome.dev'
        const testTeamName = 'E2E TEST TEAM'
        const teamPage = new TeamPage(registeredPage.page)
        const team = await registeredPage.createTeam(testTeamName)

        await teamPage.goto(team.id)
        await expect(teamPage.page.locator('h1')).toContainText(testTeamName)

        await teamPage.page.locator('[data-testid="user-add"]').click()
        await teamPage.page
            .locator('input[name="userEmail"]')
            .fill(verifiedEmail)
        await teamPage.page
            .locator('select[name="userRole"]')
            .selectOption('MEMBER')
        await teamPage.page.locator('[data-testid="useradd-confirm"]').click()

        await expect(
            teamPage.page.locator('[data-testid="user-email"]', {
                hasText: verifiedEmail,
            }),
        ).toBeVisible()
    })
})
