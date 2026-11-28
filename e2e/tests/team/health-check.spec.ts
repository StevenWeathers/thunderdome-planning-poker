import { expect, test } from "@fixtures/user-sessions";
import { HealthChecksPage } from "@fixtures/pages/health-checks-page";

test.describe("Team Health Check pages", { tag: "@healthcheck" }, () => {
  test.describe("Unauthenticated user", { tag: "@unauthenticated" }, () => {
    test("redirects to login from a shared health check link", async ({ page }) => {
      const healthChecksPage = new HealthChecksPage(page);
      await healthChecksPage.gotoHealthCheck(
        "bbaf82ef-a2d3-4e9a-b824-5e56a03ac3aa",
        "8c1e7a52-2b0f-4f3e-9a57-0d9d1c8e4b11",
      );

      await expect(page.locator('form[name="login"]')).toBeVisible();
      await expect(page).toHaveURL(/\/login\/health-check\//);
    });
  });

  test.describe("Registered user", { tag: "@registered" }, () => {
    test("creates, responds to and closes a health check", async ({ registeredPage }) => {
      const team = await registeredPage.createTeam("E2E HEALTH CHECK TEAM");
      const healthChecksPage = new HealthChecksPage(registeredPage.page);
      const page = healthChecksPage.page;

      await healthChecksPage.goto(team.id);
      await healthChecksPage.create("E2E Sprint Health");

      await expect(page.locator("h1")).toHaveText("E2E Sprint Health");
      await expect(page.getByTestId("health-check-pending")).toBeVisible();

      await healthChecksPage.respondToAll("Good");
      await expect(page.getByTestId("health-check-response-count")).toHaveText("Responses: 1");
      await expect(page.getByTestId("health-check-pending")).toBeHidden();

      await healthChecksPage.close();
      await expect(page.getByText("fewer than 3 people responded")).toBeVisible();
      await expect(page.getByTestId("health-check-submit")).toBeHidden();

      await healthChecksPage.goto(team.id);
      await expect(page.getByTestId("health-check-name")).toHaveText("E2E Sprint Health");
    });
  });
});
