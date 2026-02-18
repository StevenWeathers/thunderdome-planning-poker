<script lang="ts">
  import SolidButton from '../global/SolidButton.svelte';
  import Modal from '../global/Modal.svelte';
  import LL from '../../i18n/i18n-svelte';
  import { user } from '../../stores';
  import TextInput from '../forms/TextInput.svelte';
  import { ClipboardCopy } from '@lucide/svelte';

  import type { NotificationService } from '../../types/notifications';
  import { onMount } from 'svelte';

  interface Props {
    handleApiKeyCreate?: any;
    toggleCreateApiKey?: any;
    xfetch?: any;
    notifications: NotificationService;
  }

  let {
    handleApiKeyCreate = () => {},
    toggleCreateApiKey = () => {},
    xfetch = () => {},
    notifications,
  }: Props = $props();

  let keyName = $state('');
  let apiKey = $state('');

  function handleSubmit(event: Event) {
    event.preventDefault();

    if (keyName === '') {
      notifications.danger($LL.apiKeyNameInvalid());
      return false;
    }

    const body = {
      name: keyName,
    };

    xfetch(`/api/users/${$user.id}/apikeys`, { body })
      .then((res: Response) => res.json())
      .then(function (result: any) {
        handleApiKeyCreate();
        apiKey = result.data.apiKey;
      })
      .catch(function (error: any, response: any) {
        if (Array.isArray(error)) {
          error[1].json().then(function (result: any) {
            let errMessage;
            switch (result.error) {
              case 'USER_APIKEY_LIMIT_REACHED':
                errMessage = $LL.apiKeyLimitReached();
                break;
              case 'REQUIRES_VERIFIED_USER':
                errMessage = $LL.apiKeyUnverifiedUser();
                break;
              default:
                errMessage = $LL.apiKeyCreateFailed();
            }

            notifications.danger(errMessage);
          });
        } else {
          notifications.danger($LL.apiKeyCreateFailed());
        }
      });
  }

  function copyKey() {
    const apk = document.getElementById('apiKey') as HTMLInputElement;

    if (!navigator.clipboard) {
      apk?.select();
      document.execCommand('copy');
    } else {
      navigator.clipboard
        .writeText(apk?.value ?? '')
        .then(function () {
          notifications.success($LL.apikeyCopySuccess());
        })
        .catch(function () {
          notifications.danger($LL.apikeyCopyFailure());
        });
    }
  }

  let focusInput: any = $state();
  onMount(() => {
    focusInput?.focus();
  });
</script>

<Modal closeModal={toggleCreateApiKey} ariaLabel={$LL.modalCreateApiKey()}>
  {#if apiKey === ''}
    <form onsubmit={handleSubmit} name="createApiKey">
      <div class="mb-4">
        <label class="block dark:text-gray-400 font-bold mb-2" for="keyName">
          {$LL.apiKeyName()}
        </label>
        <TextInput
          id="keyName"
          name="keyName"
          bind:value={keyName}
          bind:this={focusInput}
          placeholder={$LL.apiKeyNamePlaceholder()}
          required
        />
      </div>
      <div class="text-right">
        <div>
          <SolidButton type="submit">
            {$LL.create()}
          </SolidButton>
        </div>
      </div>
    </form>
  {:else}
    <div class="mb-4">
      <p class="mb-3 mt-3 dark:text-white">
        {@html $LL.apiKeyCreateSuccess({
          keyName: `<span class="font-bold">${keyName}</span>`,
          onlyNowOpen: '<span class="font-bold">',
          onlyNowClose: '</span>',
        })}
      </p>
      <div class="flex flex-wrap items-stretch w-full mb-3">
        <input
          class="flex-shrink flex-grow flex-auto leading-normal w-px
                    flex-1 border-2 h-10 bg-gray-100 dark:bg-gray-900 border-gray-200 dark:border-gray-900 rounded
                    rounded-e-none px-4 appearance-none text-gray-800 dark:text-gray-400 font-bold
                    focus:outline-none focus:bg-white dark:focus:bg-gray-800 focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
          type="text"
          value={apiKey}
          id="apiKey"
          readonly
        />
        <div class="invisible md:visible md:flex md:-me-px">
          <SolidButton
            color="blue-copy"
            onClick={copyKey}
            additionalClasses="flex items-center leading-normal
                        whitespace-no-wrap text-sm"
          >
            <ClipboardCopy />
          </SolidButton>
        </div>
      </div>
      <p class="dark:text-white">
        {$LL.apiKeyStoreWarning()}
      </p>
    </div>
    <div class="text-right">
      <div>
        <SolidButton onClick={toggleCreateApiKey} testid="apikey-close">
          {$LL.close()}
        </SolidButton>
      </div>
    </div>
  {/if}
</Modal>
