import { expect, test } from "@fixtures/user-sessions";
import { LoginPage } from "@fixtures/pages/login-page";
import { registeredUser } from "@fixtures/db/registered-user";

test.describe("The Login Page", { tag: "@login" }, () => {
  test("should navigate to my games page and reflect name in header", async ({
    page,
  }) => {
    const loginPage = new LoginPage(page);
    await loginPage.goto();
    await loginPage.login(registeredUser.email, "kentRules!");
    await expect(loginPage.page.locator("h1")).toHaveText("My Games");

    // UI should reflect this user being logged in
    await expect(
      loginPage.page.locator('[data-testid="usernav-name"]'),
    ).toHaveText(registeredUser.name);
  });
});
