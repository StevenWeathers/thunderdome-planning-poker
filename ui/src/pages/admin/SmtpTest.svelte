<script lang="ts">
  import { onMount } from 'svelte';
  import AdminPageLayout from '../../components/admin/AdminPageLayout.svelte';
  import SmtpConfigDisplay from '../../components/admin/SmtpConfigDisplay.svelte';
  import SmtpTestForm from '../../components/admin/SmtpTestForm.svelte';
  import LL from '../../i18n/i18n-svelte';
  import { user } from '../../stores';
  import { validateUserIsAdmin } from '../../validationUtils';

  import type { NotificationService } from '../../types/notifications';
  import type { ApiClient } from '../../types/apiclient';

  interface Props {
    xfetch: ApiClient;
    router: any;
    notifications: NotificationService;
  }

  let { xfetch, router, notifications }: Props = $props();

  onMount(() => {
    if (!validateUserIsAdmin($user)) {
      router.route('/');
      return;
    }
  });
</script>

<svelte:head>
  <title>{$LL.smtpTestTitle()} | {$LL.appName()}</title>
</svelte:head>

<AdminPageLayout activePage="smtp-test">
  <div class="mb-6">
    <h1 class="text-3xl font-bold text-gray-900 dark:text-white">
      {$LL.smtpTestTitle()}
    </h1>
    <p class="mt-2 text-gray-600 dark:text-gray-400">
      {$LL.smtpTestDescription()}
    </p>
  </div>

  <!-- Configuration Display -->
  <SmtpConfigDisplay {xfetch} {notifications} />

  <!-- Test Form -->
  <SmtpTestForm {xfetch} {notifications} />

  <!-- Documentation Section -->
  <div class="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-700 rounded-lg p-6">
    <h3 class="text-lg font-semibold text-blue-900 dark:text-blue-200 mb-3">
      {$LL.smtpDocumentationTitle()}
    </h3>
    <div class="text-blue-800 dark:text-blue-300 space-y-2">
      <p>{$LL.smtpDocumentation1()}</p>
      <p>{$LL.smtpDocumentation2()}</p>
      <p>{$LL.smtpDocumentation3()}</p>
    </div>
  </div>
</AdminPageLayout>