<script lang="ts">
  import { user } from '../../stores.ts';
  import { appRoutes } from '../../config';
  interface Props {
    children?: import('svelte').Snippet;
  }

  let { children }: Props = $props();
</script>

{#if !$user.name || !$user.rank || $user.rank === 'GUEST'}
  <p
    class="w-full py-3 px-6 mb-8 text-blue-600 dark:text-sky-400 text-center font-bold"
  >
    <a class="underline font-bold" href="{appRoutes.login}/subscription"
      >Login</a
    >
    to subscribe.
  </p>
{:else if $user.subscribed}
  <p
    class="w-full py-3 px-6 mb-8 text-green-600 dark:text-lime-400 text-center font-bold"
  >
    Already subscribed.
  </p>
{:else}
  {@render children?.()}
{/if}
