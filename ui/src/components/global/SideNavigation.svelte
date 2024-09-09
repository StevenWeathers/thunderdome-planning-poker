<script lang="ts">
  import { Menu, X } from 'lucide-svelte';

  let isCollapsed = true;

  export let menuItems = [];
  export let activePage = '';
  export let menuType = '';
</script>

<div
  class="bg-gray-800 text-white min-h-screen {isCollapsed
    ? 'w-16'
    : 'w-64'} transition-all duration-300 ease-in-out"
>
  <div class="flex justify-end p-4">
    <button on:click="{() => (isCollapsed = !isCollapsed)}" class="text-white">
      {#if isCollapsed}
        <Menu size="{24}" />
      {:else}
        <X size="{24}" />
      {/if}
    </button>
  </div>
  <nav>
    <ul data-testid="{menuType}-nav">
      {#each menuItems as item}
        {#if item.enabled}
          <li class="mb-2">
            <a
              href="{item.path}"
              class="flex items-center p-2 hover:bg-gray-700 {isCollapsed
                ? 'justify-center'
                : 'px-4'}"
              data-testid="{menuType}-nav-item"
              class:bg-indigo-500="{activePage ===
                item.name.toLowerCase().replace(' ', '-')}"
              title="{item.label}"
            >
              <span class="mr-3">
                <svelte:component this="{item.icon}" size="{24}" />
              </span>
              <span class:hidden="{isCollapsed}">{item.label}</span>
            </a>
          </li>
        {/if}
      {/each}
    </ul>
  </nav>
</div>
