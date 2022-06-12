<script>
    export let user = {}
    export let showBorder = 'true'
    import UserAvatar from '../user/UserAvatar.svelte'
</script>

<div
    class="shrink text-center px-2 "
    data-testId="userCard"
    data-userName="{user.name}"
>
    <UserAvatar
        warriorId="{user.id}"
        avatar="{user.avatar}"
        gravatarHash="{user.gravatarHash}"
        class="mx-auto"
    />
    <div
        class="text-l font-bold leading-tight truncate dark:text-white"
        data-testId="userName"
        title="{user.name}"
    >
        {user.name}
    </div>
</div>
