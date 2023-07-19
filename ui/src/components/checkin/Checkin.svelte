<script lang="ts">
  import Modal from '../Modal.svelte';
  import SolidButton from '../SolidButton.svelte';
  import { quill } from '../../quill';
  import LL from '../../i18n/i18n-svelte';

  export let toggleCheckin = () => {};
  export let handleCheckin = () => {};
  export let handleCheckinEdit = () => {};
  export let userId;
  export let checkinId;
  export let today = '';
  export let yesterday = '';
  export let blockers = '';
  export let discuss = '';
  export let goalsMet = true;

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
</script>

<Modal closeModal="{toggleCheckin}" widthClasses="md:w-2/3">
  <form on:submit="{onSubmit}" name="teamCheckin" class="flex flex-wrap mt-8">
    <div class="w-full md:grid md:grid-cols-2 md:gap-4">
      <div>
        <div class="mb-2">
          <div
            class="text-gray-500 dark:text-gray-300 uppercase font-rajdhani tracking-wide text-2xl mb-2"
          >
            {$LL.yesterday()}
          </div>
          <div class="bg-white">
            <div
              class="w-full"
              use:quill="{{
                placeholder: `${$LL.yesterdayPlaceholder()}`,
                content: yesterday,
              }}"
              on:text-change="{e => (yesterday = e.detail.html)}"
              id="yesterday"
            ></div>
          </div>
        </div>
        <div class="mb-4">
          <div
            class="text-gray-600 dark:text-gray-400 uppercase font-rajdhani text-xl tracking-wide mb-2"
          >
            {$LL.checkinMeetYesterdayGoalsQuestion()}
          </div>
          <div
            class="relative inline-block w-16 me-2 align-middle select-none transition duration-200 ease-in"
          >
            <input
              type="checkbox"
              name="goalsMet"
              id="goalsMet"
              bind:checked="{goalsMet}"
              class="toggle-checkbox absolute block w-8 h-8 rounded-full bg-white border-4 border-gray-300 appearance-none cursor-pointer transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500 shadow"
            />
            <label
              for="goalsMet"
              class="toggle-label block overflow-hidden h-8 rounded-full bg-gray-300 cursor-pointer transition-colors duration-200 ease-in-out"
            >
            </label>
          </div>
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
            <div
              class="w-full"
              use:quill="{{
                placeholder: `${$LL.todayPlaceholder()}`,
                content: today,
              }}"
              on:text-change="{e => (today = e.detail.html)}"
              id="today"
            ></div>
          </div>
        </div>
      </div>
    </div>

    <div class="w-full mb-4">
      <div
        class="text-red-500 uppercase font-rajdhani tracking-wide text-2xl mb-2"
      >
        {$LL.blockers()}
      </div>
      <div class="bg-white">
        <div
          class="w-full"
          use:quill="{{
            placeholder: `${$LL.blockersPlaceholder()}`,
            content: blockers,
          }}"
          on:text-change="{e => (blockers = e.detail.html)}"
          id="blockers"
        ></div>
      </div>
    </div>

    <div class="w-full mb-4">
      <div
        class="text-green-500 uppercase font-rajdhani tracking-wide text-2xl mb-2"
      >
        {$LL.discuss()}
      </div>
      <div class="bg-white">
        <div
          class="w-full"
          use:quill="{{
            placeholder: `${$LL.discussPlaceholder()}`,
            content: discuss,
          }}"
          on:text-change="{e => (discuss = e.detail.html)}"
          id="discuss"
        ></div>
      </div>
    </div>

    <div class="w-full">
      <div class="text-right">
        <SolidButton type="submit" testid="save">{$LL.save()}</SolidButton>
      </div>
    </div>
  </form>
</Modal>
