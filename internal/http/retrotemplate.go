package http

import (
	"encoding/json"
	"io"
	"net/http"

	"go.uber.org/zap"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
)

type retroTemplateFormatRequestBody struct {
	Columns []struct {
		Name  string `json:"name" validate:"required,alpha,lowercase,min=1,max=16"`
		Label string `json:"label" validate:"required"`
		Color string `json:"color" validate:"omitempty,oneof=red green blue yellow purple orange teal"`
		Icon  string `json:"icon" validate:"omitempty,oneof=smiley frown angry question"`
	} `json:"columns" validate:"required,min=2,max=5,dive,required"`
}

type retroTemplateRequestBody struct {
	Name            string                         `json:"name" validate:"required"`
	Description     string                         `json:"description"`
	Format          retroTemplateFormatRequestBody `json:"format" validate:"required"`
	IsPublic        bool                           `json:"isPublic"`
	DefaultTemplate bool                           `json:"defaultTemplate"`
	OrganizationID  *string                        `json:"organizationId"`
	TeamID          *string                        `json:"teamId"`
}

// handleGetRetroTemplates gets a list of retro templates
//
//	@Summary		Get Retro Templates
//	@Description	get list of retro templates
//	@Tags			retroTemplate
//	@Produce		json
//	@Param			limit	query	int	false	"Max number of results to return"
//	@Param			offset	query	int	false	"Starting point to return rows from, should be multiplied by limit or 0"
//	@Success		200		object	standardJsonResponse{data=[]thunderdome.RetroTemplate}
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/admin/retro-templates [get]
func (s *Service) handleGetRetroTemplates() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID, _ := ctx.Value(contextKeyUserID).(*string)
		limit, offset := getLimitOffsetFromRequest(r)
		templates, count, err := s.RetroTemplateDataSvc.ListTemplates(ctx, limit, offset)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleGetRetroTemplates error", zap.Error(err),
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

		s.Success(w, r, http.StatusOK, templates, meta)
	}
}

// handleRetroTemplateCreate creates a new retro template
//
//	@Summary		Create Retro Template
//	@Description	Creates a retro template
//	@Tags			retroTemplate
//	@Produce		json
//	@Param			template	body	retroTemplateRequestBody	true	"new retro template object"
//	@Success		200			object	standardJsonResponse{data=thunderdome.RetroTemplate}
//	@Failure		400			object	standardJsonResponse{}
//	@Failure		500			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/admin/retro-templates [post]
func (s *Service) handleRetroTemplateCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		var template = retroTemplateRequestBody{}
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		jsonErr := json.Unmarshal(body, &template)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(template)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		newTemplate := &thunderdome.RetroTemplate{
			Name:            template.Name,
			Description:     template.Description,
			Format:          retroTemplateBuildFormatFromRequest(template.Format),
			IsPublic:        template.IsPublic,
			DefaultTemplate: template.DefaultTemplate,
			CreatedBy:       sessionUserID,
			OrganizationID:  template.OrganizationID,
			TeamID:          template.TeamID,
		}

		err := s.RetroTemplateDataSvc.CreateTemplate(ctx, newTemplate)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleRetroTemplateCreate error", zap.Error(err),
				zap.String("template_name", template.Name),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, newTemplate, nil)
	}
}

// handleRetroTemplateUpdate updates a retro template
//
//	@Summary		Update Retro Template
//	@Description	Updates a Retro Template
//	@Tags			retroTemplate
//	@Produce		json
//	@Param			templateId	path	string						true	"the retro template ID to update"
//	@Param			template	body	retroTemplateRequestBody	true	"retro template object to update"
//	@Success		200			object	standardJsonResponse{data=thunderdome.RetroTemplate}
//	@Failure		400			object	standardJsonResponse{}
//	@Failure		500			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/admin/retro-templates/{templateId} [put]
func (s *Service) handleRetroTemplateUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		templateID := r.PathValue("templateId")
		idErr := validate.Var(templateID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		var template = retroTemplateRequestBody{}
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		jsonErr := json.Unmarshal(body, &template)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(template)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		updatedTemplate := &thunderdome.RetroTemplate{
			ID:              templateID,
			Name:            template.Name,
			Description:     template.Description,
			Format:          retroTemplateBuildFormatFromRequest(template.Format),
			IsPublic:        template.IsPublic,
			DefaultTemplate: template.DefaultTemplate,
			OrganizationID:  template.OrganizationID,
			TeamID:          template.TeamID,
		}

		err := s.RetroTemplateDataSvc.UpdateTemplate(ctx, updatedTemplate)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleRetroTemplateUpdate error", zap.Error(err), zap.String("template_id", templateID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, updatedTemplate, nil)
	}
}

// handleRetroTemplateDelete handles deleting a retro template
//
//	@Summary		Delete Retro Template
//	@Description	Deletes a Retro Template
//	@Tags			retroTemplate
//	@Produce		json
//	@Param			templateId	path	string	true	"the retro template ID to delete"
//	@Success		200			object	standardJsonResponse{}
//	@Failure		400			object	standardJsonResponse{}
//	@Failure		500			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/admin/retro-templates/{templateId} [delete]
func (s *Service) handleRetroTemplateDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		templateID := r.PathValue("templateId")
		idErr := validate.Var(templateID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		err := s.RetroTemplateDataSvc.DeleteTemplate(ctx, templateID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleRetroTemplateDelete error", zap.Error(err), zap.String("template_id", templateID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleGetRetroTemplateByID gets a specific retro template by ID
//
//	@Summary		Get Retro Template by ID
//	@Description	get a specific retro template by its ID
//	@Tags			retroTemplate
//	@Produce		json
//	@Param			templateId	path	string	true	"the retro template ID"
//	@Success		200			object	standardJsonResponse{data=thunderdome.RetroTemplate}
//	@Failure		404			object	standardJsonResponse{}
//	@Failure		500			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/admin/retro-templates/{templateId} [get]
func (s *Service) handleGetRetroTemplateByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID, _ := ctx.Value(contextKeyUserID).(*string)

		templateID := r.PathValue("templateId")
		idErr := validate.Var(templateID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}

		template, err := s.RetroTemplateDataSvc.GetTemplateByID(ctx, templateID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleGetRetroTemplateByID error", zap.Error(err), zap.String("template_id", templateID),
				zap.Stringp("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		if template == nil {
			s.Failure(w, r, http.StatusNotFound, Errorf(ENOTFOUND, "Retro template not found"))
			return
		}

		s.Success(w, r, http.StatusOK, template, nil)
	}
}

// handleGetPublicRetroTemplates gets a list of public retro templates
//
//	@Summary		Get Public Retro Templates
//	@Description	get list of public retro templates
//	@Tags			retroTemplate
//	@Produce		json
//	@Success		200	{object}	standardJsonResponse{data=[]thunderdome.RetroTemplate}
//	@Failure		500	{object}	standardJsonResponse{}
//	@Router			/retro-templates/public [get]
func (s *Service) handleGetPublicRetroTemplates() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID, _ := ctx.Value(contextKeyUserID).(*string)

		templates, err := s.RetroTemplateDataSvc.GetPublicTemplates(ctx)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleGetPublicRetroTemplates error", zap.Error(err),
				zap.Stringp("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, templates, nil)
	}
}

type privateRetroTemplateRequestBody struct {
	Name            string                         `json:"name" validate:"required"`
	Description     string                         `json:"description"`
	Format          retroTemplateFormatRequestBody `json:"format" validate:"required"`
	DefaultTemplate bool                           `json:"defaultTemplate"`
}

// handleGetOrganizationRetroTemplates gets a list of retro templates for an organization
//
//	@Summary		Get Organization Retro Templates
//	@Description	get list of retro templates for an organization
//	@Tags			retroTemplate
//	@Produce		json
//	@Param			orgId	path		string	true	"Organization ID"
//	@Success		200		{object}	standardJsonResponse{data=[]thunderdome.RetroTemplate}
//	@Failure		400		{object}	standardJsonResponse{}
//	@Failure		500		{object}	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/retro-templates [get]
func (s *Service) handleGetOrganizationRetroTemplates() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID, _ := ctx.Value(contextKeyUserID).(*string)

		organizationID := r.PathValue("orgId")
		orgIDErr := validate.Var(organizationID, "required,uuid")
		if orgIDErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, orgIDErr.Error()))
			return
		}

		// Validate organizationID
		if err := validate.Var(organizationID, "required,uuid"); err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "Invalid organization ID"))
			return
		}

		templates, err := s.RetroTemplateDataSvc.GetTemplatesByOrganization(ctx, organizationID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleGetOrganizationRetroTemplates error", zap.Error(err),
				zap.String("organization_id", organizationID),
				zap.Stringp("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, templates, nil)
	}
}

// handleGetTeamRetroTemplates gets a list of retro templates for a team
//
//	@Summary		Get Team Retro Templates
//	@Description	get list of retro templates for a team
//	@Tags			retroTemplate
//	@Produce		json
//	@Param			teamId	path		string	true	"Team ID"
//	@Success		200		{object}	standardJsonResponse{data=[]thunderdome.RetroTemplate}
//	@Failure		400		{object}	standardJsonResponse{}
//	@Failure		500		{object}	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/retro-templates [get]
func (s *Service) handleGetTeamRetroTemplates() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID, _ := ctx.Value(contextKeyUserID).(*string)

		teamID := r.PathValue("teamId")
		teamIDErr := validate.Var(teamID, "required,uuid")
		if teamIDErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, teamIDErr.Error()))
			return
		}

		// Validate teamID
		if err := validate.Var(teamID, "required,uuid"); err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "Invalid team ID"))
			return
		}

		templates, err := s.RetroTemplateDataSvc.GetTemplatesByTeam(ctx, teamID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleGetTeamRetroTemplates error", zap.Error(err),
				zap.String("team_id", teamID),
				zap.Stringp("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, templates, nil)
	}
}

// handleTeamRetroTemplateCreate creates a new team retro template
//
//	@Summary		Create Team Retro Template
//	@Description	Creates a Team retro template
//	@Tags			retroTemplate
//	@Produce		json
//	@Param			teamId		path	string							true	"the team ID"
//	@Param			template	body	privateRetroTemplateRequestBody	true	"new retro template object"
//	@Success		200			object	standardJsonResponse{data=thunderdome.RetroTemplate}
//	@Failure		400			object	standardJsonResponse{}
//	@Failure		500			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/retro-templates [post]
func (s *Service) handleTeamRetroTemplateCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		teamID := r.PathValue("teamId")
		teamIDErr := validate.Var(teamID, "required,uuid")
		if teamIDErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, teamIDErr.Error()))
			return
		}
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		var template = privateRetroTemplateRequestBody{}
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		jsonErr := json.Unmarshal(body, &template)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(template)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		newTemplate := &thunderdome.RetroTemplate{
			Name:            template.Name,
			Description:     template.Description,
			Format:          retroTemplateBuildFormatFromRequest(template.Format),
			IsPublic:        false,
			DefaultTemplate: template.DefaultTemplate,
			CreatedBy:       sessionUserID,
			TeamID:          &teamID,
		}

		err := s.RetroTemplateDataSvc.CreateTemplate(ctx, newTemplate)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleTeamRetroTemplateCreate error", zap.Error(err),
				zap.String("template_name", template.Name),
				zap.String("team_id", teamID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, newTemplate, nil)
	}
}

// handleOrganizationRetroTemplateCreate creates a new organization retro template
//
//	@Summary		Create Organization Retro Template
//	@Description	Creates an Organization retro template
//	@Tags			retroTemplate
//	@Produce		json
//	@Param			orgId		path	string							true	"the organization ID"
//	@Param			template	body	privateRetroTemplateRequestBody	true	"new retro template object"
//	@Success		200			object	standardJsonResponse{data=thunderdome.RetroTemplate}
//	@Failure		400			object	standardJsonResponse{}
//	@Failure		500			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/retro-templates [post]
func (s *Service) handleOrganizationRetroTemplateCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		orgID := r.PathValue("orgId")
		orgIDErr := validate.Var(orgID, "required,uuid")
		if orgIDErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, orgIDErr.Error()))
			return
		}
		sessionUserID := ctx.Value(contextKeyUserID).(string)
		var template = privateRetroTemplateRequestBody{}
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		jsonErr := json.Unmarshal(body, &template)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(template)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		newTemplate := &thunderdome.RetroTemplate{
			Name:            template.Name,
			Description:     template.Description,
			Format:          retroTemplateBuildFormatFromRequest(template.Format),
			IsPublic:        false,
			DefaultTemplate: template.DefaultTemplate,
			CreatedBy:       sessionUserID,
			OrganizationID:  &orgID,
		}

		err := s.RetroTemplateDataSvc.CreateTemplate(ctx, newTemplate)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleOrganizationRetroTemplateCreate error", zap.Error(err),
				zap.String("template_name", template.Name),
				zap.String("organization_id", orgID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, newTemplate, nil)
	}
}

// handleTeamRetroTemplateUpdate updates a team retro template
//
//	@Summary		Update Team Retro Template
//	@Description	Updates a Team Retro Template
//	@Tags			retroTemplate
//	@Produce		json
//	@Param			teamId		path	string							true	"the team ID"
//	@Param			templateId	path	string							true	"the retro template ID to update"
//	@Param			template	body	privateRetroTemplateRequestBody	true	"retro template object to update"
//	@Success		200			object	standardJsonResponse{data=thunderdome.RetroTemplate}
//	@Failure		400			object	standardJsonResponse{}
//	@Failure		500			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/retro-templates/{templateId} [put]
func (s *Service) handleTeamRetroTemplateUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		templateID := r.PathValue("templateId")
		idErr := validate.Var(templateID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		teamID := r.PathValue("teamId")
		teamIDErr := validate.Var(teamID, "required,uuid")
		if teamIDErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, teamIDErr.Error()))
			return
		}

		var template = privateRetroTemplateRequestBody{}
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		jsonErr := json.Unmarshal(body, &template)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(template)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		updatedTemplate := &thunderdome.RetroTemplate{
			ID:              templateID,
			Name:            template.Name,
			Description:     template.Description,
			Format:          retroTemplateBuildFormatFromRequest(template.Format),
			DefaultTemplate: template.DefaultTemplate,
			TeamID:          &teamID,
		}

		err := s.RetroTemplateDataSvc.UpdateTeamTemplate(ctx, updatedTemplate)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleTeamRetroTemplateUpdate error", zap.Error(err),
				zap.String("template_id", templateID),
				zap.String("team_id", teamID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, updatedTemplate, nil)
	}
}

// handleOrganizationRetroTemplateUpdate updates an organization retro template
//
//	@Summary		Update Organization Retro Template
//	@Description	Updates an Organization Retro Template
//	@Tags			retroTemplate
//	@Produce		json
//	@Param			orgId		path	string							true	"the organization ID"
//	@Param			templateId	path	string							true	"the retro template ID to update"
//	@Param			template	body	privateRetroTemplateRequestBody	true	"retro template object to update"
//	@Success		200			object	standardJsonResponse{data=thunderdome.RetroTemplate}
//	@Failure		400			object	standardJsonResponse{}
//	@Failure		500			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organization/{orgId}/retro-templates/{templateId} [put]
func (s *Service) handleOrganizationRetroTemplateUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		templateID := r.PathValue("templateId")
		idErr := validate.Var(templateID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		orgID := r.PathValue("orgId")
		orgIDErr := validate.Var(orgID, "required,uuid")
		if orgIDErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, orgIDErr.Error()))
			return
		}

		var template = privateRetroTemplateRequestBody{}
		body, bodyErr := io.ReadAll(r.Body)
		if bodyErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, bodyErr.Error()))
			return
		}

		jsonErr := json.Unmarshal(body, &template)
		if jsonErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, jsonErr.Error()))
			return
		}

		inputErr := validate.Struct(template)
		if inputErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, inputErr.Error()))
			return
		}

		updatedTemplate := &thunderdome.RetroTemplate{
			ID:              templateID,
			Name:            template.Name,
			Description:     template.Description,
			Format:          retroTemplateBuildFormatFromRequest(template.Format),
			DefaultTemplate: template.DefaultTemplate,
			OrganizationID:  &orgID,
		}

		err := s.RetroTemplateDataSvc.UpdateOrganizationTemplate(ctx, updatedTemplate)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleOrganizationRetroTemplateUpdate error", zap.Error(err),
				zap.String("template_id", templateID),
				zap.String("organization_id", orgID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, updatedTemplate, nil)
	}
}

// handleOrganizationRetroTemplateDelete handles deleting an organization retro template
//
//	@Summary		Delete Organization Retro Template
//	@Description	Deletes an Organization Retro Template
//	@Tags			retroTemplate
//	@Produce		json
//	@Param			orgId		path	string	true	"the Organization ID"
//	@Param			templateId	path	string	true	"the retro template ID to delete"
//	@Success		200			object	standardJsonResponse{}
//	@Failure		400			object	standardJsonResponse{}
//	@Failure		500			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/retro-templates/{templateId} [delete]
func (s *Service) handleOrganizationRetroTemplateDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		templateID := r.PathValue("templateId")
		idErr := validate.Var(templateID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		organizationID := r.PathValue("orgId")
		oidErr := validate.Var(organizationID, "required,uuid")
		if oidErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, oidErr.Error()))
			return
		}

		err := s.RetroTemplateDataSvc.DeleteOrganizationTemplate(ctx, organizationID, templateID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleOrganizationRetroTemplateDelete error", zap.Error(err),
				zap.String("template_id", templateID),
				zap.String("organization_id", organizationID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleTeamRetroTemplateDelete handles deleting a team retro template
//
//	@Summary		Delete Team Retro Template
//	@Description	Deletes a Team Retro Template
//	@Tags			retroTemplate
//	@Produce		json
//	@Param			teamId		path	string	true	"the team ID"
//	@Param			templateId	path	string	true	"the retro template ID to delete"
//	@Success		200			object	standardJsonResponse{}
//	@Failure		400			object	standardJsonResponse{}
//	@Failure		500			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/retro-templates/{templateId} [delete]
func (s *Service) handleTeamRetroTemplateDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		templateID := r.PathValue("templateId")
		idErr := validate.Var(templateID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		teamID := r.PathValue("teamId")
		teamIDErr := validate.Var(teamID, "required,uuid")
		if teamIDErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, teamIDErr.Error()))
			return
		}

		err := s.RetroTemplateDataSvc.DeleteTeamTemplate(ctx, teamID, templateID)
		if err != nil {
			s.Logger.Ctx(ctx).Error("handleTeamRetroTemplateDelete error", zap.Error(err),
				zap.String("template_id", templateID),
				zap.String("team_id", teamID),
				zap.String("session_user_id", sessionUserID))
			s.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}
