<script lang="ts">
  interface Props {
    value?: string;
    placeholder?: string;
    id?: string;
    name?: string;
    required?: boolean;
    class?: string;
    onkeydown?: (e: KeyboardEvent) => void;
  }

  let {
    value = $bindable(''),
    placeholder = '',
    id = '',
    name = '',
    required = false,
    class: className = '',
    onkeydown,
  }: Props = $props();

  let textareaElement: HTMLTextAreaElement;
  let initialHeight = 0;

  $effect(() => {
    // Watch value and recalculate height when it changes
    value;
    recalculateHeight();
  });

  const handleInput = (e: Event) => {
    const target = e.currentTarget as HTMLTextAreaElement;
    if (initialHeight === 0) {
      initialHeight = target.offsetHeight;
    }
    target.style.height = 'auto';
    const newHeight = Math.max(initialHeight, target.scrollHeight);
    target.style.height = newHeight + 'px';
  };

  const resetHeight = () => {
    if (textareaElement) {
      textareaElement.style.height = '';
      initialHeight = 0;
    }
  };

  const recalculateHeight = () => {
    if (textareaElement) {
      if (initialHeight === 0) {
        initialHeight = textareaElement.offsetHeight;
      }
      textareaElement.style.height = 'auto';
      const newHeight = Math.max(initialHeight, textareaElement.scrollHeight);
      textareaElement.style.height = newHeight + 'px';
    }
  };

  const focus = () => {
    textareaElement?.focus();
  };

  export { resetHeight, focus };
</script>

<textarea
  bind:this={textareaElement}
  bind:value
  {placeholder}
  {id}
  {name}
  {required}
  rows="1"
  class="dark:bg-gray-800 border-gray-300 dark:border-gray-700 border-2 appearance-none rounded py-2
            px-3 text-gray-700 dark:text-gray-400 leading-tight focus:outline-none
            focus:bg-white dark:focus:bg-gray-700 focus:border-indigo-500 dark:focus:border-yellow-400 w-full resize-none overflow-hidden box-border {className}"
  oninput={handleInput}
  {onkeydown}
></textarea>
