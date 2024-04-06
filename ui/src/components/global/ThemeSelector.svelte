<script lang="ts">
  export let currentPage;

  let selectedTheme = localStorage.getItem('theme') || 'auto';
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
    on:click="{toggleMenu}"
  >
    <svg
      aria-hidden="true"
      viewBox="0 0 16 16"
      class="h-5 w-5 dark:hidden"
      fill="currentColor"
    >
      <path
        fill-rule="evenodd"
        clip-rule="evenodd"
        d="M7 1a1 1 0 0 1 2 0v1a1 1 0 1 1-2 0V1Zm4 7a3 3 0 1 1-6 0 3 3 0 0 1 6 0Zm2.657-5.657a1 1 0 0 0-1.414 0l-.707.707a1 1 0 0 0 1.414 1.414l.707-.707a1 1 0 0 0 0-1.414Zm-1.415 11.313-.707-.707a1 1 0 0 1 1.415-1.415l.707.708a1 1 0 0 1-1.415 1.414ZM16 7.999a1 1 0 0 0-1-1h-1a1 1 0 1 0 0 2h1a1 1 0 0 0 1-1ZM7 14a1 1 0 1 1 2 0v1a1 1 0 1 1-2 0v-1Zm-2.536-2.464a1 1 0 0 0-1.414 0l-.707.707a1 1 0 0 0 1.414 1.414l.707-.707a1 1 0 0 0 0-1.414Zm0-8.486A1 1 0 0 1 3.05 4.464l-.707-.707a1 1 0 0 1 1.414-1.414l.707.707ZM3 8a1 1 0 0 0-1-1H1a1 1 0 0 0 0 2h1a1 1 0 0 0 1-1Z"
      ></path>
    </svg>
    <svg
      aria-hidden="true"
      viewBox="0 0 16 16"
      class="hidden h-5 w-5 dark:block"
      fill="currentColor"
    >
      <path
        fill-rule="evenodd"
        clip-rule="evenodd"
        d="M7.23 3.333C7.757 2.905 7.68 2 7 2a6 6 0 1 0 0 12c.68 0 .758-.905.23-1.332A5.989 5.989 0 0 1 5 8c0-1.885.87-3.568 2.23-4.668ZM12 5a1 1 0 0 1 1 1 1 1 0 0 0 1 1 1 1 0 1 1 0 2 1 1 0 0 0-1 1 1 1 0 1 1-2 0 1 1 0 0 0-1-1 1 1 0 1 1 0-2 1 1 0 0 0 1-1 1 1 0 0 1 1-1Z"
      ></path>
    </svg>
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
        on:click="{setTheme('light')}"
      >
        <div
          class="rounded-md bg-white p-1 shadow ring-1 ring-slate-900/5 dark:bg-slate-700 dark:ring-inset dark:ring-white/5 text-slate-700 dark:text-slate-400 text-slate-900"
        >
          <svg
            aria-hidden="true"
            viewBox="0 0 16 16"
            class="h-4 w-4"
            fill="currentColor"
          >
            <path
              fill-rule="evenodd"
              clip-rule="evenodd"
              d="M7 1a1 1 0 0 1 2 0v1a1 1 0 1 1-2 0V1Zm4 7a3 3 0 1 1-6 0 3 3 0 0 1 6 0Zm2.657-5.657a1 1 0 0 0-1.414 0l-.707.707a1 1 0 0 0 1.414 1.414l.707-.707a1 1 0 0 0 0-1.414Zm-1.415 11.313-.707-.707a1 1 0 0 1 1.415-1.415l.707.708a1 1 0 0 1-1.415 1.414ZM16 7.999a1 1 0 0 0-1-1h-1a1 1 0 1 0 0 2h1a1 1 0 0 0 1-1ZM7 14a1 1 0 1 1 2 0v1a1 1 0 1 1-2 0v-1Zm-2.536-2.464a1 1 0 0 0-1.414 0l-.707.707a1 1 0 0 0 1.414 1.414l.707-.707a1 1 0 0 0 0-1.414Zm0-8.486A1 1 0 0 1 3.05 4.464l-.707-.707a1 1 0 0 1 1.414-1.414l.707.707ZM3 8a1 1 0 0 0-1-1H1a1 1 0 0 0 0 2h1a1 1 0 0 0 1-1Z"
            ></path>
          </svg>
        </div>
        <div class="ms-3">Light</div>
      </li>
      <li
        class="flex cursor-pointer select-none items-center rounded-[0.625rem] p-1 text-slate-700 dark:text-slate-300 {selectedTheme ===
        'dark'
          ? 'bg-indigo-600 text-white'
          : ''}"
        tabindex="-1"
        on:click="{setTheme('dark')}"
      >
        <div
          class="rounded-md bg-white p-1 shadow ring-1 ring-slate-900/5 dark:bg-slate-700 dark:ring-inset dark:ring-white/5"
        >
          <svg
            aria-hidden="true"
            viewBox="0 0 16 16"
            class="h-4 w-4"
            fill="currentColor"
          >
            <path
              fill-rule="evenodd"
              clip-rule="evenodd"
              d="M7.23 3.333C7.757 2.905 7.68 2 7 2a6 6 0 1 0 0 12c.68 0 .758-.905.23-1.332A5.989 5.989 0 0 1 5 8c0-1.885.87-3.568 2.23-4.668ZM12 5a1 1 0 0 1 1 1 1 1 0 0 0 1 1 1 1 0 1 1 0 2 1 1 0 0 0-1 1 1 1 0 1 1-2 0 1 1 0 0 0-1-1 1 1 0 1 1 0-2 1 1 0 0 0 1-1 1 1 0 0 1 1-1Z"
            ></path>
          </svg>
        </div>
        <div class="ms-3">Dark</div>
      </li>
      <li
        class="flex cursor-pointer select-none items-center rounded-[0.625rem] p-1 text-slate-700 dark:text-slate-300 text-slate-900 dark:hover:text-white hover:bg-slate-100 dark:hover:bg-slate-900/40 {selectedTheme ===
        'auto'
          ? 'bg-indigo-600 text-white'
          : ''}"
        tabindex="-1"
        on:click="{setTheme('auto')}"
      >
        <div
          class="rounded-md bg-white p-1 shadow ring-1 ring-slate-900/5 dark:bg-slate-700 dark:ring-inset dark:ring-white/5 text-slate-700 dark:text-slate-400 text-slate-900"
        >
          <svg
            aria-hidden="true"
            viewBox="0 0 16 16"
            class="h-4 w-4"
            fill="currentColor"
          >
            <path
              fill-rule="evenodd"
              clip-rule="evenodd"
              d="M1 4a3 3 0 0 1 3-3h8a3 3 0 0 1 3 3v4a3 3 0 0 1-3 3h-1.5l.31 1.242c.084.333.36.573.63.808.091.08.182.158.264.24A1 1 0 0 1 11 15H5a1 1 0 0 1-.704-1.71c.082-.082.173-.16.264-.24.27-.235.546-.475.63-.808L5.5 11H4a3 3 0 0 1-3-3V4Zm3-1a1 1 0 0 0-1 1v4a1 1 0 0 0 1 1h8a1 1 0 0 0 1-1V4a1 1 0 0 0-1-1H4Z"
            ></path>
          </svg>
        </div>
        <div class="ms-3">System</div>
      </li>
    </ul>
  {/if}
</div>
