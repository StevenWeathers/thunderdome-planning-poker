import { Locator, Page } from "@playwright/test";

export class PokerGamesPage {
  readonly page: Page;
  readonly gameNameFormField: Locator;
  readonly gameCardName: Locator;

  constructor(page: Page) {
    this.page = page;
    this.gameNameFormField = page.locator(
      'form[name="createBattle"] [name="battleName"]',
    );
    this.gameCardName = page.locator('[data-testid="battle-name"]');
  }

  async goto() {
    await this.page.goto("/battles");
  }

  async createBattle({ name }) {
    await this.gameNameFormField.fill(name);
    await this.gameNameFormField.press("Enter");
  }

  async createBattleWithStories({ name, stories }) {
    await this.gameNameFormField.fill(name);

    for (const story of stories) {
      await this.page
        .getByRole("button", {
          name: "Add Plan",
        })
        .click();
      await this.page.getByPlaceholder("Enter a plan name").first().fill(story);
    }

    await this.gameNameFormField.press("Enter");
  }
}
