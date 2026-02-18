<script lang="ts">
  import BrowserMock from '../BrowserMock.svelte';

  interface Props {
    testContent?: string;
    customClass?: string;
  }

  let { testContent = 'Test content', customClass }: Props = $props();
</script>

<BrowserMock class={customClass}>
  <div class="test-child-content" data-testid="browser-child-content">
    {testContent}
  </div>
</BrowserMock>
