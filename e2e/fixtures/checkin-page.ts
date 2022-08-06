import { Page } from '@playwright/test'

export class TeamCheckinPage {
    readonly page: Page

    constructor(page: Page) {
        this.page = page
    }

    async goto(id) {
        await this.page.goto(`/team/${id}/checkin`)
    }

    async gotoOrg(orgId, teamId) {
        await this.page.goto(`/organization/${orgId}/team/${teamId}/checkin`)
    }

    async gotoOrgDept(orgId, deptId, teamId) {
        await this.page.goto(
            `/organization/${orgId}/department/${deptId}/team/${teamId}/checkin`,
        )
    }
}
