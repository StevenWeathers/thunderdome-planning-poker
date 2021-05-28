<script>
    import SolidButton from './SolidButton.svelte'
    import Modal from './Modal.svelte'
    import { _ } from '../i18n'

    export let handleDeleteAccount = () => {}
    export let toggleDeleteAccount = () => {}
</script>

<Modal closeModal="{toggleDeleteAccount}">
    <div class="mb-4">
        <p class="font-bold text-xl text-red-600">
            {$_('pages.warriorProfile.delete.warningStatement')}
        </p>
    </div>
    <div class="text-right">
        <button
            type="button"
            class="inline-block align-baseline font-bold text-sm text-blue-500
            hover:text-blue-800 mr-4"
            on:click="{toggleDeleteAccount}">
            {$_('pages.warriorProfile.delete.cancelButton')}
        </button>
        <SolidButton onClick="{handleDeleteAccount}" color="red">
            {$_('pages.warriorProfile.delete.confirmButton')}
        </SolidButton>
    </div>
</Modal>
