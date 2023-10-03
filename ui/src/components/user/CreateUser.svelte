<script lang="ts">
  import Modal from '../global/Modal.svelte';
  import WarriorRegisterForm from './UserRegisterForm.svelte';

  export let notifications;
  export let toggleCreate = () => {};
  export let handleCreate = () => {};
</script>

<Modal closeModal="{toggleCreate}">
  <WarriorRegisterForm
    handleSubmit="{handleCreate}"
    notifications="{notifications}"
  />
</Modal>
