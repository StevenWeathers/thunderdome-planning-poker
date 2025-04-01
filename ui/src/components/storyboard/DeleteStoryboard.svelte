<script lang="ts">
  import SolidButton from '../global/SolidButton.svelte';
  import Modal from '../global/Modal.svelte';

  interface Props {
    handleDelete?: any;
    toggleDelete?: any;
  }

  let { handleDelete = () => {}, toggleDelete = () => {} }: Props = $props();
</script>

<Modal closeModal={toggleDelete}>
  <div class="mb-4">
    <p class="font-bold text-xl text-red-600">
      Are you sure you want to delete this Storyboard?
      <br />
      This cannot be undone.
    </p>
  </div>
  <div class="text-right">
    <button
      type="button"
      class="inline-block align-baseline font-bold text-sm text-blue-500
            hover:text-blue-800 me-4"
      onclick={toggleDelete}
      data-testid="confirm-cancel"
    >
      Cancel
    </button>
    <SolidButton onClick={handleDelete} color="red" testid="confirm-confirm">
      Delete Storyboard
    </SolidButton>
  </div>
</Modal>
