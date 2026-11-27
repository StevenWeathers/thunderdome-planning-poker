<script lang="ts">
  import { onMount } from 'svelte';
  import { user } from '../../../stores';
  import LL from '../../../i18n/i18n-svelte';
  import { appRoutes } from '../../../config';
  import { validateUserIsAdmin } from '../../../validationUtils';
  import AdminPageLayout from '../../../components/admin/AdminPageLayout.svelte';
  import RetroTagsList from '../../../components/retrotag/RetroTagsList.svelte';

  import type { NotificationService } from '../../../types/notifications';
  import type { ApiClient } from '../../../types/apiclient';

  interface Props {
    xfetch: ApiClient;
    router: any;
    notifications: NotificationService;
  }

  let { xfetch, router, notifications }: Props = $props();

  const isAdmin = $derived(!!$user.id && validateUserIsAdmin($user));

  onMount(() => {
    if (!$user.id) {
      router.route(appRoutes.login);
      return;
    }
    if (!validateUserIsAdmin($user)) {
      router.route(appRoutes.landing);
    }
  });
</script>

<svelte:head>
  <title>Retro Tags {$LL.admin()} | {$LL.appName()}</title>
</svelte:head>

<AdminPageLayout activePage="retro-tags">
  {#if isAdmin}
    <RetroTagsList {xfetch} {notifications} apiPrefix="/api/admin" />
  {/if}
</AdminPageLayout>
