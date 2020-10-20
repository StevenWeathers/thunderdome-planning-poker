<script>
    import HollowButton from './HollowButton.svelte'
    import {_} from '../i18n'
    import he from 'he'
    import JiraImportServerData from './JiraImportServerData.svelte'

    let showJiraRestConfig = false

    export let xfetch
    export let notifications
    export let eventTag
    export let handlePlanAdd = () => {}

    const allowJiraImport = appConfig.AllowJiraImportRest
    const jiraServerUrl = appConfig.JiraServerUrl

    function toggleJiraRestConfig() {
        showJiraRestConfig = !showJiraRestConfig
    }

    function processResponse(response) {
        if (response != null) {
            for (let i = 0; i < response.issues.length; i++) {
                const item = response.issues[i]
                // decode description and acceptance criteria
                const decodedDescription = he.decode(
                    item.fields.description,
                )
                const acceptanceCriteria = he.decode(
                    item.fields.acceptanceCriteria,
                )
                const ticketLink = jiraServerUrl + "/browse/" + item.key
                const plan = {
                    id: '',
                    planName: item.fields.summary,
                    type: item.fields.issuetype.name.toLowerCase(),
                    referenceId: item.key,
                    link: ticketLink,
                    description: decodedDescription,
                    acceptanceCriteria,
                }
                handlePlanAdd(plan)
            }
            eventTag(
                'jira_import_success',
                'battle',
                `total tickets imported: ${response.Total}`,
            )
        }
    }

    function handleImportFromRest(restCfg) {
        const body = {
            userName: restCfg.userName,
            password: restCfg.password,
            endpoint: restCfg.apiEndpoint,
            jql: restCfg.jql,
        }

        xfetch('/api/jira/tickets', { body })
            .then(res => res.json())
            .then(function (response) {
                processResponse(response)
            })
            .catch(function (error){
                notifications.danger(
                    $_('pages.jiraRestCfg.sendRequest.error')
                )
                eventTag('jira_rest_request', 'battle', 'failure')
            })


        // hide the dialogue
        toggleJiraRestConfig();
    }
</script>

{#if allowJiraImport}
    <HollowButton
            color="blue"
            onClick="{toggleJiraRestConfig}"
            additionalClasses="mr-2">
        {$_('actions.plan.importJiraRest.button')}
    </HollowButton>
{/if}
{#if showJiraRestConfig}
    <JiraImportServerData
            {handleImportFromRest}
            {toggleJiraRestConfig}
    />
{/if}