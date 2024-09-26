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

  test.describe.serial("board actions", () => {
    const storyboardName = "Test API Storyboard actions";
    let storyboard;

    test.beforeAll(async ({ registeredApiUser }) => {
      const response = await registeredApiUser.context.post(
        `users/${registeredApiUser.user.id}/storyboards`,
        {
          data: {
            storyboardName,
          },
        },
      );
      expect(response.ok()).toBeTruthy();
      const res = await response.json();
      storyboard = res.data;
    });

    test("POST /storyboards/{storyboardId}/goals creates storyboard goal", async ({
      request,
      registeredApiUser,
    }) => {
      const goalName = "Test API Create Goal";

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
      storyboard = storyboardWithGoal.data;
    });

    test("POST /storyboards/{storyboardId}/columns creates storyboard column", async ({
      request,
      registeredApiUser,
    }) => {
      const columnResp = await registeredApiUser.context.post(
        `storyboards/${storyboard.id}/columns`,
        {
          data: {
            goalId: storyboard.goals[0].id,
          },
        },
      );
      expect(columnResp.ok()).toBeTruthy();

      const updatedStoryboard = await registeredApiUser.context.get(
        `storyboards/${storyboard.id}`,
      );
      expect(updatedStoryboard.ok()).toBeTruthy();
      const storyboardWithCol = await updatedStoryboard.json();
      expect(storyboardWithCol.data).toMatchObject(
        expect.objectContaining({
          goals: expect.arrayContaining([
            expect.objectContaining({
              name: storyboard.goals[0].name,
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

      storyboard = storyboardWithCol.data;
    });

    test("POST /storyboards/{storyboardId}/stories creates storyboard story", async ({
      request,
      registeredApiUser,
    }) => {
      const resp = await registeredApiUser.context.post(
        `storyboards/${storyboard.id}/stories`,
        {
          data: {
            goalId: storyboard.goals[0].id,
            columnId: storyboard.goals[0].columns[0].id,
          },
        },
      );
      expect(resp.ok()).toBeTruthy();

      const updatedStoryboard = await registeredApiUser.context.get(
        `storyboards/${storyboard.id}`,
      );
      expect(updatedStoryboard.ok()).toBeTruthy();
      const storyboardWithStory = await updatedStoryboard.json();
      expect(storyboardWithStory.data).toMatchObject(
        expect.objectContaining({
          goals: expect.arrayContaining([
            expect.objectContaining({
              name: storyboard.goals[0].name,
              columns: expect.arrayContaining([
                expect.objectContaining({
                  name: "",
                  id: storyboard.goals[0].columns[0].id,
                  stories: expect.arrayContaining([
                    expect.objectContaining({
                      id: expect.any(String),
                    }),
                  ]),
                }),
              ]),
            }),
          ]),
        }),
      );
    });
  });
});
