<script>
    import { onMount } from 'svelte'

    import PageLayout from '../components/PageLayout.svelte'
    import { warrior } from '../stores.js'

    export let notifications
    export let router
    let battles = []
    let battleName = ''

    fetch('/api/battles', {
        method: 'GET',
        credentials: 'same-origin'
    })
        .then(function(response) {
            return response.json()
        })
        .then(function(bs) {
            console.log(bs)
            battles = bs
        })
        .catch(function(error) {
            notifications.danger("Error finding your battles")
        })

    function createBattle(e) {
        e.preventDefault()
        const data = {
            battleName,
            leaderId: $warrior.id
        }
        
        fetch('/api/battle', {
            method: 'POST',
            credentials: 'same-origin',
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(data)
        })
            .then(function(response) {
                return response.json()
            })
            .then(function(battle) {
                router.route(`/battle/${battle.id}`)
            })
            .catch(function(error) {
                notifications.danger("Error encountered creating battle")
            })
    }

    onMount(() => {
        if (!$warrior.id) {
            router.route('/enlist')
        }
    })
</script>

<PageLayout>
    <h1 class="mb-4">My Battles</h1>

    <div class="bg-white shadow-md rounded mb-8">
        {#each battles as battle}
            <div class="flex flex-wrap p-4 border-grey-light border-b">
                <div class="w-full md:w-1/2 mb-4 md:mb-0 font-semibold md:text-xl">
                    {battle.name}
                </div>
                <div class="w-full md:w-1/2 mb-4 md:mb-0 md:text-right">
                    <a
                        href="/battle/{battle.id}"
                        class="inline-block bg-transparent hover:bg-green text-green-dark font-semibold hover:text-white no-underline py-2 px-2 border border-green hover:border-transparent rounded"
                    >
                        Go To Battle
                    </a>
                </div>
            </div>
        {/each}
    </div>

    <div class="bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4">
        <h2 class="mb-4">Create a Battle</h2>
        <form on:submit={createBattle} name="createBattle">
            <div class="mb-4">
                <label class="block text-grey-darker text-sm font-bold mb-2" for="battleName">Battle Name</label>
                <div class="control">
                    <input
                        name="battleName"
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
    </div>
</PageLayout>
