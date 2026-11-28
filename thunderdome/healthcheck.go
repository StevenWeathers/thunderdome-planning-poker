package thunderdome

import (
	"errors"
	"slices"
	"time"
)

var (
	ErrHealthCheckNotFound         = errors.New("HEALTH_CHECK_NOT_FOUND")
	ErrHealthCheckClosed           = errors.New("HEALTH_CHECK_CLOSED")
	ErrHealthCheckTemplateNotFound = errors.New("HEALTH_CHECK_TEMPLATE_NOT_FOUND")
)

// HealthCheckMinResponses is the fewest responses for which results are shown, so that answers stay anonymous
const HealthCheckMinResponses = 3

// Health check ratings
const (
	HealthCheckRatingGreen  = "green"
	HealthCheckRatingYellow = "yellow"
	HealthCheckRatingRed    = "red"
)

// Health check trends
const (
	HealthCheckTrendImproving = "improving"
	HealthCheckTrendStable    = "stable"
	HealthCheckTrendDeclining = "declining"
)

// HealthCheckArea is a topic that team members rate, described by what good and bad look like
type HealthCheckArea struct {
	Name      string `json:"name"`
	Healthy   string `json:"healthy"`
	Unhealthy string `json:"unhealthy"`
}

// HealthCheckTemplate is a reusable set of health check areas
type HealthCheckTemplate struct {
	ID             string             `json:"id"`
	Name           string             `json:"name"`
	Description    string             `json:"description"`
	Areas          []*HealthCheckArea `json:"areas"`
	IsPublic       bool               `json:"isPublic"`
	CreatedBy      string             `json:"createdBy"`
	OrganizationID *string            `json:"organizationId"`
	TeamID         *string            `json:"teamId"`
	CreatedAt      time.Time          `json:"createdAt"`
	UpdatedAt      time.Time          `json:"updatedAt"`
}

// HealthCheck is a survey a team fills out before its deadline
type HealthCheck struct {
	ID            string             `json:"id"`
	TeamID        string             `json:"teamId"`
	Name          string             `json:"name"`
	Areas         []*HealthCheckArea `json:"areas"`
	Deadline      time.Time          `json:"deadline"`
	Closed        bool               `json:"closed"`
	CreatedBy     string             `json:"createdBy"`
	ResponseCount int                `json:"responseCount"`
	CreatedAt     time.Time          `json:"createdAt"`
	UpdatedAt     time.Time          `json:"updatedAt"`
}

// HealthCheckAnswer is one member's rating of one area
type HealthCheckAnswer struct {
	Area    string `json:"area"`
	Rating  string `json:"rating"`
	Trend   string `json:"trend"`
	Comment string `json:"comment"`
}

// HealthCheckParticipant identifies a team member who hasn't responded yet
type HealthCheckParticipant struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// HealthCheckRatingCounts totals the ratings of an area
type HealthCheckRatingCounts struct {
	Green  int `json:"green"`
	Yellow int `json:"yellow"`
	Red    int `json:"red"`
}

// HealthCheckTrendCounts totals the trends of an area
type HealthCheckTrendCounts struct {
	Improving int `json:"improving"`
	Stable    int `json:"stable"`
	Declining int `json:"declining"`
}

// HealthCheckAreaResult is the anonymous summary of one area
type HealthCheckAreaResult struct {
	HealthCheckArea
	Ratings  HealthCheckRatingCounts `json:"ratings"`
	Trends   HealthCheckTrendCounts  `json:"trends"`
	Comments []string                `json:"comments"`
}

// SummarizeHealthCheck totals responses per area. Comments are sorted so their order can't be used to tell
// which comments came from the same person.
func SummarizeHealthCheck(areas []*HealthCheckArea, responses [][]*HealthCheckAnswer) []*HealthCheckAreaResult {
	results := make([]*HealthCheckAreaResult, 0, len(areas))
	byName := make(map[string]*HealthCheckAreaResult, len(areas))

	for _, area := range areas {
		result := &HealthCheckAreaResult{HealthCheckArea: *area, Comments: make([]string, 0)}
		results = append(results, result)
		byName[area.Name] = result
	}

	for _, answers := range responses {
		for _, answer := range answers {
			result, ok := byName[answer.Area]
			if !ok {
				continue
			}

			switch answer.Rating {
			case HealthCheckRatingGreen:
				result.Ratings.Green++
			case HealthCheckRatingYellow:
				result.Ratings.Yellow++
			case HealthCheckRatingRed:
				result.Ratings.Red++
			}

			switch answer.Trend {
			case HealthCheckTrendImproving:
				result.Trends.Improving++
			case HealthCheckTrendStable:
				result.Trends.Stable++
			case HealthCheckTrendDeclining:
				result.Trends.Declining++
			}

			if answer.Comment != "" {
				result.Comments = append(result.Comments, answer.Comment)
			}
		}
	}

	for _, result := range results {
		slices.Sort(result.Comments)
	}

	return results
}
