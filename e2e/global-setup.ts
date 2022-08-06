import { chromium, expect, FullConfig } from '@playwright/test'
import { RegisterPage } from './fixtures/register-page'
import { LoginPage } from './fixtures/login-page'
import { setupDB } from './fixtures/db/setup'
import { setupAdminUser } from './fixtures/db/admin-user'
import { setupRegisteredUser } from './fixtures/db/registered-user'
import { setupVerifiedUser } from './fixtures/db/verified-user'
import { setupAPIUser } from './fixtures/db/api-user'
import { setupAdminAPIUser } from './fixtures/db/adminapi-user'

async function globalSetup(config: FullConfig) {
    const pool = setupDB({
        name: process.env.DB_NAME || 'thunderdome',
        user: process.env.DB_USER || 'thor',
        pass: process.env.DB_PASS || 'odinson',
        port: process.env.DB_PORT || '5432',
        host: process.env.DB_HOST || 'localhost',
    })

    const baseUrl = config.projects[0].use.baseURL
    const browser = await chromium.launch()

    await setupAdminAPIUser.teardown(pool)
    await setupAdminAPIUser.seed(pool)

    await setupAPIUser.teardown(pool)
    await setupAPIUser.seed(pool)

    const adminPage = await browser.newPage({
        baseURL: baseUrl,
    })
    await setupAdminUser.teardown(pool)
    const au = await setupAdminUser.seed(pool)
    const adminLoginPage = new LoginPage(adminPage)
    await adminLoginPage.goto()
    await adminLoginPage.login(au.email, au.password)
    await expect(adminLoginPage.page.locator('h1')).toHaveText('My Battles')
    await adminLoginPage.page
        .context()
        .storageState({ path: 'storage/adminStorageState.json' })

    const registeredPage = await browser.newPage({
        baseURL: baseUrl,
    })
    await setupRegisteredUser.teardown(pool)
    const ru = await setupRegisteredUser.seed(pool)
    const registeredRegisterPage = new LoginPage(registeredPage)
    await registeredRegisterPage.goto()
    await registeredRegisterPage.login(ru.email, ru.password)
    await expect(registeredRegisterPage.page.locator('h1')).toHaveText(
        'My Battles',
    )
    await registeredRegisterPage.page
        .context()
        .storageState({ path: 'storage/registeredStorageState.json' })

    const verifiedPage = await browser.newPage({
        baseURL: baseUrl,
    })
    await setupVerifiedUser.teardown(pool)
    const vu = await setupVerifiedUser.seed(pool)
    const userVerifiedPage = new LoginPage(verifiedPage)
    await userVerifiedPage.goto()
    await userVerifiedPage.login(vu.email, vu.password)
    await expect(userVerifiedPage.page.locator('h1')).toHaveText('My Battles')
    await userVerifiedPage.page
        .context()
        .storageState({ path: 'storage/verifiedStorageState.json' })

    const guestPage = await browser.newPage({
        baseURL: baseUrl,
    })
    const guestRegisterPage = new RegisterPage(guestPage)
    await guestRegisterPage.goto()
    await guestRegisterPage.createGuestUser('E2E Guest')
    await expect(guestRegisterPage.page.locator('h1')).toHaveText('My Battles')
    await guestRegisterPage.page
        .context()
        .storageState({ path: 'storage/guestStorageState.json' })

    await browser.close()
}

export default globalSetup
