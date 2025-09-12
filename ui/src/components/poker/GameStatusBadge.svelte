<script lang="ts">
  import Badge from '../global/Badge.svelte';
  import LL from '../../i18n/i18n-svelte';

  interface Props {
    isActive: boolean;
    endedDate?: Date;
    class?: string;
  }

  let { isActive, endedDate, class: klass = '' }: Props = $props();

  function formatEndedDate(date: Date): string {
    return new Intl.DateTimeFormat('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    }).format(date);
  }

  let badgeLabel = $derived(() => {
    if (isActive) {
      return $LL.gameActive();
    } else if (endedDate) {
      return `${$LL.gameStopped()} ${formatEndedDate(endedDate)}`;
    } else {
      return $LL.gameStopped();
    }
  });

  let badgeColor = $derived(isActive ? 'green' : 'orange');
</script>

<Badge 
  label={badgeLabel} 
  color={badgeColor} 
  testId="game-status-badge"
  class="{klass}"
/>