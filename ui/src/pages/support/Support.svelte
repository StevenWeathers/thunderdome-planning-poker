<script lang="ts">
  import PageLayout from '../../components/global/PageLayout.svelte';
  import GithubIcon from '../../components/icons/Github.svelte';
  import LL from '../../i18n/i18n-svelte';
  import { AppConfig } from '../../config';

  const { AppVersion, RepoURL } = AppConfig;
</script>

<svelte:head>
  <title>{$LL.support()} | {$LL.appName()}</title>
</svelte:head>

<PageLayout>
  <div class="dark:text-gray-200">
    <div class="mb-4">
      <h1
        class="text-3xl font-semibold font-rajdhani dark:text-white uppercase"
      >
        {$LL.support()}
      </h1>
    </div>

    <div class="mb-8">
      <h2 class="py-2 text-2xl font-semibold font-rajdhani dark:text-white">
        Have a bug or issue you want to report?
      </h2>
      <p>
        Create an issue on the
        <a
          href="{RepoURL}/issues"
          target="_blank"
          class="text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600"
        >
          <GithubIcon />
          Github Repo</a
        >
        and mention the <span class="font-semibold">current version</span>
        <span
          class="ms-2 font-bold dark:text-gray-800 bg-lime-800 dark:bg-lime-200 py-1 px-2 rounded inline-block"
          >{AppVersion}</span
        >
      </p>
    </div>
    <div>
      <h2
        class="py-2 text-2xl font-semibold font-rajdhani dark:text-white uppercase"
      >
        Contact Us
      </h2>
      <div
        class="py-4 font-rajdhani text-2xl text-center font-bold text-green-500"
      >
        coming soon...
      </div>
    </div>
  </div>
</PageLayout>
