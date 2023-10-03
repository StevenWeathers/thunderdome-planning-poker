<script lang="ts">
  import { onMount } from 'svelte';
  import { user } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import { appRoutes } from '../../config';
  import { validateUserIsAdmin } from '../../validationUtils';
  import HeadCol from '../../components/global/table/HeadCol.svelte';
  import Table from '../../components/global/table/Table.svelte';
  import AdminPageLayout from '../../components/global/AdminPageLayout.svelte';
  import TableRow from '../../components/global/table/TableRow.svelte';
  import RowCol from '../../components/global/table/RowCol.svelte';
  import Pagination from '../../components/global/Pagination.svelte';

  export let xfetch;
  export let router;
  export let notifications;

  const teamsPageLimit = 100;

  let teamCount = 0;
  let teams = [];
  let teamsPage = 1;

  function getTeams() {
    const teamsOffset = (teamsPage - 1) * teamsPageLimit;
    xfetch(`/api/admin/teams?limit=${teamsPageLimit}&offset=${teamsOffset}`)
      .then(res => res.json())
      .then(function (result) {
        teams = result.data;
        teamCount = result.meta.count;
      })
      .catch(function () {
        notifications.danger($LL.getTeamsError());
      });
  }

  const changePage = evt => {
    teamsPage = evt.detail;
    getTeams();
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

    getTeams();
  });
</script>

<svelte:head>
  <title>{$LL.teams()} {$LL.admin()} | {$LL.appName()}</title>
</svelte:head>

<AdminPageLayout activePage="teams">
  <div class="text-center px-2 mb-4">
    <h1
      class="text-3xl md:text-4xl font-semibold font-rajdhani uppercase dark:text-white"
    >
      {$LL.teams()}
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
        {#each teams as team, i}
          <TableRow itemIndex="{i}">
            <RowCol>
              <a
                href="{appRoutes.adminTeams}/{team.id}"
                class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
                >{team.name}</a
              >
            </RowCol>
            <RowCol>
              {new Date(team.createdDate).toLocaleString()}
            </RowCol>
            <RowCol>
              {new Date(team.updatedDate).toLocaleString()}
            </RowCol>
          </TableRow>
        {/each}
      </tbody>
    </Table>

    {#if teamCount > teamsPageLimit}
      <div class="pt-6 flex justify-center">
        <Pagination
          bind:current="{teamsPage}"
          num_items="{teamCount}"
          per_page="{teamsPageLimit}"
          on:navigate="{changePage}"
        />
      </div>
    {/if}
  </div>
</AdminPageLayout>
