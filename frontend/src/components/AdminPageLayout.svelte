<script>
    import { appRoutes } from '../config'
    import { _ } from '../i18n'

    export let activePage = 'admin'

    const { ExternalAPIEnabled } = appConfig

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

    let activePillClasses = 'border-blue-500 bg-blue-500 text-white'
    let nonActivePillClasses =
        'border-gray-300 hover:border-gray-400 bg-gray-300 text-blue-500 hover:bg-gray-400'
</script>

<style>
    :global(.admin-nav-pill) {
        @apply inline-block;
        @apply border;
        @apply rounded;
        @apply py-1;
        @apply px-3;
    }
</style>

<div class="flex px-6 py-2 border-b-2 bg-gray-300 border-gray-400">
    <div class="w-full">
        <ul class="flex justify-end">
            {#each pages as page}
                <li class="ml-3">
                    <a
                        class="admin-nav-pill {activePage === page.name.toLowerCase() ? activePillClasses : nonActivePillClasses}"
                        href="{appRoutes.admin}{page.path}">
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
