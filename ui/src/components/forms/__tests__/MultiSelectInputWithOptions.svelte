<script lang="ts">
  import MultiSelectInput from '../MultiSelectInput.svelte';

  interface Props {
    value?: string[];
    name?: string;
    id?: string;
    required?: boolean;
    disabled?: boolean;
    onchange?: (event: Event) => void;
    size?: number;
  }

  let {
    value = $bindable<string[]>([]),
    name = '',
    id = '',
    required = false,
    disabled = false,
    onchange,
    size = 4,
  }: Props = $props();

  let selectRef: MultiSelectInput;

  export function focus() {
    selectRef?.focus();
  }
</script>

<MultiSelectInput {value} {name} {id} {required} {disabled} {onchange} {size} bind:this={selectRef}>
  <option value="one">One</option>
  <option value="two">Two</option>
  <option value="three">Three</option>
</MultiSelectInput>
