import { Page } from "@playwright/test";

export class AdminGamesPage {
  readonly page: Page;

  constructor(page: Page) {
    this.page = page;
  }

  async goto() {
    await this.page.goto(`/admin/battles`);
  }
}
