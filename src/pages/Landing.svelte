<script>
    import { warrior } from '../stores.js'

    let battleName = ''

    function createBattle(e) {
        e.preventDefault()
        const data = {
            battleName,
            leaderId: $warrior.id
        }
        
        fetch('/api/battle', {
            method: 'POST',
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(data)
        })
            .then(function(response) {
                return response.json()
            })
            .then(function(battle) {
                window.location.href = `/battle/${battle.id}`
            });
    }
</script>

<div class="columns">
    <div class="column">
        <form on:submit={createBattle}>
            <div class="field">
                <label class="label">Battle Name</label>
                <div class="control">
                    <input bind:value={battleName} placeholder="Enter a battle name" class="input" required />
                </div>
            </div>
            
            <div class="field">
                <div class="control">
                    <button class="button is-success" type="submit">Create a Story Battle</button>
                </div>
            </div>
        </form>
    </div>
</div>
