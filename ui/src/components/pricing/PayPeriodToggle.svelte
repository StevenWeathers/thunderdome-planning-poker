<script lang="ts">
  export let selectedPaymentPeriod = 'month';
  export let togglePaymentPeriod = () => {};

  $: yearPayPeriod = selectedPaymentPeriod === 'year';

  function handleChange() {
    togglePaymentPeriod();
  }
</script>

<style>
  .toggle {
    @apply relative;
    @apply whitespace-nowrap;
    @apply rounded-md;
    @apply py-2;
    @apply px-6;
    @apply text-sm;
    @apply font-medium;
    @apply border;
    @apply border-transparent;
    @apply text-indigo-200;
  }

  .toggle:hover {
    @apply bg-indigo-800;
  }

  .toggle:focus {
    @apply z-10;
    @apply outline-none;
    @apply ring-2;
    @apply ring-white;
    @apply ring-offset-2;
    @apply ring-offset-indigo-700;
  }

  .active {
    @apply border-indigo-700;
    @apply bg-white;
    @apply text-indigo-700;
    @apply shadow-sm;
  }

  .active:hover {
    @apply bg-indigo-50;
  }
</style>

<div class="flex justify-center items-center space-x-4 mb-12">
  <span class="text-lg font-semibold text-gray-600 dark:text-gray-300"
    >Monthly</span
  >
  <label class="relative inline-flex items-center cursor-pointer">
    <input
      type="checkbox"
      id="billingToggle"
      class="sr-only peer"
      bind:checked="{yearPayPeriod}"
      on:change="{handleChange}"
    />
    <div
      class="w-14 h-7 bg-gray-300 dark:bg-gray-700 rounded-full peer peer-checked:after:translate-x-full after:content-[''] after:absolute after:top-0.5 after:left-[4px] after:bg-white after:rounded-full after:h-6 after:w-6 after:transition-all peer-checked:bg-purple-600"
    ></div>
  </label>
  <span class="text-lg font-semibold text-gray-600 dark:text-gray-300"
    >Yearly <span class="text-purple-600 dark:text-purple-400"
      >(2 Months Free!)</span
    ></span
  >
</div>
