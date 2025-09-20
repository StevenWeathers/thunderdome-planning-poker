import { chromium, expect, FullConfig } from "@playwright/test";
import * as fs from "fs";
import * as path from "path";
import { RegisterPage } from "@fixtures/pages/register-page";
import { LoginPage } from "@fixtures/pages/login-page";
import { setupDB } from "@fixtures/db/setup";
import { TestUsers } from "@fixtures/types";
import { setupAdminUser } from "@fixtures/db/admin-user";
import { setupRegisteredUser } from "@fixtures/db/registered-user";
import { setupDeleteRegisteredUser } from "@fixtures/db/registered-delete-user";
import { setupVerifiedUser } from "@fixtures/db/verified-user";
import { setupAPIUser } from "@fixtures/db/api-user";
import { setupAdminAPIUser } from "@fixtures/db/adminapi-user";
import { ThunderdomeSeeder } from "@fixtures/db/seeder";
import * as process from "process";

async function globalSetup(config: FullConfig) {
  // Ensure storage directory exists (cross-platform)
  const storageDir = path.resolve(__dirname, "storage");
  if (!fs.existsSync(storageDir)) {
    fs.mkdirSync(storageDir, { recursive: true });
  }

  const pool = setupDB();
  const seeder = new ThunderdomeSeeder(pool);
  const testUsers: TestUsers = {};

  // Create an organization owner with a department and teams
  // then create users associated to different entities of that users organization and non-org team

  // Seed Org Owner
  const orgOwnerUserName = "E2EOrgOwner";
  await seeder.deleteUserByName(orgOwnerUserName);
  const orgOwnerUserEmail = seeder.generateEmail(orgOwnerUserName);
  const orgOwnerUser = await seeder.createUser(
    orgOwnerUserName,
    orgOwnerUserEmail,
  );
  await seeder.addUserAPIKey(
    orgOwnerUser.id,
    "fi5gnBQt.a8c02d29f3984bd33b2a80bb88257e92d56886b840d51897211fc5159e75c668",
    "E2E API Key",
  );
  const org = await seeder.createOrganization(
    `${orgOwnerUserName} Org`,
    orgOwnerUser.id,
  );
  const dept = await seeder.createDepartment(
    `${orgOwnerUserName} Org Dept`,
    org.id,
  );
  const orgTeam = await seeder.createOrgTeam(
    `${orgOwnerUserName} Org Team`,
    org.id,
  );
  const deptTeam = await seeder.createOrgTeam(
    `${orgOwnerUserName} Org Dept Team`,
    org.id,
    dept.id,
  );
  const nonOrgTeam = await seeder.createUserTeam(
    `${orgOwnerUserName} Non-Org Team`,
    orgOwnerUser.id,
  );
  testUsers.orgOwner = {
    ...orgOwnerUser,
    apikey: "fi5gnBQt.iuiU9V2k_GfNmqjgw1d_m_RIoTRA6uKA",
    orgs: [org],
    orgTeams: [orgTeam],
    depts: [dept],
    deptTeams: [deptTeam],
    teams: [nonOrgTeam],
  };
  console.log(`${orgOwnerUserName} seeded successfully`);

  // Seed User associated to org, org team, and non-org team
  const orgTeamUserName = "E2EOrgTeamUser1";
  await seeder.deleteUserByName(orgTeamUserName);
  const orgTeamUserEmail = seeder.generateEmail(orgTeamUserName);
  const orgTeamUser = await seeder.createUser(
    orgTeamUserName,
    orgTeamUserEmail,
  );
  await seeder.addUserAPIKey(
    orgTeamUser.id,
    "0DIvRK3H.54e51a06bac4e9b3dcf815282b7955f6e8e2a2cfef7b191ef6d596ea88eb6e0a",
    "E2E API Key",
  );
  await seeder.addUserToOrg(orgTeamUser.id, org.id);
  await seeder.addUserToTeam(orgTeamUser.id, orgTeam.id);
  await seeder.addUserToTeam(orgTeamUser.id, nonOrgTeam.id);
  testUsers.orgTeamMember = {
    ...orgTeamUser,
    apikey: "0DIvRK3H.arPJj6FmanALShjgRUX0E_nukNXhqRT=",
    orgs: [org],
    orgTeams: [orgTeam],
    depts: [],
    deptTeams: [],
    teams: [nonOrgTeam],
  };
  console.log(`${orgTeamUserName} seeded successfully`);

  // Seed User associated to org, dept, and dept team
  const deptTeamUserName = "E2EOrgDeptTeamUser1";
  await seeder.deleteUserByName(deptTeamUserName);
  const deptTeamUserEmail = seeder.generateEmail(deptTeamUserName);
  const deptTeamUser = await seeder.createUser(
    deptTeamUserName,
    deptTeamUserEmail,
  );
  await seeder.addUserAPIKey(
    deptTeamUser.id,
    "rSFqCQ_6.24bf955b6f66556171e111a8ce492a0b24aa7dbfbeb6c3fa0ec9a7e66f048763",
    "E2E API Key",
  );
  await seeder.addUserToOrg(deptTeamUser.id, org.id);
  await seeder.addUserToDept(deptTeamUser.id, dept.id);
  await seeder.addUserToTeam(deptTeamUser.id, deptTeam.id);
  testUsers.deptTeamMember = {
    ...deptTeamUser,
    apikey: "rSFqCQ_6.3kzJ=ZMHHfS-o69DF!RhMa+iKQPKZ!6!",
    orgs: [org],
    orgTeams: [],
    depts: [dept],
    deptTeams: [deptTeam],
    teams: [],
  };
  console.log(`${deptTeamUserName} seeded successfully`);

  // Seed User associated to org as ADMIN
  const orgAdminUserName = "E2EOrgAdmin1";
  await seeder.deleteUserByName(orgAdminUserName);
  const orgAdminUserEmail = seeder.generateEmail(orgAdminUserName);
  const orgAdminUser = await seeder.createUser(
    orgAdminUserName,
    orgAdminUserEmail,
  );
  await seeder.addUserAPIKey(
    orgAdminUser.id,
    "Kgc_ujA_.19b8ddc3dcf9af62b91ec00e0f170b2b99ae01392aad9c2fc5140553463219f4",
    "E2E API Key",
  );
  await seeder.addUserToOrg(orgAdminUser.id, org.id, "ADMIN");
  testUsers.orgAdmin = {
    ...orgAdminUser,
    apikey: "Kgc_ujA_.E$ZgDdR4uT+Pvl0Qh=-ZGzz0xobbCoad",
    orgs: [org],
    orgTeams: [],
    depts: [],
    deptTeams: [],
    teams: [],
  };
  console.log(`${orgAdminUserName} seeded successfully`);

  // Seed User associated to dept as ADMIN
  const deptAdminUserName = "E2EOrgDeptAdmin1";
  await seeder.deleteUserByName(deptAdminUserName);
  const deptAdminUserEmail = seeder.generateEmail(deptAdminUserName);
  const deptAdminUser = await seeder.createUser(
    deptAdminUserName,
    deptAdminUserEmail,
  );
  await seeder.addUserAPIKey(
    deptAdminUser.id,
    "+JEDJS1o.b787661dff2eb536eaa99b217ebf6747e4c511a0a5a97def388a67d449f68c65",
    "E2E API Key",
  );
  await seeder.addUserToOrg(deptAdminUser.id, org.id);
  await seeder.addUserToDept(deptAdminUser.id, dept.id, "ADMIN");
  testUsers.deptAdmin = {
    ...deptAdminUser,
    apikey: "+JEDJS1o.4FNyoRq4-V0TOyt+dGF93xDUGR_HBOwQ",
    orgs: [org],
    orgTeams: [],
    depts: [dept],
    deptTeams: [],
    teams: [],
  };
  console.log(`${deptAdminUserName} seeded successfully`);

  // Seed User associated to team as ADMIN
  const teamAdminUserName = "E2ETeamAdmin1";
  await seeder.deleteUserByName(teamAdminUserName);
  const teamAdminUserEmail = seeder.generateEmail(teamAdminUserName);
  const teamAdminUser = await seeder.createUser(
    teamAdminUserName,
    teamAdminUserEmail,
  );
  await seeder.addUserAPIKey(
    teamAdminUser.id,
    "WSpBNDzh.19a1813c398a2f56da96ad0567341e779b37d36b85ccc508511331eaee9d9a3c",
    "E2E API Key",
  );
  await seeder.addUserToTeam(teamAdminUser.id, nonOrgTeam.id, "ADMIN");
  testUsers.teamAdmin = {
    ...teamAdminUser,
    apikey: "WSpBNDzh.dTkh7e$8w$54WaxybQk9ObZje4$3sY0C",
    orgs: [],
    orgTeams: [],
    depts: [],
    deptTeams: [],
    teams: [nonOrgTeam],
  };
  console.log(`${teamAdminUserName} seeded successfully`);

  // Seed User associated to Org team as ADMIN
  const orgTeamAdminName = "E2EOrgTeamAdmin1";
  await seeder.deleteUserByName(orgTeamAdminName);
  const orgTeamAdminEmail = seeder.generateEmail(orgTeamAdminName);
  const orgTeamAdmin = await seeder.createUser(
    orgTeamAdminName,
    orgTeamAdminEmail,
  );
  await seeder.addUserAPIKey(
    orgTeamAdmin.id,
    "Q+PSuFMl.74f05388d298e72ad4b20b6c581303e3c3230cb0242376a07e025a8d3bbb4039",
    "E2E API Key",
  );
  await seeder.addUserToOrg(orgTeamAdmin.id, org.id);
  await seeder.addUserToTeam(orgTeamAdmin.id, orgTeam.id, "ADMIN");
  testUsers.orgTeamAdmin = {
    ...orgTeamAdmin,
    apikey: "Q+PSuFMl.DRWrKdk82WkaxSmCCntSe35NrVsqYqMj",
    orgs: [org],
    orgTeams: [orgTeam],
    depts: [],
    deptTeams: [],
    teams: [],
  };
  console.log(`${orgTeamAdminName} seeded successfully`);

  const deptTeamAdminUserName = "E2EOrgDeptTeamAdmin1";
  await seeder.deleteUserByName(deptTeamAdminUserName);
  const deptTeamAdminUserEmail = seeder.generateEmail(deptTeamAdminUserName);
  const deptTeamAdminUser = await seeder.createUser(
    deptTeamAdminUserName,
    deptTeamAdminUserEmail,
  );
  await seeder.addUserAPIKey(
    deptTeamAdminUser.id,
    "s8GQVacj.566b998dd2b0966bef4092fdda4a5d40b68ddcca46860ed512ce34d2a0acef59",
    "E2E API Key",
  );
  await seeder.addUserToOrg(deptTeamAdminUser.id, org.id);
  await seeder.addUserToDept(deptTeamAdminUser.id, dept.id);
  await seeder.addUserToTeam(deptTeamAdminUser.id, deptTeam.id, "ADMIN");
  testUsers.deptTeamAdmin = {
    ...deptTeamAdminUser,
    apikey: "s8GQVacj.osP4_y8GsGmGgkfq8-T!FAjjDDYMISM$",
    orgs: [org],
    orgTeams: [],
    depts: [dept],
    deptTeams: [deptTeam],
    teams: [],
  };
  console.log(`${deptTeamAdminUserName} seeded successfully`);

  const baseUrl = config.projects[0].use.baseURL;
  const browser = await chromium.launch();

  await setupAdminAPIUser.teardown(pool);
  testUsers.adminUser = await setupAdminAPIUser.seed(pool);

  await setupAPIUser.teardown(pool);
  testUsers.registeredUser = await setupAPIUser.seed(pool);

  const adminPage = await browser.newPage({
    baseURL: baseUrl,
  });
  await setupAdminUser.teardown(pool);
  const au = await setupAdminUser.seed(pool);
  const adminLoginPage = new LoginPage(adminPage);
  await adminLoginPage.goto();
  await adminLoginPage.login(au.email, "kentRules!");
  await expect(adminLoginPage.page.locator("h1")).toHaveText(
    "Welcome back, E2E ADMIN",
  );
  await adminLoginPage.page.context().storageState({
    path: path.resolve(__dirname, "storage", "adminStorageState.json"),
  });

  const registeredPage = await browser.newPage({
    baseURL: baseUrl,
  });
  await setupRegisteredUser.teardown(pool);
  const ru = await setupRegisteredUser.seed(pool);
  const registeredRegisterPage = new LoginPage(registeredPage);
  await registeredRegisterPage.goto();
  await registeredRegisterPage.login(ru.email, "kentRules!");
  await expect(registeredRegisterPage.page.locator("h1")).toHaveText(
    "Welcome back, E2E Registered User",
  );
  await registeredRegisterPage.page.context().storageState({
    path: path.resolve(__dirname, "storage", "registeredStorageState.json"),
  });

  const verifiedPage = await browser.newPage({
    baseURL: baseUrl,
  });
  await setupVerifiedUser.teardown(pool);
  const vu = await setupVerifiedUser.seed(pool);
  const userVerifiedPage = new LoginPage(verifiedPage);
  await userVerifiedPage.goto();
  await userVerifiedPage.login(vu.email, "kentRules!");
  await expect(userVerifiedPage.page.locator("h1")).toHaveText(
    "Welcome back, E2E Verified User",
  );
  await userVerifiedPage.page.context().storageState({
    path: path.resolve(__dirname, "storage", "verifiedStorageState.json"),
  });

  const guestPage = await browser.newPage({
    baseURL: baseUrl,
  });
  const guestRegisterPage = new RegisterPage(guestPage);
  await guestRegisterPage.goto();
  await guestRegisterPage.createGuestUser("E2E Guest");
  await expect(guestRegisterPage.page.locator("h1")).toHaveText(
    "Welcome back, E2E Guest",
  );
  await guestRegisterPage.page.context().storageState({
    path: path.resolve(__dirname, "storage", "guestStorageState.json"),
  });

  const deleteGuestPage = await browser.newPage({
    baseURL: baseUrl,
  });
  const deleteGuestRegisterPage = new RegisterPage(deleteGuestPage);
  await deleteGuestRegisterPage.goto();
  await deleteGuestRegisterPage.createGuestUser("E2E Delete Guest");
  await expect(deleteGuestRegisterPage.page.locator("h1")).toHaveText(
    "Welcome back, E2E Delete Guest",
  );
  await deleteGuestRegisterPage.page.context().storageState({
    path: path.resolve(__dirname, "storage", "deleteGuestStorageState.json"),
  });

  const deleteRegPage = await browser.newPage({
    baseURL: baseUrl,
  });
  await setupDeleteRegisteredUser.teardown(pool);
  const dru = await setupDeleteRegisteredUser.seed(pool);
  const deleteRegisteredPage = new LoginPage(deleteRegPage);
  await deleteRegisteredPage.goto();
  await deleteRegisteredPage.login(dru.email, "kentRules!");
  await expect(deleteRegisteredPage.page.locator("h1")).toHaveText(
    "Welcome back, E2E Delete User",
  );
  await deleteRegisteredPage.page.context().storageState({
    path: path.resolve(
      __dirname,
      "storage",
      "deleteRegisteredStorageState.json",
    ),
  });

  await browser.close();

  process.env.TEST_USERS = JSON.stringify(testUsers);

  console.log("All users seeded successfully");
}

export default globalSetup;
