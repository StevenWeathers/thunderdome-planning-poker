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
    test(`GET /users/{userId}/retros should return empty array when no retros associated to user`, async () => {
        const b = await adminApiContext.get(`users/${adminUser.id}/retros`)
        expect(b.ok()).toBeTruthy()

        const retros = await b.json()
        expect(retros.data).toMatchObject([])
    })

    test(`POST /users/{userId}/retros should create retro`, async () => {
        const retroName = 'Test API Create Retro'
        const brainstormVisibility = 'visible'
        const maxVotes = 3
        const format = 'worked_improve_question'

        const b = await apiContext.post(`users/${user.id}/retros`, {
            data: {
                retroName,
                brainstormVisibility,
                maxVotes,
                format,
            },
        })
        expect(b.ok()).toBeTruthy()
        const retro = await b.json()
        expect(retro.data).toMatchObject({
            name: retroName,
            brainstormVisibility,
            format,
        })
    })

    test(`GET /users/{userId}/retros should return object in array when retros associated to user`, async () => {
        const retroName = 'Test API Retros'
        const brainstormVisibility = 'hidden'
        const maxVotes = 3
        const format = 'worked_improve_question'

        const b = await apiContext.post(`users/${user.id}/retros`, {
            data: {
                retroName,
                brainstormVisibility,
                maxVotes,
                format,
            },
        })
        expect(b.ok()).toBeTruthy()

        const bs = await apiContext.get(`users/${user.id}/retros`)
        expect(bs.ok()).toBeTruthy()
        const retros = await bs.json()
        expect(retros.data).toContainEqual(
            expect.objectContaining({
                name: retroName,
            }),
        )
    })

    test(`POST /teams/{teamId}/users/{userId}/retros should create retro`, async () => {
        const retroName = 'Test API Create Team Retro'
        const brainstormVisibility = 'hidden'
        const maxVotes = 3
        const format = 'worked_improve_question'

        const t = await apiContext.post(`users/${user.id}/teams`, {
            data: {
                name: 'test team create retro',
            },
        })
        const { data: team } = await t.json()

        const b = await apiContext.post(
            `teams/${team.id}/users/${user.id}/retros`,
            {
                data: {
                    retroName,
                    brainstormVisibility,
                    maxVotes,
                    format,
                },
            },
        )
        expect(b.ok()).toBeTruthy()
        const retro = await b.json()
        expect(retro.data).toMatchObject({
            name: retroName,
        })

        const bs = await apiContext.get(`teams/${team.id}/retros`)
        expect(bs.ok()).toBeTruthy()
        const retros = await bs.json()
        expect(retros.data).toContainEqual(
            expect.objectContaining({
                name: retroName,
            }),
        )
    })
})
