<script>
    import he from 'he'
    import TurndownService from 'turndown'

    import HollowButton from '../components/HollowButton.svelte'
    import { _ } from '../i18n'

    export let notifications
    export let handlePlanAdd = () => {}

    const turndownService = new TurndownService()

    let files

    function showDialog() {
        const fileInput = document.querySelector('[data-jira-import]')
        if (fileInput) {
            fileInput.click()
        }
    }

    function uploadFile() {
        let file = this.files[0]
        if (!file) {
            return
        }
        if (file.type !== 'text/xml') {
            notifications.danger($_('actions.plan.importJiraXML.badFileType'))
            return
        }

        let reader = new FileReader()

        reader.readAsText(file)

        reader.onload = () => {
            try {
                const docParser = new DOMParser()
                const domContent = reader.result.replace(/<!--.*?-->/gis, '')
                const doc = docParser.parseFromString(
                    domContent,
                    'application/xml',
                )
                const items = doc.querySelectorAll('channel>item')
                if (items) {
                    for (let i = 0; i < items.length; i++) {
                        const item = items[i]
                        const decodedDescription = he.decode(
                            item.querySelector('description').innerHTML,
                        )
                        const markdownDescription = turndownService.turndown(
                            decodedDescription,
                        )
                        const plan = {
                            id: '',
                            planName: item.querySelector('summary').innerHTML,
                            type: item
                                .querySelector('type')
                                .innerHTML.toLowerCase(),
                            referenceId: item.querySelector('key').innerHTML,
                            link: item.querySelector('link').innerHTML,
                            description: markdownDescription,
                            acceptanceCriteria: '',
                        }
                        handlePlanAdd(plan)
                    }
                }
            } catch (e) {
                notifications.danger(
                    $_('actions.plan.importJiraXML.errorReadingFile'),
                )
            }
        }

        reader.onerror = () => {
            notifications.danger(
                $_('actions.plan.importJiraXML.errorReadingFile'),
            )
        }
    }
</script>

<input type="file" on:change="{uploadFile}" class="hidden" />
<HollowButton onClick="{showDialog}">
    {$_('actions.plan.importJiraXML.button')}
</HollowButton>
