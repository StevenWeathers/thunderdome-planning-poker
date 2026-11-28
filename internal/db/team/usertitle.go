package team

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/jackc/pgx/v5/pgconn"
)

// UserTitleList gets the titles offered by the organization that owns the given scope.
// A teamID takes precedence, then a departmentID (which must belong to organizationID), then organizationID.
// Teams that don't belong to an organization have no titles.
func (d *OrganizationService) UserTitleList(ctx context.Context, organizationID string, departmentID string, teamID string) ([]*thunderdome.UserTitle, error) {
	titles := make([]*thunderdome.UserTitle, 0)

	rows, err := d.DB.QueryContext(ctx,
		`SELECT ut.id, ut.organization_id, ut.name, ut.created_at, ut.updated_at
		FROM thunderdome.organization_user_title ut
		WHERE ut.organization_id = CASE
			WHEN $3::text <> '' THEN (
				SELECT COALESCE(t.organization_id, od.organization_id)
				FROM thunderdome.team t
				LEFT JOIN thunderdome.organization_department od ON od.id = t.department_id
				WHERE t.id = NULLIF($3::text, '')::uuid
			)
			WHEN $2::text <> '' THEN (
				SELECT od.organization_id
				FROM thunderdome.organization_department od
				WHERE od.id = NULLIF($2::text, '')::uuid AND od.organization_id = NULLIF($1::text, '')::uuid
			)
			ELSE NULLIF($1::text, '')::uuid
		END
		ORDER BY lower(ut.name);`,
		organizationID, departmentID, teamID,
	)
	if err != nil {
		return nil, fmt.Errorf("user title list query error: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		title := &thunderdome.UserTitle{}
		if err := rows.Scan(&title.ID, &title.OrganizationID, &title.Name, &title.CreatedAt, &title.UpdatedAt); err != nil {
			return nil, fmt.Errorf("user title list scan error: %v", err)
		}
		titles = append(titles, title)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("user title list rows error: %v", err)
	}

	return titles, nil
}

// UserTitleCreate adds a title to an organization
func (d *OrganizationService) UserTitleCreate(ctx context.Context, title *thunderdome.UserTitle) error {
	err := d.DB.QueryRowContext(ctx,
		`INSERT INTO thunderdome.organization_user_title (organization_id, name)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at;`,
		title.OrganizationID, title.Name,
	).Scan(&title.ID, &title.CreatedAt, &title.UpdatedAt)
	if err != nil {
		return userTitleWriteError("user title create", err)
	}

	return nil
}

// UserTitleUpdate renames an organization title, titles already assigned to members are not changed
func (d *OrganizationService) UserTitleUpdate(ctx context.Context, title *thunderdome.UserTitle) error {
	err := d.DB.QueryRowContext(ctx,
		`UPDATE thunderdome.organization_user_title
		SET name = $3, updated_at = CURRENT_TIMESTAMP
		WHERE organization_id = $1 AND id = $2
		RETURNING created_at, updated_at;`,
		title.OrganizationID, title.ID, title.Name,
	).Scan(&title.CreatedAt, &title.UpdatedAt)
	if err != nil {
		return userTitleWriteError("user title update", err)
	}

	return nil
}

// UserTitleDelete removes a title from an organization, titles already assigned to members are not changed
func (d *OrganizationService) UserTitleDelete(ctx context.Context, organizationID string, titleID string) error {
	var id string
	err := d.DB.QueryRowContext(ctx,
		`DELETE FROM thunderdome.organization_user_title WHERE organization_id = $1 AND id = $2 RETURNING id;`,
		organizationID, titleID,
	).Scan(&id)
	if err != nil {
		return userTitleWriteError("user title delete", err)
	}

	return nil
}

func userTitleWriteError(action string, err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return thunderdome.ErrUserTitleNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return thunderdome.ErrDuplicateUserTitle
	}

	return fmt.Errorf("%s query error: %v", action, err)
}
