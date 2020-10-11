<script>
    import TurndownService from 'turndown'
    import marked from 'marked'

    import { quill } from '../quill'
    import SolidButton from './SolidButton.svelte'
    import CloseIcon from './icons/CloseIcon.svelte'
    import DownCarrotIcon from './icons/DownCarrotIcon.svelte'

    export let handlePlanAdd = () => {}
    export let toggleAddPlan = () => {}
    export let handlePlanRevision = () => {}

    // going by common Jira issue types for now
    const planTypes = ['story', 'bug', 'spike', 'epic', 'task', 'subtask']
    const turndownService = new TurndownService()

    export let planId = ''
    export let planName = ''
    export let planType = 'story'
    export let referenceId = ''
    export let planLink = ''
    export let description = ''
    export let acceptanceCriteria = ''

    function handleSubmit(event) {
        event.preventDefault()
        const plan = {
            planName,
            type: planType,
            referenceId,
            link: planLink,
            description: turndownService.turndown(description),
            acceptanceCriteria: turndownService.turndown(acceptanceCriteria),
        }
        if (planId === '') {
            handlePlanAdd(plan)
        } else {
            plan.planId = planId
            handlePlanRevision(plan)
        }

        toggleAddPlan()
    }
</script>

<div
    class="fixed inset-0 flex items-center z-40 max-h-screen overflow-y-scroll">
    <div class="fixed inset-0 bg-gray-900 opacity-75"></div>

    <div
        class="relative mx-4 md:mx-auto w-full md:w-2/3 lg:w-3/5 xl:w-1/2 z-50
        max-h-full">
        <div class="py-8">
            <div class="shadow-xl bg-white rounded-lg p-4 xl:p-6 max-h-full">
                <div class="flex justify-end mb-2">
                    <button
                        aria-label="close"
                        on:click="{toggleAddPlan}"
                        class="text-gray-800">
                        <CloseIcon />
                    </button>
                </div>

                <form on:submit="{handleSubmit}" name="addPlan">
                    <div class="mb-4">
                        <label
                            class="block text-sm font-bold mb-2"
                            for="planName">
                            Plan Type
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
                                <option value="" disabled>Types</option>
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
                            Plan Name
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
                            placeholder="Enter a plan name" />
                    </div>
                    <div class="mb-4">
                        <label
                            class="block text-sm font-bold mb-2"
                            for="referenceId">
                            Reference ID
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
                            placeholder="Enter a reference ID" />
                    </div>
                    <div class="mb-4">
                        <label
                            class="block text-sm font-bold mb-2"
                            for="planLink">
                            Link
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
                            placeholder="Enter a link to story" />
                    </div>
                    <div class="mb-16">
                        <div class="text-sm font-bold mb-2">Description</div>
                        <div class="h-48">
                            <div
                                class="w-full"
                                use:quill="{{ placeholder: 'Enter a plan description', content: marked(description) }}"
                                on:text-change="{e => (description = e.detail.html)}"
                                id="description"></div>
                        </div>
                    </div>
                    <div class="mb-16">
                        <div class="text-sm font-bold mb-2">
                            Acceptance Criteria
                        </div>
                        <div class="h-48">
                            <div
                                class="w-full"
                                use:quill="{{ placeholder: 'Enter plan acceptance criteria', content: marked(acceptanceCriteria) }}"
                                on:text-change="{e => (acceptanceCriteria = e.detail.html)}"
                                id="acceptanceCriteria"></div>
                        </div>
                    </div>
                    <div class="text-right">
                        <div>
                            <SolidButton type="submit">Save</SolidButton>
                        </div>
                    </div>
                </form>
            </div>
        </div>
    </div>
</div>
