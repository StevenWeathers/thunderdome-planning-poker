<script lang="ts">
  import PageLayout from '../../components/PageLayout.svelte';
  import HollowButton from '../../components/HollowButton.svelte';
  import { AppConfig, appRoutes } from '../../config';
  import { user } from '../../stores';

  const { SubscriptionCheckoutLink } = AppConfig;
</script>

<PageLayout>
  <div
    class="space-y-8 lg:w-3/4 lg:grid lg:grid-cols-2 sm:gap-6 xl:gap-10 lg:space-y-0 mx-auto"
  >
    <div
      class="flex flex-col p-6 mx-auto max-w-lg text-center text-gray-900 bg-white rounded-lg border border-gray-100 shadow dark:border-gray-600 xl:p-8 dark:bg-gray-800 dark:text-white"
    >
      <h3 class="mb-4 text-2xl font-semibold">Starter</h3>
      <p class="font-light text-gray-500 sm:text-lg dark:text-gray-400">
        All the primary features at no cost.
      </p>
      <div class="flex justify-center items-baseline my-8">
        <span class="mr-2 text-5xl font-extrabold">Free</span>
      </div>
      <ul role="list" class="mb-8 space-y-4 text-left">
        <li class="flex items-center space-x-3">
          <svg
            class="flex-shrink-0 w-5 h-5 text-green-500 dark:text-green-400"
            fill="currentColor"
            viewBox="0 0 20 20"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              fill-rule="evenodd"
              d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
              clip-rule="evenodd"></path>
          </svg>
          <span>Agile Poker Planning</span>
        </li>
        <li class="flex items-center space-x-3">
          <svg
            class="flex-shrink-0 w-5 h-5 text-green-500 dark:text-green-400"
            fill="currentColor"
            viewBox="0 0 20 20"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              fill-rule="evenodd"
              d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
              clip-rule="evenodd"></path>
          </svg>
          <span>Poker planning import stories from Jira XML or CSV files</span>
        </li>
        <li class="flex items-center space-x-3">
          <svg
            class="flex-shrink-0 w-5 h-5 text-green-500 dark:text-green-400"
            fill="currentColor"
            viewBox="0 0 20 20"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              fill-rule="evenodd"
              d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
              clip-rule="evenodd"></path>
          </svg>
          <span>Agile Sprint Retrospective</span>
        </li>
        <li class="flex items-center space-x-3">
          <svg
            class="flex-shrink-0 w-5 h-5 text-green-500 dark:text-green-400"
            fill="currentColor"
            viewBox="0 0 20 20"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              fill-rule="evenodd"
              d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
              clip-rule="evenodd"></path>
          </svg>
          <span>Team Check-ins or asynchronous stand ups</span>
        </li>
        <li class="flex items-center space-x-3">
          <svg
            class="flex-shrink-0 w-5 h-5 text-green-500 dark:text-green-400"
            fill="currentColor"
            viewBox="0 0 20 20"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              fill-rule="evenodd"
              d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
              clip-rule="evenodd"></path>
          </svg>
          <span>Agile Story Mapping (Storybaords)</span>
        </li>
        <li class="flex items-center space-x-3">
          <svg
            class="flex-shrink-0 w-5 h-5 text-green-500 dark:text-green-400"
            fill="currentColor"
            viewBox="0 0 20 20"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              fill-rule="evenodd"
              d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
              clip-rule="evenodd"></path>
          </svg>
          <span>Team organization</span>
        </li>
      </ul>
    </div>

    <div
      class="flex flex-col p-6 mx-auto max-w-lg text-center text-gray-900 bg-white rounded-lg border border-gray-100 shadow dark:border-gray-600 xl:p-8 dark:bg-gray-800 dark:text-white"
    >
      <h3 class="mb-4 text-2xl font-semibold">Premium</h3>
      <p class="font-light text-gray-500 sm:text-lg dark:text-gray-400">
        Primary features plus enhanced ones.
      </p>
      <div class="flex justify-center items-baseline my-8">
        <span class="mr-2 text-5xl font-extrabold">$5</span>
        <span class="text-gray-500 dark:text-gray-400" dark:text-gray-400
          >/month</span
        >
      </div>

      <ul role="list" class="mb-8 space-y-4 text-left">
        <li class="flex items-center space-x-3">
          <svg
            class="flex-shrink-0 w-5 h-5 text-green-500 dark:text-green-400"
            fill="currentColor"
            viewBox="0 0 20 20"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              fill-rule="evenodd"
              d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
              clip-rule="evenodd"></path>
          </svg>
          <span
            >Poker planning import stories from Jira Cloud using JQL search</span
          >
        </li>
      </ul>
      <ul role="list" class="mb-8 space-y-4 text-left">
        <li class="flex items-center space-x-3">
          <svg
            class="flex-shrink-0 w-5 h-5 text-green-500 dark:text-green-400"
            fill="currentColor"
            viewBox="0 0 20 20"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              fill-rule="evenodd"
              d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
              clip-rule="evenodd"></path>
          </svg>
          <span>Future premium features</span>
        </li>
      </ul>
      <ul role="list" class="mb-8 space-y-4 text-left">
        <li class="flex items-center space-x-3">
          <svg
            class="flex-shrink-0 w-5 h-5 text-green-500 dark:text-green-400"
            fill="currentColor"
            viewBox="0 0 20 20"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              fill-rule="evenodd"
              d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
              clip-rule="evenodd"></path>
          </svg>
          <span>Supporting the developer and hosting costs</span>
        </li>
      </ul>
      {#if !$user.name || $user.rank === 'GUEST'}
        <p class="bg-yellow-thunder text-gray-900 px-4 py-6 rounded">
          Must <a
            class="underline font-bold"
            href="{appRoutes.register}/subscription">register</a
          >
          or
          <a class="underline font-bold" href="{appRoutes.login}/subscription"
            >login</a
          > to subscribe
        </p>
      {:else if $user.subscribed}
        <p class="bg-lime-300 text-green-800 px-4 py-6 rounded font-bold">
          Already subscribed, thank you!
        </p>
      {:else}
        <HollowButton
          color="green"
          href="{SubscriptionCheckoutLink}?prefilled_email={$user.email}&client_reference_id={$user.id}"
        >
          <span class="block py-2 uppercase">Subscribe today</span>
        </HollowButton>
      {/if}
    </div>
  </div>
</PageLayout>
