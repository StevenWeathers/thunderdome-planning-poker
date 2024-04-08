<script lang="ts">
  import { onMount } from 'svelte';
  import { user } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import { appRoutes } from '../../config';
  import { validateUserIsAdmin } from '../../validationUtils';
  import AdminPageLayout from '../../components/AdminPageLayout.svelte';
  import Table from '../../components/table/Table.svelte';
  import HeadCol from '../../components/table/HeadCol.svelte';
  import TableRow from '../../components/table/TableRow.svelte';
  import RowCol from '../../components/table/RowCol.svelte';
  import TableNav from '../../components/table/TableNav.svelte';
  import TableContainer from '../../components/table/TableContainer.svelte';
  import TableFooter from '../../components/table/TableFooter.svelte';

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
        notifications.danger($LL.applicationStatsError());
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
        notifications.danger($LL.getOrganizationsError());
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
  <TableContainer>
    <TableNav title="{$LL.organizations()}" createBtnEnabled="{false}" />
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
    <TableFooter
      bind:current="{organizationsPage}"
      num_items="{appStats.organizationCount}"
      per_page="{organizationsPageLimit}"
      on:navigate="{changePage}"
    />
  </TableContainer>
</AdminPageLayout>
