<script lang="ts">
    import Badge from './Badge.svelte';
    import { LL } from '../../i18n/i18n-svelte';
    
    interface Props {
        endTime: string;
        endReason: string;
        class: string;
    }
    
    let { endTime, endReason, class: klass }: Props = $props();
    
    function endReasonColor(reason: string) {
        switch (reason.toLowerCase()) {
        case 'completed':
            return 'green';
        case 'cancelled':
        case 'abandoned':
            return 'red';
        default:
            return 'gray';
        }
    }
    

    function getLocaleReason(reason: string) {
        switch (reason?.toLowerCase()) {
            case 'completed':
                return $LL.completed();
            case 'cancelled':
                return $LL.cancelled();
            case 'abandoned':
                return $LL.abandoned();
            default:
                // 'ended' does not exist, use a generic fallback
                return $LL.ended();
        }
    }

    function badgeTitle(endTime: string, reason: string) {
        return `${getLocaleReason(reason)} on ${new Date(endTime).toLocaleString()}`;
    }
</script>

<Badge
    label={getLocaleReason(endReason)}
    title={badgeTitle(endTime, endReason)}
    class="text-sm {klass}"
    color={endReason ? endReasonColor(endReason) : 'gray'}
/>