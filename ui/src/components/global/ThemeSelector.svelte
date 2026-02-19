<script lang="ts">
  import { MonitorCog, MoonStar, Sun } from '@lucide/svelte';
  import SubMenu from './SubMenu.svelte';
  import SubMenuItem from './SubMenuItem.svelte';

  let selectedTheme = localStorage.getItem('theme') || 'auto';

  const setTheme = (theme: string, toggleSubmenu?: () => void) => () => {
    selectedTheme = theme;
    if (selectedTheme !== 'auto') {
      localStorage.setItem('theme', selectedTheme);
    } else {
      localStorage.removeItem('theme');
    }
    (window as any).setTheme();
    toggleSubmenu?.();
  };
</script>

<SubMenu relativeClass="z-10">
  {#snippet button({ toggleSubmenu })}
    <button
      class="relative z-10 flex h-8 w-8 items-center justify-center rounded-lg text-slate-700 hover:bg-slate-100 hover:text-slate-900 dark:text-slate-400 dark:hover:bg-slate-700"
      aria-label="Theme"
      type="button"
      onclick={toggleSubmenu}
    >
      <span class="sr-only">Theme</span>
      <Sun class="dark:hidden h-5 w-5" />
      <MoonStar class="hidden h-5 w-5 dark:block" />
    </button>
  {/snippet}

  {#snippet children({ toggleSubmenu })}
    <SubMenuItem
      onClickHandler={setTheme('light', toggleSubmenu)}
      testId="theme-light"
      icon={Sun}
      label="Light"
      active={selectedTheme === 'light'}
    />
    <SubMenuItem
      onClickHandler={setTheme('dark', toggleSubmenu)}
      testId="theme-dark"
      icon={MoonStar}
      label="Dark"
      active={selectedTheme === 'dark'}
    />
    <SubMenuItem
      onClickHandler={setTheme('auto', toggleSubmenu)}
      testId="theme-auto"
      icon={MonitorCog}
      label="System"
      active={selectedTheme === 'auto'}
    />
  {/snippet}
</SubMenu>
