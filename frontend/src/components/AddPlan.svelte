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
        $_('planTypeStory'),
        $_('planTypeBug'),
        $_('planTypeSpike'),
        $_('planTypeEpic'),
        $_('planTypeTask'),
        $_('planTypeSubtask'),
    ]

    export let planId = ''
    export let planName = ''
    export let planType = $_('planTypeStory')
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
            notifications.danger($_('planLinkInvalid'))
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

<Modal closeModal="{toggleAddPlan}" widthClasses="md:w-2/3 lg:w-3/5 xl:w-1/2">
    <form on:submit="{handleSubmit}" name="addPlan">
        <div class="mb-4">
            <label class="block text-sm font-bold mb-2" for="planName">
                {$_('planType')}
            </label>
            <div class="relative">
                <select
                    name="planType"
                    bind:value="{planType}"
                    required
                    class="block appearance-none w-full border-2 border-gray-400
                    text-gray-700 py-3 px-4 pr-8 rounded leading-tight
                    focus:outline-none focus:border-purple-500">
                    <option value="" disabled>
                        {$_('planTypePlaceholder')}
                    </option>
                    {#each planTypes as pType}
                        <option value="{pType}">{pType}</option>
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
            <label class="block text-sm font-bold mb-2" for="planName">
                {$_('planName')}
            </label>
            <input
                class="bg-gray-200 border-gray-200 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 leading-tight
                focus:outline-none focus:bg-white focus:border-purple-500"
                type="text"
                id="planName"
                name="planName"
                bind:value="{planName}"
                placeholder="{$_('planNamePlaceholder')}" />
        </div>
        <div class="mb-4">
            <label class="block text-sm font-bold mb-2" for="referenceId">
                {$_('planReferenceId')}
            </label>
            <input
                class="bg-gray-200 border-gray-200 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 leading-tight
                focus:outline-none focus:bg-white focus:border-purple-500"
                type="text"
                id="referenceId"
                name="referenceId"
                bind:value="{referenceId}"
                placeholder="{$_('planReferenceIdPlaceholder')}" />
        </div>
        <div class="mb-4">
            <label class="block text-sm font-bold mb-2" for="planLink">
                {$_('planLink')}
            </label>
            <input
                class="bg-gray-200 border-gray-200 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 leading-tight
                focus:outline-none focus:bg-white focus:border-purple-500"
                type="text"
                id="planLink"
                name="planLink"
                bind:value="{planLink}"
                placeholder="{$_('planLinkPlaceholder')}" />
        </div>
        <div class="mb-16">
            <div class="text-sm font-bold mb-2">
                {$_('planDescription')}
            </div>
            <div class="h-48">
                <div
                    class="w-full"
                    use:quill="{{ placeholder: $_('planDescriptionPlaceholder'), content: description }}"
                    on:text-change="{e => (description = e.detail.html)}"
                    id="description"></div>
            </div>
        </div>
        <div class="mb-16">
            <div class="text-sm font-bold mb-2">
                {$_('planAcceptanceCriteria')}
            </div>
            <div class="h-48">
                <div
                    class="w-full"
                    use:quill="{{ placeholder: $_('planAcceptanceCriteriaPlaceholder'), content: acceptanceCriteria }}"
                    on:text-change="{e => (acceptanceCriteria = e.detail.html)}"
                    id="acceptanceCriteria"></div>
            </div>
        </div>
        <div class="text-right">
            <div>
                <SolidButton type="submit">
                    {$_('save')}
                </SolidButton>
            </div>
        </div>
    </form>
</Modal>
