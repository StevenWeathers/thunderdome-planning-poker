package jiradatacenter

import (
	"context"

	jira_data_center "github.com/andygrunwald/go-jira/v2/onpremise"
)

// StoriesJQLSearch searches for stories in Jira using JQL
func (c *Client) StoriesJQLSearch(ctx context.Context, jql string, fields []string, startAt int, maxResults int) (*IssuesSearchResult, error) {
	iss := IssuesSearchResult{}
	opt := &jira_data_center.SearchOptions{
		MaxResults: maxResults, // Max results can go up to 1000
		StartAt:    startAt,
	}
	issues, respo, err := c.instance.Issue.Search(ctx, jql, opt)

	if err != nil {
		return nil, err
	}

	iss.Total = respo.Total
	iss.Issues = issues

	return &iss, err
}
