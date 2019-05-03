<script>
    let creatorName = ''
    let battleName = ''

    function createBattle(e) {
        e.preventDefault()
        
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
                window.location.href = `/battle/${battle.id}`
            });
    }
</script>

<div class="columns">
    <div class="column">
        <form on:submit={createBattle}>
            <div>
                <input bind:value={creatorName} placeholder="Enter your name" required />
            </div>
            <div>
                <input bind:value={battleName} placeholder="Enter a battle name" required />
            </div>
            <button class="button is-primary" type="submit">Create a Story Battle</button>
        </form>
    </div>
</div>