package http

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"

	"go.uber.org/zap"

	"github.com/gorilla/mux"
)

type pokerSettingsRequestBody struct {
	AutoFinishVoting     bool    `json:"autoFinishVoting"`
	PointAverageRounding string  `json:"pointAverageRounding" validate:"oneof=ceil floor round"`
	HideVoterIdentity    bool    `json:"hideVoterIdentity"`
	EstimationScaleID    *string `json:"estimationScaleId" validate:"omitempty,uuid"`
	JoinCode             string  `json:"joinCode"`
	FacilitatorCode      string  `json:"facilitatorCode"`
}

// handleCreateOrganizationPokerSettings creates new poker settings for an organization
//
//	@Summary		Create Organization Poker Settings
//	@Description	Creates new poker settings for an organization
//	@Tags			poker-settings
//	@Produce		json
//	@Param			orgId		path	string													true	"Organization ID"
//	@Param			settings	body	pokerSettingsRequestBody								true	"poker settings object to create"
//	@Success		201			object	standardJsonResponse{data=thunderdome.PokerSettings}	"returns created poker settings"
//	@Failure		400			object	standardJsonResponse{}
//	@Failure		500			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/poker-settings [post]
func (s *Service) handleCreateOrganizationPokerSettings() http.HandlerFunc {
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

		var settingsReq pokerSettingsRequestBody
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

		settings := &thunderdome.PokerSettings{
			OrganizationID:       &orgID,
			AutoFinishVoting:     settingsReq.AutoFinishVoting,
			PointAverageRounding: settingsReq.PointAverageRounding,
			HideVoterIdentity:    settingsReq.HideVoterIdentity,
			EstimationScaleID:    settingsReq.EstimationScaleID,
			JoinCode:             settingsReq.JoinCode,
			FacilitatorCode:      settingsReq.FacilitatorCode,
		}

		createdSettings, err := s.PokerDataSvc.CreateSettings(ctx, settings)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleCreateOrganizationPokerSettings error", zap.Error(err),
				zap.String("org_id", orgID), zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusCreated, createdSettings, nil)
	}
}

// handleCreateTeamPokerSettings creates new poker settings for a team
//
//	@Summary		Create Team Poker Settings
//	@Description	Creates new poker settings for a team
//	@Tags			poker-settings
//	@Produce		json
//	@Param			teamId		path	string													true	"Team ID"
//	@Param			settings	body	pokerSettingsRequestBody								true	"poker settings object to create"
//	@Success		201			object	standardJsonResponse{data=thunderdome.PokerSettings}	"returns created poker settings"
//	@Failure		400			object	standardJsonResponse{}
//	@Failure		500			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/poker-settings [post]
func (s *Service) handleCreateTeamPokerSettings() http.HandlerFunc {
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

		var settingsReq pokerSettingsRequestBody
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

		settings := &thunderdome.PokerSettings{
			TeamID:               &teamID,
			AutoFinishVoting:     settingsReq.AutoFinishVoting,
			PointAverageRounding: settingsReq.PointAverageRounding,
			HideVoterIdentity:    settingsReq.HideVoterIdentity,
			EstimationScaleID:    settingsReq.EstimationScaleID,
			JoinCode:             settingsReq.JoinCode,
			FacilitatorCode:      settingsReq.FacilitatorCode,
		}

		createdSettings, err := s.PokerDataSvc.CreateSettings(ctx, settings)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleCreateTeamPokerSettings error", zap.Error(err),
				zap.String("team_id", teamID), zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusCreated, createdSettings, nil)
	}
}

// handleCreateDepartmentPokerSettings creates new poker settings for a department
//
//	@Summary		Create Department Poker Settings
//	@Description	Creates new poker settings for a department
//	@Tags			poker-settings
//	@Produce		json
//	@Param			deptId		path	string													true	"Department ID"
//	@Param			settings	body	pokerSettingsRequestBody								true	"poker settings object to create"
//	@Success		201			object	standardJsonResponse{data=thunderdome.PokerSettings}	"returns created poker settings"
//	@Failure		400			object	standardJsonResponse{}
//	@Failure		500			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/departments/{deptId}/poker-settings [post]
func (s *Service) handleCreateDepartmentPokerSettings() http.HandlerFunc {
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

		var settingsReq pokerSettingsRequestBody
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

		settings := &thunderdome.PokerSettings{
			DepartmentID:         &deptID,
			AutoFinishVoting:     settingsReq.AutoFinishVoting,
			PointAverageRounding: settingsReq.PointAverageRounding,
			HideVoterIdentity:    settingsReq.HideVoterIdentity,
			EstimationScaleID:    settingsReq.EstimationScaleID,
			JoinCode:             settingsReq.JoinCode,
			FacilitatorCode:      settingsReq.FacilitatorCode,
		}

		createdSettings, err := s.PokerDataSvc.CreateSettings(ctx, settings)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleCreateDepartmentPokerSettings error", zap.Error(err),
				zap.String("dept_id", deptID), zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusCreated, createdSettings, nil)
	}
}

// handleGetOrganizationPokerSettings gets poker settings for a specific organization
//
//	@Summary		Get Organization Poker Settings
//	@Description	get poker settings for a specific organization
//	@Tags			poker-settings
//	@Produce		json
//	@Param			orgId	path	string	true	"Organization ID"
//	@Success		200		object	standardJsonResponse{data=thunderdome.PokerSettings}
//	@Failure		404		object	standardJsonResponse{}
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/poker-settings [get]
func (s *Service) handleGetOrganizationPokerSettings() http.HandlerFunc {
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

		settings, err := s.PokerDataSvc.GetSettingsByOrganization(ctx, orgID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleGetOrganizationPokerSettings error", zap.Error(err),
				zap.String("org_id", orgID), zap.Stringp("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, settings, nil)
	}
}

// handleOrganizationPokerSettingsUpdate updates poker settings for a specific organization
//
//	@Summary		Update Organization Poker Settings
//	@Description	Updates poker settings for a specific organization
//	@Tags			poker-settings
//	@Produce		json
//	@Param			orgId		path	string													true	"Organization ID"
//	@Param			settings	body	pokerSettingsRequestBody								true	"poker settings object to update"
//	@Success		200			object	standardJsonResponse{data=thunderdome.PokerSettings}	"returns updated poker settings"
//	@Failure		400			object	standardJsonResponse{}
//	@Failure		500			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/poker-settings [put]
func (s *Service) handleOrganizationPokerSettingsUpdate() http.HandlerFunc {
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

		var settingsReq pokerSettingsRequestBody
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

		settings := &thunderdome.PokerSettings{
			OrganizationID:       &orgID,
			AutoFinishVoting:     settingsReq.AutoFinishVoting,
			PointAverageRounding: settingsReq.PointAverageRounding,
			HideVoterIdentity:    settingsReq.HideVoterIdentity,
			EstimationScaleID:    settingsReq.EstimationScaleID,
			JoinCode:             settingsReq.JoinCode,
			FacilitatorCode:      settingsReq.FacilitatorCode,
		}

		updatedSettings, err := s.PokerDataSvc.UpdateOrganizationSettings(ctx, settings)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleOrganizationPokerSettingsUpdate error", zap.Error(err),
				zap.String("org_id", orgID), zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, updatedSettings, nil)
	}
}

// handleGetTeamPokerSettings gets poker settings for a specific team
//
//	@Summary		Get Team Poker Settings
//	@Description	get poker settings for a specific team
//	@Tags			poker-settings
//	@Produce		json
//	@Param			teamId	path	string	true	"Team ID"
//	@Success		200		object	standardJsonResponse{data=thunderdome.PokerSettings}
//	@Failure		404		object	standardJsonResponse{}
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/poker-settings [get]
func (s *Service) handleGetTeamPokerSettings() http.HandlerFunc {
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

		settings, err := s.PokerDataSvc.GetSettingsByTeam(ctx, teamID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleGetTeamPokerSettings error", zap.Error(err),
				zap.String("team_id", teamID), zap.Stringp("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, settings, nil)
	}
}

// handleTeamPokerSettingsUpdate updates poker settings for a specific team
//
//	@Summary		Update Team Poker Settings
//	@Description	Updates poker settings for a specific team
//	@Tags			poker-settings
//	@Produce		json
//	@Param			teamId		path	string													true	"Team ID"
//	@Param			settings	body	pokerSettingsRequestBody								true	"poker settings object to update"
//	@Success		200			object	standardJsonResponse{data=thunderdome.PokerSettings}	"returns updated poker settings"
//	@Failure		400			object	standardJsonResponse{}
//	@Failure		500			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/poker-settings [put]
func (s *Service) handleTeamPokerSettingsUpdate() http.HandlerFunc {
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

		var settingsReq pokerSettingsRequestBody
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

		settings := &thunderdome.PokerSettings{
			TeamID:               &teamID,
			AutoFinishVoting:     settingsReq.AutoFinishVoting,
			PointAverageRounding: settingsReq.PointAverageRounding,
			HideVoterIdentity:    settingsReq.HideVoterIdentity,
			EstimationScaleID:    settingsReq.EstimationScaleID,
			JoinCode:             settingsReq.JoinCode,
			FacilitatorCode:      settingsReq.FacilitatorCode,
		}

		updatedSettings, err := s.PokerDataSvc.UpdateTeamSettings(ctx, settings)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleTeamPokerSettingsUpdate error", zap.Error(err),
				zap.String("team_id", teamID), zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, updatedSettings, nil)
	}
}

// handleGetDepartmentPokerSettings gets poker settings for a specific department
//
//	@Summary		Get Department Poker Settings
//	@Description	get poker settings for a specific department
//	@Tags			poker-settings
//	@Produce		json
//	@Param			orgId	path	string	true	"Organization ID"
//	@Param			deptId	path	string	true	"Department ID"
//	@Success		200		object	standardJsonResponse{data=thunderdome.PokerSettings}
//	@Failure		404		object	standardJsonResponse{}
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/departments/{deptId}/poker-settings [get]
func (s *Service) handleGetDepartmentPokerSettings() http.HandlerFunc {
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

		settings, err := s.PokerDataSvc.GetSettingsByDepartment(ctx, deptID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleGetDepartmentPokerSettings error", zap.Error(err),
				zap.String("dept_id", deptID), zap.Stringp("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, settings, nil)
	}
}

// handleDepartmentPokerSettingsUpdate updates poker settings for a specific department
//
//	@Summary		Update Department Poker Settings
//	@Description	Updates poker settings for a specific department
//	@Tags			poker-settings
//	@Produce		json
//	@Param			orgId		path	string													true	"Organization ID"
//	@Param			deptId		path	string													true	"Department ID"
//	@Param			settings	body	pokerSettingsRequestBody								true	"poker settings object to update"
//	@Success		200			object	standardJsonResponse{data=thunderdome.PokerSettings}	"returns updated poker settings"
//	@Failure		400			object	standardJsonResponse{}
//	@Failure		500			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/departments/{deptId}/poker-settings [put]
func (s *Service) handleDepartmentPokerSettingsUpdate() http.HandlerFunc {
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

		var settingsReq pokerSettingsRequestBody
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

		settings := &thunderdome.PokerSettings{
			DepartmentID:         &deptID,
			AutoFinishVoting:     settingsReq.AutoFinishVoting,
			PointAverageRounding: settingsReq.PointAverageRounding,
			HideVoterIdentity:    settingsReq.HideVoterIdentity,
			EstimationScaleID:    settingsReq.EstimationScaleID,
			JoinCode:             settingsReq.JoinCode,
			FacilitatorCode:      settingsReq.FacilitatorCode,
		}

		updatedSettings, err := s.PokerDataSvc.UpdateDepartmentSettings(ctx, settings)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleDepartmentPokerSettingsUpdate error", zap.Error(err),
				zap.String("dept_id", deptID), zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, updatedSettings, nil)
	}
}

// handleDeletePokerSettings deletes poker settings
//
//	@Summary		Delete Poker Settings
//	@Description	Deletes poker settings for an organization, department, or team
//	@Tags			admin, poker-settings
//	@Produce		json
//	@Param			id	path	string					true	"Settings ID"
//	@Success		200	object	standardJsonResponse{}	"returns success message"
//	@Failure		400	object	standardJsonResponse{}
//	@Failure		500	object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/admin/poker-settings/{id} [delete]
func (s *Service) handleDeletePokerSettings() http.HandlerFunc {
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

		err := s.PokerDataSvc.DeleteSettings(ctx, settingsID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleDeletePokerSettings error", zap.Error(err),
				zap.String("settings_id", settingsID), zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, "Poker settings deleted successfully", nil)
	}
}

// handleGetPokerSettingsByID gets poker settings by ID
//
//	@Summary		Get Poker Settings by ID
//	@Description	get poker settings by ID
//	@Tags			admin, poker-settings
//	@Produce		json
//	@Param			id	path	string	true	"Settings ID"
//	@Success		200	object	standardJsonResponse{data=thunderdome.PokerSettings}
//	@Failure		404	object	standardJsonResponse{}
//	@Failure		500	object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/admin/poker-settings/{id} [get]
func (s *Service) handleGetPokerSettingsByID() http.HandlerFunc {
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

		settings, err := s.PokerDataSvc.GetSettingsByID(ctx, settingsID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleGetPokerSettingsByID error", zap.Error(err),
				zap.String("settings_id", settingsID), zap.Stringp("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		if settings == nil {
			s.Failure(w, r, http.StatusNotFound, Errorf(ENOTFOUND, "Poker settings not found"))
			return
		}

		s.Success(w, r, http.StatusOK, settings, nil)
	}
}
