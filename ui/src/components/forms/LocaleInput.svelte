<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { locales } from '../../config';
  import SelectInput from './SelectInput.svelte';

  interface Props {
    selectedLocale?: string;
    class?: string;
  }

  let { selectedLocale = 'en', class: klass = '' }: Props = $props();
  
  const supportedLocales = [];

  for (const [key, value] of Object.entries(locales)) {
    supportedLocales.push({
      name: value,
      value: key,
    });
  }

  const dispatch = createEventDispatcher();

  function switchLocale(event) {
    event.preventDefault();

    dispatch('locale-changed', event.target.value);
  }
</script>

<div class="{klass} inline-block">
  <SelectInput
    name="locale"
    on:change={switchLocale}
    value={selectedLocale}
  >
    {#each supportedLocales as locale}
      <option value="{locale.value}">{locale.name}</option>
    {/each}
  </SelectInput>
</div>
