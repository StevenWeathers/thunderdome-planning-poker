import { expect, test } from "@fixtures/user-sessions";
import { AdminGamesPage } from "@fixtures/admin/games-page";

test.describe(
  "The Admin Poker Games Page",
  { tag: ["@administration", "@poker"] },
  () => {
    test.describe("Unauthenticated user", { tag: ["@unauthenticated"] }, () => {
      test("redirects to login", async ({ page }) => {
        const adminPage = new AdminGamesPage(page);

        await adminPage.goto();

        const loginForm = adminPage.page.locator('form[name="login"]');
        await expect(loginForm).toBeVisible();
      });
    });

    test.describe("Guest user", { tag: ["@guest"] }, () => {
      test("redirects to landing", async ({ guestPage }) => {
        const adminPage = new AdminGamesPage(guestPage.page);

        await adminPage.goto();

        const title = adminPage.page.locator("h1 + p");
        await expect(title).toHaveText(
          "Elevate your agile practices, foster seamless collaboration, and unlock your team's full potential with our innovative suite of tools.",
        );
      });
    });

    test.describe("Non Admin Registered User", { tag: ["@registered"] }, () => {
      test("redirects to landing", async ({ registeredPage }) => {
        const adminPage = new AdminGamesPage(registeredPage.page);

        await adminPage.goto();

        const title = adminPage.page.locator("h1 + p");
        await expect(title).toHaveText(
          "Elevate your agile practices, foster seamless collaboration, and unlock your team's full potential with our innovative suite of tools.",
        );
      });
    });

    test.describe("Admin User", { tag: ["@admin"] }, () => {
      test("loads Games page", async ({ adminPage }) => {
        const ap = new AdminGamesPage(adminPage.page);

        await ap.goto();

        const title = ap.page.locator('[data-testid="tablenav-title"]');
        await expect(title).toHaveText("Games");
      });
    });
  },
);
