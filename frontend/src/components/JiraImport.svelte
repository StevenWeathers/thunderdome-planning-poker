<script>
    import he from 'he'
    import TurndownService from 'turndown'

    import HollowButton from '../components/HollowButton.svelte'

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
            notifications.danger('Error bad file type')
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
                notifications.danger('Error reading file')
            }
        }

        reader.onerror = () => {
            notifications.danger('Error reading file')
        }
    }
</script>

<style>
    [data-jira-import] {
        display: none;
    }
</style>

<input type="file" on:change="{uploadFile}" data-jira-import />
<HollowButton onClick="{showDialog}">Import plans from Jira XML</HollowButton>
