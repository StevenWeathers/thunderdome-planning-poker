<script>
    import { countryMap } from '../../country.js'
    import { _ } from '../../i18n.js'

    export let country = ''
    export let additionalClass = ''
    export let width = '64'
    export let height = '48'
</script>

<img
    src="https://flagcdn.com/64x48/{country.toLowerCase()}.png"
    srcset="https://flagcdn.com/128x96/{country.toLowerCase()}.png 2x,
    https://flagcdn.com/192x144/{country.toLowerCase()}.png 3x"
    width="{width}"
    height="{height}"
    alt="{countryMap[country]}
    {$_('flag')}"
    title="{countryMap[country]}"
    class="{additionalClass}"
/>
