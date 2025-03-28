import { expect, test } from "@fixtures/user-sessions";
import { AdminTeamPage } from "@fixtures/admin/team-page";

test.describe(
  "The Admin Team Page",
  { tag: ["@administration", "@team"] },
  () => {
    test.describe("Unauthenticated user", { tag: ["@unauthenticated"] }, () => {
      test("redirects to login", async ({ page }) => {
        const adminPage = new AdminTeamPage(page);

        await adminPage.goto("bbaf82ef-a2d3-4e9a-b824-5e56a03ac3aa");

        const loginForm = adminPage.page.locator('form[name="login"]');
        await expect(loginForm).toBeVisible();
      });
    });

    test.describe("Guest user", { tag: ["@guest"] }, () => {
      test("redirects to landing", async ({ guestPage }) => {
        const adminPage = new AdminTeamPage(guestPage.page);

        await adminPage.goto("bbaf82ef-a2d3-4e9a-b824-5e56a03ac3aa");

        const title = adminPage.page.locator("h1");
        await expect(title).toHaveText(
          "Elevate your agile practices, foster seamless collaboration, and unlock your team's full potential with our innovative suite of tools.",
        );
      });
    });

    test.describe("Non Admin Registered User", { tag: ["@registered"] }, () => {
      test("redirects to landing", async ({ registeredPage }) => {
        const adminPage = new AdminTeamPage(registeredPage.page);

        await adminPage.goto("bbaf82ef-a2d3-4e9a-b824-5e56a03ac3aa");

        const title = adminPage.page.locator("h1");
        await expect(title).toHaveText(
          "Elevate your agile practices, foster seamless collaboration, and unlock your team's full potential with our innovative suite of tools.",
        );
      });
    });

    test.describe("Admin User", { tag: ["@admin"] }, () => {
      test("loads Team page", async ({ registeredPage, adminPage }) => {
        const ap = new AdminTeamPage(adminPage.page);
        const testTeamName = "E2E TEST ADMIN TEAM";
        const team = await registeredPage.createTeam(testTeamName);

        await ap.goto(team.id);
        await expect(
          ap.page.locator('[data-testid="tablenav-title"]').nth(0),
        ).toContainText(testTeamName);
      });
    });
  },
);
