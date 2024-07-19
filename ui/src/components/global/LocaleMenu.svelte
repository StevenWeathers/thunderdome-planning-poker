<script lang="ts">
  import { locales } from '../../config';
  import { createEventDispatcher } from 'svelte';
  import WorldIcon from '../icons/WorldIcon.svelte';

  export let currentPage;
  export let selectedLocale = 'en';

  const supportedLocales = [];
  let showMenu = false;

  function toggleMenu() {
    showMenu = !showMenu;
  }

  function pageChangeHandler() {
    if (showMenu === true) {
      toggleMenu();
    }
  }

  $: {
    if (typeof currentPage !== 'undefined') {
      pageChangeHandler();
    }
  }

  for (const [key, value] of Object.entries(locales)) {
    supportedLocales.push({
      name: value,
      value: key,
    });
  }

  const dispatch = createEventDispatcher();

  const switchLocale = locale => () => {
    toggleMenu();
    dispatch('locale-changed', locale);
  };
</script>

<div class="relative z-10">
  <label class="sr-only">Locale</label>

  <button
    class="relative z-10 flex h-8 w-8 items-center justify-center rounded-lg text-slate-700 hover:bg-slate-100 dark:text-slate-400 dark:hover:bg-slate-700"
    aria-label="Locale"
    type="button"
    on:click="{toggleMenu}"
  >
    <WorldIcon />
  </button>

  {#if showMenu}
    <ul
      class="absolute left-1/2 top-full mt-3 w-36 -translate-x-1/2 space-y-1 rounded-xl bg-white p-3 text-sm font-medium shadow-md shadow-black/5 ring-1 ring-black/5 dark:bg-slate-700 dark:ring-white/5"
      tabindex="0"
    >
      {#each supportedLocales as locale}
        <li
          class="flex cursor-pointer select-none items-center rounded-[0.625rem] p-1 text-slate-700
                         dark:text-slate-300 text-slate-900 dark:hover:text-white hover:bg-slate-100
                          dark:hover:bg-slate-900/40 {locale.value ===
          selectedLocale
            ? 'bg-indigo-600 text-white'
            : ''}"
          tabindex="-1"
        >
          <div class="ms-3">
            <button on:click="{switchLocale(locale.value)}"
              >{locale.name}</button
            >
          </div>
        </li>
      {/each}
    </ul>
  {/if}
</div>
