<script lang="ts">
  import Modal from '../global/Modal.svelte';
  import LL from '../../i18n/i18n-svelte';
  import CsvImport from './CsvImport.svelte';
  import JiraImport from './JiraImport.svelte';
  import SolidButton from '../global/SolidButton.svelte';
  import { AppConfig, appRoutes } from '../../config';
  import { user } from '../../stores';
  import { onMount } from 'svelte';
  import SelectInput from '../global/SelectInput.svelte';
  import TextInput from '../global/TextInput.svelte';

  export let notifications;
  export let eventTag;
  export let xfetch;
  export let toggleImport = () => {};
  export let handlePlanAdd = handleAdd => {};

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
  let showJiraCloudSearch = false;
  let searchJQL = '';

  const handleAdd = newPlan => {
    handlePlanAdd(newPlan);
    toggleImport();
  };

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

  function displayJiraCloudSearch() {
    showJiraCloudSearch = true;
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
      const link = handlePlanAdd({
        planName: story.fields.summary,
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
    if ($user.subscribed) {
      getJiraInstances();
    }
  });
</script>

<Modal closeModal="{toggleImport}" widthClasses="md:w-2/3 lg:w-3/5 xl:w-1/2">
  <div class="mt-8 mb-4">
    <div class="mb-4 dark:text-gray-300">
      <h3 class="font-bold mb-2 text-xl">Import from Jira Cloud</h3>
      {#if !$user.subscribed}
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
          Visit your <a
            href="{appRoutes.profile}"
            class="underline"
            target="_blank">profile page</a
          > to setup instances of Jira Cloud.
        </p>
      {:else}
        <div class="mb-4">
          <SelectInput
            id="jirainstance"
            bind:value="{selectedJiraInstance}"
            on:change="{displayJiraCloudSearch}"
          >
            >
            <option value="" disabled
              >Select Jira Instance to import from
            </option>
            {#each jiraInstances as ji, idx}
              <option value="{idx}">{ji.host}</option>
            {/each}
          </SelectInput>
        </div>

        {#if showJiraCloudSearch}
          <form on:submit="{handleJQLSearch}" class="mb-4">
            <TextInput
              placeholder="Enter Search JQL..."
              bind:value="{searchJQL}"
              type="search"
            />
          </form>

          <div class="flex flex-wrap">
            {#each jiraStories as story, idx}
              <div
                class="p-4 w-full flex bg-gray-200 dark:bg-gray-600 dark:text-white rounded shadow mb-2"
              >
                <div class="w-3/4">
                  [{story.key}] {story.fields.summary}
                </div>
                <div class="w-1/4 text-right">
                  <SolidButton onClick="{importStory(idx)}">Import</SolidButton>
                </div>
              </div>
            {/each}
          </div>
        {/if}
      {/if}
    </div>
    <div class="mb-4">
      <h3 class="font-bold mb-2 dark:text-gray-300 text-lg">
        {$LL.importJiraXML({ friendly: AppConfig.FriendlyUIVerbs })}
      </h3>
      <JiraImport
        handlePlanAdd="{handleAdd}"
        notifications="{notifications}"
        eventTag="{eventTag}"
        testid="plans-importjira"
      />
    </div>
    <div class="mb-4">
      <h3 class="font-bold mb-2 dark:text-gray-300 text-lg">
        {$LL.importCsv({ friendly: AppConfig.FriendlyUIVerbs })}
      </h3>
      <p class="dark:text-gray-400 mb-2">
        The CSV file must include all the following fields with no header row:
      </p>
      <div class="mb-2 dark:text-gray-300">
        Type,Title,ReferenceId,Link,Description,AcceptanceCriteria
      </div>
      <CsvImport
        handlePlanAdd="{handleAdd}"
        notifications="{notifications}"
        eventTag="{eventTag}"
        testid="plans-Csvimport"
      />
    </div>
  </div>
  <div class="text-right">
    <SolidButton onClick="{toggleImport}">
      {$LL.cancel()}
    </SolidButton>
  </div>
</Modal>
