<script lang="ts">
  import { AppConfig, appRoutes } from '../config';
  import { user } from '../stores';
  import { validateUserIsRegistered } from '../validationUtils';
  import { Github, Users, Zap } from 'lucide-svelte';
  import LL from '../i18n/i18n-svelte';
  import Countries from '../components/user/Countries.svelte';
  import BrowserMock from '../components/global/BrowserMock.svelte';
  import { onMount } from 'svelte';

  export let xfetch;
  export let eventTag;

  const {
    ShowActiveCountries,
    PathPrefix,
    SubscriptionsEnabled,
    SubscriptionCheckoutLink,
    RepoURL,
  } = AppConfig;

  $: isRegisteredUser = $user && !!$user.id && validateUserIsRegistered($user);

  onMount(() => window.scrollTo(0, 0));
</script>

<svelte:head>
  <title>{$LL.appName()} - {$LL.appSubtitle()}</title>
</svelte:head>

<main class="bg-gray-100 dark:bg-gray-900">
  <header class="bg-gradient-to-r from-blue-600 to-indigo-700 text-white">
    <div class="container mx-auto px-4 py-16">
      <div class="max-w-3xl mx-auto text-center">
        <h1
          class="text-4xl sm:text-5xl lg:text-6xl font-bold mb-6 leading-tight"
        >
          Thunderdome: Empower Your Agile Teams
        </h1>
        <p class="text-xl sm:text-2xl text-blue-100">
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
            class="inline-block bg-blue-500 text-white hover:bg-blue-600 dark:bg-gray-700 dark:hover:bg-gray-600 dark:text-teal-300 font-semibold py-3 px-8 rounded-lg transition duration-300 ease-in-out transform hover:-translate-y-1 hover:scale-105"
          >
            {$user.id ? $LL.battleCreate() : 'Learn More'}
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
            class="inline-block bg-blue-500 text-white hover:bg-blue-600 dark:bg-gray-800 dark:hover:bg-gray-700 dark:text-teal-300 font-semibold py-3 px-8 rounded-lg transition duration-300 ease-in-out transform hover:-translate-y-1 hover:scale-105"
          >
            {$user.id ? 'Create Retro' : 'Learn More'}
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
            class="inline-block bg-blue-500 text-white hover:bg-blue-600 dark:bg-gray-700 dark:hover:bg-gray-600 dark:text-teal-300 font-semibold py-3 px-8 rounded-lg transition duration-300 ease-in-out transform hover:-translate-y-1 hover:scale-105"
          >
            {$user.id ? 'Create Storyboard' : 'Learn More'}
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
              href="{RepoURL}#running-in-production"
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
        <Countries xfetch="{xfetch}" eventTag="{eventTag}" />
      </div>
    </section>
  {/if}
</main>
