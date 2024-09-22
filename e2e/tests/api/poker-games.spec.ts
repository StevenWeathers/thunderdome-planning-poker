import { expect, test } from "@fixtures/test-setup";

test.describe("Poker API", { tag: ["@api", "@poker"] }, () => {
  test("GET /users/{userId}/battles returns empty array when no games associated to user", async ({
    request,
    adminApiUser,
  }) => {
    const response = await adminApiUser.context.get(
      `users/${adminApiUser.user.id}/battles`,
    );
    expect(response.ok()).toBeTruthy();
    const battles = await response.json();
    expect(battles.data).toEqual([]);
  });

  test("POST /users/{userId}/battles creates game", async ({
    request,
    registeredApiUser,
  }) => {
    const pointValuesAllowed = ["0", "1/2", "1", "2", "3", "5", "8", "13"];
    const battleName = "Test API Create Game";
    const pointAverageRounding = "floor";
    const autoFinishVoting = false;

    const response = await registeredApiUser.context.post(
      `users/${registeredApiUser.user.id}/battles`,
      {
        data: {
          name: battleName,
          pointValuesAllowed,
          pointAverageRounding,
          autoFinishVoting,
        },
      },
    );
    expect(response.ok()).toBeTruthy();
    const battle = await response.json();
    expect(battle.data).toMatchObject({
      name: battleName,
      pointValuesAllowed,
      pointAverageRounding,
      autoFinishVoting,
    });
  });

  test("GET /users/{userId}/battles returns object in array when games associated to user", async ({
    request,
    registeredApiUser,
  }) => {
    const pointValuesAllowed = ["1", "2", "3", "5", "8", "13"];
    const battleName = "Test API Games";
    const pointAverageRounding = "ceil";
    const autoFinishVoting = true;

    await registeredApiUser.context.post(
      `users/${registeredApiUser.user.id}/battles`,
      {
        data: {
          name: battleName,
          pointValuesAllowed,
          pointAverageRounding,
          autoFinishVoting,
        },
      },
    );

    const response = await registeredApiUser.context.get(
      `users/${registeredApiUser.user.id}/battles`,
    );
    expect(response.ok()).toBeTruthy();
    const battles = await response.json();
    expect(battles.data).toContainEqual(
      expect.objectContaining({
        name: battleName,
        pointValuesAllowed,
        pointAverageRounding,
        autoFinishVoting,
      }),
    );
  });

  test("POST /teams/{teamId}/users/{userId}/battles creates game", async ({
    request,
    registeredApiUser,
  }) => {
    const pointValuesAllowed = ["0", "1/2", "1", "2", "3", "5", "8", "13"];
    const battleName = "Test API Create Game";
    const pointAverageRounding = "floor";
    const autoFinishVoting = false;

    const teamResponse = await registeredApiUser.context.post(
      `users/${registeredApiUser.user.id}/teams`,
      {
        data: {
          name: "test team create game",
        },
      },
    );
    const { data: team } = await teamResponse.json();

    const battleResponse = await registeredApiUser.context.post(
      `teams/${team.id}/users/${registeredApiUser.user.id}/battles`,
      {
        data: {
          name: battleName,
          pointValuesAllowed,
          pointAverageRounding,
          autoFinishVoting,
        },
      },
    );
    expect(battleResponse.ok()).toBeTruthy();
    const battle = await battleResponse.json();
    expect(battle.data).toMatchObject({
      name: battleName,
      pointValuesAllowed,
      pointAverageRounding,
      autoFinishVoting,
    });

    const battlesResponse = await registeredApiUser.context.get(
      `teams/${team.id}/battles`,
    );
    expect(battlesResponse.ok()).toBeTruthy();
    const battles = await battlesResponse.json();
    expect(battles.data).toContainEqual(
      expect.objectContaining({
        name: battleName,
      }),
    );
  });
});
