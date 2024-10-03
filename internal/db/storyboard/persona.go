package storyboard

import (
	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"go.uber.org/zap"
)

// GetStoryboardPersonas retrieves the personas for a given storyboard from db
func (d *Service) GetStoryboardPersonas(storyboardID string) []*thunderdome.StoryboardPersona {
	var personas = make([]*thunderdome.StoryboardPersona, 0)
	rows, err := d.DB.Query(
		`SELECT
			p.id, p.name, p.role, p.description
		FROM thunderdome.storyboard_persona p
		WHERE p.storyboard_id = $1;`,
		storyboardID,
	)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var p thunderdome.StoryboardPersona
			if err := rows.Scan(&p.ID, &p.Name, &p.Role, &p.Description); err != nil {
				d.Logger.Error("get_storyboard_personas query scan error", zap.Error(err))
			} else {
				personas = append(personas, &p)
			}
		}
	}

	return personas
}

// AddStoryboardPersona adds a persona to a storyboard
func (d *Service) AddStoryboardPersona(storyboardID string, userID string, name string, role string, description string) ([]*thunderdome.StoryboardPersona, error) {
	if _, err := d.DB.Exec(
		`INSERT INTO thunderdome.storyboard_persona (storyboard_id, name, role, description) VALUES ($1, $2, $3, $4);`,
		storyboardID,
		name,
		role,
		description,
	); err != nil {
		d.Logger.Error("CALL thunderdome.persona_add error", zap.Error(err))
	}

	personas := d.GetStoryboardPersonas(storyboardID)

	return personas, nil
}

// UpdateStoryboardPersona updates a storyboard persona
func (d *Service) UpdateStoryboardPersona(storyboardID string, userID string, personaID string, name string, role string, description string) ([]*thunderdome.StoryboardPersona, error) {
	if _, err := d.DB.Exec(
		`UPDATE thunderdome.storyboard_persona SET name = $2, role = $3, description = $4, updated_date = NOW() WHERE id = $1;`,
		personaID,
		name,
		role,
		description,
	); err != nil {
		d.Logger.Error("CALL thunderdome.persona_edit error", zap.Error(err))
	}

	personas := d.GetStoryboardPersonas(storyboardID)

	return personas, nil
}

// DeleteStoryboardPersona deletes a storyboard persona
func (d *Service) DeleteStoryboardPersona(storyboardID string, userID string, personaID string) ([]*thunderdome.StoryboardPersona, error) {
	if _, err := d.DB.Exec(
		`DELETE FROM thunderdome.storyboard_persona WHERE id = $1;`,
		personaID,
	); err != nil {
		d.Logger.Error("CALL thunderdome.persona_delete error", zap.Error(err))
	}

	personas := d.GetStoryboardPersonas(storyboardID)

	return personas, nil
}
