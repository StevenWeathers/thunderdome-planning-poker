import { expect, test } from '../fixtures/user-sessions'
import { DepartmentPage } from '../fixtures/department-page'

test.describe('Department page', () => {
    test.describe('Unauthenticated user', () => {
        test('redirects to login', async ({ page }) => {
            const departmentPage = new DepartmentPage(page)
            await departmentPage.goto(
                'bbaf82ef-a2d3-4e9a-b824-5e56a03ac3aa',
                'bbaf82ef-a2d3-4e9a-b824-5e56a03ac3aa',
            )

            const title = departmentPage.page.locator(
                '[data-formtitle="login"]',
            )
            await expect(title).toHaveText('Login')
        })
    })

    test.describe('Registered user', () => {
        test('loads page successfully', async ({ registeredPage }) => {
            const testOrgName = 'E2E TEST ORGANIZATION'
            const testDepartmentName = 'E2E TEST DEPARTMENT'
            const departmentPage = new DepartmentPage(registeredPage.page)
            const org = await registeredPage.createOrg(testOrgName)
            const dept = await registeredPage.createOrgDepartment(
                org.id,
                testDepartmentName,
            )

            await departmentPage.goto(org.id, dept.id)
            await expect(departmentPage.page.locator('h1')).toContainText(
                testDepartmentName,
            )
        })
    })
})
