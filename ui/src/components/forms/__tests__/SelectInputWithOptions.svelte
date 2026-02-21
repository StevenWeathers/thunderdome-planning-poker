<script lang="ts">
  import SelectInput from '../SelectInput.svelte';

  export let value: string = '';
  export let name: string = '';
  export let id: string = '';
  export let required: boolean = false;
  export let disabled: boolean = false;

  let selectRef: SelectInput;

  export function focus() {
    selectRef?.focus();
  }
</script>

<SelectInput bind:value {name} {id} {required} {disabled} bind:this={selectRef}>
  <option value="" disabled>Select one</option>
  <option value="one">One</option>
  <option value="two">Two</option>
</SelectInput>
