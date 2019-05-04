<script>
    import { onDestroy } from 'svelte'
    import PointCard from '../components/PointCard.svelte'

    export let battleId = 0
    let points = ['1', '2', '3', '5', '8', '13', '?']
    let vote = '3'

    const socketExtension = window.location.protocol === 'https:' ? 'wss' : 'ws'

    const ws = new WebSocket(`${socketExtension}://${window.location.host}/api/battle/${battleId}`)

    let message = ''
    let responses = []

    ws.onmessage = function (evt) {
        responses[responses.length] = `${evt.data}`
    }

    ws.onerror = function (e) {
        console.log(`ERROR: ${e}`)
    }

    onDestroy(() => ws.close())

    function handleVote(event) {
        vote = event.detail.point
        ws.send(vote)
    }
</script>

<div class="columns">
    <div class="column">
        <h1 class="title">Battle: {battleId}</h1>

        <div class="columns">
            {#each points as point}
                <div class="column">
                    <PointCard point={point} active={vote === point} on:voted={handleVote} />
                </div>
            {/each}
        </div>
    </div>
    <div class="column">            
        <h2>Votes:</h2>

        {#each responses as response}
            <div class="notification is-success">
                {response}
            </div>
        {/each}
    </div>
  </div>