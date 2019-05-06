<script>
    import { warrior } from '../stores.js'

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
            });
    }
</script>

<div class="columns">
    <div class="column">
        <form on:submit={createWarrior}>
            <div class="field">
                <label class="label">Name</label>
                <div class="control">
                    <input bind:value={warriorName} placeholder="Enter your name" class="input" required />
                </div>
            </div>
            
            <div class="field">
                <div class="control">
                    <button class="button is-success" type="submit">Register</button>
                </div>
            </div>
        </form>
    </div>
</div>