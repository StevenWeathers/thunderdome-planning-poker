<script lang="ts">
  import { Eye, EyeOff } from 'lucide-svelte';
  import TextInput from './TextInput.svelte';

  interface Props {
    value?: string;
    [key: string]: any
  }

  let { value = $bindable(''), ...rest }: Props = $props();

  let showPassword = $state(false);

  function togglePassword() {
    showPassword = !showPassword;
  }
</script>

<div class="relative">
  {#if showPassword}
    <TextInput
      type="text"
      autocomplete="current-password"
      bind:value="{value}"
      {...rest}
    />
  {:else}
    <TextInput
      type="password"
      autocomplete="current-password"
      bind:value="{value}"
      {...rest}
    />
  {/if}
  <button
    type="button"
    onclick={togglePassword}
    class="absolute top-3 right-3 text-gray-500 dark:text-gray-400 focus:outline-none hover:text-indigo-500"
  >
    {#if showPassword}
      <EyeOff size="{24}" />
    {:else}
      <Eye size="{24}" />
    {/if}
  </button>
</div>
