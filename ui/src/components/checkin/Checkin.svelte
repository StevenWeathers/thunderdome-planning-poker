<script lang="ts">
  import Modal from '../global/Modal.svelte';
  import SolidButton from '../global/SolidButton.svelte';
  import LL from '../../i18n/i18n-svelte';
  import Editor from '../forms/Editor.svelte';
  import Toggle from '../forms/Toggle.svelte';
  import { AppConfig } from '../../config';
  import { user } from '../../stores';
  import { onMount } from 'svelte';
  import FeatureSubscribeBanner from '../global/FeatureSubscribeBanner.svelte';

  import type { NotificationService } from '../../types/notifications';
  import type { ApiClient } from '../../types/apiclient';
  import { Check, CircleCheck, CircleCheckBig, Save } from '@lucide/svelte';

  interface Props {
    toggleCheckin?: any;
    handleCheckin?: any;
    handleCheckinEdit?: any;
    userId: any;
    checkinId?: any;
    today?: string;
    yesterday?: string;
    blockers?: string;
    discuss?: string;
    goalsMet?: boolean;
    notifications: NotificationService;
    xfetch: ApiClient;
    teamPrefix?: string;
  }

  let {
    toggleCheckin = () => {},
    handleCheckin = () => {},
    handleCheckinEdit = () => {},
    userId,
    checkinId,
    today = $bindable(''),
    yesterday = $bindable(''),
    blockers = $bindable(''),
    discuss = $bindable(''),
    goalsMet = $bindable(true),
    notifications,
    xfetch,
    teamPrefix = '',
  }: Props = $props();

  let userSubscribed = $state(false);
  let lastCheckin = $state({
    id: '',
    yesterday: '',
    today: '',
    blockers: '',
    discuss: '',
    goalsMet: false,
  });

  function onSubmit(e: Event) {
    e.preventDefault();

    if (checkinId) {
      handleCheckinEdit(checkinId, {
        yesterday,
        today,
        blockers,
        discuss,
        goalsMet,
      });
    } else {
      handleCheckin({
        userId,
        yesterday,
        today,
        blockers,
        discuss,
        goalsMet,
      });
    }
  }

  function getLastCheckin() {
    xfetch(`${teamPrefix}/checkins/users/${userId}/last`)
      .then(res => res.json())
      .then(result => {
        if (!result.data) {
          return;
        }
        lastCheckin = result.data;
      })
      .catch(([err, response]) => {
        if (response.status === 204) {
          return;
        }
        notifications.danger('Error getting last checkin');
      });
  }

  onMount(() => {
    if (!AppConfig.SubscriptionsEnabled || (AppConfig.SubscriptionsEnabled && $user.subscribed)) {
      userSubscribed = true;
      getLastCheckin();
    }
  });
</script>

<Modal
  closeModal={toggleCheckin}
  widthClasses="md:w-4/5 lg:w-3/4 xl:w-2/3 max-w-6xl"
  ariaLabel={$LL.modalTeamCheckin()}
>
  <!-- Header -->
  <div class="py-5 mb-8">
    <div class="flex items-center space-x-4">
      <div
        class="w-12 h-12 bg-gradient-to-br from-blue-500 to-indigo-600 rounded-xl flex items-center justify-center shadow-sm"
      >
        <CircleCheckBig class="w-6 h-6 text-white" />
      </div>
      <div>
        <h1 class="text-3xl font-bold text-slate-900 dark:text-white font-rajdhani tracking-wide">Team Check-in</h1>
        <p class="text-slate-600 dark:text-slate-400 mt-1">
          {checkinId ? 'Update your progress' : 'Share your daily progress with your team'}
        </p>
      </div>
    </div>
  </div>

  <form onsubmit={onSubmit} name="teamCheckin" class="space-y-6">
    {#if userSubscribed}
      {#if lastCheckin.id !== '' && lastCheckin.id !== checkinId}
        <!-- Last Check-in Card -->
        <div
          class="bg-white dark:bg-slate-800 rounded-2xl shadow-sm border border-slate-200 dark:border-slate-700 overflow-hidden"
        >
          <div
            class="bg-gradient-to-r from-slate-50 to-slate-100 dark:from-slate-800 dark:to-slate-700 px-6 py-4 border-b border-slate-200 dark:border-slate-600"
          >
            <div class="flex items-center space-x-2">
              <svg
                class="w-5 h-5 text-slate-500 dark:text-slate-400"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
              <h2 class="text-lg font-semibold text-slate-900 dark:text-white font-rajdhani tracking-wide">
                Previous Check-in
              </h2>
            </div>
          </div>
          <div class="p-6">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div class="space-y-3">
                <div
                  class="text-xs font-medium text-slate-500 dark:text-slate-400 uppercase tracking-wider font-rajdhani"
                >
                  {$LL.yesterday()}
                </div>
                <div
                  class="text-slate-700 dark:text-slate-300 bg-slate-50 dark:bg-slate-700 rounded-lg p-4 prose prose-sm max-w-none"
                >
                  {@html lastCheckin.yesterday || '<span class="text-slate-400 italic">No content</span>'}
                </div>
              </div>
              <div class="space-y-3">
                <div
                  class="text-xs font-medium text-slate-500 dark:text-slate-400 uppercase tracking-wider font-rajdhani"
                >
                  {$LL.today()}
                </div>
                <div
                  class="text-slate-700 dark:text-slate-300 bg-slate-50 dark:bg-slate-700 rounded-lg p-4 prose prose-sm max-w-none"
                >
                  {@html lastCheckin.today || '<span class="text-slate-400 italic">No content</span>'}
                </div>
              </div>
              <div class="space-y-3">
                <div
                  class="text-xs font-medium text-red-500 dark:text-red-400 uppercase tracking-wider font-rajdhani flex items-center space-x-1"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                    />
                  </svg>
                  <span>{$LL.blockers()}</span>
                </div>
                <div
                  class="text-slate-700 dark:text-slate-300 bg-red-50 dark:bg-red-900/20 rounded-lg p-4 prose prose-sm max-w-none border border-red-100 dark:border-red-800"
                >
                  {@html lastCheckin.blockers || '<span class="text-slate-400 italic">No blockers</span>'}
                </div>
              </div>
              <div class="space-y-3">
                <div
                  class="text-xs font-medium text-emerald-600 dark:text-emerald-400 uppercase tracking-wider font-rajdhani flex items-center space-x-1"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"
                    />
                  </svg>
                  <span>{$LL.discuss()}</span>
                </div>
                <div
                  class="text-slate-700 dark:text-slate-300 bg-emerald-50 dark:bg-emerald-900/20 rounded-lg p-4 prose prose-sm max-w-none border border-emerald-100 dark:border-emerald-800"
                >
                  {@html lastCheckin.discuss || '<span class="text-slate-400 italic">Nothing to discuss</span>'}
                </div>
              </div>
            </div>
          </div>
        </div>
      {/if}
    {:else}
      <FeatureSubscribeBanner
        salesPitch="Build on yesterday's momentum - upgrade to view your last check-in and seamlessly continue your progress."
      />
    {/if}

    <!-- Main Form -->
    <div class="grid grid-cols-1 xl:grid-cols-2 gap-6">
      <!-- Yesterday & Goals Section -->
      <div class="space-y-4">
        <div
          class="bg-white dark:bg-slate-800 rounded-2xl shadow-sm border border-slate-200 dark:border-slate-700 overflow-hidden"
        >
          <div
            class="bg-gradient-to-r from-amber-50 to-orange-50 dark:from-amber-900/30 dark:to-orange-900/30 px-6 py-4 border-b border-amber-100 dark:border-amber-800"
          >
            <div class="flex items-center space-x-2">
              <svg
                class="w-5 h-5 text-amber-600 dark:text-amber-400"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
              <h3 class="text-lg font-semibold text-amber-900 dark:text-amber-100 font-rajdhani tracking-wide">
                {$LL.yesterday()}
              </h3>
            </div>
            <p class="text-sm text-amber-700 dark:text-amber-300 mt-1">What did you accomplish yesterday?</p>
          </div>
          <div class="p-6 bg-white">
            <Editor
              content={yesterday}
              placeholder={$LL.yesterdayPlaceholder()}
              id="yesterday"
              handleTextChange={(c: string) => (yesterday = c)}
            />
          </div>
        </div>

        <div class="bg-white dark:bg-slate-800 rounded-2xl shadow-sm border border-slate-200 dark:border-slate-700 p-6">
          <div class="flex items-start space-x-4">
            <div
              class="w-10 h-10 bg-gradient-to-br from-emerald-500 to-teal-600 rounded-xl flex items-center justify-center flex-shrink-0 mt-1"
            >
              <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
            </div>
            <div class="flex-1 min-w-0">
              <Toggle
                name="goalsMet"
                id="goalsMet"
                bind:checked={goalsMet}
                label={$LL.checkinMeetYesterdayGoalsQuestion()}
              />
            </div>
          </div>
        </div>
      </div>

      <!-- Today Section -->
      <div
        class="bg-white dark:bg-slate-800 rounded-2xl shadow-sm border border-slate-200 dark:border-slate-700 overflow-hidden"
      >
        <div
          class="bg-gradient-to-r from-blue-50 to-indigo-50 dark:from-blue-900/30 dark:to-indigo-900/30 px-6 py-4 border-b border-blue-100 dark:border-blue-800"
        >
          <div class="flex items-center space-x-2">
            <svg class="w-5 h-5 text-blue-600 dark:text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
            </svg>
            <h3 class="text-lg font-semibold text-blue-900 dark:text-blue-100 font-rajdhani tracking-wide">
              {$LL.today()}
            </h3>
          </div>
          <p class="text-sm text-blue-700 dark:text-blue-300 mt-1">What are your plans for today?</p>
        </div>
        <div class="p-6 bg-white">
          <Editor
            content={today}
            placeholder={$LL.todayPlaceholder()}
            id="today"
            handleTextChange={(c: string) => (today = c)}
          />
        </div>
      </div>
    </div>

    <!-- Blockers and Discussion Grid -->
    <div class="grid grid-cols-1 xl:grid-cols-2 gap-6">
      <!-- Blockers Section -->
      <div
        class="bg-white dark:bg-slate-800 rounded-2xl shadow-sm border border-slate-200 dark:border-slate-700 overflow-hidden"
      >
        <div
          class="bg-gradient-to-r from-red-50 to-rose-50 dark:from-red-900/30 dark:to-rose-900/30 px-6 py-4 border-b border-red-100 dark:border-red-800"
        >
          <div class="flex items-center space-x-2">
            <svg class="w-5 h-5 text-red-600 dark:text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
            <h3 class="text-lg font-semibold text-red-900 dark:text-red-100 font-rajdhani tracking-wide">
              {$LL.blockers()}
            </h3>
          </div>
          <p class="text-sm text-red-700 dark:text-red-300 mt-1">What's preventing you from making progress?</p>
        </div>
        <div class="p-6 bg-white">
          <Editor
            content={blockers}
            placeholder={$LL.blockersPlaceholder()}
            id="blockers"
            handleTextChange={(c: string) => (blockers = c)}
          />
        </div>
      </div>

      <!-- Discussion Section -->
      <div
        class="bg-white dark:bg-slate-800 rounded-2xl shadow-sm border border-slate-200 dark:border-slate-700 overflow-hidden"
      >
        <div
          class="bg-gradient-to-r from-emerald-50 to-teal-50 dark:from-emerald-900/30 dark:to-teal-900/30 px-6 py-4 border-b border-emerald-100 dark:border-emerald-800"
        >
          <div class="flex items-center space-x-2">
            <svg
              class="w-5 h-5 text-emerald-600 dark:text-emerald-400"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"
              />
            </svg>
            <h3 class="text-lg font-semibold text-emerald-900 dark:text-emerald-100 font-rajdhani tracking-wide">
              {$LL.discuss()}
            </h3>
          </div>
          <p class="text-sm text-emerald-700 dark:text-emerald-300 mt-1">What topics need team discussion?</p>
        </div>
        <div class="p-6 bg-white">
          <Editor
            content={discuss}
            placeholder={$LL.discussPlaceholder()}
            id="discuss"
            handleTextChange={(c: string) => (discuss = c)}
          />
        </div>
      </div>
    </div>

    <!-- Submit Button -->
    <div class="flex justify-end pt-6">
      <SolidButton type="submit" testid="save">
        <div class="flex items-center space-x-2">
          <Save class="w-5 h-5" />
          <span>{$LL.save()}</span>
        </div>
      </SolidButton>
    </div>
  </form>
</Modal>
