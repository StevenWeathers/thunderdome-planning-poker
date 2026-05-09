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

  test("GET /users/{userId}/retro-actions returns assigned open actions and supports team filtering", async ({
    registeredApiUser,
    testDatabase,
  }) => {
    const firstTeamResponse = await registeredApiUser.context.post(
      `users/${registeredApiUser.user.id}/teams`,
      {
        data: {
          name: "retro actions alpha",
        },
      },
    );
    const { data: firstTeam } = await firstTeamResponse.json();

    const secondTeamResponse = await registeredApiUser.context.post(
      `users/${registeredApiUser.user.id}/teams`,
      {
        data: {
          name: "retro actions beta",
        },
      },
    );
    const { data: secondTeam } = await secondTeamResponse.json();

    const firstRetroResponse = await registeredApiUser.context.post(
      `teams/${firstTeam.id}/users/${registeredApiUser.user.id}/retros`,
      {
        data: {
          retroName: "Alpha retro",
          brainstormVisibility: "visible",
          maxVotes: 3,
        },
      },
    );
    const { data: firstRetro } = await firstRetroResponse.json();

    const secondRetroResponse = await registeredApiUser.context.post(
      `teams/${secondTeam.id}/users/${registeredApiUser.user.id}/retros`,
      {
        data: {
          retroName: "Beta retro",
          brainstormVisibility: "visible",
          maxVotes: 3,
        },
      },
    );
    const { data: secondRetro } = await secondRetroResponse.json();

    const firstActionResult = await testDatabase.pool.query(
      `INSERT INTO thunderdome.retro_action (retro_id, content, completed)
       VALUES ($1, $2, false)
       RETURNING id`,
      [firstRetro.id, "Ship dashboard action summary"],
    );
    const secondActionResult = await testDatabase.pool.query(
      `INSERT INTO thunderdome.retro_action (retro_id, content, completed)
       VALUES ($1, $2, false)
       RETURNING id`,
      [secondRetro.id, "Clean up beta follow-up"],
    );
    const completedActionResult = await testDatabase.pool.query(
      `INSERT INTO thunderdome.retro_action (retro_id, content, completed)
       VALUES ($1, $2, true)
       RETURNING id`,
      [firstRetro.id, "Completed item should be hidden"],
    );

    await testDatabase.pool.query(
      `INSERT INTO thunderdome.retro_action_assignee (action_id, user_id) VALUES ($1, $2), ($3, $2), ($4, $2)`,
      [
        firstActionResult.rows[0].id,
        registeredApiUser.user.id,
        secondActionResult.rows[0].id,
        completedActionResult.rows[0].id,
      ],
    );

    const response = await registeredApiUser.context.get(
      `users/${registeredApiUser.user.id}/retro-actions`,
    );
    expect(response.ok()).toBeTruthy();
    const retroActions = await response.json();
    expect(retroActions.meta.count).toBe(2);
    expect(retroActions.data).toEqual(
      expect.arrayContaining([
        expect.objectContaining({
          content: "Ship dashboard action summary",
          teamId: firstTeam.id,
          teamName: firstTeam.name,
        }),
        expect.objectContaining({
          content: "Clean up beta follow-up",
          teamId: secondTeam.id,
          teamName: secondTeam.name,
        }),
      ]),
    );

    const filteredResponse = await registeredApiUser.context.get(
      `users/${registeredApiUser.user.id}/retro-actions?teamId=${firstTeam.id}`,
    );
    expect(filteredResponse.ok()).toBeTruthy();
    const filteredRetroActions = await filteredResponse.json();
    expect(filteredRetroActions.meta.count).toBe(1);
    expect(filteredRetroActions.data).toEqual([
      expect.objectContaining({
        content: "Ship dashboard action summary",
        teamId: firstTeam.id,
        teamName: firstTeam.name,
      }),
    ]);
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
});
