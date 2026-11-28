import { Page } from "@playwright/test";

export class HealthChecksPage {
  readonly page: Page;

  constructor(page: Page) {
    this.page = page;
  }

  async goto(teamId) {
    await this.page.goto(`/team/${teamId}/health-checks`);
  }

  async gotoHealthCheck(teamId, healthCheckId) {
    await this.page.goto(`/team/${teamId}/health-checks/${healthCheckId}`);
  }

  async create(name: string) {
    await this.page.getByTestId("health-check-create").click();
    await this.page.locator('input[name="healthCheckName"]').fill(name);
    await this.page.getByTestId("health-check-save").click();
  }

  async respondToAll(rating: string) {
    for (const radio of await this.page.getByRole("radio", { name: rating }).all()) {
      await radio.check();
    }
    await this.page.getByTestId("health-check-submit").click();
  }

  async close() {
    await this.page.getByTestId("health-check-close").click();
    await this.page.getByTestId("confirm-confirm").click();
  }
}
