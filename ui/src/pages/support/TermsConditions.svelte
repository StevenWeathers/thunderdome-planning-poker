<script lang="ts">
  import PageLayout from '../../components/PageLayout.svelte';
  import LL from '../../i18n/i18n-svelte';
  import { onMount } from 'svelte';
  let htmlContent = $state('');

  onMount(() => {
    window.scrollTo(0, 0);
    fetch('/content/terms-and-conditions.html')
      .then(res => res.text())
      .then(data => {
        htmlContent = data;
      });
  });
</script>

<svelte:head>
  <title>{$LL.termsConditions()} | {$LL.appName()}</title>
</svelte:head>

<PageLayout>
  <div class="dark:text-gray-200">
    <h1 class="py-2 text-2xl">{$LL.termsConditions()}</h1>
    <p class="py-4 italic text-gray-500 dark:text-gray-400">Updated at 2024-03-27</p>

    {@html htmlContent}
  </div>
</PageLayout>
