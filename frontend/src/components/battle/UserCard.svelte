<script>
    import VoteIcon from '../icons/VoteIcon.svelte'
    import LeaderIcon from '../icons/LeaderIcon.svelte'
    import WarriorRankPrivate from '../icons/UserRankGuest.svelte'
    import WarriorRankCorporal from '../icons/UserRankRegistered.svelte'
    import WarriorRankGeneral from '../icons/UserRankAdmin.svelte'
    import UserAvatar from '../user/UserAvatar.svelte'
    import BecomeLeader from './BecomeLeader.svelte'
    import { AppConfig } from '../../config.js'
    import { _ } from '../../i18n.js'
    import { warrior as activeWarrior } from '../../stores.js'

    export let voted = false
    export let warrior = {}
    export let isLeader = false
    export let autoFinishVoting = false
    export let leaders = []
    export let points = ''
    export let sendSocketEvent = () => {}
    export let eventTag

    const showRank = AppConfig.ShowWarriorRank
    let nameStyleClass = showRank ? 'text-lg' : 'text-xl'
    let showBecomeLeader = false

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

    function becomeLeader(leaderCode) {
        sendSocketEvent('become_leader', leaderCode)
        eventTag('become_leader', 'battle', '')
        toggleBecomeLeader()
    }

    function toggleBecomeLeader() {
        showBecomeLeader = !showBecomeLeader
        eventTag('toggle_become_leader', 'battle', '')
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
    class="border-b border-gray-300 dark:border-gray-700 p-4 flex items-center"
    data-testid="user-card"
    data-username="{warrior.name}"
    data-userid="{warrior.id}"
>
    <div class="w-1/4 mr-2">
        <UserAvatar
            warriorId="{warrior.id}"
            avatar="{warrior.avatar}"
            gravatarHash="{warrior.gravatarHash}"
            class="rounded-full"
            width="68"
        />
    </div>
    <div class="w-3/4">
        <div class="flex items-center">
            <div class="w-3/4">
                <p
                    class="{nameStyleClass} font-bold leading-tight truncate dark:text-gray-300"
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
                        <span
                            class="text-gray-600"
                            title="{$_('spectator')}"
                            data-testid="user-name"
                        >
                            {warrior.name}
                        </span>
                    {:else}
                        <span data-testid="user-name">{warrior.name}</span>
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
                                data-testid="user-demote"
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
                        data-testid="user-promote"
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
                            data-testid="user-nudge"
                        >
                            {$_('warriorNudge')}
                        </button>
                    {/if}
                {:else if warrior.id === $activeWarrior.id}
                    <button
                        on:click="{toggleBecomeLeader}"
                        class="inline-block align-baseline text-sm
                        text-blue-500 hover:text-blue-800 bg-transparent
                        border-transparent"
                        data-testid="user-becomeleader"
                    >
                        {$_('becomeLeader')}
                    </button>
                {/if}
                {#if autoFinishVoting && warrior.id === $activeWarrior.id}
                    <button
                        on:click="{toggleSpectator}"
                        class="inline-block align-baseline text-sm text-blue-500
                        hover:text-blue-800 bg-transparent border-transparent"
                        data-testid="user-togglespectator"
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
                        <span class="text-green-500 dark:text-lime-400">
                            <VoteIcon class="h-8 w-8" />
                        </span>
                    {:else if voted && points !== ''}
                        <span
                            class="font-bold text-green-600 dark:text-lime-400 border-green-500 dark:border-lime-500
                            border p-2 rounded ml-2"
                            data-testid="user-points"
                        >
                            {points}
                        </span>
                    {/if}
                {/if}
            </div>
        </div>
    </div>

    {#if showBecomeLeader}
        <BecomeLeader
            handleBecomeLeader="{becomeLeader}"
            toggleBecomeLeader="{toggleBecomeLeader}"
        />
    {/if}
</div>
