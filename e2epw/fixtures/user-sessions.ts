import {Browser, Locator, Page, test as base} from '@playwright/test';

const profileLinkSelector = '[data-testid="userprofile-link"]';

export const appPageUrls = {
    Landing: '/',
    Battles: '/battles',
    Login: '/login',
    Register: '/register'
}

class AdminPage {
    page: Page;
    profileLink: Locator;

    constructor(page: Page) {
        this.page = page;
        this.profileLink = page.locator(profileLinkSelector);
    }

    static async create(browser: Browser) {
        const context = await browser.newContext({storageState: 'storage/adminStorageState.json'});
        const page = await context.newPage();
        return new AdminPage(page);
    }
}

class GuestPage {
    page: Page;
    profileLink: Locator;

    constructor(page: Page) {
        this.page = page;
        this.profileLink = page.locator(profileLinkSelector);
    }

    static async create(browser: Browser) {
        const context = await browser.newContext({storageState: 'storage/guestStorageState.json'});
        const page = await context.newPage();
        return new GuestPage(page);
    }
}

class RegisteredPage {
    page: Page;
    profileLink: Locator;

    constructor(page: Page) {
        this.page = page;
        this.profileLink = page.locator(profileLinkSelector);
    }

    static async create(browser: Browser) {
        const context = await browser.newContext({storageState: 'storage/registeredStorageState.json'});
        const page = await context.newPage();
        return new RegisteredPage(page);
    }
}

type MyFixtures = {
    adminPage: AdminPage;
    registeredPage: RegisteredPage;
    guestPage: GuestPage;
};

// Extend base test by providing predefined pages.
// This new "test" can be used in multiple test files, and each of them will get the fixtures.
export const test = base.extend<MyFixtures>({
    adminPage: async ({browser}, use) => {
        await use(await AdminPage.create(browser));
    },
    guestPage: async ({browser}, use) => {
        await use(await GuestPage.create(browser));
    },
    registeredPage: async ({browser}, use) => {
        await use(await RegisteredPage.create(browser));
    },
});

export {expect} from '@playwright/test';