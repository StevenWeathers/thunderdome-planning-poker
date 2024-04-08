<script lang="ts">
  export let title = '';
  export let createBtnEnabled = true;
  export let createBtnText = '';
  export let createButtonHandler = () => {};
  export let createBtnTestId = '';
  export let searchEnabled = false;
  export let searchPlaceholder = '';
  export let searchHandler = term => {};

  let searchTerm = '';

  function onSearchSubmit(e) {
    e.preventDefault();

    searchHandler(searchTerm);
  }
</script>

<div
  class="flex flex-col md:flex-row items-stretch md:items-center md:space-x-3 space-y-3 md:space-y-0 justify-between mx-4 py-4"
>
  <div class="w-full md:w-1/2">
    <div class="flex gap-4 lg:gap-8 items-center">
      <h5>
        <span
          class="dark:text-white font-rajdhani font-semibold text-xl lg:text-2xl"
          data-testid="tablenav-title">{title}</span
        >
      </h5>
      {#if searchEnabled}
        <form class="flex items-center" on:submit="{onSearchSubmit}">
          <div class="relative w-full">
            <div
              class="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none"
            >
              <svg
                aria-hidden="true"
                class="w-5 h-5 text-gray-500 dark:text-gray-400"
                fill="currentColor"
                viewBox="0 0 20 20"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  fill-rule="evenodd"
                  clip-rule="evenodd"
                  d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z"
                ></path>
              </svg>
            </div>
            <input
              type="text"
              placeholder="{searchPlaceholder}"
              bind:value="{searchTerm}"
              class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full pl-10 p-2 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
            />
          </div>
        </form>
      {/if}
    </div>
  </div>
  <div
    class="w-full md:w-auto flex flex-col md:flex-row space-y-2 md:space-y-0 items-stretch md:items-center justify-end md:space-x-3 flex-shrink-0"
  >
    <slot />
    {#if createBtnEnabled}
      <button
        type="button"
        on:click="{createButtonHandler}"
        data-testid="{createBtnTestId}"
        class="flex items-center justify-center text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:ring-blue-300 font-medium rounded-lg text-sm px-4 py-2 dark:bg-blue-600 dark:hover:bg-blue-700 focus:outline-none dark:focus:ring-blue-800"
      >
        <svg
          class="h-3.5 w-3.5 mr-1.5 -ml-1"
          fill="currentColor"
          viewBox="0 0 20 20"
          xmlns="http://www.w3.org/2000/svg"
          aria-hidden="true"
        >
          <path
            clip-rule="evenodd"
            fill-rule="evenodd"
            d="M10 3a1 1 0 011 1v5h5a1 1 0 110 2h-5v5a1 1 0 11-2 0v-5H4a1 1 0 110-2h5V4a1 1 0 011-1z"
          ></path>
        </svg>
        {createBtnText}
      </button>
    {/if}
  </div>
</div>
