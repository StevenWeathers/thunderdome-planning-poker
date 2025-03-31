<script lang="ts">
  import Modal from '../global/Modal.svelte';
  import LL from '../../i18n/i18n-svelte';
  import CsvImport from './CsvImport.svelte';
  import JiraImport from './JiraImport.svelte';
  import JQLImport from '../jira/JQLImport.svelte';
  import SolidButton from '../global/SolidButton.svelte';
  import StoryFromGameImport from './StoryFromGameImport.svelte';
  import { AppConfig } from '../../config';
  import { user } from '../../stores';
  import StoryFromStoryboardImport from './StoryFromStoryboardImport.svelte';
  import FeatureSubscribeBanner from '../global/FeatureSubscribeBanner.svelte';

  interface Props {
    notifications: any;
    xfetch: any;
    toggleImport?: any;
    handlePlanAdd?: any;
    gameId?: string;
  }

  let {
    notifications,
    xfetch,
    toggleImport = () => {},
    handlePlanAdd = handleAdd => {},
    gameId = ''
  }: Props = $props();

  let showJiraCloudSearch = $state(false);
  let showGameImport = $state(false);
  let showStoryboardImport = $state(false);

  const toggleGameImport = () => {
    showGameImport = !showGameImport;
  };

  const toggleStoryboardImport = () => {
    showStoryboardImport = !showStoryboardImport;
  };

  const handleAdd = newPlan => {
    handlePlanAdd(newPlan);
    toggleImport();
  };

  function importStory(story) {
    handlePlanAdd({
      planName: story.name,
      type: story.type || '',
      referenceId: story.referenceId || '',
      link: story.link || '',
      description: story.description || '',
      priority: story.priority || 99,
    });
  }
</script>

<Modal closeModal="{toggleImport}" widthClasses="md:w-full md:mx-4 lg:w-3/5">
  <div class="mt-8 mb-4">
    {#if !showJiraCloudSearch}
      <div class="mb-4 dark:text-gray-300">
        <h3 class="font-bold mb-2 text-xl">Internal Import</h3>
        {#if AppConfig.SubscriptionsEnabled && !$user.subscribed}
          <FeatureSubscribeBanner
            salesPitch="Import stories from other Poker Plannings or Storyboards."
          />
        {:else if !AppConfig.SubscriptionsEnabled || (AppConfig.SubscriptionsEnabled && $user.subscribed)}
          {#if !showGameImport && !showStoryboardImport}
            <SolidButton color="indigo" onClick="{toggleGameImport}"
              >Import from another Game
            </SolidButton>
            <SolidButton color="blue" onClick="{toggleStoryboardImport}"
              >Import from a Storyboard
            </SolidButton>
          {/if}
        {/if}
      </div>

      {#if showGameImport}
        <StoryFromGameImport
          notifications={notifications}
          xfetch={xfetch}
          handleImport="{importStory}"
          gameId="{gameId}"
        />
      {/if}

      {#if showStoryboardImport}
        <StoryFromStoryboardImport
          notifications={notifications}
          xfetch={xfetch}
          handleImport="{importStory}"
        />
      {/if}
    {/if}

    {#if !showGameImport && !showStoryboardImport}
      <div class="mb-4 dark:text-gray-300">
        <h3 class="font-bold mb-2 text-xl">Import from Jira Cloud</h3>
        <JQLImport
          notifications={notifications}
          xfetch={xfetch}
          handleImport="{importStory}"
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
              notifications={notifications}
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
              notifications={notifications}
              testid="plans-Csvimport"
            />
          </div>
        </div>
      {/if}
    {/if}
  </div>
</Modal>
