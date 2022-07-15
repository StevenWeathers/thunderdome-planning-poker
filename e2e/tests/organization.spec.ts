import {expect, test} from '../fixtures/user-sessions';
import {OrganizationPage} from "../fixtures/organization-page";

test.describe('Organization Page', () => {
    test.describe('Unauthenticated user', () => {
        test('redirects to login', async ({page}) => {
            const orgPage = new OrganizationPage(page);
            await orgPage.goto('bbaf82ef-a2d3-4e9a-b824-5e56a03ac3aa');

            const title = orgPage.page.locator('[data-formtitle="login"]');
            await expect(title).toHaveText('Login');
        })
    })

    test.describe('Registered user', () => {
        test('loads page successfully', async ({registeredPage}) => {
            const testOrgName = 'E2E TEST ORGANIZATION'
            const orgPage = new OrganizationPage(registeredPage.page);
            const org = await registeredPage.createOrg(testOrgName)

            await orgPage.goto(org.id);
            await expect(orgPage.page.locator('h1')).toContainText(testOrgName)
        })
    })
})