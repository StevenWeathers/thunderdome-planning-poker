<script lang="ts">
  export let selectedPaymentPeriod = 'month';
  export let togglePaymentPeriod = () => {};
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
    @apply text-violet-200;
  }

  .toggle:hover {
    @apply bg-violet-800;
  }

  .toggle:focus {
    @apply z-10;
    @apply outline-none;
    @apply ring-2;
    @apply ring-white;
    @apply ring-offset-2;
    @apply ring-offset-violet-700;
  }

  .active {
    @apply border-violet-700;
    @apply bg-white;
    @apply text-violet-700;
    @apply shadow-sm;
  }

  .active:hover {
    @apply bg-violet-50;
  }
</style>

<div class="flex rounded-lg bg-violet-700 p-0.5">
  <button
    on:click="{togglePaymentPeriod}"
    type="button"
    class="toggle {selectedPaymentPeriod === 'month' ? 'active' : ''}"
  >
    Monthly billing
  </button>
  <button
    on:click="{togglePaymentPeriod}"
    type="button"
    class="toggle ml-0.5 {selectedPaymentPeriod === 'year' ? 'active' : ''}"
  >
    Yearly billing
  </button>
</div>
