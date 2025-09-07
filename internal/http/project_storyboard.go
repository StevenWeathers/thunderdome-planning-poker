package http

import "net/http"

// handleGetProjectStoryboards retrieves storyboards for a specific project
//
//	@Summary		Get Project Storyboards
//	@Description	Retrieve storyboards for a specific project
//	@Tags			storyboards
//	@Produce		json
//	@Param			projectId	path	string	true	"the project ID to get storyboards for"
//	@Param			limit		query	int		false	"Max number of results to return"
//	@Param			offset		query	int		false	"Starting point to return rows from, should be multiplied by limit or 0"
//	@Success		200			object	standardJsonResponse{data=[]thunderdome.Storyboard}
//	@Failure		403			object	standardJsonResponse{}
//	@Failure		404			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/projects/{projectId}/storyboards [get]
func (s *Service) handleGetProjectStoryboards() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		projectID := r.PathValue("projectId")
		idErr := validate.Var(projectID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		limit, offset := getLimitOffsetFromRequest(r)

		storyboards, err := s.ProjectDataSvc.ListStoryboards(ctx, projectID, limit, offset)
		if err != nil {
			s.Failure(w, r, http.StatusInternalServerError, Errorf(EINTERNAL, err.Error()))
			return
		}

		s.Success(w, r, http.StatusOK, storyboards, nil)
	}
}
