<script lang="ts">
  import { onMount } from 'svelte';

  import HollowButton from '../../components/global/HollowButton.svelte';
  import { user } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import { AppConfig, appRoutes } from '../../config';
  import { validateUserIsRegistered } from '../../validationUtils';
  import {
    ChartNoAxesColumn,
    CheckSquare,
    ChevronRight,
    LayoutDashboard,
    Network,
    RefreshCcw,
    SquareDashedKanban,
    User,
    Users,
    Vote,
  } from 'lucide-svelte';
  import DeleteConfirmation from '../../components/global/DeleteConfirmation.svelte';
  import EstimationScalesList from '../../components/estimationscale/EstimationScalesList.svelte';
  import MetricsDisplay from '../../components/global/MetricsDisplay.svelte';
  import { fetchAndUpdateMetrics, type MetricItem } from '../../components/team/metrics';
  import FeatureSubscribeBanner from '../../components/global/FeatureSubscribeBanner.svelte';
  import RetroTemplatesList from '../../components/retrotemplate/RetroTemplatesList.svelte';
  import PokerSettings from '../../components/poker/PokerSettings.svelte';
  import RetroSettings from '../../components/retro/RetroSettings.svelte';

  import type { NotificationService } from '../../types/notifications';
  import type { ApiClient } from '../../types/apiclient';
  import OrgPageLayout from '../../components/organization/OrgPageLayout.svelte';

  interface Props {
    xfetch: ApiClient;
    router: any;
    notifications: NotificationService;
    organizationId: any;
  }

  let { xfetch, router, notifications, organizationId }: Props = $props();

  const orgPrefix = $derived(`/api/organizations/${organizationId}`);

  let organization = $state({
    id: '',
    name: '',
    createdDate: '',
    updateDate: '',
    subscribed: false,
  });
  let role = $state('MEMBER');
  let showDeleteOrganization = $state(false);

  $effect(() => {
    organization.id = organizationId;
  });

  const toggleDeleteOrganization = () => {
    showDeleteOrganization = !showDeleteOrganization;
  };

  function getOrganization() {
    xfetch(orgPrefix)
      .then(res => res.json())
      .then(function (result) {
        organization = result.data.organization;
        role = result.data.role;

        getEstimationScales();
        getRetroTemplates();
      })
      .catch(function () {
        notifications.danger($LL.organizationGetError());
      });
  }

  let organizationMetrics: MetricItem[] = $state([
    {
      key: 'department_count',
      name: 'Departments',
      value: 0,
      icon: Network,
    },
    { key: 'team_count', name: 'Team Count', value: 0, icon: Users },
    {
      key: 'team_checkin_count',
      name: 'Team Check-ins',
      value: 0,
      icon: CheckSquare,
    },
    { key: 'user_count', name: 'Users', value: 0, icon: User },
    { key: 'poker_count', name: 'Poker Games', value: 0, icon: Vote },
    { key: 'retro_count', name: 'Retros', value: 0, icon: RefreshCcw },
    {
      key: 'storyboard_count',
      name: 'Storyboards',
      value: 0,
      icon: LayoutDashboard,
    },
    {
      key: 'estimation_scale_count',
      name: 'Custom Estimation Scales',
      value: 0,
      icon: ChartNoAxesColumn,
    },
    {
      key: 'retro_template_count',
      name: 'Retro Templates',
      value: 0,
      icon: SquareDashedKanban,
    },
  ]);

  function handleDeleteOrganization() {
    xfetch(`${orgPrefix}`, {
      method: 'DELETE',
    })
      .then(function () {
        toggleDeleteOrganization();
        notifications.success($LL.organizationDeleteSuccess());
        router.route(appRoutes.teams);
      })
      .catch(function () {
        notifications.danger($LL.organizationDeleteError());
      });
  }

  const scalesPageLimit = 20;
  let estimationScales = $state([]);
  let scaleCount = 0;
  let scalesPage = $state(1);

  const changeScalesPage = evt => {
    scalesPage = evt.detail;
    getEstimationScales();
  };

  function getEstimationScales() {
    const scalesOffset = (scalesPage - 1) * scalesPageLimit;
    if (
      AppConfig.FeaturePoker &&
      (!AppConfig.SubscriptionsEnabled || (AppConfig.SubscriptionsEnabled && organization.subscribed))
    ) {
      xfetch(`${orgPrefix}/estimation-scales?limit=${scalesPageLimit}&offset=${scalesOffset}`)
        .then(res => res.json())
        .then(function (result) {
          estimationScales = result.data;
        })
        .catch(function () {
          notifications.danger('Failed to get estimation scales');
        });
    }
  }

  const retroTemplatePageLimit = 20;
  let retroTemplates = $state([]);
  let retroTemplateCount = 0;
  let retroTemplatesPage = $state(1);

  const changeRetroTemplatesPage = evt => {
    retroTemplatesPage = evt.detail;
    getRetroTemplates();
  };

  function getRetroTemplates() {
    const offset = (retroTemplatesPage - 1) * retroTemplatePageLimit;
    if (
      AppConfig.FeatureRetro &&
      (!AppConfig.SubscriptionsEnabled || (AppConfig.SubscriptionsEnabled && organization.subscribed))
    ) {
      xfetch(`${orgPrefix}/retro-templates?limit=${retroTemplatePageLimit}&offset=${offset}`)
        .then(res => res.json())
        .then(function (result) {
          retroTemplates = result.data;
        })
        .catch(function () {
          notifications.danger('Failed to get retro templates');
        });
    }
  }

  onMount(async () => {
    if (!$user.id || !validateUserIsRegistered($user)) {
      router.route(appRoutes.login);
      return;
    }

    getOrganization();
    try {
      organizationMetrics = await fetchAndUpdateMetrics(orgPrefix, organizationMetrics);
    } catch (e) {
      notifications.danger('Failed to get organization metrics');
    }
  });

  let isAdmin = $derived(role === 'ADMIN');
</script>

<svelte:head>
  <title>{$LL.organization()} {organization.name} | {$LL.appName()}</title>
</svelte:head>

<OrgPageLayout activePage="organization" {organizationId}>
  <h1 class="mb-4 text-3xl font-semibold font-rajdhani dark:text-white">
    <span class="uppercase">{$LL.organization()}</span>
    <ChevronRight class="w-8 h-8 inline-block" />
    {organization.name}
  </h1>

  <div class="mb-8">
    <MetricsDisplay metrics={organizationMetrics} />
  </div>

  {#if AppConfig.FeaturePoker}
    <div class="mt-8">
      {#if !AppConfig.SubscriptionsEnabled || (AppConfig.SubscriptionsEnabled && organization.subscribed)}
        <PokerSettings {xfetch} {notifications} isEntityAdmin={isAdmin} apiPrefix={orgPrefix} {organizationId} />
      {:else}
        <FeatureSubscribeBanner
          salesPitch="Optimize your Organization's estimation workflow with customized default Planning Poker settings."
        />
      {/if}
    </div>

    <div class="mt-8">
      {#if !AppConfig.SubscriptionsEnabled || (AppConfig.SubscriptionsEnabled && organization.subscribed)}
        <EstimationScalesList
          {xfetch}
          {notifications}
          isEntityAdmin={isAdmin}
          apiPrefix={orgPrefix}
          {organizationId}
          scales={estimationScales}
          getScales={getEstimationScales}
          {scaleCount}
          {scalesPage}
          {scalesPageLimit}
          changePage={changeScalesPage}
        />
      {:else}
        <FeatureSubscribeBanner
          salesPitch="Create custom poker point scales to match your Organization's estimation style."
        />
      {/if}
    </div>
  {/if}

  {#if AppConfig.FeatureRetro}
    <div class="mt-8">
      {#if !AppConfig.SubscriptionsEnabled || (AppConfig.SubscriptionsEnabled && organization.subscribed)}
        <RetroSettings {xfetch} {notifications} isEntityAdmin={isAdmin} apiPrefix={orgPrefix} {organizationId} />
      {:else}
        <FeatureSubscribeBanner
          salesPitch="Enhance your Organization's reflection process with customized default Retrospective settings."
        />
      {/if}
    </div>

    <div class="mt-8">
      {#if !AppConfig.SubscriptionsEnabled || (AppConfig.SubscriptionsEnabled && organization.subscribed)}
        <RetroTemplatesList
          {xfetch}
          {notifications}
          isEntityAdmin={isAdmin}
          apiPrefix={orgPrefix}
          {organizationId}
          templates={retroTemplates}
          getTemplates={getRetroTemplates}
          templateCount={retroTemplateCount}
          templatesPage={retroTemplatesPage}
          templatesPageLimit={retroTemplatePageLimit}
          changePage={changeRetroTemplatesPage}
        />
      {:else}
        <FeatureSubscribeBanner
          salesPitch="Tailor your Organization's reflection process with custom retrospective templates."
        />
      {/if}
    </div>
  {/if}

  {#if isAdmin}
    <div class="w-full text-center mt-8">
      <HollowButton onClick={toggleDeleteOrganization} color="red">
        {$LL.deleteOrganization()}
      </HollowButton>
    </div>
  {/if}

  {#if showDeleteOrganization}
    <DeleteConfirmation
      toggleDelete={toggleDeleteOrganization}
      handleDelete={handleDeleteOrganization}
      confirmText={$LL.deleteOrganizationConfirmText()}
      confirmBtnText={$LL.deleteOrganization()}
    />
  {/if}
</OrgPageLayout>
