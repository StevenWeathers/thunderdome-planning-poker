<script>
    import SolidButton from '../SolidButton.svelte'
    import Modal from '../Modal.svelte'
    import DownCarrotIcon from '../icons/ChevronDown.svelte'
    import { AppConfig } from '../../config.js'
    import { _ } from '../../i18n.js'

    const allowedPointValues = AppConfig.AllowedPointValues
    const allowedPointAverages = ['ceil', 'round', 'floor']

    export let toggleEditBattle = () => {}
    export let handleBattleEdit = () => {}
    export let points = []
    export let battleName = ''
    export let votingLocked = false
    export let autoFinishVoting = true
    export let pointAverageRounding = 'ceil'
    export let joinCode = ''
    export let leaderCode = ''

    let checkedPointColor =
        'border-green-500 bg-green-100 text-green-600 dark:bg-gray-900 dark:text-lime-500 dark:border-lime-500'
    let uncheckedPointColor =
        'border-gray-300 bg-white dark:bg-gray-900 dark:border-gray-600 dark:text-gray-300'

    function saveBattle(e) {
        e.preventDefault()

        const pointValuesAllowed = allowedPointValues.filter(pv => {
            return points.includes(pv)
        })

        const battle = {
            battleName,
            pointValuesAllowed,
            autoFinishVoting,
            pointAverageRounding,
            joinCode,
            leaderCode,
        }

        handleBattleEdit(battle)
    }
</script>

<Modal
    closeModal="{toggleEditBattle}"
    widthClasses="md:w-2/3 lg:w-3/5 xl:w-1/2"
>
    <form on:submit="{saveBattle}" name="createBattle">
        <div class="mb-4">
            <label
                class="block text-gray-700 dark:text-gray-400 text-sm font-bold mb-2"
                for="battleName"
            >
                {$_('pages.myBattles.createBattle.fields.name.label')}
            </label>
            <div class="control">
                <input
                    name="battleName"
                    bind:value="{battleName}"
                    placeholder="{$_(
                        'pages.myBattles.createBattle.fields.name.placeholder',
                    )}"
                    class="bg-gray-100 dark:bg-gray-900 border-gray-200 dark:border-gray-800 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 dark:text-gray-300 leading-tight
                focus:outline-none focus:bg-white dark:focus:bg-gray-700 focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
                    id="battleName"
                    required
                />
            </div>
        </div>

        <div class="mb-4">
            <h3
                class="block text-gray-700 dark:text-gray-400 text-sm font-bold mb-2"
            >
                {$_(
                    'pages.myBattles.createBattle.fields.allowedPointValues.label',
                )}
            </h3>
            <div class="control relative -mr-2 md:-mr-1">
                {#if !votingLocked}
                    <div class="font-bold text-red-500">
                        {$_('battleEditPointsDisabled')}
                    </div>
                {/if}
                {#each allowedPointValues as point}
                    <label
                        class="
                        {points.includes(point)
                            ? checkedPointColor
                            : uncheckedPointColor}
                        cursor-pointer font-bold border p-2 mr-2 xl:mr-1 mb-2
                        xl:mb-0 rounded inline-block {!votingLocked
                            ? 'opacity-25 cursor-not-allowed'
                            : 'cursor-pointer'}"
                    >
                        <input
                            type="checkbox"
                            bind:group="{points}"
                            value="{point}"
                            class="hidden"
                            disabled="{!votingLocked}"
                        />
                        {point}
                    </label>
                {/each}
            </div>
        </div>

        <div class="mb-4">
            <label
                class="text-gray-700 dark:text-gray-400 text-sm font-bold mb-2"
                for="averageRounding"
            >
                {$_(
                    'pages.myBattles.createBattle.fields.averageRounding.label',
                )}
            </label>
            <div class="relative">
                <select
                    bind:value="{pointAverageRounding}"
                    class="block appearance-none w-full border-2 border-gray-300 dark:border-gray-700
                text-gray-700 dark:text-gray-300 py-3 px-4 pr-8 rounded leading-tight
                focus:outline-none focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400 dark:bg-gray-900"
                    id="averageRounding"
                    name="averageRounding"
                >
                    {#each allowedPointAverages as item}
                        <option value="{item}">
                            {$_(
                                'pages.myBattles.createBattle.fields.averageRounding.' +
                                    item,
                            )}
                        </option>
                    {/each}
                </select>
                <div
                    class="pointer-events-none absolute inset-y-0 right-0 flex
                    items-center px-2 text-gray-700 dark:text-gray-400"
                >
                    <DownCarrotIcon />
                </div>
            </div>
        </div>

        <div class="mb-4">
            <label
                class="text-gray-700 dark:text-gray-400 text-sm font-bold mb-2"
            >
                <input
                    type="checkbox"
                    bind:checked="{autoFinishVoting}"
                    id="autoFinishVoting"
                    name="autoFinishVoting"
                    disabled="{!votingLocked}"
                    class="w-4 h-4 dark:accent-lime-400 mr-1"
                />
                {$_(
                    'pages.myBattles.createBattle.fields.autoFinishVoting.label',
                )}
            </label>
        </div>

        <div class="mb-4">
            <label
                class="block text-gray-700 dark:text-gray-400 text-sm font-bold mb-2"
                for="joinCode"
            >
                {$_('passCode')}
            </label>
            <div class="control">
                <input
                    name="joinCode"
                    bind:value="{joinCode}"
                    placeholder="{$_('optionalPasscodePlaceholder')}"
                    class="bg-gray-100 dark:bg-gray-900 border-gray-200 dark:border-gray-800 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 dark:text-gray-300 leading-tight
                focus:outline-none focus:bg-white dark:focus:bg-gray-700 focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
                    id="joinCode"
                />
            </div>
        </div>

        <div class="mb-4">
            <label
                class="block text-gray-700 dark:text-gray-400 text-sm font-bold mb-2"
                for="leaderCode"
            >
                {$_('leaderPasscode')}
            </label>
            <div class="control">
                <input
                    name="leaderCode"
                    bind:value="{leaderCode}"
                    placeholder="{$_('optionalLeadercodePlaceholder')}"
                    class="bg-gray-100 dark:bg-gray-900 border-gray-200 dark:border-gray-800 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 dark:text-gray-300 leading-tight
                focus:outline-none focus:bg-white dark:focus:bg-gray-700 focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
                    id="leaderCode"
                />
            </div>
        </div>

        <div class="text-right">
            <SolidButton type="submit">{$_('save')}</SolidButton>
        </div>
    </form>
</Modal>
