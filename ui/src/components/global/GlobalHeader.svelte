<script lang="ts">
  import { validateUserIsAdmin } from '../../validationUtils';
  import { user } from '../../stores';
  import { AppConfig, appRoutes } from '../../config';
  import LL, { locale, setLocale } from '../../i18n/i18n-svelte';
  import type { Locales } from '../../i18n/i18n-types';
  import { loadLocaleAsync } from '../../i18n/i18n-util.async';
  import ThemeSelector from './ThemeSelector.svelte';
  import NavUserMenu from './NavUserMenu.svelte';
  import { ArrowRight, Menu, X } from '@lucide/svelte';
  import LocaleMenu from './LocaleMenu.svelte';
  import DomeLogo from '../logos/DomeLogo.svelte';
  import DomeLogoLight from '../logos/DomeLogoLight.svelte';

  import type { NotificationService } from '../../types/notifications';
  import type { ApiClient } from '../../types/apiclient';
  import type { SessionUser } from '../../types/user';

  interface Props {
    xfetch: ApiClient;
    router: any;
    notifications: NotificationService;
    currentPage: any;
  }

  let { xfetch, router, notifications, currentPage }: Props = $props();

  const setupI18n = async (locale: Locales) => {
    await loadLocaleAsync(locale);
    setLocale(locale);
  };

  const { FeaturePoker, FeatureRetro, FeatureStoryboard, HeaderAuthEnabled, SubscriptionsEnabled } = AppConfig;

  const activePageClass =
    'relative flex items-center px-4 py-3 lg:py-6 text-indigo-600 dark:text-yellow-thunder font-semibold transition-all duration-300 after:absolute after:bottom-0 after:start-0 after:w-full after:h-1 after:bg-indigo-600 after:dark:bg-yellow-thunder lg:after:h-1 after:rounded-t-md';

  const pageClass =
    'relative flex items-center px-4 py-3 lg:py-6 text-gray-700 dark:text-gray-300 font-medium hover:text-indigo-600 dark:hover:text-yellow-thunder transition-all duration-300 after:absolute after:bottom-0 after:start-0 after:w-0 after:h-1 after:bg-indigo-600 after:dark:bg-yellow-thunder hover:after:w-full after:transition-all after:duration-300 after:rounded-t-md';

  const mobileActiveClass =
    'block px-4 py-3 text-indigo-600 dark:text-yellow-thunder font-semibold bg-indigo-50 dark:bg-indigo-900/20 border-s-4 border-indigo-600 dark:border-yellow-thunder transition-all duration-300';

  const mobilePageClass =
    'block px-4 py-3 text-gray-700 dark:text-gray-300 font-medium hover:text-indigo-600 dark:hover:text-yellow-thunder hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-all duration-300';

  let showMobileMenu = $state(false);

  function toggleMobileMenu(e?: Event) {
    !!e && e.preventDefault();
    showMobileMenu = !showMobileMenu;
  }

  function headerLogin() {
    xfetch('/api/auth', { skip401Redirect: true })
      .then(res => res.json())
      .then(function (result) {
        const u = result.data.user;
        const newUser = {
          id: u.id,
          name: u.name,
          email: u.email,
          rank: u.rank,
          locale: u.locale,
          notificationsEnabled: u.notificationsEnabled,
        };

        user.create(newUser as SessionUser);
        setupI18n(newUser.locale);
        router.route(appRoutes.games, true);
      })
      .catch(function () {
        notifications.danger($LL.authError());
      });
  }

  function clickOutsideMobileMenu(node: HTMLElement) {
    function handleClick(event: MouseEvent) {
      if (showMobileMenu && !node.contains(event.target as Node)) {
        toggleMobileMenu();
      }
    }

    document.addEventListener('click', handleClick, true);

    return {
      destroy() {
        document.removeEventListener('click', handleClick, true);
      },
    };
  }
</script>

<header
  class="bg-white/80 dark:bg-gray-900/80 backdrop-blur-md border-b border-gray-200 dark:border-gray-800 shadow-sm relative z-50"
>
  <nav class="mx-auto max-w-screen-2xl px-4 sm:px-6 lg:px-8" aria-label="main navigation">
    <div class="flex h-16 lg:h-20 items-center justify-between">
      <!-- Logo Section -->
      <div class="flex items-center">
        <a href={appRoutes.landing} class="group flex items-center transition-transform duration-300 hover:scale-105">
          <DomeLogo
            class="hidden dark:inline-block h-8 md:h-10 lg:h-12 transition-all duration-300 group-hover:drop-shadow-lg"
          />
          <DomeLogoLight
            class="inline-block dark:hidden h-8 md:h-10 lg:h-12 transition-all duration-300 group-hover:drop-shadow-lg"
          />
        </a>
      </div>

      <!-- Desktop Navigation -->
      <div class="hidden lg:flex lg:items-center lg:space-x-1 font-rajdhani uppercase text-lg tracking-wide">
        {#if $user.name}
          <a href={appRoutes.dashboard} class={currentPage === 'dashboard' ? activePageClass : pageClass}>
            Dashboard
          </a>
        {/if}
        {#if FeaturePoker}
          <a href={appRoutes.games} class={currentPage === 'battles' ? activePageClass : pageClass}>
            {$LL.battles()}
          </a>
        {/if}
        {#if FeatureRetro}
          <a href={appRoutes.retros} class={currentPage === 'retros' ? activePageClass : pageClass}>
            {$LL.retros()}
          </a>
        {/if}
        {#if FeatureStoryboard}
          <a href={appRoutes.storyboards} class={currentPage === 'storyboards' ? activePageClass : pageClass}>
            {$LL.storyboards()}
          </a>
        {/if}
        {#if $user.name && $user.rank !== 'GUEST' && $user.rank !== 'PRIVATE'}
          <a href={appRoutes.teams} class={currentPage === 'teams' ? activePageClass : pageClass}>
            {$LL.teams()}
          </a>
        {/if}
        {#if SubscriptionsEnabled && !$user.subscribed}
          <a href={appRoutes.subscriptionPricing} class={currentPage === 'pricing' ? activePageClass : pageClass}>
            Pricing
          </a>
        {/if}
        {#if $user.name && validateUserIsAdmin($user)}
          <a href={appRoutes.admin} class={currentPage === 'admin' ? activePageClass : pageClass}>
            {$LL.admin()}
          </a>
        {/if}
      </div>

      <!-- Right Side Actions -->
      <div class="flex items-center space-x-3 rtl:space-x-reverse">
        {#if !$user.id}
          <!-- Guest User Controls -->
          <div class="hidden sm:flex items-center space-x-3 rtl:space-x-reverse">
            <LocaleMenu selectedLocale={$locale} update={(l: Locales) => setupI18n(l)} />
            <ThemeSelector />
            {#if HeaderAuthEnabled}
              <button
                onclick={headerLogin}
                class="group bg-indigo-600 hover:bg-indigo-700 text-white font-semibold px-6 py-2.5 rounded-lg transition-all duration-300 transform hover:scale-105 flex items-center space-x-2 rtl:space-x-reverse shadow-md hover:shadow-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 dark:focus:ring-offset-gray-900"
              >
                <span>{$LL.login()}</span>
                <ArrowRight
                  class="h-4 w-4 transition-transform duration-300 group-hover:translate-x-1 rtl:group-hover:-translate-x-1"
                />
              </button>
            {:else}
              <a
                href={appRoutes.login}
                class="group bg-indigo-600 hover:bg-indigo-700 text-white font-semibold px-6 py-2.5 rounded-lg transition-all duration-300 transform hover:scale-105 flex items-center space-x-2 rtl:space-x-reverse shadow-md hover:shadow-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 dark:focus:ring-offset-gray-900"
              >
                <span>{$LL.login()}</span>
                <ArrowRight
                  class="h-4 w-4 transition-transform duration-300 group-hover:translate-x-1 rtl:group-hover:-translate-x-1"
                />
              </a>
            {/if}
          </div>
        {/if}

        {#if $user.id}
          <!-- Authenticated User Controls -->
          <NavUserMenu {xfetch} {notifications} {router} {currentPage} />
        {/if}

        <!-- Mobile Menu Button -->
        <button
          onclick={toggleMobileMenu}
          type="button"
          class="lg:hidden relative inline-flex items-center justify-center p-2 text-gray-500 dark:text-gray-400 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-gray-200 dark:focus:ring-gray-600 transition-colors duration-300"
        >
          <span class="sr-only">
            {showMobileMenu ? 'Close main menu' : 'Open main menu'}
          </span>
          {#if showMobileMenu}
            <X class="h-6 w-6 transition-transform duration-300 rotate-0" />
          {:else}
            <Menu class="h-6 w-6 transition-transform duration-300 rotate-0" />
          {/if}
        </button>
      </div>
    </div>

    <!-- Mobile Menu -->
    <div class="lg:hidden {showMobileMenu ? 'block' : 'hidden'}" inert={!showMobileMenu} use:clickOutsideMobileMenu>
      <!-- Mobile Navigation Links -->
      <div class="px-2 pt-2 pb-3 space-y-1 bg-white dark:bg-gray-900 border-t border-gray-200 dark:border-gray-700">
        {#if FeaturePoker}
          <a
            href={appRoutes.games}
            class="{currentPage === 'battles'
              ? mobileActiveClass
              : mobilePageClass} font-rajdhani uppercase tracking-wide"
            onclick={() => (showMobileMenu = false)}
          >
            {$LL.battles()}
          </a>
        {/if}
        {#if FeatureRetro}
          <a
            href={appRoutes.retros}
            class="{currentPage === 'retros'
              ? mobileActiveClass
              : mobilePageClass} font-rajdhani uppercase tracking-wide"
            onclick={() => (showMobileMenu = false)}
          >
            {$LL.retros()}
          </a>
        {/if}
        {#if FeatureStoryboard}
          <a
            href={appRoutes.storyboards}
            class="{currentPage === 'storyboards'
              ? mobileActiveClass
              : mobilePageClass} font-rajdhani uppercase tracking-wide"
            onclick={() => (showMobileMenu = false)}
          >
            {$LL.storyboards()}
          </a>
        {/if}
        {#if $user.name && $user.rank !== 'GUEST' && $user.rank !== 'PRIVATE'}
          <a
            href={appRoutes.teams}
            class="{currentPage === 'teams'
              ? mobileActiveClass
              : mobilePageClass} font-rajdhani uppercase tracking-wide"
            onclick={() => (showMobileMenu = false)}
          >
            {$LL.teams()}
          </a>
        {/if}
        {#if SubscriptionsEnabled && !$user.subscribed}
          <a
            href={appRoutes.subscriptionPricing}
            class="{currentPage === 'pricing'
              ? mobileActiveClass
              : mobilePageClass} font-rajdhani uppercase tracking-wide"
            onclick={() => (showMobileMenu = false)}
          >
            Pricing
          </a>
        {/if}
        {#if $user.name && validateUserIsAdmin($user)}
          <a
            href={appRoutes.admin}
            class="{currentPage === 'admin'
              ? mobileActiveClass
              : mobilePageClass} font-rajdhani uppercase tracking-wide"
            onclick={() => (showMobileMenu = false)}
          >
            {$LL.admin()}
          </a>
        {/if}
      </div>

      <!-- Mobile User Controls (for guests) -->
      {#if !$user.id}
        <div
          class="px-4 py-3 border-t border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800/50 space-y-3 sm:hidden"
        >
          <div class="flex items-center justify-between">
            <span class="text-sm font-medium text-gray-700 dark:text-gray-300">Language & Theme</span>
            <div class="flex items-center space-x-3 rtl:space-x-reverse">
              <LocaleMenu selectedLocale={$locale} update={(l: Locales) => setupI18n(l)} />
              <ThemeSelector />
            </div>
          </div>
          {#if HeaderAuthEnabled}
            <button
              onclick={() => {
                headerLogin();
                showMobileMenu = false;
              }}
              class="w-full group bg-indigo-600 hover:bg-indigo-700 text-white font-semibold px-6 py-3 rounded-lg transition-all duration-300 flex items-center justify-center space-x-2 rtl:space-x-reverse shadow-md hover:shadow-lg"
            >
              <span>{$LL.login()}</span>
              <ArrowRight
                class="h-4 w-4 transition-transform duration-300 group-hover:translate-x-1 rtl:group-hover:-translate-x-1"
              />
            </button>
          {:else}
            <a
              href={appRoutes.login}
              class="w-full group bg-indigo-600 hover:bg-indigo-700 text-white font-semibold px-6 py-3 rounded-lg transition-all duration-300 flex items-center justify-center space-x-2 rtl:space-x-reverse shadow-md hover:shadow-lg"
              onclick={() => (showMobileMenu = false)}
            >
              <span>{$LL.login()}</span>
              <ArrowRight
                class="h-4 w-4 transition-transform duration-300 group-hover:translate-x-1 rtl:group-hover:-translate-x-1"
              />
            </a>
          {/if}
        </div>
      {/if}
    </div>
  </nav>
</header>
