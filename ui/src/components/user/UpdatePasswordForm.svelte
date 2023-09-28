<script lang="ts">
  import LL from '../../i18n/i18n-svelte';
  import { AppConfig } from '../../config';
  import { validatePasswords } from '../../validationUtils';
  import SolidButton from '../SolidButton.svelte';
  import TextInput from '../TextInput.svelte';

  export let handleUpdate = () => {};
  export let toggleForm = () => {};
  export let notifications;

  const { LdapEnabled, HeaderAuthEnabled } = AppConfig;

  let warriorPassword1 = '';
  let warriorPassword2 = '';

  function updateWarriorPassword(e) {
    e.preventDefault();

    const validPasswords = validatePasswords(
      warriorPassword1,
      warriorPassword2,
    );

    let noFormErrors = true;

    if (!validPasswords.valid) {
      noFormErrors = false;
      notifications.danger(validPasswords.error, 1500);
    }

    if (noFormErrors) {
      handleUpdate(warriorPassword1, warriorPassword2);
    }
  }

  $: updatePasswordDisabled =
    warriorPassword1 === '' ||
    warriorPassword2 === '' ||
    LdapEnabled ||
    HeaderAuthEnabled;
</script>

<form on:submit="{updateWarriorPassword}" name="updateWarriorPassword">
  <div class="mb-4">
    <label
      class="block text-gray-700 dark:text-gray-400 font-bold mb-2"
      for="yourPassword1"
    >
      {$LL.password()}
    </label>
    <TextInput
      bind:value="{warriorPassword1}"
      placeholder="{$LL.passwordPlaceholder()}"
      id="yourPassword1"
      name="yourPassword1"
      type="password"
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
    <TextInput
      bind:value="{warriorPassword2}"
      placeholder="{$LL.confirmPasswordPlaceholder()}"
      id="yourPassword2"
      name="yourPassword2"
      type="password"
      required
    />
  </div>

  <div class="text-right">
    <button
      type="button"
      class="inline-block align-baseline font-bold text-sm
            text-blue-500 hover:text-blue-800 me-4"
      on:click="{toggleForm}"
    >
      {$LL.cancel()}
    </button>
    <SolidButton type="submit" disabled="{updatePasswordDisabled}">
      {$LL.update()}
    </SolidButton>
  </div>
</form>
