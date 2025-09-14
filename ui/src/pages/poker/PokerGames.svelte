<script lang="ts">
  import { onMount } from 'svelte';

  import PageLayout from '../../components/PageLayout.svelte';
  import CreateBattle from '../../components/poker/CreatePokerGame.svelte';
  import { user } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import { appRoutes } from '../../config';
  import Pagination from '../../components/global/Pagination.svelte';
  import BoxList from '../../components/BoxList.svelte';

  import type { NotificationService } from '../../types/notifications';
  import type { ApiClient } from '../../types/apiclient';

  interface Props {
    xfetch: ApiClient;
    notifications: NotificationService;
    router: any;
  }

  let { xfetch, notifications, router }: Props = $props();

  const battlesPageLimit = 10;
  let battleCount = $state(0);
  let battlesPage = $state(1);
  let battles = $state([]);
  let loading = $state(true);

  // Stop game functionality
  function stopGame(gameId: string) {
    return () => {
      xfetch(`/api/games/${gameId}/stop`, { method: 'POST' })
        .then(res => {
          if (res.ok) {
            notifications.success($LL.battleStopped());
            // Update the local state to reflect the stopped game
            battles = battles.map(battle => 
              battle.id === gameId 
                ? { ...battle, endedDate: new Date() }
                : battle
            );
          } else {
            throw new Error('Failed to stop game');
          }
        })
        .catch(err => {
          console.error('Error stopping game:', err);
          notifications.danger($LL.stopGameError());
        });
    };
  }

  function getBattles() {
    const battlesOffset = (battlesPage - 1) * battlesPageLimit;

    xfetch(
      `/api/users/${$user.id}/battles?limit=${battlesPageLimit}&offset=${battlesOffset}`,
    )
      .then(res => res.json())
      .then(function (result) {
        // Convert endedDate strings to Date objects for consistency
        battles = result.data.map(battle => ({
          ...battle,
          endedDate: battle.endedDate ? new Date(battle.endedDate) : battle.endedDate
        }));
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
          items={battles}
          itemType="battle"
          pageRoute={appRoutes.game}
          joinBtnText={$LL.battleJoin()}
          showOwner={false}
          showOwnerName={true}
          ownerNameField="teamName"
          showFacilitatorIcon={true}
          facilitatorsKey="leaders"
          showCompletedStories={true}
          showStopButton={true}
          showStatusBadge={true}
          toggleStop={stopGame}
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
            bind:current={battlesPage}
            num_items={battleCount}
            per_page={battlesPageLimit}
            on:navigate={changePage}
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
          notifications={notifications}
          router={router}
          xfetch={xfetch}
        />
      </div>
    </div>
  </div>

  <div class="w-full text-gray-600 dark:text-gray-300">
    <p class="py-8 md:text-lg italic">{$LL.pokerDescription()}</p>
  </div>
</PageLayout>
