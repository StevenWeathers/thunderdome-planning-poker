<script lang="ts">
    import Modal from '../global/Modal.svelte';
    import LL from "../../i18n/i18n-svelte";

    interface Props {
        personas: Array<any>;
        isFacilitator: boolean;
        toggle: any;
        handleAdd: any;
        handleEdit: any;
        handleDelete: any;
    }

    let { personas = [], isFacilitator = false, toggle = () => {}, handleAdd = () => {}, handleEdit = () => {}, handleDelete = () => {} }: Props = $props(); 
  
</script>

<Modal closeModal={toggle}>
    {#each personas as persona}
    <div class="mb-1 w-full">
        <div>
        <span class="font-bold">
            {persona.name}
        </span>
        {#if isFacilitator}
            &nbsp;|&nbsp;
            <button
            onclick={toggle(persona)}
            class="text-orange-500
                                        hover:text-orange-800"
            data-testid="persona-edit"
            >
            {$LL.edit()}
            </button>
            &nbsp;|&nbsp;
            <button
            onclick={handleDelete(persona.id)}
            class="text-red-500
                                        hover:text-red-800"
            data-testid="persona-delete"
            >
            {$LL.delete()}
            </button>
        {/if}
        </div>
        <span class="text-sm">
        {persona.role}
        </span>
    </div>
    {/each}

    {#if isFacilitator}
        <div>Add persona...</div>
    {/if}
</Modal>