<script lang="ts">
  import PageLayout from '../../components/PageLayout.svelte';
  import GithubIcon from '../../components/icons/Github.svelte';
  import LL from '../../i18n/i18n-svelte';
  import { AppConfig } from '../../config';
  import { user } from '../../stores';
  import type { ApiClient } from '../../types/apiclient';
  import type { NotificationService } from '../../types/notifications';
  import {
    AlertTriangle,
    Mail,
    Info,
    CheckCircle,
    Send,
    ExternalLink,
    Sparkles,
    Code,
    Lock,
    BugIcon,
  } from 'lucide-svelte';
  import TextInput from '../../components/forms/TextInput.svelte';

  const { AppVersion, RepoURL } = AppConfig;

  interface Props {
    xfetch: ApiClient;
    router: any;
    notifications: NotificationService;
  }

  let { xfetch, router, notifications }: Props = $props();

  let inquiry: string = $state('');
  let fullName: string = $state($user.name || '');
  let email: string = $state($user.email || '');
  let isSubmitting: boolean = $state(false);
  let submitted: boolean = $state(false);

  function submitSupportTicket() {
    if (!$user || !$user.id) {
      notifications.error('You must be logged in to submit a support ticket.');
      return;
    }

    if (!inquiry || !fullName || !email) {
      notifications.error('Please fill in all fields.');
      return;
    }

    isSubmitting = true;

    xfetch(`/api/users/${$user.id}/support-ticket`, {
      body: {
        inquiry,
        fullName,
        email,
      },
    })
      .then(() => {
        notifications.success('Support ticket submitted successfully.');
        inquiry = '';
        submitted = true;
      })
      .catch((error: any) => {
        notifications.error('Failed to submit support ticket.');
        console.error(error);
      })
      .finally(() => {
        isSubmitting = false;
      });
  }
</script>

<svelte:head>
  <title>{$LL.support()} | {$LL.appName()}</title>
</svelte:head>

<PageLayout>
  <div class="max-w-8xl mx-auto px-4 py-12 text-gray-900 dark:text-gray-100 relative">
    <!-- Header Section -->
    <div class="text-center mb-16">
      <h1
        class="text-5xl lg:text-7xl font-bold font-rajdhani dark:text-white uppercase tracking-wider mb-4 bg-gradient-to-r from-blue-600 via-purple-600 to-emerald-600 bg-clip-text text-transparent"
      >
        {$LL.support()}
      </h1>

      <p class="text-xl text-gray-600 dark:text-gray-300 max-w-2xl mx-auto leading-relaxed">
        We're here to help you succeed. Get support for bugs, questions, or feedback.
      </p>

      <div class="mt-8 flex justify-center">
        <div class="h-1 w-32 bg-gradient-to-r from-blue-500 via-purple-500 to-emerald-500 rounded-full shadow-lg"></div>
      </div>
    </div>

    <div class="grid lg:grid-cols-2 gap-8">
      <!-- Bug Report Section -->
      <div class="relative">
        <div
          class="bg-white/80 dark:bg-gray-800/80 backdrop-blur-xl rounded-3xl shadow-2xl border border-white/20 dark:border-gray-700/50 p-8"
        >
          <div class="flex items-start space-x-4 mb-6">
            <div
              class="w-14 h-14 bg-gradient-to-br from-orange-500 to-red-500 rounded-2xl flex items-center justify-center shadow-lg"
            >
              <BugIcon class="w-8 h-8 text-white" />
            </div>
            <div>
              <h2 class="text-2xl lg:text-3xl font-bold font-rajdhani dark:text-white mb-2">Found a Bug?</h2>
              <p class="text-gray-600 dark:text-gray-300">Help us improve by reporting issues</p>
            </div>
          </div>

          <div class="space-y-6">
            <p class="text-gray-700 dark:text-gray-300 leading-relaxed">
              Create an issue on our GitHub repository and help us make the platform better for everyone.
            </p>

            <a
              href="{RepoURL}/issues"
              target="_blank"
              class="group relative inline-flex items-center gap-3 px-7 py-4 bg-white dark:bg-gray-900 border-2 border-gray-200 dark:border-gray-700 hover:border-gray-400 dark:hover:border-gray-500 text-gray-900 dark:text-white font-semibold rounded-xl shadow-lg hover:shadow-xl transform hover:scale-[1.02] transition-all duration-300 ease-out focus:outline-none focus:ring-4 focus:ring-gray-300/50 dark:focus:ring-gray-600/50 overflow-hidden"
            >
              <!-- Content with enhanced interactions -->
              <div class="relative flex items-center gap-3">
                <div
                  class="p-1 rounded-lg bg-gray-100 dark:bg-gray-800 group-hover:bg-gray-200 dark:group-hover:bg-gray-700 transition-colors duration-300"
                >
                  <GithubIcon width={24} class="transform group-hover:scale-110 transition-transform duration-300" />
                </div>
                <span class="group-hover:text-blue-600 dark:group-hover:text-blue-400 transition-colors duration-300"
                  >Open GitHub Issues</span
                >
                <ExternalLink
                  class="w-6 h-6 text-gray-500 group-hover:text-blue-500 transform group-hover:translate-x-1 group-hover:-translate-y-1 transition-all duration-300"
                />
              </div>
            </a>

            <div
              class="flex items-center justify-between p-4 bg-emerald-50 dark:bg-emerald-900/30 rounded-xl border border-emerald-200 dark:border-emerald-700"
            >
              <div class="flex items-center gap-3">
                <div class="w-8 h-8 bg-emerald-500 rounded-lg flex items-center justify-center">
                  <Code class="w-6 h-6 text-white" />
                </div>
                <span class="font-medium text-emerald-800 dark:text-emerald-200">Mention the Current Version</span>
              </div>
              <span class="px-4 py-2 bg-emerald-500 text-white font-bold rounded-lg shadow-md">
                {AppVersion}
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- Contact Form Section -->
      <div class="group relative">
        <div
          class="bg-white/80 dark:bg-gray-800/80 backdrop-blur-xl rounded-3xl shadow-2xl border border-white/20 dark:border-gray-700/50 p-8"
        >
          <div class="flex items-start space-x-4 mb-8">
            <div
              class="w-14 h-14 bg-gradient-to-br from-blue-500 to-emerald-500 rounded-2xl flex items-center justify-center shadow-lg"
            >
              <Mail class="w-8 h-8 text-white" />
            </div>
            <div>
              <h2 class="text-2xl lg:text-3xl font-bold font-rajdhani dark:text-white uppercase mb-2">Contact Us</h2>
              <p class="text-gray-600 dark:text-gray-300">Get personalized support from our team</p>
            </div>
          </div>

          {#if !$user.name}
            <div class="text-center py-16">
              <div
                class="w-20 h-20 bg-gradient-to-br from-amber-400 to-orange-500 rounded-full flex items-center justify-center mx-auto shadow-xl mb-6"
              >
                <Lock class="w-10 h-10 text-white" />
              </div>
              <h3 class="text-2xl font-bold text-gray-900 dark:text-white mb-3">Login Required</h3>
              <p class="text-gray-600 dark:text-gray-400 text-lg">
                Please log in to submit a support ticket and get personalized assistance.
              </p>
            </div>
          {:else if submitted}
            <div class="text-center py-16">
              <div
                class="w-20 h-20 bg-gradient-to-br from-emerald-400 to-green-500 rounded-full flex items-center justify-center mx-auto shadow-xl mb-6"
              >
                <CheckCircle class="w-10 h-10 text-white" />
              </div>
              <h3 class="text-2xl font-bold text-emerald-800 dark:text-emerald-200 mb-3">Success!</h3>
              <p class="text-emerald-600 dark:text-emerald-400 text-lg">
                Thank you for your submission! We'll get back to you within 24 hours.
              </p>
            </div>
          {:else}
            <form
              onsubmit={(event: Event) => {
                event.preventDefault();
                submitSupportTicket();
              }}
              class="space-y-6"
            >
              <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div class="space-y-2">
                  <label for="fullName" class="block font-bold text-gray-700 dark:text-gray-300 mb-4">
                    Full Name <span class="text-red-500">*</span>
                  </label>
                  <TextInput
                    type="text"
                    id="fullName"
                    bind:value={fullName}
                    required
                    placeholder="Enter your full name"
                  />
                </div>

                <div class="space-y-2">
                  <label for="email" class="block font-bold text-gray-700 dark:text-gray-300 mb-4">
                    Email Address <span class="text-red-500">*</span>
                  </label>
                  <TextInput
                    type="email"
                    id="email"
                    bind:value={email}
                    required
                    placeholder="Enter your email address"
                  />
                </div>
              </div>

              <div class="space-y-2">
                <label for="inquiry" class="block font-bold text-gray-700 dark:text-gray-300 mb-2">
                  Your Inquiry <span class="text-red-500">*</span>
                </label>
                <textarea
                  id="inquiry"
                  bind:value={inquiry}
                  required
                  rows="6"
                  class="w-full px-5 py-4 border-2 border-gray-200 dark:border-gray-600 rounded-xl shadow-lg focus:ring-4 focus:ring-blue-500/25 focus:border-blue-500 dark:focus:ring-blue-400/25 dark:focus:border-blue-400 bg-white/90 dark:bg-gray-700/90 text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400 resize-none backdrop-blur-sm"
                  placeholder="Please describe your inquiry in detail... We're here to help!"
                ></textarea>
              </div>

              <div class="pt-6 flex justify-end">
                <button
                  type="submit"
                  disabled={isSubmitting}
                  class="group relative w-full md:w-auto inline-flex items-center justify-center px-8 py-4 bg-gradient-to-r from-blue-600 via-purple-600 to-emerald-600 hover:from-blue-500 hover:via-purple-500 hover:to-emerald-500 disabled:from-gray-400 disabled:via-gray-500 disabled:to-gray-600 text-white font-bold rounded-xl shadow-xl hover:shadow-2xl disabled:shadow-lg transform hover:scale-[1.02] disabled:scale-100 focus:outline-none focus:ring-4 focus:ring-purple-500/50 transition-all duration-300 ease-out overflow-hidden disabled:cursor-not-allowed"
                >
                  <!-- Loading spinner background -->
                  {#if isSubmitting}
                    <div
                      class="absolute inset-0 bg-gradient-to-r from-gray-400 via-gray-500 to-gray-600 flex items-center justify-center"
                    >
                      <div class="w-6 h-6 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
                    </div>
                  {/if}

                  <!-- Dynamic background effects -->
                  <div
                    class="absolute inset-0 bg-gradient-to-r from-blue-400/30 via-purple-400/30 to-emerald-400/30 opacity-0 group-hover:opacity-100 transition-opacity duration-300"
                  ></div>

                  <!-- Content -->
                  <div
                    class="relative flex items-center {isSubmitting
                      ? 'opacity-0'
                      : 'opacity-100'} transition-opacity duration-300"
                  >
                    <Send class="w-5 h-5 mr-3 transform group-hover:translate-x-1 transition-transform duration-300" />
                    <span class="tracking-wide">
                      {isSubmitting ? 'Submitting...' : 'Submit Support Ticket'}
                    </span>
                  </div>

                  <!-- Success animation overlay -->
                  <div
                    class="absolute inset-0 bg-gradient-to-r from-emerald-500 via-green-500 to-emerald-600 opacity-0 flex items-center justify-center transform scale-0 transition-all duration-500"
                  >
                    <CheckCircle class="w-6 h-6 text-white" />
                  </div>

                  <!-- Enhanced shine effect -->
                  <div
                    class="absolute inset-0 -top-2 -bottom-2 bg-gradient-to-r from-transparent via-white/25 to-transparent skew-x-12 -translate-x-full group-hover:translate-x-full transition-transform duration-1000 ease-out"
                  ></div>
                </button>
              </div>
            </form>
          {/if}
        </div>
      </div>
    </div>

    <!-- Decorative Elements -->
    <div class="mt-20 text-center">
      <div
        class="inline-flex items-center gap-2 px-6 py-3 bg-white/60 dark:bg-gray-800/60 backdrop-blur-xl rounded-full border border-white/20 dark:border-gray-700/50 shadow-lg"
      >
        <div class="w-2 h-2 bg-green-500 rounded-full"></div>
        <span class="font-medium text-gray-600 dark:text-gray-400"
          >We typically respond within 24 hours, thank you for your patience!</span
        >
      </div>
    </div>
  </div>
</PageLayout>
