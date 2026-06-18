package team

import (
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
)

func TestIsCheckinAlreadyExistsError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "matching unique violation",
			err: &pgconn.PgError{
				Code:           "23505",
				ConstraintName: "team_checkin_team_id_user_id_checkin_date_uidx",
			},
			want: true,
		},
		{
			name: "wrapped matching unique violation",
			err: fmt.Errorf("wrapped: %w", &pgconn.PgError{
				Code:           "23505",
				ConstraintName: "team_checkin_team_id_user_id_checkin_date_uidx",
			}),
			want: true,
		},
		{
			name: "different unique constraint",
			err: &pgconn.PgError{
				Code:           "23505",
				ConstraintName: "other_constraint",
			},
			want: false,
		},
		{
			name: "different postgres error",
			err: &pgconn.PgError{
				Code:           "23503",
				ConstraintName: "team_checkin_team_id_user_id_checkin_date_uidx",
			},
			want: false,
		},
		{
			name: "non postgres error",
			err:  fmt.Errorf("boom"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isCheckinAlreadyExistsError(tt.err); got != tt.want {
				t.Fatalf("isCheckinAlreadyExistsError() = %v, want %v", got, tt.want)
			}
		})
	}
}
