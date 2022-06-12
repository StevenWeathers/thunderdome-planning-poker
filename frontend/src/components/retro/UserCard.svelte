<script>
    import UserAvatar from '../user/UserAvatar.svelte'

    export let user = {}
    export let votes = []
    export let maxVotes = 3

    $: reachedMaxVotes =
        votes && votes.filter(v => v.userId === user.id).length === maxVotes
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
        class="mx-auto mb-2"
    />
    <div
        class="text-l font-bold leading-tight truncate dark:text-white"
        data-testId="userName"
        title="{user.name}"
    >
        {user.name}
        {#if reachedMaxVotes}
            <div class="text-green-600 dark:text-green-400">All votes in</div>
        {/if}
    </div>
</div>
