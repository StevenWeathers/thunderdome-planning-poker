<script>
    import HollowButton from './HollowButton.svelte'
    import SolidButton from './SolidButton.svelte'
    import LocaleSwitcher from './LocaleSwitcher.svelte'
    import UserIcon from './icons/UserIcon.svelte'
    import { validateUserIsAdmin } from '../validationUtils'
    import { _, locale, setupI18n } from '../i18n.js'
    import { warrior } from '../stores.js'
    import { AppConfig, appRoutes } from '../config.js'

    export let xfetch
    export let router
    export let eventTag
    export let notifications
    export let currentPage

    const {
        AllowRegistration,
        PathPrefix,
        FeaturePoker,
        FeatureRetro,
        FeatureStoryboard,
        OrganizationsEnabled,
    } = AppConfig

    const activePageClass =
        'text-purple-700 border-purple-700 dark:text-yellow-400 dark:border-yellow-400'
    const pageClass =
        'text-gray-700 border-white dark:border-gray-800 dark:hover:border-yellow-400 hover:text-purple-700 dark:text-gray-300 dark:hover:text-yellow-400 transition duration-300'

    let showMobileMenu = false

    function toggleMobileMenu(e) {
        !!e && e.preventDefault()
        showMobileMenu = !showMobileMenu
    }

    function logoutWarrior() {
        xfetch('/api/auth/logout', { method: 'DELETE' })
            .then(function () {
                eventTag('logout', 'engagement', 'success', () => {
                    warrior.delete()
                    router.route(appRoutes.landing, true)
                })
            })
            .catch(function () {
                notifications.danger($_('logoutError'))
                eventTag('logout', 'engagement', 'failure')
            })
    }
</script>

<style>
    .nav-logo {
        height: 3.5rem;
    }
</style>

<nav class="bg-white dark:bg-gray-800" aria-label="main navigation">
    <div class="px-8">
        <div class="flex justify-between">
            <div class="flex space-x-7">
                <div>
                    <a href="{appRoutes.landing}" class="block my-3 mr-10">
                        <img
                            src="{PathPrefix}/img/logo.svg"
                            alt="Thunderdome"
                            class="nav-logo"
                        />
                    </a>
                </div>
                <nav
                    class="hidden lg:flex items-center space-x-1 font-semibold font-rajdhani uppercase text-lg lg:text-xl"
                >
                    {#if $warrior.name}
                        {#if FeaturePoker}
                            <a
                                href="{appRoutes.battles}"
                                class="pt-6 pb-4 px-4 border-b-4 {currentPage ===
                                'battles'
                                    ? activePageClass
                                    : pageClass}"
                            >
                                {$_('battles')}
                            </a>
                        {/if}
                        {#if FeatureRetro}
                            <a
                                href="{appRoutes.retros}"
                                class="pt-6 pb-4 px-4 border-b-4 {currentPage ==
                                'retros'
                                    ? activePageClass
                                    : pageClass}"
                            >
                                {$_('retros')}
                            </a>
                        {/if}
                        {#if FeatureStoryboard}
                            <a
                                href="{appRoutes.storyboards}"
                                class="pt-6 pb-4 px-4 border-b-4 {currentPage ==
                                'storyboards'
                                    ? activePageClass
                                    : pageClass}"
                            >
                                {$_('storyboards')}
                            </a>
                        {/if}
                        {#if $warrior.rank !== 'GUEST' && $warrior.rank !== 'PRIVATE'}
                            <a
                                href="{appRoutes.teams}"
                                class="pt-6 pb-4 px-4  border-b-4 {currentPage ===
                                'teams'
                                    ? activePageClass
                                    : pageClass}"
                            >
                                {$_('teams')}
                            </a>
                        {/if}
                        {#if validateUserIsAdmin($warrior)}
                            <a
                                href="{appRoutes.admin}"
                                class="pt-6 pb-4 px-4 border-b-4 {currentPage ===
                                'admin'
                                    ? activePageClass
                                    : pageClass}"
                            >
                                {$_('pages.admin.nav')}
                            </a>
                        {/if}
                    {/if}
                </nav>
            </div>
            <div
                class="hidden lg:flex items-center space-x-3 font-rajdhani font-semibold dark:text-gray-300"
            >
                {#if !$warrior.name}
                    <div class="uppercase">
                        <a
                            href="{appRoutes.login}"
                            class="py-2 px-2 text-gray-700 dark:text-gray-300 hover:text-green-600 transition duration-300 text-lg lg:text-xl"
                            >{$_('pages.login.nav')}</a
                        >
                        {#if AllowRegistration}
                            <SolidButton
                                href="{appRoutes.register}"
                                additionalClasses="uppercase text-md lg:text-lg"
                                >{$_('pages.createAccount.nav')}</SolidButton
                            >
                        {/if}
                    </div>
                {:else}
                    <span
                        class="font-bold mr-2 text-lg lg:text-xl inline-block max-w-48 truncate"
                    >
                        <UserIcon class="h-5 w-5" />
                        <a
                            href="{appRoutes.profile}"
                            data-testid="userprofile-link">{$warrior.name}</a
                        >
                    </span>
                    {#if !$warrior.rank || $warrior.rank === 'GUEST' || $warrior.rank === 'PRIVATE'}
                        <a
                            href="{appRoutes.login}"
                            class="py-2 px-2 text-gray-700 dark:text-gray-300 hover:text-green-600 transition duration-300 uppercase text-xl"
                            >{$_('pages.login.nav')}</a
                        >
                        {#if AllowRegistration}
                            <SolidButton
                                href="{appRoutes.register}"
                                additionalClasses="uppercase text-md lg:text-lg"
                                >{$_('pages.createAccount.nav')}</SolidButton
                            >
                        {/if}
                    {:else}
                        <HollowButton
                            color="red"
                            onClick="{logoutWarrior}"
                            additionalClasses="uppercase text-md lg:text-lg"
                        >
                            {$_('logout')}
                        </HollowButton>
                    {/if}
                {/if}
                <LocaleSwitcher
                    class="ml-2 text-lg lg:text-xl"
                    selectedLocale="{$locale}"
                    on:locale-changed="{e =>
                        setupI18n({
                            withLocale: e.detail,
                        })}"
                />
            </div>
            <div class="lg:hidden flex items-center">
                <button
                    class="outline-none mobile-menu-button"
                    on:click="{toggleMobileMenu}"
                >
                    <svg
                        class="w-6 h-6 text-gray-500 dark:text-200 hover:text-green-500 dark:hover:text-yellow-400"
                        x-show="!showMenu"
                        fill="none"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                    >
                        <path d="M4 6h16M4 12h16M4 18h16"></path>
                    </svg>
                </button>
            </div>
        </div>
    </div>
    {#if showMobileMenu}
        <div
            class="lg:hidden py-2 border-t-2 border-gray-200 dark:border-gray-700"
        >
            <ul
                class="font-rajdhani font-semibold uppercase text-lg dark:text-white"
            >
                {#if $warrior.name}
                    {#if FeaturePoker}
                        <li>
                            <a
                                href="{appRoutes.battles}"
                                class="block p-4 hover:bg-green-500 dark:hover:bg-yellow-400 hover:text-white dark:hover:text-gray-800 transition duration-300"
                            >
                                {$_('battles')}
                            </a>
                        </li>
                    {/if}
                    {#if FeatureRetro}
                        <li>
                            <a
                                href="{appRoutes.retros}"
                                class="block p-4 hover:bg-green-500 dark:hover:bg-yellow-400 hover:text-white dark:hover:text-gray-800 transition duration-300"
                            >
                                {$_('retros')}
                            </a>
                        </li>
                    {/if}
                    {#if FeatureStoryboard}
                        <li>
                            <a
                                href="{appRoutes.storyboards}"
                                class="block p-4 hover:bg-green-500 dark:hover:bg-yellow-400 hover:text-white dark:hover:text-gray-800 transition duration-300"
                            >
                                {$_('storyboards')}
                            </a>
                        </li>
                    {/if}
                    {#if $warrior.rank !== 'GUEST' && $warrior.rank !== 'PRIVATE'}
                        <li>
                            <a
                                href="{appRoutes.organizations}"
                                class="block p-4 hover:bg-green-500 dark:hover:bg-yellow-400 hover:text-white dark:hover:text-gray-800 transition duration-300"
                            >
                                {$_('teams')}
                            </a>
                        </li>
                    {/if}
                    {#if validateUserIsAdmin($warrior)}
                        <li>
                            <a
                                href="{appRoutes.admin}"
                                class="block p-4 hover:bg-green-500 dark:hover:bg-yellow-400 hover:text-white dark:hover:text-gray-800 transition duration-300"
                            >
                                {$_('pages.admin.nav')}
                            </a>
                        </li>
                    {/if}
                {:else}
                    <li>
                        <a
                            href="{appRoutes.login}"
                            class="block p-4 hover:bg-green-500 dark:hover:bg-yellow-400 hover:text-white dark:hover:text-gray-800 transition duration-300"
                        >
                            {$_('pages.login.nav')}
                        </a>
                    </li>
                    {#if AllowRegistration}
                        <li>
                            <a
                                href="{appRoutes.register}"
                                class="block p-4 hover:bg-green-500 dark:hover:bg-yellow-400 hover:text-white dark:hover:text-gray-800 transition duration-300"
                            >
                                {$_('pages.createAccount.nav')}
                            </a>
                        </li>
                    {/if}
                {/if}
            </ul>
            <div class="font-rajdhani font-semibold mx-4 my-2 dark:text-white">
                {#if $warrior.name}
                    <span class="font-bold mr-2 text-lg lg:text-xl">
                        <UserIcon class="h-5 w-5" />
                        <a
                            href="{appRoutes.profile}"
                            data-testid="m-userprofile-link">{$warrior.name}</a
                        >
                    </span>
                    {#if !$warrior.rank || $warrior.rank === 'GUEST' || $warrior.rank === 'PRIVATE'}
                        <a
                            href="{appRoutes.login}"
                            class="py-2 px-2 text-gray-700 dark:text-gray-400 hover:text-green-600 transition duration-300 uppercase text-lg uppercase"
                            >{$_('pages.login.nav')}</a
                        >
                        {#if AllowRegistration}
                            <SolidButton
                                href="{appRoutes.register}"
                                additionalClasses="uppercase text-md lg:text-lg"
                                >{$_('pages.createAccount.nav')}</SolidButton
                            >
                        {/if}
                    {:else}
                        <HollowButton
                            color="red"
                            onClick="{logoutWarrior}"
                            additionalClasses="uppercase text-md lg:text-lg"
                        >
                            {$_('logout')}
                        </HollowButton>
                    {/if}
                {/if}
                <LocaleSwitcher
                    class="mt-4 block text-xl w-full"
                    selectedLocale="{$locale}"
                    on:locale-changed="{e =>
                        setupI18n({
                            withLocale: e.detail,
                        })}"
                />
            </div>
        </div>
    {/if}
</nav>
