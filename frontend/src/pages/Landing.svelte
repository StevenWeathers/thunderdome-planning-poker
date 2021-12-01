<script>
    import Countries from '../components/user/Countries.svelte'
    import CheckIcon from '../components/icons/CheckIcon.svelte'
    import GithubIcon from '../components/icons/Github.svelte'

    import { _ } from '../i18n'
    import { AppConfig, appRoutes } from '../config'
    import { warrior } from '../stores'

    export let xfetch
    export let eventTag

    const { ShowActiveCountries, PathPrefix } = AppConfig
</script>

<style>
    :global(.feature-link) {
        @apply no-underline;
        @apply text-blue-600;
    }

    :global(.feature-link:hover) {
        @apply text-blue-900;
    }

    :global(.text-yellow-thunder) {
        color: #ffdd57;
    }

    :global(.bg-yellow-thunder) {
        background-color: #ffdd57;
    }
</style>

<svelte:head>
    <title>{$_('appName')} - {$_('appSubtitle')}</title>
</svelte:head>

<section class="bg-yellow-thunder text-gray-800">
    <div class="container mx-auto px-4 py-6 lg:py-10">
        <div class="flex flex-wrap items-center -mx-4">
            <div class="w-full md:w-1/2 mb-4 lg:mb-0 px-4">
                <h1 class="mb-2 lg:mb-4 text-4xl font-bold leading-tight">
                    {$_('landingTitle')}
                </h1>

                <ul class="text-lg lg:text-2xl mb-4">
                    <li class="py-2 lg:py-4">
                        <CheckIcon />
                        {$_('landingBullet1')}
                    </li>
                    <li class="py-2 lg:py-4">
                        <CheckIcon />
                        {@html $_('landingBullet2', {
                            values: {
                                crossoutOpen: `<span class="line-through">`,
                                crossoutClose: '</span>',
                            },
                        })}
                    </li>
                    <li class="py-2 lg:py-4">
                        <CheckIcon />
                        {@html $_('landingBullet3', {
                            values: {
                                mbOpen: `<span class="font-semibold">`,
                                mbClose: '</span>',
                            },
                        })}
                        <br />
                        <span class="text-sm">
                            {$_('landingBullet3NotSorry')}
                        </span>
                    </li>
                </ul>

                <div class="text-center mb-4 md:mb-0">
                    <a
                        class="w-full text-xl md:w-auto inline-block
                        no-underline bg-transparent bg-gray-800
                        hover:bg-transparent hover:text-gray-800 font-semibold
                        text-yellow-thunder py-4 px-10 border
                        hover:border-gray-800 border-transparent rounded"
                        href="{$warrior.id
                            ? appRoutes.battles
                            : appRoutes.register}"
                    >
                        {$_('battleCreate')}
                    </a>
                </div>
            </div>
            <div class="w-full md:w-1/2 px-4">
                <img
                    class="w-full"
                    src="https://user-images.githubusercontent.com/846933/118144911-fdf70980-b3da-11eb-92b6-63de3990a8b8.png"
                    alt="{$_('appPreviewAlt')}"
                />
            </div>
        </div>
    </div>
</section>

<section class="bg-white">
    <div class="container mx-auto px-4 py-6 lg:py-10">
        <div class="text-center mb-6">
            <h2 class="text-4xl font-bold">{$_('landingFeatures')}</h2>
        </div>
        <div class="flex text-center mb-8">
            <div class="w-1/2">
                <h3 class="text-2xl font-bold">{$_('openSource')}</h3>
                <p class="px-2">
                    Check out the
                    <GithubIcon />
                    <a
                        href="https://github.com/StevenWeathers/thunderdome-planning-poker"
                        class="feature-link"
                    >
                        repository
                    </a>
                    to request or contribute enhancements, locale translations, and
                    bug fixes or to
                    <a
                        href="https://github.com/StevenWeathers/thunderdome-planning-poker#donations"
                        class="feature-link"
                    >
                        Donate
                    </a>
                </p>
            </div>
            <div class="w-1/2">
                <h3 class="text-2xl font-bold">
                    {$_('hostedOrSelfHosted')}
                </h3>
                <p class="px-2">
                    <a href="https://thunderdome.dev" class="feature-link">
                        Thunderdome.dev
                    </a>
                    is a free hosted service (SaaS) without Ads, alternatively you
                    can
                    <a
                        href="https://github.com/StevenWeathers/thunderdome-planning-poker#running-in-production"
                        class="feature-link"
                    >
                        run it on your own server
                    </a>
                </p>
            </div>
        </div>
        <div class="flex text-center mb-8">
            <div class="w-1/2">
                <h3 class="text-2xl font-bold">
                    {$_('organizationsAndTeams')}
                </h3>
                <p class="px-2">
                    {$_('organizationsAndTeamsDescription')}
                </p>
            </div>
            <div class="w-1/2">
                <h3 class="text-2xl font-bold">
                    {$_('customizableBattleOptions')}
                </h3>
                <p class="px-2">
                    {$_('customizableBattleOptionsDescription')}
                </p>
            </div>
        </div>
        <div class="flex text-center mb-8">
            <div class="w-1/2">
                <h3 class="text-2xl font-bold">{$_('conciseVotingResults')}</h3>
                <p class="px-2">
                    {$_('conciseVotingResultsDescription')}
                </p>
            </div>
            <div class="w-1/2">
                <h3 class="text-2xl font-bold">{$_('jiraXmlImport')}</h3>
                <p class="px-2">
                    To expedite setting up a grooming session importing
                    <a
                        href="https://www.atlassian.com/software/jira"
                        class="feature-link"
                    >
                        Jira
                    </a>
                    XML to quickly add stories to a battle with links back to Jira.
                    Don't use Jira?, create an API Key and integrate your system
                    of choice.
                </p>
            </div>
        </div>
        <div class="flex text-center mb-8">
            <div class="w-1/2">
                <h3 class="text-2xl font-bold">{$_('fast')}</h3>
                <p class="px-2">
                    Thunderdome is fast and small (single executable) thanks to
                    <a href="https://golang.org/" class="feature-link">Go</a>
                    .
                    <br />
                    <a href="https://svelte.dev/" class="feature-link">
                        Svelte
                    </a>
                    and Websockets provide a reactive web app experience.
                    <br />
                    <span class="font-bold">
                        Over 30 active users in a battle tested!
                    </span>
                </p>
            </div>
            <div class="w-1/2">
                <h3 class="text-2xl font-bold">{$_('secure')}</h3>
                <p class="px-2">
                    {@html $_('secureDescription')}
                </p>
            </div>
        </div>
    </div>
</section>

{#if ShowActiveCountries}
    <Countries xfetch="{xfetch}" eventTag="{eventTag}" />
{/if}

<section class="mb-8">
    <div class="container mx-auto px-4 py-2 lg:py-4">
        <div class="flex text-center">
            <div class="w-full">
                <h2 class="font-bold text-4xl mb-4">{$_('additionalTools')}</h2>
                <div class="font-bold text-xl align-middle flex justify-center">
                    <div class="flex-initial px-4">
                        <a href="https://exothermic.dev">
                            <img
                                src="{PathPrefix}/img/exothermic-logo-2.png"
                                alt="Exothermic logo"
                                title="Exothermic - Open Source Agile Story
                                Mapping tool"
                                class="inline-block mx-2 mb-1"
                            />
                        </a>
                        <br />
                        {$_('exothermicDescription')}
                    </div>
                    <div class="flex-initial px-4">
                        <a href="https://wakita.dev">
                            <img
                                src="{PathPrefix}/img/wakita-logo.png"
                                alt="Wakita logo"
                                title="Wakita - Open Source Agile Retrospective
                                tool"
                                class="inline-block mx-2 mb-1"
                            />
                            <br />
                            {$_('wakitaDescription')}
                        </a>
                    </div>
                </div>
            </div>
        </div>
    </div>
</section>
