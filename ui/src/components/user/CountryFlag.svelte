<script lang="ts">
  import { countryMap } from '../../country';
  import LL from '../../i18n/i18n-svelte';

  interface Props {
    country?: string;
    additionalClass?: string;
    width?: string;
    height?: string;
  }

  let { country = '', additionalClass = '', width = '48', height = '36' }: Props = $props();
</script>

<img
  src="https://flagcdn.com/64x48/{country.toLowerCase()}.png"
  srcset="https://flagcdn.com/128x96/{country.toLowerCase()}.png 2x,
    https://flagcdn.com/192x144/{country.toLowerCase()}.png 3x"
  {width}
  {height}
  alt="{countryMap[country]}
    {$LL.flag()}"
  title={countryMap[country]}
  class={additionalClass}
/>
