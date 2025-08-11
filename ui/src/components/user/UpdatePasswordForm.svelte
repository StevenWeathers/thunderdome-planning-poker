<script lang="ts">
  import LL from '../../i18n/i18n-svelte';
  import { AppConfig } from '../../config';
  import { validatePasswords } from '../../validationUtils';
  import SolidButton from '../global/SolidButton.svelte';
  import PasswordInput from '../forms/PasswordInput.svelte';

  import type { NotificationService } from '../../types/notifications';

  interface Props {
    handleUpdate?: any;
    toggleForm?: any;
    notifications: NotificationService;
  }

  let { handleUpdate = () => {}, toggleForm = () => {}, notifications }: Props = $props();

  const { LdapEnabled, HeaderAuthEnabled } = AppConfig;

  let password1 = $state('');
  let password2 = $state('');

  function updatePassword(e) {
    e.preventDefault();

    const validPasswords = validatePasswords(password1, password2);

    let noFormErrors = true;

    if (!validPasswords.valid) {
      noFormErrors = false;
      notifications.danger(validPasswords.error, 1500);
    }

    if (noFormErrors) {
      handleUpdate(password1, password2);
    }
  }

  let updatePasswordDisabled =
    $derived(password1 === '' || password2 === '' || LdapEnabled || HeaderAuthEnabled);
</script>

<form onsubmit={updatePassword} name="updatePassword">
  <div class="mb-4">
    <label
      class="block text-gray-700 dark:text-gray-400 font-bold mb-2"
      for="yourPassword1"
    >
      {$LL.password()}
    </label>
    <PasswordInput
      bind:value="{password1}"
      placeholder={$LL.passwordPlaceholder()}
      id="yourPassword1"
      name="yourPassword1"
      required
    />
  </div>

  <div class="mb-4">
    <label
      class="block text-gray-700 dark:text-gray-400 font-bold mb-2"
      for="yourPassword2"
    >
      {$LL.confirmPassword()}
    </label>
    <PasswordInput
      bind:value="{password2}"
      placeholder={$LL.confirmPasswordPlaceholder()}
      id="yourPassword2"
      name="yourPassword2"
      required
    />
  </div>

  <div class="text-right">
    <button
      type="button"
      class="inline-block align-baseline font-bold text-sm
            text-blue-500 hover:text-blue-800 me-4"
      onclick={toggleForm}
    >
      {$LL.cancel()}
    </button>
    <SolidButton type="submit" disabled={updatePasswordDisabled}>
      {$LL.update()}
    </SolidButton>
  </div>
</form>
