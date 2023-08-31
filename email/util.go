package email

import (
	"fmt"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
)

func formatRetroItemForMarkdownList(item string) string {
	return fmt.Sprintf("- %s\n", item)
}

func formatRetroActionWithAssignee(action *thunderdome.RetroAction) string {
	var actionItem = action.Content
	var assignees = ""

	for _, assignee := range action.Assignees {
		assignees += fmt.Sprintf("[%s]", assignee.Name)
	}

	if len(assignees) > 0 {
		actionItem = assignees + " " + actionItem
	}

	return formatRetroItemForMarkdownList(actionItem)
}
