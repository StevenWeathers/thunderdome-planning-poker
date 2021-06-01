<script>
    import { createEventDispatcher } from 'svelte'

    import DownCarrotIcon from './icons/DownCarrotIcon.svelte'
    import { locales } from '../i18n'

    export let selectedLocale = 'en'
    const supportedLocales = []

    for (const [key, value] of Object.entries(locales)) {
        supportedLocales.push({
            name: value,
            value: key,
        })
    }

    const dispatch = createEventDispatcher()

    function switchLocale(event) {
        event.preventDefault()

        dispatch('locale-changed', event.target.value)
    }
</script>

<div class="ml-2 inline-block">
    <div class="relative">
        <select
            name="locale"
            on:change="{switchLocale}"
            value="{selectedLocale}"
            class="block appearance-none w-full border-2 border-gray-400
            text-gray-700 py-2 px-4 pr-8 rounded leading-tight
            focus:outline-none focus:border-purple-500">
            {#each supportedLocales as locale}
                <option value="{locale.value}">{locale.name}</option>
            {/each}
        </select>
        <div
            class="pointer-events-none absolute inset-y-0 right-0 flex
            items-center px-2 text-gray-700">
            <DownCarrotIcon />
        </div>
    </div>
</div>
