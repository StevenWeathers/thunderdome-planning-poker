<script>
    import { countryMap } from '../country'
    import { _ } from '../i18n'

    export let country = ''
    export let additionalClass = ''
</script>

<img
    src="https://flagcdn.com/64x48/{country.toLowerCase()}.png"
    srcset="https://flagcdn.com/128x96/{country.toLowerCase()}.png 2x,
    https://flagcdn.com/192x144/{country.toLowerCase()}.png 3x"
    width="64"
    height="48"
    alt="{countryMap[country]}
    {$_('flag')}"
    title="{countryMap[country]}"
    class="{additionalClass}"
/>
