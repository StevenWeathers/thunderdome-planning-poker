<script>
    import CountryFlag from './CountryFlag.svelte'
    import { _ } from '../i18n'

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

<section>
    <div class="container mx-auto px-4 py-6 lg:py-10">
        <div class="text-center mb-6">
            <h2 class="text-4xl font-bold">
                {$_('landingCountries', {
                    values: { count: activeCountries.length },
                })}
            </h2>
        </div>

        <ul class="grid grid-cols-8 lg:grid-cols-12 gap-x-4">
            {#each activeCountries as country}
                <li>
                    <CountryFlag
                        country="{country}"
                        additionalClass="mx-auto"
                    />
                </li>
            {/each}
        </ul>
    </div>
</section>
