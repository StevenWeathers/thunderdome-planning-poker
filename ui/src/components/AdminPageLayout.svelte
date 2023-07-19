<script lang="ts">
  import { AppConfig, appRoutes } from '../config';
  import LL from '../i18n/i18n-svelte';

  export let activePage = 'admin';

  const {
    ExternalAPIEnabled,
    FeaturePoker,
    FeatureRetro,
    FeatureStoryboard,
    OrganizationsEnabled,
  } = AppConfig;

  $: pages = $LL && [
    {
      name: 'Admin',
      label: $LL.adminPageAdmin(),
      path: '',
      enabled: true,
    },
    {
      name: 'Alerts',
      label: $LL.adminPageAlerts(),
      path: '/alerts',
      enabled: true,
    },
    {
      name: 'Battles',
      label: $LL.battles({ friendly: AppConfig.FriendlyUIVerbs }),
      path: '/battles',
      enabled: FeaturePoker,
    },
    {
      name: 'Retros',
      label: $LL.retros(),
      path: '/retros',
      enabled: FeatureRetro,
    },
    {
      name: 'Storyboards',
      label: $LL.storyboards(),
      path: '/storyboards',
      enabled: FeatureStoryboard,
    },
    {
      name: 'Organizations',
      label: $LL.adminPageOrganizations(),
      path: '/organizations',
      enabled: OrganizationsEnabled,
    },
    {
      name: 'Teams',
      label: $LL.adminPageTeams(),
      path: '/teams',
      enabled: true,
    },
    {
      name: 'Users',
      label: $LL.adminPageUsers(),
      path: '/users',
      enabled: true,
    },
    {
      name: 'API Keys',
      label: $LL.adminPageApi(),
      path: '/apikeys',
      enabled: ExternalAPIEnabled,
    },
  ];

  let activePillClasses =
    'border-blue-500 bg-blue-500 text-white dark:bg-sky-500 dark:border-sky-500 dark:text-gray-900';
  let nonActivePillClasses =
    'border-gray-200 hover:border-gray-300 bg-gray-200 text-blue-500 dark:text-gray-400 dark:hover:text-gray-900 hover:bg-gray-300 dark:border-gray-700 dark:bg-gray-700 dark:hover:bg-gray-400 dark:hover:border-gray-400';
</script>

<style>
  .admin-nav-pill {
    @apply inline-block;
    @apply border;
    @apply rounded;
    @apply py-1;
    @apply px-3;
  }
</style>

<div class="grow-0 w-full">
  <div
    class="w-full flex justify-end px-6 py-2 border-b-2 bg-gray-200 dark:bg-gray-700 border-gray-300 dark:border-gray-600"
  >
    <ul class="flex" data-testid="admin-nav">
      {#each pages as page}
        {#if page.enabled}
          <li class="ms-3">
            <a
              class="admin-nav-pill {activePage ===
              page.name.toLowerCase().replace(' ', '')
                ? activePillClasses
                : nonActivePillClasses}"
              href="{appRoutes.admin}{page.path}"
              data-testid="admin-nav-item"
            >
              {page.label}
            </a>
          </li>
        {/if}
      {/each}
    </ul>
  </div>
</div>

<section class="grow w-full">
  <div class="px-4 py-4 md:py-6 md:px-6 lg:py-8 lg:px-8">
    <slot />
  </div>
</section>
