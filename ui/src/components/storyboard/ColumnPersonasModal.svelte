<script lang="ts">
  import Modal from '../global/Modal.svelte';
  import HollowButton from '../global/HollowButton.svelte';
  import LL from '../../i18n/i18n-svelte';
  import { TrashIcon, Plus, User } from '@lucide/svelte';

  interface Props {
    toggleModal?: () => void;
    column: any;
    personas: any[];
    onPersonaAdd: (data: { column_id: string; persona_id: string }) => void;
    onPersonaRemove: (data: { column_id: string; persona_id: string }) => void;
  }

  let {
    toggleModal = () => {},
    column = { id: '', name: '', personas: [] },
    personas = [],
    onPersonaAdd = () => {},
    onPersonaRemove = () => () => {},
  }: Props = $props();

  function addPersona(personaId: string) {
    onPersonaAdd({ column_id: column.id, persona_id: personaId });
  }

  function handlePersonaRemove(personaId: string) {
    onPersonaRemove({ column_id: column.id, persona_id: personaId });
  }

  const addedPersonaIds = $derived(new Set(column.personas.map((p: any) => p.id)));
  const availablePersonas = $derived(personas.filter((p: any) => !addedPersonaIds.has(p.id)));
</script>

<Modal closeModal={toggleModal} ariaLabel={$LL.modalStoryboardColumnSettings()}>
  <div class="mainContainer">
    <div class="sectionTitle">{$LL.columnPersonasTitle()}</div>

    <div class="flexColGap3">
      {#if addedPersonaIds.size > 0}
        <div class="spaceY2">
          {#each column.personas as persona}
            <div class="personaCard">
              <div class="personaIconContainer">
                <User size={18} class="personaIcon" />
                <div class="personaContent">
                  <div class="personaName">{persona.name}</div>
                  <div class="personaRole">{persona.role}</div>
                </div>
              </div>
              <HollowButton color="red" onClick={() => handlePersonaRemove(persona.id)}>
                <TrashIcon size={18} />
              </HollowButton>
            </div>
          {/each}
        </div>
      {:else}
        <div class="emptyState">{$LL.columnPersonasEmpty()}</div>
      {/if}
    </div>

    <div class="flexColGap3">
      <div class="sectionHeader">{$LL.availablePersonas()}</div>
      {#if availablePersonas.length}
        <div class="spaceY2">
          {#each availablePersonas as persona}
            <div class="personaCardHoverable">
              <div class="personaIconContainer">
                <User size={18} class="personaIcon" />
                <div class="personaContent">
                  <div class="personaName">{persona.name}</div>
                  <div class="personaRole">{persona.role}</div>
                </div>
              </div>
              <HollowButton onClick={() => addPersona(persona.id)}>
                <Plus size={18} />
              </HollowButton>
            </div>
          {/each}
        </div>
      {:else if personas.length === 0}
        <div class="emptyState">{$LL.availablePersonasEmpty()}</div>
      {:else}
        <div class="emptyState">{$LL.allPersonasAdded()}</div>
      {/if}
    </div>
  </div>
</Modal>

<style lang="postcss">
  .mainContainer {
    @apply mt-4 pt-2 flex flex-col gap-6;
  }

  .sectionTitle {
    @apply block text-gray-700 font-bold text-lg;
  }

  :root.dark .sectionTitle {
    @apply text-gray-400;
  }

  .flexColGap3 {
    @apply flex flex-col gap-3;
  }

  .spaceY2 {
    @apply space-y-2;
  }

  .personaCard {
    @apply flex items-center justify-between bg-gray-50 p-3 rounded-lg;
  }

  :root.dark .personaCard {
    @apply bg-gray-800;
  }

  .personaCardHoverable {
    @apply flex items-center justify-between bg-gray-50 p-3 rounded-lg transition-colors;
  }

  .personaCardHoverable:hover {
    @apply bg-gray-100;
  }

  :root.dark .personaCardHoverable {
    @apply bg-gray-800;
  }

  :root.dark .personaCardHoverable:hover {
    @apply bg-gray-700;
  }

  .personaIconContainer {
    @apply text-gray-500 flex items-center gap-3;
  }

  :root.dark .personaIconContainer {
    @apply text-gray-400;
  }

  .personaContent {
    @apply flex-1;
  }

  .personaName {
    @apply text-gray-700 font-medium;
  }

  :root.dark .personaName {
    @apply text-gray-300;
  }

  .personaRole {
    @apply text-sm text-gray-500;
  }

  :root.dark .personaRole {
    @apply text-gray-400;
  }

  .sectionHeader {
    @apply font-semibold text-gray-600 bg-gray-50;
  }

  :root.dark .sectionHeader {
    @apply text-gray-400 bg-gray-800;
  }

  .emptyState {
    @apply text-center py-4 text-gray-500 bg-gray-200 rounded-lg;
  }

  :root.dark .emptyState {
    @apply text-gray-400 bg-gray-700;
  }
</style>
