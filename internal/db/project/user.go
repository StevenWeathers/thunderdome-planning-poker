package project

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

// IsUserProjectMember returns whether the user has access to the project and the inherited role
// Role inheritance precedence: Project role > Team role > Department role > Organization role.
// If the user isn't associated at any scope, returns (false, "").
func (s *Service) IsUserProjectMember(ctx context.Context, userID, projectID string) (bool, string, error) {
	var (
		isMember bool
		role     string
	)

	err := s.DB.QueryRowContext(ctx, `
			SELECT
				(COALESCE(pu.user_id, tu.user_id, du.user_id, ou.user_id) IS NOT NULL) AS is_member,
				COALESCE(pu.role, tu.role, du.role, ou.role, '') AS inherited_role
			FROM thunderdome.project p
			LEFT JOIN thunderdome.project_user pu ON p.id = pu.project_id AND pu.user_id = $1
			LEFT JOIN thunderdome.team_user tu ON p.team_id = tu.team_id AND tu.user_id = $1
			LEFT JOIN thunderdome.department_user du ON p.department_id = du.department_id AND du.user_id = $1
			LEFT JOIN thunderdome.organization_user ou ON p.organization_id = ou.organization_id AND ou.user_id = $1
			WHERE p.id = $2;`,
		userID, projectID).Scan(&isMember, &role)
	if err != nil {
		s.Logger.Ctx(ctx).Error("IsUserProjectMember query error", zap.Error(err))
		return false, "", fmt.Errorf("error checking project membership: %v", err)
	}

	if !isMember {
		return false, "", nil
	}

	return isMember, role, nil
}
