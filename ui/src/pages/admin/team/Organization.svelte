<script lang="ts">
  import { onMount } from 'svelte';

  import AdminPageLayout from '../../../components/admin/AdminPageLayout.svelte';
  import HollowButton from '../../../components/global/HollowButton.svelte';
  import DeleteConfirmation from '../../../components/global/DeleteConfirmation.svelte';
  import { user } from '../../../stores';
  import LL from '../../../i18n/i18n-svelte';
  import { appRoutes } from '../../../config';
  import RowCol from '../../../components/table/RowCol.svelte';
  import TableRow from '../../../components/table/TableRow.svelte';
  import HeadCol from '../../../components/table/HeadCol.svelte';
  import Table from '../../../components/table/Table.svelte';
  import { validateUserIsAdmin } from '../../../validationUtils';
  import UserAvatar from '../../../components/user/UserAvatar.svelte';
  import CountryFlag from '../../../components/user/CountryFlag.svelte';
  import { ChevronRight } from '@lucide/svelte';
  import TableContainer from '../../../components/table/TableContainer.svelte';
  import TableNav from '../../../components/table/TableNav.svelte';
  import CrudActions from '../../../components/table/CrudActions.svelte';

  import type { NotificationService } from '../../../types/notifications';
  import type { ApiClient } from '../../../types/apiclient';

  interface Props {
    xfetch: ApiClient;
    router: any;
    notifications: NotificationService;
    organizationId: any;
  }

  let { xfetch, router, notifications, organizationId }: Props = $props();

  const departmentsPageLimit = 1000;
  const teamsPageLimit = 1000;
  const usersPageLimit = 1000;

  let organization = $state({
    id: '',
    name: '',
    createdDate: '',
    updateDate: '',
  });
  let role = 'MEMBER';
  let departments = $state([]);
  let teams = $state([]);
  let users = $state([]);
  let showDeleteTeam = $state(false);
  let showDeleteDepartment = $state(false);
  let showDeleteOrganization = $state(false);
  let deleteTeamId: string | null = null;
  let deleteDeptId: string | null = null;
  let teamsPage = 1;
  let usersPage = 1;
  let departmentsPage = 1;

  $effect(() => {
    organization.id = organizationId;
  });

  const apiPrefix = $derived(`/api/organizations/${organizationId}`);

  const toggleDeleteTeam = (teamId: string | null) => () => {
    showDeleteTeam = !showDeleteTeam;
    deleteTeamId = teamId;
  };

  const toggleDeleteDepartment = (deptId: string | null) => () => {
    showDeleteDepartment = !showDeleteDepartment;
    deleteDeptId = deptId;
  };

  const toggleDeleteOrganization = () => {
    showDeleteOrganization = !showDeleteOrganization;
  };

  function getOrganization() {
    xfetch(`${apiPrefix}`)
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
    xfetch(`${apiPrefix}/departments?limit=${departmentsPageLimit}&offset=${departmentsOffset}`)
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
    xfetch(`${apiPrefix}/teams?limit=${teamsPageLimit}&offset=${teamsOffset}`)
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
    xfetch(`${apiPrefix}/users?limit=${usersPageLimit}&offset=${usersOffset}`)
      .then(res => res.json())
      .then(function (result) {
        users = result.data;
      })
      .catch(function () {
        notifications.danger($LL.organizationGetUsersError());
      });
  }

  function handleDeleteTeam() {
    xfetch(`${apiPrefix}/teams/${deleteTeamId}`, {
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

  function handleDeleteDepartment() {
    xfetch(`${apiPrefix}/departments/${deleteDeptId}`, {
      method: 'DELETE',
    })
      .then(function () {
        toggleDeleteDepartment(null)();
        notifications.success($LL.departmentDeleteSuccess());
        getDepartments();
      })
      .catch(function () {
        notifications.danger($LL.departmentDeleteError());
      });
  }

  function handleDeleteOrganization() {
    xfetch(`${apiPrefix}`, {
      method: 'DELETE',
    })
      .then(function () {
        toggleDeleteOrganization();
        notifications.success($LL.organizationDeleteSuccess());
        router.route(appRoutes.adminOrganizations);
      })
      .catch(function () {
        notifications.danger($LL.organizationDeleteError());
      });
  }

  onMount(() => {
    if (!$user.id) {
      router.route(appRoutes.login);
      return;
    }
    if (!validateUserIsAdmin($user)) {
      router.route(appRoutes.landing);
      return;
    }

    getOrganization();
  });
</script>

<svelte:head>
  <title>{$LL.organization()} {organization.name} | {$LL.appName()}</title>
</svelte:head>

<AdminPageLayout activePage="organizations">
  <h1 class="mb-4 text-3xl font-semibold font-rajdhani dark:text-white">
    <span class="uppercase">{$LL.organization()}</span>
    <ChevronRight class="w-8 h-8 inline-block" />
    {organization.name}
  </h1>

  <div class="w-full mb-6 lg:mb-8">
    <TableContainer>
      <TableNav title={$LL.departments()} createBtnEnabled={false} />
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
            {#each departments as department, i}
              <TableRow itemIndex={i}>
                <RowCol>
                  <a
                    href="{appRoutes.adminOrganizations}/{organizationId}/department/{department.id}"
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
                  <CrudActions editBtnEnabled={false} deleteBtnClickHandler={toggleDeleteDepartment(department.id)} />
                </RowCol>
              </TableRow>
            {/each}
          </tbody>
        {/snippet}
      </Table>
    </TableContainer>
  </div>

  <div class="w-full mb-6 lg:mb-8">
    <TableContainer>
      <TableNav title={$LL.teams()} createBtnEnabled={false} />
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
                    href="{appRoutes.adminOrganizations}/{organizationId}/team/{team.id}"
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
                  <CrudActions editBtnEnabled={false} deleteBtnClickHandler={toggleDeleteTeam(team.id)} />
                </RowCol>
              </TableRow>
            {/each}
          </tbody>
        {/snippet}
      </Table>
    </TableContainer>
  </div>

  <TableContainer>
    <TableNav title={$LL.users()} createBtnEnabled={false} />
    <Table>
      {#snippet header()}
        <tr>
          <HeadCol>
            {$LL.name()}
          </HeadCol>
          <HeadCol>
            {$LL.email()}
          </HeadCol>
          <HeadCol>
            {$LL.role()}
          </HeadCol>
        </tr>
      {/snippet}
      {#snippet body({ class: className })}
        <tbody class={className}>
          {#each users as user, i}
            <TableRow itemIndex={i}>
              <RowCol>
                <div class="flex items-center">
                  <div class="flex-shrink-0 h-10 w-10">
                    <UserAvatar
                      warriorId={user.id}
                      avatar={user.avatar}
                      gravatarHash={user.gravatarHash}
                      userName={user.name}
                      width={48}
                      class="h-10 w-10 rounded-full"
                    />
                  </div>
                  <div class="ms-4">
                    <div class="font-medium text-gray-900 dark:text-gray-200">
                      <a
                        data-testid="user-name"
                        href="{appRoutes.adminUsers}/{user.id}"
                        class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
                        >{user.name}</a
                      >
                      {#if user.country}
                        &nbsp;
                        <CountryFlag country={user.country} additionalClass="inline-block" width="32" height="24" />
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
            </TableRow>
          {/each}
        </tbody>
      {/snippet}
    </Table>
  </TableContainer>

  <div class="w-full text-center mt-8">
    <HollowButton onClick={toggleDeleteOrganization} color="red">
      {$LL.deleteOrganization()}
    </HollowButton>
  </div>

  {#if showDeleteTeam}
    <DeleteConfirmation
      toggleDelete={toggleDeleteTeam(null)}
      handleDelete={handleDeleteTeam}
      confirmText={$LL.deleteTeamConfirmText()}
      confirmBtnText={$LL.deleteTeam()}
    />
  {/if}

  {#if showDeleteDepartment}
    <DeleteConfirmation
      toggleDelete={toggleDeleteDepartment(null)}
      handleDelete={handleDeleteDepartment}
      confirmText={$LL.deleteDepartmentConfirmText()}
      confirmBtnText={$LL.deleteDepartment()}
    />
  {/if}

  {#if showDeleteOrganization}
    <DeleteConfirmation
      toggleDelete={toggleDeleteOrganization}
      handleDelete={handleDeleteOrganization}
      confirmText={$LL.deleteOrganizationConfirmText()}
      confirmBtnText={$LL.deleteOrganization()}
    />
  {/if}
</AdminPageLayout>
