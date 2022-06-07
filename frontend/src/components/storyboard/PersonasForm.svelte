<script>
  import SolidButton from '../SolidButton.svelte'
  import Modal from '../Modal.svelte'

  export let toggleEditPersona = () => () => {}
  export let handlePersonaAdd = () => {}
  export let handlePersonaRevision = () => {}

  export let persona = {
    id: '',
    name: '',
    role: '',
    description: '',
  }

  function handleSubmit (event) {
    event.preventDefault()

    if (persona.id === '') {
      handlePersonaAdd({
        name: persona.name,
        role: persona.role,
        description: persona.description,
      })
    } else {
      handlePersonaRevision(persona)
    }
    toggleEditPersona()
  }
</script>

<Modal closeModal="{toggleEditPersona}">
    <form on:submit="{handleSubmit}" name="addPersona">
        <div class="mb-4">
            <label class="block text-sm font-bold mb-2" for="personaName">
                Persona Name
            </label>
            <input
                    class="bg-gray-200 border-gray-200 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 leading-tight
                focus:outline-none focus:bg-white focus:border-purple-500"
                    id="personaName"
                    type="text"
                    bind:value="{persona.name}"
                    placeholder="Enter a persona name e.g. Ricky Bobby"
                    name="personaName"
            />
        </div>
        <div class="mb-4">
            <label class="block text-sm font-bold mb-2" for="personaRole">
                Persona Role
            </label>
            <input
                    class="bg-gray-200 border-gray-200 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 leading-tight
                focus:outline-none focus:bg-white focus:border-purple-500"
                    id="personaRole"
                    type="text"
                    bind:value="{persona.role}"
                    placeholder="Enter a persona role e.g. Author, Developer, Admin"
                    name="personaRole"
            />
        </div>
        <div class="mb-4">
            <label
                    class="block text-sm font-bold mb-2"
                    for="personaDescription"
            >
                Persona Description
            </label>
            <textarea
                    class="bg-gray-200 border-gray-200 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 leading-tight
                focus:outline-none focus:bg-white focus:border-purple-500"
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
