<script lang="ts">
  import { onMount } from 'svelte';
  import { user } from '../../../stores';
  import LL from '../../../i18n/i18n-svelte';
  import { appRoutes } from '../../../config';
  import { validateUserIsAdmin } from '../../../validationUtils';
  import HeadCol from '../../../components/table/HeadCol.svelte';
  import Table from '../../../components/table/Table.svelte';
  import AdminPageLayout from '../../../components/admin/AdminPageLayout.svelte';
  import TableRow from '../../../components/table/TableRow.svelte';
  import RowCol from '../../../components/table/RowCol.svelte';
  import TableNav from '../../../components/table/TableNav.svelte';
  import TableContainer from '../../../components/table/TableContainer.svelte';
  import TableFooter from '../../../components/table/TableFooter.svelte';

  import type { NotificationService } from '../../../types/notifications'; 
  import type { ApiClient } from '../../../types/apiclient';

  interface Props {
    xfetch: ApiClient;
    router: any;
    notifications: NotificationService;
  }

  let { xfetch, router, notifications }: Props = $props();

  const teamsPageLimit = 100;

  let teamCount = $state(0);
  let teams = $state([]);
  let teamsPage = $state(1);

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
  <TableContainer>
    <TableNav title={$LL.teams()} createBtnEnabled={false} />
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
        </tr>
          {/snippet}
      {#snippet body({ class: className })}
            <tbody   class="{className}">
          {#each teams as team, i}
            <TableRow itemIndex={i}>
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
          {/snippet}
    </Table>
    <TableFooter
      bind:current={teamsPage}
      num_items={teamCount}
      per_page={teamsPageLimit}
      on:navigate={changePage}
    />
  </TableContainer>
</AdminPageLayout>
