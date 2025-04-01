<script lang="ts">
  import { AppConfig, appRoutes } from '../config';
  import { user } from '../stores';
  import { validateUserIsRegistered } from '../validationUtils';
  import { Github, Users, Zap } from 'lucide-svelte';
  import LL from '../i18n/i18n-svelte';
  import Countries from '../components/user/Countries.svelte';
  import BrowserMock from '../components/global/BrowserMock.svelte';
  import { onMount } from 'svelte';

  interface Props {
    xfetch: any;
  }

  let { xfetch }: Props = $props();

  const {
    ShowActiveCountries,
    PathPrefix,
    SubscriptionsEnabled,
    SubscriptionCheckoutLink,
    RepoURL,
  } = AppConfig;

  let isRegisteredUser = $derived($user && !!$user.id && validateUserIsRegistered($user));

  const slogans = [
    'Empower Your Agile Teams',
    'Release the Thunder in Your Agile Process',
    'Storm Through Your Teams Sprints',
    'Electrifying Your Agile Experience',
    'The Arena Where Agile Teams Excel',
    'Bringing the Energy Back to Agile',
    'Harness the Storm, Master the Sprint',
    'Amplify Your Teams Agile Potential',
  ];

  let randomSlogan = $derived(slogans[Math.floor(Math.random() * slogans.length)]);

  onMount(() => window.scrollTo(0, 0));
</script>

<svelte:head>
  <title>{$LL.appName()} - {$LL.appSubtitle()}</title>
</svelte:head>

<main class="bg-gray-100 dark:bg-gray-900">
  <header class="bg-gradient-to-r from-blue-600 to-indigo-700 text-white">
    <div class="container mx-auto px-4 py-16">
      <div class="max-w-7xl mx-auto text-center">
        <h1
          class="text-4xl sm:text-5xl lg:text-6xl font-bold mb-6 leading-tight"
        >
          <span
            class="block bg-clip-text text-transparent bg-gradient-to-r from-yellow-thunder to-orange-500"
          >
            Thunderdome
          </span>
          {randomSlogan}
        </h1>
        <p class="max-w-4xl mx-auto text-xl sm:text-2xl text-blue-100">
          Elevate your agile practices, foster seamless collaboration, and
          unlock your team's full potential with our innovative suite of tools.
        </p>
      </div>
    </div>
  </header>

  <section class="bg-white dark:bg-gray-800 py-20">
    <div class="container mx-auto px-4">
      <div class="flex flex-col md:flex-row items-center justify-between">
        <div class="md:w-1/2 md:pr-8 mb-8 md:mb-0">
          <div class="title-line bg-yellow-thunder"></div>
          <h2
            class="text-4xl font-semibold font-rajdhani uppercase dark:text-white mb-6"
          >
            Agile Poker Planning
          </h2>
          <ul
            class="list-disc pl-5 space-y-2 text-gray-700 dark:text-gray-300 mb-8"
          >
            <li>Streamline sprint planning with collaborative estimation</li>
            <li>Customize estimation scales to fit your team's needs</li>
            <li>Real-time voting and result visualization</li>
          </ul>
          <a
            href="{appRoutes.games}"
            class="group relative inline-flex items-center justify-center p-0.5 mb-2 mr-2 overflow-hidden font-medium text-gray-900 rounded-lg group bg-gradient-to-br from-purple-600 to-blue-500 group-hover:from-purple-600 group-hover:to-blue-500 hover:text-white dark:text-white focus:ring-4 focus:outline-none focus:ring-blue-300 dark:focus:ring-blue-800"
          >
            <span
              class="relative px-5 py-2.5 transition-all ease-in duration-75 bg-white dark:bg-gray-900 rounded-md group-hover:bg-opacity-0"
            >
              {$user.id ? $LL.battleCreate() : 'Learn More'}
            </span>
          </a>
        </div>
        <div class="md:w-1/2">
          <BrowserMock>
            <img
              class="rounded-b-lg hidden dark:block"
              src="{PathPrefix}/img/previews/poker_game.png"
              alt="{$LL.appPreviewAlt()}"
            />
            <img
              class="rounded-b-lg dark:hidden"
              src="{PathPrefix}/img/previews/poker_game_light.png"
              alt="{$LL.appPreviewAlt()}"
            />
          </BrowserMock>
        </div>
      </div>
    </div>
  </section>

  <section class="bg-gray-100 dark:bg-gray-900 py-20">
    <div class="container mx-auto px-4">
      <div
        class="flex flex-col md:flex-row-reverse items-center justify-between"
      >
        <div class="md:w-1/2 md:pl-8 mb-8 md:mb-0">
          <div class="title-line bg-yellow-thunder"></div>
          <h2
            class="text-4xl font-semibold font-rajdhani uppercase dark:text-white mb-6"
          >
            Sprint Retrospectives
          </h2>
          <ul
            class="list-disc pl-5 space-y-2 text-gray-700 dark:text-gray-300 mb-8"
          >
            <li>Facilitate engaging and productive retrospectives</li>
            <li>Choose from customizable templates or create your own</li>
            <li>Track action items and improvement initiatives</li>
          </ul>
          <a
            href="{appRoutes.retros}"
            class="group relative inline-flex items-center justify-center p-0.5 mb-2 mr-2 overflow-hidden font-medium text-gray-900 rounded-lg group bg-gradient-to-br from-purple-600 to-blue-500 group-hover:from-purple-600 group-hover:to-blue-500 hover:text-white dark:text-white focus:ring-4 focus:outline-none focus:ring-blue-300 dark:focus:ring-blue-800"
          >
            <span
              class="relative px-5 py-2.5 transition-all ease-in duration-75 bg-white dark:bg-gray-900 rounded-md group-hover:bg-opacity-0"
            >
              {$user.id ? 'Create Retro' : 'Learn More'}
            </span>
          </a>
        </div>
        <div class="md:w-1/2">
          <BrowserMock>
            <img
              class="rounded-b-lg hidden dark:block"
              src="{PathPrefix}/img/previews/retro.png"
              alt="Sprint Retrospectives Preview"
            />
            <img
              class="rounded-b-lg dark:hidden"
              src="{PathPrefix}/img/previews/retro_light.png"
              alt="Sprint Retrospectives Preview"
            />
          </BrowserMock>
        </div>
      </div>
    </div>
  </section>

  <section class="bg-white dark:bg-gray-800 py-20">
    <div class="container mx-auto px-4">
      <div class="flex flex-col md:flex-row items-center justify-between">
        <div class="md:w-1/2 md:pr-8 mb-8 md:mb-0">
          <div class="title-line bg-yellow-thunder"></div>
          <h2
            class="text-4xl font-semibold font-rajdhani uppercase dark:text-white mb-6"
          >
            Story Mapping
          </h2>
          <ul
            class="list-disc pl-5 space-y-2 text-gray-700 dark:text-gray-300 mb-8"
          >
            <li>Drag-and-drop interface for easy story organization</li>
            <li>
              Customizable columns and swimlanes for user activities and
              iterations
            </li>
            <li>Create, edit, and prioritize user stories within the map</li>
            <li>Color-coding and tagging for improved story categorization</li>
          </ul>
          <a
            href="{appRoutes.storyboards}"
            class="group relative inline-flex items-center justify-center p-0.5 mb-2 mr-2 overflow-hidden font-medium text-gray-900 rounded-lg group bg-gradient-to-br from-purple-600 to-blue-500 group-hover:from-purple-600 group-hover:to-blue-500 hover:text-white dark:text-white focus:ring-4 focus:outline-none focus:ring-blue-300 dark:focus:ring-blue-800"
          >
            <span
              class="relative px-5 py-2.5 transition-all ease-in duration-75 bg-white dark:bg-gray-900 rounded-md group-hover:bg-opacity-0"
            >
              {$user.id ? 'Create Storyboard' : 'Learn More'}
            </span>
          </a>
        </div>
        <div class="md:w-1/2">
          <BrowserMock>
            <img
              class="rounded-b-lg hidden dark:block"
              src="{PathPrefix}/img/previews/storyboard.png"
              alt="Story Mapping Preview"
            />
            <img
              class="rounded-b-lg dark:hidden"
              src="{PathPrefix}/img/previews/storyboard_light.png"
              alt="Story Mapping Preview"
            />
          </BrowserMock>
        </div>
      </div>
    </div>
  </section>

  <section class="bg-gray-100 dark:bg-gray-900 py-20">
    <div class="container mx-auto px-4">
      <div
        class="flex flex-col md:flex-row-reverse items-center justify-between"
      >
        <div class="md:w-1/2 md:pl-8 mb-8 md:mb-0">
          <div class="title-line bg-yellow-thunder"></div>
          <h2
            class="text-4xl font-semibold font-rajdhani uppercase dark:text-white mb-6"
          >
            Team Checkins
          </h2>
          <p class="text-lg dark:text-gray-300 mb-8">
            Streamline your team's agile stand-up. Focus on blockers and
            critical details instead of spending time discussing yesterday's and
            today's tasks.
          </p>
        </div>
        <div class="md:w-1/2">
          <img
            class="rounded-b-lg hidden dark:block"
            src="{PathPrefix}/img/previews/checkin.png"
            alt="Team Checkins Preview"
          />
          <img
            class="rounded-b-lg dark:hidden"
            src="{PathPrefix}/img/previews/checkin_light.png"
            alt="Team Checkins Preview"
          />
        </div>
      </div>
    </div>
  </section>

  <section class="bg-indigo-600 text-white py-20">
    <div class="container mx-auto px-4 text-center">
      <h2 class="text-4xl font-bold mb-6 font-rajdhani uppercase">
        Why Choose Thunderdome?
      </h2>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
        <div
          class="bg-white dark:bg-gray-800 rounded-lg p-6 text-gray-800 dark:text-white"
        >
          <div class="text-indigo-500 text-4xl mb-4">
            <Zap class="h-12 w-12 mx-auto" />
          </div>
          <h3 class="text-xl font-semibold mb-2">Boost Productivity</h3>
          <p>
            Streamline your agile processes and eliminate bottlenecks to deliver
            value faster.
          </p>
        </div>
        <div
          class="bg-white dark:bg-gray-800 rounded-lg p-6 text-gray-800 dark:text-white"
        >
          <div class="text-indigo-500 text-4xl mb-4">
            <Users class="h-12 w-12 mx-auto" />
          </div>
          <h3 class="text-xl font-semibold mb-2">Enhanced Collaboration</h3>
          <p>
            Foster teamwork and communication with our suite of collaborative
            tools.
          </p>
        </div>
        <div
          class="bg-white dark:bg-gray-800 rounded-lg p-6 text-gray-800 dark:text-white"
        >
          <div class="text-indigo-500 text-4xl mb-4">
            <Github class="h-12 w-12 mx-auto" />
          </div>
          <h3 class="text-xl font-semibold mb-2">Open Source</h3>
          <p>
            Choose between our <a
              href="{appRoutes.subscriptionPricing}"
              class="text-indigo-500 dark:text-indigo-400 hover:underline"
              >premium</a
            >
            hosted solution or
            <a
              href="{RepoURL}/blob/main/docs/INSTALLATION.md"
              class="text-indigo-500 dark:text-indigo-400 hover:underline"
              >self-host Thunderdome</a
            >.
          </p>
        </div>
      </div>
    </div>
  </section>

  {#if ShowActiveCountries}
    <section class="bg-slate-100 dark:bg-gray-900 py-20">
      <div class="container mx-auto px-4">
        <Countries xfetch={xfetch} />
      </div>
    </section>
  {/if}
</main>
