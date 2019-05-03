<script>
    import { onDestroy } from 'svelte'

    export let battleId = 0

    const ws = new WebSocket(`ws://${window.location.host}/api/battle/${battleId}`)

    let message = ''
    let responseOutput = ''

    ws.onmessage = function (evt) {
        responseOutput = `RESPONSE: ${evt.data}`
    }

    ws.onerror = function (e) {
        responseOutput = `ERROR: ${e}`
    }

    function sendMessage(e) {
        e.preventDefault()
        console.log("SEND: " + message)
            
        ws.send(message)
    }

    onDestroy(() => ws.close())
</script>

<style>
    p  {
        font-size: 14px;
    }
</style>

<div class="columns">
    <div class="column">
        <h1>Battle: {battleId}</h1>
        <p>
            Click "Send" to send a message to the server.<br />
            You can change the message and send multiple times.
        </p>

        <form on:submit={sendMessage}>
            <div>
                <input bind:value={message} placeholder="Enter a message" />
            </div>

            <button class="button is-success" type="submit">Send</button>
        </form>
    </div>
    <div class="column">            
        <p>{responseOutput}</p>
    </div>
  </div>