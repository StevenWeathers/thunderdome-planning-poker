import {expect, test} from '@playwright/test'

const userProfileEndpoint = `auth/user`;

// Request context is reused by all tests in the file.
let apiContext;

test.beforeAll(async ({playwright}) => {
    apiContext = await playwright.request.newContext({
        baseURL: 'http://localhost:8080/api/',
        extraHTTPHeaders: {
            'X-API-Key': `8MenPkY8.Vqvkh030vv7$rSyYs1gt++L0v7wKuVgR`,
        },
    });
})

test.afterAll(async ({}) => {
    // Dispose all responses.
    await apiContext.dispose();
});

test.describe('registered user', () => {
    test(`${userProfileEndpoint} should return session user profile`, async () => {
        const u = await apiContext.get(userProfileEndpoint);
        expect(u.ok()).toBeTruthy();

        const user = await u.json()
        expect(user.data).toMatchObject({
            name: 'E2E API User',
            email: 'e2eapi@thunderdome.dev',
            rank: 'REGISTERED',
            verified: true,
            disabled: false,
            mfaEnabled: false
        });
    });
});