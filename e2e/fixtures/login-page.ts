import { Locator, Page } from '@playwright/test';

export class LoginPage {
  readonly page: Page;
  readonly emailField: Locator;
  readonly passwordField: Locator;

  constructor(page: Page) {
    this.emailField = page.getByPlaceholder('Enter your email');
    this.passwordField = page.getByPlaceholder('Enter a password');
    this.page = page;
  }

  async goto() {
    await this.page.goto(`/login`);
  }

  async login(email, password) {
    await this.emailField.fill(email);
    await this.passwordField.fill(password);
    await this.page.getByRole('button', { name: 'Login' }).click();
  }
}
