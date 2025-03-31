<script lang="ts">
  import { AppConfig, appRoutes } from '../../config';
  import LL from '../../i18n/i18n-svelte';
  import LoginForm from '../../components/auth/LoginForm.svelte';

  interface Props {
    router: any;
    xfetch: any;
    notifications: any;
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
  if (AllowRegistration) {
    if (teamInviteId) {
      registerLink = `${registerLink}/team/${teamInviteId}`;
    }
    if (orgInviteId) {
      registerLink = `${registerLink}/organization/${orgInviteId}`;
    }
    if (battleId) {
      registerLink = `${registerLink}/battle/${battleId}`;
    }
    if (retroId) {
      registerLink = `${registerLink}/retro/${retroId}`;
    }
    if (storyboardId) {
      registerLink = `${registerLink}/storyboard/${storyboardId}`;
    }
    if (subscription) {
      registerLink = `${registerLink}/subscription`;
    }
  }

  let targetPage = $state(appRoutes.games);
  if (teamInviteId) {
    targetPage = `${appRoutes.invite}/team/${teamInviteId}`;
  }
  if (orgInviteId) {
    targetPage = `${appRoutes.invite}/organization/${orgInviteId}`;
  }
  if (subscription) {
    targetPage = `${appRoutes.subscriptionPricing}`;
  }
  if (battleId) {
    targetPage = `${appRoutes.game}/${battleId}`;
  }
  if (retroId) {
    targetPage = `${appRoutes.retro}/${retroId}`;
  }
  if (storyboardId) {
    targetPage = `${appRoutes.storyboard}/${storyboardId}`;
  }
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
      router="{router}"
      registerLink="{registerLink}"
      targetPage="{targetPage}"
    />
  </div>
</div>
