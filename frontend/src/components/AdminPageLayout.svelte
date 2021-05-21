<script>
    import { appRoutes } from '../config'
    import { _ } from '../i18n'

    export let activePage = 'admin'

    const { APIEnabled } = appConfig

    const pages = [
        {
            name: 'Admin',
            path: '',
        },
        {
            name: 'Alerts',
            path: '/alerts',
        },
        {
            name: 'Organizations',
            path: '/organizations',
        },
        {
            name: 'Teams',
            path: '/teams',
        },
        {
            name: 'Users',
            path: '/users',
        },
    ]

    if (APIEnabled) {
        pages.push({
            name: 'API Keys',
            path: '/apikeys',
        })
    }

    let activePillClasses = 'border-blue-500 bg-blue-500 text-white'
    let nonActivePillClasses =
        'border-gray-300 hover:border-gray-400 bg-gray-300 text-blue-500 hover:bg-gray-400'
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

<div class="flex px-6 py-2 border-b-2 bg-gray-300 border-gray-400">
    <div class="w-1/3">{$_('pages.admin.title')}</div>
    <div class="w-2/3">
        <ul class="flex justify-end">
            {#each pages as page}
                <li class="ml-3">
                    <a
                        class="admin-nav-pill {activePage === page.name.toLowerCase() ? activePillClasses : nonActivePillClasses}"
                        href="{appRoutes.admin}{page.path}">
                        {page.name}
                    </a>
                </li>
            {/each}
        </ul>
    </div>
</div>

<section>
    <div class="container mx-auto px-4 py-4 md:py-6 lg:py-8">
        <slot />
    </div>
</section>
