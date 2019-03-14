<template>
  <div class="columns">
    <div class="column">
        <h1>Battle: {{ $route.params.id }}</h1>
        <p>Click "Send" to send a message to the server.<br />
        You can change the message and send multiple times.
        </p>

        <form>
            <b-field label="Name">
                <b-input v-model="name" id="input"></b-input>
            </b-field>
            <button class="button is-success" v-on:click="sendMessage">Send</button>
        </form>
    </div>
    <div class="column">            
        <div id="output"></div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Battle',
  data() {
    return {
        name: 'Thunderdome',
        ws: null
    }
  },
  created: function () {
      this.ws = new WebSocket("ws://localhost:8080/api/battle/1");

      this.ws.onmessage = function (evt) {
        console.log("RESPONSE: " + evt.data);
      }
      this.ws.onerror = function (evt) {
        console.log("ERROR: " + evt.data);
      }
  },
  beforeDestroy: function () {
      this.ws.close()
  },
  methods: {
      sendMessage: function (event) {
        event.preventDefault()

        console.log("SEND: " + input.value);
        this.ws.send(input.value);
      }
  }
}
</script>