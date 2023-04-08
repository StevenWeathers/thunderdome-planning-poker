import {chromium, FullConfig} from '@playwright/test'

async function globalTeardown(config: FullConfig) {
    const baseUrl = config.projects[0].use.baseURL
    const browser = await chromium.launch()

    const adminContext = await browser.newContext({
        storageState: 'storage/adminStorageState.json',
        baseURL: baseUrl,
    })
    const adminPage = await adminContext.newPage()
    const au = await adminPage.request.get(`/api/auth/user`)
    const adminUser = await au.json()
    await adminPage.request.delete(`/api/users/${adminUser.data.id}`)

    const registeredContext = await browser.newContext({
        storageState: 'storage/registeredStorageState.json',
        baseURL: baseUrl,
    })
    const registeredPage = await registeredContext.newPage()
    const ru = await registeredPage.request.get(`/api/auth/user`)
    const registeredUser = await ru.json()
    await registeredPage.request.delete(`/api/users/${registeredUser.data.id}`)

    const verifiedContext = await browser.newContext({
        storageState: 'storage/verifiedStorageState.json',
        baseURL: baseUrl,
    })
    const verifiedPage = await verifiedContext.newPage()
    const vu = await verifiedPage.request.get(`/api/auth/user`)
    const verifiedUser = await vu.json()
    await verifiedPage.request.delete(`/api/users/${verifiedUser.data.id}`)

    const guestContext = await browser.newContext({
        storageState: 'storage/guestStorageState.json',
        baseURL: baseUrl,
    })
    const guestPage = await guestContext.newPage()
    const gu = await guestPage.request.get(`/api/auth/user`)
    const guestUser = await gu.json()
    await guestPage.request.delete(`/api/users/${guestUser.data.id}`)

    const deleteGuestContext = await browser.newContext({
        storageState: 'storage/deleteGuestStorageState.json',
        baseURL: baseUrl,
    })
    const deleteGuestPage = await deleteGuestContext.newPage()
    const dgu = await deleteGuestPage.request.get(`/api/auth/user`)
    const deleteGuestUser = await dgu.json()
    await deleteGuestPage.request.delete(`/api/users/${deleteGuestUser.data.id}`)

    const deleteRegisteredContext = await browser.newContext({
        storageState: 'storage/deleteRegisteredStorageState.json',
        baseURL: baseUrl,
    })
    const deleteRegisteredPage = await deleteRegisteredContext.newPage()
    const dru = await deleteRegisteredPage.request.get(`/api/auth/user`)
    const deleteRegisteredUser = await dru.json()
    await deleteRegisteredPage.request.delete(`/api/users/${deleteRegisteredUser.data.id}`)

    await browser.close()
}

export default globalTeardown
