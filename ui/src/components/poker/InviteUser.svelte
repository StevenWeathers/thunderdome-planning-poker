<script lang="ts">
  import SolidButton from '../global/SolidButton.svelte';
  import LL from '../../i18n/i18n-svelte';
  import { appRoutes } from '../../config';
  import { ClipboardCopy } from 'lucide-svelte';

  interface Props {
    notifications: any;
    hostname?: string;
    battleId?: string;
    joinCode?: string;
  }

  let {
    notifications,
    hostname = '',
    battleId = '',
    joinCode = ''
  }: Props = $props();

  function copyBattleLink() {
    const bl = document.getElementById('BattleLink');

    if (!navigator.clipboard) {
      bl.select();
      document.execCommand('copy');
    } else {
      navigator.clipboard
        .writeText(bl.value)
        .then(function () {
          notifications.success($LL.inviteLinkCopySuccess());
        })
        .catch(function () {
          notifications.danger($LL.inviteLinkCopyFailure());
        });
    }
  }

  function copyJoinCode() {
    const jc = document.getElementById('JoinCode');

    if (!navigator.clipboard) {
      jc.select();
      document.execCommand('copy');
    } else {
      navigator.clipboard
        .writeText(jc.value)
        .then(function () {
          notifications.success($LL.joinCodeCopySuccess());
        })
        .catch(function () {
          notifications.danger($LL.joinCodeCopyFailure());
        });
    }
  }
</script>

<div class="w-full">
  <h4
    class="text-2xl mb-2 leading-tight font-semibold font-rajdhani uppercase dark:text-white"
  >
    {$LL.warriorInvite()}
  </h4>
  <div class="flex flex-wrap items-stretch w-full">
    <input
      class="flex-shrink flex-grow flex-auto leading-normal w-px flex-1
        border-2 h-10 bg-gray-100 dark:bg-gray-900 dark:focus:bg-gray-800 border-gray-200 dark:border-gray-900 rounded rounded-e-none px-3
        appearance-none text-gray-700 dark:text-gray-400 focus:outline-none focus:bg-white
        focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
      type="text"
      value="{hostname}{appRoutes.game}/{battleId}"
      id="BattleLink"
      readonly
    />
    <div class="flex -mr-px">
      <SolidButton
        color="blue-copy"
        onClick={copyBattleLink}
        additionalClasses="flex items-center leading-normal
            whitespace-no-wrap text-sm"
      >
        <ClipboardCopy />
      </SolidButton>
    </div>
  </div>
  {#if joinCode !== ''}
    <div class="mt-4">
      <label for="JoinCode" class="font-bold dark:text-gray-300"
        >{$LL.passCode()}</label
      >
      <div class="flex flex-wrap items-stretch w-full">
        <input
          class="flex-shrink flex-grow flex-auto leading-normal w-px flex-1
            border-2 h-10 bg-gray-100 dark:bg-gray-900 border-gray-200 dark:border-gray-900 rounded rounded-e-none px-3
            appearance-none text-gray-700 dark:text-gray-400 focus:outline-none focus:bg-white dark:focus:bg-gray-800
            focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
          type="text"
          value="{joinCode}"
          id="JoinCode"
          readonly
        />
        <div class="flex -mr-px">
          <SolidButton
            color="blue-copy"
            onClick={copyJoinCode}
            additionalClasses="flex items-center leading-normal
                whitespace-no-wrap text-sm"
          >
            <ClipboardCopy />
          </SolidButton>
        </div>
      </div>
    </div>
  {/if}
</div>
