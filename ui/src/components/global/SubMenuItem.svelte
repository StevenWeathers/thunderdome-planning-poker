<script lang="ts">
    import LL from "../../i18n/i18n-svelte";

    interface Props {
        onClickHandler: any;
        testId?: String;
        children?: import('svelte').Snippet;
    }

    let { onClickHandler = () => {}, testId = '', children }: Props = $props();
  
</script>

<li class="flex">
    <button
        class="inline-flex items-center w-full px-2 py-1 font-semibold transition-colors duration-150 rounded-md hover:bg-gray-100 hover:text-gray-800 dark:hover:bg-gray-800 dark:hover:text-gray-200"
        onclick={onClickHandler}
        data-testid={testId}
    >
        {@render children?.()}
    </button>
</li>