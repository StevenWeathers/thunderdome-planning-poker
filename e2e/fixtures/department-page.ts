import { Locator, Page } from '@playwright/test'

export class DepartmentPage {
    readonly page: Page
    readonly createTeamButton: Locator
    readonly teamNameField: Locator

    constructor(page: Page) {
        this.page = page
        this.createTeamButton = page.locator('button', {
            hasText: 'Create Team',
        })
        this.teamNameField = page.locator(
            'form[name="createTeam"] [name="teamName"]',
        )
    }

    async goto(orgId, deptId) {
        await this.page.goto(`/organization/${orgId}/department/${deptId}`)
    }

    async createTeam({ name }) {
        await this.createTeamButton.click()
        await this.teamNameField.fill(name)
        await this.teamNameField.press('Enter')
    }
}
