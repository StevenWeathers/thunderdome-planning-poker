<script lang="ts">
  import HollowButton from '../global/HollowButton.svelte';
  import LL from '../../i18n/i18n-svelte';

  interface Props {
    isFacilitator: boolean;
    gameEnded: boolean;
    onEditGame: () => void;
    onStopGame: () => void;
    onDeleteGame: () => void;
    onAbandonGame: () => void;
  }

  let { 
    isFacilitator, 
    gameEnded, 
    onEditGame, 
    onStopGame, 
    onDeleteGame, 
    onAbandonGame 
  }: Props = $props();
</script>

<div class="bg-white dark:bg-gray-800 shadow-lg p-4 mb-4 rounded-lg">
  <h3 class="text-xl font-semibold font-rajdhani uppercase mb-4 text-gray-800 dark:text-white">
    {$LL.gameControls()}
  </h3>
  
  <div class="space-y-2">
    {#if isFacilitator}
      <div class="flex flex-wrap gap-2">
        <HollowButton
          color="blue"
          onClick={onEditGame}
          testid="battle-edit"
        >
          {$LL.battleEdit()}
        </HollowButton>
        
        {#if !gameEnded}
          <HollowButton
            color="orange"
            onClick={onStopGame}
            testid="battle-stop"
          >
            {$LL.battleStop()}
          </HollowButton>
        {/if}
        
        <HollowButton
          color="red"
          onClick={onDeleteGame}
          testid="battle-delete"
        >
          {$LL.battleDelete()}
        </HollowButton>
      </div>
    {:else}
      <HollowButton
        color="red"
        onClick={onAbandonGame}
        testid="battle-abandon"
      >
        {$LL.battleAbandon()}
      </HollowButton>
    {/if}
  </div>
</div>