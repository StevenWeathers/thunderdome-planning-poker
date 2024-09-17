<script>
  import { appRoutes } from '../../config';
  import { onMount } from 'svelte';

  onMount(() => window.scrollTo(0, 0));
</script>

<section class="bg-gradient-to-r from-blue-600 to-indigo-700 text-white">
  <div class="container mx-auto px-4 py-12">
    <div class="max-w-4xl mx-auto text-center">
      <h1 class="text-4xl sm:text-5xl lg:text-6xl font-bold mb-6 leading-tight">
        Sprint Retrospectives with <span
          class="bg-clip-text text-transparent bg-gradient-to-r from-yellow-thunder to-orange-500"
        >
          Thunderdome
        </span>
      </h1>
      <p class="text-xl sm:text-2xl text-blue-100">
        Elevate your team's performance with powerful, insightful retrospective
        tools.
      </p>
    </div>
  </div>
</section>

<div class="py-16 mx-auto max-w-2xl px-6 lg:max-w-7xl lg:px-8">
  <section class="mb-20 text-gray-700 dark:text-gray-300">
    <div class="mb-4">
      <div class="title-line bg-yellow-thunder"></div>
      <h2
        class="text-4xl font-semibold font-rajdhani uppercase dark:text-white"
      >
        Empower Your Team's Continuous Improvement
      </h2>
    </div>
    <p class="mb-4">
      Sprint retrospectives are crucial for agile teams to reflect, learn, and
      adapt. Thunderdome's Retrospective tool transforms this essential practice
      into an engaging and actionable experience, helping teams to:
    </p>
    <ul class="list-disc pl-8 mb-6 space-y-2">
      <li>Identify and address challenges quickly</li>
      <li>Celebrate successes and boost team morale</li>
      <li>Foster a culture of open communication</li>
      <li>Track progress on improvement initiatives</li>
      <li>Drive continuous enhancement of processes and practices</li>
    </ul>
  </section>

  <div class="space-y-24">
    <section class="flex flex-col md:flex-row items-center">
      <div class="md:w-1/2 md:pr-8 mb-8 md:mb-0">
        <div class="mb-4">
          <div class="title-line bg-yellow-thunder"></div>
          <h3
            class="text-4xl font-semibold font-rajdhani uppercase dark:text-white"
          >
            Customizable Retro Formats
          </h3>
        </div>
        <ul class="list-disc pl-5 space-y-2 text-gray-700 dark:text-gray-300">
          <li>
            Choose from popular templates like Start/Stop/Continue, 4Ls, and
            Mad/Sad/Glad
            <!--          </li>-->
            <!--          <li>-->
            <!--            Create and save custom retro templates tailored to your team's needs-->
            <!--          </li>-->
          </li>
          <li>
            Easily switch between formats to keep retrospectives fresh and
            engaging
          </li>
          <li>
            Flexible column structure adapts to different retrospective
            approaches
          </li>
        </ul>
      </div>
      <div class="md:w-1/2">
        <img
          src="{appConfig.PathPrefix}/img/previews/retro_template.png"
          alt="Customizable Retro Formats Screenshot"
          class="rounded-lg shadow-lg max-h-64 max-w-64 mx-auto hidden dark:block"
        />
        <img
          src="{appConfig.PathPrefix}/img/previews/retro_template_light.png"
          alt="Customizable Retro Formats Screenshot"
          class="rounded-lg shadow-lg max-h-64 max-w-64 mx-auto dark:hidden"
        />
      </div>
    </section>

    <section class="flex flex-col md:flex-row-reverse items-center">
      <div class="md:w-1/2 md:pl-8 mb-8 md:mb-0">
        <div class="mb-4">
          <div class="title-line bg-yellow-thunder"></div>
          <h3
            class="text-4xl font-semibold font-rajdhani uppercase dark:text-white"
          >
            Real-time Collaboration
          </h3>
        </div>
        <ul class="list-disc pl-5 space-y-2 text-gray-700 dark:text-gray-300">
          <li>
            Enable seamless participation for both co-located and distributed
            teams
          </li>
          <li>
            Facilitate discussions with built-in commenting and voting features
          </li>
          <li>
            Assign facilitator roles to guide the retro process effectively
          </li>
          <li>Use anonymity options to encourage open and honest feedback</li>
        </ul>
      </div>
      <div class="md:w-1/2">
        <img
          src="{appConfig.PathPrefix}/img/previews/retro.png"
          alt="Real-time Collaboration Screenshot"
          class="rounded-lg shadow-lg max-h-64 max-w-64 mx-auto hidden dark:block"
        />
        <img
          src="{appConfig.PathPrefix}/img/previews/retro_light.png"
          alt="Real-time Collaboration Screenshot"
          class="rounded-lg shadow-lg max-h-64 max-w-64 mx-auto dark:hidden"
        />
      </div>
    </section>

    <section class="flex flex-col md:flex-row items-center">
      <div class="md:w-1/2 md:pr-8 mb-8 md:mb-0">
        <div class="mb-4">
          <div class="title-line bg-yellow-thunder"></div>
          <h3
            class="text-4xl font-semibold font-rajdhani uppercase dark:text-white"
          >
            Action Item Tracking
          </h3>
        </div>
        <ul class="list-disc pl-5 space-y-2 text-gray-700 dark:text-gray-300">
          <li>
            Create and assign action items directly from retro discussions
          </li>
          <li>Track progress on action items across multiple sprints</li>
        </ul>
      </div>
      <div class="md:w-1/2">
        <img
          src="{appConfig.PathPrefix}/img/previews/action_items.png"
          alt="Action Item Tracking Screenshot"
          class="rounded-lg shadow-lg max-h-64 max-w-64 mx-auto hidden dark:block"
        />
        <img
          src="{appConfig.PathPrefix}/img/previews/action_items_light.png"
          alt="Action Item Tracking Screenshot"
          class="rounded-lg shadow-lg max-h-64 max-w-64 mx-auto dark:hidden"
        />
      </div>
    </section>
  </div>
</div>

<div class="bg-indigo-600 py-12">
  <div class="mx-auto max-w-7xl px-6 lg:px-8 text-center">
    <h2 class="text-3xl font-bold tracking-tight text-white sm:text-4xl">
      Ready to supercharge your retrospectives?
    </h2>
    <p class="mt-6 text-lg leading-8 text-indigo-100">
      Join thousands of agile teams using Thunderdome's Retrospective tool to
      drive continuous improvement and team excellence.
    </p>
    <div class="mt-10 flex items-center justify-center gap-x-6">
      <div class="text-center">
        <div
          class="flex flex-col sm:flex-row items-center justify-center space-y-4 sm:space-y-0 sm:space-x-4"
        >
          <a
            href="{appRoutes.login}"
            class="bg-white bg-indigo-600 hover:bg-gray-100 font-semibold py-2 px-6 rounded-full transition duration-300"
          >
            Login
          </a>
          <div class="text-white font-medium">or</div>
          <a
            href="{appRoutes.register}"
            class="bg-transparent text-white hover:bg-white/10 border border-white font-semibold py-2 px-6 rounded-full transition duration-300"
          >
            Register
          </a>
        </div>
        <p class="text-white mt-2">to get started</p>
      </div>
    </div>
  </div>
</div>
