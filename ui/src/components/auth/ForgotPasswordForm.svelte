<script lang="ts">
  import TextInput from '../forms/TextInput.svelte';
  import { Mail } from '@lucide/svelte';
  import LL from '../../i18n/i18n-svelte';

  interface Props {
    toggleForgotPassword: () => void;
    onSubmit: (email: string) => void;
  }

  let { toggleForgotPassword, onSubmit }: Props = $props();
  let resetEmail = $state('');

  function handlerSubmit(e: Event) {
    e.preventDefault();
    onSubmit(resetEmail);
  }
</script>

<form onsubmit={handlerSubmit} class="space-y-6 max-w-md" name="resetPassword">
  <div class="mb-4">
    <h2
      class="font-rajdhani uppercase text-4xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-purple-600 to-pink-600 dark:from-purple-500 dark:to-pink-500 mb-2"
    >
      {$LL.forgotPassword()}
    </h2>
    <p class="font-light text-gray-600 dark:text-gray-300">
      {$LL.forgotPasswordSubtext()}
    </p>
  </div>

  <TextInput
    data-testid="resetemail"
    bind:value={resetEmail}
    placeholder={$LL.enterYourEmail()}
    id="yourResetEmail"
    name="yourResetEmail"
    type="email"
    required
    icon={Mail}
  />

  <div class="flex justify-between items-center">
    <button
      type="button"
      class="font-medium text-purple-600 dark:text-purple-400 hover:text-purple-500 dark:hover:text-purple-300 transition-all duration-300"
      onclick={toggleForgotPassword}
    >
      {$LL.returnToLogin()}
    </button>
    <button
      type="submit"
      class="px-6 py-2 text-lg font-medium rounded-lg text-white transition-all duration-300 transform hover:scale-105 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-purple-500 bg-gradient-to-r from-purple-500 to-indigo-500 hover:from-purple-600 hover:to-indigo-600 disabled:opacity-50 disabled:cursor-not-allowed"
    >
      {$LL.resetPassword()}
    </button>
  </div>
</form>
