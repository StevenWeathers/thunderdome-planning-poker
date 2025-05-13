<script lang="ts">
    import { onMount } from 'svelte';

  import PageLayout from '../../components/PageLayout.svelte';
  import LL from '../../i18n/i18n-svelte';
  import { user } from '../../stores';
  import { appRoutes } from '../../config';

  interface Props {
    xfetch: any;
    verifyId: any;
    notifications: any;
    router: any;
    changeId: string;
  }

  let { xfetch, notifications, changeId, router }: Props = $props();

    onMount(() => {
    if (!$user.id) {
      router.route(appRoutes.login);
      return;
    }
  });
</script>

<svelte:head>
  <title>{$LL.verifyAccount()} | {$LL.appName()}</title>
</svelte:head>

<PageLayout>
  <div class="flex justify-center">
    <div class="w-full md:w-1/2 xl:w-1/3 py-4">
      Change Email Form here...
    </div>
  </div>
</PageLayout>
