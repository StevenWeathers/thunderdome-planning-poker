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
    test(`GET /users/{userId}/storyboards should return empty array when no storyboards associated to user`, async () => {
        const b = await adminApiContext.get(`users/${adminUser.id}/storyboards`)
        expect(b.ok()).toBeTruthy()

        const storyboards = await b.json()
        expect(storyboards.data).toMatchObject([])
    })

    test(`POST /users/{userId}/storyboards should create storyboard`, async () => {
        const storyboardName = 'Test API Create Storyboard'

        const b = await apiContext.post(`users/${user.id}/storyboards`, {
            data: {
                storyboardName,
            },
        })
        expect(b.ok()).toBeTruthy()
        const storyboard = await b.json()
        expect(storyboard.data).toMatchObject({
            name: storyboardName,
        })
    })

    test(`GET /users/{userId}/storyboards should return object in array when storyboards associated to user`, async () => {
        const storyboardName = 'Test API Storyboards'

        const b = await apiContext.post(`users/${user.id}/storyboards`, {
            data: {
                storyboardName,
            },
        })
        expect(b.ok()).toBeTruthy()

        const bs = await apiContext.get(`users/${user.id}/storyboards`)
        expect(bs.ok()).toBeTruthy()
        const storyboards = await bs.json()
        expect(storyboards.data).toContainEqual(
            expect.objectContaining({
                name: storyboardName,
            }),
        )
    })

    test(`POST /teams/{teamId}/users/{userId}/storyboards should create storyboard`, async () => {
        const storyboardName = 'Test API Create Team Storyboard'

        const t = await apiContext.post(`users/${user.id}/teams`, {
            data: {
                name: 'test team create retro',
            },
        })
        const { data: team } = await t.json()

        const b = await apiContext.post(
            `teams/${team.id}/users/${user.id}/storyboards`,
            {
                data: {
                    storyboardName,
                },
            },
        )
        expect(b.ok()).toBeTruthy()
        const storyboard = await b.json()
        expect(storyboard.data).toMatchObject({
            name: storyboardName,
        })

        const bs = await apiContext.get(`teams/${team.id}/storyboards`)
        expect(bs.ok()).toBeTruthy()
        const storyboards = await bs.json()
        expect(storyboards.data).toContainEqual(
            expect.objectContaining({
                name: storyboardName,
            }),
        )
    })
})
