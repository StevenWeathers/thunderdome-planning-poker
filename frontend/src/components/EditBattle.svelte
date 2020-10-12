<script>
    import SolidButton from './SolidButton.svelte'
    import CloseIcon from './icons/CloseIcon.svelte'
    import { _ } from '../i18n'

    const allowedPointValues = appConfig.AllowedPointValues

    export let toggleEditBattle = () => {}
    export let handleBattleEdit = () => {}
    export let points = []
    export let battleName = ''
    export let votingLocked = false

    let checkedPointColor = 'border-green-500 bg-green-100 text-green-600'
    let uncheckedPointColor = 'border-gray-300 bg-white'

    function saveBattle(e) {
        e.preventDefault()

        const pointValuesAllowed = allowedPointValues.filter(pv => {
            return points.includes(pv)
        })

        const battle = {
            battleName,
            pointValuesAllowed,
        }

        handleBattleEdit(battle)
    }
</script>

<div
    class="fixed inset-0 flex items-center z-40 max-h-screen overflow-y-scroll">
    <div class="fixed inset-0 bg-gray-900 opacity-75"></div>

    <div
        class="relative mx-4 md:mx-auto w-full md:w-2/3 lg:w-3/5 xl:w-1/2 z-50
        max-h-full">
        <div class="py-8">
            <div class="shadow-xl bg-white rounded-lg p-4 xl:p-6 max-h-full">
                <div class="flex justify-end mb-2">
                    <button
                        aria-label="close"
                        on:click="{toggleEditBattle}"
                        class="text-gray-800">
                        <CloseIcon />
                    </button>
                </div>
                <form on:submit="{saveBattle}" name="createBattle">
                    <div class="mb-4">
                        <label
                            class="block text-gray-700 text-sm font-bold mb-2"
                            for="battleName">
                            {$_('pages.myBattles.createBattle.fields.name.label')}
                        </label>
                        <div class="control">
                            <input
                                name="battleName"
                                bind:value="{battleName}"
                                placeholder="{$_('pages.myBattles.createBattle.fields.name.placeholder')}"
                                class="bg-gray-200 border-gray-200 border-2
                                appearance-none rounded w-full py-2 px-3
                                text-gray-700 leading-tight focus:outline-none
                                focus:bg-white focus:border-purple-500"
                                id="battleName"
                                required />
                        </div>
                    </div>

                    <div class="mb-4">
                        <h3 class="block text-gray-700 text-sm font-bold mb-2">
                            {$_('pages.myBattles.createBattle.fields.allowedPointValues.label')}
                        </h3>
                        <div class="control relative -mr-2 md:-mr-1">
                            {#if !votingLocked}
                                <div class="font-bold text-red-500">
                                    {$_('actions.battle.editPointsDisabled')}
                                </div>
                            {/if}
                            {#each allowedPointValues as point}
                                <label
                                    class="
                                    {points.includes(point) ? checkedPointColor : uncheckedPointColor}
                                    cursor-pointer font-bold border p-2 mr-2
                                    xl:mr-1 mb-2 xl:mb-0 rounded inline-block {!votingLocked ? 'opacity-25 cursor-not-allowed' : 'cursor-pointer'}">
                                    <input
                                        type="checkbox"
                                        bind:group="{points}"
                                        value="{point}"
                                        class="hidden"
                                        disabled="{!votingLocked}" />
                                    {point}
                                </label>
                            {/each}
                        </div>
                    </div>

                    <div class="text-right">
                        <SolidButton type="submit">
                            {$_('actions.battle.save')}
                        </SolidButton>
                    </div>
                </form>
            </div>
        </div>
    </div>
</div>
