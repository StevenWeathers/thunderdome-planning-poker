<script lang="ts">
  import Modal from '../global/Modal.svelte';
  import UserRegisterForm from './UserRegisterForm.svelte';
  import LL from '../../i18n/i18n-svelte';

  import type { NotificationService } from '../../types/notifications';

  interface Props {
    notifications: NotificationService;
    toggleCreate?: any;
    handleCreate?: any;
  }

  let { notifications, toggleCreate = () => {}, handleCreate = () => {} }: Props = $props();
</script>

<Modal closeModal={toggleCreate} ariaLabel={$LL.modalCreateUser()}>
  <div class="pt-10">
    <UserRegisterForm
      handleFullAccountRegistration={handleCreate}
      handleGuestRegistration={() => {}}
      {notifications}
      fullOnly={true}
      isAdmin={true}
    />
  </div>
</Modal>
