<script lang="ts">
  import GithubIcon from '../icons/Github.svelte';
  import { AppConfig, appRoutes } from '../../config';
  import LL from '../../i18n/i18n-svelte';
  import FullLogoVertical from './FullLogoVertical.svelte';
  import FullLogoVerticalDarkText from './FullLogoVerticalDarkText.svelte';

  const { AppVersion, RepoURL, PathPrefix } = AppConfig;
  const footerLinkClasses =
    'no-underline text-blue-500 hover:text-blue-800 dark:text-sky-400 dark:hover:text-sky-600';
</script>

<footer
  class="w-full border-t px-16 py-10 text-center dark:text-white dark:border-gray-800"
  style="box-shadow: inset 0 7px 30px -7px rgba(0,0,0,0.1);"
>
  <div class="">
    <div class="md:flex md:items-center md:justify-between">
      <a
        href="{appRoutes.landing}"
        class="hidden md:inline-block flex items-center mb-4 md:mb-0 space-x-3 rtl:space-x-reverse"
      >
        <FullLogoVertical
          class="hidden h-8 lg:h-12 dark:block grayscale hover:grayscale-0"
        />
        <FullLogoVerticalDarkText
          class="h-8 lg:h-12 dark:hidden grayscale hover:grayscale-0"
        />
      </a>
      <ul
        class="text-sm md:text-base flex flex-wrap items-center mb-6 font-medium text-gray-500 md:mb-0 dark:text-gray-400"
      >
        <li>
          <a
            href="{appRoutes.subscriptionPricing}"
            class="hover:underline me-4 md:me-6">Pricing</a
          >
        </li>
        <li>
          <a
            href="{RepoURL}/blob/main/docs/GUIDE.md"
            target="_blank"
            class="hover:underline me-4 md:me-6">{$LL.userGuide()}</a
          >
        </li>
        <li>
          <a
            href="{appRoutes.privacyPolicy}"
            class="hover:underline me-4 md:me-6">{$LL.privacyPolicy()}</a
          >
        </li>
        <li>
          <a
            href="{appRoutes.termsConditions}"
            class="hover:underline me-4 md:me-6">{$LL.termsConditions()}</a
          >
        </li>
        <li>
          <a href="{appRoutes.support}" class="hover:underline"
            >{$LL.support()}</a
          >
        </li>
      </ul>
    </div>
    <hr class="my-6 border-gray-200 sm:mx-auto dark:border-gray-700 lg:my-8" />
    <div class="text-gray-500 sm:text-center dark:text-gray-400 font-rajdhani">
      <GithubIcon />
      <a href="{RepoURL}" class="{footerLinkClasses}" target="_blank">
        {$LL.appName()}
      </a>
      {@html $LL.footerAuthoredBy({
        authorOpen: `<a href="http://stevenweathers.com" class="${footerLinkClasses}" target="_blank">`,
        authorClose: `</a>`,
      })}
      <br />
      {@html $LL.footerLicense({
        licenseOpen: `<a href="http://www.apache.org/licenses/" class="${footerLinkClasses}" target="_blank">`,
        licenseClose: `</a>`,
      })}
      <br />
      <div class="text-sm text-gray-500 dark:text-gray-400">
        {$LL.appVersion({ version: AppVersion })}
      </div>
    </div>
  </div>
</footer>
