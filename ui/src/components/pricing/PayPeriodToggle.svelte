<script lang="ts">
  import { run } from 'svelte/legacy';

  interface Props {
    selectedPaymentPeriod?: string;
    togglePaymentPeriod?: any;
  }

  let { selectedPaymentPeriod = 'month', togglePaymentPeriod = () => {} }: Props = $props();

  let yearPayPeriod = $state(false);
  run(() => {
    yearPayPeriod = selectedPaymentPeriod === 'year';
  });

  function handleChange() {
    togglePaymentPeriod();
  }
</script>

<div class="flex justify-center items-center space-x-4 mb-12">
  <span class="text-lg font-semibold text-gray-600 dark:text-gray-300"
    >Monthly</span
  >
  <label class="relative inline-flex items-center cursor-pointer">
    <input
      type="checkbox"
      id="billingToggle"
      class="sr-only peer"
      bind:checked={yearPayPeriod}
      onchange={handleChange}
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
