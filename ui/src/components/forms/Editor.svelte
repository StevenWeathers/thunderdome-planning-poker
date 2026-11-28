<script lang="ts">
  import { quill } from '../../quill';
  import '../../../node_modules/quill/dist/quill.core.css';
  import '../../../node_modules/quill/dist/quill.snow.css';

  interface Props {
    content?: string;
    placeholder?: string;
    handleTextChange?: any;
    id?: string;
  }

  let { content = '', placeholder = '', handleTextChange = () => {}, id = '' }: Props = $props();

  function onTextChange(e) {
    handleTextChange(e.detail.html);
  }
</script>

<div
  class="w-full"
  use:quill={{
    placeholder,
    content,
  }}
  on:textchange={onTextChange}
  {id}
></div>
