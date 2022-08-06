import { Browser, Page, test as base } from '@playwright/test'

class AdminPage {
    page: Page
    user: {
        id: string
        email: string
    }

    constructor(page: Page) {
        this.page = page
    }

    static async create(browser: Browser) {
        const context = await browser.newContext({
            storageState: 'storage/adminStorageState.json',
        })
        const page = await context.newPage()
        const adminPage = new AdminPage(page)
        await adminPage.loadUser()
        return adminPage
    }

    public async loadUser() {
        const u = await this.page.request.get('/api/auth/user')
        const user = await u.json()
        this.user = user.data
    }
}

class RegisteredPage {
    page: Page
    user: {
        id: string
        name: string
        email: string
    }

    constructor(page: Page) {
        this.page = page
    }

    static async create(browser: Browser) {
        const context = await browser.newContext({
            storageState: 'storage/registeredStorageState.json',
        })
        const page = await context.newPage()
        const regPage = new RegisteredPage(page)
        await regPage.loadUser() // bootstrap the user info

        return regPage
    }

    public async loadUser() {
        const u = await this.page.request.get('/api/auth/user')
        const user = await u.json()
        this.user = user.data
    }

    public async createBattle(battle) {
        const b = await this.page.request.post(
            `/api/users/${this.user.id}/battles`,
            {
                data: battle,
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
}

class VerifiedPage {
    page: Page
    user: {
        id: string
        email: string
    }

    constructor(page: Page) {
        this.page = page
    }

    static async create(browser: Browser) {
        const context = await browser.newContext({
            storageState: 'storage/verifiedStorageState.json',
        })
        const page = await context.newPage()
        const vPage = new VerifiedPage(page)
        await vPage.loadUser()
        return vPage
    }

    public async loadUser() {
        const u = await this.page.request.get('/api/auth/user')
        const user = await u.json()
        this.user = user.data
    }

    public async createBattle(battle) {
        const b = await this.page.request.post(
            `/api/users/${this.user.id}/battles`,
            {
                data: battle,
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
}

class GuestPage {
    page: Page
    user: {
        id: string
        name: string
    }

    constructor(page: Page) {
        this.page = page
    }

    public async loadUser() {
        const u = await this.page.request.get('/api/auth/user')
        const user = await u.json()
        this.user = user.data
    }

    static async create(browser: Browser) {
        const context = await browser.newContext({
            storageState: 'storage/guestStorageState.json',
        })
        const page = await context.newPage()
        const guestPage = new GuestPage(page)
        await guestPage.loadUser()
        return guestPage
    }
}

type MyFixtures = {
    adminPage: AdminPage
    registeredPage: RegisteredPage
    verifiedPage: VerifiedPage
    guestPage: GuestPage
}

export const test = base.extend<MyFixtures>({
    adminPage: async ({ browser }, use) => {
        await use(await AdminPage.create(browser))
    },
    registeredPage: async ({ browser }, use) => {
        await use(await RegisteredPage.create(browser))
    },
    verifiedPage: async ({ browser }, use) => {
        await use(await VerifiedPage.create(browser))
    },
    guestPage: async ({ browser }, use) => {
        await use(await GuestPage.create(browser))
    },
})

export { expect } from '@playwright/test'
