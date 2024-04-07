<script lang="ts">
  import { onMount } from 'svelte';

  import PageLayout from '../../components/global/PageLayout.svelte';
  import { user } from '../../stores';
  import { appRoutes } from '../../config';
  import LL from '../../i18n/i18n-svelte';
  import CreateStoryboard from '../../components/storyboard/CreateStoryboard.svelte';
  import BoxList from '../../components/BoxList.svelte';
  import Pagination from '../../components/global/Pagination.svelte';

  export let xfetch;
  export let notifications;
  export let router;
  export let eventTag;

  let storyboards = [];
  const storyboardsPageLimit = 10;
  let storyboardCount = 0;
  let storyboardsPage = 1;

  function getStoryboards() {
    const retrosOffset = (storyboardsPage - 1) * storyboardsPageLimit;

    xfetch(
      `/api/users/${$user.id}/storyboards?limit=${storyboardsPageLimit}&offset=${retrosOffset}`,
    )
      .then(res => res.json())
      .then(function (result) {
        storyboards = result.data;
        storyboardCount = result.meta.count;
      })
      .catch(function (error) {
        notifications.danger(`${$LL.getStoryboardsErrorMessage()}`);
        eventTag('fetch_storyboards', 'engagement', 'failure');
      });
  }

  const changePage = evt => {
    storyboardsPage = evt.detail;
    getStoryboards();
  };

  onMount(() => {
    if (!$user.id) {
      router.route(appRoutes.login);
    }
    getStoryboards();
  });
</script>

<svelte:head>
  <title>{$LL.yourStoryboards()} | {$LL.appName()}</title>
</svelte:head>

<PageLayout>
  <h1
    class="mb-4 text-4xl font-semibold font-rajdhani uppercase dark:text-white"
  >
    {$LL.myStoryboards()}
  </h1>

  <div class="flex flex-wrap">
    <div class="mb-4 md:mb-6 w-full md:w-1/2 lg:w-3/5 md:pe-4">
      <BoxList
        items="{storyboards}"
        itemType="storyboard"
        showOwnerName="{true}"
        ownerNameField="teamName"
        pageRoute="{appRoutes.storyboard}"
        joinBtnText="{$LL.joinStoryboard()}"
      />
      {#if storyboardCount > storyboardsPageLimit}
        <div class="mt-6 pt-1 flex justify-center">
          <Pagination
            bind:current="{storyboardsPage}"
            num_items="{storyboardCount}"
            per_page="{storyboardsPageLimit}"
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
          {$LL.createAStoryboard()}
        </h2>
        <CreateStoryboard
          notifications="{notifications}"
          router="{router}"
          eventTag="{eventTag}"
          xfetch="{xfetch}"
        />
      </div>
    </div>
  </div>
</PageLayout>
