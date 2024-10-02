import { expect, test } from "@fixtures/user-sessions";
import { DepartmentPage } from "@fixtures/pages/department-page";

test.describe("Department page", { tag: "@department" }, () => {
  test.describe("Unauthenticated user", { tag: "@unauthenticated" }, () => {
    test("redirects to login", async ({ page }) => {
      const departmentPage = new DepartmentPage(page);
      await departmentPage.goto(
        "bbaf82ef-a2d3-4e9a-b824-5e56a03ac3aa",
        "bbaf82ef-a2d3-4e9a-b824-5e56a03ac3aa",
      );

      const loginForm = departmentPage.page.locator('form[name="login"]');
      await expect(loginForm).toBeVisible();
    });
  });

  test.describe("Registered user", { tag: "@registered" }, () => {
    test("loads page successfully", async ({ registeredPage }) => {
      const testOrgName = "E2E TEST ORGANIZATION";
      const testDepartmentName = "E2E TEST DEPARTMENT";
      const departmentPage = new DepartmentPage(registeredPage.page);
      const org = await registeredPage.createOrg(testOrgName);
      const dept = await registeredPage.createOrgDepartment(
        org.id,
        testDepartmentName,
      );

      await departmentPage.goto(org.id, dept.id);
      await expect(departmentPage.page.locator("h1")).toContainText(
        testDepartmentName,
      );
    });
  });
});
