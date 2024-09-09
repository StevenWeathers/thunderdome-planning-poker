<script lang="ts">
  import { user } from '../../stores.ts';
  import { appRoutes } from '../../config';
  import { ArrowRight } from 'lucide-svelte';

  export let subscribeLink = '';
  export let planEnabled = false;
</script>

<div>
  {#if planEnabled}
    {#if !$user.name || !$user.rank || $user.rank === 'GUEST'}
      <p class="bg-yellow-thunder text-gray-900 p-2 rounded">
        <a class="underline font-bold" href="{appRoutes.login}/subscription"
          >Login</a
        >
        or
        <a class="underline font-bold" href="{appRoutes.register}/subscription"
          >Register</a
        > to subscribe
      </p>
    {:else if $user.subscribed}
      <p class="text-green-600 dark:text-lime-400 p-2 font-bold">
        Already subscribed, thank you!
      </p>
    {:else}
      <a
        href="{subscribeLink}?prefilled_email={$user.email}&client_reference_id={$user.id}"
        class="flex items-center mt-auto text-white bg-indigo-500 border-0 py-2 px-4 w-full focus:outline-none hover:bg-indigo-600 rounded"
      >
        Subscribe Today
        <ArrowRight class="h-4 w-4 inline-block" />
      </a>
      <p class="text-xs text-gray-500 dark:text-gray-400 mt-3">
        Payments processed by Stripe.
      </p>
    {/if}
  {:else}
    <p class="text-blue-500 dark:text-cyan-400 p-2 font-bold text-lg">
      Coming soon!
    </p>
  {/if}
</div>
