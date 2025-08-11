<script lang="ts">
  import { AppConfig, appRoutes } from '../../config';
  import LL from '../../i18n/i18n-svelte';
  import { user } from '../../stores';
  import UserAvatar from '../user/UserAvatar.svelte';
  import { onMount } from 'svelte';
  import { Lock, LogOut, User, Vote } from 'lucide-svelte';
  import SubMenu from './SubMenu.svelte';
  import SubMenuItem from './SubMenuItem.svelte';

  import type { NotificationService } from '../../types/notifications';
  import type { ApiClient } from '../../types/apiclient';

  interface Props {
    currentPage: any;
    notifications: NotificationService;
    router: any;
    xfetch: ApiClient;
  }

  let { currentPage, notifications, router, xfetch }: Props = $props();

  let profile = $state({});

  function goToProfile(toggleSubmenu?: () => void) {
    return () => {
      toggleSubmenu?.();
      router.route(appRoutes.profile, true);
    }
  }

  function goToRegister(toggleSubmenu?: () => void) {
    return () => {
      toggleSubmenu?.();
      router.route(appRoutes.register, true);
    }
  }

  function goToLogin(toggleSubmenu?: () => void) {
    return () => {
      toggleSubmenu?.();
      router.route(appRoutes.login, true);
    }
  }

  function getProfile() {
    xfetch(`/api/users/${$user.id}`)
      .then(res => res.json())
      .then(function (result) {
        profile = result.data;
      })
      .catch(function () {
        notifications.danger($LL.profileErrorRetrieving());
      });
  }

  function logoutUser(toggleSubmenu?: () => void) {
    return () => {
      xfetch('/api/auth/logout', { method: 'DELETE' })
        .then(function () {
          user.delete();
          localStorage.removeItem('theme');
          window.setTheme();
          toggleSubmenu?.();
          router.route(appRoutes.landing, true);
        })
        .catch(function () {
          notifications.danger($LL.logoutError());
        });
    }
  }

  onMount(() => {
    getProfile();
  });
</script>

<SubMenu relativeClass="z-10">
  {#snippet button({ toggleSubmenu })}
    <span
      data-testid="usernav-name"
      class="text-gray-600 dark:text-gray-300 font-semibold me-2"
      >{$user.name}</span
    >
    <button
      class="align-middle rounded-full focus:ring focus:outline-none focus:ring-indigo-600"
      aria-label="Account"
      aria-haspopup="true"
      onclick={toggleSubmenu}
    >
      <UserAvatar
        warriorId={$user.id}
        pictureUrl={profile.picture}
        gravatarHash={profile.gravatarHash}
        class="object-cover w-10 h-10 rounded-full"
        userName={$user.name}
        avatar={profile.avatar}
      />
    </button>
  {/snippet}

  {#snippet children({ toggleSubmenu })}
    <SubMenuItem
      onClickHandler={goToProfile(toggleSubmenu)}
      testId="userprofile-link"
      icon={User}
      label={$LL.profile()}
      active={currentPage === 'profile'}
    />

    {#if $user.rank === 'GUEST' && !AppConfig.OIDCAuthEnabled}
      <SubMenuItem
        onClickHandler={goToRegister(toggleSubmenu)}
        testId="create-account-link"
        icon={Vote}
        label={$LL.createAccount()}
        active={currentPage === 'register'}
      />
      <SubMenuItem
        onClickHandler={goToLogin(toggleSubmenu)}
        testId="login-link"
        icon={Lock}
        label={$LL.login()}
        active={currentPage === 'login'}
      />
    {:else}
      <SubMenuItem
        onClickHandler={logoutUser(toggleSubmenu)}
        testId="logout-link"
        icon={LogOut}
        label={$LL.logout()}
      />
    {/if}
  {/snippet}
</SubMenu>