<script>
    import Countries from '../components/Countries.svelte'
    import CheckIcon from '../components/icons/CheckIcon.svelte'

    import { _ } from '../i18n'
    import { appRoutes } from '../config'
    import { warrior } from '../stores'

    export let xfetch
    export let eventTag

    const { ShowActiveCountries } = appConfig
</script>

<style>
    .bg-exo-grey {
        background-color: #2c2d2c;
    }
    .feature-link {
        @apply no-underline;
        @apply text-blue-600;
    }
    .feature-link:hover {
        @apply text-blue-900;
    }
</style>

<section class="bg-yellow-thunder text-gray-800">
    <div class="container mx-auto px-4 py-6 lg:py-10">
        <div class="flex flex-wrap items-center -mx-4">
            <div class="w-full md:w-1/2 mb-4 lg:mb-0 px-4">
                <h1 class="mb-2 lg:mb-4 text-4xl font-bold leading-tight">
                    {$_('pages.landing.title')}
                </h1>

                <ul class="text-lg lg:text-2xl mb-4">
                    <li class="py-2 lg:py-4">
                        <CheckIcon />
                        {$_('pages.landing.bullet1')}
                    </li>
                    <li class="py-2 lg:py-4">
                        <CheckIcon />
                        {@html $_('pages.landing.bullet2', {
                            values: {
                                crossoutOpen: `<span class="line-through">`,
                                crossoutClose: '</span>',
                            },
                        })}
                    </li>
                    <li class="py-2 lg:py-4">
                        <CheckIcon />
                        {@html $_('pages.landing.bullet3', {
                            values: {
                                mbOpen: `<span class="font-semibold">`,
                                mbClose: '</span>',
                            },
                        })}
                        <br />
                        <span class="text-sm">
                            {$_('pages.landing.bullet3NotSorry')}
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
                        href="{$warrior.id ? appRoutes.battles : appRoutes.register}">
                        {$_('actions.battle.create')}
                    </a>
                </div>
            </div>
            <div class="w-full md:w-1/2 px-4">
                <img
                    class="w-full"
                    src="https://user-images.githubusercontent.com/846933/118144911-fdf70980-b3da-11eb-92b6-63de3990a8b8.png"
                    alt="Preview of Thunderdome, Agile Planning Poker UI
                    Experience" />
            </div>
        </div>
    </div>
</section>

<section class="bg-white">
    <div class="container mx-auto px-4 py-6 lg:py-10">
        <div class="text-center mb-6">
            <h2 class="text-4xl font-bold">
                {$_('pages.landing.features.title')}
            </h2>
        </div>
        <div class="flex text-center mb-8">
            <div class="w-1/2">
                <h3 class="text-2xl font-bold">Open Source</h3>
                <p class="px-2">
                    Check out the
                    <a
                        href="https://github.com/StevenWeathers/thunderdome-planning-poker"
                        class="feature-link">
                        repository
                    </a>
                    at GitHub to request or contribute enhancements, locale
                    translations, and bug fixes or to
                    <a
                        href="https://github.com/StevenWeathers/thunderdome-planning-poker#donations"
                        class="feature-link">
                        Donate
                    </a>
                </p>
            </div>
            <div class="w-1/2">
                <h3 class="text-2xl font-bold">
                    Hosted service or self-hosted
                </h3>
                <p class="px-2">
                    <a href="https://thunderdome.dev" class="feature-link">
                        Thunderdome.dev
                    </a>
                    is a free hosted service (SaaS) without Ads, alternatively
                    you can
                    <a
                        href="https://github.com/StevenWeathers/thunderdome-planning-poker#running-in-production"
                        class="feature-link">
                        run it on your own server
                    </a>
                </p>
            </div>
        </div>
        <div class="flex text-center mb-8">
            <div class="w-1/2">
                <h3 class="text-2xl font-bold">Organizations and Teams</h3>
                <p class="px-2">
                    Create organizations with departments and teams at each
                    level to associate battle's to a team for quick access. Just
                    a small team?, no problem! Stand alone teams can be created
                    without the need for organizations.
                </p>
            </div>
            <div class="w-1/2">
                <h3 class="text-2xl font-bold">Customizable Battle Options</h3>
                <p class="px-2">
                    Each battle can set the allowed point values from the
                    configured (default 0-100 scale) set as well as whether or
                    not to automatically finish voting when all votes have been
                    cast.
                </p>
            </div>
        </div>
        <div class="flex text-center mb-8">
            <div class="w-1/2">
                <h3 class="text-2xl font-bold">Concise Voting Results</h3>
                <p class="px-2">
                    Voting results show how many and whom voted for each point
                    value. A tally of total votes, highest vote, and a rounded
                    (customizable per battle) average vote helps quickly come to
                    a consensus.
                </p>
            </div>
            <div class="w-1/2">
                <h3 class="text-2xl font-bold">Jira XML Import</h3>
                <p class="px-2">
                    To expedite setting up a grooming session importing
                    <a
                        href="https://www.atlassian.com/software/jira"
                        class="feature-link">
                        Jira
                    </a>
                    XML to quickly add stories to a battle with links back to
                    Jira. Don't use Jira?, self-hosted instances can enable API
                    access for easier integrations with other systems.
                </p>
            </div>
        </div>
        <div class="flex text-center mb-8">
            <div class="w-1/2">
                <h3 class="text-2xl font-bold">Fast</h3>
                <p class="px-2">
                    Thunderdome is fast and small (under 20mb single executable)
                    thanks to
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
                <h3 class="text-2xl font-bold">Secure</h3>
                <p class="px-2">
                    Every battle has a randomly generated, not incremetal ID
                    which keeps strangers from joining your battle without the
                    share link.
                    <br />
                    HTTPS, backend only cookie, and salted/hashed password
                    storage keep your authenticated sessions secure.
                </p>
            </div>
        </div>
    </div>
</section>

{#if ShowActiveCountries}
    <Countries {xfetch} {eventTag} />
{/if}

<section class="bg-exo-grey text-white mb-8">
    <div class="container mx-auto px-4 py-2 lg:py-4">
        <div class="flex text-center">
            <div class="w-full">
                <h4 class="font-bold text-xl">
                    Check out
                    <a href="https://exothermic.dev">
                        <img
                            src="/img/exothermic-logo.png"
                            alt="Exothermic - Open Source Agile Story Mapping
                            tool"
                            title="Exothermic - Open Source Agile Story Mapping
                            tool"
                            class="inline-block mx-2" />
                    </a>
                    an open source agile story mapping tool brought to you by
                    yours truly!
                </h4>
            </div>
        </div>
    </div>
</section>
