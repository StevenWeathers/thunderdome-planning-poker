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

  interface Props {
    toggleCheckin?: any;
    handleCheckin?: any;
    handleCheckinEdit?: any;
    userId: any;
    checkinId: any;
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
    teamPrefix = ''
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

  function onSubmit(e) {
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
    if (
      !AppConfig.SubscriptionsEnabled ||
      (AppConfig.SubscriptionsEnabled && $user.subscribed)
    ) {
      userSubscribed = true;
      getLastCheckin();
    }
  });
</script>

<Modal closeModal={toggleCheckin} widthClasses="md:w-2/3">
  <form onsubmit={onSubmit} name="teamCheckin" class="flex flex-wrap mt-8">
    {#if userSubscribed}
      {#if lastCheckin.id !== ''}
        <div
          class="w-full mb-4 p-4 border border-gray-300 dark:border-gray-600 rounded"
        >
          <div class="text-gray-500 dark:text-gray-300">
            <div
              class="mb-2 font-bold text-lg uppercase font-rajdhani tracking-wide border-b border-gray-200 dark:border-gray-700"
            >
              Last Checkin
            </div>
            <div class="w-full md:grid md:grid-cols-2 md:gap-4">
              <div>
                <div
                  class="text-gray-500 dark:text-gray-300 uppercase font-rajdhani tracking-wide mb-2"
                >
                  {$LL.yesterday()}
                </div>
                <div>{@html lastCheckin.yesterday}</div>
              </div>
              <div>
                <div
                  class="text-gray-500 dark:text-gray-300 uppercase font-rajdhani tracking-wide mb-2"
                >
                  {$LL.today()}
                </div>
                <div>{@html lastCheckin.today}</div>
              </div>
              <div>
                <div
                  class="text-gray-500 dark:text-gray-300 uppercase font-rajdhani tracking-wide mb-2"
                >
                  {$LL.blockers()}
                </div>
                <div>{@html lastCheckin.blockers}</div>
              </div>
              <div>
                <div
                  class="text-gray-500 dark:text-gray-300 uppercase font-rajdhani tracking-wide mb-2"
                >
                  {$LL.discuss()}
                </div>
                <div>{@html lastCheckin.discuss}</div>
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

    <div class="w-full md:grid md:grid-cols-2 md:gap-4">
      <div>
        <div class="mb-2">
          <div
            class="text-gray-500 dark:text-gray-300 uppercase font-rajdhani tracking-wide text-2xl mb-2"
          >
            {$LL.yesterday()}
          </div>
          <div class="bg-white">
            <Editor
              content={yesterday}
              placeholder={$LL.yesterdayPlaceholder()}
              id="yesterday"
              handleTextChange={c => (yesterday = c)}
            />
          </div>
        </div>
        <div class="mb-4">
          <Toggle
            name="goalsMet"
            id="goalsMet"
            bind:checked={goalsMet}
            label={$LL.checkinMeetYesterdayGoalsQuestion()}
          />
        </div>
      </div>
      <div>
        <div class="mb-4">
          <div
            class="text-gray-500 dark:text-gray-300 uppercase font-rajdhani tracking-wide text-2xl mb-2"
          >
            {$LL.today()}
          </div>
          <div class="bg-white">
            <Editor
              content={today}
              placeholder={$LL.todayPlaceholder()}
              id="today"
              handleTextChange={c => (today = c)}
            />
          </div>
        </div>
      </div>

      <div class="mb-4">
        <div
          class="text-red-500 uppercase font-rajdhani tracking-wide text-2xl mb-2"
        >
          {$LL.blockers()}
        </div>
        <div class="bg-white">
          <Editor
            content={blockers}
            placeholder={$LL.blockersPlaceholder()}
            id="blockers"
            handleTextChange={c => (blockers = c)}
          />
        </div>
      </div>

      <div class="mb-4">
        <div
          class="text-green-500 uppercase font-rajdhani tracking-wide text-2xl mb-2"
        >
          {$LL.discuss()}
        </div>
        <div class="bg-white">
          <Editor
            content={discuss}
            placeholder={$LL.discussPlaceholder()}
            id="discuss"
            handleTextChange={c => (discuss = c)}
          />
        </div>
      </div>
    </div>

    <div class="w-full">
      <div class="text-right">
        <SolidButton type="submit" testid="save">{$LL.save()}</SolidButton>
      </div>
    </div>
  </form>
</Modal>
