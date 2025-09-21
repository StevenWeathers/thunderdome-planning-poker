<script lang="ts">
  import Modal from '../global/Modal.svelte';
  import TextInput from '../forms/TextInput.svelte';
  import type { ApiClient } from '../../types/apiclient';
  import type { NotificationService } from '../../types/notifications';
  import Checkbox from '../forms/Checkbox.svelte';
  import HollowButton from '../global/HollowButton.svelte';
  import type { User } from '../../types/user';

  interface Props {
    xfetch: ApiClient;
    notifications: NotificationService;
    toggleUpdate?: any;
    handleUpdate?: any;
    ticket: any;
    adminUsers: User[];
  }

  let {
    xfetch,
    notifications,
    toggleUpdate = () => {},
    handleUpdate = () => {},
    ticket = {
      id: '',
      fullName: '',
      email: '',
      inquiry: '',
      assignedTo: null,
      notes: null,
      resolvedAt: null,
      resolvedBy: null,
    },
    adminUsers = [],
  }: Props = $props();

  let fullName = $state(ticket.fullName);
  let email = $state(ticket.email);
  let inquiry = $state(ticket.inquiry);
  let assignedTo = $state(ticket.assignedTo);
  let notes = $state(ticket.notes);
  let markResolved = $state(
    ticket.resolvedAt !== '' && ticket.resolvedAt !== null,
  );
  

  function onSubmit() {
    handleUpdate(ticket.id, { fullName, email, inquiry, assignedTo, notes, markResolved });
  }
</script>

<Modal closeModal={toggleUpdate}>
  <h3
    class="text-xl font-medium leading-6 text-gray-900 dark:text-gray-100 mb-6"
  >
    Update Support Ticket
  </h3>
  <form
    onsubmit={(event: SubmitEvent) => {
      event.preventDefault();
      onSubmit();
    }}
    class="space-y-6"
  >
    <div>
      <label
        for="fullName"
        class="block font-medium text-gray-700 dark:text-gray-200 mb-2"
        >Full Name <span class="text-red-500">*</span></label
      >
      <TextInput
        id="fullName"
        name="fullName"
        bind:value={fullName}
        required
        placeholder="Enter full name"
      />
    </div>
    <div>
      <label
        for="email"
        class="block font-medium text-gray-700 dark:text-gray-200 mb-2"
        >Email <span class="text-red-500">*</span></label
      >
      <TextInput
        id="email"
        name="email"
        bind:value={email}
        required
        placeholder="Enter email"
        type="email"
      />
    </div>
    <div>
      <label
        for="inquiry"
        class="block font-medium text-gray-700 dark:text-gray-200 mb-2"
        >Inquiry <span class="text-red-500">*</span></label
      >
      <textarea
        class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700 dark:text-gray-200"
        rows="4"
        id="inquiry"
        name="inquiry"
        required
        placeholder="Enter inquiry"
        bind:value={inquiry}
      ></textarea>
    </div>
    <div>
      <label
        for="assignedTo"
        class="block font-medium text-gray-700 dark:text-gray-200 mb-2"
        >Assigned To</label
      >
        <select
            id="assignedTo"
            name="assignedTo"
            bind:value={assignedTo}
            class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700 dark:text-gray-200"
        >
          <option value="" selected={assignedTo === null || assignedTo === ''}>Select an admin user</option>
          {#each adminUsers as user}
            <option value={user.id} selected={assignedTo === user.id}>{user.name}</option>
          {/each}
        </select>
    </div>
    <div>
      <label
        for="notes"
        class="block font-medium text-gray-700 dark:text-gray-200 mb-2"
        >Notes</label
      >
      <textarea
        class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700 dark:text-gray-200"
        rows="4"
        id="notes"
        name="notes"
        bind:value={notes}
        placeholder="Enter any internal notes"
      ></textarea>
    </div>
    <div>
      <Checkbox
        bind:checked={markResolved}
        label="Mark as Resolved"
        id="markResolved"
        name="markResolved"
      />
    </div>
    <div class="flex justify-end gap-2 pt-4">
        <HollowButton onClick={onSubmit}>Update Ticket</HollowButton>
    </div>
  </form>
</Modal>
