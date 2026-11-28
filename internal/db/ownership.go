package db

import "strings"

// TeamAdminCondition is SQL matching when the user is a site admin, or an admin of the team or of its
// department or organization. teamID and userID are SQL expressions such as "c.team_id" and "$3".
func TeamAdminCondition(teamID string, userID string) string {
	return strings.NewReplacer("{team}", teamID, "{user}", userID).Replace(`(
		EXISTS (SELECT 1 FROM thunderdome.users au WHERE au.id = {user} AND au.type = 'ADMIN')
		OR EXISTS (SELECT 1 FROM thunderdome.team_user atu
			WHERE atu.team_id = {team} AND atu.user_id = {user} AND atu.role = 'ADMIN')
		OR EXISTS (SELECT 1 FROM thunderdome.team at
			JOIN thunderdome.department_user adu ON adu.department_id = at.department_id
			WHERE at.id = {team} AND adu.user_id = {user} AND adu.role = 'ADMIN')
		OR EXISTS (SELECT 1 FROM thunderdome.team at
			LEFT JOIN thunderdome.organization_department aod ON aod.id = at.department_id
			JOIN thunderdome.organization_user aou
				ON aou.organization_id = COALESCE(at.organization_id, aod.organization_id)
			WHERE at.id = {team} AND aou.user_id = {user} AND aou.role = 'ADMIN')
	)`)
}

// AuthorOrTeamAdminCondition is SQL matching when the user wrote the record or is a TeamAdminCondition admin
func AuthorOrTeamAdminCondition(authorID string, teamID string, userID string) string {
	return "(" + authorID + " = " + userID + " OR " + TeamAdminCondition(teamID, userID) + ")"
}
