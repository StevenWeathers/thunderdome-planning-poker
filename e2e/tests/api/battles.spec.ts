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
    test(`GET /users/{userId}/battles should return empty array when no battles associated to user`, async () => {
        const b = await adminApiContext.get(`users/${adminUser.id}/battles`)
        expect(b.ok()).toBeTruthy()

        const battles = await b.json()
        expect(battles.data).toMatchObject([])
    })

    test(`POST /users/{userId}/battles should create battle`, async () => {
        const pointValuesAllowed = ['0', '0.5', '1', '2', '3', '5', '8', '13']
        const battleName = 'Test API Create Battle'
        const pointAverageRounding = 'floor'
        const autoFinishVoting = false

        const b = await apiContext.post(`users/${user.id}/battles`, {
            data: {
                name: battleName,
                pointValuesAllowed,
                pointAverageRounding,
                autoFinishVoting,
            },
        })
        expect(b.ok()).toBeTruthy()
        const battle = await b.json()
        expect(battle.data).toMatchObject({
            name: battleName,
            pointValuesAllowed,
            pointAverageRounding,
            autoFinishVoting,
        })
    })

    test(`GET /users/{userId}/battles should return object in array when battles associated to user`, async () => {
        const pointValuesAllowed = ['1', '2', '3', '5', '8', '13']
        const battleName = 'Test API Battles'
        const pointAverageRounding = 'ceil'
        const autoFinishVoting = true

        const b = await apiContext.post(`users/${user.id}/battles`, {
            data: {
                name: battleName,
                pointValuesAllowed,
                pointAverageRounding,
                autoFinishVoting,
            },
        })
        expect(b.ok()).toBeTruthy()

        const bs = await apiContext.get(`users/${user.id}/battles`)
        expect(bs.ok()).toBeTruthy()
        const battles = await bs.json()
        expect(battles.data).toContainEqual(
            expect.objectContaining({
                name: battleName,
                pointValuesAllowed,
                pointAverageRounding,
                autoFinishVoting,
            }),
        )
    })

    test(`POST /teams/{teamId}/users/{userId}/battles should create battle`, async () => {
        const pointValuesAllowed = ['0', '0.5', '1', '2', '3', '5', '8', '13']
        const battleName = 'Test API Create Battle'
        const pointAverageRounding = 'floor'
        const autoFinishVoting = false

        const t = await apiContext.post(`users/${user.id}/teams`, {
            data: {
                name: 'test team create battle',
            },
        })
        const { data: team } = await t.json()

        const b = await apiContext.post(
            `teams/${team.id}/users/${user.id}/battles`,
            {
                data: {
                    name: battleName,
                    pointValuesAllowed,
                    pointAverageRounding,
                    autoFinishVoting,
                },
            },
        )
        expect(b.ok()).toBeTruthy()
        const battle = await b.json()
        expect(battle.data).toMatchObject({
            name: battleName,
            pointValuesAllowed,
            pointAverageRounding,
            autoFinishVoting,
        })

        const bs = await apiContext.get(`teams/${team.id}/battles`)
        expect(bs.ok()).toBeTruthy()
        const battles = await bs.json()
        expect(battles.data).toContainEqual(
            expect.objectContaining({
                name: battleName,
            }),
        )
    })
})
