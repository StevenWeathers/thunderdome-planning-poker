<script>
    import Navaid from 'navaid'
    import { onDestroy } from 'svelte'

    import { _, locale, setupI18n, isLocaleLoaded } from './i18n'
    import { appRoutes } from './config'
    import Notifications from './components/Notifications.svelte'
    import WarriorIcon from './components/icons/WarriorIcon.svelte'
    import HollowButton from './components/HollowButton.svelte'
    import LocaleSwitcher from './components/LocaleSwitcher.svelte'
    import Landing from './pages/Landing.svelte'
    import Battles from './pages/Battles.svelte'
    import Battle from './pages/Battle.svelte'
    import Organizations from './pages/Organizations.svelte'
    import Organization from './pages/Organization.svelte'
    import Department from './pages/Department.svelte'
    import Team from './pages/Team.svelte'
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

    const { AllowRegistration, AppVersion, PathPrefix } = appConfig
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
        .on(appRoutes.landing, () => {
            currentPage = {
                route: Landing,
                params: {},
            }
        })
        .on(`${appRoutes.register}/:battleId?`, params => {
            currentPage = {
                route: Register,
                params,
            }
        })
        .on(`${appRoutes.login}/:battleId?`, params => {
            currentPage = {
                route: Login,
                params,
            }
        })
        .on(`${appRoutes.resetPwd}/:resetId`, params => {
            currentPage = {
                route: ResetPassword,
                params,
            }
        })
        .on(`${appRoutes.verifyAct}/:verifyId`, params => {
            currentPage = {
                route: VerifyAccount,
                params,
            }
        })
        .on(appRoutes.profile, params => {
            currentPage = {
                route: WarriorProfile,
                params,
            }
        })
        .on(appRoutes.battles, () => {
            currentPage = {
                route: Battles,
                params: {},
            }
        })
        .on(`${appRoutes.battle}/:battleId`, params => {
            currentPage = {
                route: Battle,
                params,
            }
        })
        .on(appRoutes.organizations, () => {
            currentPage = {
                route: Organizations,
                params: {},
            }
        })
        .on(`${appRoutes.organization}/:organizationId`, params => {
            currentPage = {
                route: Organization,
                params,
            }
        })
        .on(
            `${appRoutes.organization}/:organizationId/team/:teamId`,
            params => {
                currentPage = {
                    route: Team,
                    params,
                }
            },
        )
        .on(
            `${appRoutes.organization}/:organizationId/department/:departmentId`,
            params => {
                currentPage = {
                    route: Department,
                    params,
                }
            },
        )
        .on(`${appRoutes.team}/:teamId`, params => {
            currentPage = {
                route: Team,
                params,
            }
        })
        .on(appRoutes.admin, () => {
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
            router.route(appRoutes.login)
        })
    }

    function logoutWarrior() {
        xfetch('/api/auth/logout', { method: 'POST' })
            .then(function() {
                eventTag('logout', 'engagement', 'success', () => {
                    warrior.delete()
                    router.route(appRoutes.landing, true)
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
            <a href="{appRoutes.landing}">
                <img
                    src="{PathPrefix}/img/logo.svg"
                    alt="Thunderdome"
                    class="nav-logo" />
            </a>
        </div>
        <div class="text-right mt-4 md:mt-0">
            {#if activeWarrior.name}
                <span class="font-bold mr-2 text-xl">
                    <WarriorIcon />
                    <a href="{appRoutes.profile}">{activeWarrior.name}</a>
                </span>
                <HollowButton
                    color="teal"
                    href="{appRoutes.battles}"
                    additionalClasses="mr-2">
                    {$_('pages.myBattles.nav')}
                </HollowButton>
                {#if activeWarrior.rank !== 'PRIVATE'}
                    <HollowButton
                        color="blue"
                        href="{appRoutes.organizations}"
                        additionalClasses="mr-2">
                        Organizations &amp; Teams
                    </HollowButton>
                {/if}
                {#if !activeWarrior.rank || activeWarrior.rank === 'PRIVATE'}
                    {#if AllowRegistration}
                        <HollowButton
                            color="teal"
                            href="{appRoutes.register}"
                            additionalClasses="mr-2">
                            {$_('pages.createAccount.nav')}
                        </HollowButton>
                    {/if}
                    <HollowButton href="{appRoutes.login}">
                        {$_('pages.login.nav')}
                    </HollowButton>
                {:else}
                    {#if activeWarrior.rank === 'GENERAL'}
                        <HollowButton
                            color="purple"
                            href="{appRoutes.admin}"
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
                        href="{appRoutes.register}"
                        additionalClasses="mr-2">
                        {$_('pages.createAccount.nav')}
                    </HollowButton>
                {/if}
                <HollowButton href="{appRoutes.login}">
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
