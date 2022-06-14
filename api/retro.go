package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/StevenWeathers/thunderdome-planning-poker/model"
	"github.com/gorilla/mux"
)

// handleRetroCreate handles creating a retro (arena)
func (a *api) handleRetroCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(contextKeyUserID).(string)
		vars := mux.Vars(r)

		body, bodyErr := ioutil.ReadAll(r.Body) // check for errors
		if bodyErr != nil {
			a.logger.Error("error in reading request body: " + bodyErr.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var keyVal struct {
			RetroName string `json:"retroName"`
			Format    string `json:"format"`
			JoinCode  string `json:"joinCode"`
		}
		json.Unmarshal(body, &keyVal) // check for errors

		newRetro, err := a.db.RetroCreate(userID, keyVal.RetroName, keyVal.Format, keyVal.JoinCode)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// if retro created with team association
		TeamID, ok := vars["teamId"]
		if ok {
			OrgRole := r.Context().Value(contextKeyOrgRole)
			DepartmentRole := r.Context().Value(contextKeyDepartmentRole)
			TeamRole := r.Context().Value(contextKeyTeamRole).(string)
			var isAdmin bool
			if DepartmentRole != nil && DepartmentRole.(string) == "ADMIN" {
				isAdmin = true
			}
			if OrgRole != nil && OrgRole.(string) == "ADMIN" {
				isAdmin = true
			}

			if isAdmin == true || TeamRole != "" {
				err := a.db.TeamAddRetro(TeamID, newRetro.Id)

				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}
		}

		a.Success(w, r, http.StatusOK, newRetro, nil)
	}
}

// handleRetroGet looks up retro or returns notfound status
func (a *api) handleRetroGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		RetroID := vars["retroId"]

		retro, err := a.db.RetroGet(RetroID)

		if err != nil {
			http.NotFound(w, r)
			return
		}

		a.Success(w, r, http.StatusOK, retro, nil)
	}
}

// handleRetrosGetByUser looks up retros associated with userID
func (a *api) handleRetrosGetByUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(contextKeyUserID).(string)

		retros, err := a.db.RetroGetByUser(userID)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		a.Success(w, r, http.StatusOK, retros, nil)
	}
}

// handleGetRetros gets a list of retros
// @Summary Get Retros
// @Description get list of retros
// @Tags retro
// @Produce  json
// @Param limit query int false "Max number of results to return"
// @Param offset query int false "Starting point to return rows from, should be multiplied by limit or 0"
// @Param active query boolean false "Only active retros"
// @Success 200 object standardJsonResponse{data=[]model.Retro}
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /retros [get]
func (a *api) handleGetRetros() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Limit, Offset := getLimitOffsetFromRequest(r)
		query := r.URL.Query()
		var err error
		var Count int
		var Retros []*model.Retro
		Active, _ := strconv.ParseBool(query.Get("active"))

		if Active {
			Retros, Count, err = a.db.GetActiveRetros(Limit, Offset)
		} else {
			Retros, Count, err = a.db.GetRetros(Limit, Offset)
		}

		if err != nil {
			a.Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		Meta := &pagination{
			Count:  Count,
			Offset: Offset,
			Limit:  Limit,
		}

		a.Success(w, r, http.StatusOK, Retros, Meta)
	}
}
