<script lang="ts">
  import SubscribeButton from '../../components/pricing/SubscribeButton.svelte';
  import { AppConfig } from '../../config.ts';
  import {
    Check,
    ShieldCheck,
    Zap,
    Users,
    Target,
    Calendar,
    Globe,
    Settings,
    BarChart3,
    Package,
    RefreshCw,
  } from '@lucide/svelte';
  import PayPeriodToggle from '../../components/pricing/PayPeriodToggle.svelte';
  import { onMount } from 'svelte';
  import { user } from '../../stores';
  import { appRoutes } from '../../config';
  import PlanBoxRecommended from '../../components/pricing/PlanBoxRecommended.svelte';
  import PlanBox from '../../components/pricing/PlanBox.svelte';
  import Github from '../../components/icons/Github.svelte';

  const { RepoURL } = AppConfig;

  let selectedPaymentPeriod = $state('month');

  function togglePaymentPeriod() {
    selectedPaymentPeriod = selectedPaymentPeriod === 'month' ? 'year' : 'month';
  }

  onMount(() => window.scrollTo(0, 0));
</script>

<div
  class="bg-gradient-to-br from-white to-slate-200 dark:from-gray-900 dark:to-gray-800 min-h-screen text-gray-900 dark:text-white font-sans transition-colors duration-300"
>
  <div class="container mx-auto px-4 py-16">
    <h1 class="text-5xl sm:text-6xl font-extrabold text-center mb-4 leading-tight">
      <span
        class="bg-clip-text text-transparent bg-gradient-to-r from-purple-600 to-pink-600 dark:from-purple-400 dark:to-pink-600"
      >
        Elevate Your Agile Workflow
      </span>
    </h1>
    <p class="text-xl text-center mb-12 text-gray-600 dark:text-gray-300">
      Choose the plan that fits your team's ambitions. Unlock premium features to streamline your agile process.
    </p>

    <PayPeriodToggle {selectedPaymentPeriod} {togglePaymentPeriod} />

    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
      <!-- Starter Plan -->
      <PlanBox>
        <h2 class="text-2xl font-bold mb-4 text-purple-600 dark:text-purple-300">STARTER</h2>
        <p class="text-5xl font-bold mb-4">Free</p>
        <p class="mb-6 text-gray-600 dark:text-gray-400">Perfect for small teams just getting started</p>
        {#if $user.rank && $user.rank !== 'GUEST'}
          <p class="w-full py-3 px-6 mb-8 text-green-600 dark:text-lime-400 text-center font-bold">
            Already registered.
          </p>
        {:else}
          <a
            href="{appRoutes.register}/subscription"
            class="w-full inline-block text-center bg-gray-200 dark:bg-gray-700 text-gray-900 dark:text-white font-bold py-3 px-6 rounded-full mb-8 hover:bg-gray-300 dark:hover:bg-gray-600 transition duration-300"
            >Register
          </a>
        {/if}
        <ul class="space-y-4">
          <li
            class="group flex items-center p-3 rounded-lg bg-purple-50 dark:bg-purple-900/20 border border-purple-100 dark:border-purple-800/50 hover:bg-purple-100 dark:hover:bg-purple-900/30 transition-all duration-200"
          >
            <div
              class="flex-shrink-0 w-8 h-8 flex items-center justify-center bg-purple-100 dark:bg-purple-800 rounded-full mr-3 group-hover:scale-110 transition-transform duration-200"
            >
              <Target class="w-4 h-4 text-purple-600 dark:text-purple-300" />
            </div>
            <span class="text-gray-700 dark:text-gray-300 font-medium">Poker Planning</span>
          </li>
          <li
            class="group flex items-center p-3 rounded-lg bg-purple-50 dark:bg-purple-900/20 border border-purple-100 dark:border-purple-800/50 hover:bg-purple-100 dark:hover:bg-purple-900/30 transition-all duration-200"
          >
            <div
              class="flex-shrink-0 w-8 h-8 flex items-center justify-center bg-purple-100 dark:bg-purple-800 rounded-full mr-3 group-hover:scale-110 transition-transform duration-200"
            >
              <RefreshCw class="w-4 h-4 text-purple-600 dark:text-purple-300" />
            </div>
            <span class="text-gray-700 dark:text-gray-300 font-medium">Sprint Retrospectives</span>
          </li>
          <li
            class="group flex items-center p-3 rounded-lg bg-purple-50 dark:bg-purple-900/20 border border-purple-100 dark:border-purple-800/50 hover:bg-purple-100 dark:hover:bg-purple-900/30 transition-all duration-200"
          >
            <div
              class="flex-shrink-0 w-8 h-8 flex items-center justify-center bg-purple-100 dark:bg-purple-800 rounded-full mr-3 group-hover:scale-110 transition-transform duration-200"
            >
              <Globe class="w-4 h-4 text-purple-600 dark:text-purple-300" />
            </div>
            <span class="text-gray-700 dark:text-gray-300 font-medium">Story Mapping</span>
          </li>
          <li
            class="group flex items-center p-3 rounded-lg bg-purple-50 dark:bg-purple-900/20 border border-purple-100 dark:border-purple-800/50 hover:bg-purple-100 dark:hover:bg-purple-900/30 transition-all duration-200"
          >
            <div
              class="flex-shrink-0 w-8 h-8 flex items-center justify-center bg-purple-100 dark:bg-purple-800 rounded-full mr-3 group-hover:scale-110 transition-transform duration-200"
            >
              <Calendar class="w-4 h-4 text-purple-600 dark:text-purple-300" />
            </div>
            <span class="text-gray-700 dark:text-gray-300 font-medium">Async Daily Stand-ups</span>
          </li>
          <li
            class="group flex items-center p-3 rounded-lg bg-purple-50 dark:bg-purple-900/20 border border-purple-100 dark:border-purple-800/50 hover:bg-purple-100 dark:hover:bg-purple-900/30 transition-all duration-200"
          >
            <div
              class="flex-shrink-0 w-8 h-8 flex items-center justify-center bg-purple-100 dark:bg-purple-800 rounded-full mr-3 group-hover:scale-110 transition-transform duration-200"
            >
              <Users class="w-4 h-4 text-purple-600 dark:text-purple-300" />
            </div>
            <span class="text-gray-700 dark:text-gray-300 font-medium">No User Limits</span>
          </li>
        </ul>
      </PlanBox>

      <!-- Individual Plan -->
      {#if AppConfig.Subscription.Individual.Enabled}
        <PlanBox>
          <h2 class="text-2xl font-bold mb-4 text-blue-600 dark:text-blue-300">Individual</h2>
          <p class="text-5xl font-bold mb-4">
            ${AppConfig.Subscription.Individual[`${selectedPaymentPeriod === 'month' ? 'MonthPrice' : 'YearPrice'}`]}
            &nbsp;<span class="text-2xl">/{selectedPaymentPeriod === 'month' ? 'mo' : 'yr'}</span>
          </p>
          <p class="mb-6 text-gray-600 dark:text-gray-400">Ideal for growing teams and professionals</p>
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
          <ul class="space-y-4">
            <li
              class="group flex items-center p-3 rounded-lg bg-gradient-to-r from-blue-50 to-purple-50 dark:from-blue-900/20 dark:to-purple-900/20 border border-blue-200 dark:border-blue-800/50 hover:from-blue-100 hover:to-purple-100 dark:hover:from-blue-900/30 dark:hover:to-purple-900/30 transition-all duration-200"
            >
              <div
                class="flex-shrink-0 w-8 h-8 flex items-center justify-center bg-gradient-to-r from-blue-100 to-purple-100 dark:from-blue-800 dark:to-purple-800 rounded-full mr-3 group-hover:scale-110 transition-transform duration-200"
              >
                <Zap class="w-4 h-4 text-blue-600 dark:text-blue-300" />
              </div>
              <span class="text-gray-700 dark:text-gray-300 font-bold">All Starter Features plus</span>
            </li>
            <li
              class="group flex items-center p-3 rounded-lg bg-blue-50 dark:bg-blue-900/20 border border-blue-100 dark:border-blue-800/50 hover:bg-blue-100 dark:hover:bg-blue-900/30 transition-all duration-200"
            >
              <div
                class="flex-shrink-0 w-8 h-8 flex items-center justify-center bg-blue-100 dark:bg-blue-800 rounded-full mr-3 group-hover:scale-110 transition-transform duration-200"
              >
                <Check class="w-4 h-4 text-blue-600 dark:text-blue-300" />
              </div>
              <span class="text-gray-700 dark:text-gray-300 font-medium">Import stories from Jira Cloud</span>
            </li>
            <li
              class="group flex items-center p-3 rounded-lg bg-blue-50 dark:bg-blue-900/20 border border-blue-100 dark:border-blue-800/50 hover:bg-blue-100 dark:hover:bg-blue-900/30 transition-all duration-200"
            >
              <div
                class="flex-shrink-0 w-8 h-8 flex items-center justify-center bg-blue-100 dark:bg-blue-800 rounded-full mr-3 group-hover:scale-110 transition-transform duration-200"
              >
                <Check class="w-4 h-4 text-blue-600 dark:text-blue-300" />
              </div>
              <span class="text-gray-700 dark:text-gray-300 font-medium"
                >Import stories from other sessions or even Storyboards</span
              >
            </li>
            <li
              class="group flex items-center p-3 rounded-lg bg-blue-50 dark:bg-blue-900/20 border border-blue-100 dark:border-blue-800/50 hover:bg-blue-100 dark:hover:bg-blue-900/30 transition-all duration-200"
            >
              <div
                class="flex-shrink-0 w-8 h-8 flex items-center justify-center bg-blue-100 dark:bg-blue-800 rounded-full mr-3 group-hover:scale-110 transition-transform duration-200"
              >
                <Check class="w-4 h-4 text-blue-600 dark:text-blue-300" />
              </div>
              <span class="text-gray-700 dark:text-gray-300 font-medium"
                >See previous check-in during today's stand-up</span
              >
            </li>
          </ul>
        </PlanBox>
      {/if}

      <!-- Team Plan -->
      {#if AppConfig.Subscription.Team.Enabled}
        <PlanBoxRecommended>
          <h2 class="text-2xl font-bold mb-4 text-purple-600 dark:text-purple-300">TEAM</h2>
          <p class="text-5xl font-bold mb-4">
            ${AppConfig.Subscription.Team[`${selectedPaymentPeriod === 'month' ? 'MonthPrice' : 'YearPrice'}`]}
            &nbsp;<span class="text-2xl">/{selectedPaymentPeriod === 'month' ? 'mo' : 'yr'}</span>
          </p>
          <p class="mb-6 text-gray-700 dark:text-gray-300">Empower your entire team with premium features</p>
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
          <ul class="space-y-4">
            <li
              class="group flex items-center p-3 rounded-lg bg-gradient-to-r from-purple-50 to-pink-50 dark:from-purple-900/20 dark:to-pink-900/20 border border-purple-200 dark:border-purple-800/50 hover:from-purple-100 hover:to-pink-100 dark:hover:from-purple-900/30 dark:hover:to-pink-900/30 transition-all duration-200 shadow-sm hover:shadow-md"
            >
              <div
                class="flex-shrink-0 w-8 h-8 flex items-center justify-center bg-gradient-to-r from-purple-100 to-pink-100 dark:from-purple-800 dark:to-pink-800 rounded-full mr-3 group-hover:scale-110 transition-transform duration-200"
              >
                <Zap class="w-4 h-4 text-purple-600 dark:text-purple-300" />
              </div>
              <span class="text-gray-700 dark:text-gray-300 font-bold">All Individual Features plus</span>
            </li>
            <li
              class="group flex items-center p-3 rounded-lg bg-purple-50 dark:bg-purple-900/20 border border-purple-100 dark:border-purple-800/50 hover:bg-purple-100 dark:hover:bg-purple-900/30 transition-all duration-200 shadow-sm hover:shadow-md"
            >
              <div
                class="flex-shrink-0 w-8 h-8 flex items-center justify-center bg-purple-100 dark:bg-purple-800 rounded-full mr-3 group-hover:scale-110 transition-transform duration-200"
              >
                <Settings class="w-4 h-4 text-purple-600 dark:text-purple-300" />
              </div>
              <span class="text-gray-700 dark:text-gray-300 font-medium"
                >Configure Team default settings for creating Retrospectives and Planning Poker sessions</span
              >
            </li>
            <li
              class="group flex items-center p-3 rounded-lg bg-purple-50 dark:bg-purple-900/20 border border-purple-100 dark:border-purple-800/50 hover:bg-purple-100 dark:hover:bg-purple-900/30 transition-all duration-200 shadow-sm hover:shadow-md"
            >
              <div
                class="flex-shrink-0 w-8 h-8 flex items-center justify-center bg-purple-100 dark:bg-purple-800 rounded-full mr-3 group-hover:scale-110 transition-transform duration-200"
              >
                <BarChart3 class="w-4 h-4 text-purple-600 dark:text-purple-300" />
              </div>
              <span class="text-gray-700 dark:text-gray-300 font-medium">Custom Poker Estimation Scales</span>
            </li>
            <li
              class="group flex items-center p-3 rounded-lg bg-purple-50 dark:bg-purple-900/20 border border-purple-100 dark:border-purple-800/50 hover:bg-purple-100 dark:hover:bg-purple-900/30 transition-all duration-200 shadow-sm hover:shadow-md"
            >
              <div
                class="flex-shrink-0 w-8 h-8 flex items-center justify-center bg-purple-100 dark:bg-purple-800 rounded-full mr-3 group-hover:scale-110 transition-transform duration-200"
              >
                <RefreshCw class="w-4 h-4 text-purple-600 dark:text-purple-300" />
              </div>
              <span class="text-gray-700 dark:text-gray-300 font-medium">Custom Retrospective Templates</span>
            </li>
            <li
              class="group flex items-center p-3 rounded-lg bg-purple-50 dark:bg-purple-900/20 border border-purple-100 dark:border-purple-800/50 hover:bg-purple-100 dark:hover:bg-purple-900/30 transition-all duration-200 shadow-sm hover:shadow-md"
            >
              <div
                class="flex-shrink-0 w-8 h-8 flex items-center justify-center bg-purple-100 dark:bg-purple-800 rounded-full mr-3 group-hover:scale-110 transition-transform duration-200"
              >
                <Check class="w-4 h-4 text-purple-600 dark:text-purple-300" />
              </div>
              <span class="text-gray-700 dark:text-gray-300 font-medium"
                >Review open action items from previous retrospectives</span
              >
            </li>
            {#if AppConfig.FeatureProject}
              <li
                class="group flex items-center p-3 rounded-lg bg-purple-50 dark:bg-purple-900/20 border border-purple-100 dark:border-purple-800/50 hover:bg-purple-100 dark:hover:bg-purple-900/30 transition-all duration-200 shadow-sm hover:shadow-md"
              >
                <div
                  class="flex-shrink-0 w-8 h-8 flex items-center justify-center bg-purple-100 dark:bg-purple-800 rounded-full mr-3 group-hover:scale-110 transition-transform duration-200"
                >
                  <Package class="w-4 h-4 text-purple-600 dark:text-purple-300" />
                </div>
                <span class="text-gray-700 dark:text-gray-300 font-medium"
                  >Create Projects to manage your team's work</span
                >
              </li>
            {/if}
          </ul>
        </PlanBoxRecommended>
      {/if}

      <!-- Organization Plan -->
      {#if AppConfig.Subscription.Organization.Enabled}
        <PlanBox>
          <h2 class="text-2xl font-bold mb-4 text-green-600 dark:text-green-300">Organization</h2>
          <p class="text-5xl font-bold mb-4">
            ${AppConfig.Subscription.Organization[`${selectedPaymentPeriod === 'month' ? 'MonthPrice' : 'YearPrice'}`]}
            &nbsp;<span class="text-2xl">/{selectedPaymentPeriod === 'month' ? 'mo' : 'yr'}</span>
          </p>
          <p class="mb-6 text-gray-600 dark:text-gray-400">Cost savings for companies with multiple teams</p>
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
          <ul class="space-y-4">
            <li
              class="group flex items-center p-3 rounded-lg bg-gradient-to-r from-green-50 to-blue-50 dark:from-green-900/20 dark:to-blue-900/20 border border-green-200 dark:border-green-800/50 hover:from-green-100 hover:to-blue-100 dark:hover:from-green-900/30 dark:hover:to-blue-900/30 transition-all duration-200"
            >
              <div
                class="flex-shrink-0 w-8 h-8 flex items-center justify-center bg-gradient-to-r from-green-100 to-blue-100 dark:from-green-800 dark:to-blue-800 rounded-full mr-3 group-hover:scale-110 transition-transform duration-200"
              >
                <Zap class="w-4 h-4 text-green-600 dark:text-green-300" />
              </div>
              <span class="text-gray-700 dark:text-gray-300 font-bold">All Team Features plus</span>
            </li>
            <li
              class="group flex items-center p-3 rounded-lg bg-green-50 dark:bg-green-900/20 border border-green-100 dark:border-green-800/50 hover:bg-green-100 dark:hover:bg-green-900/30 transition-all duration-200"
            >
              <div
                class="flex-shrink-0 w-8 h-8 flex items-center justify-center bg-green-100 dark:bg-green-800 rounded-full mr-3 group-hover:scale-110 transition-transform duration-200"
              >
                <Users class="w-4 h-4 text-green-600 dark:text-green-300" />
              </div>
              <span class="text-gray-700 dark:text-gray-300 font-medium">Organization User Management and Roles</span>
            </li>
            <li
              class="group flex items-center p-3 rounded-lg bg-green-50 dark:bg-green-900/20 border border-green-100 dark:border-green-800/50 hover:bg-green-100 dark:hover:bg-green-900/30 transition-all duration-200"
            >
              <div
                class="flex-shrink-0 w-8 h-8 flex items-center justify-center bg-green-100 dark:bg-green-800 rounded-full mr-3 group-hover:scale-110 transition-transform duration-200"
              >
                <Settings class="w-4 h-4 text-green-600 dark:text-green-300" />
              </div>
              <span class="text-gray-700 dark:text-gray-300 font-medium"
                >Configure Organization default settings for creating Team Retrospectives and Planning Poker sessions</span
              >
            </li>
            {#if AppConfig.FeatureProject}
              <li
                class="group flex items-center p-3 rounded-lg bg-green-50 dark:bg-green-900/20 border border-green-100 dark:border-green-800/50 hover:bg-green-100 dark:hover:bg-green-900/30 transition-all duration-200"
              >
                <div
                  class="flex-shrink-0 w-8 h-8 flex items-center justify-center bg-green-100 dark:bg-green-800 rounded-full mr-3 group-hover:scale-110 transition-transform duration-200"
                >
                  <Package class="w-4 h-4 text-green-600 dark:text-green-300" />
                </div>
                <span class="text-gray-700 dark:text-gray-300 font-medium">
                  Create Projects to manage your organization's work
                </span>
              </li>
              <li
                class="group flex items-center p-3 rounded-lg bg-green-50 dark:bg-green-900/20 border border-green-100 dark:border-green-800/50 hover:bg-green-100 dark:hover:bg-green-900/30 transition-all duration-200"
              >
                <div
                  class="flex-shrink-0 w-8 h-8 flex items-center justify-center bg-green-100 dark:bg-green-800 rounded-full mr-3 group-hover:scale-110 transition-transform duration-200"
                >
                  <Package class="w-4 h-4 text-green-600 dark:text-green-300" />
                </div>
                <span class="text-gray-700 dark:text-gray-300 font-medium">
                  Collaborate with clients, vendors, and stakeholders outside your organization in shared projects
                </span>
              </li>
            {/if}
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
              <ShieldCheck class="w-8 h-8 inline mr-3" />
              Self-Hosted Thunderdome
            </h2>
            <p class="text-xl text-indigo-100 mb-8">
              Deploy Thunderdome on your own infrastructure with complete data ownership. Perfect for enterprises,
              government agencies, and teams with strict security, compliance, or data residency requirements.
            </p>
            <div class="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
              <ul class="space-y-3">
                <li class="flex items-start text-white">
                  <Check class="w-5 h-5 mr-3 mt-1 flex-shrink-0 text-green-300" />
                  <span>Complete data sovereignty and privacy control</span>
                </li>
                <li class="flex items-start text-white">
                  <Check class="w-5 h-5 mr-3 mt-1 flex-shrink-0 text-green-300" />
                  <span>Integration with existing SSO and LDAP systems</span>
                </li>
                <li class="flex items-start text-white">
                  <Check class="w-5 h-5 mr-3 mt-1 flex-shrink-0 text-green-300" />
                  <span>No subscription fees or user limits</span>
                </li>
              </ul>
              <ul class="space-y-3">
                <li class="flex items-start text-white">
                  <Check class="w-5 h-5 mr-3 mt-1 flex-shrink-0 text-green-300" />
                  <span>Single binary or containerized deployment</span>
                </li>
                <li class="flex items-start text-white">
                  <Check class="w-5 h-5 mr-3 mt-1 flex-shrink-0 text-green-300" />
                  <span>Open source with Apache 2.0 license</span>
                </li>
                <li class="flex items-start text-white">
                  <Check class="w-5 h-5 mr-3 mt-1 flex-shrink-0 text-green-300" />
                  <span>Consider donating to support ongoing development ❤️</span>
                </li>
              </ul>
            </div>
            <div class="flex flex-wrap gap-4">
              <a
                href="{RepoURL}/blob/main/docs/INSTALLATION.md"
                target="_blank"
                class="bg-white text-indigo-600 hover:bg-indigo-50 font-bold py-3 px-6 rounded-full transition duration-300 shadow-lg hover:shadow-xl transform hover:-translate-y-1"
                >📚 Installation Guide</a
              >
              <a
                href={RepoURL}
                target="_blank"
                class="border-2 border-white text-white hover:bg-white hover:text-indigo-600 font-bold py-3 px-6 rounded-full transition duration-300"
                ><Github class="me-2" width={22} /> View on GitHub</a
              >
            </div>
          </div>
          <div class="md:w-1/2 p-12 flex items-center justify-center">
            <div class="relative w-full max-w-sm">
              <!-- Server Icon Illustration -->
              <div class="relative">
                <div
                  class="absolute inset-0 bg-gradient-to-br from-white/20 to-white/5 rounded-2xl transform rotate-6 scale-105"
                ></div>
                <div class="relative bg-white/15 backdrop-blur-sm rounded-2xl p-8 border border-white/20">
                  <div class="text-center">
                    <div class="inline-flex items-center justify-center w-20 h-20 bg-white/20 rounded-full mb-4">
                      <Settings class="w-10 h-10 text-white" />
                    </div>
                    <h3 class="text-2xl font-bold text-white mb-2">Your Infrastructure</h3>
                    <p class="text-indigo-100 mb-4">Complete control & customization</p>
                    <div class="space-y-2">
                      <div class="flex items-center justify-center space-x-2 text-white/80">
                        <div class="w-2 h-2 bg-green-400 rounded-full"></div>
                        <span>Docker Ready</span>
                      </div>
                      <div class="flex items-center justify-center space-x-2 text-white/80">
                        <div class="w-2 h-2 bg-blue-400 rounded-full"></div>
                        <span>Cloud Native</span>
                      </div>
                      <div class="flex items-center justify-center space-x-2 text-white/80">
                        <div class="w-2 h-2 bg-purple-400 rounded-full"></div>
                        <span>Secure by Design</span>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>
