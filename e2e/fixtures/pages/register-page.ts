import { Locator, Page } from "@playwright/test";

export class RegisterPage {
  readonly page: Page;
  readonly userNameField: Locator;
  readonly createFullAccountCheckbox: Locator;
  readonly userEmailField: Locator;
  readonly userPassword1Field: Locator;
  readonly userPassword2Field: Locator;

  constructor(page: Page) {
    this.page = page;
    this.userNameField = page.locator('[name="yourName"]');
    this.userEmailField = page.locator('[name="yourEmail"]');
    this.userPassword1Field = page.locator('[name="yourPassword1"]');
    this.userPassword2Field = page.locator('[name="yourPassword2"]');
    this.createFullAccountCheckbox = page.locator('[for="createFullAccount"]');
  }

  async goto() {
    await this.page.goto("/register");
  }

  async createGuestUser(name) {
    await this.userNameField.fill(name);
    await this.userNameField.press("Enter");
  }

  async createRegisteredUser(name, email, password1, password2) {
    await this.userNameField.fill(name);
    await this.createFullAccountCheckbox.click();
    await this.userEmailField.fill(email);
    await this.userPassword1Field.fill(password1);
    await this.userPassword2Field.fill(password2);
    await this.userPassword2Field.press("Enter");
  }

  async createRegisteredUserFromGuest(email, password1, password2) {
    await this.userEmailField.fill(email);
    await this.userPassword1Field.fill(password1);
    await this.userPassword2Field.fill(password2);
    await this.userPassword2Field.press("Enter");
  }
}
