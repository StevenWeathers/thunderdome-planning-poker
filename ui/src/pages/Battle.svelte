<script>
    import { onDestroy } from 'svelte'

    import PointCard from '../components/PointCard.svelte'
    import WarriorCard from '../components/WarriorCard.svelte'
    import { warrior } from '../stores.js'

    export let battleId = 0
    const socketExtension = window.location.protocol === 'https:' ? 'wss' : 'ws'
    let points = ['1', '2', '3', '5', '8', '13', '?']
    let vote = ''
    let battle = {
        name: '',
        leaderId: '',
        warriors: [],
        votes: []
    }
    let responses = []
    let handleVote = () => {}

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

            const ws = new WebSocket(`${socketExtension}://${window.location.host}/api/arena/${battleId}`)

            ws.onmessage = function (evt) {
                const parsedEvent = JSON.parse(evt.data)
                const warriorId = parsedEvent.id
                let eventWarrior = battle.warriors.find(w => w.id === warriorId)
                let response = ''

                switch(parsedEvent.type) {
                    case "joined":
                        const joinedWarrior = {
                            name: parsedEvent.value,
                            id: warriorId
                        }
                        if (!eventWarrior) {
                            battle.warriors[battle.warriors.length] = joinedWarrior
                        }                
                        
                        response = `${joinedWarrior.name} has joined the battle.`
                        break;
                    case "retreat":
                        battle.warriors = battle.warriors.filter(w => w.id !== warriorId)
                        
                        response = `${eventWarrior.name} has retreated from battle.`
                        break;
                    case "vote":
                        const vote = parsedEvent.value
                        const currentVote = battle.votes.find(v => v.warriorId === warriorId)
                        if (currentVote) {
                            currentVote.vote = vote
                        } else {
                            battle.votes[battle.votes.length] = {
                                warriorId,
                                vote
                            }
                        }

                        response = `${eventWarrior.name} voted ${parsedEvent.value}.`
                    default:
                        break;
                }

                responses[responses.length] = response;
            }

            ws.onerror = function (e) {
                console.log(`ERROR: ${e}`)
            }

            handleVote = (event) => {
                vote = event.detail.point

                ws.send(JSON.stringify({
                    type: 'vote',
                    id: $warrior.id,
                    value: vote
                }))
            }
        })
        .catch(function(e) {
            console.log(e.message)
            window.location.href = '/'
        })
</script>

<svelte:head>
    <title>Battle {battle.name} | Thunderdome</title>
</svelte:head>

{#if battle.name}
    <h1 class="title">Battle: {battle.name}</h1>

    <div class="columns">
        <div class="column is-three-quarters">
            <div class="columns">
                {#each points as point}
                    <div class="column">
                        <PointCard point={point} active={vote === point} on:voted={handleVote} />
                    </div>
                {/each}
            </div>
            
            <h3 class="is-size-2">Votes:</h3>
            {#each responses as response}
                <div class="notification is-success">
                    {response}
                </div>
            {/each}
        </div>
        
        <div class="column">
            <h3 class="is-size-2">Users</h3>

            {#each battle.warriors as war (war.id)}
                <WarriorCard warrior={war} isLeader={war.id === battle.leaderId} voted={battle.votes.find(v => v.warriorId === war.id)} />
            {/each}
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