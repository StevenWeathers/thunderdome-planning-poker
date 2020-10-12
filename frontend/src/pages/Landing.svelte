<script>
    import CheckIcon from '../components/icons/CheckIcon.svelte'

    import { _ } from '../i18n'
    import { warrior } from '../stores'
</script>

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
                        href="{$warrior.id ? '/battles' : '/enlist'}">
                        {$_('actions.battle.create')}
                    </a>
                </div>
            </div>
            <div class="w-full md:w-1/2 px-4">
                <img
                    class="w-full"
                    src="https://user-images.githubusercontent.com/846933/95778842-eb76ef00-0c96-11eb-99d8-af5d098c12ee.png"
                    alt="Preview of Thunderdome, Agile Planning Poker UI
                    Experience" />
            </div>
        </div>
    </div>
</section>

<section class="bg-white">
    <div class="container mx-auto px-4 py-6 lg:py-10">
        <div class="text-center mb-6">
            <h2 class="text-2xl font-bold">
                {$_('pages.landing.features.title')}
            </h2>
        </div>
        <div class="flex text-center mb-8">
            <div class="w-1/3">
                <CheckIcon />
                {$_('pages.landing.features.bullet1')}
            </div>
            <div class="w-1/3">
                <CheckIcon />
                {$_('pages.landing.features.bullet2')}
            </div>
            <div class="w-1/3">
                <CheckIcon />
                {$_('pages.landing.features.bullet3')}
            </div>
        </div>
        <div class="flex text-center mb-8">
            <div class="w-1/3">
                <CheckIcon />
                {$_('pages.landing.features.bullet4')}
            </div>
            <div class="w-1/3">
                <CheckIcon />
                {$_('pages.landing.features.bullet5')}
            </div>
            <div class="w-1/3">
                <CheckIcon />
                {@html $_('pages.landing.features.bullet6', {
                    values: {
                        contribOpen: `<a
                    href="https://github.com/StevenWeathers/thunderdome-planning-poker"
                    class="no-underline text-teal-500 hover:text-teal-800\">`,
                        contribClose: '</a>',
                    },
                })}
            </div>
        </div>
        <div class="flex text-center">
            <div class="w-full">
                <CheckIcon />
                {$_('pages.landing.features.bullet7')}
            </div>
        </div>
    </div>
</section>
