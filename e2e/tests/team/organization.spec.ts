import { expect, test } from "@fixtures/user-sessions";
import { OrganizationPage } from "@fixtures/pages/organization-page";

test.describe("Organization Page", { tag: "@organization" }, () => {
  test(
    "Unauthenticated user redirects to login",
    { tag: "@unauthenticated" },
    async ({ page }) => {
      const orgPage = new OrganizationPage(page);
      await orgPage.goto("bbaf82ef-a2d3-4e9a-b824-5e56a03ac3aa");

      const loginForm = orgPage.page.locator('form[name="login"]');
      await expect(loginForm).toBeVisible();
    },
  );

  test(
    "Registered user loads page successfully",
    { tag: "@registered" },
    async ({ registeredPage }) => {
      const testOrgName = "E2E TEST ORGANIZATION";
      const orgPage = new OrganizationPage(registeredPage.page);
      const org = await registeredPage.createOrg(testOrgName);

      await orgPage.goto(org.id);
      await expect(orgPage.page.locator("h1")).toContainText(testOrgName);
    },
  );

  test("can invite user to organization", async ({ registeredPage }) => {
    const verifiedEmail = "e2everified@thunderdome.dev";
    const testOrgName = "E2E TEST ORGANIZATION UA";
    const orgPage = new OrganizationPage(registeredPage.page);
    const org = await registeredPage.createOrg(testOrgName);

    await orgPage.goto(org.id);
    await expect(orgPage.page.locator("h1")).toContainText(testOrgName);

    await orgPage.page.locator('[data-testid="user-add"]').click();
    await orgPage.page.locator('input[name="userEmail"]').fill(verifiedEmail);
    await orgPage.page
      .locator('select[name="userRole"]')
      .selectOption("MEMBER");
    await orgPage.page.locator('[data-testid="useradd-confirm"]').click();

    await expect(
      orgPage.page.locator('[data-testid="invite-user-email"]', {
        hasText: verifiedEmail,
      }),
    ).toBeVisible();
  });
});
