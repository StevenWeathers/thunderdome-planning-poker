<script>
    import SolidButton from '../SolidButton.svelte'
    import ClipboardIcon from '../icons/ClipboardIcon.svelte'
    import { _ } from '../../i18n'
    import { appRoutes } from '../../config'

    export let notifications
    export let hostname = ''
    export let battleId = ''
    export let joinCode = ''

    function copyBattleLink() {
        const bl = document.getElementById('BattleLink')

        if (!navigator.clipboard) {
            bl.select()
            document.execCommand('copy')
        } else {
            navigator.clipboard
                .writeText(bl.value)
                .then(function () {
                    notifications.success($_('inviteLinkCopySuccess'))
                })
                .catch(function () {
                    notifications.danger($_('inviteLinkCopyFailure'))
                })
        }
    }

    function copyJoinCode() {
        const jc = document.getElementById('JoinCode')

        if (!navigator.clipboard) {
            jc.select()
            document.execCommand('copy')
        } else {
            navigator.clipboard
                .writeText(jc.value)
                .then(function () {
                    notifications.success($_('joinCodeCopySuccess'))
                })
                .catch(function () {
                    notifications.danger($_('joinCodeCopyFailure'))
                })
        }
    }
</script>

<div class="w-full">
    <h4 class="text-xl mb-2 leading-tight font-bold">
        {$_('pages.battle.warriorInvite')}
    </h4>
    <div class="flex flex-wrap items-stretch w-full">
        <input
            class="flex-shrink flex-grow flex-auto leading-normal w-px flex-1
        border-2 h-10 bg-gray-200 border-gray-200 rounded rounded-r-none px-3
        appearance-none text-gray-700 focus:outline-none focus:bg-white
        focus:border-purple-500"
            type="text"
            value="{hostname}{appRoutes.battle}/{battleId}"
            id="BattleLink"
            readonly
        />
        <div class="flex -mr-px">
            <SolidButton
                color="blue-copy"
                onClick="{copyBattleLink}"
                additionalClasses="flex items-center leading-normal
            whitespace-no-wrap text-sm"
            >
                <ClipboardIcon />
            </SolidButton>
        </div>
    </div>
    {#if joinCode !== ''}
        <div class="mt-4">
            <label for="JoinCode" class="font-bold">{$_('passCode')}</label>
            <div class="flex flex-wrap items-stretch w-full">
                <input
                    class="flex-shrink flex-grow flex-auto leading-normal w-px flex-1
            border-2 h-10 bg-gray-200 border-gray-200 rounded rounded-r-none px-3
            appearance-none text-gray-700 focus:outline-none focus:bg-white
            focus:border-purple-500"
                    type="text"
                    value="{joinCode}"
                    id="JoinCode"
                    readonly
                />
                <div class="flex -mr-px">
                    <SolidButton
                        color="blue-copy"
                        onClick="{copyJoinCode}"
                        additionalClasses="flex items-center leading-normal
                whitespace-no-wrap text-sm"
                    >
                        <ClipboardIcon />
                    </SolidButton>
                </div>
            </div>
        </div>
    {/if}
</div>
