<script lang="ts">
  import SolidButton from '../global/SolidButton.svelte';
  import HollowButton from '../global/HollowButton.svelte';
  import { countryList } from '../../country';
  import { AppConfig } from '../../config';
  import { validateName, validateUserIsAdmin } from '../../validationUtils';
  import LL, { locale, setLocale } from '../../i18n/i18n-svelte';
  import UserAvatar from './UserAvatar.svelte';
  import SetupMFA from '../auth/SetupMFA.svelte';
  import DeleteConfirmation from '../global/DeleteConfirmation.svelte';
  import { user } from '../../stores';
  import LocaleSwitcher from '../forms/LocaleInput.svelte';
  import type { Locales } from '../../i18n/i18n-types';
  import { loadLocaleAsync } from '../../i18n/i18n-util.async';
  import TextInput from '../forms/TextInput.svelte';
  import SelectInput from '../forms/SelectInput.svelte';
  import Checkbox from '../forms/Checkbox.svelte';
  import { BadgeCheck, Building, Mail } from '@lucide/svelte';

  import type { NotificationService } from '../../types/notifications';
  import type { ApiClient } from '../../types/apiclient';

  const setupI18n = async (locale: Locales) => {
    await loadLocaleAsync(locale);
    setLocale(locale);
  };

  interface Props {
    profile?: any;
    credential: any;
    handleUpdate?: any;
    toggleUpdatePassword: any;
    notifications: NotificationService;
    xfetch: ApiClient;
  }

  let {
    profile = $bindable({
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
      theme: 'auto',
    }),
    credential = $bindable(),
    handleUpdate = () => {},
    toggleUpdatePassword,
    notifications,
    xfetch,
  }: Props = $props();

  const { AvatarService, LdapEnabled, HeaderAuthEnabled, OIDCAuthEnabled } = AppConfig;

  const themes = ['auto', 'light', 'dark'];
  const configurableAvatarServices = ['gravatar', 'robohash', 'govatar'];
  const isAvatarConfigurable = configurableAvatarServices.includes(AvatarService);
  const avatarOptions = {
    gravatar: ['mp', 'identicon', 'monsterid', 'wavatar', 'retro', 'robohash', 'none'],
    robohash: ['set1', 'set2', 'set3', 'set4', 'none'],
    govatar: ['male', 'female', 'none'],
  };

  let avatars = isAvatarConfigurable ? avatarOptions[AvatarService as keyof typeof avatarOptions] : [];

  function handleSubmit(e: Event) {
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

  let showMFASetup = $state(false);

  function toggleMfaSetup() {
    showMFASetup = !showMFASetup;
  }

  function handleMfaSetupCompletion() {
    credential.mfa_enabled = true;
    toggleMfaSetup();
  }

  let showMfaRemove = $state(false);

  function toggleMfaRemove() {
    showMfaRemove = !showMfaRemove;
  }

  function handleMfaRemove() {
    xfetch('/api/auth/mfa', { method: 'DELETE' })
      .then(() => {
        credential.mfa_enabled = false;
        toggleMfaRemove();
        notifications.success($LL.mfa2faRemoveSuccess());
      })
      .catch(() => {
        notifications.danger($LL.mfa2faRemoveFailure());
      });
  }

  function requestVerifyEmail(e: Event) {
    e.preventDefault();
    xfetch(`/api/users/${profile.id}/request-verify`, { method: 'POST' })
      .then(function () {
        notifications.success($LL.requestVerifyEmailSuccess());
      })
      .catch(function () {
        notifications.danger($LL.requestVerifyEmailFailure());
      });
  }

  function requestEmailChange(e: Event) {
    e.preventDefault();
    xfetch(`/api/users/${profile.id}/email-change`, { method: 'POST' })
      .then(function () {
        notifications.success($LL.requestEmailChangeSuccess(), 2500);
      })
      .catch(function () {
        notifications.danger($LL.requestEmailChangeError());
      });
  }

  let updateDisabled = $derived(profile.name === '');
  let userIsAdmin = $derived(validateUserIsAdmin($user));
</script>

<form onsubmit={handleSubmit} name="updateProfile">
  <div class="mb-4">
    <label class="block text-gray-700 dark:text-gray-400 font-bold mb-2" for="yourName">
      {$LL.name()}
    </label>
    <TextInput
      bind:value={profile.name}
      placeholder={$LL.yourNamePlaceholder()}
      id="yourName"
      name="yourName"
      disabled={LdapEnabled || HeaderAuthEnabled || OIDCAuthEnabled}
      required
    />
  </div>

  <div class="mb-4">
    <label
      class="flex items-center justify-between gap-2 text-gray-700 dark:text-gray-400 font-bold mb-2"
      for="yourEmail"
    >
      <span class="inline-flex items-center">
        {$LL.email()}
        {#if profile.rank && profile.rank !== 'GUEST'}
          {#if profile.verified}
            <span
              class="inline-block font-bold text-green-600
                                        border-green-500 border py-1 px-2 rounded
                                        ms-1"
              data-testid="user-verified"
            >
              {$LL.verified()}
              <BadgeCheck class="inline h-4 w-4" />
            </span>
          {/if}
        {/if}
      </span>

      {#if profile.rank && profile.rank !== 'GUEST'}
        {#if !profile.verified || (!LdapEnabled && !HeaderAuthEnabled && !OIDCAuthEnabled)}
          <span class="inline-flex items-center gap-2">
            {#if !profile.verified}
              <button
                class="inline-block align-baseline font-bold text-sm text-blue-500
                                                hover:text-blue-800"
                onclick={requestVerifyEmail}
                data-testid="request-verify"
                type="button"
                >{$LL.requestVerifyEmail()}
              </button>
            {/if}

            {#if !LdapEnabled && !HeaderAuthEnabled && !OIDCAuthEnabled}
              <button
                class="inline-block align-baseline font-bold text-sm text-blue-500
                                                hover:text-blue-800"
                onclick={requestEmailChange}
                data-testid="request-email-change"
                type="button"
                >{$LL.changeEmail()}
              </button>
            {/if}
          </span>
        {/if}
      {/if}
    </label>
    <TextInput
      bind:value={profile.email}
      id="yourEmail"
      name="yourEmail"
      type="email"
      disabled={!userIsAdmin}
      icon={Mail}
    />
  </div>

  {#if profile.rank !== 'GUEST' && credential}
    <div class="mb-4">
      <p class="block text-gray-700 dark:text-gray-400 font-bold mb-2">
        {$LL.mfa2faLabel()}
      </p>
      {#if !credential.mfa_enabled}
        <HollowButton color="teal" onClick={toggleMfaSetup}>{$LL.mfa2faSetup()}</HollowButton>
      {:else}
        <HollowButton color="red" onClick={toggleMfaRemove}>{$LL.mfa2faRemove()}</HollowButton>
      {/if}
    </div>
  {/if}

  <div class="mb-4">
    <label class="block text-gray-700 dark:text-gray-400 font-bold mb-2" for="yourCountry">
      {$LL.country()}
    </label>

    <SelectInput bind:value={profile.country} id="yourCountry" name="yourCountry">
      <option value="">
        {$LL.chooseCountryPlaceholder()}
      </option>
      {#each countryList as item}
        <option value={item.abbrev}>
          {item.name} [{item.abbrev}]
        </option>
      {/each}
    </SelectInput>
  </div>

  <div class="mb-4">
    <div class="text-gray-700 dark:text-gray-400 font-bold mb-2">
      {$LL.locale()}
    </div>
    <LocaleSwitcher selectedLocale={$locale} on:locale-changed={e => setupI18n(e.detail)} />
  </div>

  <div class="mb-4">
    <label class="block text-gray-700 dark:text-gray-400 font-bold mb-2" for="yourCompany">
      {$LL.company()}
    </label>
    <TextInput
      bind:value={profile.company}
      placeholder={$LL.companyPlaceholder()}
      id="yourCompany"
      name="yourCompany"
      icon={Building}
    />
  </div>

  <div class="mb-4">
    <label class="block text-gray-700 dark:text-gray-400 font-bold mb-2" for="yourJobTitle">
      {$LL.jobTitle()}
    </label>
    <TextInput
      bind:value={profile.jobTitle}
      placeholder={$LL.jobTitlePlaceholder()}
      id="yourJobTitle"
      name="yourJobTitle"
    />
  </div>

  <div class="mb-4">
    <label class="block text-gray-700 dark:text-gray-400 font-bold mb-2" for="yourCountry">
      {$LL.theme()}
    </label>
    <SelectInput bind:value={profile.theme} id="theme" name="theme">
      {#each themes as theme}
        <option value={theme}>
          {theme}
        </option>
      {/each}
    </SelectInput>
  </div>

  <div class="mb-4">
    <Checkbox bind:checked={profile.notificationsEnabled} label={$LL.enableBattleNotifications()} />
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
          <SelectInput bind:value={profile.avatar} id="yourAvatar" name="yourAvatar">
            {#each avatars as item}
              <option value={item}>
                {item}
              </option>
            {/each}
          </SelectInput>
        </div>
        <div class="shrink ms-4">
          <UserAvatar
            warriorId={profile.id}
            avatar={profile.avatar}
            gravatarHash={profile.gravatarHash}
            userName={profile.name}
            width={48}
            class="rounded-full"
          />
        </div>
      </div>
    </div>
  {/if}

  <div>
    <div class="text-right">
      {#if profile.rank && profile.rank !== 'GUEST'}
        {#if !LdapEnabled && !HeaderAuthEnabled && !OIDCAuthEnabled && toggleUpdatePassword}
          <button
            type="button"
            class="inline-block align-baseline font-bold
                                    text-sm text-blue-500 hover:text-blue-800 me-4"
            onclick={toggleUpdatePassword}
            data-testid="toggle-updatepassword"
          >
            {$LL.updatePassword()}
          </button>
        {/if}
      {/if}
      <SolidButton type="submit" disabled={updateDisabled}>
        {$LL.updateProfile()}
      </SolidButton>
    </div>
  </div>
</form>

{#if showMFASetup}
  <SetupMFA {notifications} {xfetch} toggleSetup={toggleMfaSetup} handleComplete={handleMfaSetupCompletion} />
{/if}

{#if showMfaRemove}
  <DeleteConfirmation
    toggleDelete={toggleMfaRemove}
    handleDelete={handleMfaRemove}
    confirmText={$LL.mfa2faRemoveText()}
    confirmBtnText={$LL.remove()}
  />
{/if}
