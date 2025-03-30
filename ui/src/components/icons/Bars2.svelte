<script lang="ts">
  interface Props {
    class?: string;
  }

  let { class: klass = 'h-6 w-6' }: Props = $props();
  
</script>

<svg
  xmlns="http://www.w3.org/2000/svg"
  fill="none"
  viewBox="0 0 24 24"
  stroke-width="1.5"
  stroke="currentColor"
  class="{klass}"
>
  <path
    stroke-linecap="round"
    stroke-linejoin="round"
    d="M3.75 9h16.5m-16.5 6.75h16.5"></path>
</svg>
