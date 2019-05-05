<script>
    import { user } from '../stores.js'

    let userName = ''

    function createUser(e) {
        e.preventDefault()
        
        fetch('/api/user', {
            method: 'POST',
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                userName
            })
        })
            .then(function(response) {
                return response.json()
            })
            .then(function(createdUser) {
                user.create({
                    id: createdUser.id,
                    name: createdUser.name
                })
            });
    }
</script>

<div class="columns">
    <div class="column">
        <form on:submit={createUser}>
            <div class="field">
                <label class="label">Name</label>
                <div class="control">
                    <input bind:value={userName} placeholder="Enter your name" class="input" required />
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