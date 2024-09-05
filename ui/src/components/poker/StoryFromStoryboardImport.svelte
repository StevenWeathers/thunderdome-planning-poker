<script lang="ts">
  import SelectInput from '../forms/SelectInput.svelte';
  import { user } from '../../stores';
  import LL from '../../i18n/i18n-svelte';
  import SolidButton from '../global/SolidButton.svelte';
  import { onMount } from 'svelte';

  export let eventTag;
  export let notifications;
  export let xfetch;
  export let handleImport;

  let selectedStoryboardIdx = '';
  let storyboards = [];
  let storyboard = {
    id: '',
    goals: [],
  };
  let selectedGoalIdx = '';

  function getStoryboards() {
    xfetch(`/api/users/${$user.id}/storyboards?limit=${9999}&offset=${0}`)
      .then(res => res.json())
      .then(function (result) {
        storyboards = result.data;
      })
      .catch(function () {
        notifications.danger($LL.getStoryboardsErrorMessage());
        eventTag('fetch_storyboards', 'engagement', 'failure');
      });
  }

  function getStoryboardStories() {
    if (selectedStoryboardIdx === '') {
      notifications.danger('Storyboard not selected');
      return;
    }
    const storyboardId = storyboards[selectedStoryboardIdx].id;
    xfetch(`/api/storyboards/${storyboardId}`)
      .then(res => res.json())
      .then(function (result) {
        storyboard = result.data;
        if (storyboard.goals.length === 1) {
          selectedGoalIdx = 0;
        }
      })
      .catch(function () {
        notifications.danger($LL.getStoryboardErrorMessage());
        eventTag('fetch_storyboard', 'engagement', 'failure');
      });
  }

  const importStory = (columnIdx, storyIdx) => () => {
    const selectedStory =
      storyboard.goals[selectedGoalIdx].columns[columnIdx].stories[storyIdx];

    handleImport({
      name: selectedStory.name,
      type: 'Story',
      link: selectedStory.link,
      description: selectedStory.content,
    });
  };

  onMount(() => {
    getStoryboards();
  });
</script>

<div class="mb-4">
  <SelectInput
    id="selectedStoryboard"
    bind:value="{selectedStoryboardIdx}"
    on:change="{getStoryboardStories}"
  >
    >
    <option value="" disabled>Select storyboard to import from</option>
    {#each storyboards as storyboard, idx}
      <option value="{idx}">{storyboard.name}</option>
    {/each}
  </SelectInput>
</div>

<div class="mb-4">
  <SelectInput
    id="selectedGoal"
    bind:value="{selectedGoalIdx}"
    on:change="{getStoryboardStories}"
    disabled="{selectedStoryboardIdx === ''}"
  >
    >
    <option value="" disabled>Select goal to import from</option>
    {#if selectedStoryboardIdx !== ''}
      {#each storyboard.goals as goal, idx}
        <option value="{idx}">{goal.name}</option>
      {/each}
    {/if}
  </SelectInput>
</div>

{#if selectedGoalIdx !== ''}
  <div class="mb-4">
    {#each storyboard.goals[selectedGoalIdx].columns as column, cIdx}
      {#each column.stories as story, sIdx}
        {#if story.name !== ''}
          <div
            class="p-2 w-full flex flex-wrap justify-between bg-gray-200 dark:bg-gray-900 dark:text-white rounded shadow mb-2 items-center"
          >
            <div>
              {story.name}
            </div>
            <div>
              <SolidButton onClick="{importStory(cIdx, sIdx)}"
                >Import
              </SolidButton>
            </div>
          </div>
        {/if}
      {/each}
    {/each}
  </div>
{/if}
