<script lang="ts">
  import GithubIcon from '../icons/Github.svelte';
  import { AppConfig } from '../../config';
  import LL from '../../i18n/i18n-svelte';

  const { AppVersion } = AppConfig;
  const footerLinkClasses =
    'no-underline text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600';
</script>

<footer
  class="w-full border-t px-16 py-10 text-center font-rajdhani text-xl dark:text-white dark:border-gray-800"
  style="box-shadow: inset 0 7px 30px -7px rgba(0,0,0,0.1);"
>
  <GithubIcon />
  <a
    href="https://github.com/StevenWeathers/thunderdome-planning-poker"
    class="{footerLinkClasses}"
  >
    {$LL.appName()}
  </a>
  {@html $LL.footerAuthoredBy({
    authorOpen: `<a href="http://stevenweathers.com" class="${footerLinkClasses}">`,
    authorClose: `</a>`,
  })}
  <br />
  {@html $LL.footerLicense({
    licenseOpen: `<a href="http://www.apache.org/licenses/" class="${footerLinkClasses}">`,
    licenseClose: `</a>`,
  })}
  <br />
  <div class="text-sm text-gray-500 dark:text-gray-400">
    {$LL.appVersion({ version: AppVersion })}
  </div>
</footer>
