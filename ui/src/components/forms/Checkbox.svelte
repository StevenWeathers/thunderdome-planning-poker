<script lang="ts">
  export let label = '';
  export let id = '';
  export let name = '';
  export let disabled: boolean = false;
  export let value: string | number = '';
  export let checked: boolean | undefined = undefined;
</script>

<div class="flex items-center mb-4">
  <input
    id="{id}"
    type="checkbox"
    value="{value}"
    bind:checked="{checked}"
    name="{name}"
    disabled="{disabled}"
    class="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-blue-500 dark:focus:ring-blue-600 dark:ring-offset-gray-800 focus:ring-2 dark:bg-gray-700 dark:border-gray-600"
  />
  <label
    for="{id}"
    class="ms-2 text-sm font-medium text-gray-900 dark:text-gray-300"
    >{label}</label
  >
</div>
