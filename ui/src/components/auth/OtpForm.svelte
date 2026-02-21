<script lang="ts">
  import TextInput from '../forms/TextInput.svelte';
  import { ShieldIcon } from '@lucide/svelte';
  import LL from '../../i18n/i18n-svelte';

  interface Props {
    authMfa: (token: string) => void;
  }

  let { authMfa }: Props = $props();
  let otpToken = $state('');

  function handleSubmit(event: Event) {
    event.preventDefault();
    authMfa(otpToken);
  }
</script>

<form onsubmit={handleSubmit} class="space-y-6" name="authMfa">
  <div class="space-y-2">
    <label class="block text-sm font-medium text-gray-700 dark:text-gray-300" for="otp">
      {$LL.mfaTokenLabel()}
    </label>
    <TextInput
      bind:value={otpToken}
      placeholder="Enter code"
      id="otp"
      name="otp"
      required
      icon={ShieldIcon}
      inputmode="numeric"
      pattern="[0-9]*"
      autocomplete="one-time-code"
    />
  </div>

  <div class="pt-4">
    <button
      type="submit"
      class="w-full group relative flex justify-center py-3 px-4 border border-transparent text-lg font-medium rounded-lg text-white transition-all duration-300 transform hover:scale-105 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-purple-500 bg-gradient-to-r from-purple-500 to-indigo-500 hover:from-purple-600 hover:to-indigo-600 disabled:opacity-50 disabled:cursor-not-allowed"
    >
      <span class="absolute left-0 inset-y-0 flex items-center ps-3">
        <ShieldIcon class="h-5 w-5 text-purple-300 group-hover:text-purple-200" aria-hidden="true" />
      </span>
      {$LL.login()}
    </button>
  </div>
</form>
