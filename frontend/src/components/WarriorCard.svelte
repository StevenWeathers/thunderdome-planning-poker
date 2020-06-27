<script>
    import VoteIcon from './icons/VoteIcon.svelte'
    import LeaderIcon from './icons/LeaderIcon.svelte'
    import WarriorRankPrivate from './icons/WarriorRankPrivate.svelte'
    import WarriorRankCorporal from './icons/WarriorRankCorporal.svelte'
    import WarriorRankGeneral from './icons/WarriorRankGeneral.svelte'

    export let voted = false
    export let warrior = {}
    export let isLeader = false
    export let leaderId = ''
    export let points = ''
    export let sendSocketEvent = () => {}
    export let eventTag

    const showRank = appConfig.ShowWarriorRank
    let nameStyleClass = showRank ? "text-lg" : "text-xl"

    function promoteLeader() {
        sendSocketEvent('promote_leader', warrior.id)
        eventTag('promote_leader', 'battle', '')
    }

    function jabWarrior() {
        sendSocketEvent('jab_warrior', warrior.id)
        eventTag('jab_warrior', 'battle', '')
    }
</script>

<div
    class="border-b border-gray-500 p-4 flex items-center"
    data-testId="warriorCard"
    data-warriorName="{warrior.name}">
    <div class="w-1/4">
        <img
            src="https://api.adorable.io/avatars/48/{warrior.id}.png"
            alt="Placeholder Avatar" />
    </div>
    <div class="w-3/4">
        <div class="flex items-center">
            <div class="w-3/4">
                <p
                    class="{nameStyleClass} font-bold leading-tight truncate"
                    data-testId="warriorName"
                    title="{warrior.name}">
                    {#if showRank}
                        {#if warrior.rank == 'GENERAL'}
                            <WarriorRankGeneral />
                        {:else if warrior.rank == 'CORPORAL'}
                            <WarriorRankCorporal />
                        {:else}
                            <WarriorRankPrivate />
                        {/if}
                    {/if}
                    {warrior.name}
                </p>
                {#if leaderId === warrior.id}
                    <p class="text-l text-gray-700 leading-tight">
                        <LeaderIcon />
                        &nbsp;Leader
                    </p>
                {:else if isLeader}
                    <button
                        on:click="{promoteLeader}"
                        class="inline-block align-baseline text-sm text-blue-500
                        hover:text-blue-600 bg-transparent border-transparent">
                        Promote
                    </button>
                    &nbsp;|&nbsp;
                    <button
                        on:click="{jabWarrior}"
                        class="inline-block align-baseline text-sm text-blue-500
                        hover:text-blue-600 bg-transparent border-transparent">
                        Nudge
                    </button>
                {/if}
            </div>
            <div class="w-1/4 text-right">
                {#if voted && points === ''}
                    <span class="text-green-500">
                        <VoteIcon />
                    </span>
                {:else if voted && points !== ''}
                    <span
                        class="font-bold text-green-600 border-green-500 border
                        p-2 rounded ml-2"
                        data-testId="warriorPoints">
                        {points}
                    </span>
                {/if}
            </div>
        </div>
    </div>
</div>
