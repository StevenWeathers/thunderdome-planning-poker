<script lang="ts">
  import { onMount } from 'svelte';

  import AdminPageLayout from '../../../components/admin/AdminPageLayout.svelte';
  import DeleteConfirmation from '../../../components/global/DeleteConfirmation.svelte';
  import { ChevronRight } from 'lucide-svelte';
  import CountryFlag from '../../../components/user/CountryFlag.svelte';
  import UserAvatar from '../../../components/user/UserAvatar.svelte';
  import { user } from '../../../stores';
  import LL from '../../../i18n/i18n-svelte';
  import { appRoutes } from '../../../config';
  import RowCol from '../../../components/table/RowCol.svelte';
  import TableRow from '../../../components/table/TableRow.svelte';
  import HeadCol from '../../../components/table/HeadCol.svelte';
  import { validateUserIsAdmin } from '../../../validationUtils';
  import Table from '../../../components/table/Table.svelte';
  import TableNav from '../../../components/table/TableNav.svelte';
  import TableContainer from '../../../components/table/TableContainer.svelte';
  import CrudActions from '../../../components/table/CrudActions.svelte';

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

  let organization = $state({
    id: organizationId,
    name: '',
  });
  let department = $state({
    id: departmentId,
    name: '',
  });
  let departmentRole = '';
  let organizationRole = '';
  let teams = $state([]);
  let users = $state([]);
  let showDeleteTeam = $state(false);
  let deleteTeamId = null;
  let teamsPage = 1;
  let usersPage = 1;

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

        getTeams();
        getUsers();
      })
      .catch(function () {
        notifications.danger($LL.departmentGetError());
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
        notifications.danger($LL.departmentUsersGetError());
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

  onMount(() => {
    if (!$user.id) {
      router.route(appRoutes.login);
      return;
    }
    if (!validateUserIsAdmin($user)) {
      router.route(appRoutes.landing);
      return;
    }

    getDepartment();
  });
</script>

<svelte:head>
  <title>{$LL.department()} {department.name} | {$LL.appName()}</title>
</svelte:head>

<AdminPageLayout activePage="organizations">
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
        href="{appRoutes.adminOrganizations}/{organization.id}"
      >
        {organization.name}
      </a>
    </div>
  </div>

  <div class="w-full mb-6 lg:mb-8">
    <TableContainer>
      <TableNav title="{$LL.teams()}" createBtnEnabled="{false}" />
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
              <TableRow itemIndex="{i}">
                <RowCol>
                  <a
                    href="{appRoutes.adminOrganizations}/{organizationId}/department/{departmentId}/team/{team.id}"
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
                  <CrudActions
                    editBtnEnabled="{false}"
                    deleteBtnClickHandler="{toggleDeleteTeam(team.id)}"
                  />
                </RowCol>
              </TableRow>
            {/each}
          </tbody>
              {/snippet}
      </Table>
    </TableContainer>
  </div>

  <TableContainer>
    <TableNav title="{$LL.users()}" createBtnEnabled="{false}" />
    <Table>
      {#snippet header()}
            <tr >
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
            <tbody   class="{className}">
          {#each users as user, i}
            <TableRow itemIndex="{i}">
              <RowCol>
                <div class="flex items-center">
                  <div class="flex-shrink-0 h-10 w-10">
                    <UserAvatar
                      warriorId="{user.id}"
                      avatar="{user.avatar}"
                      gravatarHash="{user.gravatarHash}"
                      userName="{user.name}"
                      width="48"
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
            </TableRow>
          {/each}
        </tbody>
          {/snippet}
    </Table>
  </TableContainer>

  {#if showDeleteTeam}
    <DeleteConfirmation
      toggleDelete="{toggleDeleteTeam(null)}"
      handleDelete="{handleDeleteTeam}"
      confirmText="{$LL.deleteTeamConfirmText()}"
      confirmBtnText="{$LL.deleteTeam()}"
    />
  {/if}
</AdminPageLayout>
