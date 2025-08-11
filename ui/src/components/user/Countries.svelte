<script lang="ts">
  import CountryFlag from './CountryFlag.svelte';
  import type { ApiClient } from '../../types/apiclient';

  interface Props {
    xfetch: ApiClient;
  }

  let { xfetch }: Props = $props();

  let activeCountries = $state([]);

  xfetch('/api/active-countries')
    .then(res => res.json())
    .then(function (result) {
      activeCountries = result.data.sort();
    })
    .catch(function () {});
</script>

<div class="text-center">
  <div class="mx-auto title-line bg-yellow-thunder"></div>
  <h2
    class="text-4xl font-semibold mb-4 font-rajdhani uppercase dark:text-white"
  >
    Transforming teamwork in {activeCountries.length} countries and counting
  </h2>

  <ul class="flex flex-wrap gap-3">
    {#each activeCountries as country}
      <li class="flex-none w-12">
        <CountryFlag country={country} additionalClass="mx-auto" />
      </li>
    {/each}
  </ul>
</div>
