<script>
    export let hostname = ''
    export let battleId = ''

    function copyBattleLink() {
        const battleLink = document.getElementById('BattleLink')
        battleLink.select()
        document.execCommand('copy')
    }
</script>


<h4 class="text-xl mb-2">Invite a warrior</h4>

<div class="flex flex-wrap items-stretch w-full">
    <input
        class="flex-shrink flex-grow flex-auto leading-normal w-px flex-1 border h-10 border-grey-light rounded rounded-r-none px-3"
        type="text"
        value="{hostname}/battle/{battleId}"
        id="BattleLink"
        readonly
    />
    <div class="flex -mr-px">
        <button
            on:click={copyBattleLink}
            class="flex items-center leading-normal rounded rounded-l-none border border-blue-darker px-3 whitespace-no-wrap bg-blue hover:bg-blue-dark text-white font-bold text-sm"
        >
            <svg class="fill-current" width="20" height="20" viewBox="0 0 1792 1792" xmlns="http://www.w3.org/2000/svg">
                <path d="M768 1664h896v-640h-416q-40 0-68-28t-28-68v-416h-384v1152zm256-1440v-64q0-13-9.5-22.5t-22.5-9.5h-704q-13 0-22.5 9.5t-9.5 22.5v64q0 13 9.5 22.5t22.5 9.5h704q13 0 22.5-9.5t9.5-22.5zm256 672h299l-299-299v299zm512 128v672q0 40-28 68t-68 28h-960q-40 0-68-28t-28-68v-160h-544q-40 0-68-28t-28-68v-1344q0-40 28-68t68-28h1088q40 0 68 28t28 68v328q21 13 36 28l408 408q28 28 48 76t20 88z"/>
            </svg>
        </button>
    </div>	
</div>