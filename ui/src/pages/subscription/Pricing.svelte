<script lang="ts">
  import SubscribeButton from '../../components/pricing/SubscribeButton.svelte';
  import { AppConfig } from '../../config.ts';
  import { Check, ShieldCheck } from 'lucide-svelte';
  import PayPeriodToggle from '../../components/pricing/PayPeriodToggle.svelte';
  import { onMount } from 'svelte';
  import { user } from '../../stores';
  import { appRoutes } from '../../config';
  import PlanBoxRecommended from '../../components/pricing/PlanBoxRecommended.svelte';
  import PlanBox from '../../components/pricing/PlanBox.svelte';

  const { RepoURL } = AppConfig;

  let selectedPaymentPeriod = 'month';

  function togglePaymentPeriod() {
    selectedPaymentPeriod =
      selectedPaymentPeriod === 'month' ? 'year' : 'month';
  }

  onMount(() => window.scrollTo(0, 0));
</script>

<div
  class="bg-gradient-to-br from-white to-slate-200 dark:from-gray-900 dark:to-gray-800 min-h-screen text-gray-900 dark:text-white font-sans transition-colors duration-300"
>
  <div class="container mx-auto px-4 py-16">
    <h1
      class="text-5xl sm:text-6xl font-extrabold text-center mb-4 leading-tight"
    >
      <span
        class="bg-clip-text text-transparent bg-gradient-to-r from-purple-600 to-pink-600 dark:from-purple-400 dark:to-pink-600"
      >
        Elevate Your Agile Workflow
      </span>
    </h1>
    <p class="text-xl text-center mb-12 text-gray-600 dark:text-gray-300">
      Choose the plan that fits your team's ambitions. Unlock premium features
      to streamline your agile process.
    </p>

    <PayPeriodToggle
      selectedPaymentPeriod="{selectedPaymentPeriod}"
      togglePaymentPeriod="{togglePaymentPeriod}"
    />

    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
      <!-- Starter Plan -->
      <PlanBox>
        <h2
          class="text-2xl font-bold mb-4 text-purple-600 dark:text-purple-300"
        >
          STARTER
        </h2>
        <p class="text-5xl font-bold mb-4">Free</p>
        <p class="mb-6 text-gray-600 dark:text-gray-400">
          Perfect for small teams just getting started
        </p>
        {#if $user.rank && $user.rank !== 'GUEST'}
          <p
            class="w-full py-3 px-6 mb-8 text-green-600 dark:text-lime-400 text-center font-bold"
          >
            Already registered.
          </p>
        {:else}
          <a
            href="{appRoutes.register}/subscription"
            class="w-full inline-block text-center bg-gray-200 dark:bg-gray-700 text-gray-900 dark:text-white font-bold py-3 px-6 rounded-full mb-8 hover:bg-gray-300 dark:hover:bg-gray-600 transition duration-300"
            >Register
          </a>
        {/if}
        <ul class="space-y-3">
          <li class="flex items-center">🚀 Poker Planning</li>
          <li class="flex items-center">🔄 Sprint Retrospectives</li>
          <li class="flex items-center">🗺️ Story Mapping</li>
          <li class="flex items-center">📅 Async Daily Stand-ups</li>
          <li class="flex items-center">👥 No User Limits</li>
        </ul>
      </PlanBox>

      <!-- Individual Plan -->
      {#if AppConfig.Subscription.Individual.Enabled}
        <PlanBox>
          <h2 class="text-2xl font-bold mb-4 text-blue-600 dark:text-blue-300">
            Individual
          </h2>
          <p class="text-5xl font-bold mb-4">
            ${AppConfig.Subscription.Individual[
              `${
                selectedPaymentPeriod === 'month' ? 'MonthPrice' : 'YearPrice'
              }`
            ]}
            &nbsp;<span class="text-2xl"
              >/{selectedPaymentPeriod === 'month' ? 'mo' : 'yr'}</span
            >
          </p>
          <p class="mb-6 text-gray-600 dark:text-gray-400">
            Ideal for growing teams and professionals
          </p>
          <SubscribeButton>
            <a
              class="w-full inline-block text-center bg-gradient-to-r from-blue-500 to-purple-600 text-white font-bold py-3 px-6 rounded-full mb-8 hover:from-blue-600 hover:to-purple-700 transition duration-300"
              href="{selectedPaymentPeriod === 'month'
                ? AppConfig.Subscription.Individual.MonthCheckoutLink
                : AppConfig.Subscription.Individual
                    .YearCheckoutLink}?prefilled_email={$user.email}&client_reference_id={$user.id}"
              target="_blank"
            >
              Subscribe Now
            </a>
          </SubscribeButton>
          <ul class="space-y-3">
            <li class="flex items-center font-bold">✨ All Starter Features</li>
            <li class="flex items-center">🔗 Import stories from Jira Cloud</li>
            <li class="flex items-center">
              🔗 Import stories from other sessions or even Storyboards
            </li>
            <li class="flex items-center">
              ✅ See previous check-in during today's stand-up
            </li>
          </ul>
        </PlanBox>
      {/if}

      <!-- Team Plan -->
      {#if AppConfig.Subscription.Team.Enabled}
        <PlanBoxRecommended>
          <h2
            class="text-2xl font-bold mb-4 text-purple-600 dark:text-purple-300"
          >
            TEAM
          </h2>
          <p class="text-5xl font-bold mb-4">
            ${AppConfig.Subscription.Team[
              `${
                selectedPaymentPeriod === 'month' ? 'MonthPrice' : 'YearPrice'
              }`
            ]}
            &nbsp;<span class="text-2xl"
              >/{selectedPaymentPeriod === 'month' ? 'mo' : 'yr'}</span
            >
          </p>
          <p class="mb-6 text-gray-700 dark:text-gray-300">
            Empower your entire team with premium features
          </p>
          <SubscribeButton>
            <a
              class="w-full inline-block text-center bg-gradient-to-r from-purple-500 to-pink-500 text-white font-bold py-3 px-6 rounded-full mb-8 hover:from-purple-600 hover:to-pink-600 transition duration-300 transform hover:-translate-y-1"
              href="{selectedPaymentPeriod === 'month'
                ? AppConfig.Subscription.Team.MonthCheckoutLink
                : AppConfig.Subscription.Team
                    .YearCheckoutLink}?prefilled_email={$user.email}&client_reference_id={$user.id}"
              target="_blank"
            >
              Subscribe Now
            </a>
          </SubscribeButton>
          <ul class="space-y-3">
            <li class="flex items-center font-bold">
              🌟 All Individual Features
            </li>
            <li class="flex items-center">
              ⚙️ Configure default settings for creating Retrospectives and
              Planning Poker sessions within the subscribed team
            </li>
            <li class="flex items-center">📊 Custom Poker Estimation Scales</li>
            <li class="flex items-center">🔄 Custom Retrospective Templates</li>
            <li class="flex items-center">
              🔄 Review open action items from previous retrospectives
            </li>
          </ul>
        </PlanBoxRecommended>
      {/if}

      <!-- Organization Plan -->
      {#if AppConfig.Subscription.Organization.Enabled}
        <PlanBox>
          <h2
            class="text-2xl font-bold mb-4 text-green-600 dark:text-green-300"
          >
            Organization
          </h2>
          <p class="text-5xl font-bold mb-4">
            ${AppConfig.Subscription.Organization[
              `${
                selectedPaymentPeriod === 'month' ? 'MonthPrice' : 'YearPrice'
              }`
            ]}
            &nbsp;<span class="text-2xl"
              >/{selectedPaymentPeriod === 'month' ? 'mo' : 'yr'}</span
            >
          </p>
          <p class="mb-6 text-gray-600 dark:text-gray-400">
            Cost savings for companies with multiple teams
          </p>
          <SubscribeButton>
            <a
              class="w-full inline-block text-center bg-gradient-to-r from-green-500 to-blue-500 text-white font-bold py-3 px-6 rounded-full mb-8 hover:from-green-600 hover:to-blue-600 transition duration-300"
              href="{selectedPaymentPeriod === 'month'
                ? AppConfig.Subscription.Organization.MonthCheckoutLink
                : AppConfig.Subscription.Organization
                    .YearCheckoutLink}?prefilled_email={$user.email}&client_reference_id={$user.id}"
              target="_blank"
            >
              Subscribe Now
            </a>
          </SubscribeButton>
          <ul class="space-y-3">
            <li class="flex items-center font-bold">💼 All Team Features</li>
            <li class="flex items-center">
              ⚙️ Configure default settings for creating Retrospectives and
              Planning Poker sessions within the subscribed organization
            </li>
          </ul>
        </PlanBox>
      {/if}
    </div>

    <div class="mt-24 mb-16">
      <div
        class="bg-gradient-to-r from-indigo-500 to-purple-600 dark:from-indigo-700 dark:to-purple-800 rounded-3xl overflow-hidden shadow-2xl"
      >
        <div class="md:flex">
          <div class="md:w-1/2 p-12">
            <h2 class="text-4xl font-bold text-white mb-6">
              <ShieldCheck class="w-8 h-8 inline" />
              Self-Hosted Thunderdome
            </h2>
            <p class="text-xl text-indigo-100 mb-8">
              Deploy Thunderdome on your own infrastructure. Perfect for teams
              that require complete control over their data and customization.
            </p>
            <ul class="space-y-4 mb-8">
              <li class="flex items-center text-white">
                <Check class="inline" />
                Full Control Over Your Data
              </li>
              <li class="flex items-center text-white">
                <Check class="inline" />
                Customizable to Your Needs
              </li>
            </ul>
            <div class="flex space-x-4">
              <a
                href="{RepoURL}#running-in-production"
                target="_blank"
                class="bg-white text-indigo-600 hover:bg-indigo-100 font-bold py-3 px-6 rounded-full transition duration-300"
                >Get Started</a
              >
            </div>
          </div>
          <div class="md:w-1/2 p-12 flex items-center justify-center">
            <div class="relative w-full h-full">
              <div
                class="absolute inset-0 bg-white opacity-10 rounded-full transform rotate-12 scale-110"
              ></div>
              <!--              <img-->
              <!--                      src="/api/placeholder/600/400"-->
              <!--                      alt="Self-hosted Thunderdome"-->
              <!--                      class="relative z-10 w-full h-auto rounded-lg shadow-lg"-->
              <!--              />-->
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>
