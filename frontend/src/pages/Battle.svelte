<script>
    import Sockette from 'sockette'
    import { onDestroy, onMount } from 'svelte'

    import PageLayout from '../components/PageLayout.svelte'
    import PointCard from '../components/PointCard.svelte'
    import WarriorCard from '../components/WarriorCard.svelte'
    import BattlePlans from '../components/BattlePlans.svelte'
    import VotingControls from '../components/VotingControls.svelte'
    import InviteWarrior from '../components/InviteWarrior.svelte'
    import VoteResults from '../components/VoteResults.svelte'
    import HollowButton from '../components/HollowButton.svelte'
    import ExternalLinkIcon from '../components/icons/ExternalLinkIcon.svelte'
    import EditBattle from '../components/EditBattle.svelte'
    import DeleteBattle from '../components/DeleteBattle.svelte'
    import { warrior } from '../stores.js'
    import { _ } from '../i18n'
    import { appRoutes, PathPrefix } from '../config'

    export let battleId
    export let notifications
    export let eventTag
    export let router

    const { AllowRegistration } = appConfig
    const loginOrRegister = AllowRegistration
        ? appRoutes.register
        : appRoutes.login

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

    $: countdown =
        battle.currentPlanId !== '' && battle.votingLocked === false
            ? timeUnitsBetween(voteStartTime, currentTime)
            : {}

    const onSocketMessage = function (evt) {
        const parsedEvent = JSON.parse(evt.data)

        switch (parsedEvent.type) {
            case 'init': {
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
                        router.route(`${appRoutes.register}/${battleId}`)
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

    // get hightest vote from active plan
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
            router.route(`${loginOrRegister}/${battleId}`)
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
                <h1 class="text-3xl font-bold leading-tight">
                    {#if currentPlan.link}
                        <a
                            href="{currentPlan.link}"
                            target="_blank"
                            class="text-blue-800"
                        >
                            <ExternalLinkIcon />
                        </a>
                        &nbsp;
                    {/if}
                    {#if currentPlan.type}
                        <span
                            class="inline-block text-lg text-gray-500
                            border-gray-400 border px-1 rounded"
                            data-testId="battlePlanType"
                        >
                            {currentPlan.type}
                        </span>
                        &nbsp;
                    {/if}
                    {#if currentPlan.referenceId}
                        [{currentPlan.referenceId}]&nbsp;
                    {/if}
                    {currentPlan.name}
                </h1>
                <h2 class="text-gray-700 text-2xl font-bold leading-tight">
                    {battle.name}
                </h2>
            </div>
            <div
                class="w-full md:w-1/3 text-center md:text-right font-semibold
                text-3xl md:text-4xl text-gray-700"
                data-testId="votingTimer"
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
                <div class="bg-white shadow-lg mb-4 rounded">
                    <div class="bg-blue-500 p-4 rounded-t">
                        <h3 class="text-2xl text-white leading-tight font-bold">
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

                <div class="bg-white shadow-lg p-4 mb-4 rounded">
                    <InviteWarrior
                        hostname="{hostname}"
                        battleId="{battle.id}"
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
            />
        {/if}
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
        <DeleteBattle
            toggleDelete="{toggleDeleteBattle}"
            handleDelete="{concedeBattle}"
        />
    {/if}
</PageLayout>
