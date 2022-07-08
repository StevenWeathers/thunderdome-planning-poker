<script>
    import './app.css'
    import './unreset.css'
    import '../../node_modules/quill/dist/quill.core.css'
    import '../../node_modules/quill/dist/quill.snow.css'
    import Navaid from 'navaid'
    import { onDestroy } from 'svelte'

    import { isLocaleLoaded, setupI18n } from './i18n.js'
    import { AppConfig, appRoutes } from './config.js'
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
    import Retros from './pages/Retros.svelte'
    import Retro from './pages/Retro.svelte'
    import Storyboards from './pages/Storyboards.svelte'
    import Storyboard from './pages/Storyboard.svelte'
    import Teams from './pages/Teams.svelte'
    import Organization from './pages/Organization.svelte'
    import Department from './pages/Department.svelte'
    import Team from './pages/Team.svelte'
    import TeamCheckin from './pages/TeamCheckin.svelte'
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
    import AdminRetros from './pages/admin/Retros.svelte'
    import AdminRetro from './pages/admin/Retro.svelte'
    import AdminStoryboards from './pages/admin/Storyboards.svelte'
    import AdminStoryboard from './pages/admin/Storyboard.svelte'

    const { FeaturePoker, FeatureRetro, FeatureStoryboard } = AppConfig

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
        name: 'landing',
    }

    const router = Navaid('/')

    router.on(appRoutes.landing, () => {
        currentPage = {
            route: Landing,
            params: {},
            name: 'landing',
        }
    })
    router.on(`${appRoutes.register}`, params => {
        currentPage = {
            route: Register,
            params,
            name: 'register',
        }
    })
    router.on(`${appRoutes.login}`, params => {
        currentPage = {
            route: Login,
            params,
            name: 'login',
        }
    })
    router.on(`${appRoutes.resetPwd}/:resetId`, params => {
        currentPage = {
            route: ResetPassword,
            params,
            name: 'reset-password',
        }
    })
    router.on(`${appRoutes.verifyAct}/:verifyId`, params => {
        currentPage = {
            route: VerifyAccount,
            params,
            name: 'verify-account',
        }
    })
    router.on(appRoutes.profile, params => {
        currentPage = {
            route: WarriorProfile,
            params,
            name: 'profile',
        }
    })
    router.on(appRoutes.teams, () => {
        currentPage = {
            route: Teams,
            params: {},
            name: 'Teams',
        }
    })
    router.on(`${appRoutes.organization}/:organizationId`, params => {
        currentPage = {
            route: Organization,
            params,
            name: 'organizations',
        }
    })
    router.on(
        `${appRoutes.organization}/:organizationId/team/:teamId`,
        params => {
            currentPage = {
                route: Team,
                params,
                name: 'team',
            }
        },
    )
    router.on(
        `${appRoutes.organization}/:organizationId/team/:teamId/checkin`,
        params => {
            currentPage = {
                route: TeamCheckin,
                params,
                name: 'team',
            }
        },
    )
    router.on(
        `${appRoutes.organization}/:organizationId/department/:departmentId`,
        params => {
            currentPage = {
                route: Department,
                params,
                name: 'department',
            }
        },
    )
    router.on(
        `${appRoutes.organization}/:organizationId/department/:departmentId/team/:teamId`,
        params => {
            currentPage = {
                route: Team,
                params,
                name: 'team',
            }
        },
    )
    router.on(
        `${appRoutes.organization}/:organizationId/department/:departmentId/team/:teamId/checkin`,
        params => {
            currentPage = {
                route: TeamCheckin,
                params,
                name: 'team',
            }
        },
    )
    router.on(`${appRoutes.team}/:teamId`, params => {
        currentPage = {
            route: Team,
            params,
            name: 'team',
        }
    })
    router.on(`${appRoutes.team}/:teamId/checkin`, params => {
        currentPage = {
            route: TeamCheckin,
            params,
            name: 'team',
        }
    })
    router.on(appRoutes.admin, () => {
        currentPage = {
            route: Admin,
            params: {},
            name: 'admin',
        }
    })
    router.on(`${appRoutes.admin}/users/:userId`, params => {
        currentPage = {
            route: AdminUser,
            params: params,
            name: 'admin',
        }
    })
    router.on(`${appRoutes.admin}/users`, () => {
        currentPage = {
            route: AdminUsers,
            params: {},
            name: 'admin',
        }
    })
    router.on(`${appRoutes.admin}/organizations`, () => {
        currentPage = {
            route: AdminOrganizations,
            params: {},
            name: 'admin',
        }
    })
    router.on(`${appRoutes.admin}/teams`, () => {
        currentPage = {
            route: AdminTeams,
            params: {},
            name: 'admin',
        }
    })
    router.on(`${appRoutes.admin}/apikeys`, () => {
        currentPage = {
            route: AdminApikeys,
            params: {},
            name: 'admin',
        }
    })
    router.on(`${appRoutes.admin}/alerts`, () => {
        currentPage = {
            route: AdminAlerts,
            params: {},
            name: 'admin',
        }
    })

    if (FeaturePoker) {
        router.on(appRoutes.battles, () => {
            currentPage = {
                route: Battles,
                params: {},
                name: 'battles',
            }
        })
        router.on(`${appRoutes.battle}/:battleId`, params => {
            currentPage = {
                route: Battle,
                params,
                name: 'battle',
            }
        })
        router.on(`${appRoutes.admin}/battles`, () => {
            currentPage = {
                route: AdminBattles,
                params: {},
                name: 'admin',
            }
        })
        router.on(`${appRoutes.admin}/battles/:battleId`, params => {
            currentPage = {
                route: AdminBattle,
                params: params,
                name: 'admin',
            }
        })
        router.on(`${appRoutes.register}/battle/:battleId`, params => {
            currentPage = {
                route: Register,
                params,
                name: 'register',
            }
        })
        router.on(`${appRoutes.login}/battle/:battleId`, params => {
            currentPage = {
                route: Login,
                params,
                name: 'login',
            }
        })
    }

    if (FeatureRetro) {
        router.on(appRoutes.retros, () => {
            currentPage = {
                route: Retros,
                params: {},
                name: 'retros',
            }
        })
        router.on(`${appRoutes.retro}/:retroId`, params => {
            currentPage = {
                route: Retro,
                params,
                name: 'retro',
            }
        })
        router.on(`${appRoutes.admin}/retros`, () => {
            currentPage = {
                route: AdminRetros,
                params: {},
                name: 'admin',
            }
        })
        router.on(`${appRoutes.admin}/retros/:retroId`, params => {
            currentPage = {
                route: AdminRetro,
                params: params,
                name: 'admin',
            }
        })
        router.on(`${appRoutes.register}/retro/:retroId`, params => {
            currentPage = {
                route: Register,
                params,
                name: 'register',
            }
        })
        router.on(`${appRoutes.login}/retro/:retroId`, params => {
            currentPage = {
                route: Login,
                params,
                name: 'login',
            }
        })
    }

    if (FeatureStoryboard) {
        router.on(appRoutes.storyboards, () => {
            currentPage = {
                route: Storyboards,
                params: {},
                name: 'storyboards',
            }
        })
        router.on(`${appRoutes.storyboard}/:storyboardId`, params => {
            currentPage = {
                route: Storyboard,
                params,
                name: 'storyboard',
            }
        })
        router.on(`${appRoutes.admin}/storyboards`, () => {
            currentPage = {
                route: AdminStoryboards,
                params: {},
                name: 'admin',
            }
        })
        router.on(`${appRoutes.admin}/storyboards/:storyboardId`, params => {
            currentPage = {
                route: AdminStoryboard,
                params: params,
                name: 'admin',
            }
        })
        router.on(`${appRoutes.register}/storyboard/:storyboardId`, params => {
            currentPage = {
                route: Register,
                params,
                name: 'register',
            }
        })
        router.on(`${appRoutes.login}/storyboard/:storyboardId`, params => {
            currentPage = {
                route: Login,
                params,
                name: 'login',
            }
        })
    }

    router.listen()

    const xfetch = apiclient(handle401)

    function handle401(skipRedirect) {
        eventTag('session_expired', 'engagement', 'unauthorized', () => {
            warrior.delete()
            if (!skipRedirect) {
                router.route(appRoutes.login)
            }
        })
    }

    onDestroy(router.unlisten)
</script>

<style lang="postcss" global>
    @tailwind base;
    @tailwind components;
    @tailwind utilities;
</style>

<Notifications bind:this="{notifications}" />

{#if $isLocaleLoaded}
    <header class="w-full">
        <GlobalAlerts registered="{!!activeWarrior.name}" />

        <GlobalHeader
            router="{router}"
            eventTag="{eventTag}"
            xfetch="{xfetch}"
            notifications="{notifications}"
            currentPage="{currentPage.name}"
        />
    </header>

    <main class="flex-grow flex flex-wrap flex-col">
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
