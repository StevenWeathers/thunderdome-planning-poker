import { test } from "@fixtures/user-sessions";
import { expect } from "@playwright/test";
import { TeamsPage } from "@fixtures/pages/teams-page";

test.describe("Teams page", { tag: ["@team"] }, () => {
  test.describe("Unauthenticated user", { tag: ["@unauthenticated"] }, () => {
    test("redirects to login", async ({ page }) => {
      const teamsPage = new TeamsPage(page);
      await teamsPage.goto();

      const loginForm = teamsPage.page.locator('form[name="login"]');
      await expect(loginForm).toBeVisible();
    });
  });

  test.describe("Registered user", { tag: ["@registered"] }, () => {
    test("successfully loads page", async ({ registeredPage }) => {
      const teamsPage = new TeamsPage(registeredPage.page);
      await teamsPage.goto();

      await expect(
        teamsPage.page.locator('[data-testid="tablenav-title"]', {
          hasText: "Organizations",
        }),
      ).toBeVisible();
      await expect(
        teamsPage.page.locator('[data-testid="tablenav-title"]', {
          hasText: "Teams",
        }),
      ).toBeVisible();
    });

    test.describe("Create Organization", { tag: ["@organization"] }, () => {
      test("should successfully submit and navigate to new organization page", async ({
        registeredPage,
      }) => {
        const teamsPage = new TeamsPage(registeredPage.page);
        await teamsPage.goto();

        await teamsPage.createOrganization({
          name: "Test Organization",
        });

        await expect(
          teamsPage.page.locator('[data-testid="tablenav-title"]', {
            hasText: "Departments",
          }),
        ).toBeVisible();
        await expect(
          teamsPage.page.locator('[data-testid="tablenav-title"]', {
            hasText: "Teams",
          }),
        ).toBeVisible();
        await expect(
          teamsPage.page.locator('[data-testid="tablenav-title"]', {
            hasText: "Users",
          }),
        ).toBeVisible();
      });
    });

    test.describe("Create Team", () => {
      test("should successfully submit and navigate to new team page", async ({
        registeredPage,
      }) => {
        const teamsPage = new TeamsPage(registeredPage.page);
        await teamsPage.goto();

        await teamsPage.createTeam({ name: "Test Team" });

        await expect(
          teamsPage.page.locator("h2", { hasText: "Games" }),
        ).toBeVisible();
        await expect(
          teamsPage.page.locator("h2", { hasText: "Retros" }),
        ).toBeVisible();
        await expect(
          teamsPage.page.locator("h2", { hasText: "Storyboards" }),
        ).toBeVisible();
        await expect(
          teamsPage.page.locator('[data-testid="tablenav-title"]', {
            hasText: "Users",
          }),
        ).toBeVisible();
      });
    });
  });
});
