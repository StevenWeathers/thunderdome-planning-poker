<script lang="ts">
  import { Check } from "lucide-svelte";

  export let label: string = '';
  export let id: string = '';
  export let name: string = '';
  export let disabled: boolean = false;
  export let value: string | number = '';
  export let checked: boolean | undefined = false;
</script>

<div class="flex items-center">
  <label for="{id}" class="flex items-center cursor-pointer select-none">
    <input
      type="checkbox"
      id="{id}"
      bind:checked="{checked}"
      value="{value}"
      class="sr-only peer"
      name="{name}"
      disabled="{disabled}"
    />
    <div
      class="w-6 h-6 border-2 rounded-md transition-all duration-200 ease-in-out flex items-center justify-center bg-white dark:bg-gray-800 border-indigo-400 peer-focus-visible:outline peer-focus-visible:outline-2 peer-focus-visible:outline-offset-2 peer-focus-visible:outline-indigo-500 peer-disabled:opacity-50"
      aria-hidden="true"
    >
      {#if checked}
        <Check class="w-4 h-4 text-green-500 dark:text-lime-300 pointer-events-none" />
      {/if}
    </div>
    <span class="ms-2 font-medium text-gray-700 dark:text-gray-300 peer-disabled:opacity-60">{label}</span>
  </label>
</div>
