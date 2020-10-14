<script>
    import Navaid from 'navaid'
    import { onDestroy } from 'svelte'
    import { _, locale, setupI18n, isLocaleLoaded } from './i18n'

    import Notifications from './components/Notifications.svelte'
    import WarriorIcon from './components/icons/WarriorIcon.svelte'
    import HollowButton from './components/HollowButton.svelte'
    import LocaleSwitcher from './components/LocaleSwitcher.svelte'
    import Landing from './pages/Landing.svelte'
    import Battles from './pages/Battles.svelte'
    import Battle from './pages/Battle.svelte'
    import Register from './pages/Register.svelte'
    import Login from './pages/Login.svelte'
    import ResetPassword from './pages/ResetPassword.svelte'
    import VerifyAccount from './pages/VerifyAccount.svelte'
    import WarriorProfile from './pages/WarriorProfile.svelte'
    import Admin from './pages/Admin.svelte'
    import { warrior } from './stores.js'
    import eventTag from './eventTag.js'
    import apiclient from './apiclient.js'

    setupI18n()

    const { AllowRegistration, AppVersion } = appConfig
    const footerLinkClasses = 'no-underline text-teal-500 hover:text-teal-800'

    let notifications

    let activeWarrior

    const unsubscribe = warrior.subscribe(w => {
        activeWarrior = w
    })

    let currentPage = {
        route: Landing,
        params: {},
    }

    const router = Navaid('/')
        .on('/', () => {
            currentPage = {
                route: Landing,
                params: {},
            }
        })
        .on('/enlist/:battleId?', params => {
            currentPage = {
                route: Register,
                params,
            }
        })
        .on('/login/:battleId?', params => {
            currentPage = {
                route: Login,
                params,
            }
        })
        .on('/reset-password/:resetId', params => {
            currentPage = {
                route: ResetPassword,
                params,
            }
        })
        .on('/verify-account/:verifyId', params => {
            currentPage = {
                route: VerifyAccount,
                params,
            }
        })
        .on('/warrior-profile', params => {
            currentPage = {
                route: WarriorProfile,
                params,
            }
        })
        .on('/battles', () => {
            currentPage = {
                route: Battles,
                params: {},
            }
        })
        .on('/battle/:battleId', params => {
            currentPage = {
                route: Battle,
                params,
            }
        })
        .on('/admin', () => {
            currentPage = {
                route: Admin,
                params: {},
            }
        })
        .listen()

    const xfetch = apiclient(handle401)

    function handle401() {
        eventTag('session_expired', 'engagement', 'unauthorized', () => {
            warrior.delete()
            router.route('/login')
        })
    }

    function logoutWarrior() {
        xfetch('/api/auth/logout', { method: 'POST' })
            .then(function() {
                eventTag('logout', 'engagement', 'success', () => {
                    warrior.delete()
                    router.route('/', true)
                })
            })
            .catch(function(error) {
                notifications.danger($_('actions.logout.failure'))
                eventTag('logout', 'engagement', 'failure')
            })
    }

    onDestroy(router.unlisten)
</script>

<style>
    :global(.nav-logo) {
        max-height: 3.75rem;
    }
    :global(.text-yellow-thunder) {
        color: #ffdd57;
    }
    :global(.bg-yellow-thunder) {
        background-color: #ffdd57;
    }
</style>

<Notifications bind:this="{notifications}" />

{#if $isLocaleLoaded}
    <nav
        class="flex items-center justify-between flex-wrap bg-white p-6"
        role="navigation"
        aria-label="main navigation">
        <div class="flex items-center flex-shrink-0 mr-6">
            <a href="/">
                <img src="/img/logo.svg" alt="Thunderdome" class="nav-logo" />
            </a>
        </div>
        <div class="text-right mt-4 md:mt-0">
            {#if activeWarrior.name}
                <span class="font-bold mr-2 text-xl">
                    <WarriorIcon />
                    <a href="/warrior-profile">{activeWarrior.name}</a>
                </span>
                <HollowButton
                    color="teal"
                    href="/battles"
                    additionalClasses="mr-2">
                    {$_('pages.myBattles.nav')}
                </HollowButton>
                {#if !activeWarrior.rank || activeWarrior.rank === 'PRIVATE'}
                    {#if AllowRegistration}
                        <HollowButton
                            color="teal"
                            href="/enlist"
                            additionalClasses="mr-2">
                            {$_('pages.createAccount.nav')}
                        </HollowButton>
                    {/if}
                    <HollowButton href="/login">
                        {$_('pages.login.nav')}
                    </HollowButton>
                {:else}
                    {#if activeWarrior.rank === 'GENERAL'}
                        <HollowButton
                            color="purple"
                            href="/admin"
                            additionalClasses="mr-2">
                            {$_('pages.admin.nav')}
                        </HollowButton>
                    {/if}
                    <HollowButton color="red" onClick="{logoutWarrior}">
                        {$_('actions.logout.button')}
                    </HollowButton>
                {/if}
            {:else}
                {#if AllowRegistration}
                    <HollowButton
                        color="teal"
                        href="/enlist"
                        additionalClasses="mr-2">
                        {$_('pages.createAccount.nav')}
                    </HollowButton>
                {/if}
                <HollowButton href="/login">
                    {$_('pages.login.nav')}
                </HollowButton>
            {/if}
            <LocaleSwitcher
                selectedLocale="{$locale}"
                on:locale-changed="{e => setupI18n({
                        withLocale: e.detail,
                    })}" />
        </div>
    </nav>

    <svelte:component
        this="{currentPage.route}"
        {...currentPage.params}
        {notifications}
        {router}
        {eventTag}
        {xfetch} />

    <footer class="p-6 text-center">
        <a
            href="https://github.com/StevenWeathers/thunderdome-planning-poker"
            class="{footerLinkClasses}">
            {$_('appName')}
        </a>
        {@html $_('footer.authoredBy', {
            values: {
                authorOpen: `<a href="http://stevenweathers.com" class="${footerLinkClasses}">`,
                authorClose: `</a>`,
            },
        })}
        {@html $_('footer.license', {
            values: {
                licenseOpen: `<a href="http://www.apache.org/licenses/" class="${footerLinkClasses}">`,
                licenseClose: `</a>`,
            },
        })}
        <br />
        {@html $_('footer.poweredBy', {
            values: {
                svelteOpen: `<a href="https://svelte.dev/" class="${footerLinkClasses}">`,
                svelteClose: `</a>`,
                goOpen: `<a href="https://golang.org/" class="${footerLinkClasses}">`,
                goClose: `</a>`,
            },
        })}
        <div class="text-sm text-gray-500">
            {$_('appVersion', { values: { version: AppVersion } })}
        </div>
    </footer>
{:else}
    <p>Loading...</p>
{/if}
