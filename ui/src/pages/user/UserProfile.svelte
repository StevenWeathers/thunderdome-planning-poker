<script lang="ts">
  import { onMount } from 'svelte';

  import PageLayout from '../../components/PageLayout.svelte';
  import UpdatePasswordForm from '../../components/user/UpdatePasswordForm.svelte';
  import HollowButton from '../../components/global/HollowButton.svelte';
  import { user } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import { AppConfig, appRoutes } from '../../config';
  import ProfileForm from '../../components/user/ProfileForm.svelte';
  import CreateApiKey from '../../components/user/CreateApiKey.svelte';
  import DeleteConfirmation from '../../components/global/DeleteConfirmation.svelte';
  import CreateJiraInstance from '../../components/jira/CreateJiraInstance.svelte';
  import SolidButton from '../../components/global/SolidButton.svelte';
  import UserSubscriptionsList from '../../components/subscription/UserSubscriptionsList.svelte';
  import BooleanDisplay from '../../components/global/BooleanDisplay.svelte';
  import FeatureSubscribeBanner from '../../components/global/FeatureSubscribeBanner.svelte';

  import type { NotificationService } from '../../types/notifications';
  import type { ApiClient } from '../../types/apiclient';

  interface ApiKey {
    id: string;
    name: string;
    prefix: string;
    active: boolean;
    updatedDate: string;
  }

  interface JiraInstance {
    id: string;
    host: string;
    client_mail: string;
  }

  interface Props {
    xfetch: ApiClient;
    router: any;
    notifications: NotificationService;
  }

  let { xfetch, router, notifications }: Props = $props();

  let userProfile = $state<any>({});
  let userCredential = $state(null);
  let apiKeys = $state<ApiKey[]>([]);
  let jiraInstances = $state<JiraInstance[]>([]);
  let showApiKeyCreate = $state(false);
  let showAccountDeletion = $state(false);

  let updatePassword = $state(false);

  const {
    ExternalAPIEnabled,
    LdapEnabled,
    HeaderAuthEnabled,
    SubscriptionsEnabled,
    Subscription,
    PathPrefix,
    OIDCAuthEnabled,
  } = AppConfig;

  function toggleUpdatePassword() {
    updatePassword = !updatePassword;
  }

  function getProfile() {
    xfetch(`/api/users/${$user.id}`)
      .then(res => res.json())
      .then(function (result) {
        userProfile = result.data;
      })
      .catch(function () {
        notifications.danger($LL.profileErrorRetrieving());
      });
  }

  function getCredential() {
    xfetch(`/api/users/${$user.id}/credential`)
      .then(res => res.json())
      .then(function (result) {
        userCredential = result.data;
      })
      .catch(function () {
        notifications.danger("Error retrieving user's credential");
      });
  }

  function updateUserProfile(p: any) {
    const body = {
      ...p,
    };

    xfetch(`/api/users/${$user.id}`, { body, method: 'PUT' })
      .then(res => res.json())
      .then(function () {
        user.update({
          id: userProfile.id,
          name: p.name,
          email: userProfile.email,
          rank: userProfile.rank,
          avatar: p.avatar,
          notificationsEnabled: p.notificationsEnabled,
          locale: p.locale,
        });

        if (p.theme !== 'auto') {
          localStorage.setItem('theme', p.theme);
        } else {
          localStorage.removeItem('theme');
        }
        (window as any).setTheme();

        notifications.success($LL.profileUpdateSuccess());
      })
      .catch(function () {
        notifications.danger($LL.profileErrorUpdating());
      });
  }

  function updateUserPassword(password1: string, password2: string) {
    const body = {
      password1,
      password2,
    };

    xfetch('/api/auth/update-password', { body, method: 'PATCH' })
      .then(function () {
        notifications.success($LL.passwordUpdated(), 1500);
        toggleUpdatePassword();
      })
      .catch(function () {
        notifications.danger($LL.passwordUpdateError());
      });
  }

  function getApiKeys() {
    xfetch(`/api/users/${$user.id}/apikeys`)
      .then(res => res.json())
      .then(function (result) {
        apiKeys = result.data;
      })
      .catch(function () {
        notifications.danger($LL.apiKeysErrorRetrieving());
      });
  }

  function getJiraInstances() {
    xfetch(`/api/users/${$user.id}/jira-instances`)
      .then(res => res.json())
      .then(function (result) {
        jiraInstances = result.data;
      })
      .catch(function (error) {
        if (Array.isArray(error)) {
          error[1].json().then(function (result: any) {
            if (result.error === 'REQUIRES_SUBSCRIBED_USER') {
              user.update({
                id: $user.id,
                name: $user.name,
                email: $user.email,
                rank: $user.rank,
                avatar: $user.avatar,
                verified: $user.verified,
                notificationsEnabled: $user.notificationsEnabled,
                locale: $user.locale,
                theme: $user.theme,
                subscribed: false,
              });
              notifications.danger('subscription(s) expired');
            } else {
              notifications.danger('error getting jira instances');
            }
          });
        } else {
          notifications.danger('error getting jira instances');
        }
      });
  }

  function deleteApiKey(apk: string) {
    return function () {
      xfetch(`/api/users/${$user.id}/apikeys/${apk}`, {
        method: 'DELETE',
      })
        .then(res => res.json())
        .then(function (result) {
          notifications.success($LL.apiKeyDeleteSuccess());
          apiKeys = result.data;
        })
        .catch(function () {
          notifications.danger($LL.apiKeyDeleteFailed());
        });
    };
  }

  function toggleApiKeyActiveStatus(apk: string, active: boolean) {
    return function () {
      const body = {
        active: !active,
      };

      xfetch(`/api/users/${$user.id}/apikeys/${apk}`, {
        body,
        method: 'PUT',
      })
        .then(res => res.json())
        .then(function (result) {
          notifications.success($LL.apiKeyUpdateSuccess());
          apiKeys = result.data;
        })
        .catch(function () {
          notifications.danger($LL.apiKeyUpdateFailed());
        });
    };
  }

  function toggleCreateApiKey() {
    showApiKeyCreate = !showApiKeyCreate;
  }

  function handleDeleteAccount() {
    xfetch(`/api/users/${$user.id}`, { method: 'DELETE' })
      .then(function () {
        user.delete();
        router.route(appRoutes.landing);
      })
      .catch(function () {
        notifications.danger($LL.profileDeleteError());
      });
  }

  function toggleDeleteAccount() {
    showAccountDeletion = !showAccountDeletion;
  }

  let showJiraInstanceCreate = $state(false);

  function toggleCreateJiraInstance() {
    showJiraInstanceCreate = !showJiraInstanceCreate;
  }

  function deleteJiraInstance(id: string) {
    return function () {
      xfetch(`/api/users/${$user.id}/jira-instances/${id}`, {
        method: 'DELETE',
      })
        .then(res => res.json())
        .then(function (result) {
          notifications.success('Deleted Jira instance');
          getJiraInstances();
        })
        .catch(function (error) {
          if (Array.isArray(error)) {
            error[1].json().then(function (result: any) {
              if (result.error === 'REQUIRES_SUBSCRIBED_USER') {
                user.update({
                  id: $user.id,
                  name: $user.name,
                  email: $user.email,
                  rank: $user.rank,
                  avatar: $user.avatar,
                  verified: $user.verified,
                  notificationsEnabled: $user.notificationsEnabled,
                  locale: $user.locale,
                  theme: $user.theme,
                  subscribed: false,
                });
                notifications.danger('subscription(s) expired');
              } else {
                notifications.danger('Failed to delete Jira instance');
              }
            });
          } else {
            notifications.danger('Failed to delete Jira instance');
          }
        });
    };
  }

  onMount(() => {
    if (!$user.id) {
      router.route(appRoutes.login);
      return;
    }

    getProfile();

    if ($user.rank !== 'GUEST') {
      getCredential();

      if (ExternalAPIEnabled) {
        getApiKeys();
      }
    }

    if ((SubscriptionsEnabled && $user.subscribed) || (!SubscriptionsEnabled && $user.rank !== 'GUEST')) {
      getJiraInstances();
    }
  });
</script>

<svelte:head>
  <title>{$LL.profileTitle()} | {$LL.appName()}</title>
</svelte:head>

<PageLayout>
  <h1
    class="font-semibold font-rajdhani uppercase text-2xl md:text-3xl mb-2 md:mb-6
                        md:leading-tight dark:text-white"
  >
    {$LL.profileTitle()}
  </h1>

  <div class="flex justify-center flex-wrap">
    <div class="w-full md:w-1/2 lg:w-1/3">
      {#if !updatePassword}
        <div class="bg-white dark:bg-gray-800 shadow-lg rounded-lg p-4 md:p-6 mb-4">
          <ProfileForm
            credential={userCredential}
            profile={userProfile}
            handleUpdate={updateUserProfile}
            {toggleUpdatePassword}
            {xfetch}
            {notifications}
          />
        </div>
      {/if}

      {#if updatePassword}
        <div class="bg-white dark:bg-gray-800 shadow-lg rounded-lg p-6 mb-4">
          <div
            class="font-semibold font-rajdhani uppercase text-2xl md:text-3xl mb-2 md:mb-6
        md:leading-tight text-center dark:text-white"
          >
            {$LL.updatePassword()}
          </div>

          <UpdatePasswordForm handleUpdate={updateUserPassword} toggleForm={toggleUpdatePassword} {notifications} />
        </div>
      {/if}
    </div>

    <div class="w-full md:w-1/2 lg:w-2/3">
      {#if SubscriptionsEnabled}
        <div class="ms-8 mb-8">
          <div class="flex w-full">
            <div class="flex-1">
              <h2 class="text-2xl md:text-3xl font-semibold font-rajdhani uppercase mb-4 dark:text-white">
                Active Subscriptions
              </h2>
            </div>
            <div class="flex-1">
              <div class="text-right">
                {#if $user.subscribed}
                  <SolidButton color="green" href={Subscription.ManageLink} options={{ target: '_blank' }}
                    >Manage subscriptions
                  </SolidButton>
                {/if}
              </div>
            </div>
          </div>
          {#if $user.subscribed}
            <UserSubscriptionsList {xfetch} {notifications} />
          {:else}
            <div class="bg-white dark:bg-gray-800 shadow-lg rounded-lg p-4">
              <p class="bg-yellow-thunder text-gray-900 p-4 rounded font-bold">
                See <a href={appRoutes.subscriptionPricing} class="underline">Pricing page</a> to subscribe today.
              </p>
            </div>
          {/if}
        </div>
      {/if}
      {#if ExternalAPIEnabled}
        <div class="ms-8 mb-8">
          <div class="flex w-full">
            <div class="flex-1">
              <h2 class="text-2xl md:text-3xl font-semibold font-rajdhani uppercase mb-4 dark:text-white">
                {$LL.apiKeys()}
              </h2>
            </div>
            <div class="flex-1">
              <div class="text-right">
                <HollowButton href="{PathPrefix}/swagger/index.html" options={{ target: '_blank' }} color="blue">
                  {$LL.apiDocumentation()}
                </HollowButton>
                <HollowButton onClick={toggleCreateApiKey} testid="apikey-create">
                  {$LL.apiKeyCreateButton()}
                </HollowButton>
              </div>
            </div>
          </div>

          <div class="flex flex-col">
            <div class="-my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
              <div class="py-2 align-middle inline-block min-w-full sm:px-6 lg:px-8">
                <div class="shadow overflow-hidden border-b border-gray-200 dark:border-gray-700 sm:rounded-lg">
                  <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
                    <thead class="bg-gray-50 dark:bg-gray-800">
                      <tr>
                        <th
                          scope="col"
                          class="px-6 py-3 text-left text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
                        >
                          {$LL.name()}
                        </th>
                        <th
                          scope="col"
                          class="px-6 py-3 text-left text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
                        >
                          {$LL.apiKeyPrefix()}
                        </th>
                        <th
                          scope="col"
                          class="px-6 py-3 text-left text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
                        >
                          {$LL.active()}
                        </th>
                        <th
                          scope="col"
                          class="px-6 py-3 text-left text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
                        >
                          {$LL.lastUpdated()}
                        </th>
                        <th scope="col" class="relative px-6 py-3">
                          <span class="sr-only">{$LL.actions()}</span>
                        </th>
                      </tr>
                    </thead>
                    <tbody
                      class="bg-white dark:bg-gray-700 divide-y divide-gray-200 dark:divide-gray-800 dark:text-white"
                    >
                      {#each apiKeys as apk, i}
                        <tr
                          class:bg-slate-100={i % 2 !== 0}
                          class:dark:bg-gray-800={i % 2 !== 0}
                          data-testid="apikey"
                          data-apikeyid={apk.id}
                        >
                          <td class="px-6 py-4 whitespace-nowrap" data-testid="apikey-name">{apk.name}</td>
                          <td class="px-6 py-4 whitespace-nowrap" data-testid="apikey-prefix">
                            {apk.prefix}
                          </td>
                          <td class="px-6 py-4 whitespace-nowrap" data-testid="apikey-active" data-active={apk.active}>
                            <BooleanDisplay boolValue={apk.active} />
                          </td>
                          <td class="px-6 py-4 whitespace-nowrap">
                            {new Date(apk.updatedDate).toLocaleString()}
                          </td>
                          <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                            <HollowButton
                              onClick={toggleApiKeyActiveStatus(apk.id, apk.active)}
                              testid="apikey-activetoggle"
                            >
                              {#if !apk.active}
                                {$LL.activate()}
                              {:else}
                                {$LL.deactivate()}
                              {/if}
                            </HollowButton>
                            <HollowButton color="red" onClick={deleteApiKey(apk.id)} testid="apikey-delete">
                              {$LL.delete()}
                            </HollowButton>
                          </td>
                        </tr>
                      {/each}
                    </tbody>
                  </table>
                </div>
              </div>
            </div>
          </div>
        </div>
      {/if}

      <div class="ms-8">
        <div class="flex w-full">
          <div class="flex-1">
            <h2 class="text-2xl md:text-3xl font-semibold font-rajdhani uppercase mb-4 dark:text-white">
              Jira Instances
            </h2>
          </div>
          {#if (SubscriptionsEnabled && $user.subscribed) || (!SubscriptionsEnabled && $user.rank !== 'GUEST')}
            <div class="flex-1">
              <div class="text-right">
                <HollowButton onClick={toggleCreateJiraInstance} testid="jirainstance-create">
                  Add Jira Instance
                </HollowButton>
              </div>
            </div>
          {/if}
        </div>
        <div class="flex flex-col">
          <div class="-my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
            <div class="py-2 align-middle inline-block min-w-full sm:px-6 lg:px-8">
              <div class="shadow overflow-hidden border-b border-gray-200 dark:border-gray-700 sm:rounded-lg">
                {#if SubscriptionsEnabled && !$user.subscribed}
                  <FeatureSubscribeBanner salesPitch="Setup Jira Cloud to import your stories in Poker Planning." />
                {:else if !SubscriptionsEnabled && $user.rank === 'GUEST'}
                  <p class="bg-sky-300 p-4 rounded text-gray-700 font-bold">
                    Must be logged in to setup Jira integrations
                  </p>
                {:else}
                  <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
                    <thead class="bg-gray-50 dark:bg-gray-800">
                      <tr>
                        <th
                          scope="col"
                          class="px-6 py-3 text-left text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
                        >
                          Host
                        </th>
                        <th
                          scope="col"
                          class="px-6 py-3 text-left text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
                        >
                          Client Mail
                        </th>
                        <th scope="col" class="relative px-6 py-3">
                          <span class="sr-only">{$LL.actions()}</span>
                        </th>
                      </tr>
                    </thead>
                    <tbody
                      class="bg-white dark:bg-gray-700 divide-y divide-gray-200 dark:divide-gray-800 dark:text-white"
                    >
                      {#each jiraInstances as ji, i}
                        <tr
                          class:bg-slate-100={i % 2 !== 0}
                          class:dark:bg-gray-800={i % 2 !== 0}
                          data-testid="apikey"
                          data-apikeyid={ji.id}
                        >
                          <td class="px-6 py-4 whitespace-nowrap" data-testid="jira-host">{ji.host}</td>
                          <td class="px-6 py-4 whitespace-nowrap" data-testid="jira-clientmail">
                            {ji.client_mail}
                          </td>
                          <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                            <HollowButton color="red" onClick={deleteJiraInstance(ji.id)} testid="jira-delete">
                              {$LL.delete()}
                            </HollowButton>
                          </td>
                        </tr>
                      {/each}
                    </tbody>
                  </table>
                {/if}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    {#if !OIDCAuthEnabled && !LdapEnabled && !HeaderAuthEnabled}
      <div class="w-full text-center mt-8">
        <HollowButton onClick={toggleDeleteAccount} color="red">
          {$LL.deleteAccount()}
        </HollowButton>
      </div>
    {/if}
  </div>
  {#if showApiKeyCreate}
    <CreateApiKey {toggleCreateApiKey} handleApiKeyCreate={getApiKeys} {notifications} {xfetch} />
  {/if}
  {#if showJiraInstanceCreate}
    <CreateJiraInstance
      toggleClose={toggleCreateJiraInstance}
      handleCreate={getJiraInstances}
      {notifications}
      {xfetch}
    />
  {/if}

  {#if showAccountDeletion}
    <DeleteConfirmation
      toggleDelete={toggleDeleteAccount}
      handleDelete={handleDeleteAccount}
      confirmText={$LL.deleteAccountWarningStatement()}
      confirmBtnText={$LL.deleteConfirmButton()}
    />
  {/if}
</PageLayout>
