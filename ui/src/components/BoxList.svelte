<script lang="ts">
  import HollowButton from './global/HollowButton.svelte';
  import LL from '../i18n/i18n-svelte';
  import { user } from '../stores';
  import { Crown } from 'lucide-svelte';

  interface Props {
    items?: Array<object>;
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
    toggleRemove = id => () => {}
  }: Props = $props();
</script>

{#each items as item}
  <div
    class="w-full bg-white dark:bg-gray-800 dark:text-white shadow-lg rounded-lg mb-2 border-gray-300 dark:border-gray-700
                        border-b"
  >
    <div class="flex flex-wrap items-center p-4">
      <div
        class="w-full md:w-1/2 mb-4 md:mb-0 font-semibold
                            md:text-xl leading-tight"
      >
        {#if showFacilitatorIcon}
          {#if item[facilitatorsKey].includes($user.id)}
            <Crown class="inline-block text-yellow-500" />&nbsp;
          {/if}
        {/if}
        <span data-testid="{itemType}-name">{item.name}</span>
        {#if showOwnerName}
          <span
            class="font-semibold md:text-sm text-gray-600 dark:text-gray-400"
          >
            {#if ownerNameField && item[ownerNameField]}
              {item[ownerNameField]}
            {/if}
          </span>
        {/if}
        {#if showOwner}
          <div
            class="font-semibold md:text-sm text-gray-600 dark:text-gray-400"
          >
            {#if $user.id === item[ownerField]}
              {$LL.owner()}
            {/if}
          </div>
        {/if}
        {#if showCompletedStories}
          <div
            class="font-semibold md:text-sm text-gray-600 dark:text-gray-400"
          >
            {$LL.countPlansPointed({
              totalPointed: item.plans.filter(p => p.points !== '').length,
              totalPlans: item.plans.length,
            })}
          </div>
        {/if}
      </div>
      <div class="w-full md:w-1/2 md:mb-0 md:text-right">
        {#if isAdmin}
          <HollowButton onClick="{toggleRemove(item.id)}" color="red">
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
