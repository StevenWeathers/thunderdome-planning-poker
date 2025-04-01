<script lang="ts">
  import { onMount } from 'svelte';
  import { user } from '../../../stores';
  import LL from '../../../i18n/i18n-svelte';
  import { appRoutes } from '../../../config';
  import { validateUserIsAdmin } from '../../../validationUtils';
  import HeadCol from '../../../components/table/HeadCol.svelte';
  import AdminPageLayout from '../../../components/admin/AdminPageLayout.svelte';
  import TableRow from '../../../components/table/TableRow.svelte';
  import RowCol from '../../../components/table/RowCol.svelte';
  import HollowButton from '../../../components/global/HollowButton.svelte';
  import Table from '../../../components/table/Table.svelte';
  import TableContainer from '../../../components/table/TableContainer.svelte';
  import TableNav from '../../../components/table/TableNav.svelte';
  import TableFooter from '../../../components/table/TableFooter.svelte';
  import Toggle from '../../../components/forms/Toggle.svelte';

  interface Props {
    xfetch: any;
    router: any;
    notifications: any;
  }

  let { xfetch, router, notifications }: Props = $props();

  const battlesPageLimit = 100;
  let battleCount = $state(0);
  let battles = $state([]);
  let battlesPage = $state(1);
  let activeBattles = $state(false);

  function getBattles() {
    const battlesOffset = (battlesPage - 1) * battlesPageLimit;
    xfetch(
      `/api/battles?limit=${battlesPageLimit}&offset=${battlesOffset}&active=${activeBattles}`,
    )
      .then(res => res.json())
      .then(function (result) {
        battles = result.data;
        battleCount = result.meta.count;
      })
      .catch(function () {
        notifications.danger($LL.getBattlesError());
      });
  }

  const changePage = evt => {
    battlesPage = evt.detail;
    getBattles();
  };

  const changeActiveBattlesToggle = () => {
    battlesPage = 1;
    getBattles();
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

    getBattles();
  });
</script>

<svelte:head>
  <title
    >{$LL.battles()}
    {$LL.admin()} | {$LL.appName()}</title
  >
</svelte:head>

<AdminPageLayout activePage="battles">
  <TableContainer>
    <TableNav title={$LL.battles()} createBtnEnabled={false}>
      <Toggle
        name="activeBattles"
        id="activeBattles"
        bind:checked={activeBattles}
        changeHandler={changeActiveBattlesToggle}
        label={$LL.showActiveBattles()}
      />
    </TableNav>
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
          {#each battles as battle, i}
            <TableRow itemIndex={i}>
              <RowCol>
                <a
                  href="{appRoutes.adminPokerGames}/{battle.id}"
                  class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
                  >{battle.name}</a
                >
              </RowCol>
              <RowCol>
                {new Date(battle.createdDate).toLocaleString()}
              </RowCol>
              <RowCol>
                {new Date(battle.updatedDate).toLocaleString()}
              </RowCol>
              <RowCol type="action">
                <HollowButton href="{appRoutes.game}/{battle.id}">
                  {$LL.battleJoin()}
                </HollowButton>
              </RowCol>
            </TableRow>
          {/each}
        </tbody>
          {/snippet}
    </Table>
    <TableFooter
      bind:current={battlesPage}
      num_items={battleCount}
      per_page={battlesPageLimit}
      on:navigate={changePage}
    />
  </TableContainer>
</AdminPageLayout>
