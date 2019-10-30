<script>
    import { onMount } from 'svelte'

    import PageLayout from '../components/PageLayout.svelte'
    import { warrior } from '../stores.js'

    export let notifications
    export let router
    let battles = []
    let battleName = ''
    let pointValuesAllowed = 2
    const possiblePointValues = [
        ["1", "2", "3", "5", "8", "13", "?"],
        ["1/2", "1", "2", "3", "5", "8", "13", "?"],
        ["0", "1/2", "1", "2", "3", "5", "8", "13", "20", "40", "100", "?"]
    ]

    fetch('/api/battles', {
        method: 'GET',
        credentials: 'same-origin'
    })
        .then(function(response) {
            return response.json()
        })
        .then(function(bs) {
            battles = bs
        })
        .catch(function(error) {
            notifications.danger("Error finding your battles")
        })

    function createBattle(e) {
        e.preventDefault()
        const data = {
            battleName,
            leaderId: $warrior.id,
            pointValuesAllowed: possiblePointValues[pointValuesAllowed],
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

    <div class="mb-8">
        {#each battles as battle}
            <div class="bg-white shadow-md rounded mb-2">
                <div class="flex flex-wrap p-4 border-grey-light border-b">
                    <div class="w-full md:w-1/2 mb-4 md:mb-0 font-semibold md:text-xl">
                        {battle.name}
                        <div class="font-semibold md:text-sm text-grey-dark">
                            {battle.plans.filter(p => p.points !== "").length} of {battle.plans.length} plans pointed
                        </div>
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
            </div>
        {/each}
    </div>

    <div class="bg-white shadow-md rounded p-6">
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

            <div class="mb-4">
                <label class="block text-grey-darker text-sm font-bold mb-2" for="pointValuesAllowed">Allowed Point Values</label>
                <div class="control relative">
                    <select
                        name="pointValuesAllowed"
                        bind:value={pointValuesAllowed}
                        class="block appearance-none w-full bg-grey-lighter border border-grey-lighter text-grey-darker py-3 px-4 pr-8 rounded leading-tight focus:outline-none focus:bg-white focus:border-grey"
                        id="pointValuesAllowed"
                        required
                    >
                        {#each possiblePointValues as points, pi}
                            <option value={pi} selected={pi === pointValuesAllowed}>{points.join(', ')}</option>
                        {/each}
                    </select>
                    <div class="pointer-events-none absolute pin-y pin-r flex items-center px-2 text-grey-darker">
                        <svg class="fill-current h-4 w-4" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20"><path d="M9.293 12.95l.707.707L15.657 8l-1.414-1.414L10 10.828 5.757 6.586 4.343 8z"/></svg>
                    </div>
                </div>
            </div>
            
            <div>
                <button class="bg-green hover:bg-green-dark text-white font-bold py-2 px-4 rounded" type="submit">Create a Story Battle</button>
            </div>
        </form>
    </div>
</PageLayout>
