<script lang="ts">
  import Modal from '../global/Modal.svelte';
  import LL from '../../i18n/i18n-svelte';
  import CsvImport from './CsvImport.svelte';
  import JiraImport from './JiraImport.svelte';
  import JQLImport from '../jira/JQLImport.svelte';
  import SolidButton from '../global/SolidButton.svelte';
  import StoryFromGameImport from './StoryFromGameImport.svelte';
  import { AppConfig, appRoutes } from '../../config';
  import { user } from '../../stores';

  export let notifications;
  export let eventTag;
  export let xfetch;
  export let toggleImport = () => {};
  export let handlePlanAdd = handleAdd => {};
  export let gameId = '';

  let showJiraCloudSearch = false;
  let showGameImport = false;

  const toggleGameImport = () => {
    showGameImport = !showGameImport;
  };

  const handleAdd = newPlan => {
    handlePlanAdd(newPlan);
    toggleImport();
  };

  function importStory(story) {
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
    {#if !showJiraCloudSearch}
      <div class="mb-4 dark:text-gray-300">
        <h3 class="font-bold mb-2 text-xl">Internal Import</h3>
        {#if AppConfig.SubscriptionsEnabled && !$user.subscribed}
          <p class="bg-yellow-thunder text-gray-900 p-4 rounded font-bold">
            Must be <a
              href="{appRoutes.subscriptionPricing}"
              class="underline"
              target="_blank">subscribed</a
            >
            to import from other Games.
          </p>
        {:else if !AppConfig.SubscriptionsEnabled || (AppConfig.SubscriptionsEnabled && $user.subscribed)}
          {#if !showGameImport}
            <SolidButton color="indigo" onClick="{toggleGameImport}"
              >Import from another Game
            </SolidButton>
          {/if}
        {/if}
      </div>

      {#if showGameImport}
        <StoryFromGameImport
          notifications="{notifications}"
          xfetch="{xfetch}"
          eventTag="{eventTag}"
          handleImport="{importStory}"
          gameId="{gameId}"
        />
      {/if}
    {/if}

    {#if !showGameImport}
      <div class="mb-4 dark:text-gray-300">
        <h3 class="font-bold mb-2 text-xl">Import from Jira Cloud</h3>
        <JQLImport
          notifications="{notifications}"
          xfetch="{xfetch}"
          eventTag="{eventTag}"
          handleImport="{importStory}"
          on:instance_selected="{() => {
            showJiraCloudSearch = true;
          }}"
        />
      </div>
    {/if}

    {#if !showJiraCloudSearch && !showGameImport}
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
