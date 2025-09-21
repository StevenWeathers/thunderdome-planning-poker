<script lang="ts">
  import { LoaderPinwheel } from 'lucide-svelte';
  interface Props {
    children?: import('svelte').Snippet;
  }

  let { children }: Props = $props();
</script>

<div class="w-full h-full fixed top-0 left-0 bg-white/75 dark:bg-gray-900/75 z-40">
  <div class="flex justify-center items-center mt-[50vh]">
    <div class="text-indigo-600 dark:text-indigo-300 text-center">
      <LoaderPinwheel class="animate-spin h-24 w-24 -mt-12 m-auto" />
      <div class="text-2xl font-rajdhani mt-6">
        {@render children?.()}
      </div>
    </div>
  </div>
</div>
