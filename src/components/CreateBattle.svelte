<script>
    import { onMount } from 'svelte'

    import DownCarrotIcon from '../components/icons/DownCarrotIcon.svelte'
    import SolidButton from '../components/SolidButton.svelte'
    import HollowButton from '../components/HollowButton.svelte'
    import { warrior } from '../stores.js'

    export let notifications
    export let router

    const possiblePointValues = [
        ['1', '2', '3', '5', '8', '13', '?'],
        ['1/2', '1', '2', '3', '5', '8', '13', '?'],
        ['0', '1/2', '1', '2', '3', '5', '8', '13', '20', '40', '100', '?'],
    ]

    let battleName = ''
    let pointValuesAllowed = 2
    let plans = []

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
        const data = {
            battleName,
            pointValuesAllowed: possiblePointValues[pointValuesAllowed],
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
                router.route(`/battle/${battle.id}`)
            })
            .catch(function(error) {
                notifications.danger('Error encountered creating battle')
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
        <label
            class="block text-gray-700 text-sm font-bold mb-2"
            for="pointValuesAllowed">
            Allowed Point Values
        </label>
        <div class="control relative">
            <select
                name="pointValuesAllowed"
                bind:value="{pointValuesAllowed}"
                class="block appearance-none w-full bg-gray-200 border-2
                border-gray-200 text-gray-darker py-3 px-4 pr-8 rounded
                leading-tight focus:outline-none focus:bg-white
                focus:border-purple-500"
                id="pointValuesAllowed"
                required>
                {#each possiblePointValues as points, pi}
                    <option value="{pi}" selected="{pi === pointValuesAllowed}">
                        {points.join(', ')}
                    </option>
                {/each}
            </select>
            <div
                class="pointer-events-none absolute inset-y-0 right-0 flex
                items-center px-2 text-gray-700">
                <DownCarrotIcon />
            </div>
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
