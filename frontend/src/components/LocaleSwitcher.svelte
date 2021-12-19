<script>
    import { createEventDispatcher } from 'svelte'

    import DownCarrotIcon from './icons/ChevronDown.svelte'
    import { locales } from '../i18n.js'

    export let selectedLocale = 'en'
    let klass = ''
    export { klass as class }
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

<div class="{klass} inline-block">
    <div class="relative">
        <select
            name="locale"
            on:change="{switchLocale}"
            value="{selectedLocale}"
            class="block appearance-none w-full border border-gray-300 dark:border-gray-700
                text-gray-700 dark:text-gray-300 py-2 px-4 pr-8 rounded leading-tight
                focus:outline-none focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400 dark:bg-gray-900"
        >
            {#each supportedLocales as locale}
                <option value="{locale.value}">{locale.name}</option>
            {/each}
        </select>
        <div
            class="pointer-events-none absolute inset-y-0 right-0 flex
            items-center px-2 text-gray-700 dark:text-gray-400"
        >
            <DownCarrotIcon />
        </div>
    </div>
</div>
