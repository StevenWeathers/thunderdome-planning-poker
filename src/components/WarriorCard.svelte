<script>
    import VoteIcon from './VoteIcon.svelte'

    export let voted = false
    export let warrior = {}
    export let isLeader = false
    export let leaderId = ''
    export let points = ''
    export let sendSocketEvent = () => {}

    function promoteLeader() {
        sendSocketEvent('promote_leader', warrior.id)
    }

    function jabWarrior() {
        sendSocketEvent('jab_warrior', warrior.id)
    }
</script>

<div class="border-b border-grey p-4 flex" data-testId="warriorCard" data-warriorName={warrior.name}>
    <div class="w-1/4">
        <img src="https://api.adorable.io/avatars/48/{warrior.id}.png" alt="Placeholder Avatar">
    </div>
    <div class="w-3/4">
        <div class="flex">
            <div class="w-3/4">
                <p class="text-xl font-bold" data-testId="warriorName">{warrior.name}</p>
                {#if leaderId === warrior.id}
                    <p class="text-l text-grey-darker">Leader</p>
                {:else if isLeader}
                    <button
                        on:click={promoteLeader}
                        class="inline-block align-baseline text-sm text-blue hover:text-blue-darker bg-transparent border-transparent"
                    >
                        Promote
                    </button>&nbsp;|&nbsp;
                    <button
                        on:click={jabWarrior}
                        class="inline-block align-baseline text-sm text-blue hover:text-blue-darker bg-transparent border-transparent"
                    >
                        Nudge
                    </button>
                {/if}
            </div>
            <div class="w-1/4 text-right">
                {#if voted && points === ''}
                    <VoteIcon />
                {:else if voted && points !== ''}
                    <span class="font-bold text-green-dark border-green border p-2 rounded ml-2" data-testId="warriorPoints">{points}</span>
                {/if}
            </div>
        </div>
    </div>
</div>