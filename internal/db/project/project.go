package project

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"

	"go.uber.org/zap"
)

// Service represents a PostgreSQL implementation of thunderdome.ProjectDataSvc.
type Service struct {
	DB     *sql.DB
	Logger *otelzap.Logger
}

// GetProjectsByOrganization retrieves all projects for a specific organization
func (s *Service) GetProjectsByOrganization(ctx context.Context, organizationID string) ([]*thunderdome.Project, error) {
	projects := make([]*thunderdome.Project, 0)

	rows, err := s.DB.QueryContext(ctx,
		`SELECT id, project_key, name, description, organization_id, department_id, team_id, created_at, updated_at
		FROM thunderdome.project
		WHERE organization_id = $1
		ORDER BY created_at DESC;`,
		organizationID,
	)

	if err != nil {
		s.Logger.Ctx(ctx).Error("GetProjectsByOrganization query error", zap.Error(err))
		return nil, fmt.Errorf("error querying projects for organization: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p thunderdome.Project
		if err := rows.Scan(
			&p.ID,
			&p.ProjectKey,
			&p.Name,
			&p.Description,
			&p.OrganizationID,
			&p.DepartmentID,
			&p.TeamID,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			s.Logger.Ctx(ctx).Error("GetProjectsByOrganization row scan error", zap.Error(err))
		} else {
			projects = append(projects, &p)
		}
	}

	return projects, nil
}

// GetProjectsByDepartment retrieves all projects for a specific department
func (s *Service) GetProjectsByDepartment(ctx context.Context, departmentID string) ([]*thunderdome.Project, error) {
	projects := make([]*thunderdome.Project, 0)

	rows, err := s.DB.QueryContext(ctx,
		`SELECT id, project_key, name, description, organization_id, department_id, team_id, created_at, updated_at
		FROM thunderdome.project
		WHERE department_id = $1
		ORDER BY created_at DESC;`,
		departmentID,
	)

	if err != nil {
		s.Logger.Ctx(ctx).Error("GetProjectsByDepartment query error", zap.Error(err))
		return nil, fmt.Errorf("error querying projects for department: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p thunderdome.Project
		if err := rows.Scan(
			&p.ID,
			&p.ProjectKey,
			&p.Name,
			&p.Description,
			&p.OrganizationID,
			&p.DepartmentID,
			&p.TeamID,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			s.Logger.Ctx(ctx).Error("GetProjectsByDepartment row scan error", zap.Error(err))
		} else {
			projects = append(projects, &p)
		}
	}

	return projects, nil
}

// GetProjectsByTeam retrieves all projects for a specific team
func (s *Service) GetProjectsByTeam(ctx context.Context, teamID string) ([]*thunderdome.Project, error) {
	projects := make([]*thunderdome.Project, 0)

	rows, err := s.DB.QueryContext(ctx,
		`SELECT id, project_key, name, description, organization_id, department_id, team_id, created_at, updated_at
		FROM thunderdome.project
		WHERE team_id = $1
		ORDER BY created_at DESC;`,
		teamID,
	)

	if err != nil {
		return nil, fmt.Errorf("error querying projects for team: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p thunderdome.Project
		if err := rows.Scan(
			&p.ID,
			&p.ProjectKey,
			&p.Name,
			&p.Description,
			&p.OrganizationID,
			&p.DepartmentID,
			&p.TeamID,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			s.Logger.Ctx(ctx).Error("GetProjectsByTeam row scan error", zap.Error(err))
		} else {
			projects = append(projects, &p)
		}
	}

	return projects, nil
}

// GetProjectByID retrieves a specific project by its ID
func (s *Service) GetProjectByID(ctx context.Context, projectID string) (*thunderdome.Project, error) {
	var p thunderdome.Project

	err := s.DB.QueryRowContext(ctx,
		`SELECT id, project_key, name, description, COALESCE(organization_id::text, ''), COALESCE(department_id::text, ''), COALESCE(team_id::text, ''), created_at, updated_at
		FROM thunderdome.project
		WHERE id = $1;`,
		projectID,
	).Scan(
		&p.ID,
		&p.ProjectKey,
		&p.Name,
		&p.Description,
		&p.OrganizationID,
		&p.DepartmentID,
		&p.TeamID,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error querying project by ID: %v", err)
	}

	return &p, nil
}

// CreateProject creates a new project
func (s *Service) CreateProject(ctx context.Context, project *thunderdome.Project) error {
	_, err := s.DB.ExecContext(ctx,
		`INSERT INTO thunderdome.project (
			project_key, name, description, organization_id, department_id, team_id)
		VALUES ($1, $2, $3, NULLIF($4, '')::uuid, NULLIF($5, '')::uuid, NULLIF($6, '')::uuid);`,
		project.ProjectKey,
		project.Name,
		project.Description,
		project.OrganizationID,
		project.DepartmentID,
		project.TeamID,
	)

	if err != nil {
		return fmt.Errorf("error creating new project: %v", err)
	}

	return nil
}

// UpdateProject updates an existing project
func (s *Service) UpdateProject(ctx context.Context, project *thunderdome.Project) error {
	_, err := s.DB.ExecContext(ctx,
		`UPDATE thunderdome.project
		SET project_key = $2, name = $3, description = $4, organization_id = NULLIF($5, '')::uuid, department_id = NULLIF($6, '')::uuid, team_id = NULLIF($7, '')::uuid, updated_at = NOW()
		WHERE id = $1;`,
		project.ID,
		project.ProjectKey,
		project.Name,
		project.Description,
		project.OrganizationID,
		project.DepartmentID,
		project.TeamID,
	)

	if err != nil {
		return fmt.Errorf("error updating project: %v", err)
	}

	return nil
}

// DeleteProject deletes a project by its ID
func (s *Service) DeleteProject(ctx context.Context, projectID string) error {
	_, err := s.DB.ExecContext(ctx,
		`DELETE FROM thunderdome.project WHERE id = $1;`,
		projectID,
	)

	if err != nil {
		return fmt.Errorf("error deleting project: %v", err)
	}

	return nil
}

// ListProjects retrieves a paginated list of projects
func (s *Service) ListProjects(ctx context.Context, limit int, offset int) ([]*thunderdome.Project, int, error) {
	projects := make([]*thunderdome.Project, 0)
	var totalCount int

	err := s.DB.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM thunderdome.project;",
	).Scan(&totalCount)

	if err != nil {
		s.Logger.Ctx(ctx).Error("ListProjects count query error", zap.Error(err))
		return nil, 0, fmt.Errorf("error counting projects: %v", err)
	}

	rows, err := s.DB.QueryContext(ctx,
		`SELECT id, project_key, name, description, COALESCE(organization_id::text, ''), COALESCE(department_id::text, ''), COALESCE(team_id::text, ''), created_at, updated_at
		FROM thunderdome.project
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2;`,
		limit,
		offset,
	)

	if err != nil {
		return nil, 0, fmt.Errorf("error querying projects: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p thunderdome.Project
		if err := rows.Scan(
			&p.ID,
			&p.ProjectKey,
			&p.Name,
			&p.Description,
			&p.OrganizationID,
			&p.DepartmentID,
			&p.TeamID,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			s.Logger.Ctx(ctx).Error("ListProjects row scan error", zap.Error(err))
		} else {
			projects = append(projects, &p)
		}
	}

	return projects, totalCount, nil
}

// UpdateOrganizationProject updates an existing organization project
func (s *Service) UpdateOrganizationProject(ctx context.Context, project *thunderdome.Project) error {
	_, err := s.DB.ExecContext(ctx,
		`UPDATE thunderdome.project
		SET project_key = $3, name = $4, description = $5, updated_at = NOW()
		WHERE id = $1 AND organization_id = $2;`,
		project.ID,
		project.OrganizationID,
		project.ProjectKey,
		project.Name,
		project.Description,
	)

	if err != nil {
		return fmt.Errorf("error updating organization project: %v", err)
	}

	return nil
}

// UpdateDepartmentProject updates an existing department project
func (s *Service) UpdateDepartmentProject(ctx context.Context, project *thunderdome.Project) error {
	_, err := s.DB.ExecContext(ctx,
		`UPDATE thunderdome.project
		SET project_key = $3, name = $4, description = $5, updated_at = NOW()
		WHERE id = $1 AND department_id = $2;`,
		project.ID,
		project.DepartmentID,
		project.ProjectKey,
		project.Name,
		project.Description,
	)

	if err != nil {
		return fmt.Errorf("error updating department project: %v", err)
	}

	return nil
}

// UpdateTeamProject updates an existing team project
func (s *Service) UpdateTeamProject(ctx context.Context, project *thunderdome.Project) error {
	_, err := s.DB.ExecContext(ctx,
		`UPDATE thunderdome.project
		SET project_key = $3, name = $4, description = $5, updated_at = NOW()
		WHERE id = $1 AND team_id = $2;`,
		project.ID,
		project.TeamID,
		project.ProjectKey,
		project.Name,
		project.Description,
	)

	if err != nil {
		return fmt.Errorf("error updating team project: %v", err)
	}

	return nil
}

// DeleteOrganizationProject deletes an organization project by its ID
func (s *Service) DeleteOrganizationProject(ctx context.Context, orgID string, projectID string) error {
	_, err := s.DB.ExecContext(ctx,
		`DELETE FROM thunderdome.project WHERE id = $1 AND organization_id = $2;`,
		projectID, orgID,
	)

	if err != nil {
		return fmt.Errorf("error deleting organization project: %v", err)
	}

	return nil
}

// DeleteDepartmentProject deletes a department project by its ID
func (s *Service) DeleteDepartmentProject(ctx context.Context, deptID string, projectID string) error {
	_, err := s.DB.ExecContext(ctx,
		`DELETE FROM thunderdome.project WHERE id = $1 AND department_id = $2;`,
		projectID, deptID,
	)

	if err != nil {
		return fmt.Errorf("error deleting department project: %v", err)
	}

	return nil
}

// DeleteTeamProject deletes a team project by its ID
func (s *Service) DeleteTeamProject(ctx context.Context, teamID string, projectID string) error {
	_, err := s.DB.ExecContext(ctx,
		`DELETE FROM thunderdome.project WHERE id = $1 AND team_id = $2;`,
		projectID, teamID,
	)

	if err != nil {
		return fmt.Errorf("error deleting team project: %v", err)
	}

	return nil
}
