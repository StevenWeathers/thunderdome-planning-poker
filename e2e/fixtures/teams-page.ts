import { Locator, Page } from '@playwright/test'

export class TeamsPage {
    readonly page: Page
    readonly createOrgButton: Locator
    readonly createTeamButton: Locator
    readonly organizationNameField: Locator
    readonly teamNameField: Locator

    constructor(page: Page) {
        this.page = page
        this.createOrgButton = page.locator('button', {
            hasText: 'Create Organization',
        })
        this.createTeamButton = page.locator('button', {
            hasText: 'Create Team',
        })
        this.organizationNameField = page.locator(
            'form[name="createOrganization"] [name="organizationName"]',
        )
        this.teamNameField = page.locator(
            'form[name="createTeam"] [name="teamName"]',
        )
    }

    async goto() {
        await this.page.goto('/teams')
    }

    async createOrganization({ name }) {
        await this.createOrgButton.click()
        await this.organizationNameField.fill(name)
        await this.organizationNameField.press('Enter')
    }

    async createTeam({ name }) {
        await this.createTeamButton.click()
        await this.teamNameField.fill(name)
        await this.teamNameField.press('Enter')
    }
}
