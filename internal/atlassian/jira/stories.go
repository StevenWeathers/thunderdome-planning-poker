package jira

import (
	"context"

	"go.uber.org/zap"
)

func (c *Client) StoriesJQLSearch(ctx context.Context, jql string, fields []string, startAt int, maxResults int) (interface{}, error) {
	logger := c.logger.Ctx(ctx)
	issues, _, err := c.instance.Issue.Search.Post(ctx, jql, fields, nil, startAt, maxResults, "")
	if err != nil {
		logger.Error(
			"Error searching for Jira stories by JQL", zap.Error(err), zap.String("jql", jql),
			zap.Any("fields", fields), zap.Int("startAt", startAt), zap.Int("maxResults", maxResults),
		)
	}

	return issues, err
}
