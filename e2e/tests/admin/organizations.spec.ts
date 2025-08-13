import { expect, test } from "@fixtures/user-sessions";
import { OrganizationsPage } from "@fixtures/admin/organizations-page";

test.describe(
  "The Admin Organizations Page",
  { tag: ["@administration", "@organization"] },
  () => {
    test.describe("Unauthenticated user", { tag: ["@unauthenticated"] }, () => {
      test("redirects to login", async ({ page }) => {
        const adminPage = new OrganizationsPage(page);

        await adminPage.goto();

        const loginForm = adminPage.page.locator('form[name="login"]');
        await expect(loginForm).toBeVisible();
      });
    });

    test.describe("Guest user", { tag: ["@guest"] }, () => {
      test("redirects to landing", async ({ guestPage }) => {
        const adminPage = new OrganizationsPage(guestPage.page);

        await adminPage.goto();

        const title = adminPage.page.locator("h1 + p");
        await expect(title).toHaveText(
          "Transform your agile ceremonies from time-wasters into team-builders. Get the tools that make planning poker, retrospectives, and story mapping actually work for remote and in-person teams.",
        );
      });
    });

    test.describe("Non Admin Registered User", { tag: ["@registered"] }, () => {
      test("redirects to landing", async ({ registeredPage }) => {
        const adminPage = new OrganizationsPage(registeredPage.page);

        await adminPage.goto();

        const title = adminPage.page.locator("h1 + p");
        await expect(title).toHaveText(
          "Transform your agile ceremonies from time-wasters into team-builders. Get the tools that make planning poker, retrospectives, and story mapping actually work for remote and in-person teams.",
        );
      });
    });

    test.describe("Admin User", { tag: ["@admin"] }, () => {
      test("loads Organizations page", async ({ adminPage }) => {
        const ap = new OrganizationsPage(adminPage.page);

        await ap.goto();

        const title = ap.page.locator('[data-testid="tablenav-title"]');
        await expect(title).toHaveText("Organizations");
      });
    });
  },
);
