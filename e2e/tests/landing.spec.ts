import { expect, test } from "@playwright/test";

test.beforeEach(async ({ page }) => {
  await page.goto("/");
});

test("Landing Page", async ({ page }) => {
  const title = page.locator("h1");
  await expect(title).toHaveText("Thunderdome: Empower Your Agile Teams");
});
