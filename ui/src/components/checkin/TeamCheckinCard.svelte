<script lang="ts">
  import { ChevronDown, ChevronUp, Pencil, Trash2 } from '@lucide/svelte';

  import ActionsMenu from '../global/ActionsMenu.svelte';
  import UserAvatar from '../user/UserAvatar.svelte';
  import Comments from './Comments.svelte';
  import LL from '../../i18n/i18n-svelte';

  import type { TeamUser, TeamCheckin as BaseTeamCheckin } from '../../types/team';
  import type { UserDisplay } from '../../types/user';

  interface TeamCheckin extends Omit<BaseTeamCheckin, 'user'> {
    user: TeamUser;
  }

  interface Props {
    checkin: TeamCheckin;
    currentUserId: string;
    isAdmin: boolean;
    isExpanded: boolean;
    userMap: Map<string, UserDisplay>;
    onToggleSummary: (checkin: TeamCheckin) => void;
    onEdit: (checkin: TeamCheckin) => void;
    onDelete: (checkinId: string) => void;
    onCreateComment: (checkinId: string, comment: any) => void;
    onEditComment: (checkinId: string, commentId: string, comment: any) => void;
    onDeleteComment: (checkinId: string, commentId: string) => void;
  }

  let {
    checkin,
    currentUserId,
    isAdmin,
    isExpanded,
    userMap,
    onToggleSummary,
    onEdit,
    onDelete,
    onCreateComment,
    onEditComment,
    onDeleteComment,
  }: Props = $props();

  const canManageCheckin = $derived(checkin.user.id === currentUserId || isAdmin);
</script>

<article
  class="group relative overflow-hidden rounded-2xl border border-gray-200/60 bg-white shadow-sm hover:shadow-lg dark:border-gray-700/60 dark:bg-gray-800 dark:shadow-gray-900/10"
  data-testid="checkin"
  aria-labelledby={`checkin-user-${checkin.user.id}`}
>
  <div class="relative p-4 sm:p-5">
    <div class="flex flex-col gap-5 lg:flex-row lg:items-start">
      <aside class="lg:w-52 lg:shrink-0">
        <div class="flex items-start gap-4 lg:flex-col lg:gap-3">
          <div class="flex items-start gap-3 lg:w-full lg:justify-between">
            <div class="relative h-16 w-16 shrink-0 sm:h-20 sm:w-20">
              <div class="h-full w-full overflow-hidden rounded-full ring-3 ring-white shadow-lg dark:ring-gray-700">
                <UserAvatar
                  width={80}
                  warriorId={checkin.user.id}
                  avatar={checkin.user.avatar}
                  gravatarHash={checkin.user.gravatarHash}
                  userName={checkin.user.name}
                  options={{
                    class:
                      'h-full w-full rounded-full object-cover transition-transform duration-300 group-hover:scale-110',
                  }}
                />
              </div>
            </div>

            {#if canManageCheckin}
              <div class="shrink-0 lg:self-start">
                <ActionsMenu
                  actions={[
                    {
                      label: $LL.edit(),
                      icon: Pencil,
                      onclick: () => onEdit(checkin),
                      testId: 'checkin-edit',
                    },
                    {
                      label: $LL.delete(),
                      icon: Trash2,
                      onclick: () => onDelete(checkin.id),
                      className: 'text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20',
                      testId: 'checkin-delete',
                    },
                  ]}
                  ariaLabel={`Settings for ${checkin.user.name}`}
                  testId="checkin-actions-menu"
                  iconSize="medium"
                />
              </div>
            {/if}
          </div>

          <div class="min-w-0 flex-1 space-y-3 lg:w-full">
            <div class="min-w-0 flex-1">
              <h3
                id={`checkin-user-${checkin.user.id}`}
                class="truncate text-lg font-semibold text-gray-900 dark:text-white sm:text-xl"
                data-testid="checkin-username"
              >
                {checkin.user.name}
              </h3>
              <div class="mt-2 space-y-2 text-xs font-semibold uppercase tracking-wide">
                <div class="flex flex-nowrap items-center gap-2">
                  {#if checkin.goalsMet}
                    <span
                      class="whitespace-nowrap rounded-full bg-emerald-50 px-2.5 py-1 text-emerald-700 dark:bg-emerald-900/20 dark:text-emerald-300"
                    >
                      {$LL.goalsMet()}
                    </span>
                  {/if}
                  {#if checkin.blockers !== ''}
                    <span
                      class="whitespace-nowrap rounded-full bg-red-50 px-2.5 py-1 text-red-700 dark:bg-red-900/20 dark:text-red-300"
                    >
                      {$LL.blocked()}
                    </span>
                  {/if}
                </div>
                {#if checkin.discuss !== ''}
                  <div class="flex items-center gap-2">
                    <span
                      class="whitespace-nowrap rounded-full bg-amber-50 px-2.5 py-1 text-amber-700 dark:bg-amber-900/20 dark:text-amber-300"
                    >
                      {$LL.discuss()}
                    </span>
                  </div>
                {/if}
              </div>
            </div>

            <Comments
              {checkin}
              {userMap}
              {isAdmin}
              handleCreate={onCreateComment}
              handleEdit={onEditComment}
              handleDelete={onDeleteComment}
            />
          </div>
        </div>
      </aside>

      <div class="min-w-0 flex-1 space-y-4">
        <section
          class="overflow-hidden rounded-xl border border-slate-200/80 bg-slate-50/80 dark:border-slate-700/70 dark:bg-slate-900/30"
        >
          <button
            class="flex w-full items-center justify-between gap-3 px-4 py-3 text-left transition-colors duration-200 hover:bg-slate-100/80 dark:hover:bg-slate-800/60"
            type="button"
            aria-expanded={isExpanded}
            aria-controls={`checkin-summary-${checkin.id}`}
            onclick={() => onToggleSummary(checkin)}
          >
            <div class="min-w-0">
              <div class="text-sm font-semibold uppercase tracking-wide text-gray-700 dark:text-gray-200">
                Daily update
              </div>
              <div class="text-xs text-gray-500 dark:text-gray-400">
                {isExpanded ? 'Hide Yesterday and Today' : 'Show Yesterday and Today'}
              </div>
            </div>

            <div
              class="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-white text-gray-500 shadow-sm dark:bg-gray-800 dark:text-gray-300"
            >
              {#if isExpanded}
                <ChevronUp class="h-4 w-4" />
              {:else}
                <ChevronDown class="h-4 w-4" />
              {/if}
            </div>
          </button>

          {#if isExpanded}
            <div
              id={`checkin-summary-${checkin.id}`}
              class="space-y-3 border-t border-slate-200/80 px-4 py-4 dark:border-slate-700/70"
            >
              <section class="space-y-2">
                <h4
                  class="flex items-center gap-2 text-sm font-semibold uppercase tracking-wide text-gray-600 dark:text-gray-300"
                >
                  <div class="h-2 w-2 rounded-full bg-blue-500"></div>
                  {$LL.yesterday()}
                </h4>
                <div
                  class="unreset whitespace-pre-wrap text-sm leading-relaxed text-gray-800 dark:text-gray-200"
                  data-testid="checkin-yesterday"
                >
                  {@html checkin.yesterday}
                </div>
              </section>

              <section class="space-y-2 border-t border-slate-200/70 pt-3 dark:border-slate-700/70">
                <h4
                  class="flex items-center gap-2 text-sm font-semibold uppercase tracking-wide text-gray-600 dark:text-gray-300"
                >
                  <div class="h-2 w-2 rounded-full bg-emerald-500"></div>
                  {$LL.today()}
                </h4>
                <div
                  class="unreset whitespace-pre-wrap text-sm leading-relaxed text-gray-800 dark:text-gray-200"
                  data-testid="checkin-today"
                >
                  {@html checkin.today}
                </div>
              </section>
            </div>
          {/if}
        </section>

        {#if checkin.blockers !== ''}
          <section
            class="space-y-2 rounded-xl border border-red-100 bg-red-50/70 p-4 dark:border-red-900/40 dark:bg-red-900/10"
          >
            <h4
              class="flex items-center gap-2 text-sm font-semibold uppercase tracking-wide text-red-600 dark:text-red-400"
            >
              <div class="h-2 w-2 animate-pulse rounded-full bg-red-500"></div>
              <span class="flex items-center gap-1">
                {$LL.blockers()}
                <svg class="h-4 w-4" fill="currentColor" viewBox="0 0 20 20" aria-hidden="true">
                  <path
                    fill-rule="evenodd"
                    d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z"
                    clip-rule="evenodd"
                  />
                </svg>
              </span>
            </h4>
            <div
              class="unreset whitespace-pre-wrap text-sm leading-relaxed text-gray-800 dark:text-gray-200"
              data-testid="checkin-blockers"
            >
              {@html checkin.blockers}
            </div>
          </section>
        {/if}

        {#if checkin.discuss !== ''}
          <section
            class="space-y-2 rounded-xl border border-amber-100 bg-amber-50/70 p-4 dark:border-amber-900/40 dark:bg-amber-900/10"
          >
            <h4
              class="flex items-center gap-2 text-sm font-semibold uppercase tracking-wide text-amber-600 dark:text-amber-400"
            >
              <div class="h-2 w-2 rounded-full bg-amber-500"></div>
              <span class="flex items-center gap-1">
                {$LL.discuss()}
                <svg class="h-4 w-4" fill="currentColor" viewBox="0 0 20 20" aria-hidden="true">
                  <path
                    fill-rule="evenodd"
                    d="M18 10c0 3.866-3.582 7-8 7a8.841 8.841 0 01-4.083-.98L2 17l1.338-3.123C2.493 12.767 2 11.434 2 10c0-3.866 3.582-7 8-7s8 3.134 8 7zM7 9H5v2h2V9zm8 0h-2v2h2V9zM9 9h2v2H9V9z"
                    clip-rule="evenodd"
                  />
                </svg>
              </span>
            </h4>
            <div
              class="unreset whitespace-pre-wrap text-sm leading-relaxed text-gray-800 dark:text-gray-200"
              data-testid="checkin-discuss"
            >
              {@html checkin.discuss}
            </div>
          </section>
        {/if}
      </div>
    </div>
  </div>
</article>
