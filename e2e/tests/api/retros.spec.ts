import { expect, test } from "@fixtures/test-setup";

test.describe("Retro API", { tag: ["@api", "@retro"] }, () => {
  test("GET /users/{userId}/retros returns empty array when no retros associated to user", async ({
    request,
    adminApiUser,
  }) => {
    const response = await adminApiUser.context.get(
      `users/${adminApiUser.user.id}/retros`,
    );
    expect(response.ok()).toBeTruthy();
    const retros = await response.json();
    expect(retros.data).toEqual([]);
  });

  test("POST /users/{userId}/retros creates retro", async ({
    request,
    registeredApiUser,
  }) => {
    const retroName = "Test API Create Retro";
    const brainstormVisibility = "visible";
    const maxVotes = 3;

    const response = await registeredApiUser.context.post(
      `users/${registeredApiUser.user.id}/retros`,
      {
        data: {
          retroName,
          brainstormVisibility,
          maxVotes,
        },
      },
    );
    expect(response.ok()).toBeTruthy();
    const retro = await response.json();
    expect(retro.data).toMatchObject({
      name: retroName,
      brainstormVisibility,
    });
  });

  test("POST /users/{userId}/retros creates retro with phase time limit and auto-advance settings", async ({
    request,
    registeredApiUser,
  }) => {
    const retroName = "Test API Create Retro with Timer";
    const brainstormVisibility = "visible";
    const maxVotes = 3;
    const phaseTimeLimitMin = 15;
    const phaseAutoAdvance = false;

    const response = await registeredApiUser.context.post(
      `users/${registeredApiUser.user.id}/retros`,
      {
        data: {
          retroName,
          brainstormVisibility,
          maxVotes,
          phaseTimeLimitMin,
          phaseAutoAdvance,
        },
      },
    );
    expect(response.ok()).toBeTruthy();
    const retro = await response.json();
    expect(retro.data).toMatchObject({
      name: retroName,
      brainstormVisibility,
      phase_time_limit_min: phaseTimeLimitMin,
      phase_auto_advance: phaseAutoAdvance,
    });
  });

  test("GET /users/{userId}/retros returns object in array when retros associated to user", async ({
    request,
    registeredApiUser,
  }) => {
    const retroName = "Test API Retros";
    const brainstormVisibility = "hidden";
    const maxVotes = 3;

    await registeredApiUser.context.post(
      `users/${registeredApiUser.user.id}/retros`,
      {
        data: {
          retroName,
          brainstormVisibility,
          maxVotes,
        },
      },
    );

    const response = await registeredApiUser.context.get(
      `users/${registeredApiUser.user.id}/retros`,
    );
    expect(response.ok()).toBeTruthy();
    const retros = await response.json();
    expect(retros.data).toContainEqual(
      expect.objectContaining({
        name: retroName,
      }),
    );
  });

  test("POST /teams/{teamId}/users/{userId}/retros creates retro", async ({
    request,
    registeredApiUser,
  }) => {
    const retroName = "Test API Create Team Retro";
    const brainstormVisibility = "hidden";
    const maxVotes = 3;

    const teamResponse = await registeredApiUser.context.post(
      `users/${registeredApiUser.user.id}/teams`,
      {
        data: {
          name: "test team create retro",
        },
      },
    );
    const { data: team } = await teamResponse.json();

    const retroResponse = await registeredApiUser.context.post(
      `teams/${team.id}/users/${registeredApiUser.user.id}/retros`,
      {
        data: {
          retroName,
          brainstormVisibility,
          maxVotes,
        },
      },
    );
    expect(retroResponse.ok()).toBeTruthy();
    const retro = await retroResponse.json();
    expect(retro.data).toMatchObject({
      name: retroName,
    });

    const retrosResponse = await registeredApiUser.context.get(
      `teams/${team.id}/retros`,
    );
    expect(retrosResponse.ok()).toBeTruthy();
    const retros = await retrosResponse.json();
    expect(retros.data).toContainEqual(
      expect.objectContaining({
        name: retroName,
      }),
    );
  });
});
