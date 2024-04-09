<script lang="ts">
  import Modal from '../global/Modal.svelte';
  import SolidButton from '../global/SolidButton.svelte';
  import LL from '../../i18n/i18n-svelte';
  import Editor from '../forms/Editor.svelte';
  import Toggle from '../forms/Toggle.svelte';

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
            <Editor
              content="{yesterday}"
              placeholder="{$LL.yesterdayPlaceholder()}"
              id="yesterday"
              handleTextChange="{c => (yesterday = c)}"
            />
          </div>
        </div>
        <div class="mb-4">
          <Toggle
            name="goalsMet"
            id="goalsMet"
            bind:checked="{goalsMet}"
            label="{$LL.checkinMeetYesterdayGoalsQuestion()}"
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
              content="{today}"
              placeholder="{$LL.todayPlaceholder()}"
              id="today"
              handleTextChange="{c => (today = c)}"
            />
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
        <Editor
          content="{blockers}"
          placeholder="{$LL.blockersPlaceholder()}"
          id="blockers"
          handleTextChange="{c => (blockers = c)}"
        />
      </div>
    </div>

    <div class="w-full mb-4">
      <div
        class="text-green-500 uppercase font-rajdhani tracking-wide text-2xl mb-2"
      >
        {$LL.discuss()}
      </div>
      <div class="bg-white">
        <Editor
          content="{discuss}"
          placeholder="{$LL.discussPlaceholder()}"
          id="discuss"
          handleTextChange="{c => (discuss = c)}"
        />
      </div>
    </div>

    <div class="w-full">
      <div class="text-right">
        <SolidButton type="submit" testid="save">{$LL.save()}</SolidButton>
      </div>
    </div>
  </form>
</Modal>
