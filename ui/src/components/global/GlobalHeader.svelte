<script lang="ts">
  import { validateUserIsAdmin } from '../../validationUtils';
  import { user } from '../../stores';
  import { AppConfig, appRoutes } from '../../config';
  import LL, { locale, setLocale } from '../../i18n/i18n-svelte';
  import type { Locales } from '../../i18n/i18n-types';
  import { loadLocaleAsync } from '../../i18n/i18n-util.async';
  import ThemeSelector from './ThemeSelector.svelte';
  import NavUserMenu from './NavUserMenu.svelte';
  import ArrowRight from '../icons/ArrowRight.svelte';
  import LocaleMenu from './LocaleMenu.svelte';

  export let xfetch;
  export let router;
  export let eventTag;
  export let notifications;
  export let currentPage;

  const setupI18n = async (locale: Locales) => {
    await loadLocaleAsync(locale);
    setLocale(locale);
  };

  const {
    AllowRegistration,
    PathPrefix,
    FeaturePoker,
    FeatureRetro,
    FeatureStoryboard,
    OrganizationsEnabled,
    HeaderAuthEnabled,
    SubscriptionsEnabled,
  } = AppConfig;

  const activePageClass =
    'block lg:pt-6 lg:pb-4 px-4 border-b lg:border-b-4 lg:dark:bg-gray-800 text-indigo-600 border-indigo-600 dark:text-yellow-400 dark:border-yellow-400 transition duration-300';
  const pageClass =
    'block lg:pt-6 lg:pb-4 px-4 border-b border-gray-200 lg:border-white dark:border-gray-600 lg:dark:border-gray-800 lg:border-b-4 text-gray-700 hover:border-indigo-600 dark:hover:border-yellow-400 hover:text-indigo-700 dark:text-gray-300 dark:hover:text-yellow-400 transition duration-300';

  let showMobileMenu = false;

  function toggleMobileMenu(e) {
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

        user.create(newUser);
        eventTag('login', 'engagement', 'success', () => {
          setupI18n(newUser.locale);
          router.route(appRoutes.games, true);
        });
      })
      .catch(function (err) {
        notifications.danger(
          $LL.authError({ friendly: AppConfig.FriendlyUIVerbs }),
        );
        eventTag('login', 'engagement', 'failure');
      });
  }
</script>

<header>
  <nav
    class="bg-white dark:bg-gray-800 px-4 lg:px-6"
    aria-label="main navigation"
  >
    <div class="flex flex-wrap justify-between mx-auto max-w-screen-2xl">
      <div class="content-center py-2 lg:py-3">
        <a href="{appRoutes.landing}" class="block">
          <img
            src="{PathPrefix}/img/logo.svg"
            alt="Thunderdome"
            class="max-h-10 lg:max-h-[3.5rem]"
          />
        </a>
      </div>

      <ul
        class="flex items-center flex-shrink-0 space-x-2 rtl:space-x-reverse lg:order-2"
      >
        {#if !$user.id}
          <li>
            <LocaleMenu
              xfetch="{xfetch}"
              notifications="{notifications}"
              router="{router}"
              eventTag="{eventTag}"
              currentPage="{currentPage}"
              selectedLocale="{$locale}"
              on:locale-changed="{e => setupI18n(e.detail)}"
            />
          </li>
          <li class="flex">
            <ThemeSelector
              xfetch="{xfetch}"
              notifications="{notifications}"
              router="{router}"
              eventTag="{eventTag}"
              currentPage="{currentPage}"
            />
          </li>
          <li>
            {#if HeaderAuthEnabled}
              <button
                on:click="{headerLogin}"
                class="block py-2 px-4 rounded transition duration-300 bg-indigo-600 hover:bg-indigo-800 text-white"
                >{$LL.login()}
                <ArrowRight class="ms-2 inline-block" />
              </button>
            {:else}
              <a
                href="{appRoutes.login}"
                class="block py-2 px-4 rounded transition duration-300 bg-indigo-600 hover:bg-indigo-800 text-white"
                >{$LL.login()}
                <ArrowRight class="ms-2 inline-block" />
              </a>
            {/if}
          </li>
        {/if}
        {#if $user.id}
          <li class="relative">
            <NavUserMenu
              xfetch="{xfetch}"
              notifications="{notifications}"
              router="{router}"
              eventTag="{eventTag}"
              currentPage="{currentPage}"
            />
          </li>
        {/if}
        <li>
          <button
            on:click="{toggleMobileMenu}"
            type="button"
            class="inline-flex items-center p-2 text-sm text-gray-500 rounded-lg lg:hidden hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-gray-200 dark:text-gray-400 dark:hover:bg-gray-700 dark:focus:ring-gray-600"
          >
            <span class="sr-only">Open main menu</span>
            <svg
              class="w-6 h-6"
              fill="currentColor"
              viewBox="0 0 20 20"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                fill-rule="evenodd"
                d="M3 5a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zM3 10a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zM3 15a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1z"
                clip-rule="evenodd"></path>
            </svg>
          </button>
        </li>
      </ul>

      <div
        class="{showMobileMenu
          ? ''
          : 'hidden'} justify-between items-center w-full lg:flex lg:w-auto lg:order-1 font-semibold font-rajdhani uppercase text-lg lg:text-xl"
      >
        <ul
          class="flex flex-col mt-4 font-medium lg:flex-row lg:space-x-1 lg:mt-0 pb-2 lg:pb-0 rtl:space-x-reverse"
        >
          {#if $user.name}
            {#if FeaturePoker}
              <li>
                <a
                  href="{appRoutes.games}"
                  class="{currentPage === 'battles'
                    ? activePageClass
                    : pageClass}"
                >
                  {$LL.battles({
                    friendly: AppConfig.FriendlyUIVerbs,
                  })}
                </a>
              </li>
            {/if}
            {#if FeatureRetro}
              <li>
                <a
                  href="{appRoutes.retros}"
                  class="{currentPage === 'retros'
                    ? activePageClass
                    : pageClass}"
                >
                  {$LL.retros()}
                </a>
              </li>
            {/if}
            {#if FeatureStoryboard}
              <li>
                <a
                  href="{appRoutes.storyboards}"
                  class="{currentPage === 'storyboards'
                    ? activePageClass
                    : pageClass}"
                >
                  {$LL.storyboards()}
                </a>
              </li>
            {/if}
            {#if $user.rank !== 'GUEST' && $user.rank !== 'PRIVATE'}
              <li>
                <a
                  href="{appRoutes.teams}"
                  class="{currentPage === 'teams'
                    ? activePageClass
                    : pageClass}"
                >
                  {$LL.teams()}
                </a>
              </li>
            {/if}
            {#if validateUserIsAdmin($user)}
              <li>
                <a
                  href="{appRoutes.admin}"
                  class="{currentPage === 'admin'
                    ? activePageClass
                    : pageClass}"
                >
                  {$LL.admin()}
                </a>
              </li>
            {/if}
          {/if}
          {#if SubscriptionsEnabled && !$user.subscribed}
            <li>
              <a
                href="{appRoutes.subscriptionPricing}"
                class="{currentPage === 'pricing'
                  ? activePageClass
                  : pageClass}"
              >
                Pricing
              </a>
            </li>
          {/if}
        </ul>
      </div>
    </div>
  </nav>
</header>
