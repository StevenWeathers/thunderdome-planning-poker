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

type userTitleRequestBody struct {
	Name string `json:"name" validate:"required,max=128"`
}

func (s *Service) userTitleFailure(w http.ResponseWriter, r *http.Request, action string, err error) {
	switch {
	case errors.Is(err, thunderdome.ErrDuplicateUserTitle):
		s.Failure(w, r, http.StatusConflict, Errorf(ECONFLICT, err.Error()))
	case errors.Is(err, thunderdome.ErrUserTitleNotFound):
		s.Failure(w, r, http.StatusNotFound, Errorf(ENOTFOUND, err.Error()))
	default:
		s.Logger.Ctx(r.Context()).Error(action+" error", zap.Error(err),
			zap.String("organization_id", r.PathValue("orgId")),
			zap.String("department_id", r.PathValue("departmentId")),
			zap.String("team_id", r.PathValue("teamId")),
			zap.String("title_id", r.PathValue("titleId")))
		s.Failure(w, r, http.StatusInternalServerError, err)
	}
}

// handleGetUserTitles gets the titles offered when editing members of an organization, department or team
//
//	@Summary		Get User Titles
//	@Description	get the organization's list of member titles, empty for teams not in an organization
//	@Tags			organization
//	@Produce		json
//	@Param			orgId	path	string	true	"the organization ID"
//	@Success		200		object	standardJsonResponse{data=[]thunderdome.UserTitle}
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/user-titles [get]
func (s *Service) handleGetUserTitles() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ids := map[string]string{}
		for _, param := range []string{"orgId", "departmentId", "teamId"} {
			id := r.PathValue(param)
			if id != "" {
				if err := validate.Var(id, "uuid"); err != nil {
					s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
					return
				}
			}
			ids[param] = id
		}

		titles, err := s.UserTitleDataSvc.UserTitleList(r.Context(), ids["orgId"], ids["departmentId"], ids["teamId"])
		if err != nil {
			s.userTitleFailure(w, r, "handleGetUserTitles", err)
			return
		}

		s.Success(w, r, http.StatusOK, titles, nil)
	}
}

// saveUserTitle creates an organization title, or renames it when the route includes a title ID
func (s *Service) saveUserTitle(update bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := &thunderdome.UserTitle{OrganizationID: r.PathValue("orgId")}
		if err := validate.Var(title.OrganizationID, "required,uuid"); err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}
		if update {
			title.ID = r.PathValue("titleId")
			if err := validate.Var(title.ID, "required,uuid"); err != nil {
				s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
				return
			}
		}

		request := userTitleRequestBody{}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}
		if err := json.Unmarshal(body, &request); err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}
		request.Name = strings.TrimSpace(request.Name)
		if err := validate.Struct(request); err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}
		title.Name = request.Name

		if update {
			err = s.UserTitleDataSvc.UserTitleUpdate(r.Context(), title)
		} else {
			err = s.UserTitleDataSvc.UserTitleCreate(r.Context(), title)
		}
		if err != nil {
			s.userTitleFailure(w, r, "saveUserTitle", err)
			return
		}

		s.Success(w, r, http.StatusOK, title, nil)
	}
}

// handleUserTitleCreate adds a title to an organization
//
//	@Summary		Create User Title
//	@Description	Adds a member title to the organization
//	@Tags			organization
//	@Produce		json
//	@Param			orgId	path	string					true	"the organization ID"
//	@Param			title	body	userTitleRequestBody	true	"new title object"
//	@Success		200		object	standardJsonResponse{data=thunderdome.UserTitle}
//	@Failure		400		object	standardJsonResponse{}
//	@Failure		409		object	standardJsonResponse{}
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/user-titles [post]
func (s *Service) handleUserTitleCreate() http.HandlerFunc {
	return s.saveUserTitle(false)
}

// handleUserTitleUpdate renames an organization title
//
//	@Summary		Update User Title
//	@Description	Renames a member title of the organization
//	@Tags			organization
//	@Produce		json
//	@Param			orgId	path	string					true	"the organization ID"
//	@Param			titleId	path	string					true	"the title ID"
//	@Param			title	body	userTitleRequestBody	true	"updated title object"
//	@Success		200		object	standardJsonResponse{data=thunderdome.UserTitle}
//	@Failure		400		object	standardJsonResponse{}
//	@Failure		404		object	standardJsonResponse{}
//	@Failure		409		object	standardJsonResponse{}
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/user-titles/{titleId} [put]
func (s *Service) handleUserTitleUpdate() http.HandlerFunc {
	return s.saveUserTitle(true)
}

// handleUserTitleDelete removes an organization title
//
//	@Summary		Delete User Title
//	@Description	Removes a member title from the organization
//	@Tags			organization
//	@Produce		json
//	@Param			orgId	path	string	true	"the organization ID"
//	@Param			titleId	path	string	true	"the title ID"
//	@Success		200		object	standardJsonResponse{}
//	@Failure		404		object	standardJsonResponse{}
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/organizations/{orgId}/user-titles/{titleId} [delete]
func (s *Service) handleUserTitleDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID := r.PathValue("orgId")
		titleID := r.PathValue("titleId")
		for _, id := range []string{orgID, titleID} {
			if err := validate.Var(id, "required,uuid"); err != nil {
				s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
				return
			}
		}

		if err := s.UserTitleDataSvc.UserTitleDelete(r.Context(), orgID, titleID); err != nil {
			s.userTitleFailure(w, r, "handleUserTitleDelete", err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}
