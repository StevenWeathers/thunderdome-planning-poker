<script lang="ts">
  import { quill } from '../../quill';

  export let content = '';
  export let placeholder = '';
  export let handleTextChange = text => {};
  export let id = '';

  function onTextChange(e) {
    handleTextChange(e.detail.html);
  }
</script>

<div
  class="w-full"
  use:quill="{{
    placeholder: `${placeholder}`,
    content: content,
  }}"
  on:text-change="{onTextChange}"
  id="{id}"
></div>
