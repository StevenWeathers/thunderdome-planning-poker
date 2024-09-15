<script lang="ts">
  import { countryMap } from '../../country';
  import LL from '../../i18n/i18n-svelte';

  export let country = '';
  export let additionalClass = '';
  export let width = '48';
  export let height = '36';
</script>

<img
  src="https://flagcdn.com/64x48/{country.toLowerCase()}.png"
  srcset="https://flagcdn.com/128x96/{country.toLowerCase()}.png 2x,
    https://flagcdn.com/192x144/{country.toLowerCase()}.png 3x"
  width="{width}"
  height="{height}"
  alt="{countryMap[country]}
    {$LL.flag()}"
  title="{countryMap[country]}"
  class="{additionalClass}"
/>
