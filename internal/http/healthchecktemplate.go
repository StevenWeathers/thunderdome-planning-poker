package http

import (
	"net/http"
	"strings"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
)

type healthCheckAreaRequestBody struct {
	Name      string `json:"name" validate:"required,max=128"`
	Healthy   string `json:"healthy" validate:"max=500"`
	Unhealthy string `json:"unhealthy" validate:"max=500"`
}

type healthCheckTemplateRequestBody struct {
	Name        string                       `json:"name" validate:"required,max=255"`
	Description string                       `json:"description" validate:"max=2000"`
	Areas       []healthCheckAreaRequestBody `json:"areas" validate:"required,min=1,max=20,unique=Name,dive"`
}

// normalize trims surrounding whitespace so that blank values fail validation
func (b *healthCheckTemplateRequestBody) normalize() {
	b.Name = strings.TrimSpace(b.Name)
	b.Description = strings.TrimSpace(b.Description)
	for i := range b.Areas {
		b.Areas[i].Name = strings.TrimSpace(b.Areas[i].Name)
		b.Areas[i].Healthy = strings.TrimSpace(b.Areas[i].Healthy)
		b.Areas[i].Unhealthy = strings.TrimSpace(b.Areas[i].Unhealthy)
	}
}

// listHealthCheckTemplates gets the templates of the scope named by pathParam, empty meaning the global defaults
func (s *Service) listHealthCheckTemplates(pathParam string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// persona template scoping applies to any template owned by an organization or team
		organizationID, teamID, err := personaTemplateScope(r, pathParam)
		if err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}

		templates, err := s.HealthCheckDataSvc.HealthCheckTemplateList(r.Context(), organizationID, teamID)
		if err != nil {
			s.healthCheckFailure(w, r, "listHealthCheckTemplates", err)
			return
		}

		s.Success(w, r, http.StatusOK, templates, nil)
	}
}

// saveHealthCheckTemplate creates a template, or updates it when update is true
func (s *Service) saveHealthCheckTemplate(pathParam string, update bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		sessionUserID := ctx.Value(contextKeyUserID).(string)

		organizationID, teamID, err := personaTemplateScope(r, pathParam)
		if err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}

		template := &thunderdome.HealthCheckTemplate{}
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

		request := healthCheckTemplateRequestBody{}
		if !s.decodeHealthCheckBody(w, r, &request) {
			return
		}
		request.normalize()
		if err := validate.Struct(request); err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}

		template.Name = request.Name
		template.Description = request.Description
		for _, area := range request.Areas {
			template.Areas = append(template.Areas, &thunderdome.HealthCheckArea{
				Name: area.Name, Healthy: area.Healthy, Unhealthy: area.Unhealthy,
			})
		}

		if update {
			err = s.HealthCheckDataSvc.HealthCheckTemplateUpdate(ctx, template)
		} else {
			template.CreatedBy = sessionUserID
			err = s.HealthCheckDataSvc.HealthCheckTemplateCreate(ctx, template)
		}
		if err != nil {
			s.healthCheckFailure(w, r, "saveHealthCheckTemplate", err)
			return
		}

		s.Success(w, r, http.StatusOK, template, nil)
	}
}

func (s *Service) deleteHealthCheckTemplate(pathParam string) http.HandlerFunc {
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

		if err := s.HealthCheckDataSvc.HealthCheckTemplateDelete(r.Context(), organizationID, teamID, templateID); err != nil {
			s.healthCheckFailure(w, r, "deleteHealthCheckTemplate", err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleGetAvailableHealthCheckTemplates gets the templates a team can create health checks from
//
//	@Summary		Get Available Health Check Templates
//	@Description	get the global, organization and team templates a team can create health checks from
//	@Tags			healthCheck
//	@Produce		json
//	@Param			teamId	path		string	true	"Team ID"
//	@Success		200		{object}	standardJsonResponse{data=[]thunderdome.HealthCheckTemplate}
//	@Failure		400		{object}	standardJsonResponse{}
//	@Failure		500		{object}	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/health-check-templates/available [get]
func (s *Service) handleGetAvailableHealthCheckTemplates() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		teamID := r.PathValue("teamId")
		if err := validate.Var(teamID, "required,uuid"); err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}

		templates, err := s.HealthCheckDataSvc.HealthCheckTemplatesAvailable(r.Context(), teamID)
		if err != nil {
			s.healthCheckFailure(w, r, "handleGetAvailableHealthCheckTemplates", err)
			return
		}

		s.Success(w, r, http.StatusOK, templates, nil)
	}
}

// handleGetAdminHealthCheckTemplates gets the global default health check templates
//
//	@Summary		Get Global Health Check Templates
//	@Description	Gets global default health check templates
//	@Tags			healthCheck
//	@Produce		json
//	@Success		200	{object}	standardJsonResponse{data=[]thunderdome.HealthCheckTemplate}
//	@Failure		400	{object}	standardJsonResponse{}
//	@Failure		500	{object}	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/admin/health-check-templates [get]
func (s *Service) handleGetAdminHealthCheckTemplates() http.HandlerFunc {
	return s.listHealthCheckTemplates("")
}

// handleAdminHealthCheckTemplateCreate creates a global default health check template
//
//	@Summary		Create Global Health Check Template
//	@Description	Creates global default health check template
//	@Tags			healthCheck
//	@Produce		json
//	@Param			template	body		healthCheckTemplateRequestBody	true	"health check template object"
//	@Success		200			{object}	standardJsonResponse{data=thunderdome.HealthCheckTemplate}
//	@Failure		400			{object}	standardJsonResponse{}
//	@Failure		500			{object}	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/admin/health-check-templates [post]
func (s *Service) handleAdminHealthCheckTemplateCreate() http.HandlerFunc {
	return s.saveHealthCheckTemplate("", false)
}

// handleAdminHealthCheckTemplateUpdate updates a global default health check template
//
//	@Summary		Update Global Health Check Template
//	@Description	Updates global default health check template
//	@Tags			healthCheck
//	@Produce		json
//	@Param			templateId	path		string							true	"the health check template ID"
//	@Param			template	body		healthCheckTemplateRequestBody	true	"health check template object"
//	@Success		200			{object}	standardJsonResponse{data=thunderdome.HealthCheckTemplate}
//	@Failure		400			{object}	standardJsonResponse{}
//	@Failure		404			{object}	standardJsonResponse{}
//	@Failure		500			{object}	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/admin/health-check-templates/{templateId} [put]
func (s *Service) handleAdminHealthCheckTemplateUpdate() http.HandlerFunc {
	return s.saveHealthCheckTemplate("", true)
}

// handleAdminHealthCheckTemplateDelete deletes a global default health check template
//
//	@Summary		Delete Global Health Check Template
//	@Description	Deletes global default health check template
//	@Tags			healthCheck
//	@Produce		json
//	@Param			templateId	path		string	true	"the health check template ID"
//	@Success		200			{object}	standardJsonResponse{}
//	@Failure		400			{object}	standardJsonResponse{}
//	@Failure		404			{object}	standardJsonResponse{}
//	@Failure		500			{object}	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/admin/health-check-templates/{templateId} [delete]
func (s *Service) handleAdminHealthCheckTemplateDelete() http.HandlerFunc {
	return s.deleteHealthCheckTemplate("")
}

// handleGetOrganizationHealthCheckTemplates gets the organization health check templates
//
//	@Summary		Get Organization Health Check Templates
//	@Description	Gets organization health check templates
//	@Tags			healthCheck
//	@Produce		json
//	@Param			orgId	path		string	true	"Organization ID"
//	@Success		200		{object}	standardJsonResponse{data=[]thunderdome.HealthCheckTemplate}
//	@Failure		400		{object}	standardJsonResponse{}
//	@Failure		500		{object}	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/health-check-templates [get]
func (s *Service) handleGetOrganizationHealthCheckTemplates() http.HandlerFunc {
	return s.listHealthCheckTemplates("orgId")
}

// handleOrganizationHealthCheckTemplateCreate creates a organization health check template
//
//	@Summary		Create Organization Health Check Template
//	@Description	Creates organization health check template
//	@Tags			healthCheck
//	@Produce		json
//	@Param			orgId		path		string							true	"Organization ID"
//	@Param			template	body		healthCheckTemplateRequestBody	true	"health check template object"
//	@Success		200			{object}	standardJsonResponse{data=thunderdome.HealthCheckTemplate}
//	@Failure		400			{object}	standardJsonResponse{}
//	@Failure		500			{object}	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/health-check-templates [post]
func (s *Service) handleOrganizationHealthCheckTemplateCreate() http.HandlerFunc {
	return s.saveHealthCheckTemplate("orgId", false)
}

// handleOrganizationHealthCheckTemplateUpdate updates a organization health check template
//
//	@Summary		Update Organization Health Check Template
//	@Description	Updates organization health check template
//	@Tags			healthCheck
//	@Produce		json
//	@Param			orgId		path		string							true	"Organization ID"
//	@Param			templateId	path		string							true	"the health check template ID"
//	@Param			template	body		healthCheckTemplateRequestBody	true	"health check template object"
//	@Success		200			{object}	standardJsonResponse{data=thunderdome.HealthCheckTemplate}
//	@Failure		400			{object}	standardJsonResponse{}
//	@Failure		404			{object}	standardJsonResponse{}
//	@Failure		500			{object}	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/health-check-templates/{templateId} [put]
func (s *Service) handleOrganizationHealthCheckTemplateUpdate() http.HandlerFunc {
	return s.saveHealthCheckTemplate("orgId", true)
}

// handleOrganizationHealthCheckTemplateDelete deletes a organization health check template
//
//	@Summary		Delete Organization Health Check Template
//	@Description	Deletes organization health check template
//	@Tags			healthCheck
//	@Produce		json
//	@Param			orgId		path		string	true	"Organization ID"
//	@Param			templateId	path		string	true	"the health check template ID"
//	@Success		200			{object}	standardJsonResponse{}
//	@Failure		400			{object}	standardJsonResponse{}
//	@Failure		404			{object}	standardJsonResponse{}
//	@Failure		500			{object}	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/health-check-templates/{templateId} [delete]
func (s *Service) handleOrganizationHealthCheckTemplateDelete() http.HandlerFunc {
	return s.deleteHealthCheckTemplate("orgId")
}

// handleGetTeamHealthCheckTemplates gets the team health check templates
//
//	@Summary		Get Team Health Check Templates
//	@Description	Gets team health check templates
//	@Tags			healthCheck
//	@Produce		json
//	@Param			teamId	path		string	true	"Team ID"
//	@Success		200		{object}	standardJsonResponse{data=[]thunderdome.HealthCheckTemplate}
//	@Failure		400		{object}	standardJsonResponse{}
//	@Failure		500		{object}	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/health-check-templates [get]
func (s *Service) handleGetTeamHealthCheckTemplates() http.HandlerFunc {
	return s.listHealthCheckTemplates("teamId")
}

// handleTeamHealthCheckTemplateCreate creates a team health check template
//
//	@Summary		Create Team Health Check Template
//	@Description	Creates team health check template
//	@Tags			healthCheck
//	@Produce		json
//	@Param			teamId		path		string							true	"Team ID"
//	@Param			template	body		healthCheckTemplateRequestBody	true	"health check template object"
//	@Success		200			{object}	standardJsonResponse{data=thunderdome.HealthCheckTemplate}
//	@Failure		400			{object}	standardJsonResponse{}
//	@Failure		500			{object}	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/health-check-templates [post]
func (s *Service) handleTeamHealthCheckTemplateCreate() http.HandlerFunc {
	return s.saveHealthCheckTemplate("teamId", false)
}

// handleTeamHealthCheckTemplateUpdate updates a team health check template
//
//	@Summary		Update Team Health Check Template
//	@Description	Updates team health check template
//	@Tags			healthCheck
//	@Produce		json
//	@Param			teamId		path		string							true	"Team ID"
//	@Param			templateId	path		string							true	"the health check template ID"
//	@Param			template	body		healthCheckTemplateRequestBody	true	"health check template object"
//	@Success		200			{object}	standardJsonResponse{data=thunderdome.HealthCheckTemplate}
//	@Failure		400			{object}	standardJsonResponse{}
//	@Failure		404			{object}	standardJsonResponse{}
//	@Failure		500			{object}	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/health-check-templates/{templateId} [put]
func (s *Service) handleTeamHealthCheckTemplateUpdate() http.HandlerFunc {
	return s.saveHealthCheckTemplate("teamId", true)
}

// handleTeamHealthCheckTemplateDelete deletes a team health check template
//
//	@Summary		Delete Team Health Check Template
//	@Description	Deletes team health check template
//	@Tags			healthCheck
//	@Produce		json
//	@Param			teamId		path		string	true	"Team ID"
//	@Param			templateId	path		string	true	"the health check template ID"
//	@Success		200			{object}	standardJsonResponse{}
//	@Failure		400			{object}	standardJsonResponse{}
//	@Failure		404			{object}	standardJsonResponse{}
//	@Failure		500			{object}	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/health-check-templates/{templateId} [delete]
func (s *Service) handleTeamHealthCheckTemplateDelete() http.HandlerFunc {
	return s.deleteHealthCheckTemplate("teamId")
}
