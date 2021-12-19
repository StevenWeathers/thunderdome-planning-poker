<script>
    import { quill } from '../../quill.js'
    import SolidButton from '../SolidButton.svelte'
    import Modal from '../Modal.svelte'
    import DownCarrotIcon from '../icons/ChevronDown.svelte'
    import { _ } from '../../i18n.js'

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
    let descriptionExpanded = false
    let acceptanceExpanded = false

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
            <label
                class="block font-bold mb-2 dark:text-gray-400"
                for="planName"
            >
                {$_('planType')}
            </label>
            <div class="relative">
                <select
                    name="planType"
                    bind:value="{planType}"
                    required
                    class="block appearance-none w-full border-2 dark:bg-gray-900 border-gray-300 dark:border-gray-600
                    text-gray-700 dark:text-gray-400 py-3 px-4 pr-8 rounded leading-tight
                    focus:outline-none focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
                >
                    <option value="" disabled>
                        {$_('planTypePlaceholder')}
                    </option>
                    {#each planTypes as pType}
                        <option value="{pType}">{pType}</option>
                    {/each}
                </select>
                <div
                    class="pointer-events-none absolute inset-y-0 right-0 flex
                    items-center px-2 text-gray-700 dark:text-gray-300"
                >
                    <DownCarrotIcon />
                </div>
            </div>
        </div>
        <div class="mb-4">
            <label
                class="block font-bold mb-2 dark:text-gray-400"
                for="planName"
            >
                {$_('planName')}
            </label>
            <input
                class="bg-gray-100 dark:bg-gray-900 dark:focus:bg-gray-800 border-gray-200 dark:border-gray-600 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 dark:text-gray-400 leading-tight
                focus:outline-none focus:bg-white focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
                type="text"
                id="planName"
                name="planName"
                bind:value="{planName}"
                placeholder="{$_('planNamePlaceholder')}"
            />
        </div>
        <div class="mb-4">
            <label
                class="block font-bold mb-2 dark:text-gray-400"
                for="referenceId"
            >
                {$_('planReferenceId')}
            </label>
            <input
                class="bg-gray-100 dark:bg-gray-900 dark:focus:bg-gray-800 border-gray-200 dark:border-gray-600 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 dark:text-gray-400 leading-tight
                focus:outline-none focus:bg-white focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
                type="text"
                id="referenceId"
                name="referenceId"
                bind:value="{referenceId}"
                placeholder="{$_('planReferenceIdPlaceholder')}"
            />
        </div>
        <div class="mb-4">
            <label
                class="block font-bold mb-2 dark:text-gray-400"
                for="planLink"
            >
                {$_('planLink')}
            </label>
            <input
                class="bg-gray-100  dark:bg-gray-900 dark:focus:bg-gray-800 border-gray-200 dark:border-gray-600 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 dark:text-gray-400 leading-tight
                focus:outline-none focus:bg-white focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
                type="text"
                id="planLink"
                name="planLink"
                bind:value="{planLink}"
                placeholder="{$_('planLinkPlaceholder')}"
            />
        </div>
        <div>
            <div class="font-bold mb-2">
                <button
                    on:click="{e => {
                        e.preventDefault()
                        descriptionExpanded = !descriptionExpanded
                    }}"
                    class="inline-block align-baseline text-sm
                        text-blue-700 dark:text-sky-400 hover:text-blue-800 dark:hover:text-sky-600 bg-transparent
                        border-transparent mr-1 font-bold text-xl"
                    type="button"
                >
                    {#if descriptionExpanded}-{:else}+{/if}
                </button>
                <span class="dark:text-gray-400">{$_('planDescription')}</span>
            </div>
            {#if descriptionExpanded}
                <div class="mb-2">
                    <div class="bg-white">
                        <div
                            class="w-full bg-white"
                            use:quill="{{
                                placeholder: $_('planDescriptionPlaceholder'),
                                content: description,
                            }}"
                            on:text-change="{e =>
                                (description = e.detail.html)}"
                            id="description"
                        ></div>
                    </div>
                </div>
            {/if}
        </div>
        <div>
            <div class="font-bold mb-2">
                <button
                    on:click="{e => {
                        e.preventDefault()
                        acceptanceExpanded = !acceptanceExpanded
                    }}"
                    class="inline-block align-baseline text-sm
                        text-blue-700 dark:text-sky-400 hover:text-blue-800 dark:hover:text-sky-600 bg-transparent
                        border-transparent mr-1 font-bold text-xl"
                    type="button"
                >
                    {#if acceptanceExpanded}-{:else}+{/if}
                </button>
                <span class="dark:text-gray-400"
                    >{$_('planAcceptanceCriteria')}</span
                >
            </div>
            {#if acceptanceExpanded}
                <div class="mb-2">
                    <div class="bg-white">
                        <div
                            class="w-full"
                            use:quill="{{
                                placeholder: $_(
                                    'planAcceptanceCriteriaPlaceholder',
                                ),
                                content: acceptanceCriteria,
                            }}"
                            on:text-change="{e =>
                                (acceptanceCriteria = e.detail.html)}"
                            id="acceptanceCriteria"
                        ></div>
                    </div>
                </div>
            {/if}
        </div>
        <div class="text-right">
            <div>
                <SolidButton type="submit" testid="plan-save"
                    >{$_('save')}</SolidButton
                >
            </div>
        </div>
    </form>
</Modal>
