<script lang="ts">
  import TextInput from '../forms/TextInput.svelte';
  import SolidButton from './SolidButton.svelte';
  import LL from '../../i18n/i18n-svelte';

  export let handleSubmit: (passcode) => {};
  export let submitText: string = 'Join';

  let joinPasscode = '';

  function onSubmit(event) {
    event.preventDefault();
    handleSubmit(joinPasscode);
  }
</script>

<div
  class="w-full h-full fixed top-0 left-0 bg-white/80 dark:bg-gray-900/80 z-50"
>
  <div class="flex justify-center items-center content-center h-full">
    <div class="w-full md:w-1/2 lg:w-1/3">
      <form
        on:submit="{onSubmit}"
        class="bg-white dark:bg-gray-800 border border-indigo-500 shadow-lg rounded-lg p-6 mb-4"
        name="joinPasscodeForm"
      >
        <div class="mb-4">
          <label
            class="block text-gray-700 dark:text-gray-400 font-bold mb-2"
            for="passCode"
          >
            {$LL.passCodeRequired()}
          </label>
          <TextInput
            bind:value="{joinPasscode}"
            placeholder="{$LL.enterPasscode()}"
            id="passCode"
            name="passCode"
            type="password"
            required
          />
        </div>

        <div class="text-right">
          <SolidButton type="submit">{submitText}</SolidButton>
        </div>
      </form>
    </div>
  </div>
</div>
