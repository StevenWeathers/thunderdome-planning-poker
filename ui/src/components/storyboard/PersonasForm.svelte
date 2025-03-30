<script lang="ts">
  import SolidButton from '../global/SolidButton.svelte';
  import Modal from '../global/Modal.svelte';
  import TextInput from '../forms/TextInput.svelte';


  interface Props {
    toggleEditPersona?: any;
    handlePersonaAdd?: any;
    handlePersonaRevision?: any;
    persona?: any;
  }

  let {
    toggleEditPersona = () => () => {},
    handlePersonaAdd = () => {},
    handlePersonaRevision = () => {},
    persona = $bindable({
    id: '',
    name: '',
    role: '',
    description: '',
  })
  }: Props = $props();

  function handleSubmit(event) {
    event.preventDefault();

    if (persona.id === '') {
      handlePersonaAdd({
        name: persona.name,
        role: persona.role,
        description: persona.description,
      });
    } else {
      handlePersonaRevision(persona);
    }
    toggleEditPersona();
  }
</script>

<Modal closeModal="{toggleEditPersona}">
  <form onsubmit={handleSubmit} name="addPersona">
    <div class="mb-4">
      <label
        class="block text-sm text-gray-700 dark:text-gray-400 font-bold mb-2"
        for="personaName"
      >
        Persona Name
      </label>
      <TextInput
        id="personaName"
        bind:value="{persona.name}"
        placeholder="Enter a persona name e.g. Ricky Bobby"
        name="personaName"
      />
    </div>
    <div class="mb-4">
      <label
        class="block text-sm text-gray-700 dark:text-gray-400 font-bold mb-2"
        for="personaRole"
      >
        Persona Role
      </label>
      <TextInput
        id="personaRole"
        bind:value="{persona.role}"
        placeholder="Enter a persona role e.g. Author, Developer, Admin"
        name="personaRole"
      />
    </div>
    <div class="mb-4">
      <label
        class="block text-sm text-gray-700 dark:text-gray-400 font-bold mb-2"
        for="personaDescription"
      >
        Persona Description
      </label>
      <textarea
        class="bg-gray-100 dark:bg-gray-900 dark:focus:bg-gray-800 border-gray-200 dark:border-gray-600 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 dark:text-gray-400 leading-tight
                focus:outline-none focus:bg-white focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
        id="personaDescription"
        bind:value="{persona.description}"
        placeholder="Enter a persona description"
        name="personaDescription"></textarea>
    </div>
    <div class="text-right">
      <div>
        <SolidButton type="submit">Save</SolidButton>
      </div>
    </div>
  </form>
</Modal>
