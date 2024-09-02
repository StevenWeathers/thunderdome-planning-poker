<script lang="ts">
  import Modal from '../global/Modal.svelte';
  import LL from '../../i18n/i18n-svelte';
  import CsvImport from './CsvImport.svelte';
  import JiraImport from './JiraImport.svelte';
  import JQLImport from '../jira/JQLImport.svelte';

  export let notifications;
  export let eventTag;
  export let xfetch;
  export let toggleImport = () => {};
  export let handlePlanAdd = handleAdd => {};

  let showJiraCloudSearch = false;

  const handleAdd = newPlan => {
    handlePlanAdd(newPlan);
    toggleImport();
  };

  function importJQLStory(story) {
    handlePlanAdd({
      planName: story.name,
      type: story.type,
      referenceId: story.referenceId,
      link: story.link,
      description: story.description,
      priority: story.priority,
    });
  }
</script>

<Modal closeModal="{toggleImport}" widthClasses="md:w-full md:mx-4 lg:w-3/5">
  <div class="mt-8 mb-4">
    <div class="mb-4 dark:text-gray-300">
      <h3 class="font-bold mb-2 text-xl">Import from Jira Cloud</h3>
      <JQLImport
        notifications="{notifications}"
        xfetch="{xfetch}"
        eventTag="{eventTag}"
        handleImport="{importJQLStory}"
        on:instance_selected="{() => {
          showJiraCloudSearch = true;
        }}"
      />
    </div>

    {#if !showJiraCloudSearch}
      <div class="md:grid md:grid-cols-2 md:gap-4">
        <div class="mb-4">
          <h3 class="font-bold mb-2 dark:text-gray-300 text-lg">
            {$LL.importJiraXML()}
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
            {$LL.importCsv()}
          </h3>
          <p class="dark:text-gray-400 mb-2">
            The CSV file must include all the following fields with no header
            row:
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
    {/if}
  </div>
</Modal>
