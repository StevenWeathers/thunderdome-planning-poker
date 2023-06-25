<script lang="ts">
    import './app.css'
    import './unreset.css'
    import '../../node_modules/quill/dist/quill.core.css'
    import '../../node_modules/quill/dist/quill.snow.css'
    import Navaid from 'navaid'
    import { onDestroy, onMount } from 'svelte'

    import { AppConfig, appRoutes } from './config'
    import apiclient from './apiclient'
    import { dir, warrior } from './stores'
    import eventTag from './eventTag'

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
    import AdminOrganization from './pages/admin/Organization.svelte'
    import AdminDepartment from './pages/admin/Department.svelte'
    import AdminTeams from './pages/admin/Teams.svelte'
    import AdminTeam from './pages/admin/Team.svelte'
    import AdminApikeys from './pages/admin/ApiKeys.svelte'
    import AdminAlerts from './pages/admin/Alerts.svelte'
    import AdminBattles from './pages/admin/Battles.svelte'
    import AdminBattle from './pages/admin/Battle.svelte'
    import AdminRetros from './pages/admin/Retros.svelte'
    import AdminRetro from './pages/admin/Retro.svelte'
    import AdminStoryboards from './pages/admin/Storyboards.svelte'
    import AdminStoryboard from './pages/admin/Storyboard.svelte'
    import { setLocale } from './i18n/i18n-svelte'
    import { detectLocale } from './i18n/i18n-util'
    import { loadLocaleAsync } from './i18n/i18n-util.async'

    const { FeaturePoker, FeatureRetro, FeatureStoryboard } = AppConfig

    let notifications

    let activeWarrior
    warrior.subscribe(w => {
        activeWarrior = w
    })

    $: if (document.dir !== $dir) {
        document.dir = $dir
    }

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
    router.on(`${appRoutes.adminUsers}/:userId`, params => {
        currentPage = {
            route: AdminUser,
            params: params,
            name: 'admin',
        }
    })
    router.on(`${appRoutes.adminUsers}`, () => {
        currentPage = {
            route: AdminUsers,
            params: {},
            name: 'admin',
        }
    })
    router.on(`${appRoutes.adminOrganizations}`, () => {
        currentPage = {
            route: AdminOrganizations,
            params: {},
            name: 'admin',
        }
    })
    router.on(`${appRoutes.adminOrganizations}/:organizationId`, params => {
        currentPage = {
            route: AdminOrganization,
            params: params,
            name: 'admin',
        }
    })
    router.on(
        `${appRoutes.adminOrganizations}/:organizationId/team/:teamId`,
        params => {
            currentPage = {
                route: AdminTeam,
                params: params,
                name: 'admin',
            }
        },
    )
    router.on(
        `${appRoutes.adminOrganizations}/:organizationId/department/:departmentId`,
        params => {
            currentPage = {
                route: AdminDepartment,
                params: params,
                name: 'admin',
            }
        },
    )
    router.on(
        `${appRoutes.adminOrganizations}/:organizationId/department/:departmentId/team/:teamId`,
        params => {
            currentPage = {
                route: AdminTeam,
                params: params,
                name: 'admin',
            }
        },
    )
    router.on(`${appRoutes.adminTeams}`, () => {
        currentPage = {
            route: AdminTeams,
            params: {},
            name: 'admin',
        }
    })
    router.on(`${appRoutes.adminTeams}/:teamId`, params => {
        currentPage = {
            route: AdminTeam,
            params: params,
            name: 'admin',
        }
    })
    router.on(`${appRoutes.adminApiKeys}`, () => {
        currentPage = {
            route: AdminApikeys,
            params: {},
            name: 'admin',
        }
    })
    router.on(`${appRoutes.adminAlerts}`, () => {
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
        router.on(`${appRoutes.adminBattles}`, () => {
            currentPage = {
                route: AdminBattles,
                params: {},
                name: 'admin',
            }
        })
        router.on(`${appRoutes.adminBattles}/:battleId`, params => {
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
        router.on(`${appRoutes.adminRetros}`, () => {
            currentPage = {
                route: AdminRetros,
                params: {},
                name: 'admin',
            }
        })
        router.on(`${appRoutes.adminRetros}/:retroId`, params => {
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
        router.on(`${appRoutes.adminStoryboards}`, () => {
            currentPage = {
                route: AdminStoryboards,
                params: {},
                name: 'admin',
            }
        })
        router.on(`${appRoutes.adminStoryboards}/:storyboardId`, params => {
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

    onMount(async () => {
        const detectedLocale = activeWarrior.locale || detectLocale()
        await loadLocaleAsync(detectedLocale)
        setLocale(detectedLocale)
    })

    onDestroy(router.unlisten)
</script>

<style lang="postcss" global>
    @tailwind base;
    @tailwind components;
    @tailwind utilities;
</style>

<Notifications bind:this="{notifications}" />

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
