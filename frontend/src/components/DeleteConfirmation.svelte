<script>
    import SolidButton from './SolidButton.svelte'
    import Modal from './Modal.svelte'
    import { _ } from '../i18n.js'

    export let confirmText = ''
    export let confirmBtnText = 'Confirm Delete'
    export let permanent = true
    export let handleDelete = () => {}
    export let toggleDelete = () => {}
</script>

<Modal closeModal="{toggleDelete}">
    <div class="mt-8 mb-4">
        <p class="font-bold text-xl text-red-600">
            {confirmText}
            {#if permanent}
                <br />
                {$_('cannotBeUndone')}
            {/if}
        </p>
    </div>
    <div class="text-right">
        <button
            type="button"
            class="inline-block align-baseline font-bold text-sm text-blue-500
            hover:text-blue-800 mr-4"
            on:click="{toggleDelete}"
            data-testid="confirm-cancel"
        >
            {$_('cancel')}
        </button>
        <SolidButton
            onClick="{handleDelete}"
            color="red"
            testid="confirm-confirm"
        >
            {confirmBtnText}
        </SolidButton>
    </div>
</Modal>
