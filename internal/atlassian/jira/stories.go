package jira

import (
	"context"
)

// StoriesJQLSearch searches for stories in Jira using JQL
func (c *Client) StoriesJQLSearch(ctx context.Context, jql string, fields []string, startAt int, maxResults int) (*IssuesSearchResult, error) {
	iss := IssuesSearchResult{}

	issues, _, err := c.instance.Issue.Search.SearchJQL(ctx, jql, fields, nil, maxResults, "")
	if err != nil {
		return nil, err
	}

	iss.Total = issues.Total
	iss.Issues = issues.Issues

	return &iss, err
}
