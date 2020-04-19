<script>
    import { onMount } from 'svelte'

    import DownCarrotIcon from '../components/icons/DownCarrotIcon.svelte'
    import SolidButton from '../components/SolidButton.svelte'
    import HollowButton from '../components/HollowButton.svelte'
    import { warrior } from '../stores.js'

    export let notifications
    export let eventTag
    export let router

    const possiblePointValues = [ '0', '1/2', '1', '2', '3', '5', '8', '13', '20', '40', '100', '?' ]

    let battleName = ''
    let points = [ '1', '2', '3', '5', '8', '13', '?' ]
    let plans = []

    let checkedPointColor = 'border-green-500 bg-green-100 text-green-600'
    let uncheckedPointColor = 'border-gray-300 bg-white'

    function addPlan() {
        plans.unshift({
            name: '',
        })
        plans = plans
    }

    function removePlan(i) {
        return function remove() {
            plans.splice(i, 1)
            plans = plans
        }
    }

    function createBattle(e) {
        e.preventDefault()

        const pointValuesAllowed = possiblePointValues.filter(pv => {
            return points.includes(pv)
        })

        const data = {
            battleName,
            pointValuesAllowed,
            plans,
        }

        fetch('/api/battle', {
            method: 'POST',
            credentials: 'same-origin',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data),
        })
            .then(function(response) {
                return response.json()
            })
            .then(function(battle) {
                eventTag('create_battle', 'engagement', 'success', () => {
                    router.route(`/battle/${battle.id}`)
                })
            })
            .catch(function(error) {
                notifications.danger('Error encountered creating battle')
                eventTag('create_battle', 'engagement', 'failure')
            })
    }

    onMount(() => {
        if (!$warrior.id) {
            router.route('/enlist')
        }
    })
</script>

<form on:submit="{createBattle}" name="createBattle">
    <div class="mb-4">
        <label
            class="block text-gray-700 text-sm font-bold mb-2"
            for="battleName">
            Battle Name
        </label>
        <div class="control">
            <input
                name="battleName"
                bind:value="{battleName}"
                placeholder="Enter a battle name"
                class="bg-gray-200 border-gray-200 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 leading-tight
                focus:outline-none focus:bg-white focus:border-purple-500"
                id="battleName"
                required />
        </div>
    </div>

    <div class="mb-4">
        <h3
            class="block text-gray-700 text-sm font-bold mb-2">
            Allowed Point Values
        </h3>
        <div class="control relative -mr-2 md:-mr-1">
            {#each possiblePointValues as point, pi}
                <label class="{points.includes(point) ? checkedPointColor : uncheckedPointColor} 
                cursor-pointer font-bold border p-2 mr-2 xl:mr-1 mb-2 xl:mb-0 rounded inline-block"
                >
                    <input type="checkbox" bind:group={points} value={point} class="hidden" />
                    {point}
                </label>
            {/each}
        </div>
    </div>

    <div class="mb-4">
        <h3 class="block text-gray-700 text-sm font-bold mb-2">Plans</h3>
        <div class="control mb-4">
            <HollowButton onClick="{addPlan}">Add Plan</HollowButton>
        </div>
        {#each plans as plan, i}
            <div class="flex flex-wrap mb-2">
                <div class="w-3/4">
                    <input
                        type="text"
                        bind:value="{plan.name}"
                        placeholder="plan name"
                        class="bg-gray-200 border-gray-200 border-2
                        appearance-none rounded w-full py-2 px-3 text-gray-700
                        leading-tight focus:outline-none focus:bg-white
                        focus:border-purple-500"
                        required />
                </div>
                <div class="w-1/4">
                    <div class="pl-2">
                        <HollowButton onClick="{removePlan(i)}" color="red">
                            Remove
                        </HollowButton>
                    </div>
                </div>
            </div>
        {/each}
    </div>

    <div class="text-right">
        <SolidButton type="submit">Create Battle</SolidButton>
    </div>
</form>
