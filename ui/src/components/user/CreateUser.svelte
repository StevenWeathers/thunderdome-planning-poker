<script lang="ts">
  import Modal from '../global/Modal.svelte';
  import UserRegisterForm from './UserRegisterForm.svelte';

  export let notifications;
  export let toggleCreate = () => {};
  export let handleCreate = () => {};
</script>

<Modal closeModal="{toggleCreate}">
  <div class="pt-10">
    <UserRegisterForm
      handleFullAccountRegistration="{handleCreate}"
      notifications="{notifications}"
      fullOnly="{true}"
      isAdmin="{true}"
    />
  </div>
</Modal>
