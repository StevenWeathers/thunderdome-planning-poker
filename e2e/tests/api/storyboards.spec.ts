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

  test("POST /teams/{teamId}/users/{userId}/storyboards creates storyboard", async ({
    request,
    registeredApiUser,
  }) => {
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
  });
});
