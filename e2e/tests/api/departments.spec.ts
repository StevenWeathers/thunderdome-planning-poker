import { expect, test } from "@fixtures/test-setup";

test.describe(
  "Organization Department API",
  { tag: ["@api", "@department"] },
  () => {
    test.describe("GET /organizations/{orgId}/departments", () => {
      test("returns departments list for org member", async ({
        request,
        orgTeamApiUser,
      }) => {
        const org = orgTeamApiUser.user.orgs[0];
        const response = await orgTeamApiUser.context.get(
          `organizations/${org.id}/departments`,
        );
        expect(response.ok()).toBeTruthy();
        expect(response.status()).toBe(200);
        const departments = await response.json();
        expect(Array.isArray(departments.data)).toBeTruthy();
        expect(departments.data.length).toBeGreaterThan(0);
        expect(departments.data[0]).toHaveProperty("id");
        expect(departments.data[0]).toHaveProperty("name");
      });

      test("returns departments list for org admin", async ({
        request,
        orgOwnerApiUser,
      }) => {
        const org = orgOwnerApiUser.user.orgs[0];
        const response = await orgOwnerApiUser.context.get(
          `organizations/${org.id}/departments`,
        );
        expect(response.ok()).toBeTruthy();
        expect(response.status()).toBe(200);
        const departments = await response.json();
        expect(Array.isArray(departments.data)).toBeTruthy();
        // Assuming org admin can see all departments, even if there are none
        expect(departments.data.length).toBeGreaterThanOrEqual(0);
        if (departments.data.length > 0) {
          expect(departments.data[0]).toHaveProperty("id");
          expect(departments.data[0]).toHaveProperty("name");
        }
      });

      test("returns departments list for global admin", async ({
        request,
        adminApiUser,
        orgTeamApiUser,
      }) => {
        const org = orgTeamApiUser.user.orgs[0]; // Using an org the admin didn't create
        const response = await adminApiUser.context.get(
          `organizations/${org.id}/departments`,
        );
        expect(response.ok()).toBeTruthy();
        expect(response.status()).toBe(200);
        const departments = await response.json();
        expect(Array.isArray(departments.data)).toBeTruthy();
        // Global admin should be able to see all departments, even if there are none
        expect(departments.data.length).toBeGreaterThanOrEqual(0);
        if (departments.data.length > 0) {
          expect(departments.data[0]).toHaveProperty("id");
          expect(departments.data[0]).toHaveProperty("name");
        }
      });

      test("returns 403 Forbidden for non-org member", async ({
        request,
        registeredApiUser,
        orgTeamApiUser,
      }) => {
        const org = orgTeamApiUser.user.orgs[0]; // An org the registered user is not part of
        const response = await registeredApiUser.context.get(
          `organizations/${org.id}/departments`,
        );
        expect(response.ok()).toBeFalsy();
        expect(response.status()).toBe(403);
      });

      test("returns empty array when no departments exist", async ({
        request,
        orgOwnerApiUser,
        testDatabase,
      }) => {
        // Create a new organization without any departments
        const newOrg = await testDatabase.seeder.createOrganization(
          "Org Without Departments",
          orgOwnerApiUser.user.id,
        );

        const response = await orgOwnerApiUser.context.get(
          `organizations/${newOrg.id}/departments`,
        );
        expect(response.ok()).toBeTruthy();
        expect(response.status()).toBe(200);
        const departments = await response.json();
        expect(Array.isArray(departments.data)).toBeTruthy();
        expect(departments.data.length).toBe(0);
      });
    });

    test.describe("POST /organizations/{orgId}/departments", () => {
      test("returns 200 for organization admin", async ({
        request,
        orgOwnerApiUser,
      }) => {
        const org = orgOwnerApiUser.user.orgs[0];
        const departmentName = "Test API Create Department Org Admin";
        const response = await orgOwnerApiUser.context.post(
          `organizations/${org.id}/departments`,
          {
            data: { name: departmentName },
          },
        );
        expect(response.ok()).toBeTruthy();
        const department = await response.json();
        expect(department.data).toMatchObject({ name: departmentName });
      });

      test("returns 403 for organization member", async ({
        request,
        orgOwnerApiUser,
        orgTeamApiUser,
      }) => {
        const org = orgOwnerApiUser.user.orgs[0];
        const departmentName = "Test API Create Department Org Member";
        const response = await orgTeamApiUser.context.post(
          `organizations/${org.id}/departments`,
          {
            data: { name: departmentName },
          },
        );
        expect(response.ok()).toBeFalsy();
        expect(response.status()).toBe(403);
      });

      test("returns 200 for global admin", async ({
        request,
        orgOwnerApiUser,
        adminApiUser,
      }) => {
        const org = orgOwnerApiUser.user.orgs[0];
        const departmentName = "Test API Create Department Global Admin";
        const response = await adminApiUser.context.post(
          `organizations/${org.id}/departments`,
          {
            data: { name: departmentName },
          },
        );
        expect(response.ok()).toBeTruthy();
        const department = await response.json();
        expect(department.data).toMatchObject({ name: departmentName });
      });
    });

    test.describe("GET /api/organizations/{orgId}/departments/{deptId}", () => {
      test("returns 200 and department data for department member", async ({
        request,
        deptTeamApiUser,
      }) => {
        const user = deptTeamApiUser.user;
        const dept = user.depts[0];
        const response = await deptTeamApiUser.context.get(
          `organizations/${dept.organization_id}/departments/${dept.id}`,
        );

        expect(response.ok()).toBeTruthy();
        expect(response.status()).toBe(200);
        const deptData = await response.json();
        expect(deptData.data.department).toMatchObject({
          id: dept.id,
          name: dept.name,
        });
      });

      test(
        "returns 200 and department data for global admin",
        { tag: ["@admin"] },
        async ({ request, adminApiUser, deptTeamApiUser }) => {
          const org = deptTeamApiUser.user.orgs[0];
          const dept = deptTeamApiUser.user.depts[0];
          const response = await adminApiUser.context.get(
            `organizations/${org.id}/departments/${dept.id}`,
          );
          expect(response.ok()).toBeTruthy();
          expect(response.status()).toBe(200);
          const deptData = await response.json();
          expect(deptData.data.department).toMatchObject({
            id: dept.id,
            name: dept.name,
          });
        },
      );

      test("returns 403 Forbidden for non department member", async ({
        request,
        registeredApiUser,
        deptTeamApiUser,
      }) => {
        const org = deptTeamApiUser.user.orgs[0];
        const dept = deptTeamApiUser.user.depts[0];
        const response = await registeredApiUser.context.get(
          `organizations/${org.id}/departments/${dept.id}`,
        );
        expect(response.ok()).toBeFalsy();
        expect(response.status()).toBe(403);
      });
    });

    test.describe("PUT /organization/{orgId}/departments/{deptId}", () => {
      test("Org Admin can update any department in their organization", async ({
        orgAdminApiUser,
        testDatabase,
      }) => {
        const org = orgAdminApiUser.user.orgs[0];
        const dept = await testDatabase.seeder.createDepartment(
          "Dept to Update",
          org.id,
        );
        const updatedName = "Updated by Org Admin";
        const response = await orgAdminApiUser.context.put(
          `organizations/${org.id}/departments/${dept.id}`,
          {
            data: { name: updatedName },
          },
        );
        expect(response.ok()).toBeTruthy();
        expect(response.status()).toBe(200);
        const updatedDept = await response.json();
        expect(updatedDept.data.name).toBe(updatedName);
      });
    });

    test.describe("DELETE /api/organizations/{orgId}/departments/{deptId}", () => {
      test("returns 200 for org admin", async ({
        request,
        orgOwnerApiUser,
        testDatabase,
      }) => {
        const org = orgOwnerApiUser.user.orgs[0];
        const dept = await testDatabase.seeder.createDepartment(
          "deptToDeleteOrgOwner",
          org.id,
        );

        const response = await orgOwnerApiUser.context.delete(
          `organizations/${org.id}/departments/${dept.id}`,
        );
        expect(response.ok()).toBeTruthy();
        expect(response.status()).toBe(200);

        const confirmDeptDeleted = await testDatabase.seeder.getDepartmentById(
          dept.id,
        );
        expect(confirmDeptDeleted).toBeNull();
      });

      test(
        "returns 200 for global admin",
        { tag: ["@admin"] },
        async ({ request, orgOwnerApiUser, adminApiUser, testDatabase }) => {
          const org = orgOwnerApiUser.user.orgs[0];
          const dept = await testDatabase.seeder.createDepartment(
            "deptToDeleteGlobalAdmin",
            org.id,
          );

          const response = await adminApiUser.context.delete(
            `organizations/${org.id}/departments/${dept.id}`,
          );
          expect(response.ok()).toBeTruthy();
          expect(response.status()).toBe(200);

          const confirmDeptDeleted =
            await testDatabase.seeder.getDepartmentById(dept.id);
          expect(confirmDeptDeleted).toBeNull();
        },
      );

      test("returns 403 Forbidden for non org admin", async ({
        request,
        deptTeamApiUser,
      }) => {
        const org = deptTeamApiUser.user.orgs[0];
        const dept = deptTeamApiUser.user.depts[0];
        const response = await deptTeamApiUser.context.delete(
          `organizations/${org.id}/departments/${dept.id}`,
        );
        expect(response.ok()).toBeFalsy();
        expect(response.status()).toBe(403);
      });
    });
  },
);
