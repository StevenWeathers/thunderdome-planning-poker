<script lang="ts">
    import HollowButton from './HollowButton.svelte'
    import UserIcon from './icons/UserIcon.svelte'
    import { validateUserIsAdmin } from '../validationUtils'
    import { warrior } from '../stores'
    import { AppConfig, appRoutes } from '../config'
    import LL, { locale, setLocale } from '../i18n/i18n-svelte'
    import SolidButton from './SolidButton.svelte'
    import type { Locales } from '../i18n/i18n-types'
    import { loadLocaleAsync } from '../i18n/i18n-util.async'
    import LocaleSwitcher from './LocaleSwitcher.svelte'

    export let xfetch
    export let router
    export let eventTag
    export let notifications
    export let currentPage

    const setupI18n = async (locale: Locales) => {
        await loadLocaleAsync(locale)
        setLocale(locale)
    }

    const {
        AllowRegistration,
        PathPrefix,
        FeaturePoker,
        FeatureRetro,
        FeatureStoryboard,
        OrganizationsEnabled,
        HeaderAuthEnabled,
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
                notifications.danger($LL.logoutError(AppConfig.FriendlyUIVerbs))
                eventTag('logout', 'engagement', 'failure')
            })
    }

    function headerLogin() {
        xfetch('/api/auth', { skip401Redirect: true })
            .then(res => res.json())
            .then(function (result) {
                const u = result.data.user
                const newUser = {
                    id: u.id,
                    name: u.name,
                    email: u.email,
                    rank: u.rank,
                    locale: u.locale,
                    notificationsEnabled: u.notificationsEnabled,
                }
                if (result.data.mfaRequired) {
                    mfaRequired = true
                    mfaUser = newUser
                    mfaSessionId = result.data.sessionId
                } else {
                    warrior.create(newUser)
                    eventTag('login', 'engagement', 'success', () => {
                        setupI18n(newUser.locale)
                        router.route(appRoutes.games, true)
                    })
                }
            })
            .catch(function (err) {
                notifications.danger(
                    $LL.authError({ friendly: AppConfig.FriendlyUIVerbs }),
                )
                eventTag('login', 'engagement', 'failure')
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
            <div class="flex space-x-7 rtl:space-x-reverse">
                <div>
                    <a href="{appRoutes.landing}" class="block my-3 me-10">
                        <img
                            src="{PathPrefix}/img/logo.svg"
                            alt="Thunderdome"
                            class="nav-logo"
                        />
                    </a>
                </div>
                <nav
                    class="hidden lg:flex items-center space-x-1 rtl:space-x-reverse font-semibold font-rajdhani uppercase text-lg lg:text-xl"
                >
                    {#if $warrior.name}
                        {#if FeaturePoker}
                            <a
                                href="{appRoutes.games}"
                                class="pt-6 pb-4 px-4 border-b-4 {currentPage ===
                                'battles'
                                    ? activePageClass
                                    : pageClass}"
                            >
                                {$LL.battles({
                                    friendly: AppConfig.FriendlyUIVerbs,
                                })}
                            </a>
                        {/if}
                        {#if FeatureRetro}
                            <a
                                href="{appRoutes.retros}"
                                class="pt-6 pb-4 px-4 border-b-4 {currentPage ===
                                'retros'
                                    ? activePageClass
                                    : pageClass}"
                            >
                                {$LL.retros()}
                            </a>
                        {/if}
                        {#if FeatureStoryboard}
                            <a
                                href="{appRoutes.storyboards}"
                                class="pt-6 pb-4 px-4 border-b-4 {currentPage ===
                                'storyboards'
                                    ? activePageClass
                                    : pageClass}"
                            >
                                {$LL.storyboards()}
                            </a>
                        {/if}
                        {#if $warrior.rank !== 'GUEST' && $warrior.rank !== 'PRIVATE'}
                            <a
                                href="{appRoutes.teams}"
                                class="pt-6 pb-4 px-4 border-b-4 {currentPage ===
                                'teams'
                                    ? activePageClass
                                    : pageClass}"
                            >
                                {$LL.teams()}
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
                                {$LL.admin()}
                            </a>
                        {/if}
                    {/if}
                </nav>
            </div>
            <div
                class="hidden lg:flex items-center space-x-3 rtl:space-x-reverse font-rajdhani font-semibold dark:text-gray-300"
            >
                {#if !$warrior.name}
                    <div class="uppercase">
                        {#if HeaderAuthEnabled}
                            <SolidButton
                                additionalClasses="uppercase text-md lg:text-lg"
                                onClick="{headerLogin}"
                                >{$LL.login()}</SolidButton
                            >
                        {:else}
                            <a
                                href="{appRoutes.login}"
                                class="py-2 px-2 text-gray-700 dark:text-gray-300 hover:text-green-600 transition duration-300 text-lg lg:text-xl"
                                >{$LL.login()}</a
                            >
                        {/if}
                        {#if AllowRegistration}
                            <SolidButton
                                href="{appRoutes.register}"
                                additionalClasses="uppercase text-md lg:text-lg"
                                >{$LL.createAccount()}</SolidButton
                            >
                        {/if}
                    </div>
                {:else}
                    <span
                        class="font-bold me-2 text-lg lg:text-xl inline-block max-w-48 truncate"
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
                            >{$LL.login()}</a
                        >
                        {#if AllowRegistration}
                            <SolidButton
                                href="{appRoutes.register}"
                                additionalClasses="uppercase text-md lg:text-lg"
                                >{$LL.createAccount()}</SolidButton
                            >
                        {/if}
                    {:else}
                        <HollowButton
                            color="red"
                            onClick="{logoutWarrior}"
                            additionalClasses="uppercase text-md lg:text-lg"
                        >
                            {$LL.logout()}
                        </HollowButton>
                    {/if}
                {/if}
                <LocaleSwitcher
                    class="ms-2 text-lg lg:text-xl"
                    selectedLocale="{$locale}"
                    on:locale-changed="{e => setupI18n(e.detail)}"
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
                                href="{appRoutes.games}"
                                class="block p-4 hover:bg-green-500 dark:hover:bg-yellow-400 hover:text-white dark:hover:text-gray-800 transition duration-300"
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
                                class="block p-4 hover:bg-green-500 dark:hover:bg-yellow-400 hover:text-white dark:hover:text-gray-800 transition duration-300"
                            >
                                {$LL.retros()}
                            </a>
                        </li>
                    {/if}
                    {#if FeatureStoryboard}
                        <li>
                            <a
                                href="{appRoutes.storyboards}"
                                class="block p-4 hover:bg-green-500 dark:hover:bg-yellow-400 hover:text-white dark:hover:text-gray-800 transition duration-300"
                            >
                                {$LL.storyboards()}
                            </a>
                        </li>
                    {/if}
                    {#if $warrior.rank !== 'GUEST' && $warrior.rank !== 'PRIVATE'}
                        <li>
                            <a
                                href="{appRoutes.teams}"
                                class="block p-4 hover:bg-green-500 dark:hover:bg-yellow-400 hover:text-white dark:hover:text-gray-800 transition duration-300"
                            >
                                {$LL.teams()}
                            </a>
                        </li>
                    {/if}
                    {#if validateUserIsAdmin($warrior)}
                        <li>
                            <a
                                href="{appRoutes.admin}"
                                class="block p-4 hover:bg-green-500 dark:hover:bg-yellow-400 hover:text-white dark:hover:text-gray-800 transition duration-300"
                            >
                                {$LL.admin()}
                            </a>
                        </li>
                    {/if}
                {:else}
                    <li>
                        <a
                            href="{appRoutes.login}"
                            class="block p-4 hover:bg-green-500 dark:hover:bg-yellow-400 hover:text-white dark:hover:text-gray-800 transition duration-300"
                        >
                            {$LL.login()}
                        </a>
                    </li>
                    {#if AllowRegistration}
                        <li>
                            <a
                                href="{appRoutes.register}"
                                class="block p-4 hover:bg-green-500 dark:hover:bg-yellow-400 hover:text-white dark:hover:text-gray-800 transition duration-300"
                            >
                                {$LL.createAccount()}
                            </a>
                        </li>
                    {/if}
                {/if}
            </ul>
            <div class="font-rajdhani font-semibold mx-4 my-2 dark:text-white">
                {#if $warrior.name}
                    <span class="font-bold me-2 text-lg lg:text-xl">
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
                            >{$LL.login()}</a
                        >
                        {#if AllowRegistration}
                            <SolidButton
                                href="{appRoutes.register}"
                                additionalClasses="uppercase text-md lg:text-lg"
                                >{$LL.createAccount()}</SolidButton
                            >
                        {/if}
                    {:else}
                        <HollowButton
                            color="red"
                            onClick="{logoutWarrior}"
                            additionalClasses="uppercase text-md lg:text-lg"
                        >
                            {$LL.logout()}
                        </HollowButton>
                    {/if}
                {/if}
                <LocaleSwitcher
                    class="mt-4 block text-xl w-full"
                    selectedLocale="{$locale}"
                    on:locale-changed="{e => setupI18n(e.detail)}"
                />
            </div>
        </div>
    {/if}
</nav>
