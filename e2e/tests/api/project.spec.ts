import { expect, test } from "@fixtures/test-setup";

test.describe.skip("Project API", { tag: ["@api", "@project"] }, () => {
  // Admin Project Operations
  test.describe("GET /admin/projects", () => {
    test("returns paginated list of projects for Global Admin", async ({
      adminApiUser,
    }) => {
      const response = await adminApiUser.context.get("admin/projects");
      expect(response.ok()).toBeTruthy();
      expect(response.status()).toBe(200);
      const projects = await response.json();
      expect(projects.data).toBeDefined();
      expect(projects.meta).toBeDefined();
      expect(projects.meta).toMatchObject({
        count: expect.any(Number),
        offset: expect.any(Number),
        limit: expect.any(Number),
      });
    });

    test("returns forbidden for non-admin user", async ({
      registeredApiUser,
    }) => {
      const response = await registeredApiUser.context.get("admin/projects");
      expect(response.ok()).toBeFalsy();
      expect(response.status()).toBe(403);
    });

    test("supports pagination parameters", async ({ adminApiUser }) => {
      const response = await adminApiUser.context.get(
        "admin/projects?limit=5&offset=0",
      );
      expect(response.ok()).toBeTruthy();
      expect(response.status()).toBe(200);
      const projects = await response.json();
      expect(projects.meta.limit).toBe(5);
      expect(projects.meta.offset).toBe(0);
    });
  });

  test.describe("POST /admin/projects", () => {
    test("creates project with organization association", async ({
      adminApiUser,
      orgOwnerApiUser,
    }) => {
      const projectData = {
        projectKey: "TESTPROJ",
        name: "Test Admin Project",
        description: "Created by admin",
        organizationId: orgOwnerApiUser.user.orgs[0].id,
      };

      const response = await adminApiUser.context.post("admin/projects", {
        data: projectData,
      });
      expect(response.ok()).toBeTruthy();
      expect(response.status()).toBe(200);
      const project = await response.json();
      expect(project.data).toMatchObject({
        projectKey: projectData.projectKey,
        name: projectData.name,
        description: projectData.description,
        organizationId: projectData.organizationId,
      });
    });

    test("creates project with team association", async ({
      adminApiUser,
      teamAdminApiUser,
    }) => {
      const projectData = {
        projectKey: "TEAMPROJ",
        name: "Test Team Project",
        description: "Created by admin for team",
        teamId: teamAdminApiUser.user.teams[0].id,
      };

      const response = await adminApiUser.context.post("admin/projects", {
        data: projectData,
      });
      expect(response.ok()).toBeTruthy();
      expect(response.status()).toBe(200);
      const project = await response.json();
      expect(project.data).toMatchObject({
        projectKey: projectData.projectKey,
        name: projectData.name,
        teamId: projectData.teamId,
      });
    });

    test("fails when no association provided", async ({ adminApiUser }) => {
      const projectData = {
        projectKey: "NOPROJ",
        name: "Project Without Association",
        description: "Should fail",
      };

      const response = await adminApiUser.context.post("admin/projects", {
        data: projectData,
      });
      expect(response.ok()).toBeFalsy();
      expect(response.status()).toBe(400);
    });

    test("fails with invalid project key format", async ({
      adminApiUser,
      orgOwnerApiUser,
    }) => {
      const projectData = {
        projectKey: "x", // Too short
        name: "Invalid Key Project",
        organizationId: orgOwnerApiUser.user.orgs[0].id,
      };

      const response = await adminApiUser.context.post("admin/projects", {
        data: projectData,
      });
      expect(response.ok()).toBeFalsy();
      expect(response.status()).toBe(400);
    });

    test("returns forbidden for non-admin user", async ({
      registeredApiUser,
      orgOwnerApiUser,
    }) => {
      const projectData = {
        projectKey: "FORBIDDEN",
        name: "Unauthorized Project",
        organizationId: orgOwnerApiUser.user.orgs[0].id,
      };

      const response = await registeredApiUser.context.post("admin/projects", {
        data: projectData,
      });
      expect(response.ok()).toBeFalsy();
      expect(response.status()).toBe(403);
    });
  });

  test.describe("GET /admin/projects/{projectId}", () => {
    test.skip("returns project details for Global Admin", async ({
      adminApiUser,
      testDatabase,
      orgOwnerApiUser,
    }) => {
      const project = await testDatabase.seeder.createProject(
        "GETTEST",
        "Get Test Project",
        orgOwnerApiUser.user.orgs[0].id,
      );

      const response = await adminApiUser.context.get(
        `admin/projects/${project.id}`,
      );
      expect(response.ok()).toBeTruthy();
      expect(response.status()).toBe(200);
      const projectData = await response.json();
      expect(projectData.data).toMatchObject({
        id: project.id,
        projectKey: "GETTEST",
        name: "Get Test Project",
      });
    });

    test("returns 404 for non-existent project", async ({ adminApiUser }) => {
      const fakeId = "00000000-0000-0000-0000-000000000000";
      const response = await adminApiUser.context.get(
        `admin/projects/${fakeId}`,
      );
      expect(response.ok()).toBeFalsy();
      expect(response.status()).toBe(404);
    });

    test("returns forbidden for non-admin user", async ({
      registeredApiUser,
      testDatabase,
      orgOwnerApiUser,
    }) => {
      const project = await testDatabase.seeder.createProject(
        "FORBID",
        "Forbidden Project",
        orgOwnerApiUser.user.orgs[0].id,
      );

      const response = await registeredApiUser.context.get(
        `admin/projects/${project.id}`,
      );
      expect(response.ok()).toBeFalsy();
      expect(response.status()).toBe(403);
    });
  });

  test.describe("PUT /admin/projects/{projectId}", () => {
    test("updates project successfully", async ({
      adminApiUser,
      testDatabase,
      orgOwnerApiUser,
    }) => {
      const project = await testDatabase.seeder.createProject(
        "UPDATE",
        "Original Name",
        orgOwnerApiUser.user.orgs[0].id,
      );

      const updateData = {
        projectKey: "UPDATED",
        name: "Updated Project Name",
        description: "Updated description",
        organizationId: orgOwnerApiUser.user.orgs[0].id,
      };

      const response = await adminApiUser.context.put(
        `admin/projects/${project.id}`,
        { data: updateData },
      );
      expect(response.ok()).toBeTruthy();
      expect(response.status()).toBe(200);
      const updatedProject = await response.json();
      expect(updatedProject.data).toMatchObject({
        projectKey: "UPDATED",
        name: "Updated Project Name",
        description: "Updated description",
      });
    });

    test("returns forbidden for non-admin user", async ({
      registeredApiUser,
      testDatabase,
      orgOwnerApiUser,
    }) => {
      const project = await testDatabase.seeder.createProject(
        "NOUPDATE",
        "No Update Project",
        orgOwnerApiUser.user.orgs[0].id,
      );

      const updateData = {
        projectKey: "NOUPDATE",
        name: "Should Not Update",
        organizationId: orgOwnerApiUser.user.orgs[0].id,
      };

      const response = await registeredApiUser.context.put(
        `admin/projects/${project.id}`,
        { data: updateData },
      );
      expect(response.ok()).toBeFalsy();
      expect(response.status()).toBe(403);
    });
  });

  test.describe("DELETE /admin/projects/{projectId}", () => {
    test("deletes project successfully", async ({
      adminApiUser,
      testDatabase,
      orgOwnerApiUser,
    }) => {
      const project = await testDatabase.seeder.createProject(
        "DELETE",
        "Delete Test Project",
        orgOwnerApiUser.user.orgs[0].id,
      );

      const response = await adminApiUser.context.delete(
        `admin/projects/${project.id}`,
      );
      expect(response.ok()).toBeTruthy();
      expect(response.status()).toBe(200);
    });

    test("returns forbidden for non-admin user", async ({
      registeredApiUser,
      testDatabase,
      orgOwnerApiUser,
    }) => {
      const project = await testDatabase.seeder.createProject(
        "NODELETE",
        "No Delete Project",
        orgOwnerApiUser.user.orgs[0].id,
      );

      const response = await registeredApiUser.context.delete(
        `admin/projects/${project.id}`,
      );
      expect(response.ok()).toBeFalsy();
      expect(response.status()).toBe(403);
    });
  });

  // Organization Project Operations
  test.describe("GET /organizations/{orgId}/projects", () => {
    test("Organization Admin can view organization projects", async ({
      orgAdminApiUser,
    }) => {
      const orgId = orgAdminApiUser.user.orgs[0].id;
      const response = await orgAdminApiUser.context.get(
        `organizations/${orgId}/projects`,
      );
      expect(response.ok()).toBeTruthy();
      expect(response.status()).toBe(200);
      const projects = await response.json();
      expect(projects.data).toBeDefined();
      expect(Array.isArray(projects.data)).toBeTruthy();
    });

    test("Organization Member can view organization projects", async ({
      orgOwnerApiUser,
    }) => {
      const orgId = orgOwnerApiUser.user.orgs[0].id;
      const response = await orgOwnerApiUser.context.get(
        `organizations/${orgId}/projects`,
      );
      expect(response.ok()).toBeTruthy();
      expect(response.status()).toBe(200);
      const projects = await response.json();
      expect(projects.data).toBeDefined();
    });

    test("Non-organization member cannot view organization projects", async ({
      registeredApiUser,
      orgOwnerApiUser,
    }) => {
      const orgId = orgOwnerApiUser.user.orgs[0].id;
      const response = await registeredApiUser.context.get(
        `organizations/${orgId}/projects`,
      );
      expect(response.ok()).toBeFalsy();
      expect(response.status()).toBe(403);
    });

    test("Global Admin can view any organization projects", async ({
      adminApiUser,
      orgOwnerApiUser,
    }) => {
      const orgId = orgOwnerApiUser.user.orgs[0].id;
      const response = await adminApiUser.context.get(
        `organizations/${orgId}/projects`,
      );
      expect(response.ok()).toBeTruthy();
      expect(response.status()).toBe(200);
    });
  });

  test.describe("POST /organizations/{orgId}/projects", () => {
    test("Organization Admin can create organization project", async ({
      orgAdminApiUser,
    }) => {
      const orgId = orgAdminApiUser.user.orgs[0].id;
      const projectData = {
        projectKey: "ORGPROJ",
        name: "Organization Project",
        description: "Created by org admin",
      };

      const response = await orgAdminApiUser.context.post(
        `organizations/${orgId}/projects`,
        { data: projectData },
      );
      expect(response.ok()).toBeTruthy();
      expect(response.status()).toBe(200);
      const project = await response.json();
      expect(project.data).toMatchObject({
        projectKey: projectData.projectKey,
        name: projectData.name,
        organizationId: orgId,
      });
    });

    test.skip("Organization Member cannot create organization project", async ({
      orgOwnerApiUser,
    }) => {
      const orgId = orgOwnerApiUser.user.orgs[0].id;
      const projectData = {
        projectKey: "NOPERM",
        name: "No Permission Project",
      };

      const response = await orgOwnerApiUser.context.post(
        `organizations/${orgId}/projects`,
        { data: projectData },
      );
      expect(response.ok()).toBeFalsy();
      expect(response.status()).toBe(403);
    });

    test("Non-organization member cannot create organization project", async ({
      registeredApiUser,
      orgOwnerApiUser,
    }) => {
      const orgId = orgOwnerApiUser.user.orgs[0].id;
      const projectData = {
        projectKey: "NOTMEMBER",
        name: "Not Member Project",
      };

      const response = await registeredApiUser.context.post(
        `organizations/${orgId}/projects`,
        { data: projectData },
      );
      expect(response.ok()).toBeFalsy();
      expect(response.status()).toBe(403);
    });
  });

  test.describe("PUT /organizations/{orgId}/projects/{projectId}", () => {
    test("Organization Admin can update organization project", async ({
      orgAdminApiUser,
      testDatabase,
    }) => {
      const orgId = orgAdminApiUser.user.orgs[0].id;
      const project = await testDatabase.seeder.createOrganizationProject(
        "ORGUPDATE",
        "Original Org Project",
        orgId,
      );

      const updateData = {
        projectKey: "ORGUPDATED",
        name: "Updated Org Project",
        description: "Updated by org admin",
      };

      const response = await orgAdminApiUser.context.put(
        `organizations/${orgId}/projects/${project.id}`,
        { data: updateData },
      );
      expect(response.ok()).toBeTruthy();
      expect(response.status()).toBe(200);
      const updatedProject = await response.json();
      expect(updatedProject.data.name).toBe("Updated Org Project");
    });

    test.skip("Organization Member cannot update organization project", async ({
      orgOwnerApiUser,
      testDatabase,
    }) => {
      const orgId = orgOwnerApiUser.user.orgs[0].id;
      const project = await testDatabase.seeder.createOrganizationProject(
        "NOUPDATE",
        "No Update Project",
        orgId,
      );

      const updateData = {
        projectKey: "NOUPDATE",
        name: "Should Not Update",
      };

      const response = await orgOwnerApiUser.context.put(
        `organizations/${orgId}/projects/${project.id}`,
        { data: updateData },
      );
      expect(response.ok()).toBeFalsy();
      expect(response.status()).toBe(403);
    });
  });

  test.describe("DELETE /organizations/{orgId}/projects/{projectId}", () => {
    test("Organization Admin can delete organization project", async ({
      orgAdminApiUser,
      testDatabase,
    }) => {
      const orgId = orgAdminApiUser.user.orgs[0].id;
      const project = await testDatabase.seeder.createOrganizationProject(
        "ORGDELETE",
        "Delete Org Project",
        orgId,
      );

      const response = await orgAdminApiUser.context.delete(
        `organizations/${orgId}/projects/${project.id}`,
      );
      expect(response.ok()).toBeTruthy();
      expect(response.status()).toBe(200);
    });

    test.skip("Organization Member cannot delete organization project", async ({
      orgOwnerApiUser,
      testDatabase,
    }) => {
      const orgId = orgOwnerApiUser.user.orgs[0].id;
      const project = await testDatabase.seeder.createOrganizationProject(
        "NODELETE",
        "No Delete Project",
        orgId,
      );

      const response = await orgOwnerApiUser.context.delete(
        `organizations/${orgId}/projects/${project.id}`,
      );
      expect(response.ok()).toBeFalsy();
      expect(response.status()).toBe(403);
    });
  });

  // Department Project Operations
  test.describe("GET /organizations/{orgId}/departments/{departmentId}/projects", () => {
    test.skip("Department Admin can view department projects", async ({
      deptAdminApiUser,
      testDatabase,
    }) => {
      const orgId = deptAdminApiUser.user.orgs[0].id;
      const department = await testDatabase.seeder.createDepartment(
        "Test Dept DAP",
        orgId,
      );

      const response = await deptAdminApiUser.context.get(
        `organizations/${orgId}/departments/${department.id}/projects`,
      );
      expect(response.ok()).toBeTruthy();
      expect(response.status()).toBe(200);
      const projects = await response.json();
      expect(projects.data).toBeDefined();
    });

    test.skip("Department Member can view department projects", async ({
      deptTeamApiUser,
      testDatabase,
    }) => {
      const orgId = deptTeamApiUser.user.orgs[0].id;
      const department = await testDatabase.seeder.createDepartment(
        "Test Dept DMDP",
        orgId,
      );

      const response = await deptTeamApiUser.context.get(
        `organizations/${orgId}/departments/${department.id}/projects`,
      );
      expect(response.ok()).toBeTruthy();
      expect(response.status()).toBe(200);
    });

    test("Non-department member cannot view department projects", async ({
      registeredApiUser,
      testDatabase,
      orgOwnerApiUser,
    }) => {
      const orgId = orgOwnerApiUser.user.orgs[0].id;
      const department = await testDatabase.seeder.createDepartment(
        "Test Dept NMM",
        orgId,
      );

      const response = await registeredApiUser.context.get(
        `organizations/${orgId}/departments/${department.id}/projects`,
      );
      expect(response.ok()).toBeFalsy();
      expect(response.status()).toBe(403);
    });
  });

  test.describe("POST /organizations/{orgId}/departments/{departmentId}/projects", () => {
    test.skip("Department Admin can create department project", async ({
      deptAdminApiUser,
      testDatabase,
    }) => {
      const orgId = deptAdminApiUser.user.orgs[0].id;
      const department = await testDatabase.seeder.createDepartment(
        "Test Dept DAD",
        orgId,
      );

      const projectData = {
        projectKey: "DEPTPROJ",
        name: "Department Project",
        description: "Created by dept admin",
      };

      const response = await deptAdminApiUser.context.post(
        `organizations/${orgId}/departments/${department.id}/projects`,
        { data: projectData },
      );
      expect(response.ok()).toBeTruthy();
      expect(response.status()).toBe(200);
      const project = await response.json();
      expect(project.data).toMatchObject({
        projectKey: projectData.projectKey,
        name: projectData.name,
        departmentId: department.id,
      });
    });

    test("Department Member cannot create department project", async ({
      deptTeamApiUser,
      testDatabase,
    }) => {
      const orgId = deptTeamApiUser.user.orgs[0].id;
      const department = await testDatabase.seeder.createDepartment(
        "Test Dept DMP",
        orgId,
      );

      const projectData = {
        projectKey: "NOPERM",
        name: "No Permission Project",
      };

      const response = await deptTeamApiUser.context.post(
        `organizations/${orgId}/departments/${department.id}/projects`,
        { data: projectData },
      );
      expect(response.ok()).toBeFalsy();
      expect(response.status()).toBe(403);
    });
  });

  // Team Project Operations
  test.describe("GET /teams/{teamId}/projects", () => {
    test("Team Admin can view team projects", async ({ teamAdminApiUser }) => {
      const teamId = teamAdminApiUser.user.teams[0].id;
      const response = await teamAdminApiUser.context.get(
        `teams/${teamId}/projects`,
      );
      expect(response.ok()).toBeTruthy();
      expect(response.status()).toBe(200);
      const projects = await response.json();
      expect(projects.data).toBeDefined();
      expect(Array.isArray(projects.data)).toBeTruthy();
    });

    test("Team Member can view team projects", async ({ orgTeamApiUser }) => {
      const teamId = orgTeamApiUser.user.teams[0].id;
      const response = await orgTeamApiUser.context.get(
        `teams/${teamId}/projects`,
      );
      expect(response.ok()).toBeTruthy();
      expect(response.status()).toBe(200);
      const projects = await response.json();
      expect(projects.data).toBeDefined();
    });

    test("Non-team member cannot view team projects", async ({
      registeredApiUser,
      teamAdminApiUser,
    }) => {
      const teamId = teamAdminApiUser.user.teams[0].id;
      const response = await registeredApiUser.context.get(
        `teams/${teamId}/projects`,
      );
      expect(response.ok()).toBeFalsy();
      expect(response.status()).toBe(403);
    });

    test("Global Admin can view any team projects", async ({
      adminApiUser,
      teamAdminApiUser,
    }) => {
      const teamId = teamAdminApiUser.user.teams[0].id;
      const response = await adminApiUser.context.get(
        `teams/${teamId}/projects`,
      );
      expect(response.ok()).toBeTruthy();
      expect(response.status()).toBe(200);
    });
  });

  test.describe("POST /teams/{teamId}/projects", () => {
    test.skip("Team Admin can create team project", async ({
      teamAdminApiUser,
    }) => {
      const teamId = teamAdminApiUser.user.teams[0].id;
      const projectData = {
        projectKey: "TEAMPROJ",
        name: "Team Project",
        description: "Created by team admin",
      };

      const response = await teamAdminApiUser.context.post(
        `teams/${teamId}/projects`,
        { data: projectData },
      );
      expect(response.ok()).toBeTruthy();
      expect(response.status()).toBe(200);
      const project = await response.json();
      expect(project.data).toMatchObject({
        projectKey: projectData.projectKey,
        name: projectData.name,
        teamId: teamId,
      });
    });

    test("Team Member cannot create team project", async ({
      orgTeamApiUser,
    }) => {
      const teamId = orgTeamApiUser.user.teams[0].id;
      const projectData = {
        projectKey: "NOPERM",
        name: "No Permission Project",
      };

      const response = await orgTeamApiUser.context.post(
        `teams/${teamId}/projects`,
        { data: projectData },
      );
      expect(response.ok()).toBeFalsy();
      expect(response.status()).toBe(403);
    });

    test("Non-team member cannot create team project", async ({
      registeredApiUser,
      teamAdminApiUser,
    }) => {
      const teamId = teamAdminApiUser.user.teams[0].id;
      const projectData = {
        projectKey: "NOTMEMBER",
        name: "Not Member Project",
      };

      const response = await registeredApiUser.context.post(
        `teams/${teamId}/projects`,
        { data: projectData },
      );
      expect(response.ok()).toBeFalsy();
      expect(response.status()).toBe(403);
    });
  });

  test.describe("PUT /teams/{teamId}/projects/{projectId}", () => {
    test.skip("Team Admin can update team project", async ({
      teamAdminApiUser,
      testDatabase,
    }) => {
      const teamId = teamAdminApiUser.user.teams[0].id;
      const project = await testDatabase.seeder.createTeamProject(
        "TEAMUPDATE",
        "Original Team Project",
        teamId,
      );

      const updateData = {
        projectKey: "TEAMUPDATED",
        name: "Updated Team Project",
        description: "Updated by team admin",
      };

      const response = await teamAdminApiUser.context.put(
        `teams/${teamId}/projects/${project.id}`,
        { data: updateData },
      );
      expect(response.ok()).toBeTruthy();
      expect(response.status()).toBe(200);
      const updatedProject = await response.json();
      expect(updatedProject.data.name).toBe("Updated Team Project");
    });

    test("Team Member cannot update team project", async ({
      orgTeamApiUser,
      testDatabase,
    }) => {
      const teamId = orgTeamApiUser.user.teams[0].id;
      const project = await testDatabase.seeder.createTeamProject(
        "NOUPDATE",
        "No Update Project",
        teamId,
      );

      const updateData = {
        projectKey: "NOUPDATE",
        name: "Should Not Update",
      };

      const response = await orgTeamApiUser.context.put(
        `teams/${teamId}/projects/${project.id}`,
        { data: updateData },
      );
      expect(response.ok()).toBeFalsy();
      expect(response.status()).toBe(403);
    });

    test("Non-team member cannot update team project", async ({
      registeredApiUser,
      teamAdminApiUser,
      testDatabase,
    }) => {
      const teamId = teamAdminApiUser.user.teams[0].id;
      const project = await testDatabase.seeder.createTeamProject(
        "NOUPDATE",
        "No Update Project",
        teamId,
      );

      const updateData = {
        projectKey: "NOUPDATE",
        name: "Should Not Update",
      };

      const response = await registeredApiUser.context.put(
        `teams/${teamId}/projects/${project.id}`,
        { data: updateData },
      );
      expect(response.ok()).toBeFalsy();
      expect(response.status()).toBe(403);
    });
  });

  test.describe("DELETE /teams/{teamId}/projects/{projectId}", () => {
    test("Team Admin can delete team project", async ({
      teamAdminApiUser,
      testDatabase,
    }) => {
      const teamId = teamAdminApiUser.user.teams[0].id;
      const project = await testDatabase.seeder.createTeamProject(
        "TEAMDELETE",
        "Delete Team Project",
        teamId,
      );

      const response = await teamAdminApiUser.context.delete(
        `teams/${teamId}/projects/${project.id}`,
      );
      expect(response.ok()).toBeTruthy();
      expect(response.status()).toBe(200);
    });

    test("Team Member cannot delete team project", async ({
      orgTeamApiUser,
      testDatabase,
    }) => {
      const teamId = orgTeamApiUser.user.teams[0].id;
      const project = await testDatabase.seeder.createTeamProject(
        "NODELETE",
        "No Delete Project",
        teamId,
      );

      const response = await orgTeamApiUser.context.delete(
        `teams/${teamId}/projects/${project.id}`,
      );
      expect(response.ok()).toBeFalsy();
      expect(response.status()).toBe(403);
    });

    test("Non-team member cannot delete team project", async ({
      registeredApiUser,
      teamAdminApiUser,
      testDatabase,
    }) => {
      const teamId = teamAdminApiUser.user.teams[0].id;
      const project = await testDatabase.seeder.createTeamProject(
        "NODELETE",
        "No Delete Project",
        teamId,
      );

      const response = await registeredApiUser.context.delete(
        `teams/${teamId}/projects/${project.id}`,
      );
      expect(response.ok()).toBeFalsy();
      expect(response.status()).toBe(403);
    });

    test("Global Admin can delete any team project", async ({
      adminApiUser,
      testDatabase,
      teamAdminApiUser,
    }) => {
      const teamId = teamAdminApiUser.user.teams[0].id;
      const project = await testDatabase.seeder.createTeamProject(
        "ADMINDELETE",
        "Admin Delete Project",
        teamId,
      );

      const response = await adminApiUser.context.delete(
        `teams/${teamId}/projects/${project.id}`,
      );
      expect(response.ok()).toBeTruthy();
      expect(response.status()).toBe(200);
    });
  });
});
