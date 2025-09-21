<script lang="ts">
  import { quill } from '../../quill';

  interface Props {
    content?: string;
    placeholder?: string;
    handleTextChange?: any;
    id?: string;
  }

  let { content = '', placeholder = '', handleTextChange = text => {}, id = '' }: Props = $props();

  function onTextChange(e) {
    handleTextChange(e.detail.html);
  }
</script>

<div
  class="w-full"
  use:quill={{
    placeholder: `${placeholder}`,
    content: content,
  }}
  ontext-change={onTextChange}
  {id}
></div>
