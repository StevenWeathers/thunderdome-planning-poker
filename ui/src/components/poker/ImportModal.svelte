<script lang="ts">
  import Modal from '../global/Modal.svelte';
  import LL from '../../i18n/i18n-svelte';
  import CsvImport from './CsvImport.svelte';
  import JiraImport from './JiraImport.svelte';
  import JQLImport from '../jira/JQLImport.svelte';
  import StoryFromGameImport from './StoryFromGameImport.svelte';
  import { AppConfig } from '../../config';
  import { user } from '../../stores';
  import StoryFromStoryboardImport from './StoryFromStoryboardImport.svelte';
  import FeatureSubscribeBanner from '../global/FeatureSubscribeBanner.svelte';
  import { DownloadCloud, FileText, FilePlus, FileCode, FileSpreadsheet } from '@lucide/svelte';

  import type { NotificationService } from '../../types/notifications';
  import type { ApiClient } from '../../types/apiclient';
  import HollowButton from '../global/HollowButton.svelte';

  interface Props {
    notifications: NotificationService;
    xfetch: ApiClient;
    toggleImport?: any;
    handlePlanAdd?: any;
    gameId?: string;
  }

  let {
    notifications,
    xfetch,
    toggleImport = () => {},
    handlePlanAdd = handleAdd => {},
    gameId = '',
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

<Modal
  closeModal={toggleImport}
  widthClasses="md:w-full md:mx-4 lg:w-4/5 xl:w-3/5"
  ariaLabel={$LL.modalImportPokerStories()}
>
  <div>
    <!-- Header -->
    <div class="mb-8">
      <h2 class="text-2xl sm:text-3xl font-bold text-gray-900 dark:text-white mb-2">Import Stories</h2>
      <p class="text-gray-600 dark:text-gray-400 text-sm sm:text-base">
        Choose your preferred method to import poker planning stories
      </p>
    </div>

    <!-- Internal Import Section -->
    {#if !showJiraCloudSearch}
      <div class="mb-8">
        <div
          class="bg-gradient-to-r from-indigo-50 to-blue-50 dark:from-indigo-950/20 dark:to-blue-950/20 rounded-xl p-6 border border-indigo-100 dark:border-indigo-800/30"
        >
          <div class="flex items-center gap-3 mb-4">
            <div class="w-8 h-8 bg-indigo-500 rounded-lg flex items-center justify-center">
              <DownloadCloud class="w-4 h-4 text-white" />
            </div>
            <h3 class="text-xl font-semibold text-gray-900 dark:text-white">Internal Import</h3>
          </div>

          {#if AppConfig.SubscriptionsEnabled && !$user.subscribed}
            <FeatureSubscribeBanner salesPitch="Import stories from other Poker Plannings or Storyboards." />
          {:else if !AppConfig.SubscriptionsEnabled || (AppConfig.SubscriptionsEnabled && $user.subscribed)}
            {#if !showGameImport && !showStoryboardImport}
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
                <HollowButton fullWidth={true} size="large" color="blue" onClick={toggleGameImport}>
                  Import from Game
                </HollowButton>
                <HollowButton fullWidth={true} size="large" color="purple" onClick={toggleStoryboardImport}>
                  Import from Storyboard
                </HollowButton>
              </div>
            {/if}
          {/if}
        </div>

        {#if showGameImport}
          <div
            class="mt-6 bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 p-6 shadow-sm"
          >
            <StoryFromGameImport {notifications} {xfetch} handleImport={importStory} {gameId} />
          </div>
        {/if}

        {#if showStoryboardImport}
          <div
            class="mt-6 bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 p-6 shadow-sm"
          >
            <StoryFromStoryboardImport {notifications} {xfetch} handleImport={importStory} />
          </div>
        {/if}
      </div>
    {/if}

    <!-- External Import Sections -->
    {#if !showGameImport && !showStoryboardImport}
      <!-- Jira Cloud Section -->
      <div class="mb-8">
        <div
          class="bg-gradient-to-r from-blue-50 to-cyan-50 dark:from-blue-950/20 dark:to-cyan-950/20 rounded-xl p-6 border border-blue-100 dark:border-blue-800/30"
        >
          <div class="flex items-center gap-3 mb-4">
            <div class="w-8 h-8 bg-blue-500 rounded-lg flex items-center justify-center">
              <FileText class="w-4 h-4 text-white" />
            </div>
            <h3 class="text-xl font-semibold text-gray-900 dark:text-white">Import from Jira Cloud</h3>
          </div>

          <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
            <JQLImport
              {notifications}
              {xfetch}
              handleImport={importStory}
              on:instance_selected={() => {
                showJiraCloudSearch = true;
              }}
            />
          </div>
        </div>
      </div>

      {#if !showJiraCloudSearch}
        <!-- File Import Section -->
        <div class="space-y-6">
          <div
            class="bg-gradient-to-r from-amber-50 to-orange-50 dark:from-amber-950/20 dark:to-orange-950/20 rounded-xl p-6 border border-amber-100 dark:border-amber-800/30"
          >
            <div class="flex items-center gap-3 mb-6">
              <div class="w-8 h-8 bg-amber-500 rounded-lg flex items-center justify-center">
                <FilePlus class="w-4 h-4 text-white" />
              </div>
              <h3 class="text-xl font-semibold text-gray-900 dark:text-white">File Import</h3>
            </div>

            <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
              <!-- Jira XML Import -->
              <div
                class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-6 h-full flex flex-col"
              >
                <div class="flex items-center gap-2 mb-4">
                  <div class="w-6 h-6 bg-blue-100 dark:bg-blue-900/30 rounded-md flex items-center justify-center">
                    <FileCode class="w-4 h-4 text-blue-600 dark:text-blue-400" />
                  </div>
                  <h4 class="text-lg font-medium text-gray-900 dark:text-white">
                    {$LL.importJiraXML()}
                  </h4>
                </div>
                <p class="text-gray-600 dark:text-gray-400 mb-4 flex-grow">Import stories from Jira XML export files</p>
                <div class="mt-auto">
                  <JiraImport handlePlanAdd={handleAdd} {notifications} />
                </div>
              </div>

              <!-- CSV Import -->
              <div
                class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-6 h-full flex flex-col"
              >
                <div class="flex items-center gap-2 mb-4">
                  <div class="w-6 h-6 bg-green-100 dark:bg-green-900/30 rounded-md flex items-center justify-center">
                    <FileSpreadsheet class="w-4 h-4 text-green-600 dark:text-green-400" />
                  </div>
                  <h4 class="text-lg font-medium text-gray-900 dark:text-white">
                    {$LL.importCsv()}
                  </h4>
                </div>

                <div class="mb-4 flex-grow">
                  <p class="text-gray-600 dark:text-gray-400 mb-3">
                    CSV file must include these fields (header row optional):
                  </p>
                  <div class="bg-gray-50 dark:bg-gray-900 rounded-md p-3 border border-gray-200 dark:border-gray-700">
                    <code class="text-gray-700 dark:text-gray-300 break-all">
                      Type,Title,ReferenceId,Link,Description,AcceptanceCriteria
                    </code>
                  </div>
                </div>

                <div class="mt-auto">
                  <CsvImport handlePlanAdd={handleAdd} {notifications} />
                </div>
              </div>
            </div>
          </div>
        </div>
      {/if}
    {/if}
  </div>
</Modal>
