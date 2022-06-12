<script>
    import { appRoutes } from '../../config'
    import SolidButton from '../SolidButton.svelte'
    import ClipboardIcon from '../icons/ClipboardIcon.svelte'

    export let hostname = ''
    export let retrospectiveId = ''

    function copyRetrospectiveLink() {
        const retrospectiveLink = document.getElementById('RetrospectiveLink')
        retrospectiveLink.select()
        document.execCommand('copy')
    }
</script>

<h4 class="text-xl mb-2 leading-tight font-bold">Invite a user</h4>

<div class="flex flex-wrap items-stretch w-full">
    <input
        class="flex-shrink flex-grow flex-auto leading-normal w-px flex-1
        border-2 h-10 bg-gray-200 border-gray-200 rounded rounded-r-none px-3
        appearance-none text-gray-700 focus:outline-none focus:bg-white
        focus:border-orange-500 "
        type="text"
        value="{hostname}{appRoutes.retro}/{retrospectiveId}"
        id="RetrospectiveLink"
        readonly
    />
    <div class="flex -mr-px">
        <SolidButton
            color="blue-copy"
            onClick="{copyRetrospectiveLink}"
            class="flex items-center leading-normal whitespace-no-wrap text-sm"
        >
            <ClipboardIcon />
        </SolidButton>
    </div>
</div>
