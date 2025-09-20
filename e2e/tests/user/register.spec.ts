import { expect, test } from "@playwright/test";
import { RegisterPage } from "@fixtures/pages/register-page";

const registerPageTitle = "Register";
const battlesPageTitle = "My Games";

test.beforeEach(async ({ page }) => {
  await page.goto("/register");
});

test.describe("Register page", () => {
  test("should display", async ({ page }) => {
    const title = page.locator("h1");
    await expect(title).toHaveText(registerPageTitle);
  });

  test.describe("Guest User", () => {
    test("should allow guest user signup", async ({ context, page }) => {
      const request = context.request;
      const registerPage = new RegisterPage(page);

      try {
        await registerPage.createGuestUser("TestGuestUser");

        const dashboardTitle = page.locator("h1");
        await expect(dashboardTitle).toHaveText("Welcome back, TestGuestUser");
      } finally {
        const u = await request.get("/api/auth/user");
        const user = await u.json();
        await request.delete(`/api/users/${user.data.id}`);
      }
    });
  });

  test.describe("Registered User", () => {
    test("should allow user registration", async ({ context, page }) => {
      const userName = "Registered Test User";
      const userEmail = "registered@thunderdome.dev";
      const userPass = "testreguserpassword";
      const request = context.request;
      const registerPage = new RegisterPage(page);

      try {
        await registerPage.createRegisteredUser(
          userName,
          userEmail,
          userPass,
          userPass,
        );

        const dashboardTitle = page.locator("h1");
        await expect(dashboardTitle).toHaveText(
          "Welcome back, Registered Test User",
        );
      } finally {
        const u = await request.get("/api/auth/user");
        const user = await u.json();
        await request.delete(`/api/users/${user.data.id}`);
      }
    });

    test("should allow user registration from guest session", async ({
      context,
      page,
    }) => {
      const userName = "Registered From Guest Test User";
      const userEmail = "registeredfromguest@thunderdome.dev";
      const userPass = "testreguserpassword";
      const request = context.request;
      const registerPage = new RegisterPage(page);

      try {
        await registerPage.createGuestUser("TestGuestUser");

        const dashboardTitle = page.locator("h1");
        await expect(dashboardTitle).toHaveText("Welcome back, TestGuestUser");

        await registerPage.goto();

        const title = page.locator("h1");
        await expect(title).toHaveText(registerPageTitle);

        await registerPage.createRegisteredUserFromGuest(
          userEmail,
          userPass,
          userPass,
        );

        const dashboardTitle2 = page.locator("h1");
        await expect(dashboardTitle2).toHaveText("Welcome back, TestGuestUser");
      } finally {
        const u = await request.get("/api/auth/user");
        const user = await u.json();
        await request.delete(`/api/users/${user.data.id}`);
      }
    });
  });
});
