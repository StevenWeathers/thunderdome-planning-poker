<script>
    import Sockette from 'sockette'
    import { onDestroy, onMount } from 'svelte'

    import PageLayout from '../components/PageLayout.svelte'
    import PointCard from '../components/battle/PointCard.svelte'
    import WarriorCard from '../components/battle/UserCard.svelte'
    import BattlePlans from '../components/battle/BattlePlans.svelte'
    import VotingControls from '../components/battle/VotingControls.svelte'
    import InviteWarrior from '../components/battle/InviteUser.svelte'
    import VoteResults from '../components/battle/VoteResults.svelte'
    import HollowButton from '../components/HollowButton.svelte'
    import SolidButton from '../components/SolidButton.svelte'
    import ExternalLinkIcon from '../components/icons/ExternalLinkIcon.svelte'
    import EditBattle from '../components/battle/EditBattle.svelte'
    import DeleteConfirmation from '../components/DeleteConfirmation.svelte'
    import { warrior } from '../stores.js'
    import { _ } from '../i18n.js'
    import { AppConfig, appRoutes, PathPrefix } from '../config.js'

    export let battleId
    export let notifications
    export let eventTag
    export let router

    const { AllowRegistration, AllowGuests } = AppConfig
    const loginOrRegister =
        AllowRegistration || AllowGuests ? appRoutes.register : appRoutes.login

    const hostname = window.location.origin
    const socketExtension = window.location.protocol === 'https:' ? 'wss' : 'ws'
    const defaultPlan = {
        id: '',
        name: `[${$_('pages.battle.votingNotStarted')}]`,
        type: '',
        referenceId: '',
        link: '',
        description: '',
        acceptanceCriteria: '',
    }

    let JoinPassRequired = false
    let socketError = false
    let socketReconnecting = false
    let points = []
    let vote = ''
    let voteStartTime = new Date()
    let battle = { leaders: [] }
    let currentPlan = { ...defaultPlan }
    let currentTime = new Date()
    let showEditBattle = false
    let showDeleteBattle = false
    let isSpectator = false
    let joinPasscode = ''

    $: countdown =
        battle.currentPlanId !== '' && battle.votingLocked === false
            ? timeUnitsBetween(voteStartTime, currentTime)
            : {}

    const onSocketMessage = function (evt) {
        const parsedEvent = JSON.parse(evt.data)

        switch (parsedEvent.type) {
            case 'join_code_required':
                JoinPassRequired = true
                break
            case 'join_code_incorrect':
                notifications.danger($_('incorrectPassCode'))
                break
            case 'init': {
                JoinPassRequired = false
                battle = JSON.parse(parsedEvent.value)
                points = battle.pointValuesAllowed
                const { spectator = false } =
                    battle.users.find(w => w.id === $warrior.id) || {}
                isSpectator = spectator

                if (battle.activePlanId !== '') {
                    const activePlan = battle.plans.find(
                        p => p.id === battle.activePlanId,
                    )
                    const warriorVote = activePlan.votes.find(
                        v => v.warriorId === $warrior.id,
                    ) || { vote: '' }
                    currentPlan = activePlan
                    voteStartTime = new Date(activePlan.voteStartTime)
                    vote = warriorVote.vote
                }

                eventTag('join', 'battle', '')
                break
            }
            case 'warrior_joined': {
                battle.users = JSON.parse(parsedEvent.value)
                const joinedWarrior = battle.users.find(
                    w => w.id === parsedEvent.warriorId,
                )
                if (joinedWarrior.id === $warrior.id) {
                    isSpectator = joinedWarrior.spectator
                }
                if ($warrior.notificationsEnabled) {
                    notifications.success(
                        `${$_('pages.battle.warriorJoined', {
                            values: { name: joinedWarrior.name },
                        })}`,
                    )
                }
                break
            }
            case 'warrior_retreated':
                const leftWarrior = battle.users.find(
                    w => w.id === parsedEvent.warriorId,
                )
                battle.users = JSON.parse(parsedEvent.value)

                if ($warrior.notificationsEnabled) {
                    notifications.danger(
                        `${$_('pages.battle.warriorRetreated', {
                            values: { name: leftWarrior.name },
                        })}`,
                    )
                }
                break
            case 'users_updated':
                battle.users = JSON.parse(parsedEvent.value)
                const updatedWarrior = battle.users.find(
                    w => w.id === $warrior.id,
                )
                isSpectator = updatedWarrior.spectator
                break
            case 'plan_added':
                battle.plans = JSON.parse(parsedEvent.value)
                break
            case 'plan_activated':
                const updatedPlans = JSON.parse(parsedEvent.value)
                const activePlan = updatedPlans.find(p => p.active)
                currentPlan = activePlan
                voteStartTime = new Date(activePlan.voteStartTime)

                battle.plans = updatedPlans
                battle.activePlanId = activePlan.id
                battle.votingLocked = false
                vote = ''
                break
            case 'plan_skipped':
                const updatedPlans2 = JSON.parse(parsedEvent.value)
                currentPlan = { ...defaultPlan }
                battle.plans = updatedPlans2
                battle.activePlanId = ''
                battle.votingLocked = true
                vote = ''
                if ($warrior.notificationsEnabled) {
                    notifications.warning($_('pages.battle.planSkipped'))
                }
                break
            case 'vote_activity':
                const votedWarrior = battle.users.find(
                    w => w.id === parsedEvent.warriorId,
                )
                if ($warrior.notificationsEnabled) {
                    notifications.success(
                        `${$_('pages.battle.warriorVoted', {
                            values: { name: votedWarrior.name },
                        })}`,
                    )
                }

                battle.plans = JSON.parse(parsedEvent.value)
                break
            case 'vote_retracted':
                const devotedWarrior = battle.users.find(
                    w => w.id === parsedEvent.warriorId,
                )
                if ($warrior.notificationsEnabled) {
                    notifications.warning(
                        `${$_('pages.battle.warriorRetractedVote', {
                            values: { name: devotedWarrior.name },
                        })}`,
                    )
                }

                battle.plans = JSON.parse(parsedEvent.value)
                break
            case 'voting_ended':
                battle.plans = JSON.parse(parsedEvent.value)
                battle.votingLocked = true
                break
            case 'plan_finalized':
                battle.plans = JSON.parse(parsedEvent.value)
                battle.activePlanId = ''
                currentPlan = { ...defaultPlan }
                vote = ''
                break
            case 'plan_revised':
                battle.plans = JSON.parse(parsedEvent.value)
                if (battle.activePlanId !== '') {
                    const activePlan = battle.plans.find(
                        p => p.id === battle.activePlanId,
                    )
                    currentPlan = activePlan
                }
                break
            case 'plan_burned':
                const postBurnPlans = JSON.parse(parsedEvent.value)

                if (
                    battle.activePlanId !== '' &&
                    postBurnPlans.filter(p => p.id === battle.activePlanId)
                        .length === 0
                ) {
                    battle.activePlanId = ''
                    currentPlan = { ...defaultPlan }
                }

                battle.plans = postBurnPlans

                break
            case 'leaders_updated':
                battle.leaders = parsedEvent.value
                break
            case 'battle_revised':
                const revisedBattle = JSON.parse(parsedEvent.value)
                battle.name = revisedBattle.battleName
                points = revisedBattle.pointValuesAllowed
                battle.autoFinishVoting = revisedBattle.autoFinishVoting
                battle.pointAverageRounding = revisedBattle.pointAverageRounding
                battle.joinCode = revisedBattle.joinCode
                break
            case 'battle_conceded':
                // battle over, goodbye.
                notifications.warning($_('pages.battle.battleDeleted'))
                router.route(appRoutes.battles)
                break
            case 'jab_warrior':
                const warriorToJab = battle.users.find(
                    w => w.id === parsedEvent.value,
                )
                notifications.info(
                    `${$_('pages.battle.warriorNudge', {
                        values: { name: warriorToJab.name },
                    })}`,
                )
                break
            default:
                break
        }
    }

    const ws = new Sockette(
        `${socketExtension}://${window.location.host}${PathPrefix}/api/arena/${battleId}`,
        {
            timeout: 2e3,
            maxAttempts: 15,
            onmessage: onSocketMessage,
            onerror: err => {
                socketError = true
                eventTag('socket_error', 'battle', '')
            },
            onclose: e => {
                if (e.code === 4004) {
                    eventTag('not_found', 'battle', '', () => {
                        router.route(appRoutes.battles)
                    })
                } else if (e.code === 4001) {
                    eventTag('socket_unauthorized', 'battle', '', () => {
                        warrior.delete()
                        router.route(`${appRoutes.register}/battle/${battleId}`)
                    })
                } else if (e.code === 4003) {
                    eventTag('socket_duplicate', 'battle', '', () => {
                        notifications.danger($_('sessionDuplicate'))
                        router.route(`${appRoutes.battles}`)
                    })
                } else if (e.code === 4002) {
                    eventTag('battle_warrior_abandoned', 'battle', '', () => {
                        router.route(appRoutes.battles)
                    })
                } else {
                    socketReconnecting = true
                    eventTag('socket_close', 'battle', '')
                }
            },
            onopen: () => {
                socketError = false
                socketReconnecting = false
                eventTag('socket_open', 'battle', '')
            },
            onmaximum: () => {
                socketReconnecting = false
                eventTag(
                    'socket_error',
                    'battle',
                    'Socket Reconnect Max Reached',
                )
            },
        },
    )

    onDestroy(() => {
        eventTag('leave', 'battle', '', () => {
            ws.close()
        })
    })

    const sendSocketEvent = (type, value) => {
        ws.send(
            JSON.stringify({
                type,
                value,
            }),
        )
    }

    const handleVote = event => {
        vote = event.detail.point
        const voteValue = {
            planId: battle.activePlanId,
            voteValue: vote,
            autoFinishVoting: battle.autoFinishVoting,
        }

        sendSocketEvent('vote', JSON.stringify(voteValue))
        eventTag('vote', 'battle', vote)
    }

    const handleUnvote = () => {
        vote = ''

        sendSocketEvent('retract_vote', battle.activePlanId)
        eventTag('retract_vote', 'battle', vote)
    }

    // Determine if the warrior has voted on active Plan yet
    function didVote(warriorId) {
        if (battle.activePlanId === '') {
            return false
        }
        const plan = battle.plans.find(p => p.id === battle.activePlanId)
        const voted = plan.votes.find(w => w.warriorId === warriorId)

        return voted !== undefined
    }

    // Determine if we are showing warriors vote
    function showVote(warriorId) {
        if (battle.activePlanId === '' || battle.votingLocked === false) {
            return ''
        }
        const plan = battle.plans.find(p => p.id === battle.activePlanId)
        const voted = plan.votes.find(w => w.warriorId === warriorId)

        return voted !== undefined ? voted.vote : ''
    }

    // get highest vote from active plan
    function getHighestVote() {
        const voteCounts = {}
        points.forEach(p => {
            voteCounts[p] = 0
        })
        const highestVote = {
            vote: '',
            count: 0,
        }
        const activePlan = battle.plans.find(p => p.id === battle.activePlanId)

        if (activePlan.votes.length > 0) {
            const reversedPoints = [...points].filter(v => v !== '?').reverse()
            reversedPoints.push('?')

            // build a count of each vote
            activePlan.votes.forEach(v => {
                const { spectator = false } = battle.users.find(
                    w => w.id === v.warriorId,
                )
                if (typeof voteCounts[v.vote] !== 'undefined' && !spectator) {
                    ++voteCounts[v.vote]
                }
            })

            // find the highest vote giving priority to higher numbers
            reversedPoints.forEach(p => {
                if (voteCounts[p] > highestVote.count) {
                    highestVote.vote = p
                    highestVote.count = voteCounts[p]
                }
            })
        }

        return highestVote.vote
    }

    $: highestVoteCount =
        battle.activePlanId !== '' && battle.votingLocked === true
            ? getHighestVote()
            : ''
    $: showVotingResults =
        battle.activePlanId !== '' && battle.votingLocked === true

    $: isLeader = battle.leaders.includes($warrior.id)

    function concedeBattle() {
        eventTag('concede_battle', 'battle', '', () => {
            sendSocketEvent('concede_battle', '')
        })
    }

    function abandonBattle() {
        eventTag('abandon_battle', 'battle', '', () => {
            sendSocketEvent('abandon_battle', '')
        })
    }

    function toggleEditBattle() {
        showEditBattle = !showEditBattle
    }

    const toggleDeleteBattle = () => {
        showDeleteBattle = !showDeleteBattle
    }

    function handleBattleEdit(revisedBattle) {
        sendSocketEvent('revise_battle', JSON.stringify(revisedBattle))
        eventTag('revise_battle', 'battle', '')
        toggleEditBattle()
        battle.leaderCode = revisedBattle.leaderCode
    }

    function authBattle(e) {
        e.preventDefault()

        sendSocketEvent('auth_battle', joinPasscode)
        eventTag('auth_battle', 'battle', '')
    }

    function timeUnitsBetween(startDate, endDate) {
        let delta = Math.abs(endDate - startDate) / 1000
        return [
            ['days', 24 * 60 * 60],
            ['hours', 60 * 60],
            ['minutes', 60],
            ['seconds', 1],
        ].reduce(
            (acc, [key, value]) => (
                (acc[key] = Math.floor(delta / value)),
                (delta -= acc[key] * value),
                acc
            ),
            {},
        )
    }

    function addTimeLeadZero(time) {
        return ('0' + time).slice(-2)
    }

    onMount(() => {
        if (!$warrior.id) {
            router.route(`${loginOrRegister}/battle/${battleId}`)
            return
        }
        const voteCounter = setInterval(() => {
            currentTime = new Date()
        }, 1000)

        return () => {
            clearInterval(voteCounter)
        }
    })
</script>

<svelte:head>
    <title>{$_('pages.battle.title')} {battle.name} | {$_('appName')}</title>
</svelte:head>

<PageLayout>
    {#if battle.name && !socketReconnecting && !socketError}
        <div class="mb-6 flex flex-wrap">
            <div class="w-full text-center md:w-2/3 md:text-left">
                <h1
                    class="text-4xl font-semibold font-rajdhani leading-tight dark:text-white flex items-center"
                >
                    {#if currentPlan.link}
                        <a
                            href="{currentPlan.link}"
                            target="_blank"
                            class="text-blue-800 dark:text-sky-400 inline-block"
                            data-testid="currentplan-link"
                        >
                            <ExternalLinkIcon class="w-8 h-8" />
                        </a>
                    {/if}
                    {#if currentPlan.type}
                        &nbsp;<span
                            class="inline-block text-lg text-gray-500
                            border-gray-300 border px-1 rounded dark:text-gray-300 dark:border-gray-500"
                            data-testid="currentplan-type"
                        >
                            {currentPlan.type}
                        </span>
                    {/if}
                    {#if currentPlan.referenceId}
                        &nbsp;<span data-testid="currentplan-refid"
                            >[{currentPlan.referenceId}]</span
                        >
                    {/if}
                    &nbsp;<span data-testid="currentplan-name"
                        >{currentPlan.name}</span
                    >
                </h1>
                <h2
                    class="text-gray-700 dark:text-gray-300 text-3xl font-semibold font-rajdhani leading-tight"
                    data-testid="battle-name"
                >
                    {battle.name}
                </h2>
            </div>
            <div
                class="w-full md:w-1/3 text-center md:text-right font-semibold
                text-3xl md:text-4xl text-gray-700 dark:text-gray-300"
                data-testid="vote-timer"
            >
                {#if countdown.seconds !== undefined}
                    {#if countdown.hours !== 0}
                        {addTimeLeadZero(countdown.hours)}:
                    {/if}
                    {addTimeLeadZero(countdown.minutes)}:{addTimeLeadZero(
                        countdown.seconds,
                    )}
                {/if}
            </div>
        </div>

        <div class="flex flex-wrap mb-4 -mx-4">
            <div class="w-full lg:w-3/4 px-4">
                {#if showVotingResults}
                    <VoteResults
                        warriors="{battle.users}"
                        plans="{battle.plans}"
                        activePlanId="{battle.activePlanId}"
                        points="{points}"
                        highestVote="{highestVoteCount}"
                        averageRounding="{battle.pointAverageRounding}"
                    />
                {:else}
                    <div class="flex flex-wrap mb-4 -mx-2 mb-4 lg:mb-6">
                        {#each points as point}
                            <div class="w-1/4 md:w-1/6 px-2 mb-4">
                                <PointCard
                                    point="{point}"
                                    active="{vote === point}"
                                    on:voted="{handleVote}"
                                    on:voteRetraction="{handleUnvote}"
                                    isLocked="{battle.votingLocked ||
                                        isSpectator}"
                                />
                            </div>
                        {/each}
                    </div>
                {/if}

                <BattlePlans
                    plans="{battle.plans}"
                    isLeader="{isLeader}"
                    sendSocketEvent="{sendSocketEvent}"
                    eventTag="{eventTag}"
                    notifications="{notifications}"
                />
            </div>

            <div class="w-full lg:w-1/4 px-4">
                <div
                    class="bg-white dark:bg-gray-800 shadow-lg mb-4 rounded-lg"
                >
                    <div class="bg-blue-500 dark:bg-gray-700 p-4 rounded-t-lg">
                        <h3
                            class="text-3xl text-white leading-tight font-semibold font-rajdhani uppercase"
                        >
                            {$_('pages.battle.warriors')}
                        </h3>
                    </div>

                    {#each battle.users as war (war.id)}
                        {#if war.active}
                            <WarriorCard
                                warrior="{war}"
                                leaders="{battle.leaders}"
                                isLeader="{isLeader}"
                                voted="{didVote(war.id)}"
                                points="{showVote(war.id)}"
                                autoFinishVoting="{battle.autoFinishVoting}"
                                sendSocketEvent="{sendSocketEvent}"
                                eventTag="{eventTag}"
                            />
                        {/if}
                    {/each}

                    {#if isLeader}
                        <VotingControls
                            points="{points}"
                            planId="{battle.activePlanId}"
                            sendSocketEvent="{sendSocketEvent}"
                            votingLocked="{battle.votingLocked}"
                            highestVote="{highestVoteCount}"
                            eventTag="{eventTag}"
                        />
                    {/if}
                </div>

                <div
                    class="bg-white dark:bg-gray-800 shadow-lg p-4 mb-4 rounded-lg"
                >
                    <InviteWarrior
                        hostname="{hostname}"
                        battleId="{battle.id}"
                        joinCode="{battle.joinCode}"
                        notifications="{notifications}"
                    />
                    {#if isLeader}
                        <div class="mt-4 text-right">
                            <HollowButton
                                color="blue"
                                onClick="{toggleEditBattle}"
                                testid="battle-edit"
                            >
                                {$_('battleEdit')}
                            </HollowButton>
                            <HollowButton
                                color="red"
                                onClick="{toggleDeleteBattle}"
                                testid="battle-delete"
                            >
                                {$_('battleDelete')}
                            </HollowButton>
                        </div>
                    {:else}
                        <div class="mt-4 text-right">
                            <HollowButton
                                color="red"
                                onClick="{abandonBattle}"
                                testid="battle-abandon"
                            >
                                {$_('battleAbandon')}
                            </HollowButton>
                        </div>
                    {/if}
                </div>
            </div>
        </div>
        {#if showEditBattle}
            <EditBattle
                battleName="{battle.name}"
                points="{points}"
                votingLocked="{battle.votingLocked}"
                autoFinishVoting="{battle.autoFinishVoting}"
                pointAverageRounding="{battle.pointAverageRounding}"
                handleBattleEdit="{handleBattleEdit}"
                toggleEditBattle="{toggleEditBattle}"
                joinCode="{battle.joinCode}"
                leaderCode="{battle.leaderCode}"
            />
        {/if}
    {:else if JoinPassRequired}
        <div class="flex justify-center">
            <div class="w-full md:w-1/2 lg:w-1/3">
                <form
                    on:submit="{authBattle}"
                    class="bg-white dark:bg-gray-800 shadow-lg rounded-lg p-6 mb-4"
                    name="authBattle"
                >
                    <div class="mb-4">
                        <label
                            class="block text-gray-700 dark:text-gray-400 font-bold mb-2"
                            for="battleJoinCode"
                        >
                            {$_('passCodeRequired')}
                        </label>
                        <input
                            bind:value="{joinPasscode}"
                            placeholder="{$_('enterPasscode')}"
                            class="bg-gray-100 dark:bg-gray-900 border-gray-200 dark:border-gray-800 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 dark:text-gray-300 leading-tight
                focus:outline-none focus:bg-white dark:focus:bg-gray-700 focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
                            id="battleJoinCode"
                            name="battleJoinCode"
                            type="password"
                            required
                        />
                    </div>

                    <div class="text-right">
                        <SolidButton type="submit"
                            >{$_('battleJoin')}</SolidButton
                        >
                    </div>
                </form>
            </div>
        </div>
    {:else if socketReconnecting}
        <div class="flex items-center">
            <div class="flex-1 text-center">
                <h1 class="text-5xl text-teal-500 leading-tight font-bold">
                    {$_('pages.battle.socketReconnecting')}
                </h1>
            </div>
        </div>
    {:else if socketError}
        <div class="flex items-center">
            <div class="flex-1 text-center">
                <h1 class="text-5xl text-red-500 leading-tight font-bold">
                    {$_('pages.battle.socketError')}
                </h1>
            </div>
        </div>
    {:else}
        <div class="flex items-center">
            <div class="flex-1 text-center">
                <h1 class="text-5xl text-green-500 leading-tight font-bold">
                    {$_('pages.battle.loading')}
                </h1>
            </div>
        </div>
    {/if}

    {#if showDeleteBattle}
        <DeleteConfirmation
            toggleDelete="{toggleDeleteBattle}"
            handleDelete="{concedeBattle}"
            confirmText="{$_('deleteBattleConfirmText')}"
            confirmBtnText="{$_('deleteBattle')}"
        />
    {/if}
</PageLayout>
