<script lang="ts">
  import SolidButton from '../global/SolidButton.svelte';
  import SelectInput from '../forms/SelectInput.svelte';
  import { AppConfig, appRoutes } from '../../config';
  import { user } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import { createEventDispatcher, onMount } from 'svelte';

  const dispatch = createEventDispatcher();

  export let handleImport = story => {};
  export let eventTag;
  export let notifications;
  export let xfetch;

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

  let jiraInstances = [];
  let jiraStories = [];
  let selectedJiraInstance = '';
  let searchJQL = '';
  let jqlError = '';

  function getJiraInstances() {
    xfetch(`/api/users/${$user.id}/jira-instances`)
      .then(res => res.json())
      .then(function (result) {
        jiraInstances = result.data;
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
              });
              notifications.danger('subscription(s) expired');
            } else {
              notifications.danger('error getting jira instances');
            }
          });
        } else {
          notifications.danger('error getting jira instances');
        }
        eventTag('fetch_profile_jira_instances', 'engagement', 'failure');
      });
  }

  function handleJQLSearch(event) {
    event.preventDefault();

    if (searchJQL === '') {
      notifications.danger(
        'Must enter a JQL search query ex: order by created DESC',
        5000,
      );
      return;
    }

    jiraStories = [];

    xfetch(
      `/api/users/${$user.id}/jira-instances/${jiraInstances[selectedJiraInstance].id}/jql-story-search`,
      {
        body: {
          jql: searchJQL,
          startAt: 0,
          maxResults: 100,
        },
      },
    )
      .then(res => res.json())
      .then(function (result) {
        jqlError = '';
        jiraStories = result.data.issues;
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
              });
              jqlError = 'subscription(s) expired';
            } else {
              jqlError = `Jira JQL Search Error: ${result.error}`;
            }
          });
        } else {
          notifications.danger('Unknown Jira JQL search error');
        }
        eventTag('fetch_profile_jira_instances', 'engagement', 'failure');
      });
  }

  function findPlanType(type) {
    return planTypes.includes(type) ? type : 'Story';
  }

  function findPriority(priority) {
    const found = priorities.find(p => p.name === priority);
    return found ? found.value : 99;
  }

  function stripTrailingSlash(str) {
    return str.endsWith('/') ? str.slice(0, -1) : str;
  }

  function importStory(idx) {
    return function () {
      const story = jiraStories[idx];
      handleImport({
        name: story.fields.summary,
        type: findPlanType(story.fields.issuetype.name),
        referenceId: story.key,
        link: `${stripTrailingSlash(
          jiraInstances[selectedJiraInstance].host,
        )}/browse/${story.key}`,
        description: '', // @TODO - get description
        priority: findPriority(story.fields.priority.name),
      });
    };
  }

  onMount(() => {
    if (
      (AppConfig.SubscriptionsEnabled && $user.subscribed) ||
      !AppConfig.SubscriptionsEnabled
    ) {
      getJiraInstances();
    }
  });
</script>

{#if AppConfig.SubscriptionsEnabled && !$user.subscribed}
  <p class="bg-yellow-thunder text-gray-900 p-4 rounded font-bold">
    Must be <a
      href="{appRoutes.subscriptionPricing}"
      class="underline"
      target="_blank">subscribed</a
    >
    to import from Jira Cloud.
  </p>
{:else if jiraInstances.length === 0}
  <p class="bg-yellow-thunder text-gray-900 p-4 rounded font-bold">
    Visit your <a href="{appRoutes.profile}" class="underline" target="_blank"
      >profile page</a
    > to setup instances of Jira Cloud.
  </p>
{:else}
  <div class="mb-4">
    <SelectInput
      id="jirainstance"
      bind:value="{selectedJiraInstance}"
      on:change="{() => {
        dispatch('instance_selected');
      }}"
    >
      >
      <option value="" disabled>Select Jira Instance to import from </option>
      {#each jiraInstances as ji, idx}
        <option value="{idx}">{ji.host}</option>
      {/each}
    </SelectInput>
  </div>

  {#if selectedJiraInstance !== ''}
    <form on:submit="{handleJQLSearch}" class="mb-4">
      <label
        for="jql-search"
        class="mb-2 text-sm font-medium text-gray-900 sr-only dark:text-white"
        >Search</label
      >
      <div class="relative">
        <div
          class="absolute inset-y-0 start-0 flex items-center ps-3 pointer-events-none"
        >
          <svg
            class="w-4 h-4 text-gray-500 dark:text-gray-400"
            aria-hidden="true"
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 20 20"
          >
            <path
              stroke="currentColor"
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="m19 19-4-4m0-7A7 7 0 1 1 1 8a7 7 0 0 1 14 0Z"></path>
          </svg>
        </div>
        <input
          type="search"
          id="jql-search"
          class="block w-full p-4 ps-10 text-sm text-gray-900 border border-gray-300 rounded-lg bg-gray-50 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
          placeholder="Enter Search JQL..."
          bind:value="{searchJQL}"
          required
        />
        <button
          type="submit"
          class="text-white absolute end-2.5 bottom-2.5 bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-4 py-2 dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800"
        >
          Search
        </button>
      </div>
    </form>

    <div class="flex flex-wrap">
      {#if jqlError !== ''}
        <div
          class="p-4 bg-red-50 border-red-500 text-red-800 font-semibold rounded"
        >
          {jqlError}
        </div>
      {/if}
      {#each jiraStories as story, idx}
        <div
          class="p-2 w-full flex flex-wrap justify-between bg-gray-200 dark:bg-gray-900 dark:text-white rounded shadow mb-2 items-center"
        >
          <div>
            [{story.key}] {story.fields.summary}
          </div>
          <div>
            <SolidButton onClick="{importStory(idx)}">Import</SolidButton>
          </div>
        </div>
      {/each}
    </div>
  {/if}
{/if}
