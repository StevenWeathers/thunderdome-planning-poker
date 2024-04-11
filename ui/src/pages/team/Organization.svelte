<script lang="ts">
  import { onMount } from 'svelte';

  import PageLayout from '../../components/PageLayout.svelte';
  import HollowButton from '../../components/global/HollowButton.svelte';
  import { user } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import { appRoutes } from '../../config';
  import { validateUserIsRegistered } from '../../validationUtils';
  import RowCol from '../../components/table/RowCol.svelte';
  import TableRow from '../../components/table/TableRow.svelte';
  import HeadCol from '../../components/table/HeadCol.svelte';
  import Table from '../../components/table/Table.svelte';
  import ChevronRight from '../../components/icons/ChevronRight.svelte';
  import CreateDepartment from '../../components/team/CreateDepartment.svelte';
  import CreateTeam from '../../components/team/CreateTeam.svelte';
  import DeleteConfirmation from '../../components/global/DeleteConfirmation.svelte';
  import UsersList from '../../components/team/UsersList.svelte';
  import TableContainer from '../../components/table/TableContainer.svelte';
  import TableNav from '../../components/table/TableNav.svelte';
  import CrudActions from '../../components/table/CrudActions.svelte';
  import InvitesList from '../../components/team/InvitesList.svelte';

  export let xfetch;
  export let router;
  export let notifications;
  export let eventTag;
  export let organizationId;

  const departmentsPageLimit = 1000;
  const teamsPageLimit = 1000;
  const usersPageLimit = 1000;
  const orgPrefix = `/api/organizations/${organizationId}`;

  let invitesList;
  let organization = {
    id: organizationId,
    name: '',
    createdDate: '',
    updateDate: '',
  };
  let role = 'MEMBER';
  let users = [];
  let departments = [];
  let teams = [];
  let invites = [];
  let showCreateDepartment = false;
  let showCreateTeam = false;
  let showDeleteTeam = false;
  let showDeleteDepartment = false;
  let showDeleteOrganization = false;
  let deleteTeamId = null;
  let deleteDeptId = null;
  let teamsPage = 1;
  let departmentsPage = 1;
  let usersPage = 1;

  function toggleCreateDepartment() {
    showCreateDepartment = !showCreateDepartment;
  }

  function toggleCreateTeam() {
    showCreateTeam = !showCreateTeam;
  }

  const toggleDeleteTeam = teamId => () => {
    showDeleteTeam = !showDeleteTeam;
    deleteTeamId = teamId;
  };

  const toggleDeleteDepartment = deptId => () => {
    showDeleteDepartment = !showDeleteDepartment;
    deleteDeptId = deptId;
  };

  const toggleDeleteOrganization = () => {
    showDeleteOrganization = !showDeleteOrganization;
  };

  function getOrganization() {
    xfetch(`/api/organizations/${organizationId}`)
      .then(res => res.json())
      .then(function (result) {
        organization = result.data.organization;
        role = result.data.role;

        getDepartments();
        getTeams();
        getUsers();
      })
      .catch(function () {
        notifications.danger($LL.organizationGetError());
      });
  }

  function getUsers() {
    const usersOffset = (usersPage - 1) * usersPageLimit;
    xfetch(`${orgPrefix}/users?limit=${usersPageLimit}&offset=${usersOffset}`)
      .then(res => res.json())
      .then(function (result) {
        users = result.data;
      })
      .catch(function () {
        notifications.danger($LL.teamGetUsersError());
      });
  }

  function getDepartments() {
    const departmentsOffset = (departmentsPage - 1) * departmentsPageLimit;
    xfetch(
      `${orgPrefix}/departments?limit=${departmentsPageLimit}&offset=${departmentsOffset}`,
    )
      .then(res => res.json())
      .then(function (result) {
        departments = result.data;
      })
      .catch(function () {
        notifications.danger($LL.organizationGetDepartmentsError());
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

  function createDepartmentHandler(name) {
    const body = {
      name,
    };

    xfetch(`${orgPrefix}/departments`, { body })
      .then(res => res.json())
      .then(function (result) {
        eventTag('create_department', 'engagement', 'success', () => {
          router.route(
            `${appRoutes.organization}/${organizationId}/department/${result.data.id}`,
          );
        });
      })
      .catch(function () {
        notifications.danger($LL.departmentCreateError());
        eventTag('create_department', 'engagement', 'failure');
      });
  }

  function createTeamHandler(name) {
    const body = {
      name,
    };

    xfetch(`${orgPrefix}/teams`, { body })
      .then(res => res.json())
      .then(function () {
        eventTag('create_organization_team', 'engagement', 'success');
        toggleCreateTeam();
        notifications.success($LL.teamCreateSuccess());
        getTeams();
      })
      .catch(function () {
        notifications.danger($LL.teamCreateError());
        eventTag('create_organization_team', 'engagement', 'failure');
      });
  }

  function handleDeleteTeam() {
    xfetch(`${orgPrefix}/teams/${deleteTeamId}`, {
      method: 'DELETE',
    })
      .then(function () {
        eventTag('organization_delete_team', 'engagement', 'success');
        toggleDeleteTeam(null)();
        notifications.success($LL.teamDeleteSuccess());
        getTeams();
      })
      .catch(function () {
        notifications.danger($LL.teamDeleteError());
        eventTag('organization_delete_team', 'engagement', 'failure');
      });
  }

  function handleDeleteDepartment() {
    xfetch(`${orgPrefix}/departments/${deleteDeptId}`, {
      method: 'DELETE',
    })
      .then(function () {
        eventTag('organization_delete_department', 'engagement', 'success');
        toggleDeleteDepartment(null)();
        notifications.success($LL.departmentDeleteSuccess());
        getDepartments();
      })
      .catch(function () {
        notifications.danger($LL.departmentDeleteError());
        eventTag('organization_delete_department', 'engagement', 'failure');
      });
  }

  function handleDeleteOrganization() {
    xfetch(`${orgPrefix}`, {
      method: 'DELETE',
    })
      .then(function () {
        eventTag('organization_delete', 'engagement', 'success');
        toggleDeleteTeam();
        notifications.success($LL.organizationDeleteSuccess());
        router.route(appRoutes.teams);
      })
      .catch(function () {
        notifications.danger($LL.organizationDeleteError());
        eventTag('organization_delete', 'engagement', 'failure');
      });
  }

  let defaultDepartment = {
    id: '',
    name: '',
  };
  let selectedDepartment = { ...defaultDepartment };
  let showDepartmentUpdate = false;

  function toggleUpdateDepartment(dept) {
    return () => {
      selectedDepartment = dept;
      showDepartmentUpdate = !showDepartmentUpdate;
    };
  }

  let defaultTeam = {
    id: '',
    name: '',
  };
  let selectedTeam = { ...defaultTeam };
  let showTeamUpdate = false;

  function toggleUpdateTeam(team) {
    return () => {
      selectedTeam = team;
      showTeamUpdate = !showTeamUpdate;
    };
  }

  function updateDepartmentHandler(name) {
    const body = {
      name,
    };

    xfetch(
      `/api/organizations/${organizationId}/departments/${selectedDepartment.id}`,
      { body, method: 'PUT' },
    )
      .then(res => res.json())
      .then(function (result) {
        eventTag('update_department', 'engagement', 'success');
        getDepartments();
        toggleUpdateDepartment(defaultDepartment)();
        notifications.success(`${$LL.deptUpdateSuccess()}`);
      })
      .catch(function () {
        notifications.danger(`${$LL.deptUpdateError()}`);
        eventTag('update_department', 'engagement', 'failure');
      });
  }

  function updateTeamHandler(name) {
    const body = {
      name,
    };

    xfetch(`/api/organizations/${organizationId}/teams/${selectedTeam.id}`, {
      body,
      method: 'PUT',
    })
      .then(res => res.json())
      .then(function () {
        eventTag('create_organization_team', 'engagement', 'success');
        toggleUpdateTeam(defaultTeam)();
        getTeams();
        notifications.success(`${$LL.teamUpdateSuccess()}`);
      })
      .catch(function () {
        notifications.danger(`${$LL.teamUpdateError()}`);
        eventTag('create_organization_team', 'engagement', 'failure');
      });
  }

  onMount(() => {
    if (!$user.id || !validateUserIsRegistered($user)) {
      router.route(appRoutes.login);
      return;
    }

    getOrganization();
  });

  $: isAdmin = role === 'ADMIN';
</script>

<svelte:head>
  <title>{$LL.organization()} {organization.name} | {$LL.appName()}</title>
</svelte:head>

<PageLayout>
  <h1 class="mb-4 text-3xl font-semibold font-rajdhani dark:text-white">
    <span class="uppercase">{$LL.organization()}</span>
    <ChevronRight class="w-8 h-8" />
    {organization.name}
  </h1>

  <div class="w-full mb-6 lg:mb-8">
    <TableContainer>
      <TableNav
        title="{$LL.departments()}"
        createBtnEnabled="{isAdmin}"
        createBtnText="{$LL.departmentCreate()}"
        createButtonHandler="{toggleCreateDepartment}"
        createBtnTestId="department-create"
      />
      <Table>
        <tr slot="header">
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
        <tbody slot="body" let:class="{className}" class="{className}">
          {#each departments as department, i}
            <TableRow itemIndex="{i}">
              <RowCol>
                <a
                  href="{appRoutes.organization}/{organizationId}/department/{department.id}"
                  class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
                >
                  {department.name}
                </a>
              </RowCol>
              <RowCol>
                {new Date(department.createdDate).toLocaleString()}
              </RowCol>
              <RowCol>
                {new Date(department.updatedDate).toLocaleString()}
              </RowCol>
              <RowCol type="action">
                {#if isAdmin}
                  <CrudActions
                    editBtnClickHandler="{toggleUpdateDepartment(department)}"
                    deleteBtnClickHandler="{toggleDeleteDepartment(
                      department.id,
                    )}"
                  />
                {/if}
              </RowCol>
            </TableRow>
          {/each}
        </tbody>
      </Table>
    </TableContainer>
  </div>

  <div class="w-full mb-6 lg:mb-8">
    <TableContainer>
      <TableNav
        title="{$LL.teams()}"
        createBtnEnabled="{isAdmin}"
        createBtnText="{$LL.teamCreate()}"
        createButtonHandler="{toggleCreateTeam}"
        createBtnTestId="team-create"
      />
      <Table>
        <tr slot="header">
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
        <tbody slot="body" let:class="{className}" class="{className}">
          {#each teams as team, i}
            <TableRow itemIndex="{i}">
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
                    editBtnClickHandler="{toggleUpdateTeam(team)}"
                    deleteBtnClickHandler="{toggleDeleteTeam(team.id)}"
                  />
                {/if}
              </RowCol>
            </TableRow>
          {/each}
        </tbody>
      </Table>
    </TableContainer>
  </div>

  {#if isAdmin}
    <div class="w-full mb-6 lg:mb-8">
      <InvitesList
        xfetch="{xfetch}"
        eventTag="{eventTag}"
        notifications="{notifications}"
        pageType="organization"
        teamPrefix="{orgPrefix}"
        bind:this="{invitesList}"
      />
    </div>
  {/if}

  <UsersList
    users="{users}"
    getUsers="{getUsers}"
    xfetch="{xfetch}"
    eventTag="{eventTag}"
    notifications="{notifications}"
    isAdmin="{isAdmin}"
    pageType="organization"
    teamPrefix="/api/organizations/{organizationId}"
    on:user-invited="{() => {
      invitesList.f('user-invited');
    }}"
  />

  {#if isAdmin}
    <div class="w-full text-center mt-8">
      <HollowButton onClick="{toggleDeleteOrganization}" color="red">
        {$LL.deleteOrganization()}
      </HollowButton>
    </div>
  {/if}

  {#if showCreateDepartment}
    <CreateDepartment
      toggleCreate="{toggleCreateDepartment}"
      handleCreate="{createDepartmentHandler}"
    />
  {/if}

  {#if showDepartmentUpdate}
    <CreateDepartment
      departmentName="{selectedDepartment.name}"
      toggleCreate="{toggleUpdateDepartment(defaultDepartment)}"
      handleCreate="{updateDepartmentHandler}"
    />
  {/if}

  {#if showCreateTeam}
    <CreateTeam
      toggleCreate="{toggleCreateTeam}"
      handleCreate="{createTeamHandler}"
    />
  {/if}

  {#if showTeamUpdate}
    <CreateTeam
      teamName="{selectedTeam.name}"
      toggleCreate="{toggleUpdateTeam(defaultTeam)}"
      handleCreate="{updateTeamHandler}"
    />
  {/if}

  {#if showDeleteTeam}
    <DeleteConfirmation
      toggleDelete="{toggleDeleteTeam(null)}"
      handleDelete="{handleDeleteTeam}"
      confirmText="{$LL.deleteTeamConfirmText()}"
      confirmBtnText="{$LL.deleteTeam()}"
    />
  {/if}

  {#if showDeleteDepartment}
    <DeleteConfirmation
      toggleDelete="{toggleDeleteDepartment(null)}"
      handleDelete="{handleDeleteDepartment}"
      confirmText="{$LL.deleteDepartmentConfirmText()}"
      confirmBtnText="{$LL.deleteDepartment()}"
    />
  {/if}

  {#if showDeleteOrganization}
    <DeleteConfirmation
      toggleDelete="{toggleDeleteOrganization}"
      handleDelete="{handleDeleteOrganization}"
      confirmText="{$LL.deleteOrganizationConfirmText()}"
      confirmBtnText="{$LL.deleteOrganization()}"
    />
  {/if}
</PageLayout>
