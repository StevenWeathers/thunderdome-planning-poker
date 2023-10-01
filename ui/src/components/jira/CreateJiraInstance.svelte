<script lang="ts">
  import SolidButton from '../SolidButton.svelte';
  import Modal from '../Modal.svelte';
  import LL from '../../i18n/i18n-svelte';
  import { user } from '../../stores';
  import TextInput from '../TextInput.svelte';

  export let handleCreate = () => {};
  export let toggleClose = () => {};
  export let eventTag = () => {};
  export let xfetch = () => {};
  export let notifications;

  let host = '';
  let client_mail = '';
  let access_token = '';

  function handleSubmit(event) {
    event.preventDefault();

    if (host === '') {
      notifications.danger('Host field required');
      eventTag('create_jira_instance_host_invalid', 'engagement', 'failure');
      return false;
    }

    if (client_mail === '') {
      notifications.danger('Host client_mail required');
      eventTag(
        'create_jira_instance_client_mail_invalid',
        'engagement',
        'failure',
      );
      return false;
    }

    if (access_token === '') {
      notifications.danger('Host access_token required');
      eventTag(
        'create_jira_instance_access_token_invalid',
        'engagement',
        'failure',
      );
      return false;
    }

    const body = {
      host,
      client_mail,
      access_token,
    };

    xfetch(`/api/users/${$user.id}/jira-instances`, { body })
      .then(res => res.json())
      .then(function () {
        handleCreate();
        toggleClose();
        eventTag('create_jira_instance', 'engagement', 'success');
      })
      .catch(function () {
        notifications.danger('failed to create jira instance');

        eventTag('create_jira_instance', 'engagement', 'failure');
      });
  }
</script>

<Modal closeModal="{toggleClose}">
  <form on:submit="{handleSubmit}" name="createjirainstance">
    <div class="mb-4">
      <label class="block dark:text-gray-400 font-bold mb-2" for="host">
        Host
      </label>
      <TextInput
        id="host"
        name="host"
        bind:value="{host}"
        placeholder="Enter the Jira Hostname..."
        required
      />
    </div>
    <div class="mb-4">
      <label class="block dark:text-gray-400 font-bold mb-2" for="client_mail">
        Jira User Email
      </label>
      <TextInput
        id="client_mail"
        name="client_mail"
        bind:value="{client_mail}"
        placeholder="Enter your Jira user email..."
        required
      />
    </div>
    <div class="mb-4">
      <label class="block dark:text-gray-400 font-bold mb-2" for="access_token">
        API Access Token
      </label>
      <TextInput
        id="access_token"
        name="access_token"
        bind:value="{access_token}"
        placeholder="Enter your Jira API Access Token..."
        required
      />
    </div>
    <div class="text-right">
      <div>
        <SolidButton type="submit">
          {$LL.create()}
        </SolidButton>
      </div>
    </div>
  </form>
</Modal>
