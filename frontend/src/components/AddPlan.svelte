<script>
    import { quill } from '../quill'
    import SolidButton from './SolidButton.svelte'
    import Modal from './Modal.svelte'
    import DownCarrotIcon from './icons/DownCarrotIcon.svelte'
    import { _ } from '../i18n'

    export let handlePlanAdd = () => {}
    export let toggleAddPlan = () => {}
    export let handlePlanRevision = () => {}
    export let eventTag = () => {}
    export let notifications

    // going by common Jira issue types for now
    const planTypes = [
        $_('actions.plan.types.story'),
        $_('actions.plan.types.bug'),
        $_('actions.plan.types.spike'),
        $_('actions.plan.types.epic'),
        $_('actions.plan.types.task'),
        $_('actions.plan.types.subtask'),
    ]

    export let planId = ''
    export let planName = ''
    export let planType = $_('actions.plan.types.story')
    export let referenceId = ''
    export let planLink = ''
    export let description = ''
    export let acceptanceCriteria = ''

    const isAbsolute = new RegExp('^([a-z]+://|//)', 'i')

    function handleSubmit(event) {
        event.preventDefault()
        let invalidPlan = false

        if (planLink !== '' && !isAbsolute.test(planLink)) {
            invalidPlan = true
            notifications.danger($_('actions.plan.fields.link.invalid'))
            eventTag('plan_add_invalid_link', 'battle', ``)
        }

        const plan = {
            planName,
            type: planType,
            referenceId,
            link: planLink,
            description,
            acceptanceCriteria,
        }

        if (!invalidPlan) {
            if (planId === '') {
                handlePlanAdd(plan)
            } else {
                plan.planId = planId
                handlePlanRevision(plan)
            }

            toggleAddPlan()
        }
    }
</script>

<Modal closeModal={toggleAddPlan} widthClasses="md:w-2/3 lg:w-3/5 xl:w-1/2">
    <form on:submit="{handleSubmit}" name="addPlan">
        <div class="mb-4">
            <label
                class="block text-sm font-bold mb-2"
                for="planName">
                {$_('actions.plan.fields.type.label')}
            </label>
            <div class="relative">
                <select
                    name="planType"
                    bind:value="{planType}"
                    required
                    class="block appearance-none w-full border-2
                    border-gray-400 text-gray-700 py-3 px-4 pr-8
                    rounded leading-tight focus:outline-none
                    focus:border-purple-500">
                    <option value="" disabled>
                        {$_('actions.plan.fields.type.placeholder')}
                    </option>
                    {#each planTypes as pType}
                        <option value="{pType}">{pType}</option>
                    {/each}
                </select>
                <div
                    class="pointer-events-none absolute inset-y-0
                    right-0 flex items-center px-2 text-gray-700">
                    <DownCarrotIcon />
                </div>
            </div>
        </div>
        <div class="mb-4">
            <label
                class="block text-sm font-bold mb-2"
                for="planName">
                {$_('actions.plan.fields.name.label')}
            </label>
            <input
                class="bg-gray-200 border-gray-200 border-2
                appearance-none rounded w-full py-2 px-3
                text-gray-700 leading-tight focus:outline-none
                focus:bg-white focus:border-purple-500"
                type="text"
                id="planName"
                name="planName"
                bind:value="{planName}"
                placeholder="{$_('actions.plan.fields.name.placeholder')}" />
        </div>
        <div class="mb-4">
            <label
                class="block text-sm font-bold mb-2"
                for="referenceId">
                {$_('actions.plan.fields.referenceId.label')}
            </label>
            <input
                class="bg-gray-200 border-gray-200 border-2
                appearance-none rounded w-full py-2 px-3
                text-gray-700 leading-tight focus:outline-none
                focus:bg-white focus:border-purple-500"
                type="text"
                id="referenceId"
                name="referenceId"
                bind:value="{referenceId}"
                placeholder="{$_('actions.plan.fields.referenceId.placeholder')}" />
        </div>
        <div class="mb-4">
            <label
                class="block text-sm font-bold mb-2"
                for="planLink">
                {$_('actions.plan.fields.link.label')}
            </label>
            <input
                class="bg-gray-200 border-gray-200 border-2
                appearance-none rounded w-full py-2 px-3
                text-gray-700 leading-tight focus:outline-none
                focus:bg-white focus:border-purple-500"
                type="text"
                id="planLink"
                name="planLink"
                bind:value="{planLink}"
                placeholder="{$_('actions.plan.fields.link.placeholder')}" />
        </div>
        <div class="mb-16">
            <div class="text-sm font-bold mb-2">
                {$_('actions.plan.fields.description.label')}
            </div>
            <div class="h-48">
                <div
                    class="w-full"
                    use:quill="{{ placeholder: $_('actions.plan.fields.description.placeholder'), content: description }}"
                    on:text-change="{e => (description = e.detail.html)}"
                    id="description"></div>
            </div>
        </div>
        <div class="mb-16">
            <div class="text-sm font-bold mb-2">
                {$_('actions.plan.fields.acceptanceCriteria.label')}
            </div>
            <div class="h-48">
                <div
                    class="w-full"
                    use:quill="{{ placeholder: $_('actions.plan.fields.acceptanceCriteria.placeholder'), content: acceptanceCriteria }}"
                    on:text-change="{e => (acceptanceCriteria = e.detail.html)}"
                    id="acceptanceCriteria"></div>
            </div>
        </div>
        <div class="text-right">
            <div>
                <SolidButton type="submit">
                    {$_('actions.plan.save')}
                </SolidButton>
            </div>
        </div>
    </form>
</Modal>
