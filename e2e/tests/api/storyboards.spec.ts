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

      storyboard = storyboardWithStory.data;
    });

    test("PUT /storyboards/{storyboardId}/stories/{storyId}/move moves storyboard story", async ({
      request,
      registeredApiUser,
    }) => {
      await test.step("Create 2 more stories", async () => {
        const story2resp = await registeredApiUser.context.post(
          `storyboards/${storyboard.id}/stories`,
          {
            data: {
              goalId: storyboard.goals[0].id,
              columnId: storyboard.goals[0].columns[0].id,
            },
          },
        );
        expect(story2resp.ok()).toBeTruthy();

        const story3resp = await registeredApiUser.context.post(
          `storyboards/${storyboard.id}/stories`,
          {
            data: {
              goalId: storyboard.goals[0].id,
              columnId: storyboard.goals[0].columns[0].id,
            },
          },
        );
        expect(story3resp.ok()).toBeTruthy();

        const updatedStoryboard = await registeredApiUser.context.get(
          `storyboards/${storyboard.id}`,
        );
        expect(updatedStoryboard.ok()).toBeTruthy();
        const storyboardWithStory = await updatedStoryboard.json();
        storyboard = storyboardWithStory.data;
      });

      await test.step("Move story to begining of column", async () => {
        const storyToMoveID = storyboard.goals[0].columns[0].stories[2].id;
        const storyPlaceBeforeID = storyboard.goals[0].columns[0].stories[0].id;
        const storyMresp = await registeredApiUser.context.put(
          `storyboards/${storyboard.id}/stories/${storyToMoveID}/move`,
          {
            data: {
              goalId: storyboard.goals[0].id,
              columnId: storyboard.goals[0].columns[0].id,
              placeBefore: storyPlaceBeforeID,
            },
          },
        );
        expect(storyMresp.ok()).toBeTruthy();

        const updatedStoryboard2 = await registeredApiUser.context.get(
          `storyboards/${storyboard.id}`,
        );
        expect(updatedStoryboard2.ok()).toBeTruthy();
        const storyboardWithMovedStory = await updatedStoryboard2.json();
        expect(
          storyboardWithMovedStory.data.goals[0].columns[0].stories[0].id,
        ).toEqual(storyToMoveID);
        expect(
          storyboardWithMovedStory.data.goals[0].columns[0].stories[1].id,
        ).toEqual(storyPlaceBeforeID);
        storyboard = storyboardWithMovedStory.data;
      });

      await test.step("Move story to between 2 stories", async () => {
        const storyToMoveID = storyboard.goals[0].columns[0].stories[2].id;
        const storyPlaceBeforeID = storyboard.goals[0].columns[0].stories[1].id;
        const storyMresp = await registeredApiUser.context.put(
          `storyboards/${storyboard.id}/stories/${storyToMoveID}/move`,
          {
            data: {
              goalId: storyboard.goals[0].id,
              columnId: storyboard.goals[0].columns[0].id,
              placeBefore: storyPlaceBeforeID,
            },
          },
        );
        expect(storyMresp.ok()).toBeTruthy();

        const updatedStoryboard2 = await registeredApiUser.context.get(
          `storyboards/${storyboard.id}`,
        );
        expect(updatedStoryboard2.ok()).toBeTruthy();
        const storyboardWithMovedStory = await updatedStoryboard2.json();
        expect(
          storyboardWithMovedStory.data.goals[0].columns[0].stories[1].id,
        ).toEqual(storyToMoveID);
        expect(
          storyboardWithMovedStory.data.goals[0].columns[0].stories[2].id,
        ).toEqual(storyPlaceBeforeID);
        storyboard = storyboardWithMovedStory.data;
      });

      await test.step("Move story to end of column", async () => {
        const storyToMoveID = storyboard.goals[0].columns[0].stories[0].id;
        const storyPlaceBeforeID = "";
        const storyMresp = await registeredApiUser.context.put(
          `storyboards/${storyboard.id}/stories/${storyToMoveID}/move`,
          {
            data: {
              goalId: storyboard.goals[0].id,
              columnId: storyboard.goals[0].columns[0].id,
              placeBefore: storyPlaceBeforeID,
            },
          },
        );
        expect(storyMresp.ok()).toBeTruthy();

        const updatedStoryboard2 = await registeredApiUser.context.get(
          `storyboards/${storyboard.id}`,
        );
        expect(updatedStoryboard2.ok()).toBeTruthy();
        const storyboardWithMovedStory = await updatedStoryboard2.json();
        expect(
          storyboardWithMovedStory.data.goals[0].columns[0].stories[2].id,
        ).toEqual(storyToMoveID);
        storyboard = storyboardWithMovedStory.data;
      });

      await test.step("Create another column", async () => {
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
        const storyboardWithAnotherColumn = await updatedStoryboard.json();
        storyboard = storyboardWithAnotherColumn.data;
      });

      await test.step("Move story to another column", async () => {
        const storyToMoveID = storyboard.goals[0].columns[0].stories[0].id;
        const storyPlaceBeforeID = "";
        const storyMresp = await registeredApiUser.context.put(
          `storyboards/${storyboard.id}/stories/${storyToMoveID}/move`,
          {
            data: {
              goalId: storyboard.goals[0].id,
              columnId: storyboard.goals[0].columns[1].id,
              placeBefore: storyPlaceBeforeID,
            },
          },
        );
        expect(storyMresp.ok()).toBeTruthy();

        const updatedStoryboard2 = await registeredApiUser.context.get(
          `storyboards/${storyboard.id}`,
        );
        expect(updatedStoryboard2.ok()).toBeTruthy();
        const storyboardWithMovedStory = await updatedStoryboard2.json();
        expect(
          storyboardWithMovedStory.data.goals[0].columns[1].stories[0].id,
        ).toEqual(storyToMoveID);
        storyboard = storyboardWithMovedStory.data;
      });

      await test.step("Create another goal", async () => {
        const columnResp = await registeredApiUser.context.post(
          `storyboards/${storyboard.id}/goals`,
          {
            data: {
              name: "second goal to move story to",
            },
          },
        );
        expect(columnResp.ok()).toBeTruthy();
        const updatedStoryboard = await registeredApiUser.context.get(
          `storyboards/${storyboard.id}`,
        );
        expect(updatedStoryboard.ok()).toBeTruthy();
        const storyboardWithAnotherGoal = await updatedStoryboard.json();
        storyboard = storyboardWithAnotherGoal.data;
      });

      await test.step("Create column in second goal", async () => {
        const columnResp = await registeredApiUser.context.post(
          `storyboards/${storyboard.id}/columns`,
          {
            data: {
              goalId: storyboard.goals[1].id,
            },
          },
        );
        expect(columnResp.ok()).toBeTruthy();
        const updatedStoryboard = await registeredApiUser.context.get(
          `storyboards/${storyboard.id}`,
        );
        expect(updatedStoryboard.ok()).toBeTruthy();
        const storyboardWithAnotherColumn = await updatedStoryboard.json();
        storyboard = storyboardWithAnotherColumn.data;
      });

      await test.step("Move story to another goals column", async () => {
        const storyToMoveID = storyboard.goals[0].columns[0].stories[0].id;
        const storyPlaceBeforeID = "";
        const storyMresp = await registeredApiUser.context.put(
          `storyboards/${storyboard.id}/stories/${storyToMoveID}/move`,
          {
            data: {
              goalId: storyboard.goals[1].id,
              columnId: storyboard.goals[1].columns[0].id,
              placeBefore: storyPlaceBeforeID,
            },
          },
        );
        expect(storyMresp.ok()).toBeTruthy();

        const updatedStoryboard2 = await registeredApiUser.context.get(
          `storyboards/${storyboard.id}`,
        );
        expect(updatedStoryboard2.ok()).toBeTruthy();
        const storyboardWithMovedStory = await updatedStoryboard2.json();
        expect(
          storyboardWithMovedStory.data.goals[1].columns[0].stories[0].id,
        ).toEqual(storyToMoveID);
        storyboard = storyboardWithMovedStory.data;
      });
    });
  });
});
