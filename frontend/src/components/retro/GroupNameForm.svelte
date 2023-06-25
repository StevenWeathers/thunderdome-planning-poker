<script lang="ts">
    import LL from '../../i18n/i18n-svelte'

    export let groupId = ''
    export let groupName = ''
    export let handleGroupNameChange = () => {}

    function handleChange(e) {
        e.preventDefault()

        handleGroupNameChange(groupId, groupName)
    }
</script>

<form on:submit="{handleChange}">
    <input
        bind:value="{groupName}"
        on:change="{handleChange}"
        placeholder="{$LL.groupName()}"
        class="dark:bg-gray-800 border-gray-300 dark:border-gray-700 border-2 appearance-none rounded py-2
                    px-3 text-gray-700 dark:text-gray-200 leading-tight focus:outline-none text-center
                    focus:bg-white dark:focus:bg-gray-700 focus:border-indigo-500 dark:focus:border-yellow-400 w-full"
        id="groupName_{groupId}"
        name="groupName"
        type="text"
    />
    <button type="submit" class="hidden"></button>
</form>
