<script lang="ts">
  import { onMount } from 'svelte';

  import PageLayout from '../../components/PageLayout.svelte';
  import HollowButton from '../../components/HollowButton.svelte';
  import SolidButton from '../../components/SolidButton.svelte';
  import { user } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import { appRoutes } from '../../config';
  import { validateUserIsRegistered } from '../../validationUtils';
  import RowCol from '../../components/table/RowCol.svelte';
  import TableRow from '../../components/table/TableRow.svelte';
  import HeadCol from '../../components/table/HeadCol.svelte';
  import Table from '../../components/table/Table.svelte';
  import ChevronRight from '../../components/icons/ChevronRight.svelte';
  import UserAvatar from '../../components/user/UserAvatar.svelte';
  import CountryFlag from '../../components/user/CountryFlag.svelte';
  import CreateDepartment from '../../components/team/CreateDepartment.svelte';
  import CreateTeam from '../../components/team/CreateTeam.svelte';
  import AddUser from '../../components/user/AddUser.svelte';
  import DeleteConfirmation from '../../components/DeleteConfirmation.svelte';

  export let xfetch;
  export let router;
  export let notifications;
  export let eventTag;
  export let organizationId;

  const departmentsPageLimit = 1000;
  const teamsPageLimit = 1000;
  const usersPageLimit = 1000;

  let organization = {
    id: organizationId,
    name: '',
    createdDate: '',
    updateDate: '',
  };
  let role = 'MEMBER';
  let departments = [];
  let teams = [];
  let users = [];
  let showCreateDepartment = false;
  let showCreateTeam = false;
  let showAddUser = false;
  let showRemoveUser = false;
  let removeUserId = null;
  let showDeleteTeam = false;
  let showDeleteDepartment = false;
  let showDeleteOrganization = false;
  let deleteTeamId = null;
  let deleteDeptId = null;
  let teamsPage = 1;
  let usersPage = 1;
  let departmentsPage = 1;

  function toggleCreateDepartment() {
    showCreateDepartment = !showCreateDepartment;
  }

  function toggleCreateTeam() {
    showCreateTeam = !showCreateTeam;
  }

  function toggleAddUser() {
    showAddUser = !showAddUser;
  }

  const toggleRemoveUser = userId => () => {
    showRemoveUser = !showRemoveUser;
    removeUserId = userId;
  };

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

  function getDepartments() {
    const departmentsOffset = (departmentsPage - 1) * departmentsPageLimit;
    xfetch(
      `/api/organizations/${organizationId}/departments?limit=${departmentsPageLimit}&offset=${departmentsOffset}`,
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
    xfetch(
      `/api/organizations/${organizationId}/teams?limit=${teamsPageLimit}&offset=${teamsOffset}`,
    )
      .then(res => res.json())
      .then(function (result) {
        teams = result.data;
      })
      .catch(function () {
        notifications.danger($LL.organizationGetTeamsError());
      });
  }

  function getUsers() {
    const usersOffset = (usersPage - 1) * usersPageLimit;
    xfetch(
      `/api/organizations/${organizationId}/users?limit=${usersPageLimit}&offset=${usersOffset}`,
    )
      .then(res => res.json())
      .then(function (result) {
        users = result.data;
      })
      .catch(function () {
        notifications.danger($LL.organizationGetUsersError());
      });
  }

  function createDepartmentHandler(name) {
    const body = {
      name,
    };

    xfetch(`/api/organizations/${organizationId}/departments`, { body })
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

    xfetch(`/api/organizations/${organizationId}/teams`, { body })
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

  function handleUserAdd(email, role) {
    const body = {
      email,
      role,
    };

    xfetch(`/api/organizations/${organizationId}/users`, { body })
      .then(function () {
        eventTag('organization_add_user', 'engagement', 'success');
        toggleAddUser();
        notifications.success($LL.userAddSuccess());
        getUsers();
      })
      .catch(function () {
        notifications.danger($LL.userAddError());
        eventTag('organization_add_user', 'engagement', 'failure');
      });
  }

  function handleUserRemove() {
    xfetch(`/api/organizations/${organizationId}/users/${removeUserId}`, {
      method: 'DELETE',
    })
      .then(function () {
        eventTag('organization_remove_user', 'engagement', 'success');
        toggleRemoveUser(null)();
        notifications.success($LL.userRemoveSuccess());
        getUsers();
      })
      .catch(function () {
        notifications.danger($LL.userRemoveError());
        eventTag('organization_remove_user', 'engagement', 'failure');
      });
  }

  function handleDeleteTeam() {
    xfetch(`/api/organizations/${organizationId}/teams/${deleteTeamId}`, {
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
    xfetch(`/api/organizations/${organizationId}/departments/${deleteDeptId}`, {
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
    xfetch(`/api/organizations/${organizationId}`, {
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
    <div class="flex w-full">
      <div class="w-4/5">
        <h2
          class="text-2xl font-semibold font-rajdhani uppercase mb-4 dark:text-white"
        >
          {$LL.departments()}
        </h2>
      </div>
      <div class="w-1/5">
        <div class="text-right">
          {#if isAdmin}
            <SolidButton onClick="{toggleCreateDepartment}">
              {$LL.departmentCreate()}
            </SolidButton>
          {/if}
        </div>
      </div>
    </div>

    <div class="w-full">
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
                  <HollowButton
                    onClick="{toggleDeleteDepartment(department.id)}"
                    color="red"
                  >
                    {$LL.delete()}
                  </HollowButton>
                {/if}
              </RowCol>
            </TableRow>
          {/each}
        </tbody>
      </Table>
    </div>
  </div>

  <div class="w-full mb-6 lg:mb-8">
    <div class="flex w-full">
      <div class="w-4/5">
        <h2
          class="text-2xl font-semibold font-rajdhani uppercase mb-4 dark:text-white"
        >
          {$LL.teams()}
        </h2>
      </div>
      <div class="w-1/5">
        <div class="text-right">
          {#if isAdmin}
            <SolidButton onClick="{toggleCreateTeam}">
              {$LL.teamCreate()}
            </SolidButton>
          {/if}
        </div>
      </div>
    </div>

    <div class="w-full">
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
                  <HollowButton
                    onClick="{toggleDeleteTeam(team.id)}"
                    color="red"
                  >
                    {$LL.delete()}
                  </HollowButton>
                {/if}
              </RowCol>
            </TableRow>
          {/each}
        </tbody>
      </Table>
    </div>
  </div>

  <div class="w-full">
    <div class="flex w-full">
      <div class="w-4/5">
        <h2
          class="text-2xl font-semibold font-rajdhani uppercase mb-4 dark:text-white"
        >
          {$LL.users()}
        </h2>
      </div>
      <div class="w-1/5">
        <div class="text-right">
          {#if isAdmin}
            <SolidButton onClick="{toggleAddUser}" testid="user-add">
              {$LL.userAdd()}
            </SolidButton>
          {/if}
        </div>
      </div>
    </div>

    <Table>
      <tr slot="header">
        <HeadCol>
          {$LL.name()}
        </HeadCol>
        <HeadCol>
          {$LL.email()}
        </HeadCol>
        <HeadCol>
          {$LL.role()}
        </HeadCol>
        <HeadCol type="action">
          <span class="sr-only">{$LL.actions()}</span>
        </HeadCol>
      </tr>
      <tbody slot="body" let:class="{className}" class="{className}">
        {#each users as user, i}
          <TableRow itemIndex="{i}">
            <RowCol>
              <div class="flex items-center">
                <div class="flex-shrink-0 h-10 w-10">
                  <UserAvatar
                    warriorId="{user.id}"
                    avatar="{user.avatar}"
                    gravatarHash="{user.gravatarHash}"
                    width="48"
                    class="h-10 w-10 rounded-full"
                  />
                </div>
                <div class="ms-4">
                  <div class="font-medium text-gray-900 dark:text-gray-200">
                    <span data-testid="user-name">{user.name}</span>
                    {#if user.country}
                      &nbsp;
                      <CountryFlag
                        country="{user.country}"
                        additionalClass="inline-block"
                        width="32"
                        height="24"
                      />
                    {/if}
                  </div>
                </div>
              </div>
            </RowCol>
            <RowCol>
              <span data-testid="user-email">{user.email}</span>
            </RowCol>
            <RowCol>
              <div class="text-sm text-gray-500 dark:text-gray-300">
                {user.role}
              </div>
            </RowCol>
            <RowCol type="action">
              {#if isAdmin}
                <HollowButton onClick="{toggleRemoveUser(user.id)}" color="red">
                  {$LL.remove()}
                </HollowButton>
              {/if}
            </RowCol>
          </TableRow>
        {/each}
      </tbody>
    </Table>
  </div>

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

  {#if showCreateTeam}
    <CreateTeam
      toggleCreate="{toggleCreateTeam}"
      handleCreate="{createTeamHandler}"
    />
  {/if}

  {#if showAddUser}
    <AddUser toggleAdd="{toggleAddUser}" handleAdd="{handleUserAdd}" />
  {/if}

  {#if showRemoveUser}
    <DeleteConfirmation
      toggleDelete="{toggleRemoveUser(null)}"
      handleDelete="{handleUserRemove}"
      permanent="{false}"
      confirmText="{$LL.removeUserConfirmText()}"
      confirmBtnText="{$LL.removeUser()}"
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
