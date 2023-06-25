<script lang="ts">
    import SolidButton from './SolidButton.svelte'
    import Modal from './Modal.svelte'
    import LL from '../i18n/i18n-svelte'

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
                {$LL.cannotBeUndone()}
            {/if}
        </p>
    </div>
    <div class="text-right">
        <button
            type="button"
            class="inline-block align-baseline font-bold text-sm text-blue-500
            hover:text-blue-800 me-4"
            on:click="{toggleDelete}"
            data-testid="confirm-cancel"
        >
            {$LL.cancel()}
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
