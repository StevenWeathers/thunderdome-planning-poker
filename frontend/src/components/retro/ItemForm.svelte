<script>
    import SmileCircle from '../icons/SmileCircle.svelte'
    import TrashIcon from '../icons/TrashIcon.svelte'
    import FrownCircle from '../icons/FrownCircle.svelte'
    import QuestionCircle from '../icons/QuestionCircle.svelte'
    import { warrior as user } from '../../stores.js'
    import { _ } from '../../i18n.js'

    export let handleSubmit = () => {}
    export let handleDelete = () => {}
    export let itemType = 'worked'
    export let content = ''
    export let newItemPlaceholder = 'What worked well...'
    export let phase = 'brainstorm'
    export let isFacilitator = false
    export let items = []
    export let feedbackVisibility = 'visible'

    const handleFormSubmit = evt => {
        evt.preventDefault()

        handleSubmit(itemType, content)
        content = ''
    }
</script>

<div class="">
    <div class="flex items-center mb-4">
        <div class="flex-shrink pr-2">
            {#if itemType === 'worked'}
                <SmileCircle
                    class="w-8 h-8 text-green-500 dark:text-lime-400"
                />
            {:else if itemType === 'improve'}
                <FrownCircle class="w-8 h-8 text-red-500" />
            {:else if itemType === 'question'}
                <QuestionCircle
                    class="w-8 h-8 text-blue-500 dark:text-sky-400"
                />
            {/if}
        </div>
        <div class="flex-grow">
            <form on:submit="{handleFormSubmit}" class="flex">
                <input
                    bind:value="{content}"
                    placeholder="{newItemPlaceholder}"
                    class="dark:bg-gray-800 border-gray-300 dark:border-gray-700 border-2 appearance-none rounded py-2
                    px-3 text-gray-700 dark:text-gray-400 leading-tight focus:outline-none
                    focus:bg-white dark:focus:bg-gray-700 focus:border-indigo-500 dark:focus:border-yellow-400 w-full"
                    id="new{itemType}"
                    name="new{itemType}"
                    type="text"
                    required
                    disabled="{phase !== 'brainstorm' && !isFacilitator}"
                />
                <button type="submit" class="hidden"></button>
            </form>
        </div>
    </div>
    <div>
        {#each items as item}
            <div
                class="p-2 mb-2 bg-white dark:bg-gray-800 shadow item-list-item border-l-4"
                class:border-green-400="{item.type === 'worked'}"
                class:dark:border-lime-400="{item.type === 'worked'}"
                class:border-red-500="{item.type === 'improve'}"
                class:border-blue-400="{item.type === 'question'}"
                class:dark:border-sky-400="{item.type === 'question'}"
                data-itemType="{itemType}"
                data-itemId="{item.id}"
            >
                <div class="flex items-center">
                    <div class="flex-grow">
                        <div class="flex items-center">
                            <div class="flex-grow dark:text-gray-200">
                                {#if feedbackVisibility === 'hidden' && item.userId !== $user.id}
                                    <span class="italic"
                                        >{$_('retroFeedbackHidden')}</span
                                    >
                                {:else if feedbackVisibility === 'concealed' && item.userId !== $user.id}
                                    <span class="italic"
                                        >{$_(
                                            'retroFeedbackConcealed',
                                        )}&nbsp;&nbsp;</span
                                    ><span class="text-white dark:text-gray-800"
                                        >{item.content}</span
                                    >
                                {:else}
                                    {item.content}
                                {/if}
                            </div>
                        </div>
                    </div>
                    <div class="flex-shrink pl-2">
                        {#if phase === 'brainstorm'}
                            <button
                                on:click="{handleDelete(itemType, item.id)}"
                                class="pr-2 pt-1 {item.userId !== $user.id
                                    ? 'text-gray-300 dark:text-gray-600 cursor-not-allowed'
                                    : 'text-gray-500 dark:text-gray-400 hover:text-red-500'}"
                                disabled="{item.userId !== $user.id}"
                            >
                                <TrashIcon />
                            </button>
                        {/if}
                    </div>
                </div>
            </div>
        {/each}
    </div>
</div>
