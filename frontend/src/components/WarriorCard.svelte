<script>
    import VoteIcon from './icons/VoteIcon.svelte'
    import LeaderIcon from './icons/LeaderIcon.svelte'
    import WarriorRankPrivate from './icons/WarriorRankPrivate.svelte'
    import WarriorRankCorporal from './icons/WarriorRankCorporal.svelte'
    import WarriorRankGeneral from './icons/WarriorRankGeneral.svelte'
    import WarriorAvatar from './WarriorAvatar.svelte'
    import { _ } from '../i18n'
    import { warrior as activeWarrior } from '../stores.js'

    export let voted = false
    export let warrior = {}
    export let isLeader = false
    export let autoFinishVoting = false
    export let leaders = []
    export let points = ''
    export let sendSocketEvent = () => {}
    export let eventTag

    const showRank = appConfig.ShowWarriorRank
    const avatarService = appConfig.AvatarService
    let nameStyleClass = showRank ? 'text-lg' : 'text-xl'

    function promoteLeader() {
        sendSocketEvent('promote_leader', warrior.id)
        eventTag('promote_leader', 'battle', '')
    }

    function demoteLeader() {
        sendSocketEvent('demote_leader', warrior.id)
        eventTag('demote_leader', 'battle', '')
    }

    function jabWarrior() {
        sendSocketEvent('jab_warrior', warrior.id)
        eventTag('jab_warrior', 'battle', '')
    }

    function toggleSpectator() {
        sendSocketEvent(
            'spectator_toggle',
            JSON.stringify({
                spectator: !warrior.spectator,
            }),
        )
        eventTag(`spectator_toggle`, 'battle', '')
    }
</script>

<div
    class="border-b border-gray-500 p-4 flex items-center"
    data-testId="warriorCard"
    data-warriorName="{warrior.name}"
>
    <div class="w-1/4 mr-1">
        <WarriorAvatar
            warriorId="{warrior.id}"
            avatar="{warrior.avatar}"
            avatarService="{avatarService}"
        />
    </div>
    <div class="w-3/4">
        <div class="flex items-center">
            <div class="w-3/4">
                <p
                    class="{nameStyleClass} font-bold leading-tight truncate"
                    data-testId="warriorName"
                    title="{warrior.name}"
                >
                    {#if showRank}
                        {#if warrior.rank == 'ADMIN'}
                            <WarriorRankGeneral />
                        {:else if warrior.rank == 'REGISTERED'}
                            <WarriorRankCorporal />
                        {:else}
                            <WarriorRankPrivate />
                        {/if}
                    {/if}
                    {#if autoFinishVoting && warrior.spectator}
                        <span class="text-gray-600" title="{$_('spectator')}">
                            {warrior.name}
                        </span>
                    {:else}
                        <span>{warrior.name}</span>
                    {/if}
                </p>
                {#if leaders.includes(warrior.id)}
                    <p class="text-l text-gray-700 leading-tight">
                        <LeaderIcon />
                        {#if isLeader}
                            &nbsp;
                            <button
                                on:click="{demoteLeader}"
                                class="inline text-sm text-red-500
                                hover:text-red-800 bg-transparent
                                border-transparent"
                            >
                                {$_('demote')}
                            </button>
                        {:else}&nbsp;{$_('pages.battle.warriorLeader')}{/if}
                    </p>
                {:else if isLeader}
                    <button
                        on:click="{promoteLeader}"
                        class="inline-block align-baseline text-sm
                        text-green-500 hover:text-green-800 bg-transparent
                        border-transparent"
                    >
                        {$_('promote')}
                    </button>
                    {#if !warrior.spectator}
                        &nbsp;|&nbsp;
                        <button
                            on:click="{jabWarrior}"
                            class="inline-block align-baseline text-sm
                            text-blue-500 hover:text-blue-800 bg-transparent
                            border-transparent"
                        >
                            {$_('warriorNudge')}
                        </button>
                    {/if}
                {/if}
                {#if autoFinishVoting && warrior.id === $activeWarrior.id}
                    <button
                        on:click="{toggleSpectator}"
                        class="inline-block align-baseline text-sm text-blue-500
                        hover:text-blue-800 bg-transparent border-transparent"
                    >
                        {#if !warrior.spectator}
                            {$_('becomeSpectator')}
                        {:else}{$_('becomeParticipant')}{/if}
                    </button>
                {/if}
            </div>
            <div class="w-1/4 text-right">
                {#if !warrior.spectator}
                    {#if voted && points === ''}
                        <span class="text-green-500">
                            <VoteIcon />
                        </span>
                    {:else if voted && points !== ''}
                        <span
                            class="font-bold text-green-600 border-green-500
                            border p-2 rounded ml-2"
                            data-testId="warriorPoints"
                        >
                            {points}
                        </span>
                    {/if}
                {/if}
            </div>
        </div>
    </div>
</div>
