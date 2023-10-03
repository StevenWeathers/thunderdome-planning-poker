<script lang="ts">
  import SolidButton from '../global/SolidButton.svelte';
  import Modal from '../global/Modal.svelte';
  import LL from '../../i18n/i18n-svelte';
  import TextInput from '../global/TextInput.svelte';

  export let handleGoalAdd = () => {};
  export let toggleAddGoal = () => {};
  export let handleGoalRevision = () => {};

  export let goalId = '';
  export let goalName = '';

  function handleSubmit(event) {
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
</script>

<Modal closeModal="{toggleAddGoal}">
  <form on:submit="{handleSubmit}" name="addGoal">
    <div class="mb-4">
      <label
        class="block text-sm text-gray-700 dark:text-gray-400 font-bold mb-2"
        for="goalName"
      >
        {$LL.storyboardGoalName()}
      </label>
      <TextInput
        id="goalName"
        bind:value="{goalName}"
        placeholder="{$LL.storyboardGoalNamePlaceholder()}"
        name="goalName"
      />
    </div>
    <div class="text-right">
      <div>
        <SolidButton type="submit">{$LL.save()}</SolidButton>
      </div>
    </div>
  </form>
</Modal>
