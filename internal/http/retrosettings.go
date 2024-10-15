package http

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"

	"go.uber.org/zap"

	"github.com/gorilla/mux"
)

type retroSettingsRequestBody struct {
	MaxVotes              int16   `json:"maxVotes" validate:"gte=1,lte=100"`
	AllowMultipleVotes    bool    `json:"allowMultipleVotes"`
	BrainstormVisibility  string  `json:"brainstormVisibility" validate:"oneof=visible hidden concealed"`
	PhaseTimeLimit        int16   `json:"phaseTimeLimit" validate:"gte=0,lte=59"`
	PhaseAutoAdvance      bool    `json:"phaseAutoAdvance"`
	AllowCumulativeVoting bool    `json:"allowCumulativeVoting"`
	TemplateID            *string `json:"templateId" validate:"omitempty,uuid"`
	JoinCode              string  `json:"joinCode"`
	FacilitatorCode       string  `json:"facilitatorCode"`
}

// handleCreateOrganizationRetroSettings creates new retro settings for an organization
//
//	@Summary		Create Organization Retro Settings
//	@Description	Creates new retro settings for an organization
//	@Tags			retro-settings
//	@Produce		json
//	@Param			orgId		path	string													true	"Organization ID"
//	@Param			settings	body	retroSettingsRequestBody								true	"retro settings object to create"
//	@Success		201			object	standardJsonResponse{data=thunderdome.RetroSettings}	"returns created retro settings"
//	@Failure		400			object	standardJsonResponse{}
//	@Failure		500			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/retro-settings [post]
func (s *Service) handleCreateOrganizationRetroSettings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		orgID := vars["orgId"]
		orgIDErr := validate.Var(orgID, "required,uuid")
		if orgIDErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, orgIDErr.Error()))
			return
		}

		var settingsReq retroSettingsRequestBody
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		jsonErr := json.Unmarshal(body, &settingsReq)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(settingsReq)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		settings := &thunderdome.RetroSettings{
			OrganizationID:        &orgID,
			MaxVotes:              settingsReq.MaxVotes,
			AllowMultipleVotes:    settingsReq.AllowMultipleVotes,
			BrainstormVisibility:  settingsReq.BrainstormVisibility,
			PhaseTimeLimit:        settingsReq.PhaseTimeLimit,
			PhaseAutoAdvance:      settingsReq.PhaseAutoAdvance,
			AllowCumulativeVoting: settingsReq.AllowCumulativeVoting,
			TemplateID:            settingsReq.TemplateID,
			JoinCode:              settingsReq.JoinCode,
			FacilitatorCode:       settingsReq.FacilitatorCode,
		}

		createdSettings, err := s.RetroDataSvc.CreateSettings(ctx, settings)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleCreateOrganizationRetroSettings error", zap.Error(err),
				zap.String("org_id", orgID), zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusCreated, createdSettings, nil)
	}
}

// handleCreateTeamRetroSettings creates new retro settings for a team
//
//	@Summary		Create Team Retro Settings
//	@Description	Creates new retro settings for a team
//	@Tags			retro-settings
//	@Produce		json
//	@Param			teamId		path	string													true	"Team ID"
//	@Param			settings	body	retroSettingsRequestBody								true	"retro settings object to create"
//	@Success		201			object	standardJsonResponse{data=thunderdome.RetroSettings}	"returns created retro settings"
//	@Failure		400			object	standardJsonResponse{}
//	@Failure		500			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/retro-settings [post]
func (s *Service) handleCreateTeamRetroSettings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		teamID := vars["teamId"]
		teamIDErr := validate.Var(teamID, "required,uuid")
		if teamIDErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, teamIDErr.Error()))
			return
		}

		var settingsReq retroSettingsRequestBody
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		jsonErr := json.Unmarshal(body, &settingsReq)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(settingsReq)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		settings := &thunderdome.RetroSettings{
			TeamID:                &teamID,
			MaxVotes:              settingsReq.MaxVotes,
			AllowMultipleVotes:    settingsReq.AllowMultipleVotes,
			BrainstormVisibility:  settingsReq.BrainstormVisibility,
			PhaseTimeLimit:        settingsReq.PhaseTimeLimit,
			PhaseAutoAdvance:      settingsReq.PhaseAutoAdvance,
			AllowCumulativeVoting: settingsReq.AllowCumulativeVoting,
			TemplateID:            settingsReq.TemplateID,
			JoinCode:              settingsReq.JoinCode,
			FacilitatorCode:       settingsReq.FacilitatorCode,
		}

		createdSettings, err := s.RetroDataSvc.CreateSettings(ctx, settings)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleCreateTeamRetroSettings error", zap.Error(err),
				zap.String("team_id", teamID), zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusCreated, createdSettings, nil)
	}
}

// handleCreateDepartmentRetroSettings creates new retro settings for a department
//
//	@Summary		Create Department Retro Settings
//	@Description	Creates new retro settings for a department
//	@Tags			retro-settings
//	@Produce		json
//	@Param			deptId		path	string													true	"Department ID"
//	@Param			settings	body	retroSettingsRequestBody								true	"retro settings object to create"
//	@Success		201			object	standardJsonResponse{data=thunderdome.RetroSettings}	"returns created retro settings"
//	@Failure		400			object	standardJsonResponse{}
//	@Failure		500			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/departments/{deptId}/retro-settings [post]
func (s *Service) handleCreateDepartmentRetroSettings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		deptID := vars["deptId"]
		deptIDErr := validate.Var(deptID, "required,uuid")
		if deptIDErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, deptIDErr.Error()))
			return
		}

		var settingsReq retroSettingsRequestBody
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		jsonErr := json.Unmarshal(body, &settingsReq)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(settingsReq)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		settings := &thunderdome.RetroSettings{
			DepartmentID:          &deptID,
			MaxVotes:              settingsReq.MaxVotes,
			AllowMultipleVotes:    settingsReq.AllowMultipleVotes,
			BrainstormVisibility:  settingsReq.BrainstormVisibility,
			PhaseTimeLimit:        settingsReq.PhaseTimeLimit,
			PhaseAutoAdvance:      settingsReq.PhaseAutoAdvance,
			AllowCumulativeVoting: settingsReq.AllowCumulativeVoting,
			TemplateID:            settingsReq.TemplateID,
			JoinCode:              settingsReq.JoinCode,
			FacilitatorCode:       settingsReq.FacilitatorCode,
		}

		createdSettings, err := s.RetroDataSvc.CreateSettings(ctx, settings)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleCreateDepartmentRetroSettings error", zap.Error(err),
				zap.String("dept_id", deptID), zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusCreated, createdSettings, nil)
	}
}

// handleGetOrganizationRetroSettings gets retro settings for a specific organization
//
//	@Summary		Get Organization Retro Settings
//	@Description	get retro settings for a specific organization
//	@Tags			retro-settings
//	@Produce		json
//	@Param			orgId	path	string	true	"Organization ID"
//	@Success		200		object	standardJsonResponse{data=thunderdome.RetroSettings}
//	@Failure		404		object	standardJsonResponse{}
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/retro-settings [get]
func (s *Service) handleGetOrganizationRetroSettings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID, _ := ctx.Value(contextKeyUserID).(*string)
		vars := mux.Vars(r)
		orgID := vars["orgId"]
		orgIDErr := validate.Var(orgID, "required,uuid")
		if orgIDErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, orgIDErr.Error()))
			return
		}

		settings, err := s.RetroDataSvc.GetSettingsByOrganization(ctx, orgID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleGetOrganizationRetroSettings error", zap.Error(err),
				zap.String("org_id", orgID), zap.Stringp("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		if settings == nil {
			s.Failure(w, r, http.StatusNotFound, Errorf(ENOTFOUND, "Retro settings not found"))
			return
		}

		s.Success(w, r, http.StatusOK, settings, nil)
	}
}

// handleOrganizationRetroSettingsUpdate updates retro settings for a specific organization
//
//	@Summary		Update Organization Retro Settings
//	@Description	Updates retro settings for a specific organization
//	@Tags			retro-settings
//	@Produce		json
//	@Param			orgId		path	string													true	"Organization ID"
//	@Param			settings	body	retroSettingsRequestBody								true	"retro settings object to update"
//	@Success		200			object	standardJsonResponse{data=thunderdome.RetroSettings}	"returns updated retro settings"
//	@Failure		400			object	standardJsonResponse{}
//	@Failure		500			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/retro-settings [put]
func (s *Service) handleOrganizationRetroSettingsUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		orgID := vars["orgId"]
		orgIDErr := validate.Var(orgID, "required,uuid")
		if orgIDErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, orgIDErr.Error()))
			return
		}

		var settingsReq retroSettingsRequestBody
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		jsonErr := json.Unmarshal(body, &settingsReq)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(settingsReq)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		settings := &thunderdome.RetroSettings{
			OrganizationID:        &orgID,
			MaxVotes:              settingsReq.MaxVotes,
			AllowMultipleVotes:    settingsReq.AllowMultipleVotes,
			BrainstormVisibility:  settingsReq.BrainstormVisibility,
			PhaseTimeLimit:        settingsReq.PhaseTimeLimit,
			PhaseAutoAdvance:      settingsReq.PhaseAutoAdvance,
			AllowCumulativeVoting: settingsReq.AllowCumulativeVoting,
			TemplateID:            settingsReq.TemplateID,
			JoinCode:              settingsReq.JoinCode,
			FacilitatorCode:       settingsReq.FacilitatorCode,
		}

		updatedSettings, err := s.RetroDataSvc.UpdateOrganizationSettings(ctx, settings)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleOrganizationRetroSettingsUpdate error", zap.Error(err),
				zap.String("org_id", orgID), zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, updatedSettings, nil)
	}
}

// handleGetTeamRetroSettings gets retro settings for a specific team
//
//	@Summary		Get Team Retro Settings
//	@Description	get retro settings for a specific team
//	@Tags			retro-settings
//	@Produce		json
//	@Param			teamId	path	string	true	"Team ID"
//	@Success		200		object	standardJsonResponse{data=thunderdome.RetroSettings}
//	@Failure		404		object	standardJsonResponse{}
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/retro-settings [get]
func (s *Service) handleGetTeamRetroSettings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID, _ := ctx.Value(contextKeyUserID).(*string)
		vars := mux.Vars(r)
		teamID := vars["teamId"]
		teamIDErr := validate.Var(teamID, "required,uuid")
		if teamIDErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, teamIDErr.Error()))
			return
		}

		settings, err := s.RetroDataSvc.GetSettingsByTeam(ctx, teamID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleGetTeamRetroSettings error", zap.Error(err),
				zap.String("team_id", teamID), zap.Stringp("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		if settings == nil {
			s.Failure(w, r, http.StatusNotFound, Errorf(ENOTFOUND, "Retro settings not found"))
			return
		}

		s.Success(w, r, http.StatusOK, settings, nil)
	}
}

// handleTeamRetroSettingsUpdate updates retro settings for a specific team
//
//	@Summary		Update Team Retro Settings
//	@Description	Updates retro settings for a specific team
//	@Tags			retro-settings
//	@Produce		json
//	@Param			teamId		path	string													true	"Team ID"
//	@Param			settings	body	retroSettingsRequestBody								true	"retro settings object to update"
//	@Success		200			object	standardJsonResponse{data=thunderdome.RetroSettings}	"returns updated retro settings"
//	@Failure		400			object	standardJsonResponse{}
//	@Failure		500			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/retro-settings [put]
func (s *Service) handleTeamRetroSettingsUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		teamID := vars["teamId"]
		teamIDErr := validate.Var(teamID, "required,uuid")
		if teamIDErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, teamIDErr.Error()))
			return
		}

		var settingsReq retroSettingsRequestBody
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		jsonErr := json.Unmarshal(body, &settingsReq)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(settingsReq)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		settings := &thunderdome.RetroSettings{
			TeamID:                &teamID,
			MaxVotes:              settingsReq.MaxVotes,
			AllowMultipleVotes:    settingsReq.AllowMultipleVotes,
			BrainstormVisibility:  settingsReq.BrainstormVisibility,
			PhaseTimeLimit:        settingsReq.PhaseTimeLimit,
			PhaseAutoAdvance:      settingsReq.PhaseAutoAdvance,
			AllowCumulativeVoting: settingsReq.AllowCumulativeVoting,
			TemplateID:            settingsReq.TemplateID,
			JoinCode:              settingsReq.JoinCode,
			FacilitatorCode:       settingsReq.FacilitatorCode,
		}

		updatedSettings, err := s.RetroDataSvc.UpdateTeamSettings(ctx, settings)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleTeamRetroSettingsUpdate error", zap.Error(err),
				zap.String("team_id", teamID), zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, updatedSettings, nil)
	}
}

// handleGetDepartmentRetroSettings gets retro settings for a specific department
//
//	@Summary		Get Department Retro Settings
//	@Description	get retro settings for a specific department
//	@Tags			retro-settings
//	@Produce		json
//	@Param			orgId	path	string	true	"Organization ID"
//	@Param			deptId	path	string	true	"Department ID"
//	@Success		200		object	standardJsonResponse{data=thunderdome.RetroSettings}
//	@Failure		404		object	standardJsonResponse{}
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/departments/{deptId}/retro-settings [get]
func (s *Service) handleGetDepartmentRetroSettings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID, _ := ctx.Value(contextKeyUserID).(*string)
		vars := mux.Vars(r)
		deptID := vars["deptId"]
		deptIDErr := validate.Var(deptID, "required,uuid")
		if deptIDErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, deptIDErr.Error()))
			return
		}

		settings, err := s.RetroDataSvc.GetSettingsByDepartment(ctx, deptID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleGetDepartmentRetroSettings error", zap.Error(err),
				zap.String("dept_id", deptID), zap.Stringp("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		if settings == nil {
			s.Failure(w, r, http.StatusNotFound, Errorf(ENOTFOUND, "Retro settings not found"))
			return
		}

		s.Success(w, r, http.StatusOK, settings, nil)
	}
}

// handleDepartmentRetroSettingsUpdate updates retro settings for a specific department
//
//	@Summary		Update Department Retro Settings
//	@Description	Updates retro settings for a specific department
//	@Tags			retro-settings
//	@Produce		json
//	@Param			orgId		path	string													true	"Organization ID"
//	@Param			deptId		path	string													true	"Department ID"
//	@Param			settings	body	retroSettingsRequestBody								true	"retro settings object to update"
//	@Success		200			object	standardJsonResponse{data=thunderdome.RetroSettings}	"returns updated retro settings"
//	@Failure		400			object	standardJsonResponse{}
//	@Failure		500			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/departments/{deptId}/retro-settings [put]
func (s *Service) handleDepartmentRetroSettingsUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		deptID := vars["deptId"]
		deptIDErr := validate.Var(deptID, "required,uuid")
		if deptIDErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, deptIDErr.Error()))
			return
		}

		var settingsReq retroSettingsRequestBody
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		jsonErr := json.Unmarshal(body, &settingsReq)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(settingsReq)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		settings := &thunderdome.RetroSettings{
			DepartmentID:          &deptID,
			MaxVotes:              settingsReq.MaxVotes,
			AllowMultipleVotes:    settingsReq.AllowMultipleVotes,
			BrainstormVisibility:  settingsReq.BrainstormVisibility,
			PhaseTimeLimit:        settingsReq.PhaseTimeLimit,
			PhaseAutoAdvance:      settingsReq.PhaseAutoAdvance,
			AllowCumulativeVoting: settingsReq.AllowCumulativeVoting,
			TemplateID:            settingsReq.TemplateID,
			JoinCode:              settingsReq.JoinCode,
			FacilitatorCode:       settingsReq.FacilitatorCode,
		}

		updatedSettings, err := s.RetroDataSvc.UpdateDepartmentSettings(ctx, settings)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleDepartmentRetroSettingsUpdate error", zap.Error(err),
				zap.String("dept_id", deptID), zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, updatedSettings, nil)
	}
}

// handleDeleteRetroSettings deletes retro settings
//
//	@Summary		Delete Retro Settings
//	@Description	Deletes retro settings for an organization, department, or team
//	@Tags			admin, retro-settings
//	@Produce		json
//	@Param			id	path	string					true	"Settings ID"
//	@Success		200	object	standardJsonResponse{}	"returns success message"
//	@Failure		400	object	standardJsonResponse{}
//	@Failure		500	object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/admin/retro-settings/{id} [delete]
func (s *Service) handleDeleteRetroSettings() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		settingsID := vars["id"]
		idErr := validate.Var(settingsID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		err := s.RetroDataSvc.DeleteSettings(ctx, settingsID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleDeleteRetroSettings error", zap.Error(err),
				zap.String("settings_id", settingsID), zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, "Retro settings deleted successfully", nil)
	}
}

// handleGetRetroSettingsByID gets retro settings by ID
//
//	@Summary		Get Retro Settings by ID
//	@Description	get retro settings by ID
//	@Tags			admin, retro-settings
//	@Produce		json
//	@Param			id	path	string	true	"Settings ID"
//	@Success		200	object	standardJsonResponse{data=thunderdome.RetroSettings}
//	@Failure		404	object	standardJsonResponse{}
//	@Failure		500	object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/admin/retro-settings/{id} [get]
func (s *Service) handleGetRetroSettingsByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID, _ := ctx.Value(contextKeyUserID).(*string)
		vars := mux.Vars(r)
		settingsID := vars["id"]
		idErr := validate.Var(settingsID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		settings, err := s.RetroDataSvc.GetSettingsByID(ctx, settingsID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleGetRetroSettingsByID error", zap.Error(err),
				zap.String("settings_id", settingsID), zap.Stringp("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		if settings == nil {
			s.Failure(w, r, http.StatusNotFound, Errorf(ENOTFOUND, "Retro settings not found"))
			return
		}

		s.Success(w, r, http.StatusOK, settings, nil)
	}
}
