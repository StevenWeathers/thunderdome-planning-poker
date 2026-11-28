<script lang="ts">
  import { onMount } from 'svelte';
  import { ChevronLeft } from '@lucide/svelte';

  import TeamPageLayout from '../../components/team/TeamPageLayout.svelte';
  import SolidButton from '../../components/global/SolidButton.svelte';
  import HollowButton from '../../components/global/HollowButton.svelte';
  import DeleteConfirmation from '../../components/global/DeleteConfirmation.svelte';
  import FeatureSubscribeBanner from '../../components/global/FeatureSubscribeBanner.svelte';
  import HealthCheckForm from '../../components/healthcheck/HealthCheckForm.svelte';
  import HealthCheckResponseForm from '../../components/healthcheck/HealthCheckResponseForm.svelte';
  import HealthCheckResults from '../../components/healthcheck/HealthCheckResults.svelte';
  import { user } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import { AppConfig, appRoutes } from '../../config';
  import { validateUserIsRegistered } from '../../validationUtils';
  import { formatTimestamp } from '../../dateUtils';

  import type { NotificationService } from '../../types/notifications';
  import type { ApiClient } from '../../types/apiclient';
  import type { HealthCheckAnswer, HealthCheckDetail } from '../../types/healthcheck';

  interface Props {
    xfetch: ApiClient;
    router: any;
    notifications: NotificationService;
    organizationId?: string;
    departmentId?: string;
    teamId: string;
    healthCheckId: string;
  }

  let { xfetch, router, notifications, organizationId, departmentId, teamId, healthCheckId }: Props = $props();

  let detail = $state<HealthCheckDetail | null>(null);
  let subscriptionRequired = $state(false);
  let showEdit = $state(false);
  let showClose = $state(false);
  let showDelete = $state(false);

  const teamPrefix = $derived(
    departmentId
      ? `/api/organizations/${organizationId}/departments/${departmentId}/teams/${teamId}`
      : organizationId
        ? `/api/organizations/${organizationId}/teams/${teamId}`
        : `/api/teams/${teamId}`,
  );
  const listPath = $derived(
    departmentId
      ? `${appRoutes.organization}/${organizationId}/department/${departmentId}/team/${teamId}/health-checks`
      : organizationId
        ? `${appRoutes.organization}/${organizationId}/team/${teamId}/health-checks`
        : `${appRoutes.team}/${teamId}/health-checks`,
  );
  const checkUrl = $derived(`${teamPrefix}/health-checks/${healthCheckId}`);
  const healthCheck = $derived(detail?.healthCheck);

  function getHealthCheck() {
    xfetch(checkUrl)
      .then(res => res.json())
      .then(result => {
        detail = result.data;
      })
      .catch(error => {
        const response = error?.[1];
        const refused = () => notifications.danger($LL.healthCheckGetError());
        if (!AppConfig.SubscriptionsEnabled || response?.status !== 403) {
          refused();
          return;
        }
        response
          .json()
          .then((result: { error: string }) => {
            if (result.error === 'TEAM_SUBSCRIPTION_REQUIRED') {
              subscriptionRequired = true;
            } else {
              refused();
            }
          })
          .catch(refused);
      });
  }

  function handleSubmit(answers: HealthCheckAnswer[]) {
    xfetch(`${checkUrl}/response`, { method: 'PUT', body: { answers } })
      .then(() => {
        notifications.success($LL.healthCheckSubmitSuccess());
        getHealthCheck();
      })
      .catch(error => {
        if (error?.[1]?.status === 409) {
          notifications.danger($LL.healthCheckClosedError());
          getHealthCheck();
        } else {
          notifications.danger($LL.healthCheckSubmitError());
        }
      });
  }

  function handleClose() {
    xfetch(`${checkUrl}/close`, { method: 'POST' })
      .then(() => {
        showClose = false;
        getHealthCheck();
      })
      .catch(() => {
        showClose = false;
        notifications.danger($LL.healthCheckCloseError());
        getHealthCheck();
      });
  }

  function handleDelete() {
    xfetch(checkUrl, { method: 'DELETE' })
      .then(() => {
        notifications.success($LL.healthCheckDeleteSuccess());
        router.route(listPath);
      })
      .catch(() => {
        showDelete = false;
        notifications.danger($LL.healthCheckDeleteError());
      });
  }

  function copyLink() {
    // the same link that is emailed, it works for any team member
    const link = `${window.location.origin}${appRoutes.team}/${teamId}/health-checks/${healthCheckId}`;
    Promise.resolve()
      .then(() => navigator.clipboard.writeText(link))
      .then(() => notifications.success($LL.healthCheckLinkCopied()))
      .catch(() => notifications.danger($LL.healthCheckLinkCopyError()));
  }

  function handleSaved() {
    showEdit = false;
    getHealthCheck();
  }

  onMount(() => {
    if (!$user.id || !validateUserIsRegistered($user)) {
      // return to this health check after logging in
      router.route(`${appRoutes.login}/health-check/${teamId}/${healthCheckId}`);
      return;
    }

    getHealthCheck();
  });
</script>

<svelte:head>
  <title>{$LL.healthCheck()} {healthCheck?.name ?? ''} | {$LL.appName()}</title>
</svelte:head>

<TeamPageLayout activePage="health-checks" {teamId} {organizationId} {departmentId}>
  <a class="inline-flex items-center text-blue-600 dark:text-sky-400 mb-4" href={listPath}>
    <ChevronLeft class="w-4 h-4 rtl:rotate-180" aria-hidden="true" />
    {$LL.healthChecks()}
  </a>

  {#if subscriptionRequired}
    <FeatureSubscribeBanner salesPitch={$LL.healthCheckSubscribePitch()} />
  {:else if detail && healthCheck}
    <div class="flex flex-wrap gap-4 items-start justify-between mb-4">
      <div>
        <h1 class="text-3xl font-semibold font-rajdhani dark:text-white">{healthCheck.name}</h1>
        <p class="text-gray-600 dark:text-gray-400">
          {healthCheck.closed ? $LL.healthCheckClosed() : $LL.healthCheckOpen()}
          &middot;
          {$LL.healthCheckDeadline()}: {formatTimestamp(healthCheck.deadline)}
          &middot;
          <span data-testid="health-check-response-count">
            {$LL.healthCheckResponses({ count: healthCheck.responseCount })}
          </span>
        </p>
      </div>
      <div class="flex flex-wrap gap-2">
        <HollowButton onClick={copyLink} color="blue">{$LL.healthCheckCopyLink()}</HollowButton>
        {#if detail.canManage}
          {#if !healthCheck.closed}
            <HollowButton onClick={() => (showEdit = true)} color="blue" testid="health-check-edit">
              {$LL.healthCheckEdit()}
            </HollowButton>
            <SolidButton onClick={() => (showClose = true)} testid="health-check-close">
              {$LL.healthCheckCloseNow()}
            </SolidButton>
          {/if}
          <HollowButton onClick={() => (showDelete = true)} color="red" testid="health-check-delete">
            {$LL.healthCheckDelete()}
          </HollowButton>
        {/if}
      </div>
    </div>

    <p class="mb-6 text-sm text-gray-600 dark:text-gray-400">
      {$LL.healthCheckAnonymousNote({ min: detail.minResponses })}
    </p>

    {#if detail.pending?.length}
      <p class="mb-6 dark:text-gray-300" data-testid="health-check-pending">
        <span class="font-bold">{$LL.healthCheckPending()}:</span>
        {detail.pending.map(p => p.name).join(', ')}
      </p>
    {/if}

    {#if healthCheck.closed}
      <h2 class="text-2xl font-semibold font-rajdhani uppercase mb-4 dark:text-white">{$LL.healthCheckResults()}</h2>
      {#if detail.results}
        <HealthCheckResults results={detail.results} />
      {:else}
        <p class="dark:text-gray-300">{$LL.healthCheckResultsTooFew({ min: detail.minResponses })}</p>
      {/if}
    {:else if detail.canRespond}
      <h2 class="text-2xl font-semibold font-rajdhani uppercase mb-1 dark:text-white">
        {$LL.healthCheckYourResponse()}
      </h2>
      <p class="mb-4 text-gray-600 dark:text-gray-400">{$LL.healthCheckUpdateUntil()}</p>
      {#key healthCheck.id}
        <HealthCheckResponseForm areas={healthCheck.areas} answers={detail.myAnswers} {handleSubmit} />
      {/key}
    {:else}
      <p class="dark:text-gray-300">{$LL.healthCheckViewOnly()}</p>
      <p class="dark:text-gray-300">{$LL.healthCheckResultsPending()}</p>
    {/if}

    {#if showEdit}
      <HealthCheckForm
        {healthCheck}
        apiPrefix={teamPrefix}
        toggleClose={() => (showEdit = false)}
        {handleSaved}
        {xfetch}
        {notifications}
      />
    {/if}

    {#if showClose}
      <DeleteConfirmation
        toggleDelete={() => (showClose = false)}
        handleDelete={handleClose}
        permanent={false}
        confirmText={$LL.healthCheckCloseConfirm()}
        confirmBtnText={$LL.healthCheckCloseNow()}
      />
    {/if}

    {#if showDelete}
      <DeleteConfirmation
        toggleDelete={() => (showDelete = false)}
        {handleDelete}
        permanent={true}
        confirmText={$LL.healthCheckDeleteConfirm()}
        confirmBtnText={$LL.healthCheckDelete()}
      />
    {/if}
  {/if}
</TeamPageLayout>
