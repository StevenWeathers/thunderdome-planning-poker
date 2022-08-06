import { Locator, Page } from '@playwright/test'

export class OrganizationPage {
    readonly page: Page
    readonly createDepartmentButton: Locator
    readonly createTeamButton: Locator
    readonly departmentNameFormField: Locator
    readonly teamNameField: Locator

    constructor(page: Page) {
        this.page = page
        this.createDepartmentButton = page.locator('button', {
            hasText: 'Create Department',
        })
        this.createTeamButton = page.locator('button', {
            hasText: 'Create Team',
        })
        this.departmentNameFormField = page.locator(
            'form[name="createDepartment"] [name="departmentName"]',
        )
        this.teamNameField = page.locator(
            'form[name="createTeam"] [name="teamName"]',
        )
    }

    async goto(id) {
        await this.page.goto(`/organization/${id}`)
    }

    async createDepartment({ name }) {
        await this.createDepartmentButton.click()
        await this.departmentNameFormField.fill(name)
        await this.departmentNameFormField.press('Enter')
    }

    async createTeam({ name }) {
        await this.createTeamButton.click()
        await this.teamNameField.fill(name)
        await this.teamNameField.press('Enter')
    }
}
