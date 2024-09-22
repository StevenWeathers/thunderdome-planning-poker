import { expect, test } from "@fixtures/test-setup";

test.describe(
  "Organization Department Team API",
  { tag: ["@api", "@department", "@organization"] },
  () => {
    test.describe("GET /teams/{teamId}", () => {
      test("Department Admin can view any team in their department", async ({
        deptAdminApiUser,
        deptTeamApiUser,
      }) => {
        const departmentTeam = deptTeamApiUser.user.deptTeams[0];
        const response = await deptAdminApiUser.context.get(
          `teams/${departmentTeam.id}`,
        );
        expect(response.ok()).toBeTruthy();
        expect(response.status()).toBe(200);
        const tr = await response.json();
        expect(tr.data.team.id).toBe(departmentTeam.id);
      });

      test("Department Team Admin can view department team details", async ({
        deptTeamApiUser,
        deptTeamAdminApiUser,
      }) => {
        const departmentTeam = deptTeamApiUser.user.deptTeams[0];
        const getResponse = await deptTeamAdminApiUser.context.get(
          `teams/${departmentTeam.id}`,
        );
        expect(getResponse.ok()).toBeTruthy();
        expect(getResponse.status()).toBe(200);
        const tr = await getResponse.json();
        expect(tr.data).toMatchObject({
          team: { id: departmentTeam.id, name: departmentTeam.name },
        });
      });

      test("Department Team Member can view department team details", async ({
        deptTeamApiUser,
      }) => {
        const departmentTeam = deptTeamApiUser.user.deptTeams[0];
        const response = await deptTeamApiUser.context.get(
          `teams/${departmentTeam.id}`,
        );
        expect(response.ok()).toBeTruthy();
        expect(response.status()).toBe(200);
        const tr = await response.json();
        expect(tr.data.team.id).toBe(departmentTeam.id);
      });

      test("Global Admin can view any department team", async ({
        adminApiUser,
        deptTeamApiUser,
      }) => {
        const departmentTeam = deptTeamApiUser.user.deptTeams[0];
        const response = await adminApiUser.context.get(
          `teams/${departmentTeam.id}`,
        );
        expect(response.ok()).toBeTruthy();
        expect(response.status()).toBe(200);
        const tr = await response.json();
        expect(tr.data.team.id).toBe(departmentTeam.id);
      });

      test("Non-department team member cannot view department team details", async ({
        registeredApiUser,
        deptTeamApiUser,
      }) => {
        const departmentTeam = deptTeamApiUser.user.deptTeams[0];
        const response = await registeredApiUser.context.get(
          `teams/${departmentTeam.id}`,
        );
        expect(response.ok()).toBeFalsy();
        expect(response.status()).toBe(403);
      });
    });

    test.describe("PUT /teams/{teamId}", () => {
      test("Department Admin can update any team in their department", async ({
        testDatabase,
        deptAdminApiUser,
        orgOwnerApiUser,
      }) => {
        const dept = orgOwnerApiUser.user.depts[0];
        const departmentTeam = await testDatabase.seeder.createOrgTeam(
          "departmentTeamToUpdateDeptAdmin",
          orgOwnerApiUser.user.orgs[0].id,
          dept.id,
        );
        const updatedName = "Dept Admin Updated Team";
        const response = await deptAdminApiUser.context.put(
          `teams/${departmentTeam.id}`,
          {
            data: { name: updatedName },
          },
        );
        expect(response.ok()).toBeTruthy();
        expect(response.status()).toBe(200);
        const updatedDepartmentTeam = await response.json();
        expect(updatedDepartmentTeam.data.name).toBe(updatedName);
      });

      test("Department Admin cannot update team outside their department", async ({
        testDatabase,
        deptAdminApiUser,
        orgOwnerApiUser,
      }) => {
        const otherDept = await testDatabase.seeder.createDepartment(
          "Other Department",
          orgOwnerApiUser.user.orgs[0].id,
        );
        const otherDeptTeam = await testDatabase.seeder.createOrgTeam(
          "OtherDeptTeam",
          orgOwnerApiUser.user.orgs[0].id,
          otherDept.id,
        );
        const updatedName = "Unauthorized Update";
        const response = await deptAdminApiUser.context.put(
          `teams/${otherDeptTeam.id}`,
          {
            data: { name: updatedName },
          },
        );
        expect(response.ok()).toBeFalsy();
        expect(response.status()).toBe(403);
      });

      test("Department Team Member cannot update department team details", async ({
        deptTeamApiUser,
      }) => {
        const updatedName = "Unauthorized Update";
        const response = await deptTeamApiUser.context.put(
          `teams/${deptTeamApiUser.user.deptTeams[0].id}`,
          {
            data: { name: updatedName },
          },
        );
        expect(response.ok()).toBeFalsy();
        expect(response.status()).toBe(403);
      });

      test("Department Team Admin can update department team details", async ({
        testDatabase,
        deptTeamAdminApiUser,
        orgOwnerApiUser,
      }) => {
        const dept = orgOwnerApiUser.user.depts[0];
        const departmentTeam = await testDatabase.seeder.createOrgTeam(
          "departmentTeamToUpdateTeamAdmin",
          orgOwnerApiUser.user.orgs[0].id,
          dept.id,
        );
        await testDatabase.seeder.addUserToTeam(
          deptTeamAdminApiUser.user.id,
          departmentTeam.id,
          "ADMIN",
        );
        const updatedName = "Admin Updated Department Team";
        const response = await deptTeamAdminApiUser.context.put(
          `teams/${departmentTeam.id}`,
          {
            data: { name: updatedName },
          },
        );
        // expect(response.ok()).toBeTruthy();
        expect(response.status()).toBe(200);
        const updatedDepartmentTeam = await response.json();
        expect(updatedDepartmentTeam.data.name).toBe(updatedName);
      });

      test("Global Admin can update any department team", async ({
        testDatabase,
        adminApiUser,
        orgOwnerApiUser,
      }) => {
        const dept = orgOwnerApiUser.user.depts[0];
        const departmentTeam = await testDatabase.seeder.createOrgTeam(
          "departmentTeamToUpdateGlobalAdmin",
          orgOwnerApiUser.user.orgs[0].id,
          dept.id,
        );
        const updatedName = "Global Admin Updated Department Team";
        const response = await adminApiUser.context.put(
          `teams/${departmentTeam.id}`,
          {
            data: { name: updatedName },
          },
        );
        expect(response.ok()).toBeTruthy();
        expect(response.status()).toBe(200);
        const updatedDepartmentTeam = await response.json();
        expect(updatedDepartmentTeam.data.name).toBe(updatedName);
      });

      test("Non-department team member cannot update department team", async ({
        deptTeamApiUser,
        registeredApiUser,
      }) => {
        const updatedName = "Unauthorized Update";
        const response = await registeredApiUser.context.put(
          `teams/${deptTeamApiUser.user.deptTeams[0].id}`,
          {
            data: { name: updatedName },
          },
        );
        expect(response.ok()).toBeFalsy();
        expect(response.status()).toBe(403);
      });
    });

    test.describe("DELETE /teams/{teamId}", () => {
      test("Department Admin can delete any team in their department", async ({
        testDatabase,
        deptAdminApiUser,
        orgOwnerApiUser,
      }) => {
        const dept = orgOwnerApiUser.user.depts[0];
        const departmentTeam = await testDatabase.seeder.createOrgTeam(
          "departmentTeamToDeleteDeptAdmin",
          orgOwnerApiUser.user.orgs[0].id,
          dept.id,
        );
        const response = await deptAdminApiUser.context.delete(
          `teams/${departmentTeam.id}`,
        );
        expect(response.ok()).toBeTruthy();
        expect(response.status()).toBe(200);
      });

      test("Department Admin cannot delete team outside their department", async ({
        testDatabase,
        deptAdminApiUser,
        orgOwnerApiUser,
      }) => {
        const otherDept = await testDatabase.seeder.createDepartment(
          "Other Department2",
          orgOwnerApiUser.user.orgs[0].id,
        );
        const otherDeptTeam = await testDatabase.seeder.createOrgTeam(
          "OtherDeptTeam2",
          orgOwnerApiUser.user.orgs[0].id,
          otherDept.id,
        );
        const response = await deptAdminApiUser.context.delete(
          `teams/${otherDeptTeam.id}`,
        );
        expect(response.ok()).toBeFalsy();
        expect(response.status()).toBe(403);
      });

      test("Department Team Admin can delete department team", async ({
        testDatabase,
        orgOwnerApiUser,
        deptTeamAdminApiUser,
      }) => {
        const dept = orgOwnerApiUser.user.depts[0];
        const departmentTeam = await testDatabase.seeder.createOrgTeam(
          "departmentTeamToDeleteTeamAdmin",
          orgOwnerApiUser.user.orgs[0].id,
          dept.id,
        );
        await testDatabase.seeder.addUserToTeam(
          deptTeamAdminApiUser.user.id,
          departmentTeam.id,
          "ADMIN",
        );
        const response = await deptTeamAdminApiUser.context.delete(
          `teams/${departmentTeam.id}`,
        );
        expect(response.ok()).toBeTruthy();
        expect(response.status()).toBe(200);
      });

      test("Global Admin can delete any department team", async ({
        testDatabase,
        adminApiUser,
        orgOwnerApiUser,
      }) => {
        const dept = orgOwnerApiUser.user.depts[0];
        const departmentTeam = await testDatabase.seeder.createOrgTeam(
          "departmentTeamToUpdateTeamAdmin",
          orgOwnerApiUser.user.orgs[0].id,
          dept.id,
        );
        const response = await adminApiUser.context.delete(
          `teams/${departmentTeam.id}`,
        );
        expect(response.ok()).toBeTruthy();
        expect(response.status()).toBe(200);
      });

      test("Department Team Member cannot delete department team", async ({
        deptTeamApiUser,
      }) => {
        const response = await deptTeamApiUser.context.delete(
          `teams/${deptTeamApiUser.user.deptTeams[0].id}`,
        );
        expect(response.ok()).toBeFalsy();
        expect(response.status()).toBe(403);
      });

      test("Non-department team member cannot delete department team", async ({
        deptTeamApiUser,
        registeredApiUser,
      }) => {
        const response = await registeredApiUser.context.delete(
          `teams/${deptTeamApiUser.user.deptTeams[0].id}`,
        );
        expect(response.ok()).toBeFalsy();
        expect(response.status()).toBe(403);
      });
    });

    test.describe("POST /departments/{departmentId}/teams", () => {
      test("Department Admin can create a new team in their department", async ({
        deptAdminApiUser,
        orgOwnerApiUser,
      }) => {
        const dept = orgOwnerApiUser.user.depts[0];
        const newTeamName = "New Team Created by Dept Admin";
        const response = await deptAdminApiUser.context.post(
          `organizations/${orgOwnerApiUser.user.orgs[0].id}/departments/${dept.id}/teams`,
          {
            data: { name: newTeamName },
          },
        );
        expect(response.ok()).toBeTruthy();
        expect(response.status()).toBe(200);
        const newTeam = await response.json();
        expect(newTeam.data.name).toBe(newTeamName);
      });
    });

    test.describe("GET /departments/{departmentId}/teams", () => {
      test("Department Admin can list all teams in their department", async ({
        deptAdminApiUser,
        orgOwnerApiUser,
      }) => {
        const dept = orgOwnerApiUser.user.depts[0];
        const response = await deptAdminApiUser.context.get(
          `organizations/${orgOwnerApiUser.user.orgs[0].id}/departments/${dept.id}/teams`,
        );
        expect(response.ok()).toBeTruthy();
        expect(response.status()).toBe(200);
        const teams = await response.json();
        expect(Array.isArray(teams.data)).toBeTruthy();
        expect(teams.data.length).toBeGreaterThan(0);

        expect(teams.data).toEqual(
          expect.arrayContaining([
            expect.objectContaining({
              id: orgOwnerApiUser.user.deptTeams[0].id,
            }),
          ]),
        );
      });
    });
  },
);
