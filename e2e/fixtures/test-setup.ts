import { APIRequestContext, test as base } from "@playwright/test";
import { setupDB } from "@fixtures/db/setup";
import { ThunderdomeSeeder } from "@fixtures/db/seeder";
import { TestUser, TestUsers } from "@fixtures/types";

export type ApiUser = {
  context: APIRequestContext;
  user: TestUser;
};

export type TestDatabase = {
  pool: ReturnType<typeof setupDB>;
  seeder: ThunderdomeSeeder;
};

export type TestFixtures = {
  testDatabase: ReturnType<TestDatabase>;
  testUsers: ReturnType<TestUsers>;
  registeredApiUser: ApiUser;
  adminApiUser: ApiUser;
  orgOwnerApiUser: ApiUser;
  orgAdminApiUser: ApiUser;
  orgTeamApiUser: ApiUser;
  orgTeamAdminApiUser: ApiUser;
  deptAdminApiUser: ApiUser;
  deptTeamApiUser: ApiUser;
  deptTeamAdminApiUser: ApiUser;
  teamAdminApiUser: ApiUser;
};

const createApiUserFixture = (userName: keyof TestUsers) => {
  return async (
    {
      playwright,
      testUsers,
      baseURL,
    }: { playwright: any; testUsers: TestUsers },
    use: (arg: ApiUser) => Promise<void>,
  ) => {
    const apiUser = testUsers[userName];
    if (!apiUser) {
      throw new Error(`${userName} not found in test users data`);
    }
    const context = await playwright.request.newContext({
      baseURL: `${baseURL}/api/`,
      extraHTTPHeaders: {
        "X-API-Key": apiUser.apikey,
      },
    });
    await use({ context, user: apiUser });
    await context.dispose();
  };
};

export const test = base.extend<TestFixtures>({
  // Existing testData fixture
  testUsers: [
    async ({}, use) => {
      // This assumes testData is generated once per worker
      // If you need it generated for each test, remove the [, {scope: 'worker'}] part
      const data: TestUsers = JSON.parse(process.env.TEST_USERS || "{}");
      await use(data);
    },
    { scope: "worker" },
  ],

  // Existing testDatabase fixture
  testDatabase: [
    async ({}, use) => {
      const pool = setupDB();
      const seeder = new ThunderdomeSeeder(pool);

      await use({ pool, seeder });
    },
    { scope: "worker" },
  ],

  registeredApiUser: createApiUserFixture("registeredUser"),
  adminApiUser: createApiUserFixture("adminUser"),
  orgOwnerApiUser: createApiUserFixture("orgOwner"),
  orgAdminApiUser: createApiUserFixture("orgAdmin"),
  orgTeamApiUser: createApiUserFixture("orgTeamMember"),
  orgTeamAdminApiUser: createApiUserFixture("orgTeamAdmin"),
  deptTeamApiUser: createApiUserFixture("deptTeamMember"),
  deptTeamAdminApiUser: createApiUserFixture("deptTeamAdmin"),
  deptAdminApiUser: createApiUserFixture("deptAdmin"),
  teamAdminApiUser: createApiUserFixture("teamAdmin"),
});

export { expect } from "@playwright/test";
