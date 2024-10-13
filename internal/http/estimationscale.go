package http

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"

	"go.uber.org/zap"

	"github.com/gorilla/mux"
)

type estimationScaleRequestBody struct {
	Name         string   `json:"name" validate:"required"`
	Description  string   `json:"description"`
	ScaleType    string   `json:"scaleType" validate:"required,oneof=modified_fibonacci fibonacci t_shirt powers_of_two custom"`
	Values       []string `json:"values" validate:"required,min=2"`
	IsPublic     bool     `json:"isPublic"`
	DefaultScale bool     `json:"defaultScale"`
}

// handleGetEstimationScales gets a list of estimation scales
//
//	@Summary		Get Estimation Scales
//	@Description	get list of estimation scales
//	@Tags			estimation-scale
//	@Produce		json
//	@Param			limit	query	int	false	"Max number of results to return"
//	@Param			offset	query	int	false	"Starting point to return rows from, should be multiplied by limit or 0"
//	@Success		200		object	standardJsonResponse{data=[]thunderdome.EstimationScale}
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/admin/estimation-scales [get]
func (s *Service) handleGetEstimationScales() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID, _ := ctx.Value(contextKeyUserID).(*string)
		limit, offset := getLimitOffsetFromRequest(r)
		scales, count, err := s.PokerDataSvc.GetEstimationScales(ctx, limit, offset)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleGetEstimationScales error", zap.Error(err),
				zap.Int("limit", limit), zap.Int("offset", offset),
				zap.Stringp("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		meta := &pagination{
			Count:  count,
			Offset: offset,
			Limit:  limit,
		}

		s.Success(w, r, http.StatusOK, scales, meta)
	}
}

// handleEstimationScaleCreate creates a new estimation scale
//
//	@Summary		Create Estimation Scale
//	@Description	Creates an estimation scale
//	@Tags			estimation-scale
//	@Produce		json
//	@Param			scale	body	estimationScaleRequestBody								true	"new estimation scale object"
//	@Success		200		object	standardJsonResponse{data=thunderdome.EstimationScale}	"returns created estimation scale"
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/admin/estimation-scales [post]
func (s *Service) handleEstimationScaleCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		var scale = estimationScaleRequestBody{}
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		jsonErr := json.Unmarshal(body, &scale)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(scale)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		es := thunderdome.EstimationScale{
			Name:         scale.Name,
			Description:  scale.Description,
			ScaleType:    scale.ScaleType,
			Values:       scale.Values,
			DefaultScale: scale.DefaultScale,
			IsPublic:     scale.IsPublic,
			CreatedBy:    sessionUserID,
		}

		createdScale, err := s.PokerDataSvc.CreateEstimationScale(ctx, &es)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleEstimationScaleCreate error", zap.Error(err),
				zap.String("scale_name", scale.Name), zap.String("scale_type", scale.ScaleType),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, createdScale, nil)
	}
}

// handleEstimationScaleUpdate updates an estimation scale
//
//	@Summary		Update Estimation Scale
//	@Description	Updates an Estimation Scale
//	@Tags			estimation-scale
//	@Produce		json
//	@Param			scaleId	path	string													true	"the estimation scale ID to update"
//	@Param			scale	body	estimationScaleRequestBody								true	"estimation scale object to update"
//	@Success		200		object	standardJsonResponse{data=thunderdome.EstimationScale}	"returns updated estimation scale"
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/admin/estimation-scales/{scaleId} [put]
func (s *Service) handleEstimationScaleUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		scaleID := vars["scaleId"]
		idErr := validate.Var(scaleID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		var scale = estimationScaleRequestBody{}
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		jsonErr := json.Unmarshal(body, &scale)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(scale)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		es := thunderdome.EstimationScale{
			ID:           scaleID,
			Name:         scale.Name,
			Description:  scale.Description,
			ScaleType:    scale.ScaleType,
			Values:       scale.Values,
			DefaultScale: scale.DefaultScale,
			IsPublic:     scale.IsPublic,
		}

		updatedScale, err := s.PokerDataSvc.UpdateEstimationScale(ctx, &es)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleEstimationScaleUpdate error", zap.Error(err), zap.String("scale_id", scaleID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, updatedScale, nil)
	}
}

// handleEstimationScaleDelete handles deleting an estimation scale
//
//	@Summary		Delete Estimation Scale
//	@Description	Deletes an Estimation Scale
//	@Tags			estimation-scale
//	@Produce		json
//	@Param			scaleId	path	string					true	"the estimation scale ID to delete"
//	@Success		200		object	standardJsonResponse{}	"returns success message"
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/admin/estimation-scales/{scaleId} [delete]
func (s *Service) handleEstimationScaleDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		scaleID := vars["scaleId"]
		idErr := validate.Var(scaleID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		err := s.PokerDataSvc.DeleteEstimationScale(ctx, scaleID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleEstimationScaleDelete error", zap.Error(err), zap.String("scale_id", scaleID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, "Estimation scale deleted successfully", nil)
	}
}

// handleGetOrganizationEstimationScales gets a list of estimation scales for a specific organization
//
//	@Summary		Get Organization Estimation Scales
//	@Description	get list of estimation scales for a specific organization
//	@Tags			estimation-scale
//	@Produce		json
//	@Param			orgId	path	string	true	"Organization ID"
//	@Param			limit	query	int		false	"Max number of results to return"
//	@Param			offset	query	int		false	"Starting point to return rows from, should be multiplied by limit or 0"
//	@Success		200		object	standardJsonResponse{data=[]thunderdome.EstimationScale}
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/estimation-scales [get]
func (s *Service) handleGetOrganizationEstimationScales() http.HandlerFunc {
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
		limit, offset := getLimitOffsetFromRequest(r)

		scales, count, err := s.PokerDataSvc.GetOrganizationEstimationScales(ctx, orgID, limit, offset)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleGetOrganizationEstimationScales error", zap.Error(err),
				zap.String("org_id", orgID), zap.Int("limit", limit), zap.Int("offset", offset),
				zap.Stringp("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		meta := &pagination{
			Count:  count,
			Offset: offset,
			Limit:  limit,
		}

		s.Success(w, r, http.StatusOK, scales, meta)
	}
}

type privateEstimationScaleRequestBody struct {
	Name         string   `json:"name" validate:"required"`
	Description  string   `json:"description"`
	Values       []string `json:"values" validate:"required,min=2"`
	DefaultScale bool     `json:"defaultScale"`
}

// handleOrganizationEstimationScaleCreate creates a new estimation scale for a specific organization
//
//	@Summary		Create Organization Estimation Scale
//	@Description	Creates an estimation scale for a specific organization
//	@Tags			estimation-scale
//	@Produce		json
//	@Param			orgId	path	string													true	"Organization ID"
//	@Param			scale	body	privateEstimationScaleRequestBody						true	"new estimation scale object"
//	@Success		200		object	standardJsonResponse{data=thunderdome.EstimationScale}	"returns created estimation scale"
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/estimation-scales [post]
func (s *Service) handleOrganizationEstimationScaleCreate() http.HandlerFunc {
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

		var scale = privateEstimationScaleRequestBody{}
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		jsonErr := json.Unmarshal(body, &scale)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(scale)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		es := thunderdome.EstimationScale{
			Name:           scale.Name,
			Description:    scale.Description,
			ScaleType:      "custom",
			Values:         scale.Values,
			DefaultScale:   scale.DefaultScale,
			OrganizationID: orgID,
			CreatedBy:      sessionUserID,
			IsPublic:       false,
		}

		createdScale, err := s.PokerDataSvc.CreateEstimationScale(ctx, &es)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleOrganizationEstimationScaleCreate error", zap.Error(err),
				zap.String("org_id", orgID), zap.String("scale_name", scale.Name),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, createdScale, nil)
	}
}

// handleGetTeamEstimationScales gets a list of estimation scales for a specific team
//
//	@Summary		Get Team Estimation Scales
//	@Description	get list of estimation scales for a specific team
//	@Tags			estimation-scale
//	@Produce		json
//	@Param			teamId	path	string	true	"Team ID"
//	@Param			limit	query	int		false	"Max number of results to return"
//	@Param			offset	query	int		false	"Starting point to return rows from, should be multiplied by limit or 0"
//	@Success		200		object	standardJsonResponse{data=[]thunderdome.EstimationScale}
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/estimation-scales [get]
func (s *Service) handleGetTeamEstimationScales() http.HandlerFunc {
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
		limit, offset := getLimitOffsetFromRequest(r)

		scales, count, err := s.PokerDataSvc.GetTeamEstimationScales(ctx, teamID, limit, offset)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleGetTeamEstimationScales error", zap.Error(err),
				zap.String("team_id", teamID), zap.Int("limit", limit), zap.Int("offset", offset),
				zap.Stringp("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		meta := &pagination{
			Count:  count,
			Offset: offset,
			Limit:  limit,
		}

		s.Success(w, r, http.StatusOK, scales, meta)
	}
}

// handleTeamEstimationScaleCreate creates a new estimation scale for a specific team
//
//	@Summary		Create Team Estimation Scale
//	@Description	Creates an estimation scale for a specific team
//	@Tags			estimation-scale
//	@Produce		json
//	@Param			teamId	path	string													true	"Team ID"
//	@Param			scale	body	privateEstimationScaleRequestBody						true	"new estimation scale object"
//	@Success		200		object	standardJsonResponse{data=thunderdome.EstimationScale}	"returns created estimation scale"
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/estimation-scales [post]
func (s *Service) handleTeamEstimationScaleCreate() http.HandlerFunc {
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

		var scale = privateEstimationScaleRequestBody{}
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		jsonErr := json.Unmarshal(body, &scale)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(scale)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		es := thunderdome.EstimationScale{
			Name:         scale.Name,
			Description:  scale.Description,
			ScaleType:    "custom",
			Values:       scale.Values,
			DefaultScale: scale.DefaultScale,
			TeamID:       teamID,
			CreatedBy:    sessionUserID,
			IsPublic:     false,
		}

		createdScale, err := s.PokerDataSvc.CreateEstimationScale(ctx, &es)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleTeamEstimationScaleCreate error", zap.Error(err),
				zap.String("team_id", teamID), zap.String("scale_name", scale.Name),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, createdScale, nil)
	}
}

// handleGetPublicEstimationScales gets a list of all public estimation scales
//
//	@Summary		Get Public Estimation Scales
//	@Description	get list of all public estimation scales
//	@Tags			estimation-scale
//	@Produce		json
//	@Param			limit	query	int	false	"Max number of results to return"
//	@Param			offset	query	int	false	"Starting point to return rows from, should be multiplied by limit or 0"
//	@Success		200		object	standardJsonResponse{data=[]thunderdome.EstimationScale}
//	@Failure		500		object	standardJsonResponse{}
//	@Router			/estimation-scales/public [get]
func (s *Service) handleGetPublicEstimationScales() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		limit, offset := getLimitOffsetFromRequest(r)

		scales, count, err := s.PokerDataSvc.GetPublicEstimationScales(ctx, limit, offset)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleGetPublicEstimationScales error", zap.Error(err),
				zap.Int("limit", limit), zap.Int("offset", offset))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		meta := &pagination{
			Count:  count,
			Offset: offset,
			Limit:  limit,
		}

		s.Success(w, r, http.StatusOK, scales, meta)
	}
}

// handleGetPublicEstimationScale gets a specific public estimation scale by ID
//
//	@Summary		Get Public Estimation Scale
//	@Description	get a specific public estimation scale by ID
//	@Tags			estimation-scale
//	@Produce		json
//	@Param			scaleId	path	string	true	"Estimation Scale ID"
//	@Success		200		object	standardJsonResponse{data=thunderdome.EstimationScale}
//	@Failure		404		object	standardJsonResponse{}
//	@Failure		500		object	standardJsonResponse{}
//	@Router			/estimation-scales/public/{scaleId} [get]
func (s *Service) handleGetPublicEstimationScale() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		vars := mux.Vars(r)
		scaleID := vars["scaleId"]
		scaleIDErr := validate.Var(scaleID, "required,uuid")
		if scaleIDErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, scaleIDErr.Error()))
			return
		}

		scale, err := s.PokerDataSvc.GetPublicEstimationScale(ctx, scaleID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleGetPublicEstimationScale error", zap.Error(err),
				zap.String("scale_id", scaleID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, scale, nil)
	}
}

// handleOrganizationEstimationScaleDelete handles deleting an organization's estimation scale
//
//	@Summary		Delete Organization Estimation Scale
//	@Description	Deletes an Organization Estimation Scale
//	@Tags			estimation-scale
//	@Produce		json
//	@Param			scaleId	path	string					true	"the estimation scale ID to delete"
//	@Success		200		object	standardJsonResponse{}	"returns success message"
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organization/{orgId}/estimation-scales/{scaleId} [delete]
func (s *Service) handleOrganizationEstimationScaleDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		scaleID := vars["scaleId"]
		scaleIDErr := validate.Var(scaleID, "required,uuid")
		if scaleIDErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, scaleIDErr.Error()))
			return
		}
		orgID := vars["orgID"]
		orgIDErr := validate.Var(orgID, "required,uuid")
		if orgIDErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, orgIDErr.Error()))
			return
		}

		err := s.PokerDataSvc.DeleteOrganizationEstimationScale(ctx, orgID, scaleID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleOrganizationEstimationScaleDelete error", zap.Error(err), zap.String("scale_id", scaleID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, "Estimation scale deleted successfully", nil)
	}
}

// handleTeamEstimationScaleDelete handles deleting a team's estimation scale
//
//	@Summary		Delete Team Estimation Scale
//	@Description	Deletes an Team Estimation Scale
//	@Tags			estimation-scale
//	@Produce		json
//	@Param			scaleId	path	string					true	"the estimation scale ID to delete"
//	@Success		200		object	standardJsonResponse{}	"returns success message"
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/team/{teamId}/estimation-scales/{scaleId} [delete]
func (s *Service) handleTeamEstimationScaleDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		scaleID := vars["scaleId"]
		scaleIDErr := validate.Var(scaleID, "required,uuid")
		if scaleIDErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, scaleIDErr.Error()))
			return
		}
		teamID := vars["teamId"]
		teamIDErr := validate.Var(teamID, "required,uuid")
		if teamIDErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, teamIDErr.Error()))
			return
		}

		err := s.PokerDataSvc.DeleteTeamEstimationScale(ctx, teamID, scaleID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleTeamEstimationScaleDelete error", zap.Error(err), zap.String("scale_id", scaleID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, "Estimation scale deleted successfully", nil)
	}
}

// handleTeamEstimationScaleUpdate updates a team estimation scale
//
//	@Summary		Update Team Estimation Scale
//	@Description	Updates a Team Estimation Scale
//	@Tags			estimation-scale
//	@Produce		json
//	@Param			teamId	path	string													true	"Team ID"
//	@Param			scaleId	path	string													true	"the estimation scale ID to update"
//	@Param			scale	body	privateEstimationScaleRequestBody						true	"estimation scale object to update"
//	@Success		200		object	standardJsonResponse{data=thunderdome.EstimationScale}	"returns updated estimation scale"
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/estimation-scales/{scaleId} [put]
func (s *Service) handleTeamEstimationScaleUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		scaleID := vars["scaleId"]
		scaleIDErr := validate.Var(scaleID, "required,uuid")
		if scaleIDErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, scaleIDErr.Error()))
			return
		}
		teamID := vars["teamId"]
		teamIDErr := validate.Var(teamID, "required,uuid")
		if teamIDErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, teamIDErr.Error()))
			return
		}

		var scale = privateEstimationScaleRequestBody{}
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		jsonErr := json.Unmarshal(body, &scale)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(scale)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		es := thunderdome.EstimationScale{
			ID:           scaleID,
			Name:         scale.Name,
			Description:  scale.Description,
			ScaleType:    "custom",
			Values:       scale.Values,
			DefaultScale: scale.DefaultScale,
			TeamID:       teamID,
			IsPublic:     false,
		}

		updatedScale, err := s.PokerDataSvc.UpdateTeamEstimationScale(ctx, &es)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleEstimationScaleUpdate error", zap.Error(err),
				zap.String("scale_id", scaleID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, updatedScale, nil)
	}
}

// handleOrganizationEstimationScaleUpdate updates an organization estimation scale
//
//	@Summary		Update Organization Estimation Scale
//	@Description	Updates an Organization Estimation Scale
//	@Tags			estimation-scale
//	@Produce		json
//	@Param			orgId	path	string													true	"Organization ID"
//	@Param			scaleId	path	string													true	"the estimation scale ID to update"
//	@Param			scale	body	privateEstimationScaleRequestBody						true	"estimation scale object to update"
//	@Success		200		object	standardJsonResponse{data=thunderdome.EstimationScale}	"returns updated estimation scale"
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{orgId}/estimation-scales/{scaleId} [put]
func (s *Service) handleOrganizationEstimationScaleUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		vars := mux.Vars(r)
		scaleID := vars["scaleId"]
		scaleIDErr := validate.Var(scaleID, "required,uuid")
		if scaleIDErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, scaleIDErr.Error()))
			return
		}
		orgID := vars["orgId"]
		orgIDErr := validate.Var(orgID, "required,uuid")
		if orgIDErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, orgIDErr.Error()))
			return
		}

		var scale = privateEstimationScaleRequestBody{}
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		jsonErr := json.Unmarshal(body, &scale)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(scale)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		es := thunderdome.EstimationScale{
			ID:             scaleID,
			Name:           scale.Name,
			Description:    scale.Description,
			ScaleType:      "custom",
			Values:         scale.Values,
			DefaultScale:   scale.DefaultScale,
			OrganizationID: orgID,
			IsPublic:       false,
		}

		updatedScale, err := s.PokerDataSvc.UpdateOrganizationEstimationScale(ctx, &es)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleEstimationScaleUpdate error", zap.Error(err),
				zap.String("scale_id", scaleID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, updatedScale, nil)
	}
}
