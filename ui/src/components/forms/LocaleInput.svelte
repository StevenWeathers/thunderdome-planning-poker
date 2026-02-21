<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { locales } from '../../config';
  import SelectInput from './SelectInput.svelte';

  interface Props {
    selectedLocale?: string;
    class?: string;
  }

  let { selectedLocale = 'en', class: klass = '' }: Props = $props();

  const supportedLocales: { name: string; value: string }[] = [];

  for (const [key, value] of Object.entries(locales)) {
    supportedLocales.push({
      name: value,
      value: key,
    });
  }

  const dispatch = createEventDispatcher();

  function switchLocale(event: Event) {
    event.preventDefault();

    dispatch('locale-changed', (event.target as HTMLSelectElement).value);
  }
</script>

<div class="{klass} inline-block">
  <SelectInput name="locale" onchange={switchLocale} value={selectedLocale}>
    {#each supportedLocales as locale}
      <option value={locale.value}>{locale.name}</option>
    {/each}
  </SelectInput>
</div>
