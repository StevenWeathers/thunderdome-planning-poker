<script lang="ts">
  import SolidButton from '../global/SolidButton.svelte';
  import Modal from '../global/Modal.svelte';
  import LL from '../../i18n/i18n-svelte';
  import HollowButton from '../global/HollowButton.svelte';
  import UserAvatar from '../user/UserAvatar.svelte';
  import SelectInput from '../forms/SelectInput.svelte';
  import Checkbox from '../forms/Checkbox.svelte';
  import { Trash2 } from 'lucide-svelte';

  interface Props {
    toggleEdit?: any;
    handleEdit?: any;
    handleDelete?: any;
    handleAssigneeAdd?: any;
    handleAssigneeRemove?: any;
    assignableUsers?: any;
    action?: any;
  }

  let {
    toggleEdit = () => {},
    handleEdit = action => {},
    handleDelete = () => {},
    handleAssigneeAdd = (retroId, actionId, userId) => {},
    handleAssigneeRemove = (retroId, actionId, userId) => () => {},
    assignableUsers = [],
    action = {
    id: '',
    retroId: '',
    content: '',
    completed: false,
    assignees: [],
  }
  }: Props = $props();

  let selectedAssignee = $state('');

  let editAction = $state({
    id: action.id,
    retroId: action.retroId,
    content: action.content,
    completed: action.completed,
  });

  const handleSubmit = e => {
    e.preventDefault();

    handleEdit(editAction);
  };
  const addAssignee = () => {
    handleAssigneeAdd(action.retroId, action.id, selectedAssignee);
    selectedAssignee = '';
  };
</script>

<Modal closeModal={toggleEdit} widthClasses="md:w-2/3 lg:w-3/5 xl:w-1/2">
  <form onsubmit={handleSubmit}>
    <div class="mb-4">
      <label
        class="block text-gray-700 dark:text-gray-400 text-sm font-bold mb-2"
        for="actionItem"
      >
        {$LL.actionItem()}
      </label>
      <div class="control">
        <input
          bind:value="{editAction.content}"
          placeholder="{$LL.actionItemPlaceholder()}"
          class="dark:bg-gray-800 border-gray-300 dark:border-gray-700 border-2 appearance-none rounded py-2
                px-3 text-gray-700 dark:text-gray-400 leading-tight focus:outline-none
                focus:bg-white dark:focus:bg-gray-700 focus:border-indigo-500 dark:focus:border-yellow-400 w-full"
          id="actionItem"
          name="actionItem"
          type="text"
          required
        />
      </div>
    </div>

    <div class="mb-4">
      <div class="control">
        <div class="flex-shrink">
          <Checkbox
            id="Completed"
            bind:checked="{editAction.completed}"
            label="{$LL.completed()}"
          />
        </div>
      </div>
    </div>

    <div class="flex w-full pt-4">
      <div class="w-1/2">
        <HollowButton color="red" onClick={handleDelete(editAction)}
          >{$LL.delete()}</HollowButton
        >
      </div>
      <div class="w-1/2 text-right">
        <SolidButton type="submit">{$LL.save()}</SolidButton>
      </div>
    </div>
  </form>

  <div class="mt-4 pt-2 border-t border-gray-400 dark:border-gray-700">
    <div class="block text-gray-700 dark:text-gray-400 font-bold mb-4">
      {$LL.assignees()}
    </div>
    <div class="flex w-full gap-4">
      <div class="w-2/3">
        <SelectInput
          bind:value="{selectedAssignee}"
          id="assignee"
          name="assignee"
        >
          <option value="" disabled>{$LL.assigneeSelectPlaceholder()}</option>
          {#each assignableUsers as user}
            <option value="{user.id}">
              {user.name}
            </option>
          {/each}
        </SelectInput>
      </div>
      <div class="w-1/3">
        <HollowButton
          onClick={addAssignee}
          disabled={selectedAssignee === ''}
        >
          {$LL.assigneeAdd()}
        </HollowButton>
      </div>
    </div>
    {#if action.assignees.length}
      <div class="grid grid-cols-2 gap-4 mt-4">
        {#each action.assignees as assignee}
          <div class="flex text-gray-700 dark:text-gray-400 mb-2">
            <div class="w-1/4">
              <UserAvatar
                warriorId="{assignee.id}"
                gravatarHash="{assignee.gravatarHash}"
                avatar="{assignee.avatar}"
                userName="{assignee.name}"
                class="inline-block me-2"
              />
            </div>
            <div class="w-2/4 text-lg">{assignee.name}</div>
            <div class="w-1/4 text-right">
              <HollowButton
                color="red"
                onClick={handleAssigneeRemove(
                  action.retroId,
                  action.id,
                  assignee.id,
                )}
              >
                <Trash2 />
              </HollowButton>
            </div>
          </div>
        {/each}
      </div>
    {/if}
  </div>
</Modal>
