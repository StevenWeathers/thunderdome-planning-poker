<script lang="ts">
  import LL from '../../i18n/i18n-svelte';
  import EditIcon from '../icons/EditIcon.svelte';
  import TrashIcon from '../icons/TrashIcon.svelte';
  import EyeIcon from '../icons/EyeIcon.svelte';

  export let detailsLink = '';
  export let editBtnEnabled = true;
  export let editBtnClickHandler = () => {};
  export let editBtnTestId = 'edit';
  export let deleteBtnEnabled = true;
  export let deleteBtnClickHandler = () => {};
  export let deleteBtnTestId = 'delete';
</script>

<div class="flex gap-2 justify-end items-center">
  <slot />
  {#if detailsLink !== ''}
    <a href="{detailsLink}" class="hover:text-blue-500">
      <EyeIcon />
      <span class="sr-only">View Details</span>
    </a>
  {/if}
  {#if editBtnEnabled}
    <button
      on:click="{editBtnClickHandler}"
      class="hover:text-green-500"
      data-testid="{editBtnTestId}"
    >
      <span class="sr-only">{$LL.edit()}</span>
      <EditIcon />
    </button>
  {/if}
  {#if deleteBtnEnabled}
    <button
      on:click="{deleteBtnClickHandler}"
      class="hover:text-red-500"
      data-testid="{deleteBtnTestId}"
    >
      <span class="sr-only">{$LL.delete()}</span>
      <TrashIcon />
    </button>
  {/if}
</div>
