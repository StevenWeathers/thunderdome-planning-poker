package thunderdome

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSummarizeHealthCheck(t *testing.T) {
	areas := []*HealthCheckArea{{Name: "Fun"}, {Name: "Speed"}}

	t.Run("totals ratings and trends and sorts comments", func(t *testing.T) {
		results := SummarizeHealthCheck(areas, [][]*HealthCheckAnswer{
			{
				{Area: "Fun", Rating: HealthCheckRatingGreen, Trend: HealthCheckTrendImproving, Comment: "zebra"},
				{Area: "Speed", Rating: HealthCheckRatingRed, Trend: HealthCheckTrendDeclining},
			},
			{
				{Area: "Fun", Rating: HealthCheckRatingYellow, Trend: HealthCheckTrendStable, Comment: "apple"},
				{Area: "Unknown", Rating: HealthCheckRatingRed, Trend: HealthCheckTrendStable},
			},
		})

		require.Len(t, results, 2)
		assert.Equal(t, "Fun", results[0].Name)
		assert.Equal(t, HealthCheckRatingCounts{Green: 1, Yellow: 1}, results[0].Ratings)
		assert.Equal(t, HealthCheckTrendCounts{Improving: 1, Stable: 1}, results[0].Trends)
		assert.Equal(t, []string{"apple", "zebra"}, results[0].Comments)
		assert.Equal(t, HealthCheckRatingCounts{Red: 1}, results[1].Ratings)
		assert.Equal(t, HealthCheckTrendCounts{Declining: 1}, results[1].Trends)
	})

	t.Run("no responses gives empty totals", func(t *testing.T) {
		results := SummarizeHealthCheck(areas, nil)

		require.Len(t, results, 2)
		assert.Equal(t, HealthCheckRatingCounts{}, results[0].Ratings)
		assert.NotNil(t, results[0].Comments)
		assert.Empty(t, results[0].Comments)
	})
}
