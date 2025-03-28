import { expect, test } from "@playwright/test";

test.beforeEach(async ({ page }) => {
  await page.goto("/");
});

test("Landing Page", { tag: "@landing" }, async ({ page }) => {
  const title = page.locator("h1 + p");
  await expect(title).toHaveText(
    "Elevate your agile practices, foster seamless collaboration, and unlock your team's full potential with our innovative suite of tools.",
  );
});
