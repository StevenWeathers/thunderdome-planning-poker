<script lang="ts">
  import Modal from '../global/Modal.svelte';
  import UserRegisterForm from './UserRegisterForm.svelte';

  import type { NotificationService } from '../../types/notifications';

  interface Props {
    notifications: NotificationService;
    toggleCreate?: any;
    handleCreate?: any;
  }

  let { notifications, toggleCreate = () => {}, handleCreate = () => {} }: Props = $props();
</script>

<Modal closeModal={toggleCreate}>
  <div class="pt-10">
    <UserRegisterForm
      handleFullAccountRegistration={handleCreate}
      notifications={notifications}
      fullOnly={true}
      isAdmin={true}
    />
  </div>
</Modal>
