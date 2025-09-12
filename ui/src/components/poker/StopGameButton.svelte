<script lang="ts">
  import HollowButton from '../global/HollowButton.svelte';
  import DeleteConfirmation from '../global/DeleteConfirmation.svelte';
  import LL from '../../i18n/i18n-svelte';

  interface Props {
    onStopGame: () => void;
    disabled?: boolean;
    testid?: string;
    class?: string;
  }

  let { onStopGame, disabled = false, testid = "stop-game-button", class: klass = '' }: Props = $props();

  let showStopConfirmation = $state(false);

  function toggleStopConfirmation() {
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
  disabled={disabled}
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