<template>
  <div class="columns">
    <div class="column">
        <p>Click "Send" to send a message to the server.<br />
        You can change the message and send multiple times.
        </p>

        <form>
            <b-field label="Name">
                <b-input v-model="name" id="input"></b-input>
            </b-field>
            <button id="send" class="button is-success">Send</button>
        </form>
    </div>
    <div class="column">            
        <div id="output"></div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Landing',
  data() {
        return {
            name: 'Thunderdome'
        }
    }
}
</script>