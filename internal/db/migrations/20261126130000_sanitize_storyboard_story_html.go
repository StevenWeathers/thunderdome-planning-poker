package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/sanitize"
)

const sanitizeStoryboardStoryBatchSize = 500

func init() {
	// the original content isn't kept, so there is nothing to restore on the way down
	goose.AddMigrationContext(upSanitizeStoryboardStoryHTML, nil)
}

type storyboardStoryHTML struct {
	id      string
	content string
}

// upSanitizeStoryboardStoryHTML sanitizes story content saved before it was sanitized on write
func upSanitizeStoryboardStoryHTML(ctx context.Context, tx *sql.Tx) error {
	policy := sanitize.RichTextPolicy()
	lastID := "00000000-0000-0000-0000-000000000000"

	for {
		stories, err := getStoryboardStoryHTMLBatch(ctx, tx, lastID)
		if err != nil {
			return err
		}
		if len(stories) == 0 {
			return nil
		}

		for _, story := range stories {
			lastID = story.id
			clean := policy.Sanitize(story.content)
			if clean == story.content {
				continue
			}
			if _, err := tx.ExecContext(ctx,
				`UPDATE thunderdome.storyboard_story SET content = $2 WHERE id = $1;`,
				story.id, clean,
			); err != nil {
				return fmt.Errorf("sanitize storyboard story %s: %w", story.id, err)
			}
		}
	}
}

// getStoryboardStoryHTMLBatch reads the whole batch before returning,
// as the transaction's connection can't run updates while rows are open
func getStoryboardStoryHTMLBatch(ctx context.Context, tx *sql.Tx, afterID string) ([]storyboardStoryHTML, error) {
	rows, err := tx.QueryContext(ctx,
		`SELECT id, content FROM thunderdome.storyboard_story
		WHERE id > $1 AND content IS NOT NULL ORDER BY id LIMIT $2;`,
		afterID, sanitizeStoryboardStoryBatchSize,
	)
	if err != nil {
		return nil, fmt.Errorf("get storyboard stories to sanitize: %w", err)
	}
	defer rows.Close()

	stories := make([]storyboardStoryHTML, 0, sanitizeStoryboardStoryBatchSize)
	for rows.Next() {
		var s storyboardStoryHTML
		if err := rows.Scan(&s.id, &s.content); err != nil {
			return nil, fmt.Errorf("scan storyboard story to sanitize: %w", err)
		}
		stories = append(stories, s)
	}

	return stories, rows.Err()
}
