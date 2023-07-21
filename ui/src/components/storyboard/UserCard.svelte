<script lang="ts">
  import UserAvatar from '../user/UserAvatar.svelte';
  import { warrior } from '../../stores';
  import LL from '../../i18n/i18n-svelte';

  export let user = {};
  export let showBorder = 'true';
  export let facilitators = [];
  export let handleAddFacilitator = () => {};
  export let handleRemoveFacilitator = () => {};

  $: borderClasses = showBorder === 'true' ? 'border-b border-gray-500' : '';
</script>

<div
  class="{borderClasses} p-2 flex items-center"
  data-testId="userCard"
  data-userName="{user.name}"
>
  <UserAvatar
    warriorId="{user.id}"
    avatar="{user.avatar}"
    gravatarHash="{user.gravatarHash}"
  />
  <p
    class="ms-2 text-l font-bold leading-tight truncate"
    data-testId="userName"
    title="{user.name}"
  >
    {user.name}
    {#if facilitators.includes(user.id)}
      <div class="text-indigo-500 dark:text-violet-400">
        {$LL.facilitator()}
        {#if facilitators.includes($warrior.id)}
          <button
            class="text-red-500 text-sm"
            on:click="{handleRemoveFacilitator(user.id)}">{$LL.remove()}</button
          >
        {/if}
      </div>
    {:else if facilitators.includes($warrior.id)}
      <div>
        <button
          class="text-blue-500 dark:text-sky-400 text-sm"
          on:click="{handleAddFacilitator(user.id)}"
          >{$LL.makeFacilitator()}</button
        >
      </div>
    {/if}
  </p>
</div>
