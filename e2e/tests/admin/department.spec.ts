import { expect, test } from "../../fixtures/user-sessions";
import { AdminDepartmentPage } from "../../fixtures/admin/department-page";

test.describe("The Admin Department Page", () => {
  test.describe("Unauthenticated user", () => {
    test("redirects to login", async ({ page }) => {
      const adminPage = new AdminDepartmentPage(page);

      await adminPage.goto(
        "bbaf82ef-a2d3-4e9a-b824-5e56a03ac3aa",
        "bbaf82ef-a2d3-4e9a-b824-5e56a03ac3aa",
      );

      const loginForm = adminPage.page.locator('form[name="login"]');
      await expect(loginForm).toBeVisible();
    });
  });

  test.describe("Guest user", () => {
    test("redirects to landing", async ({ guestPage }) => {
      const adminPage = new AdminDepartmentPage(guestPage.page);

      await adminPage.goto(
        "bbaf82ef-a2d3-4e9a-b824-5e56a03ac3aa",
        "bbaf82ef-a2d3-4e9a-b824-5e56a03ac3aa",
      );

      const title = adminPage.page.locator("h1");
      await expect(title).toHaveText(
        "Thunderdome is an Agile Planning Poker app with a fun theme",
      );
    });
  });

  test.describe("Non Admin Registered User", () => {
    test("redirects to landing", async ({ registeredPage }) => {
      const adminPage = new AdminDepartmentPage(registeredPage.page);

      await adminPage.goto(
        "bbaf82ef-a2d3-4e9a-b824-5e56a03ac3aa",
        "bbaf82ef-a2d3-4e9a-b824-5e56a03ac3aa",
      );

      const title = adminPage.page.locator("h1");
      await expect(title).toHaveText(
        "Thunderdome is an Agile Planning Poker app with a fun theme",
      );
    });
  });

  test.describe("Admin User", () => {
    test("loads Department page", async ({ registeredPage, adminPage }) => {
      const ap = new AdminDepartmentPage(adminPage.page);
      const testOrgName = "E2E TEST ADMIN ORGANIZATION";
      const testDeptName = "E2E TEST ADMIN DEPARTMENT";
      const org = await registeredPage.createOrg(testOrgName);
      const dept = await registeredPage.createOrgDepartment(
        org.id,
        testDeptName,
      );

      await ap.goto(org.id, dept.id);
      await expect(ap.page.locator("h1")).toContainText(testDeptName);
    });
  });
});
