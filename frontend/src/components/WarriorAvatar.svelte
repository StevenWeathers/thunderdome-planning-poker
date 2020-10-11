<script>
    export let avatarService = ''
    export let warriorId = ''
    export let avatar = ''
    export let width = 48
</script>

{#if avatarService === 'dicebear'}
    <img
        src="https://avatars.dicebear.com/api/{avatar}/{warriorId}.svg?w={width}"
        alt="Placeholder Avatar" />
{:else if avatarService === 'gravatar'}
    <img
        src="https://gravatar.com/avatar/{warriorId}?s={width}&d={avatar}&r=g"
        alt="Placeholder Avatar" />
{:else if avatarService === 'robohash'}
    <img
        src="https://robohash.org/{warriorId}.png?set={avatar}&size={width}x{width}"
        alt="Placeholder Avatar" />
{:else}
    <img
        src="https://api.adorable.io/avatars/{width}/{warriorId}.png"
        alt="Placeholder Avatar" />
{/if}
