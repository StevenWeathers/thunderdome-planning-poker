<script lang="ts">
  import Modal from '../Modal.svelte';
  import LL from '../../i18n/i18n-svelte';
  import CsvImport from './CsvImport.svelte';
  import JiraImport from './JiraImport.svelte';
  import SolidButton from '../SolidButton.svelte';
  import { AppConfig } from '../../config';
  import { user } from '../../stores';
  import { onMount } from 'svelte';
  import SelectInput from '../SelectInput.svelte';
  import TextInput from '../TextInput.svelte';

  export let notifications;
  export let eventTag;
  export let xfetch;
  export let toggleImport = () => {};
  export let handlePlanAdd = handleAdd => {};

  let jiraInstances = [];
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
      .catch(function () {
        notifications.danger('error getting jira instances');
        eventTag('fetch_profile_jira_instances', 'engagement', 'failure');
      });
  }

  function displayJiraCloudSearch() {
    showJiraCloudSearch = true;
  }

  function handleJQLSearch(event) {
    event.preventDefault();

    // @TODO - implement search
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
        <p class="bg-sky-300 p-4 rounded text-gray-700 font-bold">
          Must be subscribed to import from Jira Cloud
        </p>
      {:else}
        <SelectInput
          id="jirainstance"
          bind:value="{selectedJiraInstance}"
          on:change="{displayJiraCloudSearch}"
        >
          >
          <option value="" disabled>Select Jira Instance to import from</option>
          {#each jiraInstances as ji}
            <option value="{ji.id}">{ji.host}</option>
          {/each}
        </SelectInput>

        {#if showJiraCloudSearch}
          <form on:submit="{handleJQLSearch}">
            <TextInput placeholder="Enter Search JQL..." value="{searchJQL}" />
          </form>
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
