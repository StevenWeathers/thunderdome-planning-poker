import { Page } from "@playwright/test";

export class AdminDepartmentPage {
  readonly page: Page;

  constructor(page: Page) {
    this.page = page;
  }

  async goto(orgId, deptId) {
    await this.page.goto(`/admin/organizations/${orgId}/department/${deptId}`);
  }
}
