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

<form on:submit={createBattle} class="bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4">
    <div class="mb-4">
        <label class="block text-grey-darker text-sm font-bold mb-2" for="battleName">Battle Name</label>
        <div class="control">
            <input
                bind:value={battleName}
                placeholder="Enter a battle name"
                class="shadow appearance-none border rounded w-full py-2 px-3 text-grey-darker leading-tight focus:outline-none focus:shadow-outline"
                id="battleName"
                required
            />
        </div>
    </div>
    
    <div class="mb-6">
        <div class="flex items-center justify-between">
            <button class="bg-green hover:bg-green-dark text-white font-bold py-2 px-4 rounded" type="submit">Create a Story Battle</button>
        </div>
    </div>
</form>
