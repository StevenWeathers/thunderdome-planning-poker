<script lang="ts">
  import './app.css';
  import './tailwind.css';
  import './unreset.css';
  import Navaid from './lib/navaid';
  import { onDestroy } from 'svelte';

  import { AppConfig, appRoutes } from './config';
  import apiclient from './apiclient';
  import type { ApiClient } from './types/apiclient';
  import type { SessionUser } from './types/user';
  import type { NotificationService } from './types/notifications';
  import type { Locales } from './i18n/i18n-types';
  import type { ModuleLoader, Page } from './pageLoader';
  import { dir, user } from './stores';

  import Notifications from './components/global/Notifications.svelte';
  import GlobalHeader from './components/global/GlobalHeader.svelte';
  import GlobalAlerts from './components/alert/GlobalAlerts.svelte';
  import GlobalFooter from './components/global/GlobalFooter.svelte';
  import FullpageLoader from './components/global/FullpageLoader.svelte';
  import SolidButton from './components/global/SolidButton.svelte';
  import LL, { setLocale } from './i18n/i18n-svelte';
  import { detectLocale, isLocale } from './i18n/i18n-util';
  import { loadLocaleAsync } from './i18n/i18n-util.async';
  import { createPageLoader, reloadOnce } from './pageLoader';

  const {
    FeaturePoker,
    FeatureRetro,
    FeatureStoryboard,
    FeatureProject,
    FeatureHealthCheck,
    SubscriptionsEnabled,
    DefaultLocale,
  } = AppConfig;

  let notifications: NotificationService | undefined = $state();

  const detectedLocale = detectLocale() || DefaultLocale;

  // translations are blank until a locale has loaded, so nothing renders before then
  let localeReady = $state(false);
  let localeRequest = 0;

  async function applyLocale(requested: string) {
    const request = ++localeRequest;
    let locale: Locales = isLocale(requested) ? requested : 'en';

    try {
      try {
        await loadLocaleAsync(locale);
      } catch (error) {
        console.error(error);
        locale = 'en';
        await loadLocaleAsync(locale);
      }
    } catch (error) {
      console.error(error);
      // with no translations at all, reloading is the only way to recover
      if (!localeReady) reloadOnce('localeLoadReloaded');
      return;
    }

    if (request !== localeRequest) return;
    setLocale(locale);
    localeReady = true;
  }

  let activeWarrior: SessionUser | undefined = $state();
  user.subscribe((w: any) => {
    activeWarrior = w;

    applyLocale(w?.locale || detectedLocale);
  });

  $effect(() => {
    if (document.dir !== $dir) {
      document.dir = $dir;
    }
  });

  let currentPage: Page | undefined = $state.raw();
  let pagePending = $state(false);
  let pageFailed = $state(false);
  let pageRequested = false;

  const loadPage = createPageLoader({
    onPage: page => (currentPage = page),
    onPending: pending => (pagePending = pending),
    onError: failed => (pageFailed = failed),
  });

  function showPage(load: ModuleLoader, params: any, name: string) {
    pageRequested = true;
    loadPage(load, params, name);
  }

  function showLanding() {
    showPage(() => import('./pages/Landing.svelte'), {}, 'landing');
  }

  // unmatched URLs show the landing page on first load and otherwise keep the current page
  const router = Navaid('/', () => {
    if (!pageRequested) showLanding();
  });

  router.on(appRoutes.landing, showLanding);

  router.on(appRoutes.dashboard, () => {
    showPage(() => import('./pages/Dashboard.svelte'), {}, 'dashboard');
  });
  router.on(appRoutes.privacyPolicy, () => {
    showPage(() => import('./pages/support/PrivacyPolicy.svelte'), {}, 'privacy');
  });
  router.on(appRoutes.termsConditions, () => {
    showPage(() => import('./pages/support/TermsConditions.svelte'), {}, 'terms-conditions');
  });
  router.on(appRoutes.support, () => {
    showPage(() => import('./pages/support/Support.svelte'), {}, 'support');
  });
  router.on(appRoutes.openSource, () => {
    showPage(() => import('./pages/support/OpenSource.svelte'), {}, 'open-source');
  });
  router.on(`${appRoutes.register}`, params => {
    showPage(() => import('./pages/user/Register.svelte'), params, 'register');
  });
  router.on(`${appRoutes.login}`, params => {
    showPage(() => import('./pages/user/Login.svelte'), params, 'login');
  });
  router.on(`${appRoutes.login}/subscription`, () => {
    showPage(() => import('./pages/user/Login.svelte'), { subscription: true }, 'login');
  });
  router.on(`${appRoutes.login}/team/:teamInviteId`, params => {
    showPage(() => import('./pages/user/Login.svelte'), params, 'login');
  });
  router.on(`${appRoutes.login}/organization/:orgInviteId`, params => {
    showPage(() => import('./pages/user/Login.svelte'), params, 'login');
  });
  router.on(`${appRoutes.resetPwd}/:resetId`, params => {
    showPage(() => import('./pages/user/ResetPassword.svelte'), params, 'reset-password');
  });
  router.on(`${appRoutes.verifyAct}/:verifyId`, params => {
    showPage(() => import('./pages/user/VerifyAccount.svelte'), params, 'verify-account');
  });
  router.on(appRoutes.profile, params => {
    showPage(() => import('./pages/user/UserProfile.svelte'), params, 'profile');
  });
  router.on(`${appRoutes.profile}/change-email/:changeId`, params => {
    showPage(() => import('./pages/user/ChangeEmail.svelte'), params, 'change-email');
  });
  router.on(appRoutes.teams, () => {
    showPage(() => import('./pages/team/Teams.svelte'), {}, 'teams');
  });
  router.on(appRoutes.subscriptionPricing, () => {
    showPage(() => import('./pages/subscription/Pricing.svelte'), {}, 'pricing');
  });
  router.on(appRoutes.subscriptionConfirmation, () => {
    showPage(() => import('./pages/subscription/Confirmation.svelte'), {}, 'subscription-confirmation');
  });
  router.on(`${appRoutes.register}/subscription`, () => {
    showPage(() => import('./pages/user/Register.svelte'), { subscription: true }, 'register');
  });
  router.on(`${appRoutes.register}/team/:teamInviteId`, params => {
    showPage(() => import('./pages/user/Register.svelte'), params, 'register');
  });
  router.on(`${appRoutes.register}/organization/:orgInviteId`, params => {
    showPage(() => import('./pages/user/Register.svelte'), params, 'register');
  });
  router.on(`${appRoutes.invite}/organization/:inviteId`, params => {
    showPage(() => import('./pages/team/Invite.svelte'), { inviteType: 'organization', ...params }, 'invite');
  });
  router.on(`${appRoutes.invite}/department/:inviteId`, params => {
    showPage(() => import('./pages/team/Invite.svelte'), { inviteType: 'department', ...params }, 'invite');
  });
  router.on(`${appRoutes.invite}/team/:inviteId`, params => {
    showPage(() => import('./pages/team/Invite.svelte'), { inviteType: 'team', ...params }, 'invite');
  });
  router.on(`${appRoutes.organization}/:organizationId`, params => {
    showPage(() => import('./pages/organization/Organization.svelte'), params, 'organization');
  });
  router.on(`${appRoutes.organization}/:organizationId/users`, params => {
    showPage(() => import('./pages/organization/OrganizationUsers.svelte'), params, 'organization-users');
  });
  router.on(`${appRoutes.organization}/:organizationId/teams`, params => {
    showPage(() => import('./pages/organization/OrganizationTeams.svelte'), params, 'organization-teams');
  });
  router.on(`${appRoutes.organization}/:organizationId/team/:teamId`, params => {
    showPage(() => import('./pages/team/Team.svelte'), params, 'team');
  });
  router.on(`${appRoutes.organization}/:organizationId/team/:teamId/users`, params => {
    showPage(() => import('./pages/team/TeamUsers.svelte'), params, 'team-users');
  });
  router.on(`${appRoutes.organization}/:organizationId/team/:teamId/checkin`, params => {
    showPage(() => import('./pages/team/TeamCheckin.svelte'), params, 'team-checkin');
  });
  router.on(`${appRoutes.organization}/:organizationId/departments`, params => {
    showPage(() => import('./pages/organization/OrganizationDepartments.svelte'), params, 'organization-departments');
  });
  router.on(`${appRoutes.organization}/:organizationId/department/:departmentId`, params => {
    showPage(() => import('./pages/department/Department.svelte'), params, 'department');
  });
  router.on(`${appRoutes.organization}/:organizationId/department/:departmentId/team/:teamId`, params => {
    showPage(() => import('./pages/team/Team.svelte'), params, 'team');
  });
  router.on(`${appRoutes.organization}/:organizationId/department/:departmentId/team/:teamId/users`, params => {
    showPage(() => import('./pages/team/TeamUsers.svelte'), params, 'team-users');
  });
  router.on(`${appRoutes.organization}/:organizationId/department/:departmentId/team/:teamId/checkin`, params => {
    showPage(() => import('./pages/team/TeamCheckin.svelte'), params, 'team-checkin');
  });
  if (FeatureProject) {
    router.on(`${appRoutes.organization}/:organizationId/department/:departmentId/team/:teamId/projects`, params => {
      showPage(() => import('./pages/team/TeamProjects.svelte'), params, 'team-projects');
    });
    router.on(`${appRoutes.organization}/:organizationId/team/:teamId/projects`, params => {
      showPage(() => import('./pages/team/TeamProjects.svelte'), params, 'team-projects');
    });
  }
  router.on(`${appRoutes.team}/:teamId`, params => {
    showPage(() => import('./pages/team/Team.svelte'), params, 'team');
  });
  router.on(`${appRoutes.team}/:teamId/users`, params => {
    showPage(() => import('./pages/team/TeamUsers.svelte'), params, 'team-users');
  });
  if (FeatureProject) {
    router.on(`${appRoutes.team}/:teamId/projects`, params => {
      showPage(() => import('./pages/team/TeamProjects.svelte'), params, 'team-projects');
    });
  }
  router.on(`${appRoutes.team}/:teamId/checkin`, params => {
    showPage(() => import('./pages/team/TeamCheckin.svelte'), params, 'team-checkin');
  });
  router.on(appRoutes.admin, () => {
    showPage(() => import('./pages/admin/Admin.svelte'), {}, 'admin');
  });
  router.on(appRoutes.adminSupportTickets, () => {
    showPage(() => import('./pages/admin/support/ListTickets.svelte'), {}, 'support-tickets');
  });
  router.on(`${appRoutes.adminSupportTickets}/:ticketId`, params => {
    showPage(() => import('./pages/admin/support/TicketDetails.svelte'), params, 'support-tickets');
  });
  router.on(`${appRoutes.adminUsers}/:userId`, params => {
    showPage(() => import('./pages/admin/user/User.svelte'), params, 'admin-user');
  });
  router.on(`${appRoutes.adminUsers}`, () => {
    showPage(() => import('./pages/admin/user/Users.svelte'), {}, 'admin-users');
  });
  router.on(`${appRoutes.adminOrganizations}`, () => {
    showPage(() => import('./pages/admin/team/Organizations.svelte'), {}, 'admin-organizations');
  });
  router.on(`${appRoutes.adminOrganizations}/:organizationId`, params => {
    showPage(() => import('./pages/admin/team/Organization.svelte'), params, 'admin-organization');
  });
  router.on(`${appRoutes.adminOrganizations}/:organizationId/team/:teamId`, params => {
    showPage(() => import('./pages/admin/team/Team.svelte'), params, 'admin-team');
  });
  router.on(`${appRoutes.adminOrganizations}/:organizationId/department/:departmentId`, params => {
    showPage(() => import('./pages/admin/team/Department.svelte'), params, 'admin-department');
  });
  router.on(`${appRoutes.adminOrganizations}/:organizationId/department/:departmentId/team/:teamId`, params => {
    showPage(() => import('./pages/admin/team/Team.svelte'), params, 'admin-team');
  });
  router.on(`${appRoutes.adminTeams}`, () => {
    showPage(() => import('./pages/admin/team/Teams.svelte'), {}, 'admin-teams');
  });
  router.on(`${appRoutes.adminTeams}/:teamId`, params => {
    showPage(() => import('./pages/admin/team/Team.svelte'), params, 'admin-team');
  });
  router.on(`${appRoutes.adminApiKeys}`, () => {
    showPage(() => import('./pages/admin/user/ApiKeys.svelte'), {}, 'admin-apikeys');
  });
  router.on(`${appRoutes.adminAlerts}`, () => {
    showPage(() => import('./pages/admin/Alerts.svelte'), {}, 'admin-alerts');
  });

  if (FeaturePoker) {
    router.on(appRoutes.battles, () => {
      router.route(appRoutes.games);
    });
    router.on(`${appRoutes.battle}/:battleId`, params => {
      router.route(`${appRoutes.game}/${params!.battleId}`);
    });
    router.on(appRoutes.games, () => {
      showPage(
        $user.id
          ? () => import('./pages/poker/PokerGames.svelte')
          : () => import('./pages/poker/PokerLanding.svelte'),
        {},
        'battles',
      );
    });
    router.on(`${appRoutes.game}/:battleId`, params => {
      showPage(() => import('./pages/poker/PokerGame.svelte'), params, 'battle');
    });
    router.on(`${appRoutes.adminPokerGames}`, () => {
      showPage(() => import('./pages/admin/poker/PokerGames.svelte'), {}, 'admin-games');
    });
    router.on(`${appRoutes.adminPokerGames}/:battleId`, params => {
      showPage(() => import('./pages/admin/poker/PokerGame.svelte'), params, 'admin-poker');
    });
    router.on(`${appRoutes.register}/battle/:battleId`, params => {
      router.route(`${appRoutes.register}/game/${params!.battleId}`);
    });
    router.on(`${appRoutes.login}/battle/:battleId`, params => {
      router.route(`${appRoutes.login}/game/${params!.battleId}`);
    });
    router.on(`${appRoutes.register}/game/:battleId`, params => {
      showPage(() => import('./pages/user/Register.svelte'), params, 'register');
    });
    router.on(`${appRoutes.login}/game/:battleId`, params => {
      showPage(() => import('./pages/user/Login.svelte'), params, 'login');
    });
    router.on(`${appRoutes.adminEstimationScales}`, () => {
      showPage(() => import('./pages/admin/poker/EstimationScales.svelte'), {}, 'admin-estimation-scales');
    });
  }

  if (FeatureRetro) {
    router.on(appRoutes.retros, () => {
      showPage(
        $user.id
          ? () => import('./pages/retro/Retros.svelte')
          : () => import('./pages/retro/RetroLanding.svelte'),
        {},
        'retros',
      );
    });
    router.on(`${appRoutes.retro}/:retroId`, params => {
      showPage(() => import('./pages/retro/Retro.svelte'), params, 'retro');
    });
    router.on(`${appRoutes.adminRetros}`, () => {
      showPage(() => import('./pages/admin/retro/Retros.svelte'), {}, 'admin-retros');
    });
    router.on(`${appRoutes.adminRetros}/:retroId`, params => {
      showPage(() => import('./pages/admin/retro/Retro.svelte'), params, 'admin-retro');
    });
    router.on(`${appRoutes.adminPersonaTemplates}`, () => {
      showPage(() => import('./pages/admin/storyboard/PersonaTemplates.svelte'), {}, 'admin-persona-templates');
    });
    router.on(`${appRoutes.adminRetroTags}`, () => {
      showPage(() => import('./pages/admin/retro/RetroTags.svelte'), {}, 'admin-retro-tags');
    });
    router.on(`${appRoutes.adminRetroTemplates}`, () => {
      showPage(() => import('./pages/admin/retro/RetroTemplates.svelte'), {}, 'admin-retro-templates');
    });
    router.on(`${appRoutes.register}/retro/:retroId`, params => {
      showPage(() => import('./pages/user/Register.svelte'), params, 'register');
    });
    router.on(`${appRoutes.login}/retro/:retroId`, params => {
      showPage(() => import('./pages/user/Login.svelte'), params, 'login');
    });
  }

  if (FeatureStoryboard) {
    router.on(appRoutes.storyboards, () => {
      showPage(
        $user.id
          ? () => import('./pages/storyboard/Storyboards.svelte')
          : () => import('./pages/storyboard/StoryboardLanding.svelte'),
        {},
        'storyboards',
      );
    });
    router.on(`${appRoutes.storyboard}/:storyboardId`, params => {
      showPage(() => import('./pages/storyboard/Storyboard.svelte'), params, 'storyboard');
    });
    router.on(`${appRoutes.adminStoryboards}`, () => {
      showPage(() => import('./pages/admin/storyboard/Storyboards.svelte'), {}, 'admin-storyboards');
    });
    router.on(`${appRoutes.adminStoryboards}/:storyboardId`, params => {
      showPage(() => import('./pages/admin/storyboard/Storyboard.svelte'), params, 'admin-storyboard');
    });
    router.on(`${appRoutes.register}/storyboard/:storyboardId`, params => {
      showPage(() => import('./pages/user/Register.svelte'), params, 'register');
    });
    router.on(`${appRoutes.login}/storyboard/:storyboardId`, params => {
      showPage(() => import('./pages/user/Login.svelte'), params, 'login');
    });
  }

  if (FeatureProject) {
    router.on(`${appRoutes.adminProjects}`, () => {
      showPage(() => import('./pages/admin/project/Projects.svelte'), {}, 'admin-projects');
    });
    router.on(`${appRoutes.organization}/:organizationId/projects`, params => {
      showPage(() => import('./pages/organization/OrganizationProjects.svelte'), params, 'organization-projects');
    });
    router.on(`${appRoutes.projects}/:projectId`, params => {
      showPage(() => import('./pages/project/Project.svelte'), params, 'project');
    });
  }

  if (FeatureHealthCheck) {
    const teamRoutes = [
      `${appRoutes.team}/:teamId`,
      `${appRoutes.organization}/:organizationId/team/:teamId`,
      `${appRoutes.organization}/:organizationId/department/:departmentId/team/:teamId`,
    ];
    teamRoutes.forEach(teamRoute => {
      router.on(`${teamRoute}/health-checks`, params => {
        showPage(() => import('./pages/team/TeamHealthChecks.svelte'), params, 'team-health-checks');
      });
      router.on(`${teamRoute}/health-checks/:healthCheckId`, params => {
        showPage(() => import('./pages/team/TeamHealthCheck.svelte'), params, 'team-health-check');
      });
    });
    router.on(`${appRoutes.login}/health-check/:healthCheckTeamId/:healthCheckId`, params => {
      showPage(() => import('./pages/user/Login.svelte'), params, 'login');
    });
    router.on(appRoutes.adminHealthCheckTemplates, () => {
      showPage(
        () => import('./pages/admin/healthcheck/HealthCheckTemplates.svelte'),
        {},
        'admin-health-check-templates',
      );
    });
  }

  if (SubscriptionsEnabled) {
    router.on(`${appRoutes.adminSubscriptions}`, () => {
      showPage(() => import('./pages/admin/Subscriptions.svelte'), {}, 'admin-subscriptions');
    });
    router.on(`${appRoutes.adminSubscriptions}/:subscriptionId`, params => {
      showPage(() => import('./pages/admin/Subscription.svelte'), params, 'admin-subscription');
    });
  }

  router.listen();

  const xfetch: ApiClient = apiclient(handle401);

  function handle401(skipRedirect: boolean = false) {
    user.delete();
    localStorage.removeItem('theme');
    (window as any).setTheme();
    if (!skipRedirect) {
      router.route(appRoutes.login);
    }
  }

  onDestroy(() => {
    if (router.unlisten) {
      router.unlisten();
    }
  });
</script>

{#if localeReady}
  <Notifications bind:this={notifications} />

  <header class="w-full">
    <GlobalAlerts registered={!!activeWarrior?.name || false} />

    <GlobalHeader {router} {xfetch} {notifications} currentPage={currentPage?.name ?? ''} />
  </header>

  <main class="flex-grow flex flex-wrap flex-col">
    {#if pageFailed}
      <div class="w-full text-center p-8 text-gray-800 dark:text-gray-200" role="alert">
        <p class="mb-4">{$LL.pageLoadError()}</p>
        <SolidButton onClick={() => window.location.reload()}>{$LL.pageLoadReload()}</SolidButton>
      </div>
    {:else if currentPage}
      <currentPage.route {...currentPage.params} {notifications} {router} {xfetch} />
    {/if}
  </main>

  <GlobalFooter />
{/if}

{#if !localeReady || pagePending || (!currentPage && !pageFailed)}
  <FullpageLoader />
{/if}
