<script lang="ts">
  import { onMount } from 'svelte';
  import { ExternalLink } from '@lucide/svelte';

  import TeamPageLayout from '../../components/team/TeamPageLayout.svelte';
  import FeatureSubscribeBanner from '../../components/global/FeatureSubscribeBanner.svelte';
  import SolidButton from '../../components/global/SolidButton.svelte';
  import Pagination from '../../components/global/Pagination.svelte';
  import HealthCheckForm from '../../components/healthcheck/HealthCheckForm.svelte';
  import HealthCheckRatingGuide from '../../components/healthcheck/HealthCheckRatingGuide.svelte';
  import { user } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import { AppConfig, appRoutes } from '../../config';
  import { validateUserIsRegistered } from '../../validationUtils';
  import { formatTimestamp } from '../../dateUtils';

  import type { NotificationService } from '../../types/notifications';
  import type { ApiClient } from '../../types/apiclient';
  import type { HealthCheck } from '../../types/healthcheck';

  interface Props {
    xfetch: ApiClient;
    router: any;
    notifications: NotificationService;
    organizationId?: string;
    departmentId?: string;
    teamId: string;
  }

  let { xfetch, router, notifications, organizationId, departmentId, teamId }: Props = $props();

  const pageLimit = 20;
  const spotifyArticleUrl = 'https://engineering.atspotify.com/2014/9/squad-health-check-model';

  let teamName = $state('');
  let isTeamMember = $state(false);
  let healthChecks = $state<HealthCheck[]>([]);
  let healthCheckCount = $state(0);
  let currentPage = $state(1);
  let loaded = $state(false);
  let subscriptionRequired = $state(false);
  let showCreate = $state(false);

  const teamPrefix = $derived(
    departmentId
      ? `/api/organizations/${organizationId}/departments/${departmentId}/teams/${teamId}`
      : organizationId
        ? `/api/organizations/${organizationId}/teams/${teamId}`
        : `/api/teams/${teamId}`,
  );
  const pagePath = $derived(
    departmentId
      ? `${appRoutes.organization}/${organizationId}/department/${departmentId}/team/${teamId}/health-checks`
      : organizationId
        ? `${appRoutes.organization}/${organizationId}/team/${teamId}/health-checks`
        : `${appRoutes.team}/${teamId}/health-checks`,
  );

  function getTeam() {
    xfetch(teamPrefix)
      .then(res => res.json())
      .then(result => {
        teamName = result.data.team.name;
        isTeamMember = !!result.data.teamRole;
      })
      .catch(() => {
        notifications.danger($LL.teamGetError());
      });
  }

  function getHealthChecks() {
    const offset = (currentPage - 1) * pageLimit;
    xfetch(`${teamPrefix}/health-checks?limit=${pageLimit}&offset=${offset}`)
      .then(res => res.json())
      .then(result => {
        healthChecks = result.data;
        healthCheckCount = result.meta.count;
        loaded = true;
      })
      .catch(error => {
        // the team page itself loaded, so being refused here means there's no subscription
        if (AppConfig.SubscriptionsEnabled && error?.[1]?.status === 403) {
          subscriptionRequired = true;
          return;
        }
        notifications.danger($LL.healthChecksGetError());
      });
  }

  function changePage(evt: CustomEvent) {
    currentPage = evt.detail;
    getHealthChecks();
  }

  function handleCreated(healthCheck: HealthCheck | null) {
    showCreate = false;
    if (healthCheck) {
      router.route(`${pagePath}/${healthCheck.id}`);
    }
  }

  onMount(() => {
    if (!$user.id || !validateUserIsRegistered($user)) {
      router.route(appRoutes.login);
      return;
    }

    getTeam();
    getHealthChecks();
  });
</script>

<svelte:head>
  <title>{$LL.healthChecks()} {teamName} | {$LL.appName()}</title>
</svelte:head>

<TeamPageLayout activePage="health-checks" {teamId} {organizationId} {departmentId}>
  <div class="flex flex-wrap gap-4 items-center justify-between mb-6">
    <div>
      <h1 class="text-3xl font-semibold font-rajdhani uppercase dark:text-white">{$LL.healthChecks()}</h1>
      <p class="text-lg text-gray-600 dark:text-gray-400">{teamName}</p>
    </div>
    {#if isTeamMember && !subscriptionRequired}
      <SolidButton onClick={() => (showCreate = true)} testid="health-check-create">
        {$LL.healthCheckCreate()}
      </SolidButton>
    {/if}
  </div>

  <section
    class="mb-6 rounded-lg border border-blue-200 dark:border-gray-700 bg-blue-50 dark:bg-gray-800 p-4 text-gray-700 dark:text-gray-300"
    aria-labelledby="health-check-about"
  >
    <h2 id="health-check-about" class="text-xl font-semibold font-rajdhani mb-2 dark:text-white">
      {$LL.healthCheckAboutTitle()}
    </h2>
    <p class="mb-4">{$LL.healthCheckAboutIntro()}</p>
    <div class="grid gap-4 lg:grid-cols-2 mb-4">
      <HealthCheckRatingGuide />
      <div>
        <h3 class="font-bold mb-1">{$LL.healthCheckAboutTips()}</h3>
        <ul class="list-disc ps-5 space-y-1">
          <li>{$LL.healthCheckAboutTipDiscuss()}</li>
          <li>{$LL.healthCheckAboutTipSupport()}</li>
          <li>{$LL.healthCheckAboutTipCadence()}</li>
        </ul>
      </div>
    </div>
    <a
      class="inline-flex items-center gap-1 text-blue-600 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-300 underline"
      href={spotifyArticleUrl}
      target="_blank"
      rel="noopener noreferrer"
    >
      {$LL.healthCheckAboutLink()}
      <ExternalLink class="w-4 h-4" aria-hidden="true" />
      <span class="sr-only">{$LL.healthCheckOpensNewTab()}</span>
    </a>
  </section>

  {#if subscriptionRequired}
    <FeatureSubscribeBanner salesPitch={$LL.healthCheckSubscribePitch()} />
  {:else if loaded && healthChecks.length === 0}
    <p class="text-gray-600 dark:text-gray-400">{$LL.healthChecksEmpty()}</p>
  {:else}
    <ul class="space-y-3">
      {#each healthChecks as healthCheck (healthCheck.id)}
        <li
          class="flex flex-wrap gap-4 items-center justify-between rounded-lg bg-white dark:bg-gray-800 shadow p-4"
          data-testid="health-check"
        >
          <div>
            <a
              class="text-xl font-semibold text-blue-600 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-300"
              href="{pagePath}/{healthCheck.id}"
              data-testid="health-check-name"
            >
              {healthCheck.name}
            </a>
            <p class="text-sm text-gray-600 dark:text-gray-400">
              {$LL.healthCheckDeadline()}: {formatTimestamp(healthCheck.deadline)}
              &middot;
              {$LL.healthCheckResponses({ count: healthCheck.responseCount })}
            </p>
          </div>
          <span
            class="rounded-full px-3 py-1 text-sm font-semibold {healthCheck.closed
              ? 'bg-gray-200 text-gray-700 dark:bg-gray-700 dark:text-gray-300'
              : 'bg-green-100 text-green-800 dark:bg-lime-900 dark:text-lime-200'}"
          >
            {healthCheck.closed ? $LL.healthCheckClosed() : $LL.healthCheckOpen()}
          </span>
        </li>
      {/each}
    </ul>
    {#if healthCheckCount > pageLimit}
      <div class="mt-4 flex justify-center">
        <Pagination bind:current={currentPage} num_items={healthCheckCount} per_page={pageLimit} on:navigate={changePage} />
      </div>
    {/if}
  {/if}

  {#if showCreate}
    <HealthCheckForm
      apiPrefix={teamPrefix}
      toggleClose={() => (showCreate = false)}
      handleSaved={handleCreated}
      {xfetch}
      {notifications}
    />
  {/if}
</TeamPageLayout>
