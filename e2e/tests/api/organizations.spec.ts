import { expect, test } from '@playwright/test'
import { adminAPIUser } from '../../fixtures/db/adminapi-user'
import { apiUser } from '../../fixtures/db/api-user'
import { baseUrl } from '../../playwright.config'

const baseURL = `${baseUrl}/api/`
const userProfileEndpoint = `auth/user`

// Request context is reused by all tests in the file.
let apiContext
let adminApiContext
let adminUser
let user

test.beforeAll(async ({ playwright }) => {
    apiContext = await playwright.request.newContext({
        baseURL,
        extraHTTPHeaders: {
            'X-API-Key': apiUser.apikey,
        },
    })
    adminApiContext = await playwright.request.newContext({
        baseURL,
        extraHTTPHeaders: {
            'X-API-Key': adminAPIUser.apikey,
        },
    })
    const au = await adminApiContext.get(userProfileEndpoint)
    const auj = await au.json()
    adminUser = auj.data
    const u = await apiContext.get(userProfileEndpoint)
    const uj = await u.json()
    user = uj.data
})

test.afterAll(async ({}) => {
    // Dispose all responses.
    await apiContext.dispose()
})

test.describe('registered user', () => {
    test(`GET /users/{userId}/organizations should return empty array when no organizations associated to user`, async () => {
        const b = await adminApiContext.get(
            `users/${adminUser.id}/organizations`,
        )
        expect(b.ok()).toBeTruthy()

        const organizations = await b.json()
        expect(organizations.data).toMatchObject([])
    })

    test(`POST /users/{userId}/organizations should create organization`, async () => {
        const organizationName = 'Test API Create Organization'

        const b = await apiContext.post(`users/${user.id}/organizations`, {
            data: {
                name: organizationName,
            },
        })
        expect(b.ok()).toBeTruthy()
        const organization = await b.json()
        expect(organization.data).toMatchObject({
            name: organizationName,
        })
    })

    test(`GET /users/{userId}/organizations should return object in array when organizations associated to user`, async () => {
        const organizationName = 'Test API Organizations'

        const b = await apiContext.post(`users/${user.id}/organizations`, {
            data: {
                name: organizationName,
            },
        })
        expect(b.ok()).toBeTruthy()

        const bs = await apiContext.get(`users/${user.id}/organizations`)
        expect(bs.ok()).toBeTruthy()
        const organizations = await bs.json()
        expect(organizations.data).toContainEqual(
            expect.objectContaining({
                name: organizationName,
            }),
        )
    })
})
