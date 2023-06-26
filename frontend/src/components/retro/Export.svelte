<script lang="ts">
    import LL from '../../i18n/i18n-svelte'

    export let retro = {
        items: [],
        actionItems: [],
    }

    $: workedItems =
        retro.items &&
        retro.items.reduce((prev, item) => {
            if (item.type === 'worked') {
                prev.push(item)
            }
            return prev
        }, [])
    $: improveItems =
        retro.items &&
        retro.items.reduce((prev, item) => {
            if (item.type === 'improve') {
                prev.push(item)
            }
            return prev
        }, [])
    $: questionItems =
        retro.items &&
        retro.items.reduce((prev, item) => {
            if (item.type === 'question') {
                prev.push(item)
            }
            return prev
        }, [])
</script>

<div class="flex flex-grow p-4 dark:text-white">
    <div class="px-4">
        <div class="mb-4">
            <h2 class="text-3xl font-rajdhani">{$LL.works()}</h2>
            <ul class="ps-12 text-lg">
                {#each workedItems as item (item.id)}
                    <li>
                        &bull; {item.content}
                    </li>
                {/each}
            </ul>
        </div>
        <div class="mb-4">
            <h2 class="text-3xl font-rajdhani">{$LL.needsImprovement()}</h2>
            <ul class="ps-12 text-lg">
                {#each improveItems as item (item.id)}
                    <li>
                        &bull; {item.content}
                    </li>
                {/each}
            </ul>
        </div>
        <div class="mb-4">
            <h2 class="text-3xl font-rajdhani">{$LL.questions()}</h2>
            <ul class="ps-12 text-lg">
                {#each questionItems as item (item.id)}
                    <li>
                        &bull; {item.content}
                    </li>
                {/each}
            </ul>
        </div>
        <div class="mb-4">
            <h2 class="text-3xl font-rajdhani">{$LL.actionItems()}</h2>
            <ul class="ps-12 text-lg">
                {#each retro.actionItems as item (item.id)}
                    <li>&bull; {item.content}</li>
                {/each}
            </ul>
        </div>
    </div>
</div>
