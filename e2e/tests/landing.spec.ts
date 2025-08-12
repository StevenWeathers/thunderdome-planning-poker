import { expect, test } from "@playwright/test";

test.beforeEach(async ({ page }) => {
  await page.goto("/");
});

test("Landing Page", { tag: "@landing" }, async ({ page }) => {
  const title = page.locator("h1 + p");
  await expect(title).toHaveText(
    "Transform your agile ceremonies from time-wasters into team-builders. Get the tools that make planning poker, retrospectives, and story mapping actually work for remote and in-person teams.",
  );
});
