<script lang="ts">
  import Countries from '../components/user/Countries.svelte';
  import CheckIcon from '../components/icons/CheckIcon.svelte';
  import LL from '../i18n/i18n-svelte';

  import { AppConfig, appRoutes } from '../config';
  import { user } from '../stores';
  import { validateUserIsRegistered } from '../validationUtils';

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
</script>

<svelte:head>
  <title>{$LL.appName()} - {$LL.appSubtitle()}</title>
</svelte:head>

{#if SubscriptionsEnabled && isRegisteredUser && !$user.subscribed}
  <section class="w-full px-4 bg-cyan-200 border-b border-cyan-400">
    <div class="container mx-auto">
      <div class="flex flex-wrap">
        <div class="w-full py-4 px-4">
          <h2
            class="text-2xl text-gray-700 font-rajdhani uppercase font-semibold leading-none text-center"
          >
            Enjoying Thunderdome?
            <a
              class="text-blue-800 underline"
              href="{appRoutes.subscriptionPricing}"
            >
              Subscribe today
            </a>
            starting at only $5 /mo.
          </h2>
        </div>
      </div>
    </div>
  </section>
{/if}

<section
  class="w-full px-4 bg-yellow-thunder text-gray-800 border-b dark:border-gray-700"
>
  <div class="container mx-auto py-12 md:py-16 lg:py-20">
    <div class="flex flex-wrap items-center -mx-4">
      <div class="w-full md:w-1/2 mb-4 lg:mb-0 px-4">
        <h1
          class="mb-2 lg:mb-4 text-5xl font-rajdhani uppercase font-semibold leading-none"
        >
          {$LL.landingTitle()}
        </h1>

        <p class="py-2 lg:py-4 text-2xl lg:text-3xl mb-4 font-rajdhani">
          {$LL.landingSalesPitch()}
        </p>

        <div class="text-center mb-4 md:mb-0">
          <a
            class="w-full text-3xl md:w-auto inline-block
                        no-underline bg-gray-800
                        hover:bg-transparent hover:text-gray-800 font-semibold
                        text-yellow-thunder py-4 px-10 border
                        hover:border-gray-800 border-transparent rounded font-rajdhani uppercase"
            href="{$user.id ? appRoutes.games : appRoutes.register}"
          >
            {$LL.battleCreate()}
          </a>
        </div>
      </div>
      <div class="w-full md:w-1/2 px-4">
        <div class="rounded shadow-xl">
          <div class="rounded-t bg-gray-800 px-4 py-2.5">
            <div class="flex flex-row space-x-2">
              <div class="h-3 w-3 rounded-full bg-rose-500"></div>
              <div class="h-3 w-3 rounded-full bg-amber-300"></div>
              <div class="h-3 w-3 rounded-full bg-lime-400"></div>
            </div>
          </div>
          <div>
            <img
              class="w-full rounded-b"
              src="{PathPrefix}/img/web_poker_preview.png"
              alt="{$LL.appPreviewAlt()}"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</section>

<section
  class="w-full px-4 bg-white dark:bg-gray-800 border-b dark:border-gray-700"
>
  <div class="container mx-auto py-12 md:py-16 lg:py-20">
    <div class="flex items-center flex-wrap h-full">
      <div class="md:w-1/2">
        <div class="mb-4">
          <div class="title-line bg-yellow-thunder"></div>
          <h3
            class="text-4xl font-semibold font-rajdhani uppercase dark:text-white"
          >
            {$LL.customizableBattleOptions()}
          </h3>
        </div>

        <ul class="text-lg dark:text-gray-300">
          <li class="mb-2">
            <CheckIcon />
            Select allowed point values that fit your team's process
          </li>
          <li class="mb-2">
            <CheckIcon />
            Automatically end voting once all participants have voted
          </li>
          <li class="mb-2">
            <CheckIcon />
            Although every game has a secure unique ID optionally set a passcode
            to join
          </li>
          <li class="mb-2">
            <CheckIcon />
            Leader code &amp; multiple facilitator features allow team members to
            continue in your absence
          </li>
        </ul>
      </div>
      <div class="mt-4 md:mt-0 md:w-1/2">
        <img
          src="{PathPrefix}/img/poker_options.png"
          alt="customizable game options preview"
          class="w-3/4 m-auto shadow-xl"
        />
      </div>
    </div>
  </div>
</section>

<section
  class="w-full px-4 bg-slate-100 dark:bg-gray-900 border-b dark:border-gray-700"
>
  <div class="container mx-auto py-12 md:py-16 lg:py-20">
    <div class="flex items-center flex-wrap h-full">
      <div class="md:w-1/2">
        <img
          src="https://user-images.githubusercontent.com/846933/144792861-d17d532f-2235-4a4a-b38f-90be065a2447.png"
          class="w-3/4 m-auto"
          alt="concise voting results preview"
        />
      </div>
      <div class="mt-4 md:mt-0 md:w-1/2">
        <div class="mb-4">
          <div class="title-line bg-yellow-thunder"></div>
          <h3
            class="text-4xl font-semibold font-rajdhani uppercase dark:text-white"
          >
            {$LL.conciseVotingResults()}
          </h3>
        </div>

        <ul class="text-lg dark:text-gray-300">
          <li class="mb-2">
            <CheckIcon />
            Total votes, average points, and highest point metrics
          </li>
          <li class="mb-2">
            <CheckIcon />
            Participant vote transparency helps drive team discussion and aid in
            decision making
          </li>
        </ul>
      </div>
    </div>
  </div>
</section>

<section
  class="w-full px-4 bg-teal-400 text-gray-800 border-b dark:border-gray-700"
>
  <div class="container mx-auto py-12 md:py-16 lg:py-20">
    <div class="flex flex-wrap items-center -mx-4">
      <div class="w-full md:w-1/2 mb-4 lg:mb-0 px-4">
        <h2
          class="mb-2 lg:mb-4 text-4xl font-rajdhani uppercase font-semibold leading-none"
        >
          Agile Sprint Retrospectives
        </h2>

        <p class="py-2 lg:py-4 text-2xl lg:text-3xl mb-4 font-rajdhani">
          Realtime agile sprint retrospectives with grouping, voting, and action
          items.
        </p>

        <div class="text-center mb-4 md:mb-0">
          <a
            class="w-full text-3xl md:w-auto inline-block
                        no-underline bg-gray-800
                        hover:bg-transparent hover:text-gray-800 font-semibold
                        text-teal-400 py-4 px-10 border
                        hover:border-gray-800 border-transparent rounded font-rajdhani uppercase"
            href="{$user.id ? appRoutes.retros : appRoutes.register}"
          >
            Create Retro
          </a>
        </div>
      </div>
      <div class="w-full md:w-1/2 px-4">
        <div class="browser-mockup flex flex-1 bg-white rounded shadow-xl">
          <img
            class="w-full"
            src="https://user-images.githubusercontent.com/846933/173260209-3ef3299f-f1b2-41e8-802f-17d40649c66d.png"
            alt="Preview of Thunderdome.dev retrospective feature"
          />
        </div>
      </div>
    </div>
  </div>
</section>

<section
  class="w-full px-4 bg-white dark:bg-gray-900 border-b dark:border-gray-700"
>
  <div class="container mx-auto py-12 md:py-16 lg:py-20">
    <div class="flex items-center flex-wrap h-full">
      <div class="md:w-1/2">
        <img
          src="https://user-images.githubusercontent.com/846933/178159914-981b7962-f453-4b98-b3d0-df274859830a.png"
          class="w-3/4 m-auto"
          alt="team checkins preview"
        />
      </div>
      <div class="mt-4 md:mt-0 md:w-1/2">
        <div class="mb-4">
          <div class="title-line bg-yellow-thunder"></div>
          <h3
            class="text-4xl font-semibold font-rajdhani uppercase dark:text-white"
          >
            Streamline your team's agile stand-up with Team Checkins
          </h3>
        </div>

        <div class="text-lg dark:text-gray-300">
          Instead of spending time discussing what you did yesterday and what
          you're going to do today, focus on Blockers and other more critical
          details.
        </div>
      </div>
    </div>
  </div>
</section>

<section
  class="w-full px-4 bg-violet-400 text-gray-800 border-b dark:border-gray-700"
>
  <div class="container mx-auto py-12 md:py-16 lg:py-20">
    <div class="flex flex-wrap items-center -mx-4">
      <div class="w-full md:w-1/2 mb-4 lg:mb-0 px-4">
        <h2
          class="mb-2 lg:mb-4 text-4xl font-rajdhani uppercase font-semibold leading-none"
        >
          Agile Feature Story Mapping
        </h2>

        <p class="py-2 lg:py-4 text-2xl lg:text-3xl mb-4 font-rajdhani">
          Realtime agile feature story mapping with goals, columns, stories and
          more!
        </p>

        <div class="text-center mb-4 md:mb-0">
          <a
            class="w-full text-3xl md:w-auto inline-block
                        no-underline bg-gray-800
                        hover:bg-transparent hover:text-gray-800 font-semibold
                        text-violet-400 py-4 px-10 border
                        hover:border-gray-800 border-transparent rounded font-rajdhani uppercase"
            href="{$user.id ? appRoutes.storyboards : appRoutes.register}"
          >
            Create Storyboard
          </a>
        </div>
      </div>
      <div class="w-full md:w-1/2 px-4">
        <div class="browser-mockup flex flex-1 bg-white rounded shadow-xl">
          <img
            class="w-full"
            src="https://user-images.githubusercontent.com/846933/173260211-304a973d-4ede-494f-bb7d-b7e5c86a4e6e.png"
            alt="Preview of Thunderdome.dev story mapping feature"
          />
        </div>
      </div>
    </div>
  </div>
</section>

{#if ShowActiveCountries}
  <section
    class="w-full px-4 bg-slate-100 dark:bg-gray-900 border-b dark:border-gray-700"
  >
    <div class="container mx-auto py-12 md:py-16 lg:py-20">
      <Countries xfetch="{xfetch}" eventTag="{eventTag}" />
    </div>
  </section>
{/if}

<section class="w-full px-4 bg-white dark:bg-gray-800">
  <div class="container mx-auto py-12 md:py-16 lg:py-20">
    <div class="flex text-center mb-8">
      <div class="w-1/2">
        <div class="mx-auto title-line bg-yellow-thunder"></div>
        <h3
          class="text-4xl font-semibold font-rajdhani uppercase dark:text-white"
        >
          {$LL.openSource()}
        </h3>
        <p class="px-2 text-lg dark:text-gray-300">
          {@html $LL.landingFeatureOpenSourceText({
            repoOpen: `<a
                        href="${RepoURL}"
                        class="no-underline text-blue-600 dark:text-sky-400 hover:text-blue-900 dark:hover:text-sky-600"
                    >`,
            repoClose: '</a>',
            donateOpen: `<a
                        href="${RepoURL}#donations"
                        class="no-underline text-blue-600 dark:text-sky-400 hover:text-blue-900 dark:hover:text-sky-600"
                    >`,
            donateClose: '</a>',
          })}
        </p>
      </div>
      <div class="w-1/2">
        <div class="mx-auto title-line bg-yellow-thunder"></div>
        <h3
          class="text-4xl font-semibold font-rajdhani uppercase dark:text-white"
        >
          {$LL.hostedOrSelfHosted()}
        </h3>
        <p class="px-2 text-lg dark:text-gray-300">
          {@html $LL.hostedDesc({
            linkOpen: `<a href="https://thunderdome.dev" class="no-underline text-blue-600 dark:text-sky-400 hover:text-blue-900 dark:hover:text-sky-600">`,
            linkClose: '</a>',
          })}
          {@html $LL.selfHostedDesc({
            linkOpen: `<a
                            href="${RepoURL}#running-in-production"
                            class="no-underline text-blue-600 dark:text-sky-400 hover:text-blue-900 dark:hover:text-sky-600"
                    >`,
            linkClose: '</a>',
          })}
        </p>
      </div>
    </div>
  </div>
</section>
