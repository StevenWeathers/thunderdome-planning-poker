import { expect, test } from "@fixtures/user-sessions";
import { AdminPage } from "@fixtures/admin/admin-page";

test.describe("Admin page", { tag: ["@administration"] }, () => {
  test.describe("Unauthenticated user", { tag: ["@unauthenticated"] }, () => {
    test("redirects to login", async ({ page }) => {
      const adminPage = new AdminPage(page);

      await adminPage.goto();

      const loginForm = adminPage.page.locator('form[name="login"]');
      await expect(loginForm).toBeVisible();
    });
  });

  test.describe("Guest user", { tag: ["@guest"] }, () => {
    test("redirects to landing", async ({ guestPage }) => {
      const adminPage = new AdminPage(guestPage.page);

      await adminPage.goto();

      const title = adminPage.page.locator("h1");
      await expect(title).toHaveText("Thunderdome Empower Your Agile Teams");
    });
  });

  test.describe("Non Admin Registered User", { tag: ["@registered"] }, () => {
    test("redirects to landing", async ({ registeredPage }) => {
      const adminPage = new AdminPage(registeredPage.page);

      await adminPage.goto();

      const title = adminPage.page.locator("h1");
      await expect(title).toHaveText("Thunderdome Empower Your Agile Teams");
    });
  });

  test.describe("Admin User", { tag: ["@admin"] }, () => {
    test("loads admin page", async ({ adminPage }) => {
      const ap = new AdminPage(adminPage.page);

      await ap.goto();

      // admin nav items are present
      const nav = ap.page.locator('[data-testid="admin-nav-item"]');
      await expect(nav.nth(0)).toHaveText("Admin");
      await expect(nav.nth(1)).toHaveText("Alerts");
      await expect(nav.nth(2)).toHaveText("Games");
      await expect(nav.nth(3)).toHaveText("Retros");
      await expect(nav.nth(4)).toHaveText("Storyboards");
      await expect(nav.nth(5)).toHaveText("Organizations");
      await expect(nav.nth(6)).toHaveText("Teams");
      await expect(nav.nth(7)).toHaveText("Users");
      await expect(nav.nth(8)).toHaveText("API Keys");
    });
  });
});
