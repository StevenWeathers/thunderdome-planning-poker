<script lang="ts">
  import UserIcon from '../icons/UserIcon.svelte';
  import { AppConfig, appRoutes } from '../../config';
  import LL from '../../i18n/i18n-svelte';
  import { user } from '../../stores';
  import UserAvatar from '../user/UserAvatar.svelte';
  import { onMount } from 'svelte';
  import VoteIcon from '../icons/VoteIcon.svelte';
  import LockIcon from '../icons/LockIcon.svelte';

  export let currentPage;
  export let eventTag;
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
        eventTag('fetch_profile', 'engagement', 'failure');
      });
  }

  function logoutUser() {
    xfetch('/api/auth/logout', { method: 'DELETE' })
      .then(function () {
        eventTag('logout', 'engagement', 'success', () => {
          user.delete();
          localStorage.removeItem('theme');
          window.setTheme();
          router.route(appRoutes.landing, true);
        });
      })
      .catch(function () {
        notifications.danger($LL.logoutError(AppConfig.FriendlyUIVerbs));
        eventTag('logout', 'engagement', 'failure');
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
      pictureUrl="{profile.picture_url}"
      gravatarHash="{profile.gravatarHash}"
      class="object-cover w-10 h-10 rounded-full"
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
        <UserIcon class="h-4 w-4 me-3" />
        <span>{$LL.profile()}</span>
      </button>
    </li>
    {#if $user.rank === 'GUEST'}
      <li class="flex">
        <button
          class="inline-flex items-center w-full px-2 py-1 font-semibold transition-colors duration-150 rounded-md hover:bg-gray-100 hover:text-gray-800 dark:hover:bg-gray-800 dark:hover:text-gray-200"
          data-testid="create-account-link"
          on:click="{goToRegister}"
        >
          <VoteIcon class="w-4 h-4 me-3 " />
          <span>{$LL.createAccount()}</span>
        </button>
      </li>
      <li class="flex">
        <button
          class="inline-flex items-center w-full px-2 py-1 font-semibold transition-colors duration-150 rounded-md hover:bg-gray-100 hover:text-gray-800 dark:hover:bg-gray-800 dark:hover:text-gray-200"
          data-testid="login-link"
          on:click="{goToLogin}"
        >
          <LockIcon class="w-4 h-4 me-3 " />
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
          <svg
            class="w-4 h-4 me-3"
            aria-hidden="true"
            fill="none"
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              d="M11 16l-4-4m0 0l4-4m-4 4h14m-5 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h7a3 3 0 013 3v1"
            ></path>
          </svg>
          <span>{$LL.logout()}</span>
        </button>
      </li>
    {/if}
  </ul>
{/if}
