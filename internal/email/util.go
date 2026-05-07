package email

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var emailNamePattern = regexp.MustCompile(`[^\s<>"@]+@[^\s<>"]+`)

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

// sanitizeEmailName removes invalid characters from a name for use in email headers
func sanitizeEmailName(name string) (string, error) {
	// Remove accents from the name
	s, err := removeAccents(name)
	if err != nil {
		return name, err
	}

	s = strings.TrimSpace(s)

	if start := strings.Index(s, "<"); start >= 0 {
		if end := strings.Index(s[start+1:], ">"); end >= 0 {
			prefix := strings.TrimSpace(s[:start])
			if prefix != "" {
				s = prefix
			} else {
				s = strings.TrimSpace(s[start+1 : start+1+end])
			}
		}
	}

	s = emailNamePattern.ReplaceAllStringFunc(s, func(match string) string {
		parts := strings.SplitN(match, "@", 2)
		return parts[0]
	})

	// Remove any email address characters from the name
	s = strings.Map(func(r rune) rune {
		switch {
		case unicode.IsSpace(r):
			return ' '
		case unicode.IsControl(r):
			return -1
		case unicode.IsLetter(r), unicode.IsNumber(r):
			return r
		case r == '.', r == '-', r == '\'', r == '_':
			return r
		default:
			return -1
		}
	}, s)

	return strings.Join(strings.Fields(s), " "), nil
}
