<script lang="ts">
  import { Tag, X } from '@lucide/svelte';
  import { getTagSuggestions, RETRO_ITEM_MAX_TAGS, RETRO_TAG_MAX_LENGTH } from './retroTags';

  interface Props {
    itemId: string;
    tags?: string[];
    disabled?: boolean;
    sendSocketEvent?: (event: string, value: any) => void;
  }

  let { itemId, tags = [], disabled = false, sendSocketEvent = () => {} }: Props = $props();

  const suggestions = getTagSuggestions();

  let showPicker = $state(false);
  let search = $state('');

  const canAdd = $derived(!disabled && tags.length < RETRO_ITEM_MAX_TAGS);
  const applied = $derived(new Set(tags.map(tag => tag.toLowerCase())));
  const normalizedSearch = $derived(search.trim().replace(/\s+/g, ' '));
  const searchKey = $derived(normalizedSearch.toLowerCase());
  const options = $derived(
    suggestions().filter(tag => !applied.has(tag.toLowerCase()) && tag.toLowerCase().includes(searchKey)),
  );
  const isCustom = $derived(
    searchKey !== '' && !applied.has(searchKey) && !options.some(tag => tag.toLowerCase() === searchKey),
  );

  function closePicker() {
    showPicker = false;
    search = '';
  }

  function togglePicker() {
    if (showPicker) {
      closePicker();
    } else {
      showPicker = true;
    }
  }

  function addTag(tag: string) {
    sendSocketEvent('item_tag_add', JSON.stringify({ item_id: itemId, tag }));
    closePicker();
  }

  function removeTag(tag: string) {
    sendSocketEvent('item_tag_remove', JSON.stringify({ item_id: itemId, tag }));
  }

  function handleSubmit(event: Event) {
    event.preventDefault();

    if (searchKey === '' || applied.has(searchKey)) {
      return;
    }

    // reuse the casing of an existing tag when the typed value matches one
    addTag(options.find(tag => tag.toLowerCase() === searchKey) ?? normalizedSearch);
  }

  function handleKeydown(event: KeyboardEvent) {
    if (event.key === 'Escape') {
      event.stopPropagation();
      closePicker();
    }
  }

  function clickOutside(node: HTMLElement) {
    function handleClick(event: MouseEvent) {
      if (showPicker && !node.contains(event.target as Node)) {
        closePicker();
      }
    }

    document.addEventListener('click', handleClick, true);

    return {
      destroy() {
        document.removeEventListener('click', handleClick, true);
      },
    };
  }

  function focusOnMount(node: HTMLInputElement) {
    node.focus();
  }
</script>

<div class="flex flex-wrap items-center gap-2" data-testid="retro-feedback-item-tags">
  {#each tags as tag (tag)}
    <span
      class="inline-flex items-center gap-1 max-w-full rounded-md border border-indigo-200 bg-indigo-50 px-2 py-0.5 text-xs font-medium text-indigo-800 dark:border-indigo-500/50 dark:bg-indigo-500/10 dark:text-indigo-200"
      data-testid="retro-item-tag"
    >
      <span class="min-w-0 break-words">{tag}</span>
      {#if !disabled}
        <button
          type="button"
          class="leading-none hover:text-red-500 dark:hover:text-red-400"
          aria-label="Remove tag {tag}"
          onclick={() => removeTag(tag)}
        >
          <X class="h-3 w-3" aria-hidden="true" />
        </button>
      {/if}
    </span>
  {/each}

  {#if canAdd}
    <div class="relative" use:clickOutside>
      <button
        type="button"
        class="inline-flex items-center gap-1 rounded-md border border-dashed border-gray-300 px-2 py-0.5 text-xs text-gray-600 transition-colors duration-200 hover:border-blue-300 hover:text-blue-600 dark:border-gray-600 dark:text-gray-300 dark:hover:border-sky-500 dark:hover:text-sky-300"
        aria-label="Add tag"
        aria-expanded={showPicker}
        data-testid="retro-item-tag-toggle"
        onclick={togglePicker}
      >
        <Tag class="h-3 w-3" aria-hidden="true" />
        Tag
      </button>

      {#if showPicker}
        <form
          class="absolute left-0 top-full z-10 mt-2 w-56 rounded-lg border border-gray-200 bg-white p-2 shadow-lg dark:border-gray-600 dark:bg-gray-800"
          data-testid="retro-item-tag-picker"
          onsubmit={handleSubmit}
        >
          <input
            type="text"
            class="w-full rounded-md bg-gray-100 px-2 py-1 text-sm text-gray-900 outline-none focus:ring-2 focus:ring-purple-500 dark:bg-gray-700 dark:text-white dark:focus:ring-purple-400"
            aria-label="Search or create a tag"
            placeholder="Search or create a tag"
            maxlength={RETRO_TAG_MAX_LENGTH}
            bind:value={search}
            onkeydown={handleKeydown}
            use:focusOnMount
          />
          <div class="mt-1 max-h-40 overflow-y-auto">
            {#each options as option (option)}
              <button
                type="button"
                class="block w-full break-words rounded px-2 py-1 text-left text-sm text-gray-800 hover:bg-gray-100 dark:text-gray-200 dark:hover:bg-gray-700"
                data-testid="retro-item-tag-option"
                onclick={() => addTag(option)}
              >
                {option}
              </button>
            {/each}
            {#if isCustom}
              <button
                type="submit"
                class="block w-full break-words rounded px-2 py-1 text-left text-sm text-blue-600 hover:bg-gray-100 dark:text-sky-300 dark:hover:bg-gray-700"
                data-testid="retro-item-tag-create"
              >
                Add "{normalizedSearch}"
              </button>
            {/if}
          </div>
        </form>
      {/if}
    </div>
  {/if}
</div>
