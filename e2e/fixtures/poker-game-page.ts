import { Locator, Page } from "@playwright/test";

export class PokerGamePage {
  readonly page: Page;
  readonly pageTitle: Locator;
  readonly toggleSpectator: Locator;
  readonly userDemoteBtn: Locator;
  readonly gameDeleteBtn: Locator;
  readonly gameDeleteConfirmBtn: Locator;
  readonly gameDeleteCancelBtn: Locator;
  readonly addStoriesBtn: Locator;
  readonly editStoryBtn: Locator;
  readonly deleteStoryBtn: Locator;
  readonly activateStoryBtn: Locator;
  readonly abandonGameBtn: Locator;
  readonly viewStoryBtn: Locator;
  readonly storyName: Locator;
  readonly storyType: Locator;
  readonly storyNameField: Locator;
  readonly storyTypeField: Locator;
  readonly saveStoryBtn: Locator;

  constructor(page: Page) {
    this.pageTitle = page.locator("h2");
    this.toggleSpectator = page.locator('[data-testid="user-togglespectator"]');
    this.userDemoteBtn = page.locator(`[data-testid="user-demote"]`);
    this.gameDeleteBtn = page.locator('[data-testid="battle-delete"]');
    this.gameDeleteConfirmBtn = page.locator("data-testid=confirm-confirm");
    this.gameDeleteCancelBtn = page.locator("data-testid=confirm-cancel");
    this.addStoriesBtn = page.locator('[data-testid="plans-add"]');
    this.editStoryBtn = page.locator('[data-testid="plan-edit"]');
    this.deleteStoryBtn = page.locator('[data-testid="plan-delete"]');
    this.activateStoryBtn = page.locator('[data-testid="plan-activate"]');
    this.abandonGameBtn = page.locator('[data-testid="battle-abandon"]');
    this.viewStoryBtn = page.locator('[data-testid="plan-view"]');
    this.storyName = page.locator("data-testid=plan-name");
    this.storyType = page.locator("data-testid=plan-type");
    this.storyNameField = page.locator("input[name=planName]");
    this.storyTypeField = page.locator("select[name=planType]");
    this.saveStoryBtn = page.locator("data-testid=plan-save");

    this.page = page;
  }

  async goto(id) {
    await this.page.goto(`/battle/${id}`);
  }

  async addPlan(name) {
    await this.addStoriesBtn.click();
    await this.storyNameField.fill(name);
    await this.saveStoryBtn.click();
  }
}
