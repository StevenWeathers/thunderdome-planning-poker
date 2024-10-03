package email

import (
	"fmt"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
)

// formatRetroItemForMarkdownList formats a retro item for a markdown list
func formatRetroItemForMarkdownList(item string) string {
	return fmt.Sprintf("- %s\n", item)
}

// formatRetroActionWithAssignee formats a retro action with assignee for a markdown list
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
