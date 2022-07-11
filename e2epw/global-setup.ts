import {chromium, expect, FullConfig} from '@playwright/test';
import {RegisterPage} from './fixtures/register-page';

async function globalSetup(config: FullConfig) {
    const baseUrl = config.projects[0].use.baseURL;
    const browser = await chromium.launch();

    const adminPage = await browser.newPage({
        baseURL: baseUrl
    });
    const adminRegisterPage = new RegisterPage(adminPage);
    await adminRegisterPage.goto();
    await adminRegisterPage.createRegisteredUser("E2E Admin", "e2eadmin@thunderdome.dev", "e2etestpass", "e2etestpass");
    await expect(adminRegisterPage.page.locator('h1')).toHaveText('My Battles');
    await adminRegisterPage.page.context().storageState({path: 'storage/adminStorageState.json'});

    const guestPage = await browser.newPage({
        baseURL: baseUrl
    });
    const guestRegisterPage = new RegisterPage(guestPage);
    await guestRegisterPage.goto();
    await guestRegisterPage.createGuestUser("E2E Guest");
    await expect(guestRegisterPage.page.locator('h1')).toHaveText('My Battles');
    await guestRegisterPage.page.context().storageState({path: 'storage/guestStorageState.json'});

    const registeredPage = await browser.newPage({
        baseURL: baseUrl
    });
    const registeredRegisterPage = new RegisterPage(registeredPage);
    await registeredRegisterPage.goto();
    await registeredRegisterPage.createRegisteredUser("E2E Registered", "e2eregistered@thunderdome.dev", "e2etestpass", "e2etestpass");
    await expect(registeredRegisterPage.page.locator('h1')).toHaveText('My Battles');
    await registeredRegisterPage.page.context().storageState({path: 'storage/registeredStorageState.json'});

    await browser.close();
}

export default globalSetup;