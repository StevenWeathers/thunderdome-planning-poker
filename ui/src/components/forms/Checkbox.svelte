<script lang="ts">
  interface Props {
    label?: string;
    id?: string;
    name?: string;
    disabled?: boolean;
    value?: string | number;
    checked?: boolean | undefined;
  }

  let {
    label = '',
    id = '',
    name = '',
    disabled = false,
    value = '',
    checked = $bindable(false)
  }: Props = $props();
</script>

<div class="flex items-center">
  <div class="relative inline-block mr-2">
    <input
      type="checkbox"
      id="{id}"
      bind:checked="{checked}"
      value="{value}"
      class="hidden"
      name="{name}"
      disabled="{disabled}"
    />
    <label for="{id}" class="flex items-center cursor-pointer">
      <div
        class="w-6 h-6 border-2 border-indigo-400 rounded-md transition-all duration-200 ease-in-out flex items-center justify-center bg-white dark:bg-gray-800"
      >
        {#if checked}
          <svg
            class="w-4 h-4 text-white pointer-events-none"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="4"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <polyline points="20 6 9 17 4 12"></polyline>
          </svg>
        {/if}
      </div>
      <span class="ml-2 font-medium text-gray-700 dark:text-gray-300"
        >{label}</span
      >
    </label>
  </div>
</div>
