<script>
    import { countryMap } from '../country'

    export let country = ''
    export let size = 64
    export let additionalClass = ''
</script>

<img
    src="https://www.countryflags.io/{country.toLowerCase()}/shiny/{size}.png"
    alt="{countryMap[country]} Flag"
    title="{countryMap[country]}"
    class="{additionalClass}" />
