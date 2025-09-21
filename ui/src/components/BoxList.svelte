<script lang="ts">
  import HollowButton from './global/HollowButton.svelte';
  import LL from '../i18n/i18n-svelte';
  import { user } from '../stores';
  import { Crown } from 'lucide-svelte';
  import Badge from './global/Badge.svelte';
  import EndStatusBadge from './global/EndStatusBadge.svelte';

  interface Item {
    id: string;
    name: string;
    owner_id?: string;
    endTime?: Date;
    endReason?: string;
    plans?: Array<{ points: string }>;
    [key: string]: any;
  }

  interface Props {
    items?: Array<Item>;
    pageRoute?: string;
    joinBtnText?: string;
    itemType?: string;
    ownerField?: string;
    ownerNameField?: string;
    isAdmin?: boolean;
    showOwner?: boolean;
    showOwnerName?: boolean;
    showFacilitatorIcon?: boolean;
    facilitatorsKey?: string;
    showCompletedStories?: boolean;
    toggleRemove?: Function;
  }

  let {
    items = [],
    pageRoute = '',
    joinBtnText = '',
    itemType = '',
    ownerField = 'owner_id',
    ownerNameField = '',
    isAdmin = false,
    showOwner = true,
    showOwnerName = false,
    showFacilitatorIcon = false,
    facilitatorsKey = 'facilitators',
    showCompletedStories = false,
    toggleRemove = id => () => {},
  }: Props = $props();
</script>

{#each items as item}
  <div
    class="w-full bg-white dark:bg-gray-800 dark:text-white shadow-lg rounded-lg mb-2 border-gray-300 dark:border-gray-700
                        border-b"
  >
    <div class="flex flex-col md:flex-row md:items-center p-4 gap-2">
      <div class="flex-1 min-w-0 font-semibold md:text-xl leading-tight">
        {#if showFacilitatorIcon}
          {#if Array.isArray(item[facilitatorsKey]) && item[facilitatorsKey].includes($user.id)}
            <Crown class="inline-block text-yellow-500" />&nbsp;
          {/if}
        {/if}
        <span data-testid="{itemType}-name">{item.name}</span>
        {#if item.endTime}
          <EndStatusBadge endTime={item.endTime} endReason={item.endReason || 'Ended'} class="inline-block ms-2" />
        {/if}
        {#if showOwnerName}
          <span class="font-semibold md:text-sm text-gray-600 dark:text-gray-400">
            {#if ownerNameField && item[ownerNameField]}
              {item[ownerNameField]}
            {/if}
          </span>
        {/if}
        {#if showOwner}
          <div class="font-semibold md:text-sm text-gray-600 dark:text-gray-400">
            {#if $user.id === item[ownerField]}
              {$LL.owner()}
            {/if}
          </div>
        {/if}
        {#if showCompletedStories}
          <div class="mt-2 font-semibold md:text-sm text-gray-600 dark:text-gray-400">
            {$LL.countPlansPointed({
              totalPointed: item.plans?.filter(p => p.points !== '').length,
              totalPlans: item.plans?.length,
            })}
          </div>
        {/if}
      </div>
      <div class="flex-none w-full md:w-auto flex flex-wrap gap-2 justify-start md:justify-end md:items-center">
        {#if isAdmin}
          <HollowButton onClick={toggleRemove(item.id)} color="red">
            {$LL.remove()}
          </HollowButton>
        {/if}
        <HollowButton href="{pageRoute}/{item.id}">
          {joinBtnText}
        </HollowButton>
      </div>
    </div>
  </div>
{/each}
