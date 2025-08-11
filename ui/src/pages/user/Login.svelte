<script lang="ts">
  import { AppConfig, appRoutes } from '../../config';
  import LL from '../../i18n/i18n-svelte';
  import LoginForm from '../../components/auth/LoginForm.svelte';

  import type { NotificationService } from '../../types/notifications';

  interface Props {
    router: any;
    xfetch: any;
    notifications: NotificationService;
    battleId: any;
    retroId: any;
    storyboardId: any;
    subscription?: boolean;
    orgInviteId: any;
    teamInviteId: any;
  }

  let {
    router,
    xfetch,
    notifications,
    battleId,
    retroId,
    storyboardId,
    subscription = false,
    orgInviteId,
    teamInviteId
  }: Props = $props();

  const { AllowRegistration } = AppConfig;

  let registerLink = $state(AllowRegistration ? appRoutes.register : '');
  $effect(() => {
    if (!AllowRegistration) {
      registerLink = '';
      return;
    }
    let base = appRoutes.register;
    if (teamInviteId) {
      base = `${base}/team/${teamInviteId}`;
    }
    if (orgInviteId) {
      base = `${base}/organization/${orgInviteId}`;
    }
    if (battleId) {
      base = `${base}/battle/${battleId}`;
    }
    if (retroId) {
      base = `${base}/retro/${retroId}`;
    }
    if (storyboardId) {
      base = `${base}/storyboard/${storyboardId}`;
    }
    if (subscription) {
      base = `${base}/subscription`;
    }
    registerLink = base;
  });

  let targetPage = $state(appRoutes.games);
  $effect(() => {
    let next = appRoutes.games;
    if (teamInviteId) {
      next = `${appRoutes.invite}/team/${teamInviteId}`;
    }
    if (orgInviteId) {
      next = `${appRoutes.invite}/organization/${orgInviteId}`;
    }
    if (subscription) {
      next = `${appRoutes.subscriptionPricing}`;
    }
    if (battleId) {
      next = `${appRoutes.game}/${battleId}`;
    }
    if (retroId) {
      next = `${appRoutes.retro}/${retroId}`;
    }
    if (storyboardId) {
      next = `${appRoutes.storyboard}/${storyboardId}`;
    }
    targetPage = next;
  });
</script>

<svelte:head>
  <title>{$LL.login()} | {$LL.appName()}</title>
</svelte:head>

<h1 class="sr-only">
  {$LL.login()}
</h1>
<div
  class="space-y-8 flex items-center justify-center min-h-[80vh] bg-gradient-to-br from-blue-200 via-purple-200 to-indigo-300 dark:from-blue-800 dark:via-purple-800 dark:to-indigo-800"
>
  <div
    class="p-8 rounded-2xl shadow-2xl backdrop-blur-sm bg-white/70 dark:bg-gray-800/50"
  >
    <LoginForm
      xfetch={xfetch}
      notifications={notifications}
      router={router}
      registerLink={registerLink}
      targetPage={targetPage}
    />
  </div>
</div>
