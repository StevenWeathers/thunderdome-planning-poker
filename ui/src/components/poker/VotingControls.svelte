<script lang="ts">
  import SolidButton from '../SolidButton.svelte';
  import LL from '../../i18n/i18n-svelte';
  import { AppConfig } from '../../config';
  import SelectInput from '../SelectInput.svelte';

  export let sendSocketEvent = () => {};
  export let eventTag;
  export let planId = '';
  export let points = [];
  export let votingLocked = true;
  export let highestVote = '';

  $: planPoints = highestVote;

  const endPlanVoting = () => {
    sendSocketEvent('end_voting', planId);
    eventTag('vote_end', 'battle', '');
  };

  const skipPlan = () => {
    sendSocketEvent('skip_plan', planId);
    eventTag('plan_skip', 'battle', '');
  };

  const restartVoting = () => {
    sendSocketEvent('activate_plan', planId);
    eventTag('plan_restart_vote', 'battle', '');
  };

  function handleSubmit(event) {
    event.preventDefault();

    sendSocketEvent(
      'finalize_plan',
      JSON.stringify({
        planId,
        planPoints,
      }),
    );
    eventTag('plan_finalize', 'battle', planPoints);

    planPoints = '';
  }
</script>

{#if planId != ''}
  <div class="p-4" data-testId="votingControls">
    <SolidButton
      color="blue"
      additionalClasses="mb-2 w-full"
      onClick="{skipPlan}"
      testid="voting-skip"
    >
      {$LL.planSkip({ friendly: AppConfig.FriendlyUIVerbs })}
    </SolidButton>
    {#if !votingLocked}
      <SolidButton
        additionalClasses="w-full"
        onClick="{endPlanVoting}"
        testid="voting-finish"
      >
        {$LL.votingFinish()}
      </SolidButton>
    {:else}
      <SolidButton
        color="blue"
        additionalClasses="mb-2 w-full"
        onClick="{restartVoting}"
        testid="voting-restart"
      >
        {$LL.votingRestart()}
      </SolidButton>
      <form on:submit="{handleSubmit}" name="savePlanPoints">
        <legend
          class="text-xl mb-2 font-semibold leading-tight dark:text-gray-300"
        >
          {$LL.finalPoints()}
        </legend>
        <div class="flex -mx-2">
          <div class="w-1/2 px-2">
            <SelectInput name="planPoints" bind:value="{planPoints}" required>
              <option value="" disabled>
                {$LL.points()}
              </option>
              >
              {#each points as point}
                <option value="{point}">{point}</option>
              {/each}
            </SelectInput>
          </div>
          <div class="w-1/2 text-right px-2">
            <SolidButton
              additionalClasses="w-full h-full"
              type="submit"
              testid="voting-save"
            >
              {$LL.save()}
            </SolidButton>
          </div>
        </div>
      </form>
    {/if}
  </div>
{/if}
