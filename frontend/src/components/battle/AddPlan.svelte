<script lang="ts">
    import { quill } from '../../quill'
    import SolidButton from '../SolidButton.svelte'
    import Modal from '../Modal.svelte'
    import NoSymbol from '../icons/NoSymbol.svelte'
    import DoubleChevronUp from '../icons/DoubleChevronUp.svelte'
    import ChevronUp from '../icons/ChevronUp.svelte'
    import Bars2 from '../icons/Bars2.svelte'
    import ChevronDown from '../icons/ChevronDown.svelte'
    import DoubleChevronDown from '../icons/DoubleChevronDown.svelte'
    import LL from '../../i18n/i18n-svelte'
    import { AppConfig } from '../../config'

    export let handlePlanAdd = () => {}
    export let toggleAddPlan = () => {}
    export let handlePlanRevision = () => {}
    export let eventTag = () => {}
    export let notifications

    // going by common Jira issue types for now
    const planTypes = [
        $LL.planTypeStory(),
        $LL.planTypeBug(),
        $LL.planTypeSpike(),
        $LL.planTypeEpic(),
        $LL.planTypeTask(),
        $LL.planTypeSubtask(),
    ]

    // going by common Jira issue priorities for now
    const priorities = [
        { name: $LL.planPriorityBlocker(), value: 1, icon: NoSymbol },
        {
            name: $LL.planPriorityHighest(),
            value: 2,
            icon: DoubleChevronUp,
        },
        { name: $LL.planPriorityHigh(), value: 3, icon: ChevronUp },
        { name: $LL.planPriorityMedium(), value: 4, icon: Bars2 },
        { name: $LL.planPriorityLow(), value: 5, icon: ChevronDown },
        {
            name: $LL.planPriorityLowest(),
            value: 6,
            icon: DoubleChevronDown,
        },
    ]

    export let planId = ''
    export let planName = ''
    export let planType = $LL.planTypeStory()
    export let referenceId = ''
    export let planLink = ''
    export let description = ''
    export let acceptanceCriteria = ''
    export let priority = 99

    const isAbsolute = new RegExp('^([a-z]+://|//)', 'i')
    let descriptionExpanded = false
    let acceptanceExpanded = false

    function handleSubmit(event) {
        event.preventDefault()
        let invalidPlan = false

        if (planLink !== '' && !isAbsolute.test(planLink)) {
            invalidPlan = true
            notifications.danger($LL.planLinkInvalid())
            eventTag('plan_add_invalid_link', 'battle', ``)
        }

        const plan = {
            planName,
            type: planType,
            referenceId,
            link: planLink,
            description,
            acceptanceCriteria,
            priority,
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
                for="planType"
            >
                {$LL.planType({ friendly: AppConfig.FriendlyUIVerbs })}
            </label>
            <div class="relative">
                <select
                    name="planType"
                    id="planType"
                    bind:value="{planType}"
                    required
                    class="block appearance-none w-full border-2 dark:bg-gray-900 border-gray-300 dark:border-gray-600
                    text-gray-700 dark:text-gray-400 py-3 px-4 pe-8 rounded leading-tight
                    focus:outline-none focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
                >
                    <option value="" disabled>
                        {$LL.planTypePlaceholder({
                            friendly: AppConfig.FriendlyUIVerbs,
                        })}
                    </option>
                    {#each planTypes as pType}
                        <option value="{pType}">{pType}</option>
                    {/each}
                </select>
                <div
                    class="pointer-events-none absolute inset-y-0 end-0 flex
                    items-center px-2 text-gray-700 dark:text-gray-300"
                >
                    <ChevronDown />
                </div>
            </div>
        </div>
        <div class="mb-4">
            <label
                class="block font-bold mb-2 dark:text-gray-400"
                for="planName"
            >
                {$LL.planName({ friendly: AppConfig.FriendlyUIVerbs })}
            </label>
            <input
                class="bg-gray-100 dark:bg-gray-900 dark:focus:bg-gray-800 border-gray-200 dark:border-gray-600 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 dark:text-gray-400 leading-tight
                focus:outline-none focus:bg-white focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
                type="text"
                id="planName"
                name="planName"
                bind:value="{planName}"
                placeholder="{$LL.planNamePlaceholder({
                    friendly: AppConfig.FriendlyUIVerbs,
                })}"
            />
        </div>
        <div class="mb-4">
            <label
                class="block font-bold mb-2 dark:text-gray-400"
                for="referenceId"
            >
                {$LL.planReferenceId()}
            </label>
            <input
                class="bg-gray-100 dark:bg-gray-900 dark:focus:bg-gray-800 border-gray-200 dark:border-gray-600 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 dark:text-gray-400 leading-tight
                focus:outline-none focus:bg-white focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
                type="text"
                id="referenceId"
                name="referenceId"
                bind:value="{referenceId}"
                placeholder="{$LL.planReferenceIdPlaceholder()}"
            />
        </div>
        <div class="mb-4">
            <label
                class="block font-bold mb-2 dark:text-gray-400"
                for="planLink"
            >
                {$LL.planLink()}
            </label>
            <input
                class="bg-gray-100 dark:bg-gray-900 dark:focus:bg-gray-800 border-gray-200 dark:border-gray-600 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 dark:text-gray-400 leading-tight
                focus:outline-none focus:bg-white focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
                type="text"
                id="planLink"
                name="planLink"
                bind:value="{planLink}"
                placeholder="{$LL.planLinkPlaceholder({
                    friendly: AppConfig.FriendlyUIVerbs,
                })}"
            />
        </div>
        <div class="mb-4">
            <label
                class="block font-bold mb-2 dark:text-gray-400"
                for="priority"
            >
                {$LL.planPriority()}
            </label>
            <div class="relative">
                <select
                    name="priority"
                    id="priority"
                    bind:value="{priority}"
                    class="block appearance-none w-full border-2 dark:bg-gray-900 border-gray-300 dark:border-gray-600
                    text-gray-700 dark:text-gray-400 py-3 px-4 pe-8 rounded leading-tight
                    focus:outline-none focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
                >
                    <option value="{99}" disabled>
                        {$LL.planPriorityPlaceholder()}
                    </option>
                    {#each priorities as p}
                        <option value="{p.value}">
                            <svelte:component this="{p.icon}" />{p.name}</option
                        >
                    {/each}
                </select>
                <div
                    class="pointer-events-none absolute inset-y-0 end-0 flex
                    items-center px-2 text-gray-700 dark:text-gray-300"
                >
                    <ChevronDown />
                </div>
            </div>
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
                        border-transparent me-1 font-bold text-xl"
                    type="button"
                >
                    {#if descriptionExpanded}-{:else}+{/if}
                </button>
                <span class="dark:text-gray-400">{$LL.planDescription()}</span>
            </div>
            {#if descriptionExpanded}
                <div class="mb-2">
                    <div class="bg-white">
                        <div
                            class="w-full bg-white"
                            use:quill="{{
                                placeholder: $LL.planDescriptionPlaceholder(),
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
                        border-transparent me-1 font-bold text-xl"
                    type="button"
                >
                    {#if acceptanceExpanded}-{:else}+{/if}
                </button>
                <span class="dark:text-gray-400"
                    >{$LL.planAcceptanceCriteria()}</span
                >
            </div>
            {#if acceptanceExpanded}
                <div class="mb-2">
                    <div class="bg-white">
                        <div
                            class="w-full"
                            use:quill="{{
                                placeholder:
                                    $LL.planAcceptanceCriteriaPlaceholder(),
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
                    >{$LL.save()}</SolidButton
                >
            </div>
        </div>
    </form>
</Modal>
