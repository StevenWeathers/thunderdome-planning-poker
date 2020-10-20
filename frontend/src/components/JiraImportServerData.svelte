<script>
    import SolidButton from './SolidButton.svelte'
    import CloseIcon from './icons/CloseIcon.svelte'
    import {_} from '../i18n'

    const apiEndPointConfigured = appConfig.JiraServerUrl

    export let toggleJiraRestConfig = () => {}
    export let handleImportFromRest = () => {}

    export let userName = ''
    export let password = ''
    export let apiEndpoint = apiEndPointConfigured
    export let jql = ''

    function takeOverConfig(e) {
        e.preventDefault()

        const restCfg = {
            userName,
            password,
            apiEndpoint,
            jql
        }

        handleImportFromRest(restCfg)
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
                            on:click="{toggleJiraRestConfig}"
                            class="text-gray-800">
                        <CloseIcon/>
                    </button>
                </div>
                <form on:submit="{takeOverConfig}" name="JiraRestConfig">
                    <div class="mb-4">
                        <label
                                class="block text-gray-700 text-sm font-bold mb-2"
                                for="userName">
                            {$_('pages.jiraRestCfg.fields.username.label')}
                        </label>
                        <div class="control">
                            <input
                                    name="userName"
                                    bind:value="{userName}"
                                    placeholder="{$_('pages.jiraRestCfg.fields.username.placeholder')}"
                                    class="bg-gray-200 border-gray-200 border-2
                                appearance-none rounded w-full py-2 px-3
                                text-gray-700 leading-tight focus:outline-none
                                focus:bg-white focus:border-purple-500"
                                    id="userName"
                                    required/>
                        </div>
                    </div>
                    <div class="mb-4">
                        <label
                                class="block text-gray-700 text-sm font-bold mb-2"
                                for="password">
                            {$_('pages.jiraRestCfg.fields.password.label')}
                        </label>
                        <div class="control">
                            <input
                                    bind:value="{password}"
                                    placeholder="{$_('pages.jiraRestCfg.fields.password.placeholder')}"
                                    class="bg-gray-200 border-gray-200 border-2 appearance-none rounded
                                w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none
                                focus:bg-white focus:border-purple-500"
                                    id="password"
                                    name="password"
                                    type="password"
                                    required/>
                        </div>
                    </div>
                    <div class="mb-4">
                        <label
                                class="block text-gray-700 text-sm font-bold mb-2"
                                for="apiEndpoint">
                            {$_('pages.jiraRestCfg.fields.apiEndpoint.label')}
                        </label>
                        <div class="control">
                            <input
                                    name="apiEndpoint"
                                    bind:value="{apiEndpoint}"
                                    placeholder="{$_('pages.jiraRestCfg.fields.apiEndpoint.placeholder')}"
                                    class="bg-gray-200 border-gray-200 border-2
                                appearance-none rounded w-full py-2 px-3
                                text-gray-700 leading-tight focus:outline-none
                                focus:bg-white focus:border-purple-500"
                                    id="apiEndpoint"
                                    required/>
                        </div>
                    </div>
                    <div class="mb-4">
                        <label
                                class="block text-gray-700 text-sm font-bold mb-2"
                                for="jql">
                            {$_('pages.jiraRestCfg.fields.jql.label')}
                        </label>
                        <div class="control">
                            <input
                                    name="jql"
                                    bind:value="{jql}"
                                    placeholder="{$_('pages.jiraRestCfg.fields.jql.placeholder')}"
                                    class="bg-gray-200 border-gray-200 border-2
                                appearance-none rounded w-full py-2 px-3
                                text-gray-700 leading-tight focus:outline-none
                                focus:bg-white focus:border-purple-500"
                                    id="jql"
                                    type="text"
                                    required/>
                        </div>
                    </div>

                    <div class="text-right">
                        <SolidButton type="submit">
                            {$_('actions.plan.importJiraRest.button')}
                        </SolidButton>
                    </div>
                </form>
            </div>
        </div>
    </div>
</div>
