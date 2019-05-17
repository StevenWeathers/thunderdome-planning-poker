<script>
    import { onDestroy } from 'svelte'

    import PointCard from '../components/PointCard.svelte'
    import WarriorCard from '../components/WarriorCard.svelte'
    import BattlePlans from '../components/BattlePlans.svelte'
    import VotingControls from '../components/VotingControls.svelte'
    import InviteWarrior from '../components/InviteWarrior.svelte'

    import { warrior } from '../stores.js'

    export let battleId = 0
    export let notifications

    const hostname = window.location.origin
    const socketExtension = window.location.protocol === 'https:' ? 'wss' : 'ws'
    
    let socketError = false
    let points = ['0', '1/2', '1', '2', '3', '5', '8', '13', '20', '40', '100', '?']
    let vote = ''
    let battle = {}
    let currentPlanName = '[Voting not started]'
    
    let ws = {
        send: () => {},
        close: () => {}
    }

    fetch(`/api/battle/${battleId}`)
        .then(function(response) {
            if (response.status >= 200 && response.status < 300) {
                return response.json()
            } else {
                const error = new Error(response.statusText || `${response.status}`)

                return Promise.reject(error)
            }
        })
        .then(function(b) {
            battle = b
            if (battle.activePlanId !== '') {
                currentPlanName = battle.plans.find(p => p.id === battle.activePlanId).name
            }

            ws = new WebSocket(`${socketExtension}://${window.location.host}/api/arena/${battleId}`)

            ws.onmessage = function (evt) {
                const parsedEvent = JSON.parse(evt.data)

                switch(parsedEvent.type) {
                    case "user_activity":
                        battle.warriors = JSON.parse(parsedEvent.value)
                        break;
                    case "plan_added":
                        battle.plans = JSON.parse(parsedEvent.value)
                        break;
                    case "plan_activated":
                        const updatedPlans = JSON.parse(parsedEvent.value)
                        const activePlan = updatedPlans.find(p => p.active)
                        currentPlanName = activePlan.name
                        battle.plans = updatedPlans                        
                        battle.activePlanId = activePlan.id
                        battle.votingLocked = false
                        vote = ''
                        break;
                    case "plan_skipped":
                        const updatedPlans2 = JSON.parse(parsedEvent.value)
                        currentPlanName = '[Voting not started]'
                        battle.plans = updatedPlans2                        
                        battle.activePlanId = ''
                        battle.votingLocked = true
                        vote = ''
                        break;
                    case "vote_activity":
                        battle.plans = JSON.parse(parsedEvent.value)
                        break;
                    case "voting_ended":
                        battle.plans = JSON.parse(parsedEvent.value)
                        battle.votingLocked = true
                        break;
                    case "plan_finalized":
                        battle.plans = JSON.parse(parsedEvent.value)
                        battle.activePlanId = ''
                        currentPlanName = '[Voting not started]'
                        vote = ''
                        break;
                    case "plan_revised":
                        battle.plans = JSON.parse(parsedEvent.value)
                        break;
                    case "plan_burned":
                        const postBurnPlans = JSON.parse(parsedEvent.value)

                        if (battle.activePlanId !== '' && postBurnPlans.filter(p => p.id === battle.activePlanId).length === 0) {
                            battle.activePlanId = ''
                            currentPlanName = '[Voting not started]'
                        }

                        battle.plans = postBurnPlans

                        break;
                    case "battle_updated":
                        battle = JSON.parse(parsedEvent.value)
                        break;
                    case "battle_conceded":
                        // battle over, goodbye.
                        window.location.href = '/'
                        break;
                    default:
                        break;
                }
            }

            ws.onerror = function (err) {
                // @TODO - add toast or some visual error notifications...
                console.log(`ERROR: ${JSON.stringify(err, ["message", "arguments", "type", "name"])}`)
                socketError = true
            }
        })
        .catch(function(err) {
            // battle not found or server issue, redirect to landing
            window.location.href = '/'
        })

    onDestroy(() => {
        ws.close();
    })

    const sendSocketEvent = (type, value) => {
        ws.send(JSON.stringify({
            type,
            value
        }))
    }

    const handleVote = (event) => {
        vote = event.detail.point
        const voteValue = {
            planId: battle.activePlanId,
            voteValue: vote
        }

        sendSocketEvent('vote', JSON.stringify(voteValue))
    }

    // Determine if the warrior has voted on active Plan yet
    function didVote(warriorId) {
        if (battle.activePlanId === "") {
            return false
        }
        const currentPlan = battle.plans.find(p => p.id === battle.activePlanId)
        const voted = currentPlan.votes.find(w => w.warriorId === warriorId)

        return voted !== undefined
    }

    // Determine if we are showing warriors vote
    function showVote(warriorId) {
        if (battle.activePlanId === "" || battle.votingLocked === false) {
            return ''
        }
        const currentPlan = battle.plans.find(p => p.id === battle.activePlanId)
        const voted = currentPlan.votes.find(w => w.warriorId === warriorId)

        return voted !== undefined ? voted.vote : '' 
    }

    function concedeBattle() {
        sendSocketEvent("concede_battle", "")
    }
</script>

<svelte:head>
    <title>Battle {battle.name} | Thunderdome</title>
</svelte:head>

{#if battle.name && !socketError}
    <div class="mb-6">
        <h1>{currentPlanName}</h1>
        <h2 class="text-grey-darker">{battle.name}</h2>
    </div>

    <div class="flex flex-wrap mb-4 -mx-4">
        <div class="w-full lg:w-3/4 px-4">
            <div class="flex flex-wrap mb-4 -mx-2 mb-4 lg:mb-6">
                {#each points as point}
                    <div class="w-1/4 md:w-1/6 px-2 mb-4">
                        <PointCard point={point} active={vote === point} on:voted={handleVote} isLocked={battle.votingLocked} />
                    </div>
                {/each}
            </div>

            <BattlePlans
                plans={battle.plans}
                isLeader={battle.leaderId === $warrior.id}
                sendSocketEvent={sendSocketEvent}
            />
        </div>
        
        <div class="w-full lg:w-1/4 px-4">
            <div class="bg-white shadow-md mb-4 rounded">
                <div class="bg-blue p-4 rounded-t">
                    <h3 class="text-2xl text-white">Warriors</h3>
                </div>

                {#each battle.warriors as war (war.id)}
                    <WarriorCard warrior={war} leaderId={battle.leaderId} isLeader={battle.leaderId === $warrior.id} voted={didVote(war.id)} points={showVote(war.id)} sendSocketEvent={sendSocketEvent} />
                {/each}

                {#if battle.leaderId === $warrior.id}
                    <VotingControls points={points} planId={battle.activePlanId} sendSocketEvent={sendSocketEvent} votingLocked={battle.votingLocked} />
                {/if}
            </div>

            <div class="bg-white shadow-md p-5 mb-4 rounded">
                <InviteWarrior hostname={hostname} battleId={battle.id} />
                {#if battle.leaderId === $warrior.id}
                    <div class="mt-4 text-right">
                        <button
                            class="bg-transparent hover:bg-red text-red-dark font-semibold hover:text-white py-2 px-2 border border-red hover:border-transparent rounded"
                            on:click={concedeBattle}
                        >
                            Delete Battle
                        </button>
                    </div>
                {/if}
            </div>
        </div>
    </div>
{:else if socketError}
    <div class="flex items-center">
        <div class="flex-1 text-center">
            <h1 class="text-5xl text-red">Error joining battle, refresh and try again.</h1>
        </div>
    </div>
{:else}
    <div class="flex items-center">
        <div class="flex-1 text-center">
            <h1 class="text-5xl text-teal">Loading Battle Plans...</h1>
        </div>
    </div>
{/if}