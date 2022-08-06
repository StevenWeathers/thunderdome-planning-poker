import { Page } from '@playwright/test'

export class TeamPage {
    readonly page: Page

    constructor(page: Page) {
        this.page = page
    }

    async goto(id) {
        await this.page.goto(`/team/${id}`)
    }

    async gotoOrgTeam(orgId, teamId) {
        await this.page.goto(`/organization/${orgId}/team/${teamId}`)
    }

    async gotoOrgDeptTeam(orgId, deptId, teamId) {
        await this.page.goto(
            `/organization/${orgId}/department/${deptId}/team/${teamId}`,
        )
    }
}
