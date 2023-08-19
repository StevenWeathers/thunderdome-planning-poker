<script lang="ts">
  import { onMount } from 'svelte';

  import PageLayout from '../../components/PageLayout.svelte';
  import { user } from '../../stores';
  import { appRoutes } from '../../config';
  import LL from '../../i18n/i18n-svelte';
  import CreateStoryboard from '../../components/storyboard/CreateStoryboard.svelte';
  import BoxList from '../../components/BoxList.svelte';

  export let xfetch;
  export let notifications;
  export let router;
  export let eventTag;

  let storyboards = [];

  xfetch(`/api/users/${$user.id}/storyboards`)
    .then(res => res.json())
    .then(function (bs) {
      storyboards = bs.data;
    })
    .catch(function (error) {
      notifications.danger($LL.getStoryboardsErrorMessage());
      eventTag('fetch_storyboards', 'engagement', 'failure');
    });

  onMount(() => {
    if (!$user.id) {
      router.route(appRoutes.login);
    }
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
        pageRoute="{appRoutes.storyboard}"
        joinBtnText="{$LL.joinStoryboard()}"
      />
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
