import { test } from "@fixtures/user-sessions";
import { expect } from "@playwright/test";
import { StoryboardsPage } from "@fixtures/pages/storyboards-page";

const pageTitle = "My Storyboards";

test.describe("Storyboards page", { tag: ["@storyboard"] }, () => {
  test.describe("Unauthenticated user", { tag: ["@unauthenticated"] }, () => {
    test("displays storyboard landing page", async ({ page }) => {
      const storyboardsPage = new StoryboardsPage(page);
      await storyboardsPage.goto();

      const title = storyboardsPage.page.locator("h1");
      await expect(title).toHaveText("Story Mapping with Thunderdome");
    });
  });

  test.describe("Guest user", { tag: ["@guest"] }, () => {
    test("should load page", async ({ guestPage }) => {
      const storyboardsPage = new StoryboardsPage(guestPage.page);
      await storyboardsPage.goto();
      const title = storyboardsPage.page.locator("h1");
      await expect(title).toHaveText(pageTitle);
    });

    test("should allow creating a storyboard", async ({ guestPage }) => {
      const storyboardName = "Test Storyboard";
      const storyboardsPage = new StoryboardsPage(guestPage.page);
      await storyboardsPage.goto();

      await storyboardsPage.createStoryboard({ name: storyboardName });

      const storyboardTitle = storyboardsPage.page.locator("h1");
      await expect(storyboardTitle).toHaveText(storyboardName);
    });

    test("should display storyboards", async ({ guestPage }) => {
      const storyboardName = "Test Display Storyboard";

      const storyboardsPage = new StoryboardsPage(guestPage.page);
      await storyboardsPage.goto();

      await storyboardsPage.createStoryboard({ name: storyboardName });

      const storyboardTitle = storyboardsPage.page.locator("h1");
      await expect(storyboardTitle).toHaveText(storyboardName);

      await storyboardsPage.goto();

      const title = await storyboardsPage.storyboardCardName.filter({
        hasText: storyboardName,
      });
      await expect(title).toBeVisible();
    });
  });

  test.describe("Registered user", { tag: ["@registered"] }, () => {
    test("should load page", async ({ registeredPage }) => {
      const storyboardsPage = new StoryboardsPage(registeredPage.page);
      await storyboardsPage.goto();
      const title = storyboardsPage.page.locator("h1");
      await expect(title).toHaveText(pageTitle);
    });

    test("should allow creating a storyboard", async ({ registeredPage }) => {
      const storyboardName = "Test Storyboard";
      const storyboardsPage = new StoryboardsPage(registeredPage.page);
      await storyboardsPage.goto();

      await storyboardsPage.createStoryboard({ name: storyboardName });

      const storyboardTitle = storyboardsPage.page.locator("h1");
      await expect(storyboardTitle).toHaveText(storyboardName);
    });

    test("should display storyboards", async ({ registeredPage }) => {
      const storyboardName = "Test Display Storyboard";

      const storyboardsPage = new StoryboardsPage(registeredPage.page);
      await storyboardsPage.goto();

      await storyboardsPage.createStoryboard({ name: storyboardName });

      const storyboardTitle = storyboardsPage.page.locator("h1");
      await expect(storyboardTitle).toHaveText(storyboardName);

      await storyboardsPage.goto();

      const title = await storyboardsPage.storyboardCardName.filter({
        hasText: storyboardName,
      });
      await expect(title).toBeVisible();
    });
  });
});
