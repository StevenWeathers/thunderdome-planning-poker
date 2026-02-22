<script lang="ts">
  import SolidButton from '../global/SolidButton.svelte';
  import SelectInput from '../forms/SelectInput.svelte';
  import { AppConfig, appRoutes } from '../../config';
  import { user } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import { createEventDispatcher, onMount } from 'svelte';
  import FeatureSubscribeBanner from '../global/FeatureSubscribeBanner.svelte';

  import type { NotificationService } from '../../types/notifications';
  import type { ApiClient } from '../../types/apiclient';
  import type { SessionUser } from '../../types/user';
  import { SearchIcon } from '@lucide/svelte';

  const dispatch = createEventDispatcher();

  interface Props {
    handleImport?: any;
    notifications: NotificationService;
    xfetch: ApiClient;
  }

  let { handleImport = (story: any) => {}, notifications, xfetch }: Props = $props();

  // going by common Jira issue types for now
  const planTypes = [
    $LL.planTypeStory(),
    $LL.planTypeBug(),
    $LL.planTypeSpike(),
    $LL.planTypeEpic(),
    $LL.planTypeTask(),
    $LL.planTypeSubtask(),
  ];

  // going by common Jira issue priorities for now
  const priorities = [
    { name: $LL.planPriorityBlocker(), value: 1 },
    {
      name: $LL.planPriorityHighest(),
      value: 2,
    },
    { name: $LL.planPriorityHigh(), value: 3 },
    { name: $LL.planPriorityMedium(), value: 4 },
    { name: $LL.planPriorityLow(), value: 5 },
    {
      name: $LL.planPriorityLowest(),
      value: 6,
    },
  ];

  let jiraInstances = $state([]);
  let jiraStories = $state([]);
  let selectedJiraInstance: string = $state('');
  let searchJQL: string = $state('');
  let jqlError: string = $state('');
  let importedStoryKeys = $state([]);
  let searchCompleted: boolean = $state(false);

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
              } as SessionUser);
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

  function handleJQLSearch(event: Event) {
    event.preventDefault();

    if (searchJQL === '') {
      notifications.danger('Must enter a JQL search query ex: order by created DESC', 5000);
      return;
    }

    jiraStories = [];
    importedStoryKeys = [];
    searchCompleted = false;

    xfetch(`/api/users/${$user.id}/jira-instances/${jiraInstances[selectedJiraInstance].id}/jql-story-search`, {
      body: {
        jql: searchJQL,
        startAt: 0,
        maxResults: 100,
      },
    })
      .then(res => res.json())
      .then(function (result) {
        jqlError = '';
        jiraStories = result.data.issues;
        searchCompleted = true;
      })
      .catch(function (error) {
        if (Array.isArray(error)) {
          error[1].json().then(function (result) {
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
              } as SessionUser);
              jqlError = 'subscription(s) expired';
            } else {
              jqlError = `Jira JQL Search Error: ${result.error}`;
            }
            searchCompleted = true;
          });
        } else {
          notifications.danger('Unknown Jira JQL search error');
          searchCompleted = true;
        }
      });
  }

  function findPlanType(type: string) {
    return planTypes.includes(type) ? type : 'Story';
  }

  function findPriority(priority) {
    const found = priorities.find(p => p.name === priority);
    return found ? found.value : 99;
  }

  function stripTrailingSlash(str: string) {
    return str.endsWith('/') ? str.slice(0, -1) : str;
  }

  function importStory(idx: number) {
    return function () {
      const story = jiraStories[idx];
      handleImport({
        name: story.fields.summary,
        type: findPlanType(story.fields.issuetype.name),
        referenceId: story.key,
        link: `${stripTrailingSlash(jiraInstances[selectedJiraInstance].host)}/browse/${story.key}`,
        description: '', // @TODO - get description
        priority: findPriority(story.fields.priority.name),
      });
      importedStoryKeys = [...importedStoryKeys, story.key];
    };
  }

  function importAllStories() {
    const storiesToImport = jiraStories.filter(story => !importedStoryKeys.includes(story.key));
    storiesToImport.forEach(story => {
      handleImport({
        name: story.fields.summary,
        type: findPlanType(story.fields.issuetype.name),
        referenceId: story.key,
        link: `${stripTrailingSlash(jiraInstances[selectedJiraInstance].host)}/browse/${story.key}`,
        description: '', // @TODO - get description
        priority: findPriority(story.fields.priority.name),
      });
      importedStoryKeys = [...importedStoryKeys, story.key];
    });
  }

  onMount(() => {
    if ((AppConfig.SubscriptionsEnabled && $user.subscribed) || !AppConfig.SubscriptionsEnabled) {
      getJiraInstances();
    }
  });
</script>

{#if AppConfig.SubscriptionsEnabled && !$user.subscribed}
  <FeatureSubscribeBanner salesPitch="Import your stories for Poker Planning from Jira Cloud." />
{:else if !AppConfig.SubscriptionsEnabled || (AppConfig.SubscriptionsEnabled && $user.subscribed)}
  {#if jiraInstances.length === 0}
    <p class="info-banner">
      Visit your <a href={appRoutes.profile} class="info-banner-link" target="_blank">profile page</a> to setup instances
      of Jira Cloud.
    </p>
  {:else}
    <div class="select-wrapper">
      <SelectInput
        id="jirainstance"
        bind:value={selectedJiraInstance}
        on:change={() => {
          dispatch('instance_selected');
        }}
      >
        <option value="" disabled>Select Jira Instance to import from</option>
        {#each jiraInstances as ji, idx}
          <option value={idx}>{ji.host}</option>
        {/each}
      </SelectInput>
    </div>

    {#if selectedJiraInstance !== ''}
      <form onsubmit={handleJQLSearch} class="search-form">
        <label for="jql-search" class="search-label">Search</label>
        <div class="search-container">
          <div class="search-icon-wrapper">
            <SearchIcon class="w-4 h-4 text-gray-500 dark:text-gray-400" />
          </div>
          <input
            type="search"
            id="jql-search"
            class="search-input"
            placeholder="Enter Search JQL..."
            bind:value={searchJQL}
            required
          />
          <button type="submit" class="search-button"> Search </button>
        </div>
      </form>

      <div class="stories-wrapper">
        {#if jqlError !== ''}
          <div class="error-message">
            {jqlError}
          </div>
        {/if}
        {#if searchCompleted && jiraStories.length === 0 && jqlError === ''}
          <p class="no-stories-message">No stories found for this JQL query.</p>
        {/if}
        {#if jiraStories.length > 0 && importedStoryKeys.length === jiraStories.length}
          <p class="all-imported-message">All stories have been imported!</p>
        {/if}
        {#if jiraStories.length > 0 && importedStoryKeys.length < jiraStories.length}
          <div class="search-result-header">
            <span class="text-lg font-semibold">
              Search Results {jiraStories.length > 0 ? `(${jiraStories.length - importedStoryKeys.length})` : ''}
            </span>
            <SolidButton onClick={importAllStories}>Import All</SolidButton>
          </div>
        {/if}
        {#each jiraStories as story, idx}
          <div
            class="story-item"
            class:hidden={importedStoryKeys.includes(story.key)}
            aria-hidden={importedStoryKeys.includes(story.key)}
          >
            <div>
              [{story.key}] {story.fields.summary}
            </div>
            <div>
              <SolidButton onClick={importStory(idx)}>Import</SolidButton>
            </div>
          </div>
        {/each}
      </div>
    {/if}
  {/if}
{/if}

<style lang="postcss">
  .info-banner {
    @apply bg-yellow-thunder text-gray-900 p-4 rounded font-bold;
  }

  :root.dark .info-banner {
    @apply bg-yellow-600 text-white;
  }

  .info-banner-link {
    @apply underline;
  }

  .select-wrapper {
    @apply mb-4;
  }

  .search-form {
    @apply mb-4;
  }

  .search-label {
    @apply mb-2 text-sm font-medium text-gray-900 sr-only;
  }

  :root.dark .search-label {
    @apply text-white;
  }

  .search-container {
    @apply relative;
  }

  .search-icon-wrapper {
    @apply absolute inset-y-0 start-0 flex items-center ps-3 pointer-events-none;
  }

  :root.dark .search-icon-wrapper {
    @apply text-gray-400;
  }

  .search-input {
    @apply block w-full p-4 ps-10 text-sm text-gray-900 border border-gray-300 rounded-lg bg-gray-50;
  }

  :root.dark .search-input {
    @apply bg-gray-700 border-gray-600 placeholder-gray-400 text-white;
  }

  .search-input:focus {
    @apply ring-blue-500 border-blue-500;
  }

  :root.dark .search-input:focus {
    @apply ring-blue-500 border-blue-500;
  }

  .search-button {
    @apply text-white absolute end-2.5 bottom-2.5 bg-blue-700 font-medium rounded-lg text-sm px-4 py-2;
  }

  :root.dark .search-button {
    @apply bg-blue-600;
  }

  .search-button:hover {
    @apply bg-blue-800;
  }

  :root.dark .search-button:hover {
    @apply bg-blue-700;
  }

  .search-button:focus {
    @apply ring-4 outline-none ring-blue-300;
  }

  :root.dark .search-button:focus {
    @apply ring-blue-800;
  }

  .search-result-header {
    @apply flex justify-between items-center mb-4 text-gray-700 w-full;
  }

  :root.dark .search-result-header {
    @apply text-gray-300;
  }

  .stories-wrapper {
    @apply flex flex-wrap;
  }

  .error-message {
    @apply p-4 bg-red-50 border-red-500 text-red-800 font-semibold rounded-lg w-full;
  }

  :root.dark .error-message {
    @apply bg-red-900 border-red-700 text-red-200;
  }

  .no-stories-message {
    @apply p-4 text-gray-700 text-center italic w-full rounded-lg bg-gray-200;
  }

  :root.dark .no-stories-message {
    @apply text-gray-200 bg-gray-700;
  }

  .all-imported-message {
    @apply p-4 text-green-700 text-center font-semibold bg-green-50 rounded-lg w-full;
  }

  :root.dark .all-imported-message {
    @apply text-green-200 bg-green-900;
  }

  .story-item {
    padding: 0.5rem;
    width: 100%;
    display: flex;
    flex-wrap: wrap;
    justify-content: space-between;
    background-color: theme('colors.gray.200');
    border-radius: 0.5rem;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
    align-items: center;
    transition: all 200ms ease-out;
    margin-bottom: 0.5rem;
  }

  :root.dark .story-item {
    background-color: theme('colors.gray.700');
    color: white;
  }

  .story-item.hidden {
    pointer-events: none;
    max-height: 0;
    padding: 0;
    overflow: hidden;
    margin-bottom: 0;
    transition:
      padding 200ms ease-out,
      max-height 200ms ease-out,
      margin-bottom 200ms ease-out;
  }
</style>
