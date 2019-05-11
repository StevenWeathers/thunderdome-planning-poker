<script>
    import { onDestroy } from 'svelte'

    import PointCard from '../components/PointCard.svelte'
    import WarriorCard from '../components/WarriorCard.svelte'
    import BattlePlans from '../components/BattlePlans.svelte'

    import { warrior } from '../stores.js'

    export let battleId = 0
    const socketExtension = window.location.protocol === 'https:' ? 'wss' : 'ws'
    let points = ['1', '2', '3', '5', '8', '13', '?']
    let vote = ''
    let battle = {}
    let currentPlanName = '[Voting not started]'
    
    let ws = {
        send: () => {}
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
                currentPlanName = battle.plans.find(p => p.active).name
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
                        break;
                    case "vote_activity":
                        battle.plans = JSON.parse(parsedEvent.value)
                        break;
                    case "voting_ended":
                        battle.plans = JSON.parse(parsedEvent.value)
                        battle.votingLocked = true
                        break;
                    case "plan_burned":
                        const postBurnPlans = JSON.parse(parsedEvent.value)

                        if (battle.activePlanId !== '' && postBurnPlans.filter(p => p.id === battle.activePlanId).length === 0) {
                            battle.activePlanId = ''
                        }

                        battle.plans = postBurnPlans

                        break;
                    default:
                        break;
                }
            }

            ws.onerror = function (e) {
                // @TODO - add toast or some visual error notifications...
                console.log(`ERROR: ${e}`)
            }
        })
        .catch(function(e) {
            // battle not found or server issue, redirect to landing
            window.location.href = '/'
        })

    const sendSocketEvent = (type, value) => {
        ws.send(JSON.stringify({
            type,
            id: $warrior.id,
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

    const endPlanVoting = () => {
        sendSocketEvent('end_voting', battle.activePlanId)
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
</script>

<svelte:head>
    <title>Battle {battle.name} | Thunderdome</title>
</svelte:head>

{#if battle.name}
    <h1 class="title">{currentPlanName}</h1>
    <h2 class="subtitle">{battle.name}</h2>

    <div class="columns">
        <div class="column is-three-quarters">
            <div class="columns">
                {#each points as point}
                    <div class="column">
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
        
        <div class="column">
            <h3 class="is-size-3">Users</h3>

            {#each battle.warriors as war (war.id)}
                <WarriorCard warrior={war} isLeader={war.id === battle.leaderId} voted={didVote(war.id)} />
            {/each}

            {#if battle.leaderId === $warrior.id}
                <div>
                    <button class="button is-link is-outlined is-fullwidth" on:click={endPlanVoting} disabled={battle.votingLocked}>End Voting</button>
                </div>
            {/if}
        </div>
    </div>
{:else}
    <div class="columns">
        <div class="column is-half is-offset-one-quarter has-text-centered">
            <h1 class="is-size-1 has-text-primary">Loading Battle Plans...</h1>
            <progress class="progress is-large is-info" max="100">60%</progress>
        </div>
    </div>
{/if}