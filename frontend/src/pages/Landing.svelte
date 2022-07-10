<script>
    import Countries from '../components/user/Countries.svelte'
    import CheckIcon from '../components/icons/CheckIcon.svelte'

    import { _ } from '../i18n.js'
    import { AppConfig, appRoutes } from '../config.js'
    import { warrior } from '../stores.js'

    export let xfetch
    export let eventTag

    const { ShowActiveCountries, PathPrefix } = AppConfig
</script>

<style>
    :global(.text-yellow-thunder) {
        color: #ffdd57;
    }

    :global(.bg-yellow-thunder) {
        background-color: #ffdd57;
    }

    .browser-mockup {
        border-top: 2em solid;
        position: relative;
    }

    .browser-mockup:before {
        display: block;
        position: absolute;
        content: '';
        top: -1.25em;
        left: 1em;
        width: 0.5em;
        height: 0.5em;
        border-radius: 50%;
        background-color: #f44;
        box-shadow: 0 0 0 2px #f44, 1.5em 0 0 2px #9b3, 3em 0 0 2px #fb5;
    }

    :global(.title-line) {
        width: 10rem;
        height: 0.25rem;
        margin-bottom: 0.75rem;
    }
</style>

<svelte:head>
    <title>{$_('appName')} - {$_('appSubtitle')}</title>
</svelte:head>

<section
    class="w-full px-4 bg-yellow-thunder text-gray-800 border-b dark:border-gray-700"
>
    <div class="container mx-auto py-12 md:py-16 lg:py-20">
        <div class="flex flex-wrap items-center -mx-4">
            <div class="w-full md:w-1/2 mb-4 lg:mb-0 px-4">
                <h1
                    class="mb-2 lg:mb-4 text-5xl font-rajdhani uppercase font-semibold leading-none"
                >
                    {$_('landingTitle')}
                </h1>

                <p class="py-2 lg:py-4 text-2xl lg:text-3xl mb-4 font-rajdhani">
                    {$_('landingSalesPitch')}
                </p>

                <div class="text-center mb-4 md:mb-0">
                    <a
                        class="w-full text-3xl md:w-auto inline-block
                        no-underline bg-transparent bg-gray-800
                        hover:bg-transparent hover:text-gray-800 font-semibold
                        text-yellow-thunder py-4 px-10 border
                        hover:border-gray-800 border-transparent rounded font-rajdhani uppercase"
                        href="{$warrior.id
                            ? appRoutes.battles
                            : appRoutes.register}"
                    >
                        {$_('battleCreate')}
                    </a>
                </div>
            </div>
            <div class="w-full md:w-1/2 px-4">
                <div
                    class="browser-mockup flex flex-1 bg-white rounded shadow-xl"
                >
                    <img
                        class="w-full"
                        src="https://user-images.githubusercontent.com/846933/144792533-351cbc72-877e-4832-a4ce-d3b184e0c97e.png"
                        alt="{$_('appPreviewAlt')}"
                    />
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
                        {$_('customizableBattleOptions')}
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
                        Although every battle has a secure unique ID optionally set
                        a passcode to join
                    </li>
                    <li class="mb-2">
                        <CheckIcon />
                        Leader code &amp; multiple leader features allow team members
                        to continue in your absence
                    </li>
                </ul>
            </div>
            <div class="mt-4 md:mt-0 md:w-1/2">
                <img
                    src="https://user-images.githubusercontent.com/846933/144793090-a0f33f90-b679-464c-89b2-9bfb27a18e10.png"
                    alt="customizable battle options preview"
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
                        {$_('conciseVotingResults')}
                    </h3>
                </div>

                <ul class="text-lg dark:text-gray-300">
                    <li class="mb-2">
                        <CheckIcon />
                        Total votes, average points, and highest point metrics
                    </li>
                    <li class="mb-2">
                        <CheckIcon />
                        Participant vote transparency helps drive team discussion
                        and aid in decision making
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
                    Realtime agile sprint retrospectives with grouping, voting,
                    and action items.
                </p>

                <div class="text-center mb-4 md:mb-0">
                    <a
                        class="w-full text-3xl md:w-auto inline-block
                        no-underline bg-transparent bg-gray-800
                        hover:bg-transparent hover:text-gray-800 font-semibold
                        text-teal-400 py-4 px-10 border
                        hover:border-gray-800 border-transparent rounded font-rajdhani uppercase"
                        href="{$warrior.id
                            ? appRoutes.retros
                            : appRoutes.register}"
                    >
                        Create Retro
                    </a>
                </div>
            </div>
            <div class="w-full md:w-1/2 px-4">
                <div
                    class="browser-mockup flex flex-1 bg-white rounded shadow-xl"
                >
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
                    Instead of spending time discussing what you did yesterday
                    and what you're going to do today, focus on Blockers and
                    other more critical details.
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
                    Realtime agile feature story mapping with goals, columns,
                    stories and more!
                </p>

                <div class="text-center mb-4 md:mb-0">
                    <a
                        class="w-full text-3xl md:w-auto inline-block
                        no-underline bg-transparent bg-gray-800
                        hover:bg-transparent hover:text-gray-800 font-semibold
                        text-violet-400 py-4 px-10 border
                        hover:border-gray-800 border-transparent rounded font-rajdhani uppercase"
                        href="{$warrior.id
                            ? appRoutes.storyboards
                            : appRoutes.register}"
                    >
                        Create Storyboard
                    </a>
                </div>
            </div>
            <div class="w-full md:w-1/2 px-4">
                <div
                    class="browser-mockup flex flex-1 bg-white rounded shadow-xl"
                >
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
                    {$_('openSource')}
                </h3>
                <p class="px-2 text-lg dark:text-gray-300">
                    {@html $_('landingFeatureOpenSourceText', {
                        values: {
                            repoOpen: `<a
                        href="https://github.com/StevenWeathers/thunderdome-planning-poker"
                        class="no-underline text-blue-600 dark:text-sky-400 hover:text-blue-900 dark:hover:text-sky-600"
                    >`,
                            repoClose: '</a>',
                            donateOpen: `<a
                        href="https://github.com/StevenWeathers/thunderdome-planning-poker#donations"
                        class="no-underline text-blue-600 dark:text-sky-400 hover:text-blue-900 dark:hover:text-sky-600"
                    >`,
                            donateClose: '</a>',
                        },
                    })}
                </p>
            </div>
            <div class="w-1/2">
                <div class="mx-auto title-line bg-yellow-thunder"></div>
                <h3
                    class="text-4xl font-semibold font-rajdhani uppercase dark:text-white"
                >
                    {$_('hostedOrSelfHosted')}
                </h3>
                <p class="px-2 text-lg dark:text-gray-300">
                    {@html $_('hostedDesc', {
                        values: {
                            linkOpen: `<a href="https://thunderdome.dev" class="no-underline text-blue-600 dark:text-sky-400 hover:text-blue-900 dark:hover:text-sky-600">`,
                            linkClose: '</a>',
                        },
                    })}
                    {@html $_('selfHostedDesc', {
                        values: {
                            linkOpen: `<a
                            href="https://github.com/StevenWeathers/thunderdome-planning-poker#running-in-production"
                            class="no-underline text-blue-600 dark:text-sky-400 hover:text-blue-900 dark:hover:text-sky-600"
                    >`,
                            linkClose: '</a>',
                        },
                    })}
                </p>
            </div>
        </div>
    </div>
</section>
