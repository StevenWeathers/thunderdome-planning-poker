<script lang="ts">
  import SolidButton from '../global/SolidButton.svelte';
  import TextInput from '../forms/TextInput.svelte';
  import LL from '../../i18n/i18n-svelte';
  import { validateEmail } from '../../validationUtils';

  interface Props {
    xfetch: any;
    notifications: any;
  }

  let { xfetch, notifications }: Props = $props();

  let testEmail = $state('');
  let testName = $state('');
  let isTestingConnection = $state(false);
  let isSendingTestEmail = $state(false);
  let connectionTestResult = $state<string | null>(null);
  let emailTestResult = $state<string | null>(null);

  // Form validation
  let emailError = $derived(testEmail && !validateEmail(testEmail) ? 'Invalid email address' : '');
  let nameError = $derived(testName.trim() === '' ? 'Name is required' : '');
  let canSendTestEmail = $derived(testEmail && testName.trim() && !emailError && !nameError);

  async function testConnection() {
    isTestingConnection = true;
    connectionTestResult = null;

    try {
      const response = await xfetch('/api/admin/smtp/test-connection', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
      });

      const result = await response.json();

      if (!response.ok) {
        connectionTestResult = `Connection failed: ${result.error || 'Unknown error'}`;
        notifications.danger('SMTP connection test failed');
      } else {
        connectionTestResult = 'SMTP connection test successful!';
        notifications.success('SMTP connection test successful');
      }
    } catch (error) {
      connectionTestResult = `Connection failed: ${error.message}`;
      notifications.danger('SMTP connection test failed');
    } finally {
      isTestingConnection = false;
    }
  }

  async function sendTestEmail() {
    if (!canSendTestEmail) return;

    isSendingTestEmail = true;
    emailTestResult = null;

    try {
      const response = await xfetch('/api/admin/smtp/test-email', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          email: testEmail.trim(),
          name: testName.trim(),
        }),
      });

      const result = await response.json();

      if (!response.ok) {
        const errorMsg = result.error || result.message || 'Unknown error';
        emailTestResult = `Email sending failed: ${errorMsg}`;
        notifications.danger(`Test email sending failed: ${errorMsg}`);
        console.error('SMTP test email error:', {
          status: response.status,
          statusText: response.statusText,
          result: result,
          requestData: { email: testEmail.trim(), name: testName.trim() }
        });
      } else {
        emailTestResult = `Test email sent successfully to ${testEmail}!`;
        notifications.success('Test email sent successfully');
        // Clear form on success
        testEmail = '';
        testName = '';
      }
    } catch (error) {
      emailTestResult = `Email sending failed: ${error.message}`;
      notifications.danger('Test email sending failed');
      console.error('SMTP test email exception:', error);
    } finally {
      isSendingTestEmail = false;
    }
  }

  function clearResults() {
    connectionTestResult = null;
    emailTestResult = null;
  }
</script>

<div class="bg-white dark:bg-gray-800 shadow-lg rounded-lg p-6 mb-6">
  <h3 class="text-2xl font-bold text-gray-900 dark:text-white mb-6">
    {$LL.smtpTestTitle()}
  </h3>

  <!-- Connection Test Section -->
  <div class="mb-8">
    <h4 class="text-lg font-semibold text-gray-800 dark:text-gray-200 mb-4">
      {$LL.smtpConnectionTestTitle()}
    </h4>
    <p class="text-gray-600 dark:text-gray-400 mb-4">
      {$LL.smtpConnectionTestDescription()}
    </p>

    <div class="flex items-center gap-4 mb-4">
      <SolidButton
        onClick={testConnection}
        disabled={isTestingConnection}
        color="blue"
        testid="test-smtp-connection"
      >
        {isTestingConnection ? $LL.testing() : $LL.testConnection()}
      </SolidButton>

      {#if connectionTestResult}
        <button
          onclick={clearResults}
          class="text-sm text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200"
        >
          {$LL.clearResults()}
        </button>
      {/if}
    </div>

    {#if connectionTestResult}
      <div
        class="p-4 rounded-lg {connectionTestResult.includes('successful')
          ? 'bg-green-100 border border-green-300 text-green-800 dark:bg-green-900 dark:border-green-700 dark:text-green-200'
          : 'bg-red-100 border border-red-300 text-red-800 dark:bg-red-900 dark:border-red-700 dark:text-red-200'}"
      >
        {connectionTestResult}
      </div>
    {/if}
  </div>

  <!-- Test Email Section -->
  <div class="border-t pt-8">
    <h4 class="text-lg font-semibold text-gray-800 dark:text-gray-200 mb-4">
      {$LL.smtpSendTestEmailTitle()}
    </h4>
    <p class="text-gray-600 dark:text-gray-400 mb-6">
      {$LL.smtpSendTestEmailDescription()}
    </p>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
      <div>
        <TextInput
          bind:value={testEmail}
          placeholder={$LL.emailAddress()}
          id="test-email"
          name="test-email"
          type="email"
          required={true}
          testid="smtp-test-email-input"
        />
        {#if emailError}
          <p class="text-red-500 text-sm mt-1">{emailError}</p>
        {/if}
      </div>

      <div>
        <TextInput
          bind:value={testName}
          placeholder={$LL.recipientName()}
          id="test-name"
          name="test-name"
          type="text"
          required={true}
          testid="smtp-test-name-input"
        />
        {#if nameError}
          <p class="text-red-500 text-sm mt-1">{nameError}</p>
        {/if}
      </div>
    </div>

    <div class="flex items-center gap-4 mb-4">
      <SolidButton
        onClick={sendTestEmail}
        disabled={!canSendTestEmail || isSendingTestEmail}
        color="green"
        testid="send-test-email"
      >
        {isSendingTestEmail ? $LL.sending() : $LL.sendTestEmail()}
      </SolidButton>

      {#if emailTestResult}
        <button
          onclick={clearResults}
          class="text-sm text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200"
        >
          {$LL.clearResults()}
        </button>
      {/if}
    </div>

    {#if emailTestResult}
      <div
        class="p-4 rounded-lg {emailTestResult.includes('successful')
          ? 'bg-green-100 border border-green-300 text-green-800 dark:bg-green-900 dark:border-green-700 dark:text-green-200'
          : 'bg-red-100 border border-red-300 text-red-800 dark:bg-red-900 dark:border-red-700 dark:text-red-200'}"
      >
        {emailTestResult}
      </div>
    {/if}
  </div>
</div>