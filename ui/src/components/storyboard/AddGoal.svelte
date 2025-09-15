<script lang="ts">
  import SolidButton from '../global/SolidButton.svelte';
  import Modal from '../global/Modal.svelte';
  import LL from '../../i18n/i18n-svelte';
  import TextInput from '../forms/TextInput.svelte';
  import { onMount } from 'svelte';


  interface Props {
    handleGoalAdd?: any;
    toggleAddGoal?: any;
    handleGoalRevision?: any;
    goalId?: string;
    goalName?: string;
  }

  let {
    handleGoalAdd = () => {},
    toggleAddGoal = () => {},
    handleGoalRevision = () => {},
    goalId = '',
    goalName = $bindable('')
  }: Props = $props();

  let focusInput: any = $state();

  function handleSubmit(event: Event) {
    event.preventDefault();

    if (goalId === '') {
      handleGoalAdd(goalName);
    } else {
      handleGoalRevision({
        goalId,
        name: goalName,
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
      <label
        class="block text-lg text-gray-700 dark:text-gray-300 font-bold mb-2"
        for="goalName"
      >
        {$LL.storyboardGoalName()}
      </label>
      <TextInput
        id="goalName"
        bind:value="{goalName}"
        placeholder={$LL.storyboardGoalNamePlaceholder()}
        name="goalName"
        bind:this={focusInput}
      />
    </div>
    <div class="text-right">
      <div>
        <SolidButton type="submit">{$LL.save()}</SolidButton>
      </div>
    </div>
  </form>
</Modal>
