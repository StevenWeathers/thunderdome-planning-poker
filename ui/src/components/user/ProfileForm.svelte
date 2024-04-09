<script lang="ts">
  import SolidButton from '../global/SolidButton.svelte';
  import VerifiedIcon from '../icons/VerifiedIcon.svelte';
  import HollowButton from '../global/HollowButton.svelte';
  import { countryList } from '../../country';
  import { AppConfig } from '../../config';
  import { validateName, validateUserIsAdmin } from '../../validationUtils';
  import LL, { locale, setLocale } from '../../i18n/i18n-svelte';
  import UserAvatar from './UserAvatar.svelte';
  import SetupMFA from './SetupMFA.svelte';
  import DeleteConfirmation from '../global/DeleteConfirmation.svelte';
  import { user } from '../../stores';
  import LocaleSwitcher from '../forms/LocaleInput.svelte';
  import type { Locales } from '../../i18n/i18n-types';
  import { loadLocaleAsync } from '../../i18n/i18n-util.async';
  import TextInput from '../forms/TextInput.svelte';
  import SelectInput from '../forms/SelectInput.svelte';

  const setupI18n = async (locale: Locales) => {
    await loadLocaleAsync(locale);
    setLocale(locale);
  };

  export let profile = {
    id: '',
    rank: '',
    name: '',
    email: '',
    company: '',
    country: '',
    jobTitle: '',
    notificationsEnabled: true,
    avatar: '',
    gravatarHash: '',
    verified: false,
    mfaEnabled: false,
    theme: 'auto',
  };
  export let handleUpdate = () => {};
  export let toggleUpdatePassword;
  export let eventTag;
  export let notifications;
  export let xfetch;
  export let ldapEnabled;
  export let headerAuthEnabled;

  const { AvatarService } = AppConfig;

  const themes = ['auto', 'light', 'dark'];
  const configurableAvatarServices = ['gravatar', 'robohash', 'govatar'];
  const isAvatarConfigurable =
    configurableAvatarServices.includes(AvatarService);
  const avatarOptions = {
    gravatar: ['mp', 'identicon', 'monsterid', 'wavatar', 'retro', 'robohash'],
    robohash: ['set1', 'set2', 'set3', 'set4'],
    govatar: ['male', 'female'],
  };

  let avatars = isAvatarConfigurable ? avatarOptions[AvatarService] : [];

  function handleSubmit(e) {
    e.preventDefault();
    const validName = validateName(profile.name);
    let p = {
      name: profile.name,
      country: profile.country,
      company: profile.company,
      jobTitle: profile.jobTitle,
      notificationsEnabled: profile.notificationsEnabled,
      avatar: profile.avatar,
      locale: $locale,
      email: profile.email,
      theme: profile.theme,
    };

    if (!validName.valid) {
      notifications.danger(validName.error, 1500);
    } else {
      handleUpdate(p);
    }
  }

  let showMFASetup = false;

  function toggleMfaSetup() {
    showMFASetup = !showMFASetup;
  }

  function handleMfaSetupCompletion() {
    profile.mfaEnabled = true;
    toggleMfaSetup();
  }

  let showMfaRemove = false;

  function toggleMfaRemove() {
    showMfaRemove = !showMfaRemove;
  }

  function handleMfaRemove() {
    xfetch('/api/auth/mfa', { method: 'DELETE' })
      .then(() => {
        profile.mfaEnabled = false;
        toggleMfaRemove();
        notifications.success($LL.mfa2faRemoveSuccess());
      })
      .catch(() => {
        notifications.danger($LL.mfa2faRemoveFailure());
      });
  }

  function requestVerifyEmail(e) {
    e.preventDefault();
    xfetch(`/api/users/${profile.id}/request-verify`, { method: 'POST' })
      .then(function () {
        eventTag('user_verify_request', 'engagement', 'success');

        notifications.success($LL.requestVerifyEmailSuccess());
      })
      .catch(function () {
        notifications.danger($LL.requestVerifyEmailFailure());
        eventTag('user_verify_request', 'engagement', 'failure');
      });
  }

  $: updateDisabled = profile.name === '';
  $: userIsAdmin = validateUserIsAdmin($user);
</script>

<form on:submit="{handleSubmit}" name="updateProfile">
  <div class="mb-4">
    <label
      class="block text-gray-700 dark:text-gray-400 font-bold mb-2"
      for="yourName"
    >
      {$LL.name()}
    </label>
    <input
      bind:value="{profile.name}"
      placeholder="{$LL.yourNamePlaceholder()}"
      class="bg-gray-100 dark:bg-gray-900 border-gray-200 dark:border-gray-800 border-2 appearance-none
            rounded w-full py-2 px-3 text-gray-700 dark:text-gray-300 leading-tight
            focus:outline-none focus:bg-white dark:focus:bg-gray-700 focus:border-indigo-500 focus:caret-indigo-500
            dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
      id="yourName"
      name="yourName"
      type="text"
      disabled="{ldapEnabled || headerAuthEnabled}"
      required
    />
  </div>

  <div class="mb-4">
    <label
      class="block text-gray-700 dark:text-gray-400 font-bold mb-2"
      for="yourEmail"
    >
      {$LL.email()}
      {#if profile.verified}
        <span
          class="inline-block font-bold text-green-600
                                    border-green-500 border py-1 px-2 rounded
                                    ms-1"
          data-testid="user-verified"
        >
          {$LL.verified()}
          <VerifiedIcon class="inline fill-current h-4 w-4" />
        </span>
      {:else if profile.rank !== 'GUEST'}
        <button
          class="float-right inline-block align-baseline font-bold text-sm text-blue-500
                                        hover:text-blue-800"
          on:click="{requestVerifyEmail}"
          data-testid="request-verify"
          type="button"
          >{$LL.requestVerifyEmail()}
        </button>
      {/if}
    </label>
    <TextInput
      bind:value="{profile.email}"
      additionalClasses="{!userIsAdmin ? 'cursor-not-allowed' : ''}"
      id="yourEmail"
      name="yourEmail"
      type="email"
      disabled="{!userIsAdmin}"
    />
  </div>

  {#if profile.rank !== 'GUEST'}
    <div class="mb-4">
      <p class="block text-gray-700 dark:text-gray-400 font-bold mb-2">
        {$LL.mfa2faLabel()}
      </p>
      {#if !profile.mfaEnabled}
        <HollowButton color="teal" onClick="{toggleMfaSetup}"
          >{$LL.mfa2faSetup()}
        </HollowButton>
      {:else}
        <HollowButton color="red" onClick="{toggleMfaRemove}"
          >{$LL.mfa2faRemove()}
        </HollowButton>
      {/if}
    </div>
  {/if}

  <div class="mb-4">
    <label
      class="block text-gray-700 dark:text-gray-400 font-bold mb-2"
      for="yourCountry"
    >
      {$LL.country()}
    </label>

    <SelectInput
      bind:value="{profile.country}"
      id="yourCountry"
      name="yourCountry"
    >
      <option value="">
        {$LL.chooseCountryPlaceholder()}
      </option>
      {#each countryList as item}
        <option value="{item.abbrev}">
          {item.name} [{item.abbrev}]
        </option>
      {/each}
    </SelectInput>
  </div>

  <div class="mb-4">
    <div class="text-gray-700 dark:text-gray-400 font-bold mb-2">
      {$LL.locale()}
    </div>
    <LocaleSwitcher
      selectedLocale="{$locale}"
      on:locale-changed="{e => setupI18n(e.detail)}"
    />
  </div>

  <div class="mb-4">
    <label
      class="block text-gray-700 dark:text-gray-400 font-bold mb-2"
      for="yourCompany"
    >
      {$LL.company()}
    </label>
    <TextInput
      bind:value="{profile.company}"
      placeholder="{$LL.companyPlaceholder()}"
      id="yourCompany"
      name="yourCompany"
    />
  </div>

  <div class="mb-4">
    <label
      class="block text-gray-700 dark:text-gray-400 font-bold mb-2"
      for="yourJobTitle"
    >
      {$LL.jobTitle()}
    </label>
    <TextInput
      bind:value="{profile.jobTitle}"
      placeholder="{$LL.jobTitlePlaceholder()}"
      id="yourJobTitle"
      name="yourJobTitle"
    />
  </div>

  <div class="mb-4">
    <label
      class="block text-gray-700 dark:text-gray-400 font-bold mb-2"
      for="yourCountry"
    >
      {$LL.theme()}
    </label>
    <SelectInput bind:value="{profile.theme}" id="theme" name="theme">
      {#each themes as theme}
        <option value="{theme}">
          {theme}
        </option>
      {/each}
    </SelectInput>
  </div>

  <div class="mb-4">
    <label class="block text-gray-700 dark:text-gray-400 font-bold mb-2">
      <input
        bind:checked="{profile.notificationsEnabled}"
        type="checkbox"
        class="w-4 h-4 dark:accent-lime-400 me-1"
      />
      <span>
        {$LL.enableBattleNotifications({
          friendly: AppConfig.FriendlyUIVerbs,
        })}
      </span>
    </label>
  </div>

  {#if isAvatarConfigurable}
    <div class="mb-4">
      <label
        class="block text-gray-700 dark:text-gray-400 font-bold
                                mb-2"
        for="yourAvatar"
      >
        {$LL.avatar()}
      </label>
      <div class="flex items-center content-center">
        <div class="grow">
          {#if AvatarService === 'gravatar'}
            <span class="dark:text-gray-300">Optional Gravatar Fallback</span>
          {/if}
          <SelectInput
            bind:value="{profile.avatar}"
            id="yourAvatar"
            name="yourAvatar"
          >
            {#each avatars as item}
              <option value="{item}">
                {item}
              </option>
            {/each}
          </SelectInput>
        </div>
        <div class="shrink ms-4">
          <UserAvatar
            warriorId="{profile.id}"
            avatar="{profile.avatar}"
            gravatarHash="{profile.gravatarHash}"
            userName="{profile.name}"
            width="48"
            class="rounded-full"
          />
        </div>
      </div>
    </div>
  {/if}

  <div>
    <div class="text-right">
      {#if !ldapEnabled && !headerAuthEnabled && profile.rank !== 'GUEST' && toggleUpdatePassword}
        <button
          type="button"
          class="inline-block align-baseline font-bold
                                    text-sm text-blue-500 hover:text-blue-800 me-4"
          on:click="{toggleUpdatePassword}"
          data-testid="toggle-updatepassword"
        >
          {$LL.updatePassword()}
        </button>
      {/if}
      <SolidButton type="submit" disabled="{updateDisabled}">
        {$LL.updateProfile()}
      </SolidButton>
    </div>
  </div>
</form>

{#if showMFASetup}
  <SetupMFA
    notifications="{notifications}"
    xfetch="{xfetch}"
    eventTag="{eventTag}"
    toggleSetup="{toggleMfaSetup}"
    handleComplete="{handleMfaSetupCompletion}"
  />
{/if}

{#if showMfaRemove}
  <DeleteConfirmation
    toggleDelete="{toggleMfaRemove}"
    handleDelete="{handleMfaRemove}"
    confirmText="{$LL.mfa2faRemoveText()}"
    confirmBtnText="{$LL.remove()}"
  />
{/if}
