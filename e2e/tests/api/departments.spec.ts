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
    test(`GET /organizations/{orgId}/departments should return empty array when no departments associated to user`, async () => {
        const o = await adminApiContext.post(`users/${user.id}/organizations`, {
            data: {
                name: 'Test API Create Organization',
            },
        })
        const organization = await o.json()
        const b = await adminApiContext.get(
            `organizations/${organization.data.id}/departments`,
        )
        expect(b.ok()).toBeTruthy()

        const departments = await b.json()
        expect(departments.data).toMatchObject([])
    })

    test(`POST /organizations/{orgId}/departments should create department`, async () => {
        const o = await apiContext.post(`users/${user.id}/organizations`, {
            data: {
                name: 'Test API Create Organization',
            },
        })
        const organization = await o.json()
        const departmentName = 'Test API Create Department'

        const b = await apiContext.post(
            `organizations/${organization.data.id}/departments`,
            {
                data: {
                    name: departmentName,
                },
            },
        )
        expect(b.ok()).toBeTruthy()
        const department = await b.json()
        expect(department.data).toMatchObject({
            name: departmentName,
        })
    })

    test(`GET /organizations/{orgId}/departments should return object in array when departments associated to user`, async () => {
        const o = await apiContext.post(`users/${user.id}/organizations`, {
            data: {
                name: 'Test API Create Organization',
            },
        })
        const organization = await o.json()

        const departmentName = 'Test API Departments'
        const b = await apiContext.post(
            `organizations/${organization.data.id}/departments`,
            {
                data: {
                    name: departmentName,
                },
            },
        )
        expect(b.ok()).toBeTruthy()

        const bs = await apiContext.get(
            `organizations/${organization.data.id}/departments`,
        )
        expect(bs.ok()).toBeTruthy()
        const departments = await bs.json()
        expect(departments.data).toContainEqual(
            expect.objectContaining({
                name: departmentName,
            }),
        )
    })
})
