package jira

import (
	"context"
)

func (c *Client) StoriesJQLSearch(ctx context.Context, jql string, fields []string, startAt int, maxResults int) (*IssuesSearchResult, error) {
	iss := IssuesSearchResult{}
	issues, _, err := c.instance.Issue.Search.Post(ctx, jql, fields, nil, startAt, maxResults, "")
	if err != nil {
		return nil, err
	}

	iss.Total = issues.Total
	iss.Issues = issues.Issues

	return &iss, err
}
