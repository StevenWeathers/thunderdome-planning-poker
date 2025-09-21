<script lang="ts">
  import LL from '../../i18n/i18n-svelte';
  import { Eye, Pencil, Trash2 } from 'lucide-svelte';

  interface Props {
    detailsLink?: string;
    editBtnEnabled?: boolean;
    editBtnClickHandler?: any;
    editBtnTestId?: string;
    deleteBtnEnabled?: boolean;
    deleteBtnClickHandler?: any;
    deleteBtnTestId?: string;
    detailsLinkText?: string;
    children?: import('svelte').Snippet;
  }

  let {
    detailsLink = '',
    editBtnEnabled = true,
    editBtnClickHandler = () => {},
    editBtnTestId = 'edit',
    deleteBtnEnabled = true,
    deleteBtnClickHandler = () => {},
    deleteBtnTestId = 'delete',
    detailsLinkText = 'View Details',
    children,
  }: Props = $props();
</script>

<div class="flex gap-2 justify-end items-center">
  {@render children?.()}
  {#if detailsLink !== ''}
    <a href={detailsLink} class="hover:text-blue-500" title={detailsLinkText}>
      <Eye />
      <span class="sr-only">{detailsLinkText}</span>
    </a>
  {/if}
  {#if editBtnEnabled}
    <button onclick={editBtnClickHandler} class="hover:text-green-500" data-testid={editBtnTestId}>
      <span class="sr-only">{$LL.edit()}</span>
      <Pencil />
    </button>
  {/if}
  {#if deleteBtnEnabled}
    <button onclick={deleteBtnClickHandler} class="hover:text-red-500" data-testid={deleteBtnTestId}>
      <span class="sr-only">{$LL.delete()}</span>
      <Trash2 />
    </button>
  {/if}
</div>
