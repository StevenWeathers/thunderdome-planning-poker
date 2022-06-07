<script>
    import SolidButton from '../SolidButton.svelte'
    import Modal from '../Modal.svelte'

    export let handleGoalAdd = () => {}
    export let toggleAddGoal = () => {}
    export let handleGoalRevision = () => {}

    export let goalId = ''
    export let goalName = ''

    function handleSubmit(event) {
        event.preventDefault()

        if (goalId === '') {
            handleGoalAdd(goalName)
        } else {
            handleGoalRevision({
                goalId,
                name: goalName,
            })
        }
        toggleAddGoal()
    }
</script>

<Modal closeModal="{toggleAddGoal}">
    <form on:submit="{handleSubmit}" name="addGoal">
        <div class="mb-4">
            <label
                class="block text-sm text-gray-700 dark:text-gray-400 font-bold mb-2"
                for="goalName"
            >
                Goal Name
            </label>
            <input
                class="bg-gray-100  dark:bg-gray-900 dark:focus:bg-gray-800 border-gray-200 dark:border-gray-600 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 dark:text-gray-400 leading-tight
                focus:outline-none focus:bg-white focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
                id="goalName"
                type="text"
                bind:value="{goalName}"
                placeholder="Enter a goal name"
                name="goalName"
            />
        </div>
        <div class="text-right">
            <div>
                <SolidButton type="submit">Save</SolidButton>
            </div>
        </div>
    </form>
</Modal>
