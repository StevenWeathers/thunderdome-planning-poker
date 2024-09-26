import { expect, test } from "@fixtures/user-sessions";
import { RetroPage } from "@fixtures/pages/retro-page";

test.describe("Retro page", { tag: ["@retro"] }, () => {
  let retro = { id: "", name: "e2e retro page tests" };
  let retroLeave = { id: "" };
  let retroCancelDelete = { id: "" };
  let retroDelete = { id: "" };
  let retroAdvancePhases = { id: "" };
  let retroPhaseBrainstorm = { id: "" };
  let retroPhaseGroup = { id: "" };
  let retroPhaseVote = { id: "" };
  let retroPhaseActionItem = { id: "" };

  test.beforeAll(async ({ registeredPage, verifiedPage, adminPage }) => {
    const commonRetro = {
      retroName: `${retro.name}`,
      maxVotes: 3,
      brainstormVisibility: "visible",
      retroFacilitators: [`${adminPage.user.email}`],
    };
    retro = await registeredPage.createRetro({ ...commonRetro });
    retroLeave = await verifiedPage.createRetro({ ...commonRetro });
    retroCancelDelete = await registeredPage.createRetro({
      ...commonRetro,
    });
    retroDelete = await registeredPage.createRetro({ ...commonRetro });
    retroAdvancePhases = await registeredPage.createRetro({
      ...commonRetro,
    });
    retroPhaseBrainstorm = await registeredPage.createRetro({
      ...commonRetro,
    });
    retroPhaseGroup = await registeredPage.createRetro({
      ...commonRetro,
    });
    retroPhaseVote = await registeredPage.createRetro({
      ...commonRetro,
    });
    retroPhaseActionItem = await registeredPage.createRetro({
      ...commonRetro,
    });
  });

  test(
    "unauthenticated user redirects to register",
    { tag: ["@unauthenticated"] },
    async ({ page }) => {
      const bp = new RetroPage(page);
      await bp.goto(retro.id);

      const title = bp.page.locator("h1");
      await expect(title).toHaveText("Register");
    },
  );

  test(
    "guest user successfully loads",
    { tag: ["@guest"] },
    async ({ guestPage }) => {
      const bp = new RetroPage(guestPage.page);
      await bp.goto(retro.id);

      await expect(bp.retroTitle).toHaveText(retro.name);
    },
  );

  test(
    "registered user successfully loads",
    { tag: ["@registered"] },
    async ({ registeredPage }) => {
      const bp = new RetroPage(registeredPage.page);
      await bp.goto(retro.id);

      await expect(bp.retroTitle).toHaveText(retro.name);
    },
  );

  test("user can leave retro", async ({ registeredPage }) => {
    const bp = new RetroPage(registeredPage.page);
    await bp.goto(retroLeave.id);

    await bp.page.click('[data-testid="retro-leave"]');
    await expect(bp.page.locator("h1")).toHaveText("My Retros");
  });

  test("delete retro confirmation cancel does not delete retro", async ({
    registeredPage,
  }) => {
    const bp = new RetroPage(registeredPage.page);
    await bp.goto(retroCancelDelete.id);

    await bp.retroDeleteBtn.click();
    await bp.retroDeleteCancelBtn.click();

    await expect(bp.retroTitle).toHaveText(retro.name);
  });

  test("delete retro confirmation confirm deletes retro and redirects to retros page", async ({
    registeredPage,
  }) => {
    const bp = new RetroPage(registeredPage.page);
    await bp.goto(retroDelete.id);

    await bp.retroDeleteBtn.click();
    await bp.retroDeleteConfirmBtn.click();

    await expect(bp.page.locator("h1")).toHaveText("My Retros");
  });

  test("facilitator can advance phases", async ({ registeredPage }) => {
    const bp = new RetroPage(registeredPage.page);
    await bp.goto(retroAdvancePhases.id);

    await bp.retroNextPhaseBtn.click();
    await expect(bp.page.getByText("Add your comments below")).toBeVisible();

    await bp.retroNextPhaseBtn.click();
    await expect(
      bp.page.getByText("Drag and drop comments to group them together"),
    ).toBeVisible();

    await bp.retroNextPhaseBtn.click();
    await expect(
      bp.page.getByText("Vote for the groups you'd like to discuss most"),
    ).toBeVisible();

    await bp.retroNextPhaseBtn.click();
    await expect(
      bp.page.getByText(
        "Add action items, you can no longer group or vote comments",
      ),
    ).toBeVisible();

    await bp.retroNextPhaseBtn.click();
    await expect(bp.retroExportBtn).toBeVisible();
  });

  test("brainstorm phase can add items", async ({ registeredPage }) => {
    const bp = new RetroPage(registeredPage.page);
    await bp.goto(retroPhaseBrainstorm.id);
    const happyItem1 = "happy test item 1";
    const sadItem1 = "sad test item 1";
    const questionItem1 = "question test item 1";

    await bp.retroPhaseBrainstormBtn.click();

    await bp.retroWorkedWellInput.fill(happyItem1);
    await bp.retroWorkedWellInput.press("Enter");
    expect(await bp.page.getByText(happyItem1));

    await bp.retroNeedsImprovementInput.fill(sadItem1);
    await bp.retroNeedsImprovementInput.press("Enter");
    expect(await bp.page.getByText(sadItem1));

    await bp.retroQuestionInput.fill(questionItem1);
    await bp.retroQuestionInput.press("Enter");
    expect(await bp.page.getByText(questionItem1));
  });

  test("group phase can group items", async ({ registeredPage }) => {
    const bp = new RetroPage(registeredPage.page);
    await bp.goto(retroPhaseGroup.id);
    const firstGroupName = "Test Group #1";

    // @TODO replace this boilerplate with API calls in beforeAll
    const happyItem1 = "happy test item 1";
    const sadItem1 = "sad test item 1";
    const questionItem1 = "question test item 1";
    await bp.retroPhaseBrainstormBtn.click();
    await bp.retroWorkedWellInput.fill(happyItem1);
    await bp.retroWorkedWellInput.press("Enter");
    await bp.retroNeedsImprovementInput.fill(sadItem1);
    await bp.retroNeedsImprovementInput.press("Enter");
    await bp.retroQuestionInput.fill(questionItem1);
    await bp.retroQuestionInput.press("Enter");

    await bp.retroPhaseGroupBtn.click();

    const group1Input = await bp.retroGroupNameInput.first();
    await group1Input.focus();
    await group1Input.fill(firstGroupName);
    await bp.page.keyboard.press("Tab");
    expect(await bp.retroGroupNameInput.first().inputValue()).toEqual(
      firstGroupName,
    );
  });

  test("action phase can add items", async ({ registeredPage }) => {
    const bp = new RetroPage(registeredPage.page);
    await bp.goto(retroPhaseActionItem.id);
    const actionItem1 = "action test item 1";
    const actionItem2 = "action test item 2";

    await bp.retroPhaseActionItemsBtn.click();

    await bp.retroActionItemInput.fill(actionItem1);
    await bp.retroActionItemInput.press("Enter");
    expect(await bp.page.getByText(actionItem1));

    await bp.retroActionItemInput.fill(actionItem2);
    await bp.retroActionItemInput.press("Enter");
    expect(await bp.page.getByText(actionItem2));
  });
});
