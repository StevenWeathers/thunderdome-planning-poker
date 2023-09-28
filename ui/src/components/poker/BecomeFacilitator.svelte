<script lang="ts">
  import Modal from '../Modal.svelte';
  import LL from '../../i18n/i18n-svelte';
  import SolidButton from '../SolidButton.svelte';
  import TextInput from '../TextInput.svelte';

  export let toggleBecomeLeader = () => {};
  export let handleBecomeLeader = () => {};

  let leaderCode = '';

  function handleSubmit(e) {
    e.preventDefault();

    handleBecomeLeader(leaderCode);
  }
</script>

<Modal
  closeModal="{toggleBecomeLeader}"
  widthClasses="md:w-2/3 lg:w-3/5 xl:w-1/2"
>
  <form on:submit="{handleSubmit}" name="becomeLeader">
    <div class="mb-4">
      <label
        class="block text-gray-700 dark:text-gray-400 font-bold mb-2"
        for="leaderCode"
      >
        {$LL.leaderPasscode()}
      </label>
      <div class="control">
        <TextInput
          name="leaderCode"
          bind:value="{leaderCode}"
          id="leaderCode"
        />
      </div>
    </div>

    <div class="text-right">
      <SolidButton type="submit">{$LL.save()}</SolidButton>
    </div>
  </form>
</Modal>
