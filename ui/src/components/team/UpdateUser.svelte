<script lang="ts">
  import Modal from '../global/Modal.svelte';
  import SolidButton from '../global/SolidButton.svelte';
  import LL from '../../i18n/i18n-svelte';
  import TextInput from '../forms/TextInput.svelte';
  import SelectInput from '../forms/SelectInput.svelte';
  import { onMount, untrack } from 'svelte';

  interface Props {
    toggleUpdate?: any;
    handleUpdate?: any;
    userId?: string;
    userEmail?: string;
    role?: string;
    title?: string;
    notes?: string;
    titleOptions?: string[];
    inheritedTitle?: string;
  }

  let {
    toggleUpdate = () => {},
    handleUpdate = () => {},
    userId = '',
    userEmail = '',
    role = $bindable(''),
    title = $bindable(''),
    notes = $bindable(''),
    titleOptions = [],
    inheritedTitle = '',
  }: Props = $props();

  const roles = ['ADMIN', 'MEMBER'];
  const CUSTOM_TITLE = '__custom__';

  // a title that isn't in the organization's list is a custom one
  let titleChoice = $state(untrack(() => (!title || titleOptions.includes(title) ? title : CUSTOM_TITLE)));

  function onSubmit(e: Event) {
    e.preventDefault();

    const usesTextInput = titleOptions.length === 0 || titleChoice === CUSTOM_TITLE;
    handleUpdate(userId, role, { title: usesTextInput ? title : titleChoice, notes });
  }

  let updateDisabled = $derived(role === '');

  let focusInput: any = $state();
  onMount(() => {
    focusInput?.focus();
  });
</script>

<Modal closeModal={toggleUpdate} ariaLabel={$LL.modalTeamUpdateUser()}>
  <form onsubmit={onSubmit} name="teamUpdateUser">
    <div class="mb-4">
      <label class="block text-gray-700 dark:text-gray-400 font-bold mb-2 disabled" for="userEmail">
        {$LL.userEmail()}
      </label>
      <TextInput bind:this={focusInput} value={userEmail} id="userEmail" name="userEmail" disabled />
    </div>

    <div class="mb-4">
      <label class="text-gray-700 dark:text-gray-400 font-bold mb-2" for="userRole">
        {$LL.role()}
      </label>
      <SelectInput bind:value={role} id="userRole" name="userRole">
        {#each roles as userRole}
          <option value={userRole}>{userRole}</option>
        {/each}
      </SelectInput>
    </div>

    {#if titleOptions.length > 0}
      <div class="mb-4">
        <label class="block text-gray-700 dark:text-gray-400 font-bold mb-2" for="userTitleChoice">
          {$LL.userTitleLabel()}
        </label>
        <SelectInput bind:value={titleChoice} id="userTitleChoice" name="userTitleChoice">
          <option value="">
            {inheritedTitle ? $LL.userTitleInherited({ title: inheritedTitle }) : $LL.userTitleNone()}
          </option>
          {#each titleOptions as option}
            <option value={option}>{option}</option>
          {/each}
          <option value={CUSTOM_TITLE}>{$LL.userTitleCustom()}</option>
        </SelectInput>
      </div>
    {/if}

    {#if titleOptions.length === 0 || titleChoice === CUSTOM_TITLE}
      <div class="mb-4">
        <label class="block text-gray-700 dark:text-gray-400 font-bold mb-2" for="userTitle">
          {titleOptions.length === 0 ? $LL.userTitleLabel() : $LL.userTitleCustomLabel()}
        </label>
        <TextInput
          bind:value={title}
          id="userTitle"
          name="userTitle"
          maxlength={128}
          autocomplete="off"
          placeholder={inheritedTitle
            ? $LL.userTitleInherited({ title: inheritedTitle })
            : $LL.userTitlePlaceholder()}
        />
      </div>
    {/if}

    <div class="mb-4">
      <label class="block text-gray-700 dark:text-gray-400 font-bold mb-2" for="userNotes">
        {$LL.userNotesLabel()}
      </label>
      <textarea
        bind:value={notes}
        id="userNotes"
        name="userNotes"
        rows="4"
        maxlength="1000"
        aria-describedby="userNotesHelp"
        class="block w-full px-5 py-3 rounded-lg outline-none transition-all duration-300 bg-gray-100 dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-purple-500 dark:focus:ring-purple-400"
      ></textarea>
      <p id="userNotesHelp" class="mt-1 text-sm text-gray-500 dark:text-gray-400">
        {$LL.userNotesHelp()}
      </p>
    </div>

    <div>
      <div class="text-right">
        <SolidButton type="submit" disabled={updateDisabled} testid="userupdate-confirm">
          {$LL.userUpdate()}
        </SolidButton>
      </div>
    </div>
  </form>
</Modal>
