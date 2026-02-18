<script lang="ts">
  import SubMenu from './SubMenu.svelte';
  import SubMenuItem from './SubMenuItem.svelte';
  import { locales } from '../../config';
  import { Globe } from '@lucide/svelte';
  import type { Locales } from '../../i18n/i18n-types';

  interface Props {
    selectedLocale: Locales;
    update: (l: Locales) => void;
  }

  interface Locale {
    name: string;
    value: Locales;
  }

  let { update, selectedLocale = 'en' }: Props = $props();

  const supportedLocales: Array<Locale> = [];
  for (const [key, value] of Object.entries(locales)) {
    supportedLocales.push({
      name: value,
      value: key as Locales,
    });
  }

  const switchLocale = (locale: Locales, toggleSubmenu?: () => void) => () => {
    toggleSubmenu?.();
    update(locale);
  };
</script>

<SubMenu class="w-36" relativeClass="z-10">
  {#snippet button({ toggleSubmenu })}
    <button
      class="relative z-10 flex h-8 w-8 items-center justify-center rounded-lg text-slate-700 hover:bg-slate-100 dark:text-slate-400 dark:hover:bg-slate-700"
      aria-label="Locale"
      type="button"
      onclick={toggleSubmenu}
    >
      <span class="sr-only">Locale</span><Globe />
    </button>
  {/snippet}

  {#snippet children({ toggleSubmenu })}
    {#each supportedLocales as locale}
      <SubMenuItem
        onClickHandler={switchLocale(locale.value, toggleSubmenu)}
        testId="locale-{locale.name}"
        label={locale.name}
        active={locale.value === selectedLocale}
      />
    {/each}
  {/snippet}
</SubMenu>
