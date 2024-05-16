<script lang="ts">
  import './app.css';
  import './tailwind.css';
  import './unreset.css';
  import '../node_modules/quill/dist/quill.core.css';
  import '../node_modules/quill/dist/quill.snow.css';
  import Navaid from 'navaid';
  import { onDestroy } from 'svelte';

  import { AppConfig, appRoutes } from './config';
  import apiclient from './apiclient';
  import { dir, user } from './stores';
  import eventTag from './eventTag';

  import Notifications from './components/global/Notifications.svelte';
  import GlobalHeader from './components/global/GlobalHeader.svelte';
  import GlobalAlerts from './components/alert/GlobalAlerts.svelte';
  import GlobalFooter from './components/global/GlobalFooter.svelte';
  import Landing from './pages/Landing.svelte';
  import Battles from './pages/poker/PokerGames.svelte';
  import Battle from './pages/poker/PokerGame.svelte';
  import Retros from './pages/retro/Retros.svelte';
  import Retro from './pages/retro/Retro.svelte';
  import Storyboards from './pages/storyboard/Storyboards.svelte';
  import Storyboard from './pages/storyboard/Storyboard.svelte';
  import Teams from './pages/team/Teams.svelte';
  import Organization from './pages/team/Organization.svelte';
  import Department from './pages/team/Department.svelte';
  import Team from './pages/team/Team.svelte';
  import TeamCheckin from './pages/team/TeamCheckin.svelte';
  import Register from './pages/user/Register.svelte';
  import Login from './pages/user/Login.svelte';
  import ResetPassword from './pages/user/ResetPassword.svelte';
  import VerifyAccount from './pages/user/VerifyAccount.svelte';
  import WarriorProfile from './pages/user/UserProfile.svelte';
  import Admin from './pages/admin/Admin.svelte';
  import AdminUsers from './pages/admin/Users.svelte';
  import AdminUser from './pages/admin/User.svelte';
  import AdminOrganizations from './pages/admin/Organizations.svelte';
  import AdminOrganization from './pages/admin/Organization.svelte';
  import AdminDepartment from './pages/admin/Department.svelte';
  import AdminTeams from './pages/admin/Teams.svelte';
  import AdminTeam from './pages/admin/Team.svelte';
  import AdminApikeys from './pages/admin/ApiKeys.svelte';
  import AdminAlerts from './pages/admin/Alerts.svelte';
  import AdminBattles from './pages/admin/PokerGames.svelte';
  import AdminBattle from './pages/admin/PokerGame.svelte';
  import AdminRetros from './pages/admin/Retros.svelte';
  import AdminRetro from './pages/admin/Retro.svelte';
  import AdminStoryboards from './pages/admin/Storyboards.svelte';
  import AdminStoryboard from './pages/admin/Storyboard.svelte';
  import AdminSubscriptions from './pages/admin/Subscriptions.svelte';
  import AdminSubscription from './pages/admin/Subscription.svelte';
  import { setLocale } from './i18n/i18n-svelte';
  import Confirmation from './pages/subscription/Confirmation.svelte';
  import Pricing from './pages/subscription/Pricing.svelte';
  import PrivacyPolicy from './pages/support/PrivacyPolicy.svelte';
  import TermsConditions from './pages/support/TermsConditions.svelte';
  import Support from './pages/support/Support.svelte';
  import { loadLocale } from './i18n/i18n-util.sync';

  const {
    FeaturePoker,
    FeatureRetro,
    FeatureStoryboard,
    SubscriptionsEnabled,
  } = AppConfig;

  let notifications;

  let activeWarrior;
  user.subscribe(w => {
    activeWarrior = w;
  });

  const detectedLocale = activeWarrior.locale || AppConfig.DefaultLocale;
  loadLocale(detectedLocale);
  setLocale(detectedLocale);

  $: if (document.dir !== $dir) {
    document.dir = $dir;
  }

  let currentPage = {
    route: Landing,
    params: {},
    name: 'landing',
  };

  const router = Navaid('/');

  router.on(appRoutes.landing, () => {
    currentPage = {
      route: Landing,
      params: {},
      name: 'landing',
    };
  });
  router.on(appRoutes.privacyPolicy, () => {
    currentPage = {
      route: PrivacyPolicy,
      params: {},
      name: 'privacy',
    };
  });
  router.on(appRoutes.termsConditions, () => {
    currentPage = {
      route: TermsConditions,
      params: {},
      name: 'terms-conditions',
    };
  });
  router.on(appRoutes.support, () => {
    currentPage = {
      route: Support,
      params: {},
      name: 'support',
    };
  });
  router.on(`${appRoutes.register}`, params => {
    currentPage = {
      route: Register,
      params,
      name: 'register',
    };
  });
  router.on(`${appRoutes.login}`, params => {
    currentPage = {
      route: Login,
      params,
      name: 'login',
    };
  });
  router.on(`${appRoutes.login}/subscription`, () => {
    currentPage = {
      route: Login,
      params: { subscription: true },
      name: 'login',
    };
  });
  router.on(`${appRoutes.resetPwd}/:resetId`, params => {
    currentPage = {
      route: ResetPassword,
      params,
      name: 'reset-password',
    };
  });
  router.on(`${appRoutes.verifyAct}/:verifyId`, params => {
    currentPage = {
      route: VerifyAccount,
      params,
      name: 'verify-account',
    };
  });
  router.on(appRoutes.profile, params => {
    currentPage = {
      route: WarriorProfile,
      params,
      name: 'profile',
    };
  });
  router.on(appRoutes.teams, () => {
    currentPage = {
      route: Teams,
      params: {},
      name: 'teams',
    };
  });
  router.on(appRoutes.subscriptionPricing, () => {
    currentPage = {
      route: Pricing,
      params: {},
      name: 'pricing',
    };
  });
  router.on(appRoutes.subscriptionConfirmation, () => {
    currentPage = {
      route: Confirmation,
      params: {},
      name: 'subscription-confirmation',
    };
  });
  router.on(`${appRoutes.register}/subscription`, () => {
    currentPage = {
      route: Register,
      params: { subscription: true },
      name: 'register',
    };
  });
  router.on(`${appRoutes.register}/team/:teamInviteId`, params => {
    currentPage = {
      route: Register,
      params,
      name: 'register',
    };
  });
  router.on(`${appRoutes.register}/organization/:orgInviteId`, params => {
    currentPage = {
      route: Register,
      params,
      name: 'register',
    };
  });
  router.on(`${appRoutes.organization}/:organizationId`, params => {
    currentPage = {
      route: Organization,
      params,
      name: 'organization',
    };
  });
  router.on(
    `${appRoutes.organization}/:organizationId/team/:teamId`,
    params => {
      currentPage = {
        route: Team,
        params,
        name: 'team',
      };
    },
  );
  router.on(
    `${appRoutes.organization}/:organizationId/team/:teamId/checkin`,
    params => {
      currentPage = {
        route: TeamCheckin,
        params,
        name: 'team-checkin',
      };
    },
  );
  router.on(
    `${appRoutes.organization}/:organizationId/department/:departmentId`,
    params => {
      currentPage = {
        route: Department,
        params,
        name: 'department',
      };
    },
  );
  router.on(
    `${appRoutes.organization}/:organizationId/department/:departmentId/team/:teamId`,
    params => {
      currentPage = {
        route: Team,
        params,
        name: 'team',
      };
    },
  );
  router.on(
    `${appRoutes.organization}/:organizationId/department/:departmentId/team/:teamId/checkin`,
    params => {
      currentPage = {
        route: TeamCheckin,
        params,
        name: 'team-checkin',
      };
    },
  );
  router.on(`${appRoutes.team}/:teamId`, params => {
    currentPage = {
      route: Team,
      params,
      name: 'team',
    };
  });
  router.on(`${appRoutes.team}/:teamId/checkin`, params => {
    currentPage = {
      route: TeamCheckin,
      params,
      name: 'team-checkin',
    };
  });
  router.on(appRoutes.admin, () => {
    currentPage = {
      route: Admin,
      params: {},
      name: 'admin',
    };
  });
  router.on(`${appRoutes.adminUsers}/:userId`, params => {
    currentPage = {
      route: AdminUser,
      params: params,
      name: 'admin-user',
    };
  });
  router.on(`${appRoutes.adminUsers}`, () => {
    currentPage = {
      route: AdminUsers,
      params: {},
      name: 'admin-users',
    };
  });
  router.on(`${appRoutes.adminOrganizations}`, () => {
    currentPage = {
      route: AdminOrganizations,
      params: {},
      name: 'admin-organizations',
    };
  });
  router.on(`${appRoutes.adminOrganizations}/:organizationId`, params => {
    currentPage = {
      route: AdminOrganization,
      params: params,
      name: 'admin-organization',
    };
  });
  router.on(
    `${appRoutes.adminOrganizations}/:organizationId/team/:teamId`,
    params => {
      currentPage = {
        route: AdminTeam,
        params: params,
        name: 'admin-team',
      };
    },
  );
  router.on(
    `${appRoutes.adminOrganizations}/:organizationId/department/:departmentId`,
    params => {
      currentPage = {
        route: AdminDepartment,
        params: params,
        name: 'admin-department',
      };
    },
  );
  router.on(
    `${appRoutes.adminOrganizations}/:organizationId/department/:departmentId/team/:teamId`,
    params => {
      currentPage = {
        route: AdminTeam,
        params: params,
        name: 'admin-team',
      };
    },
  );
  router.on(`${appRoutes.adminTeams}`, () => {
    currentPage = {
      route: AdminTeams,
      params: {},
      name: 'admin-teams',
    };
  });
  router.on(`${appRoutes.adminTeams}/:teamId`, params => {
    currentPage = {
      route: AdminTeam,
      params: params,
      name: 'admin-team',
    };
  });
  router.on(`${appRoutes.adminApiKeys}`, () => {
    currentPage = {
      route: AdminApikeys,
      params: {},
      name: 'admin-apikeys',
    };
  });
  router.on(`${appRoutes.adminAlerts}`, () => {
    currentPage = {
      route: AdminAlerts,
      params: {},
      name: 'admin-alerts',
    };
  });

  if (FeaturePoker) {
    router.on(appRoutes.games, () => {
      currentPage = {
        route: Battles,
        params: {},
        name: 'battles',
      };
    });
    router.on(`${appRoutes.game}/:battleId`, params => {
      currentPage = {
        route: Battle,
        params,
        name: 'battle',
      };
    });
    router.on(`${appRoutes.adminPokerGames}`, () => {
      currentPage = {
        route: AdminBattles,
        params: {},
        name: 'admin-games',
      };
    });
    router.on(`${appRoutes.adminPokerGames}/:battleId`, params => {
      currentPage = {
        route: AdminBattle,
        params: params,
        name: 'admin-poker',
      };
    });
    router.on(`${appRoutes.register}/battle/:battleId`, params => {
      currentPage = {
        route: Register,
        params,
        name: 'register',
      };
    });
    router.on(`${appRoutes.login}/battle/:battleId`, params => {
      currentPage = {
        route: Login,
        params,
        name: 'login',
      };
    });
  }

  if (FeatureRetro) {
    router.on(appRoutes.retros, () => {
      currentPage = {
        route: Retros,
        params: {},
        name: 'retros',
      };
    });
    router.on(`${appRoutes.retro}/:retroId`, params => {
      currentPage = {
        route: Retro,
        params,
        name: 'retro',
      };
    });
    router.on(`${appRoutes.adminRetros}`, () => {
      currentPage = {
        route: AdminRetros,
        params: {},
        name: 'admin-retros',
      };
    });
    router.on(`${appRoutes.adminRetros}/:retroId`, params => {
      currentPage = {
        route: AdminRetro,
        params: params,
        name: 'admin-retro',
      };
    });
    router.on(`${appRoutes.register}/retro/:retroId`, params => {
      currentPage = {
        route: Register,
        params,
        name: 'register',
      };
    });
    router.on(`${appRoutes.login}/retro/:retroId`, params => {
      currentPage = {
        route: Login,
        params,
        name: 'login',
      };
    });
  }

  if (FeatureStoryboard) {
    router.on(appRoutes.storyboards, () => {
      currentPage = {
        route: Storyboards,
        params: {},
        name: 'storyboards',
      };
    });
    router.on(`${appRoutes.storyboard}/:storyboardId`, params => {
      currentPage = {
        route: Storyboard,
        params,
        name: 'storyboard',
      };
    });
    router.on(`${appRoutes.adminStoryboards}`, () => {
      currentPage = {
        route: AdminStoryboards,
        params: {},
        name: 'admin-storyboards',
      };
    });
    router.on(`${appRoutes.adminStoryboards}/:storyboardId`, params => {
      currentPage = {
        route: AdminStoryboard,
        params: params,
        name: 'admin-storyboard',
      };
    });
    router.on(`${appRoutes.register}/storyboard/:storyboardId`, params => {
      currentPage = {
        route: Register,
        params,
        name: 'register',
      };
    });
    router.on(`${appRoutes.login}/storyboard/:storyboardId`, params => {
      currentPage = {
        route: Login,
        params,
        name: 'login',
      };
    });
  }

  if (SubscriptionsEnabled) {
    router.on(`${appRoutes.adminSubscriptions}`, () => {
      currentPage = {
        route: AdminSubscriptions,
        params: {},
        name: 'admin-subscriptions',
      };
    });
    router.on(`${appRoutes.adminSubscriptions}/:subscriptionId`, params => {
      currentPage = {
        route: AdminSubscription,
        params,
        name: 'admin-subscription',
      };
    });
  }

  router.listen();

  const xfetch = apiclient(handle401);

  function handle401(skipRedirect) {
    eventTag('session_expired', 'engagement', 'unauthorized', () => {
      user.delete();
      localStorage.removeItem('theme');
      window.setTheme();
      if (!skipRedirect) {
        router.route(appRoutes.login);
      }
    });
  }

  onDestroy(router.unlisten);
</script>

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
