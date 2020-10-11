<script>
    import { _ } from '../i18n'

    export let avatarService = ''
    export let warriorId = ''
    export let avatar = ''
    export let width = 48
</script>

{#if avatarService === 'dicebear'}
    <img
        src="https://avatars.dicebear.com/api/{avatar}/{warriorId}.svg?w={width}"
        alt="{$_('avtarAltText')}" />
{:else if avatarService === 'gravatar'}
    <img
        src="https://gravatar.com/avatar/{warriorId}?s={width}&d={avatar}&r=g"
        alt="{$_('avtarAltText')}" />
{:else if avatarService === 'robohash'}
    <img
        src="https://robohash.org/{warriorId}.png?set={avatar}&size={width}x{width}"
        alt="{$_('avtarAltText')}" />
{:else}
    <img
        src="https://api.adorable.io/avatars/{width}/{warriorId}.png"
        alt="{$_('avtarAltText')}" />
{/if}
