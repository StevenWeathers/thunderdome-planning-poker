<script lang="ts">
  import SolidButton from '../global/SolidButton.svelte';
  import Modal from '../global/Modal.svelte';
  import HollowButton from '../global/HollowButton.svelte';
  import TextInput from '../forms/TextInput.svelte';
  import SelectInput from '../forms/SelectInput.svelte';
  import LL from '../../i18n/i18n-svelte';
  import { Trash2, User } from 'lucide-svelte';
  import { onMount } from 'svelte';

  interface Props {
    toggleColumnEdit?: any;
    handleColumnRevision?: any;
    deleteColumn?: any;
    handlePersonaRemove?: any;
    handlePersonaAdd?: any;
    personas?: any;
    column?: any;
  }

  let {
    toggleColumnEdit = () => {},
    handleColumnRevision = () => {},
    deleteColumn = () => () => {},
    handlePersonaRemove = () => () => {},
    handlePersonaAdd = () => {},
    personas = [],
    column = $bindable({
    id: '',
    name: '',
    personas: [],
  })
  }: Props = $props();

  let selectedPersona = $state('');
  let focusInput: any;

  function handleSubmit(event) {
    event.preventDefault();

    const c = {
      id: column.id,
      name: column.name,
    };

    handleColumnRevision(c);
    toggleColumnEdit();
  }

  function addPersona() {
    handlePersonaAdd({ column_id: column.id, persona_id: selectedPersona });
  }

  onMount(() => {
    focusInput?.focus();
  });
</script>

<Modal closeModal={toggleColumnEdit}>
  <form onsubmit={handleSubmit} name="addColumn">
    <div class="mb-4">
      <label
        class="block text-sm text-gray-700 dark:text-gray-400 font-bold mb-2"
        for="columnName"
      >
        Column Name
      </label>
      <TextInput
        id="columnName"
        bind:value="{column.name}"
        placeholder="Enter a column name"
        name="columnName"
        bind:this={focusInput}
      />
    </div>
    <div class="flex">
      <div class="md:w-1/2 text-left">
        <HollowButton color="red" onClick={deleteColumn(column.id)}>
          Delete Column
        </HollowButton>
      </div>
      <div class="md:w-1/2 text-right">
        <SolidButton type="submit">Save</SolidButton>
      </div>
    </div>
  </form>

  <div class="mt-4 pt-2 border-t border-gray-400 dark:border-gray-700">
    <div class="block text-gray-700 dark:text-gray-400 font-bold mb-4">
      {$LL.personas()}
    </div>
    <div class="flex w-full gap-4">
      <div class="w-2/3">
        <SelectInput bind:value="{selectedPersona}" id="persona" name="persona">
          <option value="" disabled>Select a persona</option>
          {#each personas as persona}
            <option value="{persona.id}">
              {persona.name} ({persona.role})
            </option>
          {/each}
        </SelectInput>
      </div>
      <div class="w-1/3">
        <HollowButton
          onClick={addPersona}
          disabled={selectedPersona === ''}
        >
          Add Persona
        </HollowButton>
      </div>
    </div>
    {#if column.personas.length}
      <div class="grid grid-cols-2 gap-4 mt-4">
        {#each column.personas as persona}
          <div class="flex text-gray-700 dark:text-gray-400 mb-2">
            <div class="w-1/4">
              <User class="inline-block" />
            </div>
            <div class="w-2/4 text-lg">{persona.name} ({persona.role})</div>
            <div class="w-1/4 text-right">
              <HollowButton
                color="red"
                onClick={handlePersonaRemove({
                  column_id: column.id,
                  persona_id: persona.id,
                })}
              >
                <Trash2 />
              </HollowButton>
            </div>
          </div>
        {/each}
      </div>
    {/if}
  </div>
</Modal>
