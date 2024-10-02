import { expect, test } from "@fixtures/user-sessions";
import { ProfilePage } from "@fixtures/pages/profile-page";

test.describe("User Profile page", { tag: ["@user"] }, () => {
  test(
    "Unauthenticated user redirects to login",
    { tag: "@unauthenticated" },
    async ({ page }) => {
      const profilePage = new ProfilePage(page);
      await profilePage.goto();

      const loginForm = profilePage.page.locator('form[name="login"]');
      await expect(loginForm).toBeVisible();
    },
  );

  test(
    "Guest user successfully loads",
    { tag: "@guest" },
    async ({ guestPage }) => {
      const profilePage = new ProfilePage(guestPage.page);
      await profilePage.goto();

      await expect(profilePage.page.locator("h1")).toHaveText("Your Profile");
      await expect(profilePage.page.locator("[name=yourName]")).toHaveValue(
        guestPage.user.name,
      );
      await expect(profilePage.page.locator("[name=yourEmail]")).toHaveValue(
        "",
      );

      await expect(
        profilePage.page.locator('[data-testid="user-verified"]'),
      ).not.toBeVisible();
      await expect(
        profilePage.page.locator('[data-testid="request-verify"]'),
      ).not.toBeVisible();

      await expect(
        profilePage.page.locator('[data-testid="toggle-updatepassword"]'),
      ).not.toBeVisible();
    },
  );

  test(
    "Guest user cannot create API keys",
    { tag: "@guest" },
    async ({ guestPage }) => {
      const profilePage = new ProfilePage(guestPage.page);
      await profilePage.goto();

      await expect(
        profilePage.page.locator("h2").filter({ hasText: "API Keys" }),
      ).toBeVisible();
      await profilePage.page.locator('[data-testid="apikey-create"]').click();
      await profilePage.page
        .locator("[name=keyName]")
        .fill("Create API Key Test");
      await profilePage.page
        .locator("[name=createApiKey] [type=submit]")
        .click();

      await expect(
        profilePage.page.locator('[data-testid="notification-msg"]'),
      ).toContainText("Only verified registered users can create API keys.");
      await expect(profilePage.page.locator("[name=keyName]")).toBeVisible();
    },
  );

  test(
    "Registered user successfully loads",
    { tag: "@registered" },
    async ({ registeredPage }) => {
      const profilePage = new ProfilePage(registeredPage.page);
      await profilePage.goto();

      await expect(profilePage.page.locator("h1")).toHaveText("Your Profile");
      await expect(profilePage.page.locator("[name=yourName]")).toHaveValue(
        registeredPage.user.name,
      );
      await expect(profilePage.page.locator("[name=yourEmail]")).toHaveValue(
        registeredPage.user.email,
      );

      await expect(
        profilePage.page.locator('[data-testid="user-verified"]'),
      ).not.toBeVisible();
      await expect(
        profilePage.page.locator('[data-testid="request-verify"]'),
      ).toBeVisible();

      await expect(
        profilePage.page.locator('[data-testid="toggle-updatepassword"]'),
      ).toBeVisible();
    },
  );

  test(
    "Registered non verified user cannot create API keys",
    { tag: "@registered" },
    async ({ registeredPage }) => {
      const profilePage = new ProfilePage(registeredPage.page);
      await profilePage.goto();

      await expect(
        profilePage.page.locator("h2").filter({ hasText: "API Keys" }),
      ).toBeVisible();
      await profilePage.page.locator('[data-testid="apikey-create"]').click();
      await profilePage.page
        .locator("[name=keyName]")
        .fill("Create API Key Test");
      await profilePage.page
        .locator("[name=createApiKey] [type=submit]")
        .click();

      await expect(
        profilePage.page.locator('[data-testid="notification-msg"]'),
      ).toContainText("Only verified registered users can create API keys.");
      await expect(profilePage.page.locator("[name=keyName]")).toBeVisible();
    },
  );

  test(
    "Verified user should have verified status next to email field label",
    { tag: ["@registered", "@verified"] },
    async ({ verifiedPage }) => {
      const profilePage = new ProfilePage(verifiedPage.page);
      await profilePage.goto();

      await expect(profilePage.page.locator("[name=yourEmail]")).toHaveValue(
        verifiedPage.user.email,
      );

      await expect(
        profilePage.page.locator('[data-testid="user-verified"]'),
      ).toBeVisible();
      await expect(
        profilePage.page.locator('[data-testid="request-verify"]'),
      ).not.toBeVisible();
    },
  );

  test(
    "Guest User can update profile",
    { tag: "@guest" },
    async ({ guestPage }) => {
      const testCompanyName = "Test Update Company Guest";
      const testJobTitle = "Test Engineer Guest";
      const profilePage = new ProfilePage(guestPage.page);
      await profilePage.goto();

      await expect(profilePage.page.locator("[name=yourCompany]")).toHaveValue(
        "",
      );

      await profilePage.page.locator("[name=yourCountry]").selectOption("US");
      await profilePage.page
        .locator("[name=yourCompany]")
        .fill(testCompanyName);
      await profilePage.page.locator("[name=yourJobTitle]").fill(testJobTitle);

      await profilePage.page
        .locator("[name=updateProfile] [type=submit]")
        .click();

      await expect(profilePage.page.locator("[name=yourCountry]")).toHaveValue(
        "US",
      );
      await expect(profilePage.page.locator("[name=yourCompany]")).toHaveValue(
        testCompanyName,
      );
      await expect(profilePage.page.locator("[name=yourJobTitle]")).toHaveValue(
        testJobTitle,
      );

      await expect(
        profilePage.page.locator(`[data-testid="notification-msg"]`),
      ).toHaveText("Profile updated");
    },
  );

  test(
    "Registered User can update profile",
    { tag: "@registered" },
    async ({ registeredPage }) => {
      const testCompanyName = "Test Update Company";
      const testJobTitle = "Test Engineer";
      const profilePage = new ProfilePage(registeredPage.page);
      await profilePage.goto();

      await expect(profilePage.page.locator("[name=yourCompany]")).toHaveValue(
        "",
      );

      await profilePage.page.locator("[name=yourCountry]").selectOption("US");
      await profilePage.page
        .locator("[name=yourCompany]")
        .fill(testCompanyName);
      await profilePage.page.locator("[name=yourJobTitle]").fill(testJobTitle);

      await profilePage.page
        .locator("[name=updateProfile] [type=submit]")
        .click();

      await expect(profilePage.page.locator("[name=yourCountry]")).toHaveValue(
        "US",
      );
      await expect(profilePage.page.locator("[name=yourCompany]")).toHaveValue(
        testCompanyName,
      );
      await expect(profilePage.page.locator("[name=yourJobTitle]")).toHaveValue(
        testJobTitle,
      );

      await expect(
        profilePage.page.locator(`[data-testid="notification-msg"]`),
      ).toHaveText("Profile updated");
    },
  );

  test(
    "Verified user should display existing API keys",
    { tag: ["@registered", "@verified"] },
    async ({ verifiedPage }) => {
      const apiKeyName = "Display API Keys Test";
      const profilePage = new ProfilePage(verifiedPage.page);
      const apk = await verifiedPage.createApikey(apiKeyName);

      await profilePage.goto();

      await expect(
        profilePage.page.locator("h2").filter({ hasText: "API Keys" }),
      ).toBeVisible();

      await expect(
        profilePage.page.locator(
          `[data-apikeyid="${apk.id}"] [data-testid="apikey-name"]`,
        ),
      ).toHaveText(apk.name);
      await expect(
        profilePage.page.locator(
          `[data-apikeyid="${apk.id}"] [data-testid="apikey-prefix"]`,
        ),
      ).toHaveText(apk.prefix);
      await expect(
        profilePage.page.locator(
          `[data-apikeyid="${apk.id}"] [data-testid="apikey-active"]`,
        ),
      ).toHaveAttribute("data-active", "true");
    },
  );

  test(
    "Verified user can create API keys",
    { tag: ["@registered", "@verified"] },
    async ({ verifiedPage }) => {
      const apiKeyName = "Create API Key Test";
      const profilePage = new ProfilePage(verifiedPage.page);
      await profilePage.goto();

      await expect(
        profilePage.page.locator("h2").filter({ hasText: "API Keys" }),
      ).toBeVisible();
      await profilePage.page.locator('[data-testid="apikey-create"]').click();
      await profilePage.page.locator("[name=keyName]").fill(apiKeyName);
      await profilePage.page
        .locator("[name=createApiKey] [type=submit]")
        .click();

      await expect(profilePage.page.locator('[id="apiKey"]')).toBeVisible();
      await profilePage.page.locator('[data-testid="apikey-close"]').click();

      await expect(
        profilePage.page.locator('[data-testid="apikey-name"]', {
          hasText: apiKeyName,
        }),
      ).toBeVisible();
    },
  );

  test(
    "Verified User can toggle api key active status",
    { tag: ["@registered", "@verified"] },
    async ({ verifiedPage }) => {
      const apiKeyName = "Toggle API Key Active Test";
      const profilePage = new ProfilePage(verifiedPage.page);
      const apk = await verifiedPage.createApikey(apiKeyName);

      await profilePage.goto();

      await expect(
        profilePage.page.locator("h2").filter({ hasText: "API Keys" }),
      ).toBeVisible();

      await expect(
        profilePage.page.locator(
          `[data-apikeyid="${apk.id}"] [data-testid="apikey-active"]`,
        ),
      ).toHaveAttribute("data-active", "true");

      await profilePage.page
        .locator(
          `[data-apikeyid="${apk.id}"] [data-testid="apikey-activetoggle"]`,
        )
        .click();

      await expect(
        profilePage.page.locator(
          `[data-apikeyid="${apk.id}"] [data-testid="apikey-active"]`,
        ),
      ).toHaveAttribute("data-active", "false");
    },
  );

  test(
    "Verified User can delete api key",
    { tag: ["@registered", "@verified"] },
    async ({ verifiedPage }) => {
      const apiKeyName = "Delete API Key Test";
      const profilePage = new ProfilePage(verifiedPage.page);
      const apk = await verifiedPage.createApikey(apiKeyName);

      await profilePage.goto();

      await expect(
        profilePage.page.locator("h2").filter({ hasText: "API Keys" }),
      ).toBeVisible();

      await expect(
        profilePage.page.locator(
          `[data-apikeyid="${apk.id}"] [data-testid="apikey-name"]`,
        ),
      ).toHaveText(apiKeyName);

      await profilePage.page
        .locator(`[data-apikeyid="${apk.id}"] [data-testid="apikey-delete"]`)
        .click();

      await expect(
        profilePage.page.locator(`[data-apikeyid="${apk.id}"]`),
      ).toHaveCount(0);
    },
  );

  test(
    "Verified User can create no more than 5 API keys (default for config)",
    { tag: ["@registered", "@verified"] },
    async ({ adminPage }) => {
      const apiKeyName = "Max API Key Test";
      const profilePage = new ProfilePage(adminPage.page);
      await adminPage.createApikey(`${apiKeyName} 1`);
      await adminPage.createApikey(`${apiKeyName} 2`);
      await adminPage.createApikey(`${apiKeyName} 3`);
      await adminPage.createApikey(`${apiKeyName} 4`);
      await adminPage.createApikey(`${apiKeyName} 5`);

      await profilePage.goto();

      await expect(
        profilePage.page.locator("h2").filter({ hasText: "API Keys" }),
      ).toBeVisible();

      await profilePage.page.locator('[data-testid="apikey-create"]').click();
      await profilePage.page.locator("[name=keyName]").fill(apiKeyName);
      await profilePage.page
        .locator("[name=createApiKey] [type=submit]")
        .click();

      await expect(
        profilePage.page.locator(`[data-testid="notification-msg"]`),
      ).toHaveText("You have the max number of API keys allowed.");

      await expect(profilePage.page.locator('[id="apiKey"]')).toHaveCount(0);
      await expect(
        profilePage.page.locator('[data-testid="apikey-close"]'),
      ).toHaveCount(0);
    },
  );

  test(
    "delete account confirmation cancel does not delete account",
    { tag: ["@registered"] },
    async ({ registeredPage }) => {
      const profilePage = new ProfilePage(registeredPage.page);
      await profilePage.goto();

      await expect(profilePage.page.locator("h1")).toHaveText("Your Profile");

      await profilePage.page
        .locator("button", { hasText: "Delete Account" })
        .click();
      await profilePage.page.locator("data-testid=confirm-cancel").click();

      await expect(profilePage.page.locator("h1")).toHaveText("Your Profile");
      await expect(
        profilePage.page.locator("data-testid=usernav-name"),
      ).toHaveText(registeredPage.user.name);
    },
  );

  test(
    "Guest User user can delete account",
    { tag: ["@guest"] },
    async ({ deleteGuestPage }) => {
      const profilePage = new ProfilePage(deleteGuestPage.page);
      await profilePage.goto();

      await expect(profilePage.page.locator("h1")).toHaveText("Your Profile");

      await expect(
        profilePage.page.locator("data-testid=usernav-name"),
      ).toBeVisible();

      await profilePage.page
        .locator("button", { hasText: "Delete Account" })
        .click();
      await profilePage.page.locator("data-testid=confirm-confirm").click();

      // should be on landing page and no longer authenticated
      await expect(profilePage.page.locator("h1")).toHaveText(
        "Thunderdome Empower Your Agile Teams",
      );
      await expect(
        profilePage.page.locator("data-testid=usernav-name"),
      ).not.toBeVisible();
    },
  );

  test(
    "Registered user can delete account",
    { tag: ["@registered"] },
    async ({ deleteRegisteredPage }) => {
      const profilePage = new ProfilePage(deleteRegisteredPage.page);
      await profilePage.goto();

      await expect(profilePage.page.locator("h1")).toHaveText("Your Profile");

      await expect(
        profilePage.page.locator("data-testid=usernav-name"),
      ).toBeVisible();

      await profilePage.page
        .locator("button", { hasText: "Delete Account" })
        .click();
      await profilePage.page.locator("data-testid=confirm-confirm").click();

      // should be on landing page and no longer authenticated
      await expect(profilePage.page.locator("h1")).toHaveText(
        "Thunderdome Empower Your Agile Teams",
      );
      await expect(
        profilePage.page.locator("data-testid=usernav-name"),
      ).not.toBeVisible();
    },
  );
});
