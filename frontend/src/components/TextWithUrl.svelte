<script>
    import Autolinker from 'autolinker'

    export let text = ''

    const { AllowTextUrls } = appConfig

</script>

{#if AllowTextUrls }
    {@html Autolinker.link(text)}
{:else}
    {text}
{/if}