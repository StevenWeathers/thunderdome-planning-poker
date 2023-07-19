import { Locator, Page } from '@playwright/test';

export class RetroPage {
  readonly page: Page;
  readonly retroTitle: Locator;
  readonly retroEditBtn: Locator;
  readonly retroDeleteBtn: Locator;
  readonly retroDeleteConfirmBtn: Locator;
  readonly retroDeleteCancelBtn: Locator;
  readonly retroNextPhaseBtn: Locator;
  readonly retroPhasePrimeDirectiveBtn: Locator;
  readonly retroPhaseBrainstormBtn: Locator;
  readonly retroPhaseGroupBtn: Locator;
  readonly retroPhaseVoteBtn: Locator;
  readonly retroPhaseActionItemsBtn: Locator;
  readonly retroPhaseDoneBtn: Locator;
  readonly retroExportBtn: Locator;
  readonly retroWorkedWellInput: Locator;
  readonly retroNeedsImprovementInput: Locator;
  readonly retroQuestionInput: Locator;
  readonly retroGroupNameInput: Locator;
  readonly retroActionItemInput: Locator;

  constructor(page: Page) {
    this.retroTitle = page.locator('h1');
    this.retroEditBtn = page.getByRole('button', { name: 'Edit Retro' });
    this.retroDeleteBtn = page.getByRole('button', {
      name: 'Delete Retro',
    });
    this.retroDeleteConfirmBtn = page.locator('data-testid=confirm-confirm');
    this.retroDeleteCancelBtn = page.locator('data-testid=confirm-cancel');
    this.retroNextPhaseBtn = page.getByRole('button', {
      name: 'Next Phase',
    });
    this.retroPhasePrimeDirectiveBtn = page.getByRole('button', {
      name: 'Prime Directive',
    });
    this.retroPhaseBrainstormBtn = page.getByRole('button', {
      name: 'Brainstorm',
    });
    this.retroPhaseGroupBtn = page.getByRole('button', { name: 'Group' });
    this.retroPhaseVoteBtn = page.getByRole('button', { name: 'Vote' });
    this.retroPhaseActionItemsBtn = page.getByRole('button', {
      name: 'Action Items',
    });
    this.retroPhaseDoneBtn = page.getByRole('button', { name: 'Done' });
    this.retroExportBtn = page.getByRole('button', { name: 'Export' });
    this.retroWorkedWellInput = page.getByPlaceholder('What worked well...');
    this.retroNeedsImprovementInput = page.getByPlaceholder(
      'What needs improvement...',
    );
    this.retroQuestionInput = page.getByPlaceholder('I want to ask...');
    this.retroGroupNameInput = page.getByPlaceholder('Group Name');
    this.retroActionItemInput = page.getByPlaceholder('Action item...');

    this.page = page;
  }

  async goto(id) {
    await this.page.goto(`/retro/${id}`);
  }
}
