<script lang="ts">
  import { AppConfig, appRoutes } from '../../config';
  import LL from '../../i18n/i18n-svelte';
  import LoginForm from '../../components/auth/LoginForm.svelte';

  export let router;
  export let xfetch;
  export let notifications;
  export let eventTag;
  export let battleId;
  export let retroId;
  export let storyboardId;
  export let subscription = false;
  export let orgInviteId;
  export let teamInviteId;

  const { AllowRegistration, LdapEnabled } = AppConfig;

  let registerLink = AllowRegistration ? appRoutes.register : '';
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

  let targetPage = appRoutes.games;
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

<div class="flex min-h-[80vh] flex-col justify-center py-12 sm:px-6 lg:px-8">
  <div class="text-center sm:mx-auto sm:w-full sm:max-w-md">
    <h1 class="sr-only">
      {$LL.login()}
    </h1>
  </div>
  <div class="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
    <div
      class="bg-white dark:bg-gray-700 px-4 pb-4 pt-8 sm:rounded-lg sm:px-10 sm:pb-6 sm:shadow"
    >
      <LoginForm
        xfetch="{xfetch}"
        notifications="{notifications}"
        eventTag="{eventTag}"
        router="{router}"
        registerLink="{registerLink}"
        targetPage="{targetPage}"
      />
    </div>
  </div>
</div>
