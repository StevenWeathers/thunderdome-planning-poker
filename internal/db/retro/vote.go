package retro

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"go.uber.org/zap"
)

type VoteLimitExceededError struct{}

func (e *VoteLimitExceededError) Error() string {
	return "Vote limit exceeded for this retro"
}

type VoteAlreadyCastError struct{}

func (e *VoteAlreadyCastError) Error() string {
	return "Vote already cast for this group and cumulative voting is not allowed"
}

type UnauthorizedUserError struct{}

func (e *UnauthorizedUserError) Error() string {
	return "User is not authorized to vote in this retro"
}

type VotePhaseNotActiveError struct{}

func (e *VotePhaseNotActiveError) Error() string {
	return "Vote phase has ended or has not started yet for this retro"
}

// Helper function to check if a user is authorized for a retro
func isUserAuthorizedForRetro(tx *sql.Tx, RetroID, UserID string) (bool, error) {
	var active bool
	err := tx.QueryRow(`
		SELECT COALESCE((SELECT coalesce(active, FALSE)
		FROM thunderdome.retro_user
		WHERE user_id = $2 AND retro_id = $1), FALSE);`,
		RetroID,
		UserID,
	).Scan(
		&active,
	)
	if err != nil {
		return false, fmt.Errorf("get retro user active status query error: %v", err)
	}

	return active, nil
}

// GroupUserVote inserts a user vote for the retro item group
func (d *Service) GroupUserVote(RetroID string, GroupID string, UserID string) ([]*thunderdome.RetroVote, error) {
	var allowCumulativeVoting bool
	var phase string
	var maxVotes int
	var totalVoteCount int
	var groupVoteCount int
	type voteGroup struct {
		RetroID   string `json:"retro_id"`
		GroupID   string `json:"group_id"`
		UserID    string `json:"user_id"`
		VoteCount int    `json:"vote_count"`
	}
	var votesString string
	voteGroups := make([]voteGroup, 0)

	// Start a transaction
	tx, err := d.DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback if not committed

	// First, check if the user is authorized to vote in this retro
	authorized, err := isUserAuthorizedForRetro(tx, RetroID, UserID)
	if err != nil {
		d.Logger.Error("error checking user authorization", zap.Error(err))
		return nil, fmt.Errorf("error checking user authorization: %w", err)
	}
	if !authorized {
		return nil, &UnauthorizedUserError{}
	}

	err = tx.QueryRow(
		`SELECT r.max_votes, r.allow_cumulative_voting, r.phase,
				COALESCE((SELECT jsonb_agg(rgv)
				FROM thunderdome.retro_group_vote rgv
				WHERE rgv.retro_id = $1 AND rgv.user_id = $2 LIMIT 1), '[]'::jsonb) as votes
				FROM thunderdome.retro r
				WHERE r.id = $1;`,
		RetroID, UserID,
	).Scan(&maxVotes, &allowCumulativeVoting, &phase, &votesString)
	if err != nil {
		d.Logger.Error("retro vote query error", zap.Error(err))
		return nil, err
	}

	if phase != "vote" {
		return nil, &VotePhaseNotActiveError{}
	}

	err = json.Unmarshal([]byte(votesString), &voteGroups)
	if err != nil {
		d.Logger.Error("retro vote json error", zap.Error(err))
		return nil, err
	}

	for _, vg := range voteGroups {
		totalVoteCount += vg.VoteCount
		if vg.GroupID == GroupID {
			groupVoteCount = vg.VoteCount
		}
	}

	if totalVoteCount >= maxVotes {
		return nil, &VoteLimitExceededError{}
	}

	if groupVoteCount > 0 && !allowCumulativeVoting {
		return nil, &VoteAlreadyCastError{}
	}

	result, err := tx.Exec(
		`INSERT INTO thunderdome.retro_group_vote (retro_id, group_id, user_id, vote_count)
				VALUES ($1, $2, $3, 1)
				ON CONFLICT (retro_id, group_id, user_id)
				DO UPDATE SET vote_count = thunderdome.retro_group_vote.vote_count + 1
				WHERE thunderdome.retro_group_vote.retro_id = $1
				  AND thunderdome.retro_group_vote.group_id = $2
				  AND thunderdome.retro_group_vote.user_id = $3
				  AND EXISTS (
					SELECT 1
					FROM thunderdome.retro
					WHERE id = $1 AND allow_cumulative_voting = true
  			);`,
		RetroID, GroupID, UserID,
	)
	if err != nil {
		d.Logger.Error("retro group vote query error", zap.Error(err))
		return nil, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		d.Logger.Error("retro group vote rows affected error", zap.Error(err))
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		d.Logger.Error("retro group vote rows affected zero error", zap.Error(err))
		return nil, fmt.Errorf("no vote found to retract")
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	votes := d.GetRetroVotes(RetroID)

	return votes, nil
}

// GroupUserSubtractVote retracts a single user vote for the retro item group
func (d *Service) GroupUserSubtractVote(RetroID string, GroupID string, UserID string) ([]*thunderdome.RetroVote, error) {
	// Start a transaction
	tx, err := d.DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback if not committed

	// Check if cumulative voting is enabled and get the current vote count
	var allowCumulativeVoting bool
	var currentVoteCount int
	err = tx.QueryRow(`
        SELECT r.allow_cumulative_voting, COALESCE(rgv.vote_count, 0)
        FROM thunderdome.retro r
        LEFT JOIN thunderdome.retro_group_vote rgv ON r.id = rgv.retro_id AND rgv.group_id = $2 AND rgv.user_id = $3
        WHERE r.id = $1
    `, RetroID, GroupID, UserID).Scan(&allowCumulativeVoting, &currentVoteCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get retro and vote information: %w", err)
	}

	if currentVoteCount == 0 {
		return nil, fmt.Errorf("no vote found to retract")
	}

	var result sql.Result
	if allowCumulativeVoting && currentVoteCount > 1 {
		// Decrement the vote count if cumulative voting is enabled and there's more than one vote
		result, err = tx.Exec(`
            UPDATE thunderdome.retro_group_vote
            SET vote_count = vote_count - 1
            WHERE retro_id = $1 AND group_id = $2 AND user_id = $3;
        `, RetroID, GroupID, UserID)
	} else {
		// Delete the vote if it's the last one or cumulative voting is not enabled
		result, err = tx.Exec(`
            DELETE FROM thunderdome.retro_group_vote
            WHERE retro_id = $1 AND group_id = $2 AND user_id = $3;
        `, RetroID, GroupID, UserID)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to update vote: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("no vote found to retract")
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	votes := d.GetRetroVotes(RetroID)
	return votes, nil
}

// GetRetroVotes gets retro votes
func (d *Service) GetRetroVotes(RetroID string) []*thunderdome.RetroVote {
	var votes = make([]*thunderdome.RetroVote, 0)

	itemRows, itemsErr := d.DB.Query(
		`SELECT group_id, user_id, vote_count FROM thunderdome.retro_group_vote WHERE retro_id = $1;`,
		RetroID,
	)
	if itemsErr == nil {
		defer itemRows.Close()
		for itemRows.Next() {
			var rv = &thunderdome.RetroVote{}
			if err := itemRows.Scan(&rv.GroupID, &rv.UserID, &rv.Count); err != nil {
				d.Logger.Error("get retro votes query scan error", zap.Error(err))
			} else {
				votes = append(votes, rv)
			}
		}
	} else {
		d.Logger.Error("get retro votes query error", zap.Error(itemsErr))
	}

	return votes
}
