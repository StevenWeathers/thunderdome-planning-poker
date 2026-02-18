<script lang="ts">
  import { onMount } from 'svelte';

  import { ChevronRight } from 'lucide-svelte';
  import { user } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import { appRoutes } from '../../config';
  import { validateUserIsRegistered } from '../../validationUtils';
  import UsersList from '../../components/team/UsersList.svelte';
  import InvitesList from '../../components/team/InvitesList.svelte';

  import type { NotificationService } from '../../types/notifications';
  import type { ApiClient } from '../../types/apiclient';
  import TeamPageLayout from '../../components/team/TeamPageLayout.svelte';

  interface Props {
    xfetch: ApiClient;
    router: any;
    notifications: NotificationService;
    organizationId: any;
    departmentId: any;
    teamId: any;
  }

  let { xfetch, router, notifications, organizationId, departmentId, teamId }: Props = $props();

  const usersPageLimit = 1000;

  let invitesList = $state();

  let team = $state({
    id: '',
    name: '',
    subscribed: false,
  });
  let organization = $state({
    id: '',
    name: '',
    subscribed: false,
  });
  let department = $state({
    id: '',
    name: '',
  });

  let users = $state([]);
  $effect(() => {
    team.id = teamId;
    organization.id = organizationId;
    department.id = departmentId;
  });

  let organizationRole = $state('');
  let departmentRole = $state('');
  let teamRole = $state('');
  let isAdmin = $state(false);
  let isTeamMember = $state(false);

  const apiPrefix = '/api';
  let orgPrefix = $derived(
    departmentId
      ? `${apiPrefix}/organizations/${organizationId}/departments/${departmentId}`
      : `${apiPrefix}/organizations/${organizationId}`,
  );
  let teamPrefix = $derived(organizationId ? `${orgPrefix}/teams/${teamId}` : `${apiPrefix}/teams/${teamId}`);

  function getTeam() {
    xfetch(teamPrefix)
      .then(res => res.json())
      .then(function (result) {
        team = result.data.team;
        teamRole = result.data.teamRole;

        if (departmentId) {
          department = result.data.department;
          departmentRole = result.data.departmentRole;
        }
        if (organizationId) {
          organization = result.data.organization;
          organizationRole = result.data.organizationRole;
        }

        isAdmin = organizationRole === 'ADMIN' || departmentRole === 'ADMIN' || teamRole === 'ADMIN';
        isTeamMember = organizationRole === 'ADMIN' || departmentRole === 'ADMIN' || teamRole !== '';

        getUsers();
      })
      .catch(function () {
        notifications.danger($LL.teamGetError());
      });
  }

  function getUsers() {
    const usersOffset = (usersPage - 1) * usersPageLimit;
    xfetch(`${teamPrefix}/users?limit=${usersPageLimit}&offset=${usersOffset}`)
      .then(res => res.json())
      .then(function (result) {
        users = result.data;
      })
      .catch(function () {
        notifications.danger($LL.teamGetUsersError());
      });
  }

  onMount(() => {
    if (!$user.id || !validateUserIsRegistered($user)) {
      router.route(appRoutes.login);
      return;
    }

    getTeam();
  });
</script>

<svelte:head>
  <title>{$LL.users()} {team.name} | {$LL.appName()}</title>
</svelte:head>

<TeamPageLayout activePage="users" {teamId} {organizationId} {departmentId}>
  <div class="flex mb-6 lg:mb-8">
    <div class="flex-1">
      <h1 class="text-3xl font-semibold font-rajdhani dark:text-white">
        <span class="uppercase">{$LL.team()}</span>
        <ChevronRight class="w-8 h-8 inline-block" />
        {team.name}
        <ChevronRight class="w-8 h-8 inline-block" />
        {$LL.users()}
      </h1>

      {#if organizationId}
        <div class="text-xl font-semibold font-rajdhani dark:text-white">
          <span class="uppercase">{$LL.organization()}</span>
          <ChevronRight class="inline-block" />
          <a
            class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
            href="{appRoutes.organization}/{organization.id}"
          >
            {organization.name}
          </a>
          {#if departmentId}
            &nbsp;
            <ChevronRight class="inline-block" />
            <span class="uppercase">{$LL.department()}</span>
            <ChevronRight class="inline-block" />
            <a
              class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
              href="{appRoutes.organization}/{organization.id}/department/{department.id}"
            >
              {department.name}
            </a>
          {/if}
        </div>
      {/if}
    </div>
  </div>

  {#if isAdmin}
    <div class="w-full mb-6 lg:mb-8">
      <InvitesList {xfetch} {notifications} pageType="team" {teamPrefix} bind:this={invitesList} />
    </div>
  {/if}

  <UsersList
    {users}
    {getUsers}
    {xfetch}
    {notifications}
    {isAdmin}
    pageType="team"
    {teamPrefix}
    orgId={organizationId}
    deptId={departmentId}
    on:user-invited={() => {
      invitesList.f('user-invited');
    }}
  />
</TeamPageLayout>
