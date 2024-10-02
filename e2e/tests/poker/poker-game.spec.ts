import { expect, test } from "@fixtures/user-sessions";
import { PokerGamePage } from "@fixtures/pages/poker-game-page";

const allowedPointValues = ["0", "1", "2", "3", "5", "8", "13", "?"];
const pointAverageRounding = "ceil";
const lokiPlan = { name: "Defeat Loki", type: "Story" };
const thanosPlan = { name: "Defeat Thanos", type: "Epic" };
const scarletPlan = { name: "Defeat Scarlet Witch", type: "Epic" };

test.describe("Poker Game page", { tag: ["@poker"] }, () => {
  let poker = { id: "", name: "e2e poker page tests" };
  let pokerWithStory = { id: "" };
  let pokerAddStory = { id: "" };
  let pokerEditStory = { id: "" };
  let pokerDeleteStory = { id: "" };
  let pokerActivateStory = { id: "" };
  let pokerSkipStory = { id: "" };
  let pokerWithoutAutoVoting = { id: "" };
  let pokerWithAutoVoting = { id: "", name: "" };
  let pokerFinishVoting = { id: "" };
  let pokerSaveVoting = { id: "" };
  let pokerAbandon = { id: "" };
  let pokerCancelDelete = { id: "" };
  let pokerDelete = { id: "" };

  test.beforeAll(async ({ registeredPage, verifiedPage, adminPage }) => {
    const commonBattle = {
      name: `${poker.name}`,
      pointValuesAllowed: [...allowedPointValues],
      pointAverageRounding: `${pointAverageRounding}`,
      plans: [],
      autoFinishVoting: false,
      battleLeaders: [`${adminPage.user.email}`],
    };
    poker = await registeredPage.createPokerGame({
      ...commonBattle,
      plans: [lokiPlan],
    });
    pokerWithStory = await registeredPage.createPokerGame({
      ...commonBattle,
      plans: [lokiPlan],
    });
    pokerAddStory = await registeredPage.createPokerGame({
      ...commonBattle,
    });
    pokerEditStory = await registeredPage.createPokerGame({
      ...commonBattle,
      plans: [thanosPlan],
    });
    pokerDeleteStory = await registeredPage.createPokerGame({
      ...commonBattle,
      plans: [scarletPlan],
    });
    pokerActivateStory = await registeredPage.createPokerGame({
      ...commonBattle,
      plans: [scarletPlan],
    });
    pokerSkipStory = await registeredPage.createPokerGame({
      ...commonBattle,
      plans: [thanosPlan],
    });
    pokerWithoutAutoVoting = await registeredPage.createPokerGame({
      ...commonBattle,
    });
    pokerWithAutoVoting = await verifiedPage.createPokerGame({
      ...commonBattle,
      autoFinishVoting: true,
    });
    pokerFinishVoting = await registeredPage.createPokerGame({
      ...commonBattle,
      plans: [lokiPlan],
    });
    pokerSaveVoting = await registeredPage.createPokerGame({
      ...commonBattle,
      plans: [lokiPlan],
    });
    pokerAbandon = await verifiedPage.createPokerGame({ ...commonBattle });
    pokerCancelDelete = await registeredPage.createPokerGame({
      ...commonBattle,
    });
    pokerDelete = await registeredPage.createPokerGame({ ...commonBattle });
  });

  test(
    "unauthenticated user redirects to register",
    { tag: "@unauthenticated" },
    async ({ page }) => {
      const bp = new PokerGamePage(page);
      await bp.goto(poker.id);

      const title = bp.page.locator("h1");
      await expect(title).toHaveText("Register");
    },
  );

  test(
    "guest user successfully loads",
    { tag: "@guest" },
    async ({ guestPage }) => {
      const bp = new PokerGamePage(guestPage.page);
      await bp.goto(poker.id);

      await expect(bp.pageTitle).toHaveText(poker.name);
    },
  );

  test(
    "registered user successfully loads",
    { tag: "@registered" },
    async ({ registeredPage }) => {
      const bp = new PokerGamePage(registeredPage.page);
      await bp.goto(poker.id);

      await expect(bp.pageTitle).toHaveText(poker.name);
    },
  );

  test("user cannot become spectator when autoFinishVoting is false", async ({
    registeredPage,
  }) => {
    const bp = new PokerGamePage(registeredPage.page);
    await bp.goto(pokerWithoutAutoVoting.id);

    await expect(bp.toggleSpectator).not.toBeVisible();
  });

  test("user can become spectator when autoFinishVoting is true", async ({
    registeredPage,
  }) => {
    const bp = new PokerGamePage(registeredPage.page);
    await bp.goto(pokerWithAutoVoting.id);

    const spectatorButton = bp.toggleSpectator;

    await spectatorButton.click();
    await expect(spectatorButton).toHaveText("Become Participant");
  });

  // @TODO - update test now that only facilitator can't be removed
  test.skip("facilitator can remove facilitator", async ({ adminPage }) => {
    const bp = new PokerGamePage(adminPage.page);
    await bp.goto(poker.id);

    const userDemoteBtn = bp.page
      .locator(`[data-testid="user-card"][data-userid="${adminPage.user.id}"]`)
      .locator('[data-testid="user-demote"]');

    await expect(userDemoteBtn).toBeVisible();
    await expect(bp.gameDeleteBtn).toBeVisible();
    await expect(bp.addStoriesBtn).toBeVisible();
    await expect(bp.editStoryBtn).toBeVisible();
    await expect(bp.deleteStoryBtn).toBeVisible();
    await expect(bp.activateStoryBtn).toBeVisible();
    await expect(bp.viewStoryBtn).toBeVisible();
    await expect(bp.abandonGameBtn).not.toBeVisible();

    // yes you can demote yourself!
    await userDemoteBtn.click();

    await expect(userDemoteBtn).not.toBeVisible();
    await expect(bp.gameDeleteBtn).not.toBeVisible();
    await expect(bp.addStoriesBtn).not.toBeVisible();
    await expect(bp.editStoryBtn).not.toBeVisible();
    await expect(bp.deleteStoryBtn).not.toBeVisible();
    await expect(bp.activateStoryBtn).not.toBeVisible();
    await expect(bp.viewStoryBtn).toBeVisible();
    await expect(bp.abandonGameBtn).toBeVisible();
  });

  test("user can abandon game", async ({ registeredPage }) => {
    const bp = new PokerGamePage(registeredPage.page);
    await bp.goto(pokerAbandon.id);

    await bp.page.click('[data-testid="battle-abandon"]');
    await expect(bp.page.locator("h1")).toHaveText("My Games");
  });

  test("should display existing stories", async ({ registeredPage }) => {
    const bp = new PokerGamePage(registeredPage.page);
    await bp.goto(pokerWithStory.id);

    await expect(bp.storyName.filter({ hasText: lokiPlan.name })).toBeVisible();
  });

  test("should allow adding stories", async ({ registeredPage }) => {
    const newPlanName = "Defeat Thanos";
    const bp = new PokerGamePage(registeredPage.page);
    await bp.goto(pokerAddStory.id);

    await bp.addPlan(newPlanName);
    await expect(bp.storyName.filter({ hasText: newPlanName })).toBeVisible();
  });

  test("should allow editing stories", async ({ registeredPage }) => {
    const newType = "Story";
    const bp = new PokerGamePage(registeredPage.page);
    await bp.goto(pokerEditStory.id);

    await expect(
      bp.storyType.filter({ hasText: thanosPlan.type }),
    ).toBeVisible();
    await bp.editStoryBtn.click();
    await bp.storyTypeField.selectOption(newType);
    await bp.saveStoryBtn.click();
    await expect(bp.storyType.filter({ hasText: newType })).toBeVisible();
  });

  test("should allow deleting stories", async ({ registeredPage }) => {
    const bp = new PokerGamePage(registeredPage.page);
    await bp.goto(pokerDeleteStory.id);

    await expect(bp.storyName).toBeVisible();
    await bp.deleteStoryBtn.click();
    await expect(bp.storyName).not.toBeVisible();
  });

  test("should allow activating stories", async ({ registeredPage }) => {
    const bp = new PokerGamePage(registeredPage.page);
    await bp.goto(pokerActivateStory.id);

    await expect(
      bp.page.locator('[data-testid="currentplan-name"]'),
    ).toContainText("[Voting not started]");
    await expect(
      bp.page.locator('[data-testid="pointCard"][data-locked="true"]'),
    ).toHaveCount(8);

    await bp.page.locator('[data-testid="plan-activate"]').click();

    await expect(
      bp.page.locator('[data-testid="currentplan-name"]'),
    ).toContainText(scarletPlan.name);
    await expect(
      bp.page.locator('[data-testid="pointCard"][data-locked="false"]'),
    ).toHaveCount(8);
  });

  test("should allow skipping story voting", async ({ registeredPage }) => {
    const bp = new PokerGamePage(registeredPage.page);
    await bp.goto(pokerSkipStory.id);

    await expect(
      bp.page.locator('[data-testid="currentplan-name"]'),
    ).toContainText("[Voting not started]");
    await expect(
      bp.page.locator('[data-testid="pointCard"][data-locked="true"]'),
    ).toHaveCount(8);

    await bp.page.locator('[data-testid="plan-activate"]').click();

    await expect(
      bp.page.locator('[data-testid="currentplan-name"]'),
    ).toContainText(thanosPlan.name);
    await expect(
      bp.page.locator('[data-testid="pointCard"][data-locked="false"]'),
    ).toHaveCount(8);

    await bp.page.locator('[data-testid="voting-skip"]').click();

    await expect(
      bp.page.locator('[data-testid="currentplan-name"]'),
    ).toContainText("[Voting not started]");
    await expect(
      bp.page.locator('[data-testid="pointCard"][data-locked="true"]'),
    ).toHaveCount(8);
  });

  test("should allow finishing story voting", async ({ registeredPage }) => {
    const bp = new PokerGamePage(registeredPage.page);
    await bp.goto(pokerFinishVoting.id);

    await expect(
      bp.page.locator('[data-testid="currentplan-name"]'),
    ).toContainText("[Voting not started]");
    await expect(
      bp.page.locator('[data-testid="pointCard"][data-locked="true"]'),
    ).toHaveCount(8);

    await bp.page.locator('[data-testid="plan-activate"]').click();

    await expect(
      bp.page.locator('[data-testid="currentplan-name"]'),
    ).toContainText(lokiPlan.name);
    await expect(
      bp.page.locator('[data-testid="pointCard"][data-locked="false"]'),
    ).toHaveCount(8);

    await expect(
      bp.page.locator('[data-testid="voteresult-total"]'),
    ).not.toBeVisible();
    await expect(
      bp.page.locator('[data-testid="voteresult-average"]'),
    ).not.toBeVisible();
    await expect(
      bp.page.locator('[data-testid="voteresult-consensus"]'),
    ).not.toBeVisible();
    await expect(
      bp.page.locator('[data-testid="voteresult-agreement"]'),
    ).not.toBeVisible();

    await bp.page.locator('[data-testid="voting-finish"]').click();

    await expect(
      bp.page.locator('[data-testid="voteresult-total"]'),
    ).toBeVisible();
    await expect(
      bp.page.locator('[data-testid="voteresult-average"]'),
    ).toBeVisible();
    await expect(
      bp.page.locator('[data-testid="voteresult-consensus"]'),
    ).toBeVisible();
    await expect(
      bp.page.locator('[data-testid="voteresult-agreement"]'),
    ).toBeVisible();
    await expect(
      bp.page.locator('[data-testid="currentplan-name"]'),
    ).toContainText(lokiPlan.name);
    await expect(
      bp.page.locator('[data-testid="pointCard"]'),
    ).not.toBeVisible();
  });

  test("should allow saving story voting final points", async ({
    registeredPage,
  }) => {
    const bp = new PokerGamePage(registeredPage.page);
    await bp.goto(pokerSaveVoting.id);

    await expect(bp.page.locator('[data-testid="plans-unpointed"]')).toHaveText(
      "Unpointed (1)",
    );
    await expect(bp.page.locator('[data-testid="plans-pointed"]')).toHaveText(
      "Pointed (0)",
    );
    await expect(
      bp.page.locator('[data-testid="plan-points"]'),
    ).not.toBeVisible();

    await expect(
      bp.page.locator('[data-testid="currentplan-name"]'),
    ).toContainText("[Voting not started]");
    await bp.page.locator('[data-testid="plan-activate"]').click();
    await expect(
      bp.page.locator('[data-testid="currentplan-name"]'),
    ).toContainText(lokiPlan.name);
    await bp.page.locator('[data-testid="voting-finish"]').click();

    await expect(
      bp.page.locator('[data-testid="voteresult-total"]'),
    ).toBeVisible();

    await bp.page.locator('select[name="planPoints"]').selectOption("1");
    await bp.page.locator('[data-testid="voting-save"]').click();

    await expect(
      bp.page.locator('[data-testid="currentplan-name"]'),
    ).toContainText("[Voting not started]");
    await expect(
      bp.page.locator('[data-testid="plan-name"]'),
    ).not.toBeVisible();

    await expect(bp.page.locator('[data-testid="plans-unpointed"]')).toHaveText(
      "Unpointed (0)",
    );
    await expect(bp.page.locator('[data-testid="plans-pointed"]')).toHaveText(
      "Pointed (1)",
    );
    await bp.page.locator('[data-testid="plans-pointed"]').click();
    await expect(bp.page.locator('[data-testid="plan-name"]')).toHaveText(
      lokiPlan.name,
    );
    await expect(bp.page.locator('[data-testid="plan-points"]')).toHaveText(
      "1",
    );
  });

  test("delete game confirmation cancel does not delete game", async ({
    registeredPage,
  }) => {
    const bp = new PokerGamePage(registeredPage.page);
    await bp.goto(pokerCancelDelete.id);

    await bp.gameDeleteBtn.click();
    await bp.gameDeleteCancelBtn.click();

    await expect(bp.pageTitle).toHaveText(poker.name);
  });

  test("delete game confirmation confirm deletes game and redirects to games page", async ({
    registeredPage,
  }) => {
    const bp = new PokerGamePage(registeredPage.page);
    await bp.goto(pokerDelete.id);

    await bp.gameDeleteBtn.click();
    await bp.gameDeleteConfirmBtn.click();

    await expect(bp.page.locator("h1")).toHaveText("My Games");
  });
});
