package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"go.uber.org/zap"
)

type retroTagRequestBody struct {
	Name string `json:"name" validate:"required"`
}

// retroTagScope reads the owning organization or team of a preset tag from the route,
// team routes nested under an organization belong to the team
func (s *Service) retroTagScope(r *http.Request, pathParam string) (organizationID string, teamID string, err error) {
	if pathParam == "" {
		return "", "", nil
	}

	id := r.PathValue(pathParam)
	if err := validate.Var(id, "required,uuid"); err != nil {
		return "", "", err
	}

	if pathParam == "teamId" {
		return "", id, nil
	}

	return id, "", nil
}

func (s *Service) retroTagFailure(w http.ResponseWriter, r *http.Request, action string, err error) {
	switch {
	case errors.Is(err, thunderdome.ErrDuplicateRetroTag):
		s.Failure(w, r, http.StatusConflict, Errorf(ECONFLICT, err.Error()))
	case errors.Is(err, thunderdome.ErrRetroTagNotFound):
		s.Failure(w, r, http.StatusNotFound, Errorf(ENOTFOUND, err.Error()))
	default:
		s.Logger.Ctx(r.Context()).Error(action+" error", zap.Error(err),
			zap.String("organization_id", r.PathValue("orgId")),
			zap.String("team_id", r.PathValue("teamId")),
			zap.String("tag_id", r.PathValue("tagId")))
		s.Failure(w, r, http.StatusInternalServerError, err)
	}
}

func (s *Service) listRetroTags(pathParam string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		organizationID, teamID, err := s.retroTagScope(r, pathParam)
		if err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}

		tags, err := s.RetroTagDataSvc.GetRetroTags(r.Context(), organizationID, teamID)
		if err != nil {
			s.retroTagFailure(w, r, "listRetroTags", err)
			return
		}

		s.Success(w, r, http.StatusOK, tags, nil)
	}
}

// saveRetroTag creates a preset tag, or updates it when the route includes a tag ID
func (s *Service) saveRetroTag(pathParam string, update bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		organizationID, teamID, err := s.retroTagScope(r, pathParam)
		if err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}

		tag := &thunderdome.RetroTag{}
		if organizationID != "" {
			tag.OrganizationID = &organizationID
		}
		if teamID != "" {
			tag.TeamID = &teamID
		}

		if update {
			tag.ID = r.PathValue("tagId")
			if err := validate.Var(tag.ID, "required,uuid"); err != nil {
				s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
				return
			}
		}

		request := retroTagRequestBody{}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}
		if err := json.Unmarshal(body, &request); err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}
		if err := validate.Struct(request); err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}

		tag.Name, err = thunderdome.NormalizeRetroTag(request.Name)
		if err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}

		if update {
			err = s.RetroTagDataSvc.UpdateRetroTag(ctx, tag)
		} else {
			tag.CreatedBy = sessionUserID
			err = s.RetroTagDataSvc.CreateRetroTag(ctx, tag)
		}
		if err != nil {
			s.retroTagFailure(w, r, "saveRetroTag", err)
			return
		}

		s.Success(w, r, http.StatusOK, tag, nil)
	}
}

func (s *Service) deleteRetroTag(pathParam string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		organizationID, teamID, err := s.retroTagScope(r, pathParam)
		if err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}

		tagID := r.PathValue("tagId")
		if err := validate.Var(tagID, "required,uuid"); err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}

		if err := s.RetroTagDataSvc.DeleteRetroTag(r.Context(), organizationID, teamID, tagID); err != nil {
			s.retroTagFailure(w, r, "deleteRetroTag", err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleGetAdminRetroTags gets the global default retro tags
//
//	@Summary		Get Global Retro Tags
//	@Description	get list of global default retro tags
//	@Tags			retroTag
//	@Produce		json
//	@Success		200	object	standardJsonResponse{data=[]thunderdome.RetroTag}
//	@Failure		500	object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/admin/retro-tags [get]
func (s *Service) handleGetAdminRetroTags() http.HandlerFunc {
	return s.listRetroTags("")
}

// handleAdminRetroTagCreate creates a global default retro tag
//
//	@Summary		Create Global Retro Tag
//	@Description	Creates a global default retro tag
//	@Tags			retroTag
//	@Produce		json
//	@Param			tag	body	retroTagRequestBody	true	"new retro tag object"
//	@Success		200	object	standardJsonResponse{data=thunderdome.RetroTag}
//	@Failure		400	object	standardJsonResponse{}
//	@Failure		409	object	standardJsonResponse{}
//	@Failure		500	object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/admin/retro-tags [post]
func (s *Service) handleAdminRetroTagCreate() http.HandlerFunc {
	return s.saveRetroTag("", false)
}

// handleAdminRetroTagUpdate updates a global default retro tag
//
//	@Summary		Update Global Retro Tag
//	@Description	Updates a global default retro tag
//	@Tags			retroTag
//	@Produce		json
//	@Param			tagId	path	string				true	"the retro tag ID to update"
//	@Param			tag		body	retroTagRequestBody	true	"retro tag object to update"
//	@Success		200		object	standardJsonResponse{data=thunderdome.RetroTag}
//	@Failure		400		object	standardJsonResponse{}
//	@Failure		404		object	standardJsonResponse{}
//	@Failure		409		object	standardJsonResponse{}
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/admin/retro-tags/{tagId} [put]
func (s *Service) handleAdminRetroTagUpdate() http.HandlerFunc {
	return s.saveRetroTag("", true)
}

// handleAdminRetroTagDelete deletes a global default retro tag
//
//	@Summary		Delete Global Retro Tag
//	@Description	Deletes a global default retro tag
//	@Tags			retroTag
//	@Produce		json
//	@Param			tagId	path	string	true	"the retro tag ID to delete"
//	@Success		200		object	standardJsonResponse{}
//	@Failure		400		object	standardJsonResponse{}
//	@Failure		404		object	standardJsonResponse{}
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/admin/retro-tags/{tagId} [delete]
func (s *Service) handleAdminRetroTagDelete() http.HandlerFunc {
	return s.deleteRetroTag("")
}

// handleGetOrganizationRetroTags gets the retro tags of an organization
//
//	@Summary		Get Organization Retro Tags
//	@Description	get list of retro tags for an organization
//	@Tags			retroTag
//	@Produce		json
//	@Param			orgId	path	string	true	"Organization ID"
//	@Success		200		object	standardJsonResponse{data=[]thunderdome.RetroTag}
//	@Failure		400		object	standardJsonResponse{}
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/retro-tags [get]
func (s *Service) handleGetOrganizationRetroTags() http.HandlerFunc {
	return s.listRetroTags("orgId")
}

// handleOrganizationRetroTagCreate creates an organization retro tag
//
//	@Summary		Create Organization Retro Tag
//	@Description	Creates an organization retro tag
//	@Tags			retroTag
//	@Produce		json
//	@Param			orgId	path	string				true	"the organization ID"
//	@Param			tag		body	retroTagRequestBody	true	"new retro tag object"
//	@Success		200		object	standardJsonResponse{data=thunderdome.RetroTag}
//	@Failure		400		object	standardJsonResponse{}
//	@Failure		409		object	standardJsonResponse{}
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/retro-tags [post]
func (s *Service) handleOrganizationRetroTagCreate() http.HandlerFunc {
	return s.saveRetroTag("orgId", false)
}

// handleOrganizationRetroTagUpdate updates an organization retro tag
//
//	@Summary		Update Organization Retro Tag
//	@Description	Updates an organization retro tag
//	@Tags			retroTag
//	@Produce		json
//	@Param			orgId	path	string				true	"the organization ID"
//	@Param			tagId	path	string				true	"the retro tag ID to update"
//	@Param			tag		body	retroTagRequestBody	true	"retro tag object to update"
//	@Success		200		object	standardJsonResponse{data=thunderdome.RetroTag}
//	@Failure		400		object	standardJsonResponse{}
//	@Failure		404		object	standardJsonResponse{}
//	@Failure		409		object	standardJsonResponse{}
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/retro-tags/{tagId} [put]
func (s *Service) handleOrganizationRetroTagUpdate() http.HandlerFunc {
	return s.saveRetroTag("orgId", true)
}

// handleOrganizationRetroTagDelete deletes an organization retro tag
//
//	@Summary		Delete Organization Retro Tag
//	@Description	Deletes an organization retro tag
//	@Tags			retroTag
//	@Produce		json
//	@Param			orgId	path	string	true	"the organization ID"
//	@Param			tagId	path	string	true	"the retro tag ID to delete"
//	@Success		200		object	standardJsonResponse{}
//	@Failure		400		object	standardJsonResponse{}
//	@Failure		404		object	standardJsonResponse{}
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/retro-tags/{tagId} [delete]
func (s *Service) handleOrganizationRetroTagDelete() http.HandlerFunc {
	return s.deleteRetroTag("orgId")
}

// handleGetTeamRetroTags gets the retro tags of a team
//
//	@Summary		Get Team Retro Tags
//	@Description	get list of retro tags for a team
//	@Tags			retroTag
//	@Produce		json
//	@Param			teamId	path	string	true	"Team ID"
//	@Success		200		object	standardJsonResponse{data=[]thunderdome.RetroTag}
//	@Failure		400		object	standardJsonResponse{}
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/retro-tags [get]
//	@Router			/organizations/{orgId}/teams/{teamId}/retro-tags [get]
//	@Router			/organizations/{orgId}/departments/{departmentId}/teams/{teamId}/retro-tags [get]
func (s *Service) handleGetTeamRetroTags() http.HandlerFunc {
	return s.listRetroTags("teamId")
}

// handleTeamRetroTagCreate creates a team retro tag
//
//	@Summary		Create Team Retro Tag
//	@Description	Creates a team retro tag
//	@Tags			retroTag
//	@Produce		json
//	@Param			teamId	path	string				true	"the team ID"
//	@Param			tag		body	retroTagRequestBody	true	"new retro tag object"
//	@Success		200		object	standardJsonResponse{data=thunderdome.RetroTag}
//	@Failure		400		object	standardJsonResponse{}
//	@Failure		409		object	standardJsonResponse{}
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/retro-tags [post]
//	@Router			/organizations/{orgId}/teams/{teamId}/retro-tags [post]
//	@Router			/organizations/{orgId}/departments/{departmentId}/teams/{teamId}/retro-tags [post]
func (s *Service) handleTeamRetroTagCreate() http.HandlerFunc {
	return s.saveRetroTag("teamId", false)
}

// handleTeamRetroTagUpdate updates a team retro tag
//
//	@Summary		Update Team Retro Tag
//	@Description	Updates a team retro tag
//	@Tags			retroTag
//	@Produce		json
//	@Param			teamId	path	string				true	"the team ID"
//	@Param			tagId	path	string				true	"the retro tag ID to update"
//	@Param			tag		body	retroTagRequestBody	true	"retro tag object to update"
//	@Success		200		object	standardJsonResponse{data=thunderdome.RetroTag}
//	@Failure		400		object	standardJsonResponse{}
//	@Failure		404		object	standardJsonResponse{}
//	@Failure		409		object	standardJsonResponse{}
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/retro-tags/{tagId} [put]
//	@Router			/organizations/{orgId}/teams/{teamId}/retro-tags/{tagId} [put]
//	@Router			/organizations/{orgId}/departments/{departmentId}/teams/{teamId}/retro-tags/{tagId} [put]
func (s *Service) handleTeamRetroTagUpdate() http.HandlerFunc {
	return s.saveRetroTag("teamId", true)
}

// handleTeamRetroTagDelete deletes a team retro tag
//
//	@Summary		Delete Team Retro Tag
//	@Description	Deletes a team retro tag
//	@Tags			retroTag
//	@Produce		json
//	@Param			teamId	path	string	true	"the team ID"
//	@Param			tagId	path	string	true	"the retro tag ID to delete"
//	@Success		200		object	standardJsonResponse{}
//	@Failure		400		object	standardJsonResponse{}
//	@Failure		404		object	standardJsonResponse{}
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/retro-tags/{tagId} [delete]
//	@Router			/organizations/{orgId}/teams/{teamId}/retro-tags/{tagId} [delete]
//	@Router			/organizations/{orgId}/departments/{departmentId}/teams/{teamId}/retro-tags/{tagId} [delete]
func (s *Service) handleTeamRetroTagDelete() http.HandlerFunc {
	return s.deleteRetroTag("teamId")
}
