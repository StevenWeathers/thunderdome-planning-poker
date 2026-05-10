package http

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"go.uber.org/zap"
)

type colorLegendTemplateColorRequestBody struct {
	Color  string `json:"color" validate:"required,lowercase,min=1,max=32"`
	Legend string `json:"legend" validate:"max=255"`
}

type colorLegendTemplateRequestBody struct {
	Name        string                                `json:"name" validate:"required"`
	Description string                                `json:"description"`
	ColorLegend []colorLegendTemplateColorRequestBody `json:"colorLegend" validate:"required,min=1,max=32,dive"`
}

func colorLegendTemplateBuildLegendFromRequest(requestColors []colorLegendTemplateColorRequestBody) []*thunderdome.Color {
	colors := make([]*thunderdome.Color, 0, len(requestColors))
	for _, color := range requestColors {
		colors = append(colors, &thunderdome.Color{
			Color:  color.Color,
			Legend: color.Legend,
		})
	}

	return colors
}

// handleGetOrganizationColorLegendTemplates gets a list of color legend templates for an organization
//
//	@Summary		Get Organization Color Legend Templates
//	@Description	get list of color legend templates for an organization
//	@Tags			colorLegendTemplate
//	@Produce		json
//	@Param			orgId	path		string	true	"Organization ID"
//	@Success		200		{object}	standardJsonResponse{data=[]thunderdome.ColorLegendTemplate}
//	@Failure		400		{object}	standardJsonResponse{}
//	@Failure		500		{object}	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/color-legend-templates [get]
func (s *Service) handleGetOrganizationColorLegendTemplates() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID, _ := ctx.Value(contextKeyUserID).(*string)

		organizationID := r.PathValue("orgId")
		if err := validate.Var(organizationID, "required,uuid"); err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}

		templates, err := s.ColorLegendTemplateDataSvc.GetColorLegendTemplatesByOrganization(ctx, organizationID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleGetOrganizationColorLegendTemplates error", zap.Error(err),
				zap.String("organization_id", organizationID),
				zap.Stringp("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, templates, nil)
	}
}

// handleGetTeamColorLegendTemplates gets a list of color legend templates for a team
//
//	@Summary		Get Team Color Legend Templates
//	@Description	get list of color legend templates for a team
//	@Tags			colorLegendTemplate
//	@Produce		json
//	@Param			teamId	path		string	true	"Team ID"
//	@Success		200		{object}	standardJsonResponse{data=[]thunderdome.ColorLegendTemplate}
//	@Failure		400		{object}	standardJsonResponse{}
//	@Failure		500		{object}	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/color-legend-templates [get]
//	@Router			/organizations/{orgId}/teams/{teamId}/color-legend-templates [get]
//	@Router			/organizations/{orgId}/departments/{departmentId}/teams/{teamId}/color-legend-templates [get]
func (s *Service) handleGetTeamColorLegendTemplates() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID, _ := ctx.Value(contextKeyUserID).(*string)

		teamID := r.PathValue("teamId")
		if err := validate.Var(teamID, "required,uuid"); err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}

		templates, err := s.ColorLegendTemplateDataSvc.GetColorLegendTemplatesByTeam(ctx, teamID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleGetTeamColorLegendTemplates error", zap.Error(err),
				zap.String("team_id", teamID),
				zap.Stringp("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, templates, nil)
	}
}

// handleOrganizationColorLegendTemplateCreate creates a new organization color legend template
//
//	@Summary		Create Organization Color Legend Template
//	@Description	Creates an organization color legend template
//	@Tags			colorLegendTemplate
//	@Produce		json
//	@Param			orgId		path	string							true	"the organization ID"
//	@Param			template	body	colorLegendTemplateRequestBody	true	"new color legend template object"
//	@Success		200			object	standardJsonResponse{data=thunderdome.ColorLegendTemplate}
//	@Failure		400			object	standardJsonResponse{}
//	@Failure		500			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/color-legend-templates [post]
func (s *Service) handleOrganizationColorLegendTemplateCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		organizationID := r.PathValue("orgId")
		if err := validate.Var(organizationID, "required,uuid"); err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}

		sessionUserID := ctx.Value(contextKeyUserID).(string)
		request := colorLegendTemplateRequestBody{}
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
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

		template := &thunderdome.ColorLegendTemplate{
			Name:           request.Name,
			Description:    request.Description,
			ColorLegend:    colorLegendTemplateBuildLegendFromRequest(request.ColorLegend),
			CreatedBy:      sessionUserID,
			OrganizationID: &organizationID,
		}

		createdTemplate, err := s.ColorLegendTemplateDataSvc.CreateColorLegendTemplate(ctx, template)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleOrganizationColorLegendTemplateCreate error", zap.Error(err),
				zap.String("organization_id", organizationID),
				zap.String("session_user_id", sessionUserID),
				zap.String("template_name", request.Name))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, createdTemplate, nil)
	}
}

// handleTeamColorLegendTemplateCreate creates a new team color legend template
//
//	@Summary		Create Team Color Legend Template
//	@Description	Creates a team color legend template
//	@Tags			colorLegendTemplate
//	@Produce		json
//	@Param			teamId		path	string							true	"the team ID"
//	@Param			template	body	colorLegendTemplateRequestBody	true	"new color legend template object"
//	@Success		200			object	standardJsonResponse{data=thunderdome.ColorLegendTemplate}
//	@Failure		400			object	standardJsonResponse{}
//	@Failure		500			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/color-legend-templates [post]
//	@Router			/organizations/{orgId}/teams/{teamId}/color-legend-templates [post]
//	@Router			/organizations/{orgId}/departments/{departmentId}/teams/{teamId}/color-legend-templates [post]
func (s *Service) handleTeamColorLegendTemplateCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		teamID := r.PathValue("teamId")
		if err := validate.Var(teamID, "required,uuid"); err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}

		sessionUserID := ctx.Value(contextKeyUserID).(string)
		request := colorLegendTemplateRequestBody{}
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
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

		template := &thunderdome.ColorLegendTemplate{
			Name:        request.Name,
			Description: request.Description,
			ColorLegend: colorLegendTemplateBuildLegendFromRequest(request.ColorLegend),
			CreatedBy:   sessionUserID,
			TeamID:      &teamID,
		}

		createdTemplate, err := s.ColorLegendTemplateDataSvc.CreateColorLegendTemplate(ctx, template)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleTeamColorLegendTemplateCreate error", zap.Error(err),
				zap.String("team_id", teamID),
				zap.String("session_user_id", sessionUserID),
				zap.String("template_name", request.Name))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, createdTemplate, nil)
	}
}

// handleOrganizationColorLegendTemplateUpdate updates an organization color legend template
//
//	@Summary		Update Organization Color Legend Template
//	@Description	Updates an organization color legend template
//	@Tags			colorLegendTemplate
//	@Produce		json
//	@Param			orgId		path	string							true	"the organization ID"
//	@Param			templateId	path	string							true	"the color legend template ID to update"
//	@Param			template	body	colorLegendTemplateRequestBody	true	"color legend template object to update"
//	@Success		200			object	standardJsonResponse{data=thunderdome.ColorLegendTemplate}
//	@Failure		400			object	standardJsonResponse{}
//	@Failure		500			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/color-legend-templates/{templateId} [put]
func (s *Service) handleOrganizationColorLegendTemplateUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		organizationID := r.PathValue("orgId")
		templateID := r.PathValue("templateId")

		if err := validate.Var(organizationID, "required,uuid"); err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}
		if err := validate.Var(templateID, "required,uuid"); err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}

		request := colorLegendTemplateRequestBody{}
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
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

		updatedTemplate := &thunderdome.ColorLegendTemplate{
			ID:             templateID,
			Name:           request.Name,
			Description:    request.Description,
			ColorLegend:    colorLegendTemplateBuildLegendFromRequest(request.ColorLegend),
			OrganizationID: &organizationID,
		}

		template, err := s.ColorLegendTemplateDataSvc.UpdateOrganizationColorLegendTemplate(ctx, updatedTemplate)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleOrganizationColorLegendTemplateUpdate error", zap.Error(err),
				zap.String("organization_id", organizationID),
				zap.String("template_id", templateID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, template, nil)
	}
}

// handleTeamColorLegendTemplateUpdate updates a team color legend template
//
//	@Summary		Update Team Color Legend Template
//	@Description	Updates a team color legend template
//	@Tags			colorLegendTemplate
//	@Produce		json
//	@Param			teamId		path	string							true	"the team ID"
//	@Param			templateId	path	string							true	"the color legend template ID to update"
//	@Param			template	body	colorLegendTemplateRequestBody	true	"color legend template object to update"
//	@Success		200			object	standardJsonResponse{data=thunderdome.ColorLegendTemplate}
//	@Failure		400			object	standardJsonResponse{}
//	@Failure		500			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/color-legend-templates/{templateId} [put]
//	@Router			/organizations/{orgId}/teams/{teamId}/color-legend-templates/{templateId} [put]
//	@Router			/organizations/{orgId}/departments/{departmentId}/teams/{teamId}/color-legend-templates/{templateId} [put]
func (s *Service) handleTeamColorLegendTemplateUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		teamID := r.PathValue("teamId")
		templateID := r.PathValue("templateId")

		if err := validate.Var(teamID, "required,uuid"); err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}
		if err := validate.Var(templateID, "required,uuid"); err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}

		request := colorLegendTemplateRequestBody{}
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
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

		updatedTemplate := &thunderdome.ColorLegendTemplate{
			ID:          templateID,
			Name:        request.Name,
			Description: request.Description,
			ColorLegend: colorLegendTemplateBuildLegendFromRequest(request.ColorLegend),
			TeamID:      &teamID,
		}

		template, err := s.ColorLegendTemplateDataSvc.UpdateTeamColorLegendTemplate(ctx, updatedTemplate)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleTeamColorLegendTemplateUpdate error", zap.Error(err),
				zap.String("team_id", teamID),
				zap.String("template_id", templateID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, template, nil)
	}
}

// handleOrganizationColorLegendTemplateDelete deletes an organization color legend template
//
//	@Summary		Delete Organization Color Legend Template
//	@Description	Deletes an organization color legend template
//	@Tags			colorLegendTemplate
//	@Produce		json
//	@Param			orgId		path	string	true	"the organization ID"
//	@Param			templateId	path	string	true	"the color legend template ID to delete"
//	@Success		200			object	standardJsonResponse{}
//	@Failure		400			object	standardJsonResponse{}
//	@Failure		500			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/color-legend-templates/{templateId} [delete]
func (s *Service) handleOrganizationColorLegendTemplateDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		organizationID := r.PathValue("orgId")
		templateID := r.PathValue("templateId")

		if err := validate.Var(organizationID, "required,uuid"); err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}
		if err := validate.Var(templateID, "required,uuid"); err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}

		if err := s.ColorLegendTemplateDataSvc.DeleteOrganizationColorLegendTemplate(ctx, organizationID, templateID); err != nil {
			s.Logger.Ctx(ctx).Error("handleOrganizationColorLegendTemplateDelete error", zap.Error(err),
				zap.String("organization_id", organizationID),
				zap.String("template_id", templateID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleTeamColorLegendTemplateDelete deletes a team color legend template
//
//	@Summary		Delete Team Color Legend Template
//	@Description	Deletes a team color legend template
//	@Tags			colorLegendTemplate
//	@Produce		json
//	@Param			teamId		path	string	true	"the team ID"
//	@Param			templateId	path	string	true	"the color legend template ID to delete"
//	@Success		200			object	standardJsonResponse{}
//	@Failure		400			object	standardJsonResponse{}
//	@Failure		500			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/color-legend-templates/{templateId} [delete]
//	@Router			/organizations/{orgId}/teams/{teamId}/color-legend-templates/{templateId} [delete]
//	@Router			/organizations/{orgId}/departments/{departmentId}/teams/{teamId}/color-legend-templates/{templateId} [delete]
func (s *Service) handleTeamColorLegendTemplateDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		teamID := r.PathValue("teamId")
		templateID := r.PathValue("templateId")

		if err := validate.Var(teamID, "required,uuid"); err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}
		if err := validate.Var(templateID, "required,uuid"); err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}

		if err := s.ColorLegendTemplateDataSvc.DeleteTeamColorLegendTemplate(ctx, teamID, templateID); err != nil {
			s.Logger.Ctx(ctx).Error("handleTeamColorLegendTemplateDelete error", zap.Error(err),
				zap.String("team_id", teamID),
				zap.String("template_id", templateID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}
