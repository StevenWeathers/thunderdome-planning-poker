<script lang="ts">
  import { MessageSquareMore } from '@lucide/svelte';

  interface Props {
    title?: string;
    description?: string;
  }
  let { title = 'No comments yet', description = 'Be the first to share your thoughts.' }: Props = $props();
</script>

<div class="text-center py-8">
  <div class="w-16 h-16 mx-auto mb-4 rounded-full bg-gray-100 dark:bg-gray-700 flex items-center justify-center">
    <MessageSquareMore class="w-8 h-8 text-gray-400 dark:text-gray-500" />
  </div>
  <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-2">{title}</h3>
  <p class="text-gray-500 dark:text-gray-400">{description}</p>
</div>
