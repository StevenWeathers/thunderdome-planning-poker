<script lang="ts">
  import SolidButton from '../global/SolidButton.svelte';
  import LL from '../../i18n/i18n-svelte';
  import SelectInput from '../forms/SelectInput.svelte';
  import TextInput from '../forms/TextInput.svelte';

  export let sendSocketEvent = () => {};
  export let planId = '';
  export let points = [];
  export let votingLocked = true;
  export let highestVote = '';

  let customPointValue = false;

  $: planPoints = highestVote;
  let customPlanPoints = '';

  const toggleCustomPointValue = () => {
    if (planPoints === 'CUSTOM') {
      customPointValue = true;
    } else {
      customPointValue = false;
    }
  };

  const endPlanVoting = () => {
    sendSocketEvent('end_voting', planId);
  };

  const skipPlan = () => {
    sendSocketEvent('skip_plan', planId);
  };

  const restartVoting = () => {
    sendSocketEvent('activate_plan', planId);
  };

  function handleSubmit(event) {
    event.preventDefault();

    sendSocketEvent(
      'finalize_plan',
      JSON.stringify({
        planId,
        planPoints: customPlanPoints === '' ? planPoints : customPlanPoints,
      }),
    );

    planPoints = '';
    customPlanPoints = '';
    toggleCustomPointValue();
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
      {$LL.planSkip()}
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
        <div class="-mx-2">
          <div class="mb-2">
            <SelectInput
              name="planPoints"
              bind:value="{planPoints}"
              on:change="{toggleCustomPointValue}"
              required
            >
              <option value="" disabled>
                {$LL.points()}
              </option>
              >
              {#each points as point}
                <option value="{point}">{point}</option>
              {/each}
              <option value="CUSTOM">Custom</option>
            </SelectInput>
            {#if customPointValue}
              <TextInput
                name="customPlanPoints"
                bind:value="{customPlanPoints}"
                placeholder="enter a custom point value..."
                id="customPlanPoints"
                class="mt-2"
              />
            {/if}
          </div>
          <div>
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
