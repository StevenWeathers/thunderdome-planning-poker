<script>
    import { onMount } from 'svelte'

    import SolidButton from './SolidButton.svelte'
    import HollowButton from './HollowButton.svelte'
    import JiraImport from './JiraImport.svelte'
    import { warrior } from '../stores.js'
    import { _ } from '../i18n'

    export let notifications
    export let eventTag
    export let router
    export let xfetch

    const allowedPointValues = appConfig.AllowedPointValues

    let points = appConfig.DefaultPointValues
    let battleName = ''
    let plans = []
    let autoFinishVoting = true

    let checkedPointColor = 'border-green-500 bg-green-100 text-green-600'
    let uncheckedPointColor = 'border-gray-300 bg-white'

    function addPlan() {
        plans.unshift({
            name: '',
            type: $_('actions.plan.types.story'),
            referenceId: '',
            link: '',
            description: '',
            acceptanceCriteria: '',
        })
        plans = plans
    }

    function handlePlanImport(newPlan) {
        const plan = {
            name: newPlan.planName,
            type: newPlan.type,
            referenceId: newPlan.referenceId,
            link: newPlan.Link,
            description: newPlan.description,
            acceptanceCriteria: newPlan.acceptanceCriteria,
        }
        plans.unshift(plan)
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

        const pointValuesAllowed = allowedPointValues.filter(pv => {
            return points.includes(pv)
        })

        const body = {
            battleName,
            pointValuesAllowed,
            plans,
            autoFinishVoting,
        }

        xfetch('/api/battle', { body })
            .then(res => res.json())
            .then(function(battle) {
                eventTag('create_battle', 'engagement', 'success', () => {
                    router.route(`/battle/${battle.id}`)
                })
            })
            .catch(function(error) {
                notifications.danger(
                    $_('pages.myBattles.createBattle.createError'),
                )
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
            {$_('pages.myBattles.createBattle.fields.name.label')}
        </label>
        <div class="control">
            <input
                name="battleName"
                bind:value="{battleName}"
                placeholder="{$_('pages.myBattles.createBattle.fields.name.placeholder')}"
                class="bg-gray-200 border-gray-200 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 leading-tight
                focus:outline-none focus:bg-white focus:border-purple-500"
                id="battleName"
                required />
        </div>
    </div>

    <div class="mb-4">
        <h3 class="block text-gray-700 text-sm font-bold mb-2">
            {$_('pages.myBattles.createBattle.fields.allowedPointValues.label')}
        </h3>
        <div class="control relative -mr-2 md:-mr-1">
            {#each allowedPointValues as point, pi}
                <label
                    class="{points.includes(point) ? checkedPointColor : uncheckedPointColor}
                    cursor-pointer font-bold border p-2 mr-2 xl:mr-1 mb-2
                    xl:mb-0 rounded inline-block">
                    <input
                        type="checkbox"
                        bind:group="{points}"
                        value="{point}"
                        class="hidden" />
                    {point}
                </label>
            {/each}
        </div>
    </div>

    <div class="mb-4">
        <h3 class="block text-gray-700 text-sm font-bold mb-2">
            {$_('pages.myBattles.createBattle.fields.plans.label')}
        </h3>
        <div class="control mb-4">
            <JiraImport handlePlanAdd="{handlePlanImport}" {notifications} />
            <HollowButton onClick="{addPlan}">
                {$_('pages.myBattles.createBattle.fields.plans.addButton')}
            </HollowButton>
        </div>
        {#each plans as plan, i}
            <div class="flex flex-wrap mb-2">
                <div class="w-3/4">
                    <input
                        type="text"
                        bind:value="{plan.name}"
                        placeholder="{$_('pages.myBattles.createBattle.fields.plans.fields.name.placeholder')}"
                        class="bg-gray-200 border-gray-200 border-2
                        appearance-none rounded w-full py-2 px-3 text-gray-700
                        leading-tight focus:outline-none focus:bg-white
                        focus:border-purple-500"
                        required />
                </div>
                <div class="w-1/4">
                    <div class="pl-2">
                        <HollowButton onClick="{removePlan(i)}" color="red">
                            {$_('pages.myBattles.createBattle.fields.plans.removeButton')}
                        </HollowButton>
                    </div>
                </div>
            </div>
        {/each}
    </div>

    <div class="mb-4">
        <label class="text-gray-700 text-sm font-bold mb-2">
            <input
                type="checkbox"
                bind:checked="{autoFinishVoting}"
                id="autoFinishVoting"
                name="autoFinishVoting" />
            {$_('pages.myBattles.createBattle.fields.autoFinishVoting.label')}
        </label>
    </div>

    <div class="text-right">
        <SolidButton type="submit">{$_('actions.battle.create')}</SolidButton>
    </div>
</form>
