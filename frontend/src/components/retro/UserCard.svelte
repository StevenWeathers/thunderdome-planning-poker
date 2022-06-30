<script>
    import UserAvatar from '../user/UserAvatar.svelte'
    import { _ } from '../../i18n.js'
    import { warrior } from '../../stores.js'

    export let user = {}
    export let votes = []
    export let maxVotes = 3
    export let facilitators = []
    export let handleAddFacilitator = () => {}
    export let handleRemoveFacilitator = () => {}

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
        {#if facilitators.includes(user.id)}
            <div class="text-indigo-500 dark:text-violet-400">
                {$_('facilitator')}
                {#if facilitators.includes($warrior.id)}
                    <button
                        class="text-red-500 text-sm"
                        on:click="{handleRemoveFacilitator(user.id)}"
                        >{$_('remove')}</button
                    >
                {/if}
            </div>
        {:else if facilitators.includes($warrior.id)}
            <div>
                <button
                    class="text-blue-500 dark:text-sky-400 text-sm"
                    on:click="{handleAddFacilitator(user.id)}"
                    >{$_('makeFacilitator')}</button
                >
            </div>
        {/if}
        {#if reachedMaxVotes}
            <div class="text-green-600 dark:text-green-400">
                {$_('allVotesIn')}
            </div>
        {/if}
    </div>
</div>
