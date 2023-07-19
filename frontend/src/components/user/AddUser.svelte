<script lang="ts">
  import Modal from '../Modal.svelte';
  import DownCarrotIcon from '../icons/ChevronDown.svelte';
  import SolidButton from '../SolidButton.svelte';
  import LL from '../../i18n/i18n-svelte';

  export let toggleAdd = () => {};
  export let handleAdd = () => {};

  const roles = ['ADMIN', 'MEMBER'];
  let userEmail = '';
  let role = '';

  function onSubmit(e) {
    e.preventDefault();

    handleAdd(userEmail, role);
  }

  $: createDisabled = userEmail === '' || role === '';
</script>

<Modal closeModal="{toggleAdd}">
  <form on:submit="{onSubmit}" name="teamAddUser">
    <div class="mb-4">
      <label
        class="block text-gray-700 dark:text-gray-400 font-bold mb-2"
        for="userEmail"
      >
        {$LL.userEmail()}
      </label>
      <input
        bind:value="{userEmail}"
        placeholder="{$LL.userEmailPlaceholder()}"
        class="bg-gray-100 dark:bg-gray-900 border-gray-200 dark:border-gray-800 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 dark:text-gray-300 leading-tight
                focus:outline-none focus:bg-white dark:focus:bg-gray-700 focus:border-indigo-500 focus:caret-indigo-500
                dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
        id="userEmail"
        name="userEmail"
        required
      />
    </div>

    <div class="mb-4">
      <label
        class="text-gray-700 dark:text-gray-400 font-bold mb-2"
        for="userRole"
      >
        {$LL.role()}
      </label>
      <div class="relative">
        <select
          bind:value="{role}"
          class="bg-gray-100 dark:bg-gray-900 border-gray-200 dark:border-gray-800 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 dark:text-gray-300 leading-tight
                focus:outline-none focus:bg-white dark:focus:bg-gray-700 focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
          id="userRole"
          name="userRole"
        >
          <option value="">{$LL.rolePlaceholder()}</option>
          {#each roles as userRole}
            <option value="{userRole}">{userRole}</option>
          {/each}
        </select>
        <div
          class="pointer-events-none absolute inset-y-0 end-0 flex
                    items-center px-2 text-gray-700 dark:text-gray-400"
        >
          <DownCarrotIcon />
        </div>
      </div>
    </div>

    <div>
      <div class="text-right">
        <SolidButton
          type="submit"
          disabled="{createDisabled}"
          testid="useradd-confirm"
        >
          {$LL.userAdd()}
        </SolidButton>
      </div>
    </div>
  </form>
</Modal>
