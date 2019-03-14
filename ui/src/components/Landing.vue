<template>
  <div class="columns">
    <div class="column">
        <form v-on:submit="createBattle">
            <b-field label="Your Name">
                <b-input v-model="creatorName" placeholder="Enter your name" required></b-input>
            </b-field>
            <b-field label="Battle Name">
                <b-input v-model="battleName" placeholder="Enter a battle name" required></b-input>
            </b-field>
            <button class="button is-primary">Create a Story Battle</button>
        </form>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Landing',
  methods: {
      createBattle: function (event) {
        event.preventDefault()
        const router = this.$router

        const {
            creatorName,
            battleName
        } = this

        fetch('/api/battle', {
            method: 'POST',
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                creatorName,
                battleName
            })
        })
            .then(function(response) {
                return response.json()
            })
            .then(function(battle) {
                router.push(`/battle/${battle.id}`)
            });
        
    }
  },
  data() {
        return {
            creatorName: '',
            battleName: ''
        }
    }
}
</script>