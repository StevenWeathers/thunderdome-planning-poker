package poker

import (
	"testing"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/sanitize"
	"github.com/stretchr/testify/assert"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
)

func TestSanitizeStory(t *testing.T) {
	d := &Service{HTMLSanitizerPolicy: sanitize.RichTextPolicy()}
	story := &thunderdome.Story{
		Description:        `<p class="ql-align-center">keep</p><script>alert(1)</script>`,
		AcceptanceCriteria: `<img src="x" onerror="alert(1)"><strong>keep</strong>`,
	}

	d.sanitizeStory(story)

	assert.Equal(t, `<p class="ql-align-center">keep</p>`, story.Description)
	assert.Equal(t, `<img src="x"><strong>keep</strong>`, story.AcceptanceCriteria)
}
