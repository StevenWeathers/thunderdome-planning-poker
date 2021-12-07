<script>
    import Navaid from 'navaid'
    import { onDestroy } from 'svelte'

    import { isLocaleLoaded, setupI18n } from './i18n'
    import { appRoutes } from './config'
    import apiclient from './apiclient.js'
    import { warrior } from './stores.js'
    import eventTag from './eventTag.js'

    import Notifications from './components/Notifications.svelte'
    import GlobalHeader from './components/GlobalHeader.svelte'
    import GlobalAlerts from './components/alert/GlobalAlerts.svelte'
    import GlobalFooter from './components/GlobalFooter.svelte'
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
    import WarriorProfile from './pages/UserProfile.svelte'
    import Admin from './pages/admin/Admin.svelte'
    import AdminUsers from './pages/admin/Users.svelte'
    import AdminUser from './pages/admin/User.svelte'
    import AdminOrganizations from './pages/admin/Organizations.svelte'
    import AdminTeams from './pages/admin/Teams.svelte'
    import AdminApikeys from './pages/admin/ApiKeys.svelte'
    import AdminAlerts from './pages/admin/Alerts.svelte'
    import AdminBattles from './pages/admin/Battles.svelte'
    import AdminBattle from './pages/admin/Battle.svelte'

    let notifications

    let activeWarrior
    warrior.subscribe(w => {
        activeWarrior = w
    })

    setupI18n({
        withLocale: activeWarrior ? activeWarrior.locale : 'en',
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
        .on(
            `${appRoutes.organization}/:organizationId/department/:departmentId/team/:teamId`,
            params => {
                currentPage = {
                    route: Team,
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
        .on(`${appRoutes.admin}/users/:userId`, params => {
            currentPage = {
                route: AdminUser,
                params: params,
            }
        })
        .on(`${appRoutes.admin}/users`, () => {
            currentPage = {
                route: AdminUsers,
                params: {},
            }
        })
        .on(`${appRoutes.admin}/organizations`, () => {
            currentPage = {
                route: AdminOrganizations,
                params: {},
            }
        })
        .on(`${appRoutes.admin}/teams`, () => {
            currentPage = {
                route: AdminTeams,
                params: {},
            }
        })
        .on(`${appRoutes.admin}/apikeys`, () => {
            currentPage = {
                route: AdminApikeys,
                params: {},
            }
        })
        .on(`${appRoutes.admin}/alerts`, () => {
            currentPage = {
                route: AdminAlerts,
                params: {},
            }
        })
        .on(`${appRoutes.admin}/battles/:battleId`, params => {
            currentPage = {
                route: AdminBattle,
                params: params,
            }
        })
        .on(`${appRoutes.admin}/battles`, () => {
            currentPage = {
                route: AdminBattles,
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

    onDestroy(router.unlisten)
</script>

<Notifications bind:this="{notifications}" />

{#if $isLocaleLoaded}
    <header>
        <GlobalAlerts registered="{!!activeWarrior.name}" />

        <GlobalHeader
            router="{router}"
            eventTag="{eventTag}"
            xfetch="{xfetch}"
            notifications="{notifications}"
        />
    </header>

    <main class="flex-grow">
        <svelte:component
            this="{currentPage.route}"
            {...currentPage.params}
            notifications="{notifications}"
            router="{router}"
            eventTag="{eventTag}"
            xfetch="{xfetch}"
        />
    </main>

    <GlobalFooter />
{:else}
    <p>Loading...</p>
{/if}
