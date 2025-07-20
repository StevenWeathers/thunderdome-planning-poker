<script lang="ts">
  import Modal from '../global/Modal.svelte';
  import SolidButton from '../global/SolidButton.svelte';

  import LL from '../../i18n/i18n-svelte';
  import TextInput from '../forms/TextInput.svelte';
  import { Shield } from 'lucide-svelte';

  interface Props {
    toggleSetup?: any;
    handleComplete?: any;
    xfetch: any;
    notifications: any;
  }

  let {
    toggleSetup = () => {},
    handleComplete = () => {},
    xfetch,
    notifications
  }: Props = $props();

  let qrCode = $state('');
  let secret = $state('');
  let passcode = $state('');

  xfetch('/api/auth/mfa/setup/generate', { method: 'POST' })
    .then(res => res.json())
    .then(r => {
      qrCode = r.data.qrCode;
      secret = r.data.secret;
    })
    .catch(() => {
      notifications.danger($LL.mfaGenerateFailed());
    });

  function onSubmit(e) {
    e.preventDefault();
    xfetch('/api/auth/mfa/setup/validate', { body: { secret, passcode } })
      .then(res => res.json())
      .then(r => {
        if (r.data.result === 'SUCCESS') {
          notifications.success($LL.mfaSetupSuccess());
          handleComplete();
        } else {
          notifications.danger(`${r.data.result}`);
        }
      })
      .catch(() => {
        notifications.danger($LL.mfaSetupFailed());
      });
  }

  let submitDisabled = $derived(passcode === '');
</script>

<Modal closeModal={toggleSetup} widthClasses="md:w-2/3 lg:w-1/2">
  <div class="pt-12">
    <div class="dark:text-gray-300 text-center">
      <p class="font-rajdhani text-lg mb-2">
        {$LL.mfaSetupIntro()}
      </p>
      {#if qrCode !== ''}
        <img
          src="data:image/png;base64,{qrCode}"
          class="m-auto"
          alt="MFA QR Code"
        />

        <p class="mt-2 font-rajdhani text-xl text-red-500">
          {$LL.mfaSecretKeyLabel()}: {secret}
        </p>
      {/if}
    </div>
    <form onsubmit={onSubmit} name="validateMFAPasscode" class="mt-8">
      <div class="mb-4">
        <label
          class="block text-gray-700 dark:text-gray-400 font-bold mb-2"
          for="mfaPasscode"
        >
          {$LL.mfaTokenLabel()}
        </label>
        <TextInput
          bind:value="{passcode}"
          placeholder={$LL.mfaTokenPlaceholder()}
          id="mfaPasscode"
          name="twofactor_token"
          type="password"
          required
          icon={Shield}
        />
      </div>

      <div>
        <div class="text-right">
          <SolidButton type="submit" disabled={submitDisabled}>
            {$LL.mfaConfirmToken()}
          </SolidButton>
        </div>
      </div>
    </form>
  </div>
</Modal>
