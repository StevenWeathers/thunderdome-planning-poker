<script>
  import SolidButton from '../SolidButton.svelte'
  import Modal from '../Modal.svelte'

  export let handleGoalAdd = () => {}
  export let toggleAddGoal = () => {}
  export let handleGoalRevision = () => {}

  export let goalId = ''
  export let goalName = ''

  function handleSubmit (event) {
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
            <label class="block text-sm font-bold mb-2" for="goalName">
                Goal Name
            </label>
            <input
                    class="bg-gray-200 border-gray-200 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 leading-tight
                focus:outline-none focus:bg-white focus:border-purple-500"
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
