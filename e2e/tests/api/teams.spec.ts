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
    test(`GET /users/{userId}/teams should return empty array when no teams associated to user`, async () => {
        const b = await adminApiContext.get(`users/${adminUser.id}/teams`)
        expect(b.ok()).toBeTruthy()

        const teams = await b.json()
        expect(teams.data).toMatchObject([])
    })

    test(`POST /users/{userId}/teams should create team`, async () => {
        const teamName = 'Test API Create Team'

        const b = await apiContext.post(`users/${user.id}/teams`, {
            data: {
                name: teamName,
            },
        })
        expect(b.ok()).toBeTruthy()
        const team = await b.json()
        expect(team.data).toMatchObject({
            name: teamName,
        })
    })

    test(`GET /users/{userId}/teams should return object in array when teams associated to user`, async () => {
        const teamName = 'Test API Teams'

        const b = await apiContext.post(`users/${user.id}/teams`, {
            data: {
                name: teamName,
            },
        })
        expect(b.ok()).toBeTruthy()

        const bs = await apiContext.get(`users/${user.id}/teams`)
        expect(bs.ok()).toBeTruthy()
        const teams = await bs.json()
        expect(teams.data).toContainEqual(
            expect.objectContaining({
                name: teamName,
            }),
        )
    })
})
