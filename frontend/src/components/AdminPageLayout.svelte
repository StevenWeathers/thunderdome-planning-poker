<script>
    import { AppConfig, appRoutes } from '../config.js'
    import { _ } from '../i18n.js'

    export let activePage = 'admin'

    const { ExternalAPIEnabled } = AppConfig

    const pages = [
        {
            name: 'Admin',
            label: $_('adminPageAdmin'),
            path: '',
        },
        {
            name: 'Alerts',
            label: $_('adminPageAlerts'),
            path: '/alerts',
        },
        {
            name: 'Battles',
            label: $_('battles'),
            path: '/battles',
        },
        {
            name: 'Organizations',
            label: $_('adminPageOrganizations'),
            path: '/organizations',
        },
        {
            name: 'Teams',
            label: $_('adminPageTeams'),
            path: '/teams',
        },
        {
            name: 'Users',
            label: $_('adminPageUsers'),
            path: '/users',
        },
    ]

    if (ExternalAPIEnabled) {
        pages.push({
            name: 'API Keys',
            label: $_('adminPageApi'),
            path: '/apikeys',
        })
    }

    let activePillClasses =
        'border-blue-500 bg-blue-500 text-white dark:bg-sky-500 dark:border-sky-500 dark:text-gray-900'
    let nonActivePillClasses =
        'border-gray-200 hover:border-gray-300 bg-gray-200 text-blue-500 dark:text-gray-400 dark:hover:text-gray-900 hover:bg-gray-300 dark:border-gray-700 dark:bg-gray-700 dark:hover:bg-gray-400 dark:hover:border-gray-400'
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

<div
    class="flex px-6 py-2 border-b-2 bg-gray-200 dark:bg-gray-700 border-gray-300 dark:border-gray-600"
>
    <div class="w-full">
        <ul class="flex justify-end">
            {#each pages as page}
                <li class="ml-3">
                    <a
                        class="admin-nav-pill {activePage ===
                        page.name.toLowerCase().replace(' ', '')
                            ? activePillClasses
                            : nonActivePillClasses}"
                        href="{appRoutes.admin}{page.path}"
                    >
                        {page.label}
                    </a>
                </li>
            {/each}
        </ul>
    </div>
</div>

<section>
    <div class="px-4 py-4 md:py-6 md:px-6 lg:py-8 lg:px-8">
        <slot />
    </div>
</section>
