import { expect, test } from "@fixtures/user-sessions";
import { AdminOrganizationPage } from "@fixtures/admin/organization-page";

test.describe(
  "The Admin Organization Page",
  { tag: ["@administration", "@organization"] },
  () => {
    test.describe("Unauthenticated user", { tag: ["@unauthenticated"] }, () => {
      test("redirects to login", async ({ page }) => {
        const adminPage = new AdminOrganizationPage(page);

        await adminPage.goto("bbaf82ef-a2d3-4e9a-b824-5e56a03ac3aa");

        const loginForm = adminPage.page.locator('form[name="login"]');
        await expect(loginForm).toBeVisible();
      });
    });

    test.describe("Guest user", { tag: ["@guest"] }, () => {
      test("redirects to landing", async ({ guestPage }) => {
        const adminPage = new AdminOrganizationPage(guestPage.page);

        await adminPage.goto("bbaf82ef-a2d3-4e9a-b824-5e56a03ac3aa");

        const title = adminPage.page.locator("h1");
        await expect(title).toHaveText("Thunderdome Empower Your Agile Teams");
      });
    });

    test.describe("Non Admin Registered User", { tag: ["@registered"] }, () => {
      test("redirects to landing", async ({ registeredPage }) => {
        const adminPage = new AdminOrganizationPage(registeredPage.page);

        await adminPage.goto("bbaf82ef-a2d3-4e9a-b824-5e56a03ac3aa");

        const title = adminPage.page.locator("h1");
        await expect(title).toHaveText("Thunderdome Empower Your Agile Teams");
      });
    });

    test.describe("Admin User", { tag: ["@admin"] }, () => {
      test("loads Organization page", async ({ registeredPage, adminPage }) => {
        const ap = new AdminOrganizationPage(adminPage.page);
        const testOrgName = "E2E TEST ADMIN ORGANIZATION";
        const org = await registeredPage.createOrg(testOrgName);

        await ap.goto(org.id);
        await expect(ap.page.locator("h1")).toContainText(testOrgName);
      });
    });
  },
);
