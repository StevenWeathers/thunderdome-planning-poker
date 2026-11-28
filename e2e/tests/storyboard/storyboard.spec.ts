import { expect, test } from "@fixtures/user-sessions";
import { StoryboardPage } from "@fixtures/pages/storyboard-page";

test.describe("Storyboard page", { tag: ["@storyboard"] }, () => {
  let storyboard = { id: "", name: "e2e storyboard page tests" };
  let storyboardLeave = { id: "" };
  let storyboardCancelDelete = { id: "" };
  let storyboardDelete = { id: "" };

  test.beforeAll(async ({ registeredPage, verifiedPage, adminPage }) => {
    const commonStoryboard = {
      storyboardName: `${storyboard.name}`,
      storyboardFacilitators: [`${adminPage.user.email}`],
    };
    storyboard = await registeredPage.createStoryboard({
      ...commonStoryboard,
    });
    storyboardLeave = await verifiedPage.createStoryboard({
      ...commonStoryboard,
    });
    storyboardCancelDelete = await registeredPage.createStoryboard({
      ...commonStoryboard,
    });
    storyboardDelete = await registeredPage.createStoryboard({
      ...commonStoryboard,
    });
  });

  test(
    "unauthenticated user redirects to register",
    { tag: ["@unauthenticated"] },
    async ({ page }) => {
      const bp = new StoryboardPage(page);
      await bp.goto(storyboard.id);

      const title = bp.page.locator("h1");
      await expect(title).toHaveText("Register");
    },
  );

  test("guest user successfully loads", { tag: ["@guest"] }, async ({ guestPage }) => {
    const bp = new StoryboardPage(guestPage.page);
    await bp.goto(storyboard.id);

    await expect(bp.storyboardTitle).toContainText(storyboard.name);
  });

  test(
    "registered user successfully loads",
    { tag: ["@registered"] },
    async ({ registeredPage }) => {
      const bp = new StoryboardPage(registeredPage.page);
      await bp.goto(storyboard.id);

      await expect(bp.storyboardTitle).toContainText(storyboard.name);
    },
  );

  test("user can leave storyboard", async ({ registeredPage }) => {
    const bp = new StoryboardPage(registeredPage.page);
    await bp.goto(storyboardLeave.id);

    await bp.page.click('[data-testid="storyboard-settings"]');
    await bp.page.click('[data-testid="storyboard-leave"]');
    await expect(bp.page.locator("h1")).toHaveText("My Storyboards");
  });

  test("facilitator can reorder goals", async ({ registeredPage }) => {
    const { id } = await registeredPage.createStoryboard({ storyboardName: "e2e goal order" });
    for (const name of ["First", "Second", "Third"]) {
      const res = await registeredPage.page.request.post(`/api/storyboards/${id}/goals`, {
        data: { name },
      });
      expect(res.ok()).toBeTruthy();
    }

    const bp = new StoryboardPage(registeredPage.page);
    await bp.goto(id);

    const goals = bp.page.getByTestId("goal-order-item");
    const startReorder = async () => {
      await bp.page.getByTestId("goal-settings").first().click();
      await bp.page.getByTestId("toggle-edit-goal-order").click();
    };

    await startReorder();
    await expect(goals).toHaveText(["First", "Second", "Third"]);

    // keyboard drag and drop: pick up the last goal, move it up twice, then drop it
    await goals.nth(2).focus();
    await bp.page.keyboard.press("Enter");
    await bp.page.keyboard.press("ArrowUp");
    await expect(goals).toHaveText(["First", "Third", "Second"]);
    await bp.page.keyboard.press("ArrowUp");
    await expect(goals).toHaveText(["Third", "First", "Second"]);
    await bp.page.keyboard.press("Enter");

    // the new order is persisted
    await bp.page.reload();
    await startReorder();
    await expect(goals).toHaveText(["Third", "First", "Second"]);

    await bp.page.getByTestId("goal-order-done").click();
    await expect(bp.page.getByTestId("storyboard-goal")).toContainText([
      "Third",
      "First",
      "Second",
    ]);
  });

  test("delete storyboard confirmation cancel does not delete storyboard", async ({
    registeredPage,
  }) => {
    const bp = new StoryboardPage(registeredPage.page);
    await bp.goto(storyboardCancelDelete.id);

    await bp.page.click('[data-testid="storyboard-settings"]');
    await bp.storyboardDeleteBtn.click();
    await bp.storyboardDeleteCancelBtn.click();

    await expect(bp.storyboardTitle).toContainText(storyboard.name);
  });

  test("delete storyboard confirmation confirm deletes storyboard and redirects to storyboards page", async ({
    registeredPage,
  }) => {
    const bp = new StoryboardPage(registeredPage.page);
    await bp.goto(storyboardDelete.id);

    await bp.page.click('[data-testid="storyboard-settings"]');
    await bp.storyboardDeleteBtn.click();
    await bp.storyboardDeleteConfirmBtn.click();

    await expect(bp.page.locator("h1")).toHaveText("My Storyboards");
  });
});
