import { expect, test } from "@fixtures/user-sessions";
import { APIKeysPage } from "@fixtures/admin/apikeys-page";

test.describe(
  "The Admin API Keys Page",
  { tag: ["@administration", "@apikeys"] },
  () => {
    test.describe("Unauthenticated user", () => {
      test("redirects to login", async ({ page }) => {
        const adminPage = new APIKeysPage(page);

        await adminPage.goto();

        const loginForm = adminPage.page.locator('form[name="login"]');
        await expect(loginForm).toBeVisible();
      });
    });

    test.describe("Guest user", { tag: ["@guest"] }, () => {
      test("redirects to landing", async ({ guestPage }) => {
        const adminPage = new APIKeysPage(guestPage.page);

        await adminPage.goto();

        const title = adminPage.page.locator("h1");
        await expect(title).toHaveText("Thunderdome Empower Your Agile Teams");
      });
    });

    test.describe("Non Admin Registered User", { tag: ["@registered"] }, () => {
      test("redirects to landing", async ({ registeredPage }) => {
        const adminPage = new APIKeysPage(registeredPage.page);

        await adminPage.goto();

        const title = adminPage.page.locator("h1");
        await expect(title).toHaveText("Thunderdome Empower Your Agile Teams");
      });
    });

    test.describe("Admin User", { tag: ["@admin"] }, () => {
      test("loads API Keys page", async ({ adminPage }) => {
        const ap = new APIKeysPage(adminPage.page);

        await ap.goto();

        const title = ap.page.locator('[data-testid="tablenav-title"]');
        await expect(title).toHaveText("API Keys");
      });
    });
  },
);
