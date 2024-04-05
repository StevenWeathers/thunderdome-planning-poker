<script lang="ts">
  import { locales } from '../../config';
  import { createEventDispatcher } from 'svelte';

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
    class="relative z-10 flex h-8 w-8 items-center justify-center rounded-lg text-slate-700 hover:bg-slate-100 hover:text-slate-900 dark:text-slate-400 dark:hover:bg-slate-700"
    aria-label="Locale"
    type="button"
    on:click="{toggleMenu}"
  >
    <svg
      xmlns="http://www.w3.org/2000/svg"
      fill="none"
      viewBox="0 0 24 24"
      stroke-width="1.5"
      stroke="currentColor"
      class="w-6 h-6"
    >
      <path
        stroke-linecap="round"
        stroke-linejoin="round"
        d="M12 21a9.004 9.004 0 0 0 8.716-6.747M12 21a9.004 9.004 0 0 1-8.716-6.747M12 21c2.485 0 4.5-4.03 4.5-9S14.485 3 12 3m0 18c-2.485 0-4.5-4.03-4.5-9S9.515 3 12 3m0 0a8.997 8.997 0 0 1 7.843 4.582M12 3a8.997 8.997 0 0 0-7.843 4.582m15.686 0A11.953 11.953 0 0 1 12 10.5c-2.998 0-5.74-1.1-7.843-2.918m15.686 0A8.959 8.959 0 0 1 21 12c0 .778-.099 1.533-.284 2.253m0 0A17.919 17.919 0 0 1 12 16.5c-3.162 0-6.133-.815-8.716-2.247m0 0A9.015 9.015 0 0 1 3 12c0-1.605.42-3.113 1.157-4.418"
      ></path>
    </svg>
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
          <div class="ml-3">
            <button on:click="{switchLocale(locale.value)}"
              >{locale.name}</button
            >
          </div>
        </li>
      {/each}
    </ul>
  {/if}
</div>
