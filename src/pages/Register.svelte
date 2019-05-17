<script>
    import { warrior } from '../stores.js'

    export let notifications
    let warriorName = ''

    function createWarrior(e) {
        e.preventDefault()
        
        fetch('/api/warrior', {
            method: 'POST',
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                warriorName
            })
        })
            .then(function(response) {
                return response.json()
            })
            .then(function(newWarrior) {
                warrior.create({
                    id: newWarrior.id,
                    name: newWarrior.name
                })
            }).catch(function(error) {
                notifications.danger("Error encountered registering warrior")
            })
    }
</script>

<form on:submit={createWarrior} class="bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4">
    <div class="mb-4">
        <label class="block text-grey-darker text-sm font-bold mb-2" for="yourName">Name</label>
        <input
            bind:value={warriorName}
            placeholder="Enter your name"
            class="shadow appearance-none border rounded w-full py-2 px-3 text-grey-darker leading-tight focus:outline-none focus:shadow-outline"
            id="yourName"
            required
        />
    </div>
    
    <div class="mb-6">
        <div class="flex items-center justify-between">
            <button class="bg-green hover:bg-green-dark text-white font-bold py-2 px-4 rounded" type="submit">Register</button>
        </div>
    </div>
</form>