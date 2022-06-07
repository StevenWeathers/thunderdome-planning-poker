<script>
  import SolidButton from '../SolidButton.svelte'
  import Modal from '../Modal.svelte'

  export let handleRemove = () => {}
  export let toggleRemove = () => {}
</script>

<Modal closeModal="{toggleRemove}">
    <div class="mb-4">
        <p class="font-bold text-xl text-red-600">
            Are you sure you want to remove this storyboard?
        </p>
    </div>
    <div class="text-right">
        <button
                type="button"
                class="inline-block align-baseline font-bold text-sm text-blue-500
            hover:text-blue-800 mr-4"
                on:click="{toggleRemove}"
        >
            Cancel
        </button>
        <SolidButton onClick="{handleRemove}" color="red">
            Remove Storyboard
        </SolidButton>
    </div>
</Modal>
