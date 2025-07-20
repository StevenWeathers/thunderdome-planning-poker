<script lang="ts">
  import { AppConfig, appRoutes } from '../../config';
  import LL from '../../i18n/i18n-svelte';
  import { user } from '../../stores';
  import UserAvatar from '../user/UserAvatar.svelte';
  import { onMount } from 'svelte';
  import { Lock, LogOut, User, Vote } from 'lucide-svelte';

  export let currentPage;
  export let notifications;
  export let router;
  export let xfetch;

  let showMenu = false;
  let profile = {};

  function toggleMenu() {
    showMenu = !showMenu;
  }

  function pageChangeHandler() {
    if (showMenu === true) {
      toggleMenu();
    }
  }

  function goToProfile() {
    toggleMenu();
    router.route(appRoutes.profile, true);
  }

  function goToRegister() {
    toggleMenu();
    router.route(appRoutes.register, true);
  }

  function goToLogin() {
    toggleMenu();
    router.route(appRoutes.login, true);
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

  function logoutUser() {
    xfetch('/api/auth/logout', { method: 'DELETE' })
      .then(function () {
        user.delete();
        localStorage.removeItem('theme');
        window.setTheme();
        router.route(appRoutes.landing, true);
      })
      .catch(function () {
        notifications.danger($LL.logoutError());
      });
  }

  $: {
    if (typeof currentPage !== 'undefined') {
      pageChangeHandler();
    }
  }

  onMount(() => {
    getProfile();
  });
</script>

<div>
  <span
    data-testid="usernav-name"
    class="text-gray-600 dark:text-gray-300 font-semibold me-2"
    >{$user.name}</span
  >
  <button
    class="align-middle rounded-full focus:ring focus:outline-none focus:ring-indigo-600"
    aria-label="Account"
    aria-haspopup="true"
    on:click="{toggleMenu}"
  >
    <UserAvatar
      warriorId="{$user.id}"
      pictureUrl="{profile.picture}"
      gravatarHash="{profile.gravatarHash}"
      class="object-cover w-10 h-10 rounded-full"
      userName="{$user.name}"
      avatar="{profile.avatar}"
    />
  </button>
</div>

{#if showMenu}
  <ul
    class="absolute right-0 w-56 p-2 mt-2 space-y-2 text-gray-600 bg-white border border-gray-100 rounded-md shadow-md dark:border-gray-700 dark:text-gray-300 dark:bg-gray-700"
    aria-label="submenu"
  >
    <li class="flex">
      <button
        class="inline-flex items-center w-full px-2 py-1 font-semibold transition-colors duration-150 rounded-md hover:bg-gray-100 hover:text-gray-800 dark:hover:bg-gray-800 dark:hover:text-gray-200"
        data-testid="userprofile-link"
        on:click="{goToProfile}"
      >
        <User class="h-6 w-6 me-3 inline-block" />
        <span>{$LL.profile()}</span>
      </button>
    </li>
    {#if $user.rank === 'GUEST' && !AppConfig.OIDCAuthEnabled}
      <li class="flex">
        <button
          class="inline-flex items-center w-full px-2 py-1 font-semibold transition-colors duration-150 rounded-md hover:bg-gray-100 hover:text-gray-800 dark:hover:bg-gray-800 dark:hover:text-gray-200"
          data-testid="create-account-link"
          on:click="{goToRegister}"
        >
          <Vote class="w-6 h-6 me-3 inline-block" />
          <span>{$LL.createAccount()}</span>
        </button>
      </li>
      <li class="flex">
        <button
          class="inline-flex items-center w-full px-2 py-1 font-semibold transition-colors duration-150 rounded-md hover:bg-gray-100 hover:text-gray-800 dark:hover:bg-gray-800 dark:hover:text-gray-200"
          data-testid="login-link"
          on:click="{goToLogin}"
        >
          <Lock class="w-6 h-6 me-3 inline-block" />
          <span>{$LL.login()}</span>
        </button>
      </li>
    {:else}
      <li class="flex">
        <button
          class="inline-flex items-center w-full px-2 py-1 font-semibold transition-colors duration-150 rounded-md hover:bg-gray-100 hover:text-gray-800 dark:hover:bg-gray-800 dark:hover:text-gray-200"
          data-testid="logout-link"
          on:click="{logoutUser}"
        >
          <LogOut class="inline-block w-6 h-6 me-3" />
          <span>{$LL.logout()}</span>
        </button>
      </li>
    {/if}
  </ul>
{/if}
