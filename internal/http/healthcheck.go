package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/StevenWeathers/thunderdome-planning-poker/thunderdome"
	"go.uber.org/zap"
)

// healthCheckInviteLimit is the most team members that are emailed about a health check
const healthCheckInviteLimit = 1000

type healthCheckCreateRequestBody struct {
	Name       string    `json:"name" validate:"required,max=255"`
	TemplateID string    `json:"templateId" validate:"required,uuid"`
	Deadline   time.Time `json:"deadline" validate:"required,gt"`
	SendEmail  bool      `json:"sendEmail"`
}

type healthCheckUpdateRequestBody struct {
	Name     string    `json:"name" validate:"required,max=255"`
	Deadline time.Time `json:"deadline" validate:"required,gt"`
}

type healthCheckAnswerRequestBody struct {
	Area    string `json:"area" validate:"required"`
	Rating  string `json:"rating" validate:"required,oneof=green yellow red"`
	Trend   string `json:"trend" validate:"required,oneof=improving stable declining"`
	Comment string `json:"comment" validate:"max=1000"`
}

type healthCheckResponseRequestBody struct {
	Answers []healthCheckAnswerRequestBody `json:"answers" validate:"required,min=1,unique=Area,dive"`
}

type healthCheckDetail struct {
	HealthCheck  *thunderdome.HealthCheck         `json:"healthCheck"`
	CanRespond   bool                             `json:"canRespond"`
	CanManage    bool                             `json:"canManage"`
	MinResponses int                              `json:"minResponses"`
	MyAnswers    []*thunderdome.HealthCheckAnswer `json:"myAnswers"`
	// Pending lists the members yet to respond, only given to those who manage the health check
	Pending []*thunderdome.HealthCheckParticipant `json:"pending"`
	// Results is null until the health check is closed with enough responses
	Results []*thunderdome.HealthCheckAreaResult `json:"results"`
}

// isTeamParticipant reports whether the user belongs to the team itself,
// organization and department admins who don't can only view results
func isTeamParticipant(roles *thunderdome.UserTeamRoleInfo) bool {
	return roles.TeamRole != nil
}

// canManageHealthCheck reports whether the user may change a health check: its creator or a team admin
// who is a team member, or an application admin
func canManageHealthCheck(userID string, userType string, roles *thunderdome.UserTeamRoleInfo, hc *thunderdome.HealthCheck) bool {
	if userType == thunderdome.AdminUserType {
		return true
	}
	if !isTeamParticipant(roles) {
		return false
	}

	isAdmin := func(role *string) bool { return role != nil && *role == thunderdome.AdminUserType }

	return hc.CreatedBy == userID || isAdmin(roles.TeamRole) || isAdmin(roles.DepartmentRole) || isAdmin(roles.OrganizationRole)
}

func (s *Service) healthCheckFailure(w http.ResponseWriter, r *http.Request, action string, err error) {
	switch {
	case errors.Is(err, thunderdome.ErrHealthCheckNotFound), errors.Is(err, thunderdome.ErrHealthCheckTemplateNotFound):
		s.Failure(w, r, http.StatusNotFound, Errorf(ENOTFOUND, err.Error()))
	case errors.Is(err, thunderdome.ErrHealthCheckClosed):
		s.Failure(w, r, http.StatusConflict, Errorf(ECONFLICT, err.Error()))
	default:
		s.Logger.Ctx(r.Context()).Error(action+" error", zap.Error(err),
			zap.String("organization_id", r.PathValue("orgId")),
			zap.String("team_id", r.PathValue("teamId")),
			zap.String("health_check_id", r.PathValue("healthCheckId")),
			zap.String("template_id", r.PathValue("templateId")))
		s.Failure(w, r, http.StatusInternalServerError, err)
	}
}

// decodeHealthCheckBody reads a JSON request body, writing the failure response when it can't
func (s *Service) decodeHealthCheckBody(w http.ResponseWriter, r *http.Request, v any) bool {
	body, err := io.ReadAll(r.Body)
	if err == nil {
		err = json.Unmarshal(body, v)
	}
	if err != nil {
		s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
		return false
	}

	return true
}

// requestedHealthCheck gets the {healthCheckId} health check of the {teamId} team, writing the failure response when it can't
func (s *Service) requestedHealthCheck(w http.ResponseWriter, r *http.Request) (*thunderdome.HealthCheck, bool) {
	teamID := r.PathValue("teamId")
	healthCheckID := r.PathValue("healthCheckId")
	for _, id := range []string{teamID, healthCheckID} {
		if err := validate.Var(id, "required,uuid"); err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return nil, false
		}
	}

	hc, err := s.HealthCheckDataSvc.HealthCheckGet(r.Context(), teamID, healthCheckID)
	if err != nil {
		s.healthCheckFailure(w, r, "requestedHealthCheck", err)
		return nil, false
	}

	return hc, true
}

// managedHealthCheck is requestedHealthCheck restricted to users who can manage it
func (s *Service) managedHealthCheck(w http.ResponseWriter, r *http.Request) (*thunderdome.HealthCheck, bool) {
	hc, ok := s.requestedHealthCheck(w, r)
	if !ok {
		return nil, false
	}

	ctx := r.Context()
	userID := ctx.Value(contextKeyUserID).(string)
	userType := ctx.Value(contextKeyUserType).(string)
	roles := ctx.Value(contextKeyUserTeamRoles).(*thunderdome.UserTeamRoleInfo)

	if !canManageHealthCheck(userID, userType, roles, hc) {
		s.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "REQUIRES_HEALTH_CHECK_MANAGER"))
		return nil, false
	}

	return hc, true
}

// handleGetTeamHealthChecks gets a list of team health checks
//
//	@Summary		Get Team Health Checks
//	@Description	get a team's health checks, newest first
//	@Tags			healthCheck
//	@Produce		json
//	@Param			teamId	path	string	true	"the team ID"
//	@Param			limit	query	int		false	"Max number of results to return"
//	@Param			offset	query	int		false	"Starting point to return rows from"
//	@Success		200		object	standardJsonResponse{data=[]thunderdome.HealthCheck}
//	@Failure		400		object	standardJsonResponse{}
//	@Failure		500		object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/health-checks [get]
func (s *Service) handleGetTeamHealthChecks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		teamID := r.PathValue("teamId")
		if err := validate.Var(teamID, "required,uuid"); err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}
		limit, offset := getLimitOffsetFromRequest(r)

		checks, count, err := s.HealthCheckDataSvc.HealthCheckList(r.Context(), teamID, limit, offset)
		if err != nil {
			s.healthCheckFailure(w, r, "handleGetTeamHealthChecks", err)
			return
		}

		s.Success(w, r, http.StatusOK, checks, &pagination{Count: count, Limit: limit, Offset: offset})
	}
}

// handleGetTeamHealthCheck gets a health check, with results once it's closed
//
//	@Summary		Get Team Health Check
//	@Description	get a health check, the requester's answers and, once closed with enough responses, the anonymous results
//	@Tags			healthCheck
//	@Produce		json
//	@Param			teamId			path	string	true	"the team ID"
//	@Param			healthCheckId	path	string	true	"the health check ID"
//	@Success		200				object	standardJsonResponse{data=healthCheckDetail}
//	@Failure		400				object	standardJsonResponse{}
//	@Failure		404				object	standardJsonResponse{}
//	@Failure		500				object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/health-checks/{healthCheckId} [get]
func (s *Service) handleGetTeamHealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hc, ok := s.requestedHealthCheck(w, r)
		if !ok {
			return
		}

		ctx := r.Context()
		userID := ctx.Value(contextKeyUserID).(string)
		userType := ctx.Value(contextKeyUserType).(string)
		roles := ctx.Value(contextKeyUserTeamRoles).(*thunderdome.UserTeamRoleInfo)

		detail := &healthCheckDetail{
			HealthCheck:  hc,
			CanRespond:   isTeamParticipant(roles),
			CanManage:    canManageHealthCheck(userID, userType, roles, hc),
			MinResponses: thunderdome.HealthCheckMinResponses,
		}

		var err error
		detail.MyAnswers, err = s.HealthCheckDataSvc.HealthCheckUserResponse(ctx, hc.ID, userID)
		if err != nil {
			s.healthCheckFailure(w, r, "handleGetTeamHealthCheck", err)
			return
		}

		if detail.CanManage {
			detail.Pending, err = s.HealthCheckDataSvc.HealthCheckPendingUsers(ctx, hc.TeamID, hc.ID)
			if err != nil {
				s.healthCheckFailure(w, r, "handleGetTeamHealthCheck", err)
				return
			}
		}

		// nothing is summarized while open, so that nobody can watch totals change as members respond
		if hc.Closed && hc.ResponseCount >= thunderdome.HealthCheckMinResponses {
			responses, err := s.HealthCheckDataSvc.HealthCheckResponses(ctx, hc.ID)
			if err != nil {
				s.healthCheckFailure(w, r, "handleGetTeamHealthCheck", err)
				return
			}
			detail.Results = thunderdome.SummarizeHealthCheck(hc.Areas, responses)
		}

		s.Success(w, r, http.StatusOK, detail, nil)
	}
}

// handleTeamHealthCheckCreate creates a team health check
//
//	@Summary		Create Team Health Check
//	@Description	Creates a health check from a template, optionally emailing the team a link to it
//	@Tags			healthCheck
//	@Produce		json
//	@Param			teamId		path	string							true	"the team ID"
//	@Param			healthCheck	body	healthCheckCreateRequestBody	true	"new health check object"
//	@Success		200			object	standardJsonResponse{data=thunderdome.HealthCheck}
//	@Failure		400			object	standardJsonResponse{}
//	@Failure		403			object	standardJsonResponse{}
//	@Failure		404			object	standardJsonResponse{}
//	@Failure		500			object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/health-checks [post]
func (s *Service) handleTeamHealthCheckCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID := ctx.Value(contextKeyUserID).(string)
		userType := ctx.Value(contextKeyUserType).(string)
		roles := ctx.Value(contextKeyUserTeamRoles).(*thunderdome.UserTeamRoleInfo)

		teamID := r.PathValue("teamId")
		if err := validate.Var(teamID, "required,uuid"); err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}
		if userType != thunderdome.AdminUserType && !isTeamParticipant(roles) {
			s.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "REQUIRES_TEAM_USER"))
			return
		}

		request := healthCheckCreateRequestBody{}
		if !s.decodeHealthCheckBody(w, r, &request) {
			return
		}
		request.Name = strings.TrimSpace(request.Name)
		if err := validate.Struct(request); err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}

		id, err := s.HealthCheckDataSvc.HealthCheckCreate(ctx, teamID, request.TemplateID, request.Name, request.Deadline, userID)
		if err != nil {
			s.healthCheckFailure(w, r, "handleTeamHealthCheckCreate", err)
			return
		}

		hc, err := s.HealthCheckDataSvc.HealthCheckGet(ctx, teamID, id)
		if err != nil {
			s.healthCheckFailure(w, r, "handleTeamHealthCheckCreate", err)
			return
		}

		if request.SendEmail {
			s.sendHealthCheckInvites(r, hc)
		}

		s.Success(w, r, http.StatusOK, hc, nil)
	}
}

// sendHealthCheckInvites emails the team a link to the health check. The health check already exists,
// so failures are only logged. Mail is sent in the background to not hold up the request.
func (s *Service) sendHealthCheckInvites(r *http.Request, hc *thunderdome.HealthCheck) {
	ctx := r.Context()

	team, err := s.TeamDataSvc.TeamGetByID(ctx, hc.TeamID)
	if err != nil {
		s.Logger.Ctx(ctx).Error("sendHealthCheckInvites team error", zap.Error(err), zap.String("team_id", hc.TeamID))
		return
	}
	users, _, err := s.TeamDataSvc.TeamUserList(ctx, hc.TeamID, healthCheckInviteLimit, 0)
	if err != nil {
		s.Logger.Ctx(ctx).Error("sendHealthCheckInvites users error", zap.Error(err), zap.String("team_id", hc.TeamID))
		return
	}

	go func() {
		for _, user := range users {
			if user.Email == "" {
				continue
			}
			// the email service logs its own failures
			_ = s.Email.SendHealthCheckInvite(user.Name, user.Email, team.Name, hc.Name, hc.Deadline, hc.TeamID, hc.ID)
		}
	}()
}

// handleTeamHealthCheckUpdate renames an open health check or changes its deadline
//
//	@Summary		Update Team Health Check
//	@Description	Updates the name and deadline of an open health check
//	@Tags			healthCheck
//	@Produce		json
//	@Param			teamId			path	string							true	"the team ID"
//	@Param			healthCheckId	path	string							true	"the health check ID"
//	@Param			healthCheck		body	healthCheckUpdateRequestBody	true	"updated health check object"
//	@Success		200				object	standardJsonResponse{}
//	@Failure		400				object	standardJsonResponse{}
//	@Failure		403				object	standardJsonResponse{}
//	@Failure		404				object	standardJsonResponse{}
//	@Failure		409				object	standardJsonResponse{}
//	@Failure		500				object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/health-checks/{healthCheckId} [put]
func (s *Service) handleTeamHealthCheckUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hc, ok := s.managedHealthCheck(w, r)
		if !ok {
			return
		}

		request := healthCheckUpdateRequestBody{}
		if !s.decodeHealthCheckBody(w, r, &request) {
			return
		}
		request.Name = strings.TrimSpace(request.Name)
		if err := validate.Struct(request); err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}

		if err := s.HealthCheckDataSvc.HealthCheckUpdate(r.Context(), hc.TeamID, hc.ID, request.Name, request.Deadline); err != nil {
			s.healthCheckFailure(w, r, "handleTeamHealthCheckUpdate", err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleTeamHealthCheckClose ends a health check now
//
//	@Summary		Close Team Health Check
//	@Description	Ends an open health check now, it can't be reopened
//	@Tags			healthCheck
//	@Produce		json
//	@Param			teamId			path	string	true	"the team ID"
//	@Param			healthCheckId	path	string	true	"the health check ID"
//	@Success		200				object	standardJsonResponse{}
//	@Failure		403				object	standardJsonResponse{}
//	@Failure		404				object	standardJsonResponse{}
//	@Failure		409				object	standardJsonResponse{}
//	@Failure		500				object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/health-checks/{healthCheckId}/close [post]
func (s *Service) handleTeamHealthCheckClose() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hc, ok := s.managedHealthCheck(w, r)
		if !ok {
			return
		}

		if err := s.HealthCheckDataSvc.HealthCheckClose(r.Context(), hc.TeamID, hc.ID); err != nil {
			s.healthCheckFailure(w, r, "handleTeamHealthCheckClose", err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleTeamHealthCheckDelete deletes a health check
//
//	@Summary		Delete Team Health Check
//	@Description	Deletes a health check and its responses
//	@Tags			healthCheck
//	@Produce		json
//	@Param			teamId			path	string	true	"the team ID"
//	@Param			healthCheckId	path	string	true	"the health check ID"
//	@Success		200				object	standardJsonResponse{}
//	@Failure		403				object	standardJsonResponse{}
//	@Failure		404				object	standardJsonResponse{}
//	@Failure		500				object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/health-checks/{healthCheckId} [delete]
func (s *Service) handleTeamHealthCheckDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hc, ok := s.managedHealthCheck(w, r)
		if !ok {
			return
		}

		if err := s.HealthCheckDataSvc.HealthCheckDelete(r.Context(), hc.TeamID, hc.ID); err != nil {
			s.healthCheckFailure(w, r, "handleTeamHealthCheckDelete", err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleTeamHealthCheckRespond saves the requester's answers
//
//	@Summary		Respond to Team Health Check
//	@Description	Saves or replaces the requesting team member's anonymous answers to an open health check
//	@Tags			healthCheck
//	@Produce		json
//	@Param			teamId			path	string							true	"the team ID"
//	@Param			healthCheckId	path	string							true	"the health check ID"
//	@Param			response		body	healthCheckResponseRequestBody	true	"an answer for every area"
//	@Success		200				object	standardJsonResponse{}
//	@Failure		400				object	standardJsonResponse{}
//	@Failure		403				object	standardJsonResponse{}
//	@Failure		404				object	standardJsonResponse{}
//	@Failure		409				object	standardJsonResponse{}
//	@Failure		500				object	standardJsonResponse{}
//	@Security		ApiKeyAuth
//	@Router			/teams/{teamId}/health-checks/{healthCheckId}/response [put]
func (s *Service) handleTeamHealthCheckRespond() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID := ctx.Value(contextKeyUserID).(string)
		roles := ctx.Value(contextKeyUserTeamRoles).(*thunderdome.UserTeamRoleInfo)

		if !isTeamParticipant(roles) {
			s.Failure(w, r, http.StatusForbidden, Errorf(EUNAUTHORIZED, "REQUIRES_TEAM_USER"))
			return
		}

		hc, ok := s.requestedHealthCheck(w, r)
		if !ok {
			return
		}

		request := healthCheckResponseRequestBody{}
		if !s.decodeHealthCheckBody(w, r, &request) {
			return
		}
		for i := range request.Answers {
			request.Answers[i].Comment = strings.TrimSpace(request.Answers[i].Comment)
		}
		if err := validate.Struct(request); err != nil {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, err.Error()))
			return
		}

		// area names are unique, so matching counts and known names mean every area is answered
		areas := make(map[string]bool, len(hc.Areas))
		for _, area := range hc.Areas {
			areas[area.Name] = true
		}
		answers := make([]*thunderdome.HealthCheckAnswer, 0, len(request.Answers))
		for _, a := range request.Answers {
			if !areas[a.Area] {
				s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "HEALTH_CHECK_UNKNOWN_AREA"))
				return
			}
			answers = append(answers, &thunderdome.HealthCheckAnswer{
				Area: a.Area, Rating: a.Rating, Trend: a.Trend, Comment: a.Comment,
			})
		}
		if len(answers) != len(hc.Areas) {
			s.Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "HEALTH_CHECK_INCOMPLETE"))
			return
		}

		if err := s.HealthCheckDataSvc.HealthCheckResponseSubmit(ctx, hc.TeamID, hc.ID, userID, answers); err != nil {
			s.healthCheckFailure(w, r, "handleTeamHealthCheckRespond", err)
			return
		}

		s.Success(w, r, http.StatusOK, nil, nil)
	}
}
