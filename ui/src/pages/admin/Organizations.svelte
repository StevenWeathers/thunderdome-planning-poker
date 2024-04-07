<script lang="ts">
  import { onMount } from 'svelte';
  import { user } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import { appRoutes } from '../../config';
  import { validateUserIsAdmin } from '../../validationUtils';
  import AdminPageLayout from '../../components/global/AdminPageLayout.svelte';
  import Table from '../../components/global/table/Table.svelte';
  import HeadCol from '../../components/global/table/HeadCol.svelte';
  import TableRow from '../../components/global/table/TableRow.svelte';
  import RowCol from '../../components/global/table/RowCol.svelte';
  import Pagination from '../../components/global/Pagination.svelte';

  export let xfetch;
  export let router;
  export let notifications;
  // export let eventTag

  const organizationsPageLimit = 100;

  let appStats = {
    unregisteredUserCount: 0,
    registeredUserCount: 0,
    battleCount: 0,
    planCount: 0,
    organizationCount: 0,
    departmentCount: 0,
    teamCount: 0,
  };
  let organizations = [];
  let organizationsPage = 1;

  function getAppStats() {
    xfetch('/api/admin/stats')
      .then(res => res.json())
      .then(function (result) {
        appStats = result.data;
      })
      .catch(function () {
        notifications.danger(`${$LL.applicationStatsError()}`);
      });
  }

  function getOrganizations() {
    const organizationsOffset =
      (organizationsPage - 1) * organizationsPageLimit;
    xfetch(
      `/api/admin/organizations?limit=${organizationsPageLimit}&offset=${organizationsOffset}`,
    )
      .then(res => res.json())
      .then(function (result) {
        organizations = result.data;
      })
      .catch(function () {
        notifications.danger(`${$LL.getOrganizationsError()}`);
      });
  }

  const changePage = evt => {
    organizationsPage = evt.detail;
    getOrganizations();
  };

  onMount(() => {
    if (!$user.id) {
      router.route(appRoutes.login);
      return;
    }
    if (!validateUserIsAdmin($user)) {
      router.route(appRoutes.landing);
      return;
    }

    getAppStats();
    getOrganizations();
  });
</script>

<svelte:head>
  <title>
    {$LL.organizations()}
    {$LL.admin()} | {$LL.appName()}
  </title>
</svelte:head>

<AdminPageLayout activePage="organizations">
  <div class="text-center px-2 mb-4">
    <h1
      class="text-3xl md:text-4xl font-semibold font-rajdhani uppercase dark:text-white"
    >
      {$LL.organizations()}
    </h1>
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
      </tr>
      <tbody slot="body" let:class="{className}" class="{className}">
        {#each organizations as org, i}
          <TableRow itemIndex="{i}">
            <RowCol>
              <a
                href="{appRoutes.adminOrganizations}/{org.id}"
                class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
                >{org.name}</a
              >
            </RowCol>
            <RowCol>
              {new Date(org.createdDate).toLocaleString()}
            </RowCol>
            <RowCol>
              {new Date(org.updatedDate).toLocaleString()}
            </RowCol>
          </TableRow>
        {/each}
      </tbody>
    </Table>

    {#if appStats.organizationCount > organizationsPageLimit}
      <div class="pt-6 flex justify-center">
        <Pagination
          bind:current="{organizationsPage}"
          num_items="{appStats.organizationCount}"
          per_page="{organizationsPageLimit}"
          on:navigate="{changePage}"
        />
      </div>
    {/if}
  </div>
</AdminPageLayout>
