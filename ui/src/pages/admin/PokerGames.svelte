<script lang="ts">
  import { onMount } from 'svelte';
  import { user } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import { AppConfig, appRoutes } from '../../config';
  import { validateUserIsAdmin } from '../../validationUtils';
  import HeadCol from '../../components/table/HeadCol.svelte';
  import AdminPageLayout from '../../components/AdminPageLayout.svelte';
  import TableRow from '../../components/table/TableRow.svelte';
  import RowCol from '../../components/table/RowCol.svelte';
  import HollowButton from '../../components/global/HollowButton.svelte';
  import Table from '../../components/table/Table.svelte';
  import TableContainer from '../../components/table/TableContainer.svelte';
  import TableNav from '../../components/table/TableNav.svelte';
  import TableFooter from '../../components/table/TableFooter.svelte';

  export let xfetch;
  export let router;
  export let notifications;
  // export let eventTag

  const battlesPageLimit = 100;
  let battleCount = 0;
  let battles = [];
  let battlesPage = 1;
  let activeBattles = false;

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
        notifications.danger(
          $LL.getBattlesError({
            friendly: AppConfig.FriendlyUIVerbs,
          }),
        );
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
    >{$LL.battles({ friendly: AppConfig.FriendlyUIVerbs })}
    {$LL.admin()} | {$LL.appName()}</title
  >
</svelte:head>

<AdminPageLayout activePage="battles">
  <div class="text-center px-2 mb-4">
    <h1
      class="text-3xl md:text-4xl font-semibold font-rajdhani uppercase dark:text-white"
    ></h1>
  </div>

  <TableContainer>
    <TableNav
      title="{$LL.battles({ friendly: AppConfig.FriendlyUIVerbs })}"
      createBtnEnabled="{false}"
    >
      <label class="inline-flex items-center cursor-pointer">
        <input
          type="checkbox"
          class="sr-only peer"
          name="activeBattles"
          id="activeBattles"
          bind:checked="{activeBattles}"
          on:change="{changeActiveBattlesToggle}"
        />
        <div
          class="relative w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 dark:peer-focus:ring-blue-800 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-blue-600"
        ></div>
        <span class="ms-3 text-sm font-medium text-gray-900 dark:text-gray-300">
          {$LL.showActiveBattles({
            friendly: AppConfig.FriendlyUIVerbs,
          })}
        </span>
      </label>
    </TableNav>
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
        {#each battles as battle, i}
          <TableRow itemIndex="{i}">
            <RowCol>
              <a
                href="{appRoutes.admin}/battles/{battle.id}"
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
                {$LL.battleJoin({
                  friendly: AppConfig.FriendlyUIVerbs,
                })}
              </HollowButton>
            </RowCol>
          </TableRow>
        {/each}
      </tbody>
    </Table>
    <TableFooter
      bind:current="{battlesPage}"
      num_items="{battleCount}"
      per_page="{battlesPageLimit}"
      on:navigate="{changePage}"
    />
  </TableContainer>
</AdminPageLayout>
