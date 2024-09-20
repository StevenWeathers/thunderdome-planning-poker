import { expect, test } from "@fixtures/test-setup";

test.describe(
  "Organization Team API",
  { tag: ["@api", "@team", "@organization"] },
  () => {
    test.describe("POST /organizations/{orgId}/teams", () => {
      test("Org Admin can create a new team in any department of their organization", async ({
        orgAdminApiUser,
        testDatabase,
      }) => {
        const org = orgAdminApiUser.user.orgs[0];
        const dept = await testDatabase.seeder.createDepartment(
          "Dept for New Team",
          org.id,
        );
        const newTeamName = "New Team Created by Org Admin";
        const response = await orgAdminApiUser.context.post(
          `organizations/${org.id}/teams`,
          {
            data: { name: newTeamName, department_id: dept.id },
          },
        );
        expect(response.ok()).toBeTruthy();
        expect(response.status()).toBe(200);
        const newTeam = await response.json();
        expect(newTeam.data.name).toBe(newTeamName);
      });
    });

    test.describe("GET /teams/{teamId}", () => {
      test("Team Admin can view team details", async ({
        orgTeamApiUser,
        orgTeamAdminApiUser,
      }) => {
        const team = orgTeamApiUser.user.orgTeams[0];
        const getResponse = await orgTeamAdminApiUser.context.get(
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
        const team = orgTeamApiUser.user.orgTeams[0];
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
        const team = orgTeamApiUser.user.orgTeams[0];
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
        const team = orgTeamApiUser.user.orgTeams[0];
        const response = await registeredApiUser.context.get(
          `teams/${team.id}`,
        );
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
          `teams/${orgTeamApiUser.user.orgTeams[0].id}`,
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
        orgOwnerApiUser,
      }) => {
        const team = await testDatabase.seeder.createOrgTeam(
          "teamToUpdateTeamAdmin",
          orgOwnerApiUser.user.orgs[1],
        );
        await testDatabase.seeder.addUserToTeam(
          teamAdminApiUser.user.id,
          team.id,
          "ADMIN",
        );
        const updatedName = "Admin Updated Team";
        const response = await teamAdminApiUser.context.put(
          `teams/${team.id}`,
          {
            data: { name: updatedName },
          },
        );
        expect(response.ok()).toBeTruthy();
        expect(response.status()).toBe(200);
        const updatedTeam = await response.json();
        expect(updatedTeam.data.name).toBe(updatedName);
      });

      test("Org Admin can update any team in their organization", async ({
        orgAdminApiUser,
        testDatabase,
      }) => {
        const org = orgAdminApiUser.user.orgs[0];
        const dept = await testDatabase.seeder.createDepartment(
          "Dept for Team Update",
          org.id,
        );
        const team = await testDatabase.seeder.createOrgTeam(
          "Team to Update",
          org.id,
          dept.id,
        );
        const updatedName = "Updated by Org Admin";
        const response = await orgAdminApiUser.context.put(`teams/${team.id}`, {
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
        orgOwnerApiUser,
      }) => {
        const team = await testDatabase.seeder.createOrgTeam(
          "teamToUpdateGlobalAdmin",
          orgOwnerApiUser.user.orgs[1],
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
          `teams/${orgTeamApiUser.user.orgTeams[0].id}`,
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
        orgOwnerApiUser,
        teamAdminApiUser,
      }) => {
        const team = await testDatabase.seeder.createOrgTeam(
          "teamToDeleteTeamAdmin",
          orgOwnerApiUser.user.orgs[0].id,
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

      test("Org Admin can delete any team in their organization", async ({
        orgAdminApiUser,
        testDatabase,
      }) => {
        const org = orgAdminApiUser.user.orgs[0];
        const dept = await testDatabase.seeder.createDepartment(
          "Dept for Team Deletion",
          org.id,
        );
        const team = await testDatabase.seeder.createOrgTeam(
          "Team to Delete",
          org.id,
          dept.id,
        );
        const response = await orgAdminApiUser.context.delete(
          `teams/${team.id}`,
        );
        expect(response.ok()).toBeTruthy();
        expect(response.status()).toBe(200);
      });

      test("Global Admin can delete any team", async ({
        testDatabase,
        adminApiUser,
        orgOwnerApiUser,
      }) => {
        const team = await testDatabase.seeder.createOrgTeam(
          "teamToDeleteGlobalAdmin",
          orgOwnerApiUser.user.orgs[0].id,
        );
        const response = await adminApiUser.context.delete(`teams/${team.id}`);
        expect(response.ok()).toBeTruthy();
        expect(response.status()).toBe(200);
      });

      test("Team Member cannot delete team", async ({ orgTeamApiUser }) => {
        const response = await orgTeamApiUser.context.delete(
          `teams/${orgTeamApiUser.user.orgTeams[0].id}`,
        );
        expect(response.ok()).toBeFalsy();
        expect(response.status()).toBe(403);
      });

      test("Non-team member cannot delete team", async ({
        orgTeamApiUser,
        registeredApiUser,
      }) => {
        const response = await registeredApiUser.context.delete(
          `teams/${orgTeamApiUser.user.orgTeams[0].id}`,
        );
        expect(response.ok()).toBeFalsy();
        expect(response.status()).toBe(403);
      });
    });
  },
);
