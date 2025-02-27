<script lang="ts">
  import { validateName, validatePasswords } from '../../validationUtils';
  import LL from '../../i18n/i18n-svelte';
  import { Mail, User } from 'lucide-svelte';
  import PasswordInput from '../forms/PasswordInput.svelte';
  import { AppConfig } from '../../config';
  import Checkbox from '../forms/Checkbox.svelte';
  import { user } from '../../stores';
  import TextInput from '../forms/TextInput.svelte';

  export let notifications;
  export let wasInvited = false;
  export let email = '';
  export let userName = '';
  export let fullOnly = false;
  export let handleGuestRegistration;
  export let handleFullAccountRegistration;
  export let isAdmin = false;

  let password1 = '';
  let password2 = '';

  /** @type {TextInput} */
  let userNameTextInput;

  function onSubmit(e) {
    e.preventDefault();

    const validName = validateName(userName);
    let noFormErrors = true;

    if (!validName.valid) {
      noFormErrors = false;
      notifications.danger(validName.error, 1500);
    }

    if (createFullAccount) {
      const validPasswords = validatePasswords(password1, password2);
      if (!validPasswords.valid) {
        noFormErrors = false;
        notifications.danger(validPasswords.error, 1500);
      }
    }

    if (noFormErrors) {
      if (!createFullAccount) {
        handleGuestRegistration(userName);
      } else {
        handleFullAccountRegistration(userName, email, password1, password2);
      }
    }
  }

  $: isGuest = $user.id !== '' && $user.rank === 'GUEST';
  $: createFullAccount =
    fullOnly ||
    isGuest ||
    (!AppConfig.AllowGuests && AppConfig.AllowRegistration);
  $: passwordsMatch = password1 === password2;

  $: submitDisabled =
    userName === '' ||
    (createFullAccount &&
      (email === '' ||
        password1 === '' ||
        password2 === '' ||
        !passwordsMatch));

  // Focus the warrior name input field if it exists
  userNameTextInput?.focus();
</script>

{#if !AppConfig.AllowRegistration && !AppConfig.AllowGuests && !isAdmin}
  <div class="text-red-500">
    {$LL.registrationDisabled()}
  </div>
{:else}
  <form on:submit="{onSubmit}" name="register" class="space-y-6">
    <div class="space-y-2">
      <TextInput
        bind:this="{userNameTextInput}"
        bind:value="{userName}"
        placeholder="{$LL.userNamePlaceholder()}"
        id="yourName"
        name="yourName"
        required
        icon="{User}"
      />
    </div>

    {#if AppConfig.AllowRegistration && !isGuest && !fullOnly}
      <div class="flex items-center">
        <Checkbox
          id="createFullAccount"
          name="createFullAccount"
          bind:checked="{createFullAccount}"
          label="{AppConfig.GuestsAllowed
            ? 'Create full account (optional)'
            : 'Create full account'}"
        />
      </div>
    {/if}

    {#if createFullAccount}
      <div class="space-y-2">
        <TextInput
          bind:value="{email}"
          placeholder="{$LL.enterYourEmail()}"
          id="yourEmail"
          name="yourEmail"
          type="email"
          required
          disabled="{wasInvited}"
          icon="{Mail}"
          autocomplete="email"
        />
      </div>

      <div class="space-y-2">
        <PasswordInput
          bind:value="{password1}"
          placeholder="{$LL.passwordPlaceholder()}"
          id="yourPassword1"
          name="yourPassword1"
          data-testid="yourPassword1"
          required
        />
      </div>

      <div class="space-y-2">
        <PasswordInput
          bind:value="{password2}"
          placeholder="{$LL.confirmPasswordPlaceholder()}"
          id="yourPassword2"
          name="yourPassword2"
          data-testid="yourPassword2"
          required
        />
      </div>

      {#if password1 !== '' && password2 !== '' && !passwordsMatch}
        <div class="text-center text-lg text-red-500 font-medium">
          Passwords do not match
        </div>
      {/if}
    {/if}

    <div class="pt-4">
      <button
        type="submit"
        disabled="{submitDisabled}"
        class="w-full group relative flex justify-center py-3 px-4 border border-transparent text-lg font-medium rounded-lg text-white transition-all duration-300 transform hover:scale-105 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-purple-500 bg-gradient-to-r from-purple-500 to-indigo-500 hover:from-purple-600 hover:to-indigo-600 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        <span class="absolute left-0 inset-y-0 flex items-center pl-3">
          <User
            class="h-5 w-5 text-purple-300 group-hover:text-purple-200"
            aria-hidden="true"
          />
        </span>
        {createFullAccount ? $LL.createAccount() : $LL.registerAsGuest()}
      </button>
    </div>
  </form>
{/if}
