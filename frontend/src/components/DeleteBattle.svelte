<script>
    import SolidButton from './SolidButton.svelte'
    import Modal from './Modal.svelte'
    import { _ } from '../i18n'

    export let handleDelete = () => {}
    export let toggleDelete = () => {}
</script>

<Modal closeModal="{toggleDelete}">
    <div class="mb-4">
        <p class="font-bold text-xl text-red-600">
            Are you sure you want to delete this battle?
            <br />
            {$_('cannotBeUndone')}
        </p>
    </div>
    <div class="text-right" data-testid="confirm-actions">
        <button
            type="button"
            class="inline-block align-baseline font-bold text-sm text-blue-500
            hover:text-blue-800 mr-4"
            on:click="{toggleDelete}"
        >
            {$_('cancel')}
        </button>
        <SolidButton onClick="{handleDelete}" color="red">
            Delete Battle
        </SolidButton>
    </div>
</Modal>
