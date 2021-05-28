<script>
    import { countryMap } from '../country'
    import { _ } from '../i18n'

    export let country = ''
    export let size = 64
    export let additionalClass = ''
</script>

<img
    src="https://www.countryflags.io/{country.toLowerCase()}/shiny/{size}.png"
    alt="{countryMap[country]}
    {$_('flag')}"
    title="{countryMap[country]}"
    class="{additionalClass}" />
