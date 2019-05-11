<script>
    export let voted = false
    export let warrior = {}
    export let isLeader = false
    export let points = ''
</script>

<div class="border-b border-grey p-4 flex">
    <div class="w-1/4">
        <img src="https://api.adorable.io/avatars/48/{warrior.id}.png" alt="Placeholder Avatar">
    </div>
    <div class="w-3/4">
        <div class="flex">
            <div class="w-3/4">
                <p class="text-xl">{warrior.name}</p>
                {#if isLeader}<p class="text-l">Leader</p>{/if}
            </div>
            <div class="w-1/4 text-right">
                {#if voted && points === ''}
                    <svg class="fill-current text-green" width="32" height="32" viewBox="0 0 1792 1792" xmlns="http://www.w3.org/2000/svg">
                        <path d="M813 1299l614-614q19-19 19-45t-19-45l-102-102q-19-19-45-19t-45 19l-467 467-211-211q-19-19-45-19t-45 19l-102 102q-19 19-19 45t19 45l358 358q19 19 45 19t45-19zm851-883v960q0 119-84.5 203.5t-203.5 84.5h-960q-119 0-203.5-84.5t-84.5-203.5v-960q0-119 84.5-203.5t203.5-84.5h960q119 0 203.5 84.5t84.5 203.5z"/>
                    </svg>
                {:else if voted && points !== ''}
                    <span class="font-bold text-green-dark border-green border px-3 py-2 rounded ml-2">{points}</span>
                {/if}
            </div>
        </div>
    </div>
</div>