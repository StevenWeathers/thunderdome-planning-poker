<script lang="ts">
  import { onMount } from 'svelte';

  import PageLayout from '../../components/PageLayout.svelte';
  import CreateTeam from '../../components/team/CreateTeam.svelte';
  import DeleteConfirmation from '../../components/global/DeleteConfirmation.svelte';
  import { ChevronRight } from 'lucide-svelte';
  import { user } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import { appRoutes } from '../../config';
  import { validateUserIsRegistered } from '../../validationUtils';
  import RowCol from '../../components/table/RowCol.svelte';
  import TableRow from '../../components/table/TableRow.svelte';
  import HeadCol from '../../components/table/HeadCol.svelte';
  import Table from '../../components/table/Table.svelte';
  import UsersList from '../../components/team/UsersList.svelte';
  import TableContainer from '../../components/table/TableContainer.svelte';
  import TableNav from '../../components/table/TableNav.svelte';
  import CrudActions from '../../components/table/CrudActions.svelte';
  import InvitesList from '../../components/team/InvitesList.svelte';

  interface Props {
    xfetch: any;
    router: any;
    notifications: any;
    organizationId: any;
    departmentId: any;
  }

  let {
    xfetch,
    router,
    notifications,
    organizationId,
    departmentId
  }: Props = $props();

  const teamsPageLimit = 1000;
  const usersPageLimit = 1000;
  const deptPrefix = `/api/organizations/${organizationId}/departments/${departmentId}`;

  let invitesList = $state();
  let organization = $state({
    id: organizationId,
    name: '',
    createdDate: '',
    updateDate: '',
    subscribed: false,
  });
  let department = $state({
    id: departmentId,
    name: '',
  });
  let departmentRole = $state('');
  let organizationRole = $state('');
  let teams = $state([]);
  let users = $state([]);
  let invites = [];
  let showCreateTeam = $state(false);
  let showDeleteTeam = $state(false);
  let deleteTeamId = null;
  let teamsPage = 1;
  let usersPage = 1;

  function toggleCreateTeam() {
    showCreateTeam = !showCreateTeam;
  }

  const toggleDeleteTeam = teamId => () => {
    showDeleteTeam = !showDeleteTeam;
    deleteTeamId = teamId;
  };

  function getDepartment() {
    xfetch(`/api/organizations/${organizationId}/departments/${departmentId}`)
      .then(res => res.json())
      .then(function (result) {
        department = result.data.department;
        organization = result.data.organization;
        organizationRole = result.data.organizationRole;
        departmentRole = result.data.departmentRole;

        getTeams();
        getUsers();
      })
      .catch(function () {
        notifications.danger($LL.departmentGetError());
      });
  }

  function getUsers() {
    const usersOffset = (usersPage - 1) * usersPageLimit;
    xfetch(
      `/api/organizations/${organizationId}/departments/${departmentId}/users?limit=${usersPageLimit}&offset=${usersOffset}`,
    )
      .then(res => res.json())
      .then(function (result) {
        users = result.data;
      })
      .catch(function () {
        notifications.danger($LL.teamGetUsersError());
      });
  }

  function getTeams() {
    const teamsOffset = (teamsPage - 1) * teamsPageLimit;
    xfetch(
      `/api/organizations/${organizationId}/departments/${departmentId}/teams?limit=${teamsPageLimit}&offset=${teamsOffset}`,
    )
      .then(res => res.json())
      .then(function (result) {
        teams = result.data;
      })
      .catch(function () {
        notifications.danger($LL.departmentTeamsGetError());
      });
  }

  function createTeamHandler(name) {
    const body = {
      name,
    };

    xfetch(
      `/api/organizations/${organizationId}/departments/${departmentId}/teams`,
      { body },
    )
      .then(res => res.json())
      .then(function () {
        toggleCreateTeam();
        notifications.success($LL.teamCreateSuccess());
        getTeams();
      })
      .catch(function () {
        notifications.danger($LL.teamCreateError());
      });
  }

  function handleDeleteTeam() {
    xfetch(
      `/api/organizations/${organizationId}/departments/${departmentId}/teams/${deleteTeamId}`,
      { method: 'DELETE' },
    )
      .then(function () {
        toggleDeleteTeam(null)();
        notifications.success($LL.teamDeleteSuccess());
        getTeams();
      })
      .catch(function () {
        notifications.danger($LL.teamDeleteError());
      });
  }

  let defaultTeam = {
    id: '',
    name: '',
  };
  let selectedTeam = $state({ ...defaultTeam });
  let showTeamUpdate = $state(false);

  function toggleUpdateTeam(team) {
    return () => {
      selectedTeam = team;
      showTeamUpdate = !showTeamUpdate;
    };
  }

  function updateTeamHandler(name) {
    const body = {
      name,
    };

    xfetch(
      `/api/organizations/${organizationId}/departments/${departmentId}/teams/${selectedTeam.id}`,
      { body, method: 'PUT' },
    )
      .then(res => res.json())
      .then(function () {
        getTeams();
        toggleUpdateTeam(defaultTeam)();
        notifications.success(`${$LL.teamUpdateSuccess()}`);
      })
      .catch(function () {
        notifications.danger(`${$LL.teamUpdateError()}`);
      });
  }

  onMount(() => {
    if (!$user.id || !validateUserIsRegistered($user)) {
      router.route(appRoutes.login);
      return;
    }

    getDepartment();
  });

  let isAdmin = $derived(organizationRole === 'ADMIN' || departmentRole === 'ADMIN');
</script>

<svelte:head>
  <title>{$LL.department()} {department.name} | {$LL.appName()}</title>
</svelte:head>

<PageLayout>
  <div class="mb-6 lg:mb-8 dark:text-white">
    <h1 class="text-3xl font-semibold font-rajdhani">
      <span class="uppercase">{$LL.department()}</span>
      <ChevronRight class="w-8 h-8 inline-block" />
      {department.name}
    </h1>
    <div class="text-xl font-semibold font-rajdhani">
      <span class="uppercase">{$LL.organization()}</span>
      <ChevronRight class="inline-block" />
      <a
        class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
        href="{appRoutes.organization}/{organization.id}"
      >
        {organization.name}
      </a>
    </div>
  </div>

  <div class="w-full mb-6 lg:mb-8">
    <TableContainer>
      <TableNav
        title={$LL.teams()}
        createBtnEnabled={isAdmin}
        createBtnText={$LL.teamCreate()}
        createButtonHandler={toggleCreateTeam}
        createBtnTestId="team-create"
      />
      <Table>
        {#snippet header()}
                <tr >
            <HeadCol>
              {$LL.name()}
            </HeadCol>
            <HeadCol>
              {$LL.dateCreated()}
            </HeadCol>
            <HeadCol>
              {$LL.dateUpdated()}
            </HeadCol>
            <HeadCol type="action">
              <span class="sr-only">{$LL.actions()}</span>
            </HeadCol>
          </tr>
              {/snippet}
        {#snippet body({ class: className })}
                <tbody   class="{className}">
            {#each teams as team, i}
              <TableRow itemIndex={i}>
                <RowCol>
                  <a
                    href="{appRoutes.organization}/{organizationId}/department/{departmentId}/team/{team.id}"
                    class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
                  >
                    {team.name}
                  </a>
                </RowCol>
                <RowCol>
                  {new Date(team.createdDate).toLocaleString()}
                </RowCol>
                <RowCol>
                  {new Date(team.updatedDate).toLocaleString()}
                </RowCol>
                <RowCol type="action">
                  {#if isAdmin}
                    <CrudActions
                      editBtnClickHandler={toggleUpdateTeam(team)}
                      deleteBtnClickHandler={toggleDeleteTeam(team.id)}
                    />
                  {/if}
                </RowCol>
              </TableRow>
            {/each}
          </tbody>
              {/snippet}
      </Table>
    </TableContainer>
  </div>

  {#if isAdmin}
    <div class="w-full mb-6 lg:mb-8">
      <InvitesList
        xfetch={xfetch}
        notifications={notifications}
        pageType="department"
        teamPrefix={deptPrefix}
        bind:this="{invitesList}"
      />
    </div>
  {/if}

  <UsersList
    users={users}
    getUsers={getUsers}
    xfetch={xfetch}
    notifications={notifications}
    isAdmin={isAdmin}
    pageType="department"
    orgId={organizationId}
    deptId={departmentId}
    teamPrefix="/api/organizations/{organizationId}/departments/{departmentId}"
  />

  {#if showCreateTeam}
    <CreateTeam
      toggleCreate={toggleCreateTeam}
      handleCreate={createTeamHandler}
    />
  {/if}

  {#if showTeamUpdate}
    <CreateTeam
      teamName={selectedTeam.name}
      toggleCreate={toggleUpdateTeam(defaultTeam)}
      handleCreate={updateTeamHandler}
    />
  {/if}

  {#if showDeleteTeam}
    <DeleteConfirmation
      toggleDelete={toggleDeleteTeam(null)}
      handleDelete={handleDeleteTeam}
      confirmText={$LL.deleteTeamConfirmText()}
      confirmBtnText={$LL.deleteTeam()}
    />
  {/if}
</PageLayout>
