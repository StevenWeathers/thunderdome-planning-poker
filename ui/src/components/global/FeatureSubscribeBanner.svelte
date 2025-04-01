<script lang="ts">
  import { Bell } from 'lucide-svelte';
  import { appRoutes } from '../../config';

  interface Props {
    salesPitch?: string;
  }

  let { salesPitch = '' }: Props = $props();
</script>

<div
  class="bg-gradient-to-r from-blue-50 to-indigo-50 dark:from-indigo-800 dark:to-indigo-900 rounded-lg shadow-md p-4 max-w-4xl mx-auto my-4 border border-blue-100 dark:border-indigo-800 transition-all duration-300"
>
  <div class="flex items-center space-x-3 text-gray-700 dark:text-gray-200">
    <Bell class="w-5 h-5 text-blue-500 dark:text-blue-400 animate-pulse" />
    <p class="text-sm sm:text-base">
      <span class="font-semibold">New:</span>
      {salesPitch}
      <a
        class="text-blue-600 dark:text-blue-400 hover:underline cursor-pointer ml-1 underline"
        href="{appRoutes.subscriptionPricing}"
        target="_blank"
      >
        Subscribe now
      </a>
    </p>
  </div>
</div>
