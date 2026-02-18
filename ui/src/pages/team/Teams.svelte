<script lang="ts">
  import { onMount } from 'svelte';

  import PageLayout from '../../components/PageLayout.svelte';
  import CreateOrganization from '../../components/team/CreateOrganization.svelte';
  import CreateTeam from '../../components/team/CreateTeam.svelte';
  import { user } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import { AppConfig, appRoutes } from '../../config';
  import { validateUserIsRegistered } from '../../validationUtils';
  import RowCol from '../../components/table/RowCol.svelte';
  import TableRow from '../../components/table/TableRow.svelte';
  import HeadCol from '../../components/table/HeadCol.svelte';
  import Table from '../../components/table/Table.svelte';
  import TableNav from '../../components/table/TableNav.svelte';
  import TableContainer from '../../components/table/TableContainer.svelte';
  import CrudActions from '../../components/table/CrudActions.svelte';
  import DeleteConfirmation from '../../components/global/DeleteConfirmation.svelte';

  import type { NotificationService } from '../../types/notifications';
  import type { ApiClient } from '../../types/apiclient';
  import type { Organization } from '../../types/organization';
  import type { Team } from '../../types/team';

  interface Props {
    xfetch: ApiClient;
    router: any;
    notifications: NotificationService;
  }

  let { xfetch, router, notifications }: Props = $props();

  const organizationsPageLimit = 1000;
  const teamsPageLimit = 1000;
  const { OrganizationsEnabled } = AppConfig;

  let defaultOrganization: Organization = {
    id: '',
    name: '',
    createdDate: '',
    updatedDate: '',
  };

  let organizations: Organization[] = $state([]);
  let teams: Team[] = $state([]);
  let showCreateOrganization = $state(false);
  let showCreateTeam = $state(false);
  let organizationsPage = 1;
  let teamsPage = 1;

  function toggleCreateOrganization() {
    showCreateOrganization = !showCreateOrganization;
  }

  let showOrganizationUpdate = $state(false);
  let selectedOrganization = $state({ ...defaultOrganization });

  function toggleUpdateOrganization(selectedOrg: Organization) {
    return () => {
      selectedOrganization = selectedOrg;
      showOrganizationUpdate = !showOrganizationUpdate;
    };
  }

  function toggleCreateTeam() {
    showCreateTeam = !showCreateTeam;
  }

  let showDeleteOrganization = $state(false);
  const toggleDeleteOrganization = (selectedOrg: Organization) => () => {
    selectedOrganization = selectedOrg;
    showDeleteOrganization = !showDeleteOrganization;
  };

  function getOrganizations() {
    const orgsOffset = (organizationsPage - 1) * organizationsPageLimit;
    xfetch(`/api/users/${$user.id}/organizations?limit=${organizationsPageLimit}&offset=${orgsOffset}`)
      .then(res => res.json())
      .then(function (result) {
        organizations = result.data;
      })
      .catch(function () {
        notifications.danger($LL.getOrganizationsError());
      });
  }

  function getTeams() {
    const teamsOffset = (teamsPage - 1) * teamsPageLimit;
    xfetch(`/api/users/${$user.id}/teams-non-org?limit=${teamsPageLimit}&offset=${teamsOffset}`)
      .then(res => res.json())
      .then(function (result) {
        teams = result.data;
      })
      .catch(function () {
        notifications.danger($LL.getTeamsError());
      });
  }

  function createOrganizationHandler(name: string) {
    const body = {
      name,
    };

    xfetch(`/api/users/${$user.id}/organizations`, { body })
      .then(res => res.json())
      .then(function (result) {
        router.route(`${appRoutes.organization}/${result.data.id}`);
      })
      .catch(function () {
        notifications.danger($LL.createOrgError());
      });
  }

  function createTeamHandler(name: string) {
    const body = {
      name,
    };

    xfetch(`/api/users/${$user.id}/teams`, { body })
      .then(res => res.json())
      .then(function (result) {
        router.route(`${appRoutes.team}/${result.data.id}`);
      })
      .catch(function () {
        notifications.danger($LL.teamCreateError());
      });
  }

  function updateOrganizationHandler(name: string) {
    const body = {
      name,
    };

    xfetch(`/api/organizations/${selectedOrganization.id}`, {
      method: 'PUT',
      body,
    })
      .then(res => res.json())
      .then(function (result) {
        notifications.success(`${$LL.orgUpdateSuccess()}`);
        getOrganizations();
        toggleUpdateOrganization(defaultOrganization)();
      })
      .catch(function () {
        notifications.danger(`${$LL.orgUpdateError()}`);
      });
  }

  function handleDeleteOrganization() {
    xfetch(`/api/organizations/${selectedOrganization.id}`, {
      method: 'DELETE',
    })
      .then(function () {
        getOrganizations();
        toggleDeleteOrganization(defaultOrganization)();
        notifications.success($LL.organizationDeleteSuccess());
      })
      .catch(function () {
        notifications.danger($LL.organizationDeleteError());
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

  let showDeleteTeam = $state(false);
  const toggleDeleteTeam = (team: Team) => () => {
    selectedTeam = team;
    showDeleteTeam = !showDeleteTeam;
  };

  function updateTeamHandler(name: string) {
    const body = {
      name,
    };

    xfetch(`/api/teams/${selectedTeam.id}`, {
      method: 'PUT',
      body,
    })
      .then(res => res.json())
      .then(function (result) {
        notifications.success(`${$LL.teamUpdateSuccess()}`);
        getTeams();
        toggleUpdateTeam(defaultTeam)();
      })
      .catch(function () {
        notifications.danger(`${$LL.teamUpdateError()}`);
      });
  }

  function handleDeleteTeam() {
    xfetch(`/api/teams/${selectedTeam.id}`, {
      method: 'DELETE',
    })
      .then(function () {
        getTeams();
        toggleDeleteTeam(defaultTeam)();
        notifications.success($LL.teamDeleteSuccess());
      })
      .catch(function () {
        notifications.danger($LL.teamDeleteError());
      });
  }

  onMount(() => {
    if (!$user.id || !validateUserIsRegistered($user)) {
      router.route(appRoutes.login);
      return;
    }

    if (OrganizationsEnabled) {
      getOrganizations();
    }
    getTeams();
  });
</script>

<svelte:head>
  <title>{$LL.organizations()} | {$LL.appName()}</title>
</svelte:head>

<PageLayout>
  {#if OrganizationsEnabled}
    <div class="w-full mb-6 lg:mb-8">
      <TableContainer>
        <TableNav
          title={$LL.organizations()}
          createBtnText={$LL.organizationCreate()}
          createButtonHandler={toggleCreateOrganization}
          createBtnTestId="organization-create"
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
              <HeadCol type="action" />
            </tr>
          {/snippet}
          {#snippet body({ class: className })}
            <tbody class={className}>
              {#each organizations as organization, i}
                <TableRow itemIndex={i}>
                  <RowCol>
                    <a
                      href="{appRoutes.organization}/{organization.id}"
                      class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
                    >
                      {organization.name}
                    </a>
                  </RowCol>
                  <RowCol>
                    {new Date(organization.createdDate).toLocaleString()}
                  </RowCol>
                  <RowCol>
                    {new Date(organization.updatedDate).toLocaleString()}
                  </RowCol>
                  <RowCol type="action">
                    {#if (organization as any).role === 'ADMIN'}
                      <CrudActions
                        editBtnClickHandler={toggleUpdateOrganization(organization)}
                        deleteBtnClickHandler={toggleDeleteOrganization(organization)}
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
  {/if}

  <TableContainer>
    <TableNav
      title={$LL.teams()}
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
          <HeadCol type="action" />
        </tr>
      {/snippet}
      {#snippet body({ class: className })}
        <tbody class={className}>
          {#each teams as team, i}
            <TableRow itemIndex={i}>
              <RowCol>
                <a
                  href="{appRoutes.team}/{team.id}"
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
                {#if (team as any).role === 'ADMIN'}
                  <CrudActions
                    editBtnClickHandler={toggleUpdateTeam(team)}
                    deleteBtnClickHandler={toggleDeleteTeam(team)}
                  />
                {/if}
              </RowCol>
            </TableRow>
          {/each}
        </tbody>
      {/snippet}
    </Table>
  </TableContainer>

  {#if showCreateOrganization}
    <CreateOrganization toggleCreate={toggleCreateOrganization} handleCreate={createOrganizationHandler} />
  {/if}

  {#if showOrganizationUpdate}
    <CreateOrganization
      toggleCreate={toggleUpdateOrganization(defaultOrganization)}
      organizationName={selectedOrganization.name}
      handleCreate={updateOrganizationHandler}
    />
  {/if}

  {#if showDeleteOrganization}
    <DeleteConfirmation
      toggleDelete={toggleDeleteOrganization(defaultOrganization)}
      handleDelete={handleDeleteOrganization}
      confirmText={$LL.deleteOrganizationConfirmText()}
      confirmBtnText={$LL.deleteOrganization()}
    />
  {/if}

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
      toggleDelete={toggleDeleteTeam(defaultTeam)}
      handleDelete={handleDeleteTeam}
      confirmText={$LL.deleteTeamConfirmText()}
      confirmBtnText={$LL.deleteTeam()}
    />
  {/if}
</PageLayout>
