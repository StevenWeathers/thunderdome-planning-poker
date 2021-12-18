<script>
    import CountryFlag from './CountryFlag.svelte'
    import { _ } from '../../i18n.js'

    export let xfetch
    export let eventTag

    let activeCountries = []

    xfetch('/api/active-countries')
        .then(res => res.json())
        .then(function (result) {
            activeCountries = result.data.sort()
        })
        .catch(function () {
            eventTag('get_active_countries', 'engagement', 'failure')
        })
</script>

<div class="text-center mb-6">
    <div class="mx-auto title-line bg-yellow-thunder"></div>
    <h2
        class="text-5xl font-semibold mb-12 font-rajdhani uppercase dark:text-white"
    >
        {$_('landingCountries', {
            values: { count: activeCountries.length },
        })}
    </h2>
</div>

<ul class="grid grid-cols-8 lg:grid-cols-12 gap-x-4 gap-y-8">
    {#each activeCountries as country}
        <li>
            <CountryFlag country="{country}" additionalClass="mx-auto" />
        </li>
    {/each}
</ul>
