<script lang="ts">
  import HollowButton from '../global/HollowButton.svelte';
  import DeleteConfirmation from '../global/DeleteConfirmation.svelte';
  import LL from '../../i18n/i18n-svelte';

  interface Props {
    onStopGame: () => void;
    disabled?: boolean;
    endedDate?: Date | string;
    testid?: string;
    class?: string;
  }

  let { onStopGame, disabled = false, endedDate, testid = "stop-game-button", class: klass = '' }: Props = $props();

  // Determine if game is stopped based on endedDate
  let isGameStopped = $derived(endedDate !== null && endedDate !== undefined);
  let computedDisabled = $derived(disabled || isGameStopped);

  let showStopConfirmation = $state(false);

  function toggleStopConfirmation() {
    if (computedDisabled) return;
    showStopConfirmation = !showStopConfirmation;
  }

  function handleStopGame() {
    onStopGame();
    showStopConfirmation = false;
  }
</script>

<HollowButton
  color="orange"
  onClick={toggleStopConfirmation}
  disabled={computedDisabled}
  testid={testid}
  additionalClasses={klass}
>
  {$LL.battleStop()}
</HollowButton>

{#if showStopConfirmation}
  <DeleteConfirmation
    toggleDelete={toggleStopConfirmation}
    handleDelete={handleStopGame}
    confirmText={$LL.stopBattleConfirmText()}
    confirmBtnText={$LL.battleStop()}
  />
{/if}