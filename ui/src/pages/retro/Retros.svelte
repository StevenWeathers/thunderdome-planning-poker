<script lang="ts">
  import { onMount } from 'svelte';

  import PageLayout from '../../components/PageLayout.svelte';
  import CreateRetro from '../../components/retro/CreateRetro.svelte';
  import { user } from '../../stores';
  import { appRoutes } from '../../config';
  import LL from '../../i18n/i18n-svelte';
  import BoxList from '../../components/BoxList.svelte';

  export let xfetch;
  export let notifications;
  export let router;
  export let eventTag;

  let retros = [];

  xfetch(`/api/users/${$user.id}/retros`)
    .then(res => res.json())
    .then(function (bs) {
      retros = bs.data;
    })
    .catch(function () {
      notifications.danger($LL.getRetrosErrorMessage());
      eventTag('fetch_retros', 'engagement', 'failure');
    });

  onMount(() => {
    if (!$user.id) {
      router.route(appRoutes.login);
    }
  });
</script>

<svelte:head>
  <title>{$LL.yourRetros()} | {$LL.appName()}</title>
</svelte:head>

<PageLayout>
  <h1
    class="mb-4 text-4xl font-semibold font-rajdhani uppercase dark:text-white"
  >
    {$LL.myRetros()}
  </h1>

  <div class="flex flex-wrap">
    <div class="mb-4 md:mb-6 w-full md:w-1/2 lg:w-3/5 md:pe-4">
      <BoxList
        items="{retros}"
        itemType="retro"
        pageRoute="{appRoutes.retro}"
        ownerField="ownerId"
        joinBtnText="{$LL.joinRetro()}"
      />
    </div>

    <div class="w-full md:w-1/2 lg:w-2/5 md:ps-2 xl:ps-4">
      <div
        class="p-6 bg-white dark:bg-gray-800 shadow-lg rounded-lg dark:text-white"
      >
        <h2
          class="mb-4 text-3xl font-semibold font-rajdhani uppercase leading-tight"
        >
          {$LL.createARetro()}
        </h2>
        <CreateRetro
          notifications="{notifications}"
          router="{router}"
          eventTag="{eventTag}"
          xfetch="{xfetch}"
        />
      </div>
    </div>
  </div>
</PageLayout>
