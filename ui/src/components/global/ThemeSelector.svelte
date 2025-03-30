<script lang="ts">
  import { run } from 'svelte/legacy';

  import { MonitorCog, MoonStar, Sun } from 'lucide-svelte';

  interface Props {
    currentPage: any;
  }

  let { currentPage }: Props = $props();

  let selectedTheme = $state(localStorage.getItem('theme') || 'auto');
  let showMenu = $state(false);

  function toggleMenu() {
    showMenu = !showMenu;
  }

  function pageChangeHandler() {
    if (showMenu === true) {
      toggleMenu();
    }
  }

  run(() => {
    if (typeof currentPage !== 'undefined') {
      pageChangeHandler();
    }
  });

  const setTheme = theme => () => {
    selectedTheme = theme;
    if (selectedTheme !== 'auto') {
      localStorage.setItem('theme', selectedTheme);
    } else {
      localStorage.removeItem('theme');
    }
    window.setTheme();
    toggleMenu();
  };
</script>

<div class="relative z-10">
  <span class="sr-only">Theme</span>
  <button
    class="relative z-10 flex h-8 w-8 items-center justify-center rounded-lg text-slate-700 hover:bg-slate-100 hover:text-slate-900 dark:text-slate-400 dark:hover:bg-slate-700"
    aria-label="Theme"
    type="button"
    onclick={toggleMenu}
  >
    <Sun class="dark:hidden h-5 w-5" />
    <MoonStar class="hidden h-5 w-5 dark:block" />
  </button>

  {#if showMenu}
    <ul
      class="absolute left-1/2 top-full mt-3 w-36 -translate-x-1/2 space-y-1 rounded-xl bg-white p-3 text-sm font-medium shadow-md shadow-black/5 ring-1 ring-black/5 dark:bg-slate-700 dark:ring-white/5"
      tabindex="0"
    >
      <li
        class="flex cursor-pointer select-none items-center rounded-[0.625rem] p-1 text-slate-700 dark:text-slate-300 text-slate-900 dark:hover:text-white hover:bg-slate-100 dark:hover:bg-slate-900/40 {selectedTheme ===
        'light'
          ? 'bg-indigo-600 text-white'
          : ''}"
        tabindex="-1"
        onclick={setTheme('light')}
      >
        <div
          class="rounded-md bg-white p-1 shadow ring-1 ring-slate-900/5 dark:bg-slate-700 dark:ring-inset dark:ring-white/5 text-slate-700 dark:text-slate-400 text-slate-900"
        >
          <Sun />
        </div>
        <div class="ms-3">Light</div>
      </li>
      <li
        class="flex cursor-pointer select-none items-center rounded-[0.625rem] p-1 text-slate-700 dark:text-slate-300 text-slate-900 dark:hover:text-white hover:bg-slate-100 dark:hover:bg-slate-900/40 {selectedTheme ===
        'dark'
          ? 'bg-indigo-600 text-white'
          : ''}"
        tabindex="-1"
        onclick={setTheme('dark')}
      >
        <div
          class="rounded-md bg-white p-1 shadow ring-1 ring-slate-900/5 dark:bg-slate-700 dark:ring-inset dark:ring-white/5 text-slate-700 dark:text-slate-400 text-slate-900"
        >
          <MoonStar />
        </div>
        <div class="ms-3">Dark</div>
      </li>
      <li
        class="flex cursor-pointer select-none items-center rounded-[0.625rem] p-1 text-slate-700 dark:text-slate-300 text-slate-900 dark:hover:text-white hover:bg-slate-100 dark:hover:bg-slate-900/40 {selectedTheme ===
        'auto'
          ? 'bg-indigo-600 text-white'
          : ''}"
        tabindex="-1"
        onclick={setTheme('auto')}
      >
        <div
          class="rounded-md bg-white p-1 shadow ring-1 ring-slate-900/5 dark:bg-slate-700 dark:ring-inset dark:ring-white/5 text-slate-700 dark:text-slate-400 text-slate-900"
        >
          <MonitorCog />
        </div>
        <div class="ms-3">System</div>
      </li>
    </ul>
  {/if}
</div>
