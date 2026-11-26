// Package migrations holds the database migrations that need to run Go code
package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/microcosm-cc/bluemonday"
	"github.com/pressly/goose/v3"

	"github.com/StevenWeathers/thunderdome-planning-poker/internal/sanitize"
)

const sanitizePokerStoryBatchSize = 500

func init() {
	// the original content isn't kept, so there is nothing to restore on the way down
	goose.AddMigrationContext(upSanitizePokerStoryHTML, nil)
}

type pokerStoryHTML struct {
	id                 string
	description        sql.NullString
	acceptanceCriteria sql.NullString
}

// sanitize cleans the HTML fields and reports whether either of them changed
func (s *pokerStoryHTML) sanitize(policy *bluemonday.Policy) bool {
	changed := false
	for _, field := range []*sql.NullString{&s.description, &s.acceptanceCriteria} {
		if clean := policy.Sanitize(field.String); field.Valid && clean != field.String {
			field.String = clean
			changed = true
		}
	}

	return changed
}

// upSanitizePokerStoryHTML sanitizes stories saved before game creation sanitized them
func upSanitizePokerStoryHTML(ctx context.Context, tx *sql.Tx) error {
	policy := sanitize.RichTextPolicy()
	lastID := "00000000-0000-0000-0000-000000000000"

	for {
		stories, err := getPokerStoryHTMLBatch(ctx, tx, lastID)
		if err != nil {
			return err
		}
		if len(stories) == 0 {
			return nil
		}

		for _, story := range stories {
			lastID = story.id
			if !story.sanitize(policy) {
				continue
			}
			if _, err := tx.ExecContext(ctx,
				`UPDATE thunderdome.poker_story SET description = $2, acceptance_criteria = $3 WHERE id = $1;`,
				story.id, story.description, story.acceptanceCriteria,
			); err != nil {
				return fmt.Errorf("sanitize poker story %s: %w", story.id, err)
			}
		}
	}
}

// getPokerStoryHTMLBatch reads the whole batch before returning,
// as the transaction's connection can't run updates while rows are open
func getPokerStoryHTMLBatch(ctx context.Context, tx *sql.Tx, afterID string) ([]pokerStoryHTML, error) {
	rows, err := tx.QueryContext(ctx,
		`SELECT id, description, acceptance_criteria FROM thunderdome.poker_story
		WHERE id > $1 ORDER BY id LIMIT $2;`,
		afterID, sanitizePokerStoryBatchSize,
	)
	if err != nil {
		return nil, fmt.Errorf("get poker stories to sanitize: %w", err)
	}
	defer rows.Close()

	stories := make([]pokerStoryHTML, 0, sanitizePokerStoryBatchSize)
	for rows.Next() {
		var s pokerStoryHTML
		if err := rows.Scan(&s.id, &s.description, &s.acceptanceCriteria); err != nil {
			return nil, fmt.Errorf("scan poker story to sanitize: %w", err)
		}
		stories = append(stories, s)
	}

	return stories, rows.Err()
}
