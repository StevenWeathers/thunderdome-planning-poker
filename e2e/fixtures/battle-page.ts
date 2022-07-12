import {Page} from '@playwright/test';

export class BattlePage {
    readonly page: Page;

    constructor(page: Page) {
        this.page = page;
    }

    async goto(id) {
        await this.page.goto(`/battle/${id}`);
    }
}