<script lang="ts">
  import GithubIcon from '../icons/Github.svelte';
  import { AppConfig, appRoutes } from '../../config';
  import LL from '../../i18n/i18n-svelte';
  import FullLogoVertical from '../logos/FullLogo.svelte';
  import FullLogoVerticalDarkText from '../logos/FullLogoLight.svelte';

  const { AppVersion, RepoURL, PathPrefix, SubscriptionsEnabled } = AppConfig;
  const footerLinkClasses =
    'no-underline text-indigo-600 dark:text-indigo-400 hover:text-yellow-thunder dark:hover:text-yellow-thunder transition-all duration-300 font-medium';
  const navLinkClasses =
    'text-gray-600 dark:text-gray-300 hover:text-indigo-600 dark:hover:text-indigo-400 transition-colors duration-300 relative after:absolute after:w-0 after:h-0.5 after:bg-indigo-600 after:dark:bg-indigo-400 after:start-0 after:-bottom-1 hover:after:w-full after:transition-all after:duration-300';
</script>

<footer
  class="relative w-full bg-gradient-to-br from-gray-50 to-gray-100 dark:from-gray-900 dark:to-gray-950 border-t border-gray-200 dark:border-gray-800 overflow-hidden"
>
  <!-- Subtle background pattern -->
  <div class="absolute inset-0 opacity-5 dark:opacity-10">
    <div
      class="absolute inset-0"
      style="background-image: radial-gradient(circle at 25% 25%, rgba(99, 102, 241, 0.3) 0%, transparent 50%), radial-gradient(circle at 75% 75%, rgba(99, 102, 241, 0.2) 0%, transparent 50%);"
    ></div>
  </div>

  <!-- Content -->
  <div class="relative">
    <!-- Main footer content -->
    <div class="max-w-7xl mx-auto px-6 sm:px-8 lg:px-16 py-10 lg:py-12">
      <div class="flex flex-col lg:flex-row lg:items-start lg:justify-between gap-8 lg:gap-16">
        <!-- Logo and branding section -->
        <div class="flex flex-col items-center lg:items-start space-y-4 lg:max-w-md">
          <a
            href={appRoutes.landing}
            class="group flex items-center space-x-3 rtl:space-x-reverse transition-transform duration-300 hover:scale-105"
          >
            <FullLogoVertical
              class="hidden h-10 lg:h-12 dark:block transition-all duration-500 grayscale group-hover:grayscale-0 group-hover:drop-shadow-lg"
            />
            <FullLogoVerticalDarkText
              class="h-10 lg:h-12 dark:hidden transition-all duration-500 grayscale group-hover:grayscale-0 group-hover:drop-shadow-lg"
            />
          </a>

          <!-- App description or tagline could go here -->
          <p class="text-gray-600 dark:text-gray-400 text-center lg:text-start leading-relaxed">
            Empowering teams with modern agile tools for better collaboration and delivery.
          </p>

          <!-- GitHub link with icon -->
          <div class="flex items-center space-x-2 rtl:space-x-reverse">
            <div
              class="flex items-center justify-center w-10 h-10 rounded-lg bg-gray-200 dark:bg-gray-800 transition-colors duration-300 group-hover:bg-indigo-100 dark:group-hover:bg-indigo-900/30"
            >
              <GithubIcon class="w-6 h-6 text-gray-600 dark:text-gray-400" />
            </div>
            <a href={RepoURL} class={footerLinkClasses} target="_blank">
              {$LL.githubRepository()}
            </a>
          </div>
        </div>

        <!-- Navigation links -->
        <div class="lg:flex-shrink-0">
          <h3
            class="font-semibold text-gray-900 dark:text-white uppercase tracking-wider mb-4 text-center lg:text-start"
          >
            Quick Links
          </h3>
          <nav>
            <ul class="grid grid-cols-2 lg:grid-cols-1 gap-3 text-center lg:text-start">
              {#if SubscriptionsEnabled}
                <li>
                  <a href={appRoutes.subscriptionPricing} class="{navLinkClasses} block py-1"> Pricing </a>
                </li>
              {/if}
              <li>
                <a href="{RepoURL}/blob/main/docs/GUIDE.md" target="_blank" class="{navLinkClasses} block py-1">
                  {$LL.userGuide()}
                </a>
              </li>
              <li>
                <a href={appRoutes.openSource} class="{navLinkClasses} block py-1">
                  {$LL.openSource()}
                </a>
              </li>
              <li>
                <a href={appRoutes.privacyPolicy} class="{navLinkClasses} block py-1">
                  {$LL.privacyPolicy()}
                </a>
              </li>
              <li>
                <a href={appRoutes.termsConditions} class="{navLinkClasses} block py-1">
                  {$LL.termsConditions()}
                </a>
              </li>
              <li>
                <a href={appRoutes.support} class="{navLinkClasses} block py-1">
                  {$LL.support()}
                </a>
              </li>
            </ul>
          </nav>
        </div>
      </div>
    </div>

    <!-- Bottom section with divider -->
    <div class="border-t border-gray-200 dark:border-gray-800 bg-gray-50/50 dark:bg-gray-900/50 backdrop-blur-sm">
      <div class="max-w-7xl mx-auto px-6 sm:px-8 lg:px-16 py-6">
        <div class="flex flex-col lg:flex-row items-center justify-between space-y-3 lg:space-y-0">
          <!-- Copyright and attribution -->
          <div class="text-gray-500 dark:text-gray-400 font-rajdhani text-center lg:text-start">
            <div class="space-y-1">
              <div>
                {@html $LL.footerAuthoredBy({
                  authorOpen: `<a href="http://stevenweathers.com" class="${footerLinkClasses}" target="_blank">`,
                  authorClose: `</a>`,
                })}
              </div>
            </div>
          </div>

          <!-- Version info -->
          <div class="flex items-center space-x-2 rtl:space-x-reverse">
            <div class="w-2 h-2 rounded-full bg-green-500 animate-pulse"></div>
            <span class="text-gray-500 dark:text-gray-400 font-mono">
              {$LL.appVersion({ version: AppVersion })}
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</footer>
