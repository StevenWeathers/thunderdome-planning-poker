<script lang="ts">
  import SideNavigation from './global/SideNavigation.svelte';
  import type { ComponentType } from 'svelte';

  // Define the page item type
  export interface PageItem {
    name: string;
    label: string;
    path: string;
    icon: ComponentType;
    enabled: boolean;
  }

  // Define the menu type for the SideNavigation component
  export type MenuType = 'admin' | 'user' | 'organization' | 'team' | 'department' | string;

  interface Props {
    pages: PageItem[];
    activePage?: string;
    menuType?: MenuType;
    expanded?: boolean;
    children?: import('svelte').Snippet;
  }

  let { 
    pages, 
    activePage = '', 
    menuType = 'default',
    expanded,
    children 
  }: Props = $props();

  // Filter enabled pages
  let enabledPages = $derived(pages.filter(page => page.enabled));
</script>

<section class="flex min-h-screen">
  <SideNavigation
    menuItems={enabledPages}
    activePage={activePage}
    menuType={menuType}
    {expanded}
  />
  <div class="flex-1 px-4 py-4 md:py-6 md:px-6 lg:py-8 lg:px-8">
    {@render children?.()}
  </div>
</section>