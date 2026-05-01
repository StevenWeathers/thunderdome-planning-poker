<script lang="ts">
  import { Pencil, Plus, Trash2, X } from '@lucide/svelte';
  import ActionsMenu from '../global/ActionsMenu.svelte';
  import HollowButton from '../global/HollowButton.svelte';
  import Modal from '../global/Modal.svelte';
  import UserAvatar from '../user/UserAvatar.svelte';
  import LL from '../../i18n/i18n-svelte';
  import type { TeamUser } from '../../types/team';

  interface TeamKudo {
    id: string;
    user: TeamUser;
    targetUser: TeamUser;
    comment: string;
  }

  interface PendingKudo {
    targetUserId: string;
    comment: string;
  }

  interface Props {
    kudos?: TeamKudo[];
    toggleKudos?: () => void;
    currentUserId?: string;
    isAdmin?: boolean;
    users?: TeamUser[];
    onCreate?: (data: { targetUserId: string; comment: string }) => Promise<boolean>;
    onUpdate?: (kudoId: string, data: { comment: string }) => void;
    onDelete?: (kudoId: string) => void;
  }

  let {
    kudos = [],
    toggleKudos = () => {},
    currentUserId = '',
    isAdmin = false,
    users = [],
    onCreate = async () => false,
    onUpdate = () => {},
    onDelete = () => {},
  }: Props = $props();

  let editingKudoId = $state('');
  let draftComment = $state('');
  let showCreateForm = $state(false);
  let kudosToCreate = $state<PendingKudo[]>([{ targetUserId: '', comment: '' }]);

  const availableKudoUsers = $derived(users.filter((teamUser: TeamUser) => teamUser.id !== currentUserId));

  function canManageKudo(kudo: TeamKudo): boolean {
    return isAdmin || kudo.user.id === currentUserId;
  }

  function startEditing(kudo: TeamKudo): void {
    editingKudoId = kudo.id;
    draftComment = kudo.comment;
  }

  function cancelEditing(): void {
    editingKudoId = '';
    draftComment = '';
  }

  function saveEditing(kudoId: string): void {
    onUpdate(kudoId, {
      comment: draftComment,
    });
    cancelEditing();
  }

  function updatePendingKudo(index: number, field: keyof PendingKudo, value: string) {
    kudosToCreate = kudosToCreate.map((kudo: PendingKudo, currentIndex: number) =>
      currentIndex === index ? { ...kudo, [field]: value } : kudo,
    );
  }

  function addKudoRow() {
    kudosToCreate = [...kudosToCreate, { targetUserId: '', comment: '' }];
  }

  function removeKudoRow(index: number) {
    kudosToCreate = kudosToCreate.filter((_: PendingKudo, currentIndex: number) => currentIndex !== index);
    if (kudosToCreate.length === 0) {
      kudosToCreate = [{ targetUserId: '', comment: '' }];
    }
  }

  function resetCreateForm() {
    showCreateForm = false;
    kudosToCreate = [{ targetUserId: '', comment: '' }];
  }

  async function saveNewKudos(): Promise<void> {
    const completedKudos = kudosToCreate.filter(
      (kudo: PendingKudo) => kudo.targetUserId !== '' && kudo.comment.trim() !== '',
    );
    if (completedKudos.length === 0) {
      return;
    }

    let allSucceeded = true;
    for (const kudo of completedKudos) {
      const created = await onCreate({
        targetUserId: kudo.targetUserId,
        comment: kudo.comment,
      });
      if (!created) {
        allSucceeded = false;
      }
    }

    if (allSucceeded) {
      resetCreateForm();
    }
  }
</script>

<Modal closeModal={toggleKudos} widthClasses="md:w-4/5 lg:w-3/4 xl:w-2/3 max-w-5xl" ariaLabel="Team kudos">
  <div class="mt-6 space-y-5">
    <div class="flex flex-wrap items-end justify-between gap-3 border-b border-slate-200 pb-4 dark:border-slate-700">
      <div>
        <div class="flex flex-wrap items-center gap-2">
          <h2 class="text-3xl font-rajdhani font-semibold uppercase tracking-wide text-slate-900 dark:text-white">
            Kudos
          </h2>
          <span
            class="inline-flex items-center rounded-full bg-indigo-100 px-3 py-1 text-sm font-semibold text-indigo-700 ring-1 ring-indigo-200 dark:bg-indigo-500/20 dark:text-indigo-300 dark:ring-indigo-400/35"
          >
            {kudos.length}
          </span>
        </div>
        <p class="mt-1 text-sm text-slate-600 dark:text-slate-300">Appreciation shared in today&apos;s check-in.</p>
      </div>
      <div class="flex flex-wrap items-center gap-2">
        {#if availableKudoUsers.length > 0}
          <HollowButton
            color="green"
            size="medium"
            class="inline-flex items-center gap-2 whitespace-nowrap"
            onClick={() => (showCreateForm = !showCreateForm)}
          >
            <Plus class="h-4 w-4 shrink-0" />
            <span class="whitespace-nowrap">Give Kudos</span>
          </HollowButton>
        {/if}
      </div>
    </div>

    {#if showCreateForm}
      <div class="rounded-2xl border border-slate-200 bg-slate-50 p-4 dark:border-slate-700 dark:bg-slate-900/30">
        <div class="space-y-4">
          {#each kudosToCreate as pendingKudo, index}
            <div class="rounded-xl border border-slate-200 bg-white p-4 dark:border-slate-700 dark:bg-slate-800/80">
              <div class="grid gap-4 lg:grid-cols-[minmax(0,220px)_1fr_auto] lg:items-start">
                <div>
                  <label
                    for={`modal-kudo-target-${index}`}
                    class="mb-2 block text-xs font-medium uppercase tracking-wider text-slate-500 dark:text-slate-400"
                  >
                    Teammate
                  </label>
                  <select
                    id={`modal-kudo-target-${index}`}
                    class="w-full rounded-xl border border-slate-300 bg-white px-3 py-2 text-sm text-slate-900 shadow-sm outline-none transition focus:border-indigo-500 focus:ring-2 focus:ring-indigo-200 dark:border-slate-600 dark:bg-slate-800 dark:text-white dark:focus:ring-indigo-500/30"
                    bind:value={pendingKudo.targetUserId}
                    onchange={(event: Event) =>
                      updatePendingKudo(index, 'targetUserId', (event.currentTarget as HTMLSelectElement).value)}
                  >
                    <option value="">Select teammate</option>
                    {#each availableKudoUsers as teamUser}
                      <option value={teamUser.id}>{teamUser.name}</option>
                    {/each}
                  </select>
                </div>

                <div>
                  <label
                    for={`modal-kudo-comment-${index}`}
                    class="mb-2 block text-xs font-medium uppercase tracking-wider text-slate-500 dark:text-slate-400"
                  >
                    Kudos note
                  </label>
                  <textarea
                    id={`modal-kudo-comment-${index}`}
                    class="min-h-[96px] w-full rounded-xl border border-slate-300 bg-white px-3 py-2 text-sm text-slate-900 shadow-sm outline-none transition focus:border-indigo-500 focus:ring-2 focus:ring-indigo-200 dark:border-slate-600 dark:bg-slate-800 dark:text-white dark:focus:ring-indigo-500/30"
                    placeholder="Call out something helpful, thoughtful, or impactful."
                    value={pendingKudo.comment}
                    oninput={(event: Event) =>
                      updatePendingKudo(index, 'comment', (event.currentTarget as HTMLTextAreaElement).value)}
                  ></textarea>
                </div>

                <div class="flex justify-end lg:pt-7">
                  {#if kudosToCreate.length > 1}
                    <button
                      type="button"
                      class="inline-flex h-10 w-10 items-center justify-center rounded-xl border border-slate-200 text-slate-500 transition hover:border-red-200 hover:bg-red-50 hover:text-red-600 dark:border-slate-600 dark:text-slate-300 dark:hover:border-red-900/40 dark:hover:bg-red-900/20 dark:hover:text-red-400"
                      onclick={() => removeKudoRow(index)}
                      aria-label="Remove kudos row"
                    >
                      <X class="h-4 w-4" />
                    </button>
                  {/if}
                </div>
              </div>
            </div>
          {/each}

          <div class="flex flex-wrap items-center justify-between gap-3">
            <button
              type="button"
              class="inline-flex items-center gap-2 rounded-xl border border-dashed border-indigo-300 px-4 py-2 text-sm font-semibold text-indigo-600 transition hover:bg-indigo-50 dark:border-indigo-500/40 dark:text-indigo-300 dark:hover:bg-indigo-500/10"
              onclick={addKudoRow}
            >
              <Plus class="h-4 w-4" />
              Add teammate
            </button>

            <div class="flex items-center gap-2">
              <button
                type="button"
                class="rounded-xl border border-slate-300 px-3 py-2 text-sm font-semibold text-slate-700 transition hover:bg-slate-100 dark:border-slate-600 dark:text-slate-200 dark:hover:bg-slate-800"
                onclick={resetCreateForm}
              >
                {$LL.cancel()}
              </button>
              <button
                type="button"
                class="rounded-xl bg-indigo-600 px-3 py-2 text-sm font-semibold text-white transition hover:bg-indigo-500"
                onclick={saveNewKudos}
              >
                {$LL.save()}
              </button>
            </div>
          </div>
        </div>
      </div>
    {/if}

    <div class="space-y-3">
      {#each kudos as kudo}
        <div
          class="grid gap-4 rounded-2xl border border-slate-200 bg-slate-50 px-4 py-4 dark:border-slate-700 dark:bg-slate-900/30 md:grid-cols-[220px_minmax(0,1fr)] md:items-start"
        >
          <div class="flex min-w-0 items-start gap-3 border-slate-200 md:border-r md:pr-4 dark:md:border-slate-700">
            <div class="h-11 w-11 shrink-0 overflow-hidden rounded-full ring-2 ring-white shadow-sm dark:ring-gray-700">
              <UserAvatar
                width={44}
                warriorId={kudo.targetUser.id}
                avatar={kudo.targetUser.avatar}
                gravatarHash={kudo.targetUser.gravatarHash}
                userName={kudo.targetUser.name}
                options={{ class: 'h-full w-full rounded-full object-cover' }}
              />
            </div>

            <div class="min-w-0 flex-1">
              <div class="truncate text-base font-semibold text-slate-900 dark:text-white">{kudo.targetUser.name}</div>
              <div class="truncate text-sm text-slate-500 dark:text-slate-400">from {kudo.user.name}</div>
            </div>

            {#if canManageKudo(kudo)}
              <div class="shrink-0">
                <ActionsMenu
                  actions={[
                    {
                      label: $LL.edit(),
                      icon: Pencil,
                      onclick: () => startEditing(kudo),
                      disabled: editingKudoId !== '' && editingKudoId !== kudo.id,
                      testId: `kudo-edit-${kudo.id}`,
                    },
                    {
                      label: $LL.delete(),
                      icon: Trash2,
                      onclick: () => onDelete(kudo.id),
                      className: 'text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20',
                      disabled: editingKudoId !== '' && editingKudoId !== kudo.id,
                      testId: `kudo-delete-${kudo.id}`,
                    },
                  ]}
                  ariaLabel={`Settings for kudos from ${kudo.user.name} to ${kudo.targetUser.name}`}
                  testId={`kudo-actions-${kudo.id}`}
                  iconSize="medium"
                />
              </div>
            {/if}
          </div>

          <div class="min-w-0 text-sm text-slate-700 dark:text-slate-200">
            {#if editingKudoId === kudo.id}
              <label
                for={`kudo-comment-${kudo.id}`}
                class="mb-2 block text-[11px] font-semibold uppercase tracking-wide text-slate-500 dark:text-slate-400"
              >
                Note
              </label>
              <textarea
                id={`kudo-comment-${kudo.id}`}
                bind:value={draftComment}
                class="min-h-[110px] w-full rounded-xl border border-slate-300 bg-white px-3 py-2 text-sm text-slate-900 shadow-sm outline-none transition focus:border-indigo-500 focus:ring-2 focus:ring-indigo-200 dark:border-slate-600 dark:bg-slate-800 dark:text-white dark:focus:ring-indigo-500/30"
              ></textarea>

              <div class="mt-3 flex justify-end gap-2">
                <button
                  type="button"
                  class="rounded-xl border border-slate-300 px-3 py-2 text-sm font-semibold text-slate-700 transition hover:bg-slate-100 dark:border-slate-600 dark:text-slate-200 dark:hover:bg-slate-800"
                  onclick={cancelEditing}
                >
                  {$LL.cancel()}
                </button>
                <button
                  type="button"
                  class="rounded-xl bg-indigo-600 px-3 py-2 text-sm font-semibold text-white transition hover:bg-indigo-500 disabled:cursor-not-allowed disabled:opacity-60"
                  onclick={() => saveEditing(kudo.id)}
                  disabled={draftComment.trim() === ''}
                >
                  {$LL.save()}
                </button>
              </div>
            {:else}
              <div class="mb-2 text-[11px] font-semibold uppercase tracking-wide text-slate-500 dark:text-slate-400">
                Note
              </div>
              <div class="unreset whitespace-pre-wrap leading-relaxed">
                {@html kudo.comment}
              </div>
            {/if}
          </div>
        </div>
      {/each}
    </div>
  </div>
</Modal>
