<script lang="ts">
  import { onMount } from 'svelte';

  import { user } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import { appRoutes } from '../../config';
  import { validateUserIsRegistered } from '../../validationUtils';
  import RowCol from '../../components/table/RowCol.svelte';
  import TableRow from '../../components/table/TableRow.svelte';
  import HeadCol from '../../components/table/HeadCol.svelte';
  import Table from '../../components/table/Table.svelte';
  import { ChevronRight } from '@lucide/svelte';
  import CreateTeam from '../../components/team/CreateTeam.svelte';
  import DeleteConfirmation from '../../components/global/DeleteConfirmation.svelte';
  import TableContainer from '../../components/table/TableContainer.svelte';
  import TableNav from '../../components/table/TableNav.svelte';
  import CrudActions from '../../components/table/CrudActions.svelte';

  import type { NotificationService } from '../../types/notifications';
  import type { ApiClient } from '../../types/apiclient';
  import type { Team } from '../../types/team';
  import OrgPageLayout from '../../components/organization/OrgPageLayout.svelte';

  interface Props {
    xfetch: ApiClient;
    router: any;
    notifications: NotificationService;
    organizationId: any;
  }

  let { xfetch, router, notifications, organizationId }: Props = $props();

  const teamsPageLimit = 1000;
  const orgPrefix = $derived(`/api/organizations/${organizationId}`);

  let organization = $state({
    id: '',
    name: '',
    createdDate: '',
    updateDate: '',
    subscribed: false,
  });

  $effect(() => {
    organization.id = organizationId;
  });
  let role = $state('MEMBER');
  let teams: Team[] = $state([]);
  let showCreateTeam = $state(false);
  let showDeleteTeam = $state(false);
  let deleteTeamId = $state<string | null>(null);
  let teamsPage = $state(1);

  function toggleCreateTeam() {
    showCreateTeam = !showCreateTeam;
  }

  const toggleDeleteTeam = (teamId: string | null) => () => {
    showDeleteTeam = !showDeleteTeam;
    deleteTeamId = teamId;
  };

  function getOrganization() {
    xfetch(orgPrefix)
      .then(res => res.json())
      .then(function (result) {
        organization = result.data.organization;
        role = result.data.role;

        getTeams();
      })
      .catch(function () {
        notifications.danger($LL.organizationGetError());
      });
  }

  function getTeams() {
    const teamsOffset = (teamsPage - 1) * teamsPageLimit;
    xfetch(`${orgPrefix}/teams?limit=${teamsPageLimit}&offset=${teamsOffset}`)
      .then(res => res.json())
      .then(function (result) {
        teams = result.data;
      })
      .catch(function () {
        notifications.danger($LL.organizationGetTeamsError());
      });
  }

  function createTeamHandler(name: string) {
    const body = {
      name,
    };

    xfetch(`${orgPrefix}/teams`, { body })
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
    xfetch(`${orgPrefix}/teams/${deleteTeamId}`, {
      method: 'DELETE',
    })
      .then(function () {
        toggleDeleteTeam(null)();
        notifications.success($LL.teamDeleteSuccess());
        getTeams();
      })
      .catch(function () {
        notifications.danger($LL.teamDeleteError());
      });
  }

  let defaultTeam: Team = {
    id: '',
    name: '',
    createdDate: '',
    updatedDate: '',
  };
  let selectedTeam = $state({ ...defaultTeam });
  let showTeamUpdate = $state(false);

  function toggleUpdateTeam(team: Team) {
    return () => {
      selectedTeam = team;
      showTeamUpdate = !showTeamUpdate;
    };
  }

  function updateTeamHandler(name: string) {
    const body = {
      name,
    };

    xfetch(`/api/organizations/${organizationId}/teams/${selectedTeam.id}`, {
      body,
      method: 'PUT',
    })
      .then(res => res.json())
      .then(function () {
        toggleUpdateTeam(defaultTeam)();
        getTeams();
        notifications.success(`${$LL.teamUpdateSuccess()}`);
      })
      .catch(function () {
        notifications.danger(`${$LL.teamUpdateError()}`);
      });
  }

  onMount(async () => {
    if (!$user.id || !validateUserIsRegistered($user)) {
      router.route(appRoutes.login);
      return;
    }

    getOrganization();
  });

  let isAdmin = $derived(role === 'ADMIN');
</script>

<svelte:head>
  <title>{$LL.teams()} {organization.name} | {$LL.appName()}</title>
</svelte:head>

<OrgPageLayout activePage="teams" {organizationId}>
  <h1 class="mb-4 text-3xl font-semibold font-rajdhani dark:text-white">
    <span class="uppercase">{$LL.organization()}</span>
    <ChevronRight class="w-8 h-8 inline-block" />
    {organization.name}
    <ChevronRight class="w-8 h-8 inline-block" />
    {$LL.teams()}
  </h1>

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
          <tr>
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
          <tbody class={className}>
            {#each teams as team, i}
              <TableRow itemIndex={i}>
                <RowCol>
                  <a
                    href="{appRoutes.organization}/{organizationId}/team/{team.id}"
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

  {#if showCreateTeam}
    <CreateTeam toggleCreate={toggleCreateTeam} handleCreate={createTeamHandler} />
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
</OrgPageLayout>
