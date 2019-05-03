<script>
    import { onDestroy } from 'svelte'

    export let battleId = 0

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

    function sendMessage(e) {
        e.preventDefault()
        console.log("SEND: " + message)
            
        ws.send(message)
    }

    onDestroy(() => ws.close())
</script>

<div class="columns">
    <div class="column">
        <h1>Battle: {battleId}</h1>
        <p>
            Click "Send" to send a message to the server.<br />
            You can change the message and send multiple times.
        </p>

        <form on:submit={sendMessage}>
            <div class="field">
                <label class="label">Message</label>
                <div class="control">
                    <input bind:value={message} placeholder="Enter a message" class="input" required />
                </div>
            </div>

            <div class="field">
                <div class="control">
                    <button class="button is-success" type="submit">Send</button>
                </div>
            </div>
        </form>
    </div>
    <div class="column">            
        <h2>Responses:</h2>

        {#each responses as response}
            <div class="notification is-success">
                {response}
            </div>
        {/each}
    </div>
  </div>