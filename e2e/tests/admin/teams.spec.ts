import { expect, test } from "../../fixtures/user-sessions";
import { AdminTeamsPage } from "../../fixtures/admin/teams-page";

test.describe("The Admin Teams Page", () => {
  test.describe("Unauthenticated user", () => {
    test("redirects to login", async ({ page }) => {
      const adminPage = new AdminTeamsPage(page);

      await adminPage.goto();

      const loginForm = adminPage.page.locator('form[name="login"]');
      await expect(loginForm).toBeVisible();
    });
  });

  test.describe("Guest user", () => {
    test("redirects to landing", async ({ guestPage }) => {
      const adminPage = new AdminTeamsPage(guestPage.page);

      await adminPage.goto();

      const title = adminPage.page.locator("h1");
      await expect(title).toHaveText("Thunderdome: Empower Your Agile Teams");
    });
  });

  test.describe("Non Admin Registered User", () => {
    test("redirects to landing", async ({ registeredPage }) => {
      const adminPage = new AdminTeamsPage(registeredPage.page);

      await adminPage.goto();

      const title = adminPage.page.locator("h1");
      await expect(title).toHaveText("Thunderdome: Empower Your Agile Teams");
    });
  });

  test.describe("Admin User", () => {
    test("loads Teams page", async ({ adminPage }) => {
      const ap = new AdminTeamsPage(adminPage.page);

      await ap.goto();

      const title = ap.page.locator('[data-testid="tablenav-title"]');
      await expect(title).toHaveText("Teams");
    });
  });
});
