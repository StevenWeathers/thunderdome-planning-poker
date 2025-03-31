<script lang="ts">
  import UserAvatar from '../user/UserAvatar.svelte';
  import { user as sessionUser } from '../../stores';
  import LL from '../../i18n/i18n-svelte';

  interface Props {
    user?: any;
    showBorder?: boolean;
    facilitators?: any;
    handleAddFacilitator?: any;
    handleRemoveFacilitator?: any;
  }

  let {
    user = {},
    showBorder = true,
    facilitators = [],
    handleAddFacilitator = () => {},
    handleRemoveFacilitator = () => {}
  }: Props = $props();

  let borderClasses = $derived(showBorder ? 'border-b border-gray-500' : '');
</script>

<div
  class="{borderClasses} p-2 flex items-center"
  data-testId="userCard"
  data-userName="{user.name}"
>
  <UserAvatar
    warriorId={user.id}
    avatar={user.avatar}
    gravatarHash={user.gravatarHash}
    userName={user.name}
  />
  <div
    class="ms-2 text-l font-bold leading-tight truncate"
    data-testId="userName"
    title="{user.name}"
  >
    {user.name}
    {#if facilitators.includes(user.id)}
      <div class="text-indigo-500 dark:text-violet-400">
        {$LL.facilitator()}
        {#if facilitators.includes($sessionUser.id)}
          <button
            class="text-red-500 text-sm"
            onclick={handleRemoveFacilitator(user.id)}>{$LL.remove()}</button
          >
        {/if}
      </div>
    {:else if facilitators.includes($sessionUser.id)}
      <div>
        <button
          class="text-blue-500 dark:text-sky-400 text-sm"
          onclick={handleAddFacilitator(user.id)}
          >{$LL.makeFacilitator()}</button
        >
      </div>
    {/if}
    </div>
</div>
