<script lang="ts">
  import { onMount } from 'svelte';

  import PageLayout from '../../components/PageLayout.svelte';
  import CreateBattle from '../../components/poker/CreatePokerGame.svelte';
  import { user } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import { appRoutes } from '../../config';
  import Pagination from '../../components/global/Pagination.svelte';
  import BoxList from '../../components/BoxList.svelte';

  export let xfetch;
  export let notifications;
  export let router;

  const battlesPageLimit = 10;
  let battleCount = 0;
  let battlesPage = 1;
  let battles = [];
  let loading = true;

  function getBattles() {
    const battlesOffset = (battlesPage - 1) * battlesPageLimit;

    xfetch(
      `/api/users/${$user.id}/battles?limit=${battlesPageLimit}&offset=${battlesOffset}`,
    )
      .then(res => res.json())
      .then(function (result) {
        battles = result.data;
        battleCount = result.meta.count;
        loading = false;
      })
      .catch(function () {
        loading = false;
        notifications.danger($LL.myBattlesError());
      });
  }

  const changePage = evt => {
    battlesPage = evt.detail;
    getBattles();
  };

  onMount(() => {
    if (!$user.id) {
      router.route(appRoutes.login);
      return;
    }
    getBattles();
  });
</script>

<svelte:head>
  <title>{$LL.myBattles()} | {$LL.appName()}</title>
</svelte:head>

<PageLayout>
  <h1
    class="mb-4 text-4xl font-semibold font-rajdhani uppercase dark:text-white"
  >
    {$LL.myBattles()}
  </h1>

  <div class="flex flex-wrap">
    <div class="mb-4 md:mb-6 w-full md:w-1/2 lg:w-3/5 md:pe-4">
      {#if battleCount > 0}
        <BoxList
          items="{battles}"
          itemType="battle"
          pageRoute="{appRoutes.game}"
          joinBtnText="{$LL.battleJoin()}"
          showOwner="{false}"
          showOwnerName="{true}"
          ownerNameField="teamName"
          showFacilitatorIcon="{true}"
          facilitatorsKey="leaders"
          showCompletedStories="{true}"
        />
      {:else if loading === false}
        <div
          class="w-full my-10 text-lg md:text-xl dark:text-white text-center"
        >
          {$LL.noGamesFound()}
        </div>
      {/if}
      {#if battleCount > battlesPageLimit}
        <div class="mt-6 pt-1 flex justify-center">
          <Pagination
            bind:current="{battlesPage}"
            num_items="{battleCount}"
            per_page="{battlesPageLimit}"
            on:navigate="{changePage}"
          />
        </div>
      {/if}
    </div>

    <div class="w-full md:w-1/2 lg:w-2/5 md:ps-2 xl:ps-4">
      <div
        class="p-6 bg-white dark:bg-gray-800 shadow-lg rounded-lg dark:text-white"
      >
        <h2
          class="mb-4 text-3xl font-semibold font-rajdhani uppercase leading-tight"
        >
          {$LL.createBattle()}
        </h2>
        <CreateBattle
          notifications="{notifications}"
          router="{router}"
          xfetch="{xfetch}"
        />
      </div>
    </div>
  </div>

  <div class="w-full text-gray-600 dark:text-gray-300">
    <p class="py-8 md:text-lg italic">{$LL.pokerDescription()}</p>
  </div>
</PageLayout>
