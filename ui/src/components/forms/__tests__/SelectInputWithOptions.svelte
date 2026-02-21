<script lang="ts">
  import type { Snippet } from 'svelte';
  import SelectInput from '../SelectInput.svelte';

  interface Props {
    value?: any;
    name?: string;
    id?: string;
    required?: boolean;
    disabled?: boolean;
    onchange?: (event: Event) => void;
  }

  let { value = $bindable(), name = '', id = '', required = false, disabled = false, onchange }: Props = $props();

  let selectRef: SelectInput;

  export function focus() {
    selectRef?.focus();
  }
</script>

<SelectInput {value} {name} {id} {required} {disabled} {onchange} bind:this={selectRef}>
  <option value="" disabled>Select one</option>
  <option value="one">One</option>
  <option value="two">Two</option>
</SelectInput>
