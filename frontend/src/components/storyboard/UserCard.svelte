<script>
    export let user = {}
    export let showBorder = 'true'
    import UserAvatar from '../user/UserAvatar.svelte'

    $: borderClasses = showBorder === 'true' ? 'border-b border-gray-500' : ''
</script>

<div
    class="{borderClasses} p-2 flex items-center"
    data-testId="userCard"
    data-userName="{user.name}"
>
    <UserAvatar
        warriorId="{user.id}"
        avatar="{user.avatar}"
        gravatarHash="{user.gravatarHash}"
    />
    <p
        class="ml-2 text-l font-bold leading-tight truncate"
        data-testId="userName"
        title="{user.name}"
    >
        {user.name}
    </p>
</div>
