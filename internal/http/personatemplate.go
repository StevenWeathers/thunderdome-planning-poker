package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"go.uber.org/zap"
)

type personaTemplatePersonaRequestBody struct {
	Name        string `json:"name" validate:"required,max=256"`
	Role        string `json:"role" validate:"max=256"`
	Description string `json:"description" validate:"max=2000"`
}

type personaTemplateRequestBody struct {
	Name        string                              `json:"name" validate:"required,max=255"`
	Description string                              `json:"description" validate:"max=2000"`
	Personas    []personaTemplatePersonaRequestBody `json:"personas" validate:"required,min=1,max=32,unique=Name,dive"`
}

// normalize trims surrounding whitespace so that blank values fail validation
func (b *personaTemplateRequestBody) normalize() {
	b.Name = strings.TrimSpace(b.Name)
	b.Description = strings.TrimSpace(b.Description)
	for i := range b.Personas {
		b.Personas[i].Name = strings.TrimSpace(b.Personas[i].Name)
		b.Personas[i].Role = strings.TrimSpace(b.Personas[i].Role)
		b.Personas[i].Description = strings.TrimSpace(b.Personas[i].Description)
	}
}

// personaTemplateScope reads the owning organization or team of a template from the route,
// team routes nested under an organization belong to the team
func personaTemplateScope(r *http.Request, pathParam string) (organizationID string, teamID string, err error) {
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

func (s *Service) personaTemplateFailure(w http.ResponseWriter, r *http.Request, action string, err error) {
	if errors.Is(err, thunderdome.ErrPersonaTemplateNotFound) {
		s.Failure(w, r, http.StatusNotFound, Errorf(ENOTFOUND, err.Error()))
		return
	}

	s.Logger.Ctx(r.Context()).Error(action+" error", zap.Error(err),
		zap.String("organization_id", r.PathValue("orgId")),
		zap.String("team_id", r.PathValue("teamId")),
		zap.String("template_id", r.PathValue("templateId")),
		zap.String("storyboard_id", r.PathValue("storyboardId")))
	s.Failure(w, r, http.StatusInternalServerError, err)
}

func (s *Service) listPersonaTemplates(pathParam string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		organizationID, teamID, err := personaTemplateScope(r, pathParam)
		if err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}

		templates, err := s.PersonaTemplateDataSvc.GetPersonaTemplates(r.Context(), organizationID, teamID)
		if err != nil {
			s.personaTemplateFailure(w, r, "listPersonaTemplates", err)
			return
		}

		s.Success(w, r, http.StatusOK, templates, nil)
	}
}

// savePersonaTemplate creates a template, or updates it when update is true
func (s *Service) savePersonaTemplate(pathParam string, update bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		organizationID, teamID, err := personaTemplateScope(r, pathParam)
		if err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}

		template := &thunderdome.PersonaTemplate{}
		if organizationID != "" {
			template.OrganizationID = &organizationID
		}
		if teamID != "" {
			template.TeamID = &teamID
		}

		if update {
			template.ID = r.PathValue("templateId")
			if err := validate.Var(template.ID, "required,uuid"); err != nil {
				s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
				return
			}
		}

		request := personaTemplateRequestBody{}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}
		if err := json.Unmarshal(body, &request); err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}
		request.normalize()
		if err := validate.Struct(request); err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}

		template.Name = request.Name
		template.Description = request.Description
		template.Personas = make([]*thunderdome.PersonaTemplatePersona, 0, len(request.Personas))
		for _, p := range request.Personas {
			template.Personas = append(template.Personas, &thunderdome.PersonaTemplatePersona{
				Name:        p.Name,
				Role:        p.Role,
				Description: p.Description,
			})
		}

		if update {
			err = s.PersonaTemplateDataSvc.UpdatePersonaTemplate(ctx, template)
		} else {
			template.CreatedBy = sessionUserID
			err = s.PersonaTemplateDataSvc.CreatePersonaTemplate(ctx, template)
		}
		if err != nil {
			s.personaTemplateFailure(w, r, "savePersonaTemplate", err)
			return
		}

		s.Success(w, r, http.StatusOK, template, nil)
	}
}

func (s *Service) deletePersonaTemplate(pathParam string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		organizationID, teamID, err := personaTemplateScope(r, pathParam)
		if err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}

		templateID := r.PathValue("templateId")
		if err := validate.Var(templateID, "required,uuid"); err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}

		if err := s.PersonaTemplateDataSvc.DeletePersonaTemplate(r.Context(), organizationID, teamID, templateID); err != nil {
			s.personaTemplateFailure(w, r, "deletePersonaTemplate", err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleGetStoryboardPersonaTemplates gets the persona templates a facilitator can load into a storyboard
//
//	@Summary		Get Storyboard Persona Templates
//	@Description	get the persona templates that can be loaded into a storyboard (facilitator only)
//	@Tags			personaTemplate
//	@Produce		json
//	@Param			storyboardId	path		string	true	"Storyboard ID"
//	@Success		200				{object}	standardJsonResponse{data=[]thunderdome.PersonaTemplate}
//	@Failure		400				{object}	standardJsonResponse{}
//	@Failure		403				{object}	standardJsonResponse{}
//	@Failure		500				{object}	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/storyboards/{storyboardId}/persona-templates [get]
func (s *Service) handleGetStoryboardPersonaTemplates() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		storyboardID := r.PathValue("storyboardId")
		if err := validate.Var(storyboardID, "required,uuid"); err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}

		if err := s.StoryboardDataSvc.ConfirmStoryboardFacilitator(storyboardID, sessionUserID); err != nil {
			s.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "REQUIRES_STORYBOARD_FACILITATOR"))
			return
		}

		templates, err := s.PersonaTemplateDataSvc.GetStoryboardPersonaTemplates(ctx, storyboardID, s.Config.SubscriptionsEnabled)
		if err != nil {
			s.personaTemplateFailure(w, r, "handleGetStoryboardPersonaTemplates", err)
			return
		}

		s.Success(w, r, http.StatusOK, templates, nil)
	}
}

// handleGetAdminPersonaTemplates gets the global default persona templates
//
//	@Summary		Get Global Persona Templates
//	@Description	get list of global default persona templates
//	@Tags			personaTemplate
//	@Produce		json
//	@Success		200	{object}	standardJsonResponse{data=[]thunderdome.PersonaTemplate}
//	@Failure		500	{object}	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/admin/persona-templates [get]
func (s *Service) handleGetAdminPersonaTemplates() http.HandlerFunc {
	return s.listPersonaTemplates("")
}

// handleAdminPersonaTemplateCreate creates a global default persona template
//
//	@Summary		Create Global Persona Template
//	@Description	Creates a global default persona template
//	@Tags			personaTemplate
//	@Produce		json
//	@Param			template	body		personaTemplateRequestBody	true	"new persona template object"
//	@Success		200			{object}	standardJsonResponse{data=thunderdome.PersonaTemplate}
//	@Failure		400			{object}	standardJsonResponse{}
//	@Failure		500			{object}	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/admin/persona-templates [post]
func (s *Service) handleAdminPersonaTemplateCreate() http.HandlerFunc {
	return s.savePersonaTemplate("", false)
}

// handleAdminPersonaTemplateUpdate updates a global default persona template
//
//	@Summary		Update Global Persona Template
//	@Description	Updates a global default persona template
//	@Tags			personaTemplate
//	@Produce		json
//	@Param			templateId	path		string						true	"the persona template ID to update"
//	@Param			template	body		personaTemplateRequestBody	true	"persona template object to update"
//	@Success		200			{object}	standardJsonResponse{data=thunderdome.PersonaTemplate}
//	@Failure		400			{object}	standardJsonResponse{}
//	@Failure		404			{object}	standardJsonResponse{}
//	@Failure		500			{object}	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/admin/persona-templates/{templateId} [put]
func (s *Service) handleAdminPersonaTemplateUpdate() http.HandlerFunc {
	return s.savePersonaTemplate("", true)
}

// handleAdminPersonaTemplateDelete deletes a global default persona template
//
//	@Summary		Delete Global Persona Template
//	@Description	Deletes a global default persona template
//	@Tags			personaTemplate
//	@Produce		json
//	@Param			templateId	path		string	true	"the persona template ID to delete"
//	@Success		200			{object}	standardJsonResponse{}
//	@Failure		400			{object}	standardJsonResponse{}
//	@Failure		404			{object}	standardJsonResponse{}
//	@Failure		500			{object}	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/admin/persona-templates/{templateId} [delete]
func (s *Service) handleAdminPersonaTemplateDelete() http.HandlerFunc {
	return s.deletePersonaTemplate("")
}

// handleGetOrganizationPersonaTemplates gets the persona templates of an organization
//
//	@Summary		Get Organization Persona Templates
//	@Description	get list of persona templates for an organization
//	@Tags			personaTemplate
//	@Produce		json
//	@Param			orgId	path		string	true	"Organization ID"
//	@Success		200		{object}	standardJsonResponse{data=[]thunderdome.PersonaTemplate}
//	@Failure		400		{object}	standardJsonResponse{}
//	@Failure		500		{object}	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/persona-templates [get]
func (s *Service) handleGetOrganizationPersonaTemplates() http.HandlerFunc {
	return s.listPersonaTemplates("orgId")
}

// handleOrganizationPersonaTemplateCreate creates an organization persona template
//
//	@Summary		Create Organization Persona Template
//	@Description	Creates an organization persona template
//	@Tags			personaTemplate
//	@Produce		json
//	@Param			orgId		path		string						true	"the organization ID"
//	@Param			template	body		personaTemplateRequestBody	true	"new persona template object"
//	@Success		200			{object}	standardJsonResponse{data=thunderdome.PersonaTemplate}
//	@Failure		400			{object}	standardJsonResponse{}
//	@Failure		500			{object}	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/persona-templates [post]
func (s *Service) handleOrganizationPersonaTemplateCreate() http.HandlerFunc {
	return s.savePersonaTemplate("orgId", false)
}

// handleOrganizationPersonaTemplateUpdate updates an organization persona template
//
//	@Summary		Update Organization Persona Template
//	@Description	Updates an organization persona template
//	@Tags			personaTemplate
//	@Produce		json
//	@Param			orgId		path		string						true	"the organization ID"
//	@Param			templateId	path		string						true	"the persona template ID to update"
//	@Param			template	body		personaTemplateRequestBody	true	"persona template object to update"
//	@Success		200			{object}	standardJsonResponse{data=thunderdome.PersonaTemplate}
//	@Failure		400			{object}	standardJsonResponse{}
//	@Failure		404			{object}	standardJsonResponse{}
//	@Failure		500			{object}	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/persona-templates/{templateId} [put]
func (s *Service) handleOrganizationPersonaTemplateUpdate() http.HandlerFunc {
	return s.savePersonaTemplate("orgId", true)
}

// handleOrganizationPersonaTemplateDelete deletes an organization persona template
//
//	@Summary		Delete Organization Persona Template
//	@Description	Deletes an organization persona template
//	@Tags			personaTemplate
//	@Produce		json
//	@Param			orgId		path		string	true	"the organization ID"
//	@Param			templateId	path		string	true	"the persona template ID to delete"
//	@Success		200			{object}	standardJsonResponse{}
//	@Failure		400			{object}	standardJsonResponse{}
//	@Failure		404			{object}	standardJsonResponse{}
//	@Failure		500			{object}	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/persona-templates/{templateId} [delete]
func (s *Service) handleOrganizationPersonaTemplateDelete() http.HandlerFunc {
	return s.deletePersonaTemplate("orgId")
}

// handleGetTeamPersonaTemplates gets the persona templates of a team
//
//	@Summary		Get Team Persona Templates
//	@Description	get list of persona templates for a team
//	@Tags			personaTemplate
//	@Produce		json
//	@Param			teamId	path		string	true	"Team ID"
//	@Success		200		{object}	standardJsonResponse{data=[]thunderdome.PersonaTemplate}
//	@Failure		400		{object}	standardJsonResponse{}
//	@Failure		500		{object}	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/persona-templates [get]
//	@Router			/organizations/{orgId}/teams/{teamId}/persona-templates [get]
//	@Router			/organizations/{orgId}/departments/{departmentId}/teams/{teamId}/persona-templates [get]
func (s *Service) handleGetTeamPersonaTemplates() http.HandlerFunc {
	return s.listPersonaTemplates("teamId")
}

// handleTeamPersonaTemplateCreate creates a team persona template
//
//	@Summary		Create Team Persona Template
//	@Description	Creates a team persona template
//	@Tags			personaTemplate
//	@Produce		json
//	@Param			teamId		path		string						true	"the team ID"
//	@Param			template	body		personaTemplateRequestBody	true	"new persona template object"
//	@Success		200			{object}	standardJsonResponse{data=thunderdome.PersonaTemplate}
//	@Failure		400			{object}	standardJsonResponse{}
//	@Failure		500			{object}	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/persona-templates [post]
//	@Router			/organizations/{orgId}/teams/{teamId}/persona-templates [post]
//	@Router			/organizations/{orgId}/departments/{departmentId}/teams/{teamId}/persona-templates [post]
func (s *Service) handleTeamPersonaTemplateCreate() http.HandlerFunc {
	return s.savePersonaTemplate("teamId", false)
}

// handleTeamPersonaTemplateUpdate updates a team persona template
//
//	@Summary		Update Team Persona Template
//	@Description	Updates a team persona template
//	@Tags			personaTemplate
//	@Produce		json
//	@Param			teamId		path		string						true	"the team ID"
//	@Param			templateId	path		string						true	"the persona template ID to update"
//	@Param			template	body		personaTemplateRequestBody	true	"persona template object to update"
//	@Success		200			{object}	standardJsonResponse{data=thunderdome.PersonaTemplate}
//	@Failure		400			{object}	standardJsonResponse{}
//	@Failure		404			{object}	standardJsonResponse{}
//	@Failure		500			{object}	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/persona-templates/{templateId} [put]
//	@Router			/organizations/{orgId}/teams/{teamId}/persona-templates/{templateId} [put]
//	@Router			/organizations/{orgId}/departments/{departmentId}/teams/{teamId}/persona-templates/{templateId} [put]
func (s *Service) handleTeamPersonaTemplateUpdate() http.HandlerFunc {
	return s.savePersonaTemplate("teamId", true)
}

// handleTeamPersonaTemplateDelete deletes a team persona template
//
//	@Summary		Delete Team Persona Template
//	@Description	Deletes a team persona template
//	@Tags			personaTemplate
//	@Produce		json
//	@Param			teamId		path		string	true	"the team ID"
//	@Param			templateId	path		string	true	"the persona template ID to delete"
//	@Success		200			{object}	standardJsonResponse{}
//	@Failure		400			{object}	standardJsonResponse{}
//	@Failure		404			{object}	standardJsonResponse{}
//	@Failure		500			{object}	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/persona-templates/{templateId} [delete]
//	@Router			/organizations/{orgId}/teams/{teamId}/persona-templates/{templateId} [delete]
//	@Router			/organizations/{orgId}/departments/{departmentId}/teams/{teamId}/persona-templates/{templateId} [delete]
func (s *Service) handleTeamPersonaTemplateDelete() http.HandlerFunc {
	return s.deletePersonaTemplate("teamId")
}
