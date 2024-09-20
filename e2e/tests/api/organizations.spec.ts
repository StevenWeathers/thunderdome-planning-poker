import { expect, test } from "@fixtures/test-setup";

test.describe("Organization API", { tag: ["@api", "@organization"] }, () => {
  test.describe("GET /users/{userId}/organizations", () => {
    test("returns empty array when no organizations associated to user", async ({
      request,
      adminApiUser,
    }) => {
      const response = await adminApiUser.context.get(
        `users/${adminApiUser.user.id}/organizations`,
      );
      expect(response.ok()).toBeTruthy();
      const organizations = await response.json();
      expect(organizations.data).toEqual([]);
    });

    test("returns object in array when organizations associated to user", async ({
      request,
      orgTeamApiUser,
    }) => {
      const response = await orgTeamApiUser.context.get(
        `users/${orgTeamApiUser.user.id}/organizations`,
      );
      expect(response.ok()).toBeTruthy();
      const organizations = await response.json();
      expect(organizations.data).toContainEqual(
        expect.objectContaining({
          id: orgTeamApiUser.user.orgs[0].id,
          name: orgTeamApiUser.user.orgs[0].name,
        }),
      );
    });
  });

  test.describe("POST /users/{userId}/organizations", () => {
    test("creates organization", async ({ request, registeredApiUser }) => {
      const organizationName = "Test API Create Organization";
      const response = await registeredApiUser.context.post(
        `users/${registeredApiUser.user.id}/organizations`,
        {
          data: { name: organizationName },
        },
      );
      expect(response.ok()).toBeTruthy();
      const organization = await response.json();
      expect(organization.data).toMatchObject({ name: organizationName });
    });
  });

  test.describe("GET /api/organizations/{orgId}", () => {
    test("returns 200 and organization data for org member", async ({
      request,
      orgTeamApiUser,
    }) => {
      const org = orgTeamApiUser.user.orgs[0];
      const response = await orgTeamApiUser.context.get(
        `organizations/${org.id}`,
      );
      expect(response.ok()).toBeTruthy();
      expect(response.status()).toBe(200);
      const orgData = await response.json();
      expect(orgData.data.organization).toMatchObject({
        id: org.id,
        name: org.name,
      });
    });

    test(
      "returns 200 and organization data for global admin",
      { tag: ["@admin"] },
      async ({ request, adminApiUser, orgTeamApiUser }) => {
        const org = orgTeamApiUser.user.orgs[0];
        const response = await adminApiUser.context.get(
          `organizations/${org.id}`,
        );
        expect(response.ok()).toBeTruthy();
        expect(response.status()).toBe(200);
        const orgData = await response.json();
        expect(orgData.data.organization).toMatchObject({
          id: org.id,
          name: org.name,
        });
      },
    );

    test("returns 403 Forbidden for non-org member", async ({
      request,
      registeredApiUser,
      orgTeamApiUser,
    }) => {
      const org = orgTeamApiUser.user.orgs[0];
      const response = await registeredApiUser.context.get(
        `organizations/${org.id}`,
      );
      expect(response.ok()).toBeFalsy();
      expect(response.status()).toBe(403);
    });
  });

  test.describe("DELETE /api/organizations/{orgId}", () => {
    test("returns 200 for org admin", async ({
      request,
      orgOwnerApiUser,
      testDatabase,
    }) => {
      const org = await testDatabase.seeder.createOrganization(
        "orgToDeleteOwner",
        orgOwnerApiUser.user.id,
      );

      const response = await orgOwnerApiUser.context.delete(
        `organizations/${org.id}`,
      );
      expect(response.ok()).toBeTruthy();
      expect(response.status()).toBe(200);

      const confirmOrgDeleted = await testDatabase.seeder.getOrganizationById(
        org.id,
      );
      expect(confirmOrgDeleted).toBeNull();
    });

    test(
      "returns 200 for global admin",
      { tag: ["@admin"] },
      async ({ request, orgOwnerApiUser, adminApiUser, testDatabase }) => {
        const org = await testDatabase.seeder.createOrganization(
          "orgToDeleteGlobalAdmin",
          orgOwnerApiUser.user.id,
        );

        const response = await adminApiUser.context.delete(
          `organizations/${org.id}`,
        );
        expect(response.ok()).toBeTruthy();
        expect(response.status()).toBe(200);

        const confirmOrgDeleted = await testDatabase.seeder.getOrganizationById(
          org.id,
        );
        expect(confirmOrgDeleted).toBeNull();
      },
    );

    test("returns 403 Forbidden for non org admin", async ({
      request,
      orgTeamApiUser,
    }) => {
      const org = orgTeamApiUser.user.orgs[0];
      const response = await orgTeamApiUser.context.delete(
        `organizations/${org.id}`,
      );
      expect(response.ok()).toBeFalsy();
      expect(response.status()).toBe(403);
    });
  });

  test.describe("GET /organizations/{orgId}/teams", () => {
    test("Org Admin can list all teams in their organization", async ({
      orgAdminApiUser,
      orgOwnerApiUser,
    }) => {
      const org = orgAdminApiUser.user.orgs[0];
      const response = await orgAdminApiUser.context.get(
        `organizations/${org.id}/teams`,
      );
      expect(response.ok()).toBeTruthy();
      expect(response.status()).toBe(200);
      const teams = await response.json();
      expect(Array.isArray(teams.data)).toBeTruthy();
      expect(teams.data.length).toBeGreaterThan(0);

      expect(teams.data).toContainEqual(
        expect.objectContaining({
          id: orgOwnerApiUser.user.orgTeams[0].id,
          name: orgOwnerApiUser.user.orgTeams[0].name,
        }),
      );
    });
  });
});
