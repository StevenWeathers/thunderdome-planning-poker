<script lang="ts">
  import Checkmark from '../../components/pricing/Checkmark.svelte';
  import SubscribeButton from '../../components/pricing/SubscribeButton.svelte';
  import Crossmark from '../../components/pricing/Crossmark.svelte';
  import { user } from '../../stores.ts';
  import { AppConfig, appRoutes } from '../../config.ts';
  import ArrowRight from '../../components/icons/ArrowRight.svelte';
  import PayPeriodToggle from '../../components/pricing/PayPeriodToggle.svelte';

  let selectedPaymentPeriod = 'month';

  function togglePaymentPeriod() {
    selectedPaymentPeriod =
      selectedPaymentPeriod === 'month' ? 'year' : 'month';
  }

  const Plans = [
    {
      Name: 'Starter',
      Enabled: true,
      Popular: false,
      Price: {
        Free: true,
        Month: 'Free',
        Year: 'Free',
        Subtext: 'All the primary features at no cost forever.',
      },
      Subscribe: {
        MonthLink: '',
        YearLink: '',
      },
      Features: [
        {
          Name: 'Primary Features',
          Included: true,
          Header: true,
          ColorClass: 'text-blue-600 dark:text-sky-400',
        },
        {
          Name: 'Poker Planning',
          Included: true,
          Header: false,
        },
        {
          Name: 'Import stories from Jira XML or CSV files',
          Included: true,
          Header: false,
        },
        {
          Name: 'Sprint Retrospective',
          Included: true,
          Header: false,
        },
        {
          Name: 'Story Mapping (Storyboards)',
          Included: true,
          Header: false,
        },
        {
          Name: 'Async Daily Stand-up (Team Checkins)',
          Included: true,
          Header: false,
        },
        {
          Name: 'Organize by Organization, Departments & Teams',
          Included: true,
          Header: false,
        },
        {
          Name: 'No user limits',
          Included: true,
          Header: false,
        },
        {
          Name: 'Premium Features',
          Included: true,
          Header: true,
          ColorClass: 'text-green-600 dark:text-lime-400',
        },
        {
          Name: 'Import Poker Stories from Jira',
          Included: false,
          Header: false,
        },
        {
          Name: 'Future premium features',
          Included: false,
          Header: false,
        },
        {
          Name: 'Compensating development efforts',
          Included: false,
          Header: false,
        },
      ],
    },
    {
      Name: 'Individual',
      Enabled: AppConfig.Subscription.Individual.Enabled,
      Popular: false,
      Price: {
        Free: false,
        Month: AppConfig.Subscription.Individual.MonthPrice,
        Year: AppConfig.Subscription.Individual.YearPrice,
        Subtext: 'Premium Features for your account only.',
      },
      Subscribe: {
        MonthLink: AppConfig.Subscription.Individual.MonthCheckoutLink,
        YearLink: AppConfig.Subscription.Individual.YearCheckoutLink,
      },
      Features: [
        {
          Name: 'Primary Features',
          Included: true,
          Header: true,
          ColorClass: 'text-blue-600 dark:text-sky-400',
        },
        {
          Name: 'Poker Planning',
          Included: true,
          Header: false,
        },
        {
          Name: 'Import stories from Jira XML or CSV files',
          Included: true,
          Header: false,
        },
        {
          Name: 'Sprint Retrospective',
          Included: true,
          Header: false,
        },
        {
          Name: 'Story Mapping (Storyboards)',
          Included: true,
          Header: false,
        },
        {
          Name: 'Async Daily Stand-up (Team Checkins)',
          Included: true,
          Header: false,
        },
        {
          Name: 'Organize by Organization, Departments & Teams',
          Included: true,
          Header: false,
        },
        {
          Name: 'No user limits',
          Included: true,
          Header: false,
        },
        {
          Name: 'Premium Features',
          Included: true,
          Header: true,
          ColorClass: 'text-green-600 dark:text-lime-400',
        },
        {
          Name: 'Import Poker Stories from Jira',
          Included: true,
          Header: false,
        },
        {
          Name: 'Future premium features',
          Included: true,
          Header: false,
        },
        {
          Name: 'Compensating development efforts',
          Included: true,
          Header: false,
        },
      ],
    },
    {
      Name: 'Team',
      Enabled: AppConfig.Subscription.Team.Enabled,
      Popular: false,
      Price: {
        Free: false,
        Month: AppConfig.Subscription.Team.MonthPrice,
        Year: AppConfig.Subscription.Team.YearPrice,
        Subtext: 'Premium Features for all Team members.',
      },
      Subscribe: {
        MonthLink: AppConfig.Subscription.Team.MonthCheckoutLink,
        YearLink: AppConfig.Subscription.Team.YearCheckoutLink,
      },
      Features: [
        {
          Name: 'Primary Features',
          Included: true,
          Header: true,
          ColorClass: 'text-blue-600 dark:text-sky-400',
        },
        {
          Name: 'Poker Planning',
          Included: true,
          Header: false,
        },
        {
          Name: 'Import stories from Jira XML or CSV files',
          Included: true,
          Header: false,
        },
        {
          Name: 'Sprint Retrospective',
          Included: true,
          Header: false,
        },
        {
          Name: 'Story Mapping (Storyboards)',
          Included: true,
          Header: false,
        },
        {
          Name: 'Async Daily Stand-up (Team Checkins)',
          Included: true,
          Header: false,
        },
        {
          Name: 'Organize by Organization, Departments & Teams',
          Included: true,
          Header: false,
        },
        {
          Name: 'No user limits',
          Included: true,
          Header: false,
        },
        {
          Name: 'Premium Features',
          Included: true,
          Header: true,
          ColorClass: 'text-green-600 dark:text-lime-400',
        },
        {
          Name: 'Import Poker Stories from Jira',
          Included: true,
          Header: false,
        },
        {
          Name: 'Future premium features',
          Included: true,
          Header: false,
        },
        {
          Name: 'Compensating development efforts',
          Included: true,
          Header: false,
        },
      ],
    },
    {
      Name: 'Organization',
      Enabled: AppConfig.Subscription.Organization.Enabled,
      Popular: false,
      Price: {
        Free: false,
        Month: AppConfig.Subscription.Organization.MonthPrice,
        Year: AppConfig.Subscription.Organization.YearPrice,
        Subtext: 'Premium Features for all Organization members.',
      },
      Subscribe: {
        MonthLink: AppConfig.Subscription.Organization.MonthCheckoutLink,
        YearLink: AppConfig.Subscription.Organization.YearCheckoutLink,
      },
      Features: [
        {
          Name: 'Primary Features',
          Included: true,
          Header: true,
          ColorClass: 'text-blue-600 dark:text-sky-400',
        },
        {
          Name: 'Poker Planning',
          Included: true,
          Header: false,
        },
        {
          Name: 'Import stories from Jira XML or CSV files',
          Included: true,
          Header: false,
        },
        {
          Name: 'Sprint Retrospective',
          Included: true,
          Header: false,
        },
        {
          Name: 'Story Mapping (Storyboards)',
          Included: true,
          Header: false,
        },
        {
          Name: 'Async Daily Stand-up (Team Checkins)',
          Included: true,
          Header: false,
        },
        {
          Name: 'Organize by Organization, Departments & Teams',
          Included: true,
          Header: false,
        },
        {
          Name: 'No user limits',
          Included: true,
          Header: false,
        },
        {
          Name: 'Premium Features',
          Included: true,
          Header: true,
          ColorClass: 'text-green-600 dark:text-lime-400',
        },
        {
          Name: 'Import Poker Stories from Jira',
          Included: true,
          Header: false,
        },
        {
          Name: 'Future premium features',
          Included: true,
          Header: false,
        },
        {
          Name: 'Compensating development efforts',
          Included: true,
          Header: false,
        },
      ],
    },
  ];
</script>

<section class="relative bg-indigo-600 py-12">
  <div class="relative mx-auto max-w-2xl px-6 text-center lg:max-w-7xl lg:px-8">
    <h1 class="text-4xl font-bold tracking-tight text-white sm:text-6xl">
      <span class="block lg:inline">Simple pricing,</span>
      <span class="block lg:inline">no commitment.</span>
    </h1>
    <p class="mt-4 text-xl text-indigo-100">
      Everything you need, nothing you don't. Pick a plan that best suits your
      needs.
    </p>
  </div>

  <h2 class="sr-only">Plans</h2>

  <div class="relative flex justify-center mt-8">
    <PayPeriodToggle
      selectedPaymentPeriod="{selectedPaymentPeriod}"
      togglePaymentPeriod="{togglePaymentPeriod}"
    />
  </div>
</section>

<section class="bg-white dark:bg-transparent text-gray-700 dark:text-gray-300">
  <div class="container px-5 pb-24 pt-12 mx-auto flex flex-wrap">
    <div class="lg:w-1/4 mt-48 hidden lg:block">
      <div
        class="mt-px border-t border-gray-300 dark:border-gray-700 border-b border-l rounded-tl-lg rounded-bl-lg overflow-hidden"
      >
        {#each Plans[0].Features as feature, idx}
          <p
            class="{idx % 2
              ? 'text-gray-900 dark:text-gray-300'
              : 'bg-gray-100 dark:bg-gray-800 text-gray-900 dark:text-gray-300'} h-12 px-4 flex items-center justify-start"
          >
            {#if feature.Header}
              <span class="text-lg font-semibold {feature.ColorClass}"
                >{feature.Name}</span
              >
            {:else}
              {feature.Name}
            {/if}
          </p>
        {/each}
      </div>
    </div>
    <div
      class="flex lg:w-3/4 w-full flex-wrap lg:border border-gray-300 dark:border-gray-700 rounded-lg"
    >
      {#each Plans as plan, pidx}
        <div
          class="lg:w-1/4 lg:mt-px w-full mb-10 lg:mb-0 border-2 rounded-lg relative {pidx %
          2
            ? 'border-gray-300 dark:border-gray-700 lg:border-transparent dark:lg:border-transparent lg:rounded-none'
            : 'border-gray-300 dark:border-gray-700'} {plan.Popular
            ? 'border-indigo-500'
            : ''}"
        >
          {#if plan.Popular}
            <span
              class="bg-indigo-500 text-white px-3 py-1 tracking-widest text-xs absolute right-0 top-0 rounded-bl uppercase"
              >Popular</span
            >
          {/if}
          <div
            class="px-2 text-center h-48 flex flex-col items-center justify-center"
          >
            <h3 class="tracking-widest uppercase">{plan.Name}</h3>
            <h2
              class="text-5xl text-gray-600 dark:text-gray-100 font-medium leading-none mb-4 mt-2"
            >
              {#if plan.Price.Free}
                Free
              {:else}
                {#if selectedPaymentPeriod === 'month'}${plan.Price
                    .Month}{:else}
                  ${plan.Price.Year}{/if}
                <span class="text-gray-600 dark:text-gray-300 text-base ml-1"
                  >{#if selectedPaymentPeriod === 'month'}/mo{:else}/yr{/if}</span
                >
              {/if}
            </h2>
            <span class="text-sm text-gray-600 dark:text-gray-300"
              >{plan.Price.Subtext}</span
            >
          </div>
          <div>
            {#each plan.Features as feature, idx}
              <p
                class="{idx % 2
                  ? 'text-gray-600 dark:text-gray-300'
                  : 'bg-gray-100 dark:bg-gray-800'} h-12 px-2 flex items-center justify-start text-center justify-center {idx ===
                0
                  ? '-mt-px border-t border-gray-300 dark:border-gray-700'
                  : ''}"
              >
                {#if feature.Header}
                  <span
                    class="text-lg font-semibold {feature.ColorClass} lg:hidden"
                    >{feature.Name}</span
                  >
                {:else}
                  {#if feature.Included}
                    <Checkmark />
                  {:else}
                    <Crossmark />
                  {/if}
                  <span class="lg:hidden ml-2">{feature.Name}</span>
                {/if}
              </p>
            {/each}
          </div>
          <div
            class="border-t border-gray-300 dark:border-gray-700 p-6 text-center rounded-bl-lg"
          >
            {#if plan.Price.Free}
              {#if $user.rank && $user.rank !== 'GUEST'}
                <p class="text-green-600 dark:text-lime-400 p-2 font-bold">
                  Already registered.
                </p>
              {:else}
                <a
                  href="{appRoutes.register}/subscription"
                  class="flex items-center mt-auto text-white bg-indigo-500 border-0 py-2 px-4 w-full focus:outline-none hover:bg-indigo-600 rounded"
                  >Register
                  <ArrowRight />
                </a>
              {/if}
            {:else}
              <SubscribeButton
                planEnabled="{plan.Enabled}"
                subscribeLink="{selectedPaymentPeriod === 'month'
                  ? plan.Subscribe.MonthLink
                  : plan.Subscribe.YearLink}"
              />
            {/if}
          </div>
        </div>
      {/each}
    </div>
  </div>
</section>
