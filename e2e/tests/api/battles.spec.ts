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
    test(`/users/{userId}/battles should return empty array when no battles associated to user`, async () => {
        const u = await apiContext.get(userProfileEndpoint);
        const user = await u.json()

        const b = await apiContext.get(`users/${user.data.id}/battles`);
        console.log(b)
        expect(b.ok()).toBeTruthy();
        const battles = await b.json()
        expect(battles.data).toMatchObject([]);
    });
});