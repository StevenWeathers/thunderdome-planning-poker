import { expect, test } from "@fixtures/test-setup";

test.describe(
  "User Profile API",
  { tag: ["@api", "@user", "@registered"] },
  () => {
    const userProfileEndpoint = `auth/user`;

    test("GET /auth/user returns session user profile for registered user", async ({
      request,
      registeredApiUser,
    }) => {
      const response = await registeredApiUser.context.get(userProfileEndpoint);
      expect(response.ok()).toBeTruthy();
      expect(response.status()).toBe(200);

      const userProfile = await response.json();
      expect(userProfile.data).toMatchObject({
        id: registeredApiUser.user.id,
        name: "E2EAPIUser",
        email: "e2eapi@thunderdome.dev",
        rank: "REGISTERED",
        verified: true,
        disabled: false,
      });
    });

    test("GET /auth/user returns session user profile for admin user", async ({
      request,
      adminApiUser,
    }) => {
      const response = await adminApiUser.context.get(userProfileEndpoint);
      expect(response.ok()).toBeTruthy();
      expect(response.status()).toBe(200);

      const userProfile = await response.json();
      expect(userProfile.data).toMatchObject({
        id: adminApiUser.user.id,
        rank: "ADMIN",
        verified: true,
        disabled: false,
      });
    });

    test("GET /auth/user returns 401 for unauthenticated request", async ({
      request,
    }) => {
      const response = await request.get(`api/${userProfileEndpoint}`, {
        headers: {
          "X-API-Key": "invalid_api_key",
        },
      });
      expect(response.ok()).toBeFalsy();
      expect(response.status()).toBe(401);
    });
  },
);
