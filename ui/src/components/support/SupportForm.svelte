<script lang="ts">
  import TextInput from '../global/TextInput.svelte';
  import SolidButton from '../global/SolidButton.svelte';
  import LL from '../../i18n/i18n-svelte';
  import { user } from '../../stores';

  export let xfetch = async (url, ...options) => {};
  export let notifications = () => {};
  export let eventTag = (eventName: string, type: string, result: string) => {};

  let submitted = false;
  let userName = $user.id && $user.name !== '' ? $user.name : '';
  let userEmail = $user.id && $user.email !== '' ? $user.email : '';
  let message = '';

  function handleSubmit(e) {
    e.preventDefault();

    if (userName === '') {
      notifications.danger('Name field required');
      eventTag('support_form_name_invalid', 'engagement', 'failure');
      return false;
    }

    if (userEmail === '') {
      notifications.danger('Email field required');
      eventTag('support_form_email_invalid', 'engagement', 'failure');
      return false;
    }

    if (message === '') {
      notifications.danger('Message field required');
      eventTag('support_form_message_invalid', 'engagement', 'failure');
      return false;
    }

    const body = {
      user_name: userName,
      user_email: userEmail,
      user_question: message,
    };

    xfetch(`/api/support-tickets`, {
      body,
      method: 'POST',
    })
      .then((res: any) => res.json())
      .then(function () {
        submitted = true;
        eventTag('support_form', 'engagement', 'success');
      })
      .catch(function () {
        notifications.danger('failed to submit support form');

        eventTag('support_form', 'engagement', 'failure');
      });
  }

  $: submitDisabled = userName === '' || userEmail === '' || message === '';
</script>

{#if !submitted}
  <form on:submit="{handleSubmit}" name="supportForm">
    <div class="mb-4">
      <label
        class="block text-gray-700 dark:text-gray-400 font-bold mb-2"
        for="yourName"
      >
        {$LL.name()}
      </label>
      <TextInput
        bind:value="{userName}"
        placeholder="{$LL.userNamePlaceholder()}"
        id="yourName"
        name="yourName"
        required
      />
    </div>

    <div class="mb-4">
      <label
        class="block text-gray-700 dark:text-gray-400 font-bold mb-2"
        for="yourEmail"
      >
        {$LL.email()}
      </label>
      <TextInput
        bind:value="{userEmail}"
        placeholder="{$LL.enterYourEmail()}"
        id="yourEmail"
        name="yourEmail"
        type="email"
        disabled="{$user.rank && $user.rank !== 'GUEST' ? 'disabled' : null}"
        required
      />
    </div>

    <div class="mb-4">
      <label
        class="block text-gray-700 dark:text-gray-400 font-bold mb-2"
        for="message"
      >
        Message
      </label>
      <textarea
        rows="4"
        class="bg-gray-100 dark:bg-gray-900 dark:focus:bg-gray-800 border-gray-200 dark:border-gray-600 border-2 appearance-none
                            rounded w-full py-2 px-3 text-gray-700 dark:text-gray-400 leading-tight
                            focus:outline-none focus:bg-white focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400 mb-2"
        placeholder="How can we assist you?"
        bind:value="{message}"
        required></textarea>
    </div>

    <div class="text-right">
      <SolidButton type="submit" disabled="{submitDisabled}"
        >Submit
      </SolidButton>
    </div>
  </form>
{:else}
  <div
    class="font-bold p-8 text-xl md:text-2xl text-green-600 dark:text-lime-400"
  >
    Thank you for contacting us, we will get back to you shortly.
  </div>
{/if}
