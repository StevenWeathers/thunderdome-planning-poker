<script lang="ts">
  import { Plus, Edit, Trash2, Users, User } from 'lucide-svelte';
  import SolidButton from '../global/SolidButton.svelte';
  import TextInput from '../forms/TextInput.svelte';
  import Modal from '../global/Modal.svelte';
  import type { StoryboardPersona } from '../../types/storyboard';
  
  interface Props {
    personas?: StoryboardPersona[];
    onAdd?: (persona: Omit<StoryboardPersona, 'id'>) => void;
    onUpdate?: (persona: StoryboardPersona) => void;
    onDelete?: (personaId: string) => void;
    closeModal?: () => void;
    isFacilitator?: boolean;
  }
  
  let { personas = [], onAdd, onUpdate, onDelete, closeModal, isFacilitator = false }: Props = $props();
  
  // Local state
  let editingId = $state<string | null>(null);
  let editForm = $state({
    name: '',
    role: '',
    description: ''
  });
  let showAddForm = $state(false);
  let addForm = $state({
    name: '',
    role: '',
    description: ''
  });



  // Start editing a persona
  function startEdit(persona: StoryboardPersona) {
    editingId = persona.id;
    editForm = {
      name: persona.name || '',
      role: persona.role || '',
      description: persona.description || ''
    };
  }

  // Cancel editing
  function cancelEdit() {
    editingId = null;
    editForm = { name: '', role: '', description: '' };
  }

  // Save edited persona
  function saveEdit(persona: StoryboardPersona) {
    if (!editForm.name.trim()) return;
    
    const updatedPersona = {
      ...persona,
      name: editForm.name.trim(),
      role: editForm.role.trim(),
      description: editForm.description.trim()
    };
    
    onUpdate?.(updatedPersona);
    cancelEdit();
  }

  // Delete persona
  function deletePersona(persona: StoryboardPersona) {
    onDelete?.(persona.id);
  }

  // Show add form
  function showAdd() {
    showAddForm = true;
    addForm = { name: '', role: '', description: '' };
  }

  // Cancel add
  function cancelAdd() {
    showAddForm = false;
    addForm = { name: '', role: '', description: '' };
  }

  // Add new persona
  function addPersona() {
    if (!addForm.name.trim()) return;
    
    const newPersona = {
      name: addForm.name.trim(),
      role: addForm.role.trim(),
      description: addForm.description.trim()
    };
    
    onAdd?.(newPersona);
    cancelAdd();
  }

  // Handle form submissions
  function handleEditSubmit(event: Event, persona: StoryboardPersona) {
    event.preventDefault();
    saveEdit(persona);
  }

  function handleAddSubmit(event: Event) {
    event.preventDefault();
    addPersona();
  }
</script>

<Modal closeModal={closeModal} ariaLabel={$LL.modalStoryboardPersonas()}>
  <div class="persona-list pt-6">
    <!-- Header -->
    <div class="flex items-center justify-between mb-4">
      <h2 class="text-xl font-bold text-gray-900 dark:text-gray-100">Personas</h2>
      {#if isFacilitator}
        <SolidButton type="button" onClick={showAdd}>
          <Plus class="w-4 h-4 me-1" />
          Add Persona
        </SolidButton>
      {/if}
    </div>

  <!-- Add Form -->
  {#if showAddForm && isFacilitator}
    <div class="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-4 mb-4">
      <form onsubmit={handleAddSubmit} class="space-y-3">
        <div>
          <label for="add-name" class="block text-gray-700 dark:text-gray-300 mb-1">Name *</label>
          <TextInput
            id="add-name"
            bind:value={addForm.name}
            placeholder="Enter persona name"
            autofocus
          />
        </div>
        <div>
          <label for="add-role" class="block text-gray-700 dark:text-gray-300 mb-1">Role</label>
          <TextInput
            id="add-role"
            bind:value={addForm.role}
            placeholder="Enter role (optional)"
          />
        </div>
        <div>
          <label for="add-description" class="block text-gray-700 dark:text-gray-300 mb-1">Description</label>
          <textarea
            id="add-description"
            bind:value={addForm.description}
            placeholder="Enter description (optional)"
            rows="2"
            class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-gray-100 resize-none"
          ></textarea>
        </div>
        <div class="flex justify-end space-x-2 mt-4">
          <button
            type="button"
            onclick={cancelAdd}
            class="px-3 py-1.5 text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md hover:bg-gray-50 dark:hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition-colors"
          >
            Cancel
          </button>
          <SolidButton type="submit" disabled={!addForm.name.trim()}>
            Add Persona
          </SolidButton>
        </div>
      </form>
    </div>
  {/if}

  <!-- Personas List -->
  <div class="space-y-3">
    {#each personas as persona (persona.id)}
      <div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-2 lg:p-4 hover:border-gray-300 dark:hover:border-gray-600 transition-colors">
        {#if editingId === persona.id && isFacilitator}
          <!-- Edit Mode -->
          <form onsubmit={(e) => handleEditSubmit(e, persona)} class="space-y-3">
            <div>
              <label for="edit-name-{persona.id}" class="block text-gray-700 dark:text-gray-300 mb-1">Name *</label>
              <TextInput
                id="edit-name-{persona.id}"
                bind:value={editForm.name}
                autofocus
              />
            </div>
            <div>
              <label for="edit-role-{persona.id}" class="block text-gray-700 dark:text-gray-300 mb-1">Role</label>
              <TextInput
                id="edit-role-{persona.id}"
                bind:value={editForm.role}
              />
            </div>
            <div>
              <label for="edit-description-{persona.id}" class="block text-gray-700 dark:text-gray-300 mb-1">Description</label>
              <textarea
                id="edit-description-{persona.id}"
                bind:value={editForm.description}
                rows="2"
                class="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-gray-100 resize-none"
              ></textarea>
            </div>
            <div class="flex justify-end space-x-2">
              <button
                type="button"
                onclick={cancelEdit}
                class="px-3 py-1.5 text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md hover:bg-gray-50 dark:hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition-colors"
              >
                Cancel
              </button>
              <SolidButton type="submit" disabled={!editForm.name.trim()}>
                Save
              </SolidButton>
            </div>
          </form>
        {:else}
          <!-- Display Mode -->
          <div class="flex justify-between items-center">
            <div class="flex-1">
              <div class="flex items-center space-x-2 mb-1">
                <User class="w-6 h-6 text-gray-500 dark:text-gray-400" />
                <h4 class="font-semibold text-gray-900 dark:text-gray-100">{persona.name}</h4>
                {#if persona.role}
                  <span class="inline-flex items-center px-2 py-0.5 rounded text-xs bg-blue-100 dark:bg-blue-900/50 text-blue-800 dark:text-blue-200">
                    {persona.role}
                  </span>
                {/if}
              </div>
              {#if persona.description}
                <p class="text-gray-600 dark:text-gray-400 mt-1">{persona.description}</p>
              {/if}
            </div>
            {#if isFacilitator}
              <div class="flex space-x-1 ms-4">
                <button
                  onclick={() => startEdit(persona)}
                  class="p-1.5 text-gray-400 dark:text-gray-500 hover:text-blue-600 dark:hover:text-blue-400 hover:bg-blue-50 dark:hover:bg-blue-900/20 rounded-md transition-colors"
                  title="Edit persona"
                >
                  <Edit class="w-5 h-5" />
                </button>
                <button
                  onclick={() => deletePersona(persona)}
                  class="p-1.5 text-gray-400 dark:text-gray-500 hover:text-red-600 dark:hover:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 rounded-md transition-colors"
                  title="Delete persona"
                >
                  <Trash2 class="w-5 h-5" />
                </button>
              </div>
            {/if}
          </div>
        {/if}
      </div>
    {:else}
      <div class="py-8 text-gray-600 dark:text-gray-300 font-bold">
        <div class="text-center">
            <Users class="w-12 h-12 mx-auto mb-4" />
            <p class="text-xl text-gray-900 dark:text-white">No personas yet</p>
            <p class="text-lg mb-4">Add your first persona to get started</p>
        </div>
        <p class="rounded-lg text-left bg-gray-200 dark:bg-gray-700 p-3">User personas in agile story mapping represent different user types whose unique goals and needs help organize story columns by ensuring each vertical slice addresses specific user journeys and priorities.</p>
      </div>
    {/each}
      </div>
  </div>
</Modal>

<style>
  .persona-list {
    @apply overflow-y-auto;
  }
</style>