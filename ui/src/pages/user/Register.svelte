<script lang="ts">
  import { user } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import { appRoutes } from '../../config';
  import UserRegisterForm from '../../components/user/UserRegisterForm.svelte';
  import { onMount } from 'svelte';

  interface Props {
    router: any;
    xfetch: any;
    notifications: any;
    battleId: any;
    retroId: any;
    storyboardId: any;
    orgInviteId: any;
    teamInviteId: any;
    subscription?: boolean;
  }

  let {
    router,
    xfetch,
    notifications,
    battleId,
    retroId,
    storyboardId,
    orgInviteId,
    teamInviteId,
    subscription = false
  }: Props = $props();

  let userName = $user.name || '';
  let wasInvited = $state(false);
  let inviteDetails = $state({
    email: '',
  });

  function targetPage() {
    let tp = appRoutes.games;

    if (teamInviteId) {
      tp = `${appRoutes.invite}/team/${teamInviteId}`;
    }
    if (orgInviteId) {
      tp = `${appRoutes.invite}/organization/${orgInviteId}`;
    }
    if (subscription) {
      tp = `${appRoutes.subscriptionPricing}`;
    }

    if (battleId) {
      tp = `${appRoutes.game}/${battleId}`;
    }

    if (retroId) {
      tp = `${appRoutes.retro}/${retroId}`;
    }

    if (storyboardId) {
      tp = `${appRoutes.storyboard}/${storyboardId}`;
    }

    return tp;
  }

  function createUserGuest(name) {
    const body = {
      name,
    };

    xfetch('/api/auth/guest', { body })
      .then(res => res.json())
      .then(function (result) {
        const newWarrior = result.data;
        user.create({
          id: newWarrior.id,
          name: newWarrior.name,
          rank: newWarrior.rank,
          notificationsEnabled: newWarrior.notificationsEnabled,
        });

        router.route(targetPage(), true);
      })
      .catch(function () {
        notifications.danger($LL.guestRegisterError());
      });
  }

  function createUserRegistered(name, email, password1, password2) {
    const body = {
      name,
      email,
      password1,
      password2,
    };

    xfetch('/api/auth/register', { body })
      .then(res => res.json())
      .then(function (result) {
        const newWarrior = result.data;
        user.create({
          id: newWarrior.id,
          name: newWarrior.name,
          email: newWarrior.email,
          rank: newWarrior.rank,
          notificationsEnabled: newWarrior.notificationsEnabled,
          subscribed: false,
        });

        router.route(targetPage(), true);
      })
      .catch(function () {
        notifications.danger($LL.registerError());
      });
  }

  function getInviteDetails() {
    const inviteType =
      typeof teamInviteId !== 'undefined' ? 'team' : 'organization';
    const inviteId = inviteType === 'team' ? teamInviteId : orgInviteId;
    xfetch(`/api/auth/invite/${inviteType}/${inviteId}`)
      .then(res => res.json())
      .then(function (result) {
        inviteDetails = result.data;
      })
      .catch(function () {
        notifications.danger(`Failed to get ${inviteType} invite`);
      });
  }

  let registerDisabled = $derived(wasInvited && inviteDetails.email === '');

  const loginLinkContextClasses =
    'font-rajdhani uppercase text-lg md:text-xl lg:text-2xl md:leading-tight';
  const loginLinkClasses =
    'font-bold underline text-blue-600 dark:text-cyan-300 hover:text-purple-700 dark:hover:text-yellow-thunder transition-colors duration-300';

  onMount(() => {
    if (orgInviteId || teamInviteId) {
      wasInvited = true;
      getInviteDetails();
    }
  });
</script>

<svelte:head>
  <title>{$LL.register()} | {$LL.appName()}</title>
</svelte:head>

<div
  class="min-h-[80vh] py-10 space-y-8 bg-gradient-to-br from-blue-200 via-purple-200 to-indigo-300 dark:from-blue-800 dark:via-purple-800 dark:to-indigo-800"
>
  <div class="w-full text-center text-gray-700 dark:text-white">
    <h1
      class="font-rajdhani leading-tight uppercase text-4xl md:text-5xl font-semibold dark:drop-shadow-[0_2px_2px_rgba(0,0,0,0.8)]"
    >
      {$LL.register()}
    </h1>
    {#if teamInviteId != null}
      <div class="{loginLinkContextClasses}">to join your Team</div>
    {/if}
    {#if orgInviteId != null}
      <div class="{loginLinkContextClasses}">to join your Organization</div>
    {/if}
    {#if battleId}
      <div class="{loginLinkContextClasses}">
        {@html $LL.loginForBattle({
          loginOpen: `<a href="${appRoutes.login}/battle/${battleId}" class="${loginLinkClasses}">`,
          loginClose: `</a>`,
        })}
      </div>
    {/if}
    {#if retroId}
      <div class="{loginLinkContextClasses}">
        {@html $LL.loginForRetro({
          loginOpen: `<a href="${appRoutes.login}/retro/${retroId}" class="${loginLinkClasses}">`,
          loginClose: `</a>`,
        })}
      </div>
    {/if}
    {#if storyboardId}
      <div class="{loginLinkContextClasses}">
        {@html $LL.loginForStoryboard({
          loginOpen: `<a href="${appRoutes.login}/storyboard/${storyboardId}" class="${loginLinkClasses}">`,
          loginClose: `</a>`,
        })}
      </div>
    {/if}
  </div>

  <div class="w-full flex flex-wrap justify-center">
    <div
      class="w-full md:max-w-lg p-8 space-y-8 rounded-2xl shadow-2xl backdrop-blur-sm bg-white/70 dark:bg-gray-800/50"
    >
      <UserRegisterForm
        userName={userName}
        handleFullAccountRegistration={createUserRegistered}
        handleGuestRegistration={createUserGuest}
        notifications={notifications}
        email={wasInvited ? inviteDetails.email : ''}
        wasInvited={wasInvited}
      />
    </div>
  </div>
</div>
