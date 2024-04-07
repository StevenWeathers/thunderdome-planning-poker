<script lang="ts">
  import PageLayout from '../../components/global/PageLayout.svelte';
  import LL from '../../i18n/i18n-svelte';
</script>

<svelte:head>
  <title>{$LL.termsConditions()} | {$LL.appName()}</title>
</svelte:head>

<PageLayout>
  {#await import('../../components/support/TermsServiceContent.svelte') then Module}
    <Module.default />
  {/await}
</PageLayout>
