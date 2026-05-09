<script lang="ts">
  import SolidButton from '../global/SolidButton.svelte';
  import Modal from '../global/Modal.svelte';
  import LL from '../../i18n/i18n-svelte';
  import TextInput from '../forms/TextInput.svelte';
  import StoryColorSelector from './StoryColorSelector.svelte';
  import type { ColorLegend } from '../../types/storyboard';
  import { onMount } from 'svelte';

  interface Props {
    handleGoalAdd?: any;
    toggleAddGoal?: any;
    handleGoalRevision?: any;
    goalId?: string;
    goalName?: string;
    goalDefaultStoryColor?: string | null;
    colorLegend?: ColorLegend[];
  }

  let {
    handleGoalAdd = () => {},
    toggleAddGoal = () => {},
    handleGoalRevision = () => {},
    goalId = '',
    goalName = $bindable(''),
    goalDefaultStoryColor = null,
    colorLegend = [],
  }: Props = $props();

  let focusInput: any = $state();
  let selectedDefaultStoryColor = $state('');

  $effect(() => {
    selectedDefaultStoryColor = goalDefaultStoryColor ?? '';
  });

  function handleSubmit(event: Event) {
    event.preventDefault();

    const defaultStoryColor = selectedDefaultStoryColor === '' ? null : selectedDefaultStoryColor;

    if (goalId === '') {
      handleGoalAdd({
        name: goalName,
        defaultStoryColor,
      });
    } else {
      handleGoalRevision({
        goalId,
        name: goalName,
        defaultStoryColor,
      });
    }
    toggleAddGoal();
  }

  onMount(() => {
    focusInput?.focus();
  });
</script>

<Modal closeModal={toggleAddGoal} ariaLabel={$LL.modalAddStoryboardGoal()}>
  <form onsubmit={handleSubmit} name="addGoal">
    <div class="mb-4">
      <label class="block text-lg text-gray-700 dark:text-gray-300 font-bold mb-2" for="goalName">
        {$LL.storyboardGoalName()}
      </label>
      <TextInput
        id="goalName"
        bind:value={goalName}
        placeholder={$LL.storyboardGoalNamePlaceholder()}
        name="goalName"
        bind:this={focusInput}
      />
    </div>
    <div class="mb-4">
      <label class="block text-lg text-gray-700 dark:text-gray-300 font-bold mb-2" for="goalDefaultStoryColor">
        Default Story Color
      </label>
      <div id="goalDefaultStoryColor">
        <StoryColorSelector
          bind:value={selectedDefaultStoryColor}
          {colorLegend}
          allowNone={true}
          noneLabel="No default"
        />
      </div>
    </div>
    <div class="text-right">
      <div>
        <SolidButton type="submit">{$LL.save()}</SolidButton>
      </div>
    </div>
  </form>
</Modal>
