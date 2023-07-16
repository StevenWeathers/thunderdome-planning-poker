import { Browser, Page, test as base } from '@playwright/test'

class UserPage {
    page: Page
    user: {
        id: string
        name: string
        email: string
    }

    constructor(page: Page) {
        this.page = page
    }

    static async create(browser: Browser, storageState: string) {
        const context = await browser.newContext({
            storageState,
        })
        const page = await context.newPage()
        const regPage = new UserPage(page)
        await regPage.loadUser() // bootstrap the user info

        return regPage
    }

    public async loadUser() {
        const u = await this.page.request.get('/api/auth/user')
        const user = await u.json()
        this.user = user.data
    }

    public async createPokerGame(game) {
        const b = await this.page.request.post(
            `/api/users/${this.user.id}/battles`,
            {
                data: game,
            },
        )
        const res = await b.json()
        return res.data
    }

    public async createRetro(retro) {
        const b = await this.page.request.post(
            `/api/users/${this.user.id}/retros`,
            {
                data: retro,
            },
        )
        const res = await b.json()
        return res.data
    }

    public async createStoryboard(storyboard) {
        const b = await this.page.request.post(
            `/api/users/${this.user.id}/storyboards`,
            {
                data: storyboard,
            },
        )
        const res = await b.json()
        return res.data
    }

    public async createOrg(name: string) {
        const ruo = await this.page.request.post(
            `/api/users/${this.user.id}/organizations`,
            {
                data: {
                    name: name,
                },
            },
        )
        const regUserOrg = await ruo.json()
        return regUserOrg.data
    }

    public async createTeam(name: string) {
        const t = await this.page.request.post(
            `/api/users/${this.user.id}/teams`,
            {
                data: {
                    name: name,
                },
            },
        )
        const team = await t.json()
        return team.data
    }

    public async createOrgTeam(orgId, name) {
        const t = await this.page.request.post(
            `/api/organizations/${orgId}/teams`,
            {
                data: {
                    name: name,
                },
            },
        )
        const team = await t.json()
        return team.data
    }

    public async createOrgDepartment(orgId, name) {
        const d = await this.page.request.post(
            `/api/organizations/${orgId}/departments`,
            {
                data: {
                    name: name,
                },
            },
        )
        const department = await d.json()
        return department.data
    }

    public async createDepartmentTeam(orgId, deptId, name) {
        const dt = await this.page.request.post(
            `/api/organizations/${orgId}/departments/${deptId}/teams`,
            {
                data: {
                    name: name,
                },
            },
        )
        const team = await dt.json()
        return team.data
    }

    public async createApikey(name) {
        const k = await this.page.request.post(
            `/api/users/${this.user.id}/apikeys`,
            {
                data: {
                    name,
                },
            },
        )
        const res = await k.json()
        return res.data
    }
}

type MyFixtures = {
    adminPage: UserPage
    registeredPage: UserPage
    verifiedPage: UserPage
    guestPage: UserPage
    deleteGuestPage: UserPage
    deleteRegisteredPage: UserPage
}

export const test = base.extend<MyFixtures>({
    adminPage: async ({ browser }, use) => {
        await use(
            await UserPage.create(browser, 'storage/adminStorageState.json'),
        )
    },
    registeredPage: async ({ browser }, use) => {
        await use(
            await UserPage.create(
                browser,
                'storage/registeredStorageState.json',
            ),
        )
    },
    verifiedPage: async ({ browser }, use) => {
        await use(
            await UserPage.create(browser, 'storage/verifiedStorageState.json'),
        )
    },
    guestPage: async ({ browser }, use) => {
        await use(
            await UserPage.create(browser, 'storage/guestStorageState.json'),
        )
    },
    deleteGuestPage: async ({ browser }, use) => {
        await use(
            await UserPage.create(
                browser,
                'storage/deleteGuestStorageState.json',
            ),
        )
    },
    deleteRegisteredPage: async ({ browser }, use) => {
        await use(
            await UserPage.create(
                browser,
                'storage/deleteRegisteredStorageState.json',
            ),
        )
    },
})

export { expect } from '@playwright/test'
