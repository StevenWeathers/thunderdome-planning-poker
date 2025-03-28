package email

import (
	"fmt"
	"unicode"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
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

// removeAccents removes accents from a string
func removeAccents(s string) (string, error) {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, e := transform.String(t, s)
	if e != nil {
		return "", e
	}
	return output, nil
}
