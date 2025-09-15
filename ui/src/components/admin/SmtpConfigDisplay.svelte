<script lang="ts">
  import { onMount } from 'svelte';
  import { CheckCircle, XCircle, Eye, EyeOff } from 'lucide-svelte';
  import LL from '../../i18n/i18n-svelte';

  interface Props {
    xfetch: any;
    notifications: any;
  }

  let { xfetch, notifications }: Props = $props();

  let config = $state<any>({});
  let isLoading = $state(true);
  let showSensitive = $state(false);

  onMount(async () => {
    await loadConfig();
  });

  async function loadConfig() {
    isLoading = true;
    try {
      const response = await xfetch('/api/admin/smtp/config');
      if (response.ok) {
        const result = await response.json();
        config = result.data || {};
      } else {
        notifications.danger('Failed to load SMTP configuration');
      }
    } catch (error) {
      notifications.danger('Error loading SMTP configuration');
    } finally {
      isLoading = false;
    }
  }

  function toggleSensitiveDisplay() {
    showSensitive = !showSensitive;
  }

  function getStatusIcon(enabled: boolean) {
    return enabled ? CheckCircle : XCircle;
  }

  function getStatusColor(enabled: boolean) {
    return enabled ? 'text-green-500' : 'text-red-500';
  }

  function formatConfigValue(key: string, value: any) {
    if (typeof value === 'boolean') {
      return value ? 'Yes' : 'No';
    }
    if (key === 'username' && !showSensitive && value) {
      return value; // Already masked by backend
    }
    if (key === 'hasPassword') {
      return value ? 'Configured' : 'Not configured';
    }
    return value?.toString() || 'Not configured';
  }
</script>

<div class="bg-white dark:bg-gray-800 shadow-lg rounded-lg p-6 mb-6">
  <div class="flex items-center justify-between mb-6">
    <h3 class="text-2xl font-bold text-gray-900 dark:text-white">
      {$LL.smtpConfiguration()}
    </h3>

    <div class="flex items-center gap-4">
      <button
        onclick={toggleSensitiveDisplay}
        class="flex items-center gap-2 px-3 py-1 text-sm bg-gray-100 hover:bg-gray-200 dark:bg-gray-700 dark:hover:bg-gray-600 rounded-lg transition-colors"
        title={showSensitive ? 'Hide sensitive info' : 'Show sensitive info'}
      >
        {#if showSensitive}
          <EyeOff class="w-4 h-4" />
        {:else}
          <Eye class="w-4 h-4" />
        {/if}
        {showSensitive ? $LL.hideSensitive() : $LL.showSensitive()}
      </button>

      <button
        onclick={loadConfig}
        class="px-3 py-1 text-sm bg-blue-100 hover:bg-blue-200 dark:bg-blue-900 dark:hover:bg-blue-800 text-blue-800 dark:text-blue-200 rounded-lg transition-colors"
        disabled={isLoading}
      >
        {isLoading ? $LL.refreshing() : $LL.refresh()}
      </button>
    </div>
  </div>

  {#if isLoading}
    <div class="flex items-center justify-center py-8">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
    </div>
  {:else}
    <!-- Status Overview -->
    <div class="mb-6">
      <div class="flex items-center gap-3 p-4 rounded-lg {config.enabled ? 'bg-green-50 border border-green-200 dark:bg-green-900/20 dark:border-green-700' : 'bg-red-50 border border-red-200 dark:bg-red-900/20 dark:border-red-700'}">
        {#if config.enabled !== undefined}
          <svelte:component
            this={getStatusIcon(config.enabled)}
            class="w-6 h-6 {getStatusColor(config.enabled)}"
          />
          <div>
            <p class="font-semibold {config.enabled ? 'text-green-800 dark:text-green-200' : 'text-red-800 dark:text-red-200'}">
              {config.enabled ? $LL.smtpEnabled() : $LL.smtpDisabled()}
            </p>
            <p class="text-sm {config.enabled ? 'text-green-600 dark:text-green-300' : 'text-red-600 dark:text-red-300'}">
              {config.enabled ? $LL.smtpEnabledDesc() : $LL.smtpDisabledDesc()}
            </p>
          </div>
        {/if}
      </div>
    </div>

    <!-- Configuration Details -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <!-- Server Configuration -->
      <div class="space-y-4">
        <h4 class="text-lg font-semibold text-gray-800 dark:text-gray-200 border-b pb-2">
          {$LL.serverConfiguration()}
        </h4>

        <div class="space-y-3">
          <div>
            <label class="block text-sm font-medium text-gray-600 dark:text-gray-400">
              {$LL.host()}
            </label>
            <p class="text-gray-900 dark:text-white font-mono">
              {formatConfigValue('host', config.host)}
            </p>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-600 dark:text-gray-400">
              {$LL.port()}
            </label>
            <p class="text-gray-900 dark:text-white font-mono">
              {formatConfigValue('port', config.port)}
            </p>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-600 dark:text-gray-400">
              {$LL.secure()}
            </label>
            <p class="text-gray-900 dark:text-white">
              {formatConfigValue('secure', config.secure)}
            </p>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-600 dark:text-gray-400">
              {$LL.skipTLSVerify()}
            </label>
            <p class="text-gray-900 dark:text-white">
              {formatConfigValue('skipTLSVerify', config.skipTLSVerify)}
            </p>
          </div>
        </div>
      </div>

      <!-- Authentication & Sender -->
      <div class="space-y-4">
        <h4 class="text-lg font-semibold text-gray-800 dark:text-gray-200 border-b pb-2">
          {$LL.authenticationAndSender()}
        </h4>

        <div class="space-y-3">
          <div>
            <label class="block text-sm font-medium text-gray-600 dark:text-gray-400">
              {$LL.authType()}
            </label>
            <p class="text-gray-900 dark:text-white font-mono">
              {formatConfigValue('authType', config.authType)}
            </p>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-600 dark:text-gray-400">
              {$LL.username()}
            </label>
            <p class="text-gray-900 dark:text-white font-mono">
              {formatConfigValue('username', config.username)}
            </p>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-600 dark:text-gray-400">
              {$LL.password()}
            </label>
            <p class="text-gray-900 dark:text-white">
              {formatConfigValue('hasPassword', config.hasPassword)}
            </p>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-600 dark:text-gray-400">
              {$LL.senderEmail()}
            </label>
            <p class="text-gray-900 dark:text-white font-mono">
              {formatConfigValue('sender', config.sender)}
            </p>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-600 dark:text-gray-400">
              {$LL.senderName()}
            </label>
            <p class="text-gray-900 dark:text-white">
              {formatConfigValue('senderName', config.senderName)}
            </p>
          </div>
        </div>
      </div>
    </div>

    {#if !config.enabled}
      <div class="mt-6 p-4 bg-yellow-50 border border-yellow-200 dark:bg-yellow-900/20 dark:border-yellow-700 rounded-lg">
        <p class="text-yellow-800 dark:text-yellow-200 text-sm">
          {$LL.smtpDisabledWarning()}
        </p>
      </div>
    {/if}
  {/if}
</div>