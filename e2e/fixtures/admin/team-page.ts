import { Page } from '@playwright/test'

export class AdminTeamPage {
    readonly page: Page

    constructor(page: Page) {
        this.page = page
    }

    async goto(id) {
        await this.page.goto(`/admin/teams/${id}`)
    }

    async gotoOrgTeam(orgId, teamId) {
        await this.page.goto(`/admin/organizations/${orgId}/team/${teamId}`)
    }

    async gotoOrgDeptTeam(orgId, deptId, teamId) {
        await this.page.goto(
            `/admin/organizations/${orgId}/department/${deptId}/team/${teamId}`,
        )
    }
}
