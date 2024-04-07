import { test } from "../fixtures/user-sessions";
import { expect } from "@playwright/test";
import { PokerGamesPage } from "../fixtures/poker-games-page";

const pageTitle = "My Battles";

test.describe("Poker Games page", () => {
  test.describe("Unauthenticated user", () => {
    test("redirects to login", async ({ page }) => {
      const battlesPage = new PokerGamesPage(page);
      await battlesPage.goto();

      const loginForm = battlesPage.page.locator('form[name="login"]');
      await expect(loginForm).toBeVisible();
    });
  });

  test.describe("Guest user", () => {
    test("should load page", async ({ guestPage }) => {
      const battlesPage = new PokerGamesPage(guestPage.page);
      await battlesPage.goto();
      const title = battlesPage.page.locator("h1");
      await expect(title).toHaveText(pageTitle);
    });

    test("should allow creating a game", async ({ guestPage }) => {
      const battleName = "Test Game";
      const battlesPage = new PokerGamesPage(guestPage.page);
      await battlesPage.goto();

      await battlesPage.createBattle({ name: battleName });

      const battleTitle = battlesPage.page.locator("h2");
      await expect(battleTitle).toHaveText(battleName);
    });

    test("should display games", async ({ guestPage }) => {
      const battleName = "Test Display Games";

      const battlesPage = new PokerGamesPage(guestPage.page);
      await battlesPage.goto();

      await battlesPage.createBattle({ name: battleName });

      const battleTitle = battlesPage.page.locator("h2");
      await expect(battleTitle).toHaveText(battleName);

      await battlesPage.goto();

      const title = await battlesPage.gameCardName.filter({
        hasText: battleName,
      });
      await expect(title).toBeVisible();
    });
  });

  test.describe("Registered user", () => {
    test("should load page", async ({ registeredPage }) => {
      const battlesPage = new PokerGamesPage(registeredPage.page);
      await battlesPage.goto();
      const title = battlesPage.page.locator("h1");
      await expect(title).toHaveText(pageTitle);
    });

    test("should allow creating a game", async ({ registeredPage }) => {
      const battleName = "Test Game";
      const battlesPage = new PokerGamesPage(registeredPage.page);
      await battlesPage.goto();

      await battlesPage.createBattle({ name: battleName });

      const battleTitle = battlesPage.page.locator("h2");
      await expect(battleTitle).toHaveText(battleName);
    });

    test.skip("should allow creating a game with stories", async ({
      registeredPage,
    }) => {
      const battleName = "Test Game with Stories";
      const story1 = "Test Story 1";
      const story2 = "Test Story 2";
      const battlesPage = new PokerGamesPage(registeredPage.page);
      await battlesPage.goto();

      await battlesPage.createBattleWithStories({
        name: battleName,
        stories: [story1, story2],
      });

      const battleTitle = battlesPage.page.locator("h2");
      await expect(battleTitle).toHaveText(battleName);
      await expect(battlesPage.page.getByText(story1)).toBeVisible();
      await expect(battlesPage.page.getByText(story2)).toBeVisible();
    });

    test("should display games", async ({ registeredPage }) => {
      const battleName = "Test Display Game";

      const battlesPage = new PokerGamesPage(registeredPage.page);
      await battlesPage.goto();

      await battlesPage.createBattle({ name: battleName });

      const battleTitle = battlesPage.page.locator("h2");
      await expect(battleTitle).toHaveText(battleName);

      await battlesPage.goto();

      const title = await battlesPage.gameCardName.filter({
        hasText: battleName,
      });
      await expect(title).toBeVisible();
    });
  });
});
