import { expect, test } from "@fixtures/test-setup";

test.describe("Storyboard API", { tag: ["@api", "@storyboard"] }, () => {
  test("GET /users/{userId}/storyboards returns empty array when no storyboards associated to user", async ({
    request,
    adminApiUser,
  }) => {
    const response = await adminApiUser.context.get(
      `users/${adminApiUser.user.id}/storyboards`,
    );
    expect(response.ok()).toBeTruthy();
    const storyboards = await response.json();
    expect(storyboards.data).toEqual([]);
  });

  test("POST /users/{userId}/storyboards creates storyboard", async ({
    request,
    registeredApiUser,
  }) => {
    const storyboardName = "Test API Create Storyboard";

    const response = await registeredApiUser.context.post(
      `users/${registeredApiUser.user.id}/storyboards`,
      {
        data: {
          storyboardName,
        },
      },
    );
    expect(response.ok()).toBeTruthy();
    const storyboard = await response.json();
    expect(storyboard.data).toMatchObject({
      name: storyboardName,
    });
  });

  test("GET /users/{userId}/storyboards returns object in array when storyboards associated to user", async ({
    request,
    registeredApiUser,
  }) => {
    const storyboardName = "Test API Storyboards";

    await registeredApiUser.context.post(
      `users/${registeredApiUser.user.id}/storyboards`,
      {
        data: {
          storyboardName,
        },
      },
    );

    const response = await registeredApiUser.context.get(
      `users/${registeredApiUser.user.id}/storyboards`,
    );
    expect(response.ok()).toBeTruthy();
    const storyboards = await response.json();
    expect(storyboards.data).toContainEqual(
      expect.objectContaining({
        name: storyboardName,
      }),
    );
  });

  test(
    "POST /teams/{teamId}/users/{userId}/storyboards creates storyboard",
    { tag: ["@team"] },
    async ({ request, registeredApiUser }) => {
      const storyboardName = "Test API Create Team Storyboard";

      const teamResponse = await registeredApiUser.context.post(
        `users/${registeredApiUser.user.id}/teams`,
        {
          data: {
            name: "test team create storyboard",
          },
        },
      );
      const { data: team } = await teamResponse.json();

      const storyboardResponse = await registeredApiUser.context.post(
        `teams/${team.id}/users/${registeredApiUser.user.id}/storyboards`,
        {
          data: {
            storyboardName,
          },
        },
      );
      expect(storyboardResponse.ok()).toBeTruthy();
      const storyboard = await storyboardResponse.json();
      expect(storyboard.data).toMatchObject({
        name: storyboardName,
      });

      const storyboardsResponse = await registeredApiUser.context.get(
        `teams/${team.id}/storyboards`,
      );
      expect(storyboardsResponse.ok()).toBeTruthy();
      const storyboards = await storyboardsResponse.json();
      expect(storyboards.data).toContainEqual(
        expect.objectContaining({
          name: storyboardName,
        }),
      );
    },
  );

  test("POST /storyboards/{storyboardId}/goals creates storyboard goal", async ({
    request,
    registeredApiUser,
  }) => {
    const storyboardName = "Test API Create Storyboard Goal Test";
    const goalName = "Test API Create Goal";

    const response = await registeredApiUser.context.post(
      `users/${registeredApiUser.user.id}/storyboards`,
      {
        data: {
          storyboardName,
        },
      },
    );
    expect(response.ok()).toBeTruthy();
    const { data: storyboard } = await response.json();
    expect(storyboard).toMatchObject({
      name: storyboardName,
    });

    const goalResp = await registeredApiUser.context.post(
      `storyboards/${storyboard.id}/goals`,
      {
        data: {
          name: goalName,
        },
      },
    );
    expect(goalResp.ok()).toBeTruthy();

    const updatedStoryboard = await registeredApiUser.context.get(
      `storyboards/${storyboard.id}`,
    );
    expect(updatedStoryboard.ok()).toBeTruthy();
    const storyboardWithGoal = await updatedStoryboard.json();
    expect(storyboardWithGoal.data).toMatchObject(
      expect.objectContaining({
        goals: expect.arrayContaining([
          expect.objectContaining({
            name: goalName,
          }),
        ]),
      }),
    );
  });

  test("POST /storyboards/{storyboardId}/columns creates storyboard column", async ({
    request,
    registeredApiUser,
  }) => {
    const storyboardName = "Test API Create Storyboard Column Test";
    const goalName = "Test API Create Column Goal";

    const response = await registeredApiUser.context.post(
      `users/${registeredApiUser.user.id}/storyboards`,
      { data: { storyboardName } },
    );
    expect(response.ok()).toBeTruthy();
    const { data: storyboard } = await response.json();
    expect(storyboard).toMatchObject({
      name: storyboardName,
    });

    const goalResp = await registeredApiUser.context.post(
      `storyboards/${storyboard.id}/goals`,
      { data: { name: goalName } },
    );
    expect(goalResp.ok()).toBeTruthy();

    let updatedStoryboard = await registeredApiUser.context.get(
      `storyboards/${storyboard.id}`,
    );
    expect(updatedStoryboard.ok()).toBeTruthy();
    const storyboardWithGoal = await updatedStoryboard.json();

    const columnResp = await registeredApiUser.context.post(
      `storyboards/${storyboard.id}/columns`,
      {
        data: {
          goalId: storyboardWithGoal.data?.goals[0].id,
        },
      },
    );
    expect(columnResp.ok()).toBeTruthy();

    updatedStoryboard = await registeredApiUser.context.get(
      `storyboards/${storyboard.id}`,
    );
    expect(updatedStoryboard.ok()).toBeTruthy();
    const storyboardWithCol = await updatedStoryboard.json();
    expect(storyboardWithCol.data).toMatchObject(
      expect.objectContaining({
        goals: expect.arrayContaining([
          expect.objectContaining({
            name: goalName,
            columns: expect.arrayContaining([
              expect.objectContaining({
                name: "",
                id: expect.any(String),
              }),
            ]),
          }),
        ]),
      }),
    );
  });
});
