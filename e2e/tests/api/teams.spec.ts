import { expect, test } from "@fixtures/test-setup";

test.describe("Team API", { tag: ["@api", "@team"] }, () => {
  test.describe("GET /users/{userId}/teams", () => {
    test("returns empty array when no teams associated to user for Entity User", async ({
      orgAdminApiUser,
    }) => {
      const response = await orgAdminApiUser.context.get(
        `users/${orgAdminApiUser.user.id}/teams`,
      );
      expect(response.ok()).toBeTruthy();
      expect(response.status()).toBe(200);
      const teams = await response.json();
      expect(teams.data).toEqual([]);
    });

    test("returns empty array when no teams associated to user for Global Admin", async ({
      orgAdminApiUser,
      adminApiUser,
    }) => {
      const response = await adminApiUser.context.get(
        `users/${orgAdminApiUser.user.id}/teams`,
      );
      expect(response.ok()).toBeTruthy();
      expect(response.status()).toBe(200);
      const teams = await response.json();
      expect(teams.data).toEqual([]);
    });

    test("returns object in array when teams associated to user for Entity User", async ({
      orgOwnerApiUser,
    }) => {
      const team = orgOwnerApiUser.user.teams[0];
      const response = await orgOwnerApiUser.context.get(
        `users/${orgOwnerApiUser.user.id}/teams`,
      );
      expect(response.ok()).toBeTruthy();
      expect(response.status()).toBe(200);
      const teams = await response.json();
      expect(teams.data).toContainEqual(
        expect.objectContaining({
          id: team.id,
          name: team.name,
        }),
      );
    });

    test("returns object in array when teams associated to user for Global Admin", async ({
      adminApiUser,
      orgOwnerApiUser,
    }) => {
      const team = orgOwnerApiUser.user.teams[0];
      const response = await adminApiUser.context.get(
        `users/${orgOwnerApiUser.user.id}/teams`,
      );
      expect(response.ok()).toBeTruthy();
      expect(response.status()).toBe(200);
      const teams = await response.json();
      expect(teams.data).toContainEqual(
        expect.objectContaining({
          id: team.id,
          name: team.name,
        }),
      );
    });

    test("returns forbidden for Non Entity User", async ({
      teamAdminApiUser,
      registeredApiUser,
    }) => {
      const response = await registeredApiUser.context.get(
        `users/${teamAdminApiUser.user.id}/teams`,
      );
      expect(response.ok()).toBeFalsy();
      expect(response.status()).toBe(403);
    });
  });

  test.describe("POST /users/{userId}/teams", () => {
    test("creates team", async ({ registeredApiUser }) => {
      const teamName = "Test API Create Team";
      const response = await registeredApiUser.context.post(
        `users/${registeredApiUser.user.id}/teams`,
        {
          data: { name: teamName },
        },
      );
      expect(response.ok()).toBeTruthy();
      expect(response.status()).toBe(200);
      const team = await response.json();
      expect(team.data).toMatchObject({ name: teamName });
    });
  });

  test.describe("GET /teams/{teamId}", () => {
    test("Team Admin can view team details", async ({ teamAdminApiUser }) => {
      const team = teamAdminApiUser.user.teams[0];
      const getResponse = await teamAdminApiUser.context.get(
        `teams/${team.id}`,
      );
      expect(getResponse.ok()).toBeTruthy();
      expect(getResponse.status()).toBe(200);
      const tr = await getResponse.json();
      expect(tr.data).toMatchObject({
        team: { id: team.id, name: team.name },
      });
    });

    test("Team Member can view team details", async ({ orgTeamApiUser }) => {
      const team = orgTeamApiUser.user.teams[0];
      const response = await orgTeamApiUser.context.get(`teams/${team.id}`);
      expect(response.ok()).toBeTruthy();
      expect(response.status()).toBe(200);
      const tr = await response.json();
      expect(tr.data.team.id).toBe(team.id);
    });

    test("Global Admin can view any team", async ({
      adminApiUser,
      orgTeamApiUser,
    }) => {
      const team = orgTeamApiUser.user.teams[0];
      const response = await adminApiUser.context.get(`teams/${team.id}`);
      expect(response.ok()).toBeTruthy();
      expect(response.status()).toBe(200);
      const tr = await response.json();
      expect(tr.data.team.id).toBe(team.id);
    });

    test("Non-team member cannot view team details", async ({
      registeredApiUser,
      orgTeamApiUser,
    }) => {
      const team = orgTeamApiUser.user.teams[0];
      const response = await registeredApiUser.context.get(`teams/${team.id}`);
      expect(response.ok()).toBeFalsy();
      expect(response.status()).toBe(403);
    });
  });

  test.describe("PUT /teams/{teamId}", () => {
    test("Team Member cannot update team details", async ({
      orgTeamApiUser,
    }) => {
      const updatedName = "Unauthorized Update";
      const response = await orgTeamApiUser.context.put(
        `teams/${orgTeamApiUser.user.teams[0].id}`,
        {
          data: { name: updatedName },
        },
      );
      expect(response.ok()).toBeFalsy();
      expect(response.status()).toBe(403);
    });

    test("Team Admin can update team details", async ({
      testDatabase,
      teamAdminApiUser,
      registeredApiUser,
    }) => {
      const team = await testDatabase.seeder.createUserTeam(
        "teamToUpdate",
        registeredApiUser.user.id,
      );
      await testDatabase.seeder.addUserToTeam(
        teamAdminApiUser.user.id,
        team.id,
        "ADMIN",
      );
      const updatedName = "Admin Updated Team";
      const response = await teamAdminApiUser.context.put(`teams/${team.id}`, {
        data: { name: updatedName },
      });
      expect(response.ok()).toBeTruthy();
      expect(response.status()).toBe(200);
      const updatedTeam = await response.json();
      expect(updatedTeam.data.name).toBe(updatedName);
    });

    test("Global Admin can update any team", async ({
      testDatabase,
      adminApiUser,
      registeredApiUser,
    }) => {
      const team = await testDatabase.seeder.createUserTeam(
        "teamToUpdate",
        registeredApiUser.user.id,
      );
      const updatedName = "Global Admin Updated Team";
      const response = await adminApiUser.context.put(`teams/${team.id}`, {
        data: { name: updatedName },
      });
      expect(response.ok()).toBeTruthy();
      expect(response.status()).toBe(200);
      const updatedTeam = await response.json();
      expect(updatedTeam.data.name).toBe(updatedName);
    });

    test("Non-team member cannot update team", async ({
      orgTeamApiUser,
      registeredApiUser,
    }) => {
      const updatedName = "Unauthorized Update";
      const response = await registeredApiUser.context.put(
        `teams/${orgTeamApiUser.user.teams[0].id}`,
        {
          data: { name: updatedName },
        },
      );
      expect(response.ok()).toBeFalsy();
      expect(response.status()).toBe(403);
    });
  });

  test.describe("DELETE /teams/{teamId}", () => {
    test("Team Admin can delete team", async ({
      testDatabase,
      registeredApiUser,
      teamAdminApiUser,
    }) => {
      const team = await testDatabase.seeder.createUserTeam(
        "teamToUpdate",
        registeredApiUser.user.id,
      );
      await testDatabase.seeder.addUserToTeam(
        teamAdminApiUser.user.id,
        team.id,
        "ADMIN",
      );
      const response = await teamAdminApiUser.context.delete(
        `teams/${team.id}`,
      );
      expect(response.ok()).toBeTruthy();
      expect(response.status()).toBe(200);
    });

    test("Global Admin can delete any team", async ({
      testDatabase,
      adminApiUser,
      registeredApiUser,
    }) => {
      const team = await testDatabase.seeder.createUserTeam(
        "teamToUpdate",
        registeredApiUser.user.id,
      );
      const response = await adminApiUser.context.delete(`teams/${team.id}`);
      expect(response.ok()).toBeTruthy();
      expect(response.status()).toBe(200);
    });

    test("Team Member cannot delete team", async ({ orgTeamApiUser }) => {
      const response = await orgTeamApiUser.context.delete(
        `teams/${orgTeamApiUser.user.teams[0].id}`,
      );
      expect(response.ok()).toBeFalsy();
      expect(response.status()).toBe(403);
    });

    test("Non-team member cannot delete team", async ({
      orgTeamApiUser,
      registeredApiUser,
    }) => {
      const response = await registeredApiUser.context.delete(
        `teams/${orgTeamApiUser.user.teams[0].id}`,
      );
      expect(response.ok()).toBeFalsy();
      expect(response.status()).toBe(403);
    });
  });
});
