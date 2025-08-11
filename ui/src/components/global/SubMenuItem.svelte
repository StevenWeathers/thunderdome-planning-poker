<script lang="ts">
    interface Props {
        onClickHandler: any;
        testId?: String;
        icon?: any;
        label: string;
        active?: boolean;
    }

    let { onClickHandler = () => {}, testId = '', label = '', icon, active = false }: Props = $props();
</script>

<li class="flex">
    <button
        class="inline-flex items-center w-full px-2 py-1 font-semibold transition-colors duration-150 rounded-md hover:bg-gray-100 hover:text-gray-800 dark:hover:bg-gray-800 dark:hover:text-gray-200"
        class:bg-indigo-600={active}
        class:text-white={active}
        onclick={onClickHandler}
        data-testid={testId}
    >
        {#if icon}{@const SvelteComponent = icon}
      <SvelteComponent class="h-5 w-5 me-2 inline-block" />{/if}{label}
    </button>
</li>