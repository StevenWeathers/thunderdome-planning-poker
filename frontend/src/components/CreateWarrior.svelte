<script>
    import Modal from './Modal.svelte'
    import WarriorRegisterForm from './WarriorRegisterForm.svelte'

    export let notifications
    export let toggleCreate = () => {}
    export let handleCreate = () => {}
</script>

<Modal closeModal={toggleCreate}>
    <WarriorRegisterForm
        handleSubmit="{handleCreate}"
        {notifications} />
</Modal>