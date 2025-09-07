package http

import "net/http"

// handleGetProjectPokerGames retrieves poker games for a specific project
//
//	@Summary		Get Project Pokers
//	@Description	Retrieve poker games for a specific project
//	@Tags			poker
//	@Produce		json
//	@Param			projectId	path	string	true	"the project ID to get pokergames for"
//	@Param			limit		query	int		false	"Max number of results to return"
//	@Param			offset		query	int		false	"Starting point to return rows from, should be multiplied by limit or 0"
//	@Success		200			object	standardJsonResponse{data=[]thunderdome.Poker}
//	@Failure		403			object	standardJsonResponse{}
//	@Failure		404			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/projects/{projectId}/poker [get]
func (s *Service) handleGetProjectPokerGames() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		projectID := r.PathValue("projectId")
		idErr := validate.Var(projectID, "required,uuid")
		if idErr != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, idErr.Error()))
			return
		}
		limit, offset := getLimitOffsetFromRequest(r)

		games, err := s.ProjectDataSvc.ListPokerGames(ctx, projectID, limit, offset)
		if err != nil {
			s.Failure(w, r, http.StatusInternalServerError, Errorf(EINTERNAL, err.Error()))
			return
		}

		s.Success(w, r, http.StatusOK, games, nil)
	}
}
