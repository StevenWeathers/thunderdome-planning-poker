import {expect, test} from '../fixtures/user-sessions';
import {TeamPage} from "../fixtures/team-page";

test.describe('Team page', () => {
    test.describe('Unauthenticated user', () => {
        test('redirects to login', async ({page}) => {
            const teamPage = new TeamPage(page);
            await teamPage.goto('bbaf82ef-a2d3-4e9a-b824-5e56a03ac3aa');

            const title = teamPage.page.locator('[data-formtitle="login"]');
            await expect(title).toHaveText('Login');
        })
    })

    test.describe('Registered user', () => {
        test('loads team page', async ({registeredPage}) => {
            const testTeamName = 'E2E TEST TEAM'
            const teamPage = new TeamPage(registeredPage.page);
            const team = await registeredPage.createTeam(testTeamName)

            await teamPage.goto(team.id);
            await expect(teamPage.page.locator('h1')).toContainText(testTeamName)
        })

        test('loads organization team page', async ({registeredPage}) => {
            const testTeamName = 'E2E TEST ORGANIZATION TEAM'
            const testOrgName = 'E2E TEST ORGANIZATION'
            const teamPage = new TeamPage(registeredPage.page);
            const org = await registeredPage.createOrg(testOrgName)
            const team = await registeredPage.createTeam(testTeamName)

            await teamPage.gotoOrgTeam(org.id, team.id);
            await expect(teamPage.page.locator('h1')).toContainText(testTeamName)
        })

        test('loads department team page', async ({registeredPage}) => {
            const testTeamName = 'E2E TEST ORGANIZATION TEAM'
            const testDepartmentName = 'E2E TEST DEPARTMENT'
            const testOrgName = 'E2E TEST ORGANIZATION'
            const teamPage = new TeamPage(registeredPage.page);
            const org = await registeredPage.createOrg(testOrgName)
            const dept = await registeredPage.createOrgDepartment(org.id, testDepartmentName)
            const team = await registeredPage.createDepartmentTeam(org.id, dept.id, testTeamName)

            await teamPage.gotoOrgTeam(org.id, team.id);
            await expect(teamPage.page.locator('h1')).toContainText(testTeamName)
        })
    })
})