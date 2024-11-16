package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildWhereClause(t *testing.T) {
	uc := &ListReviewsUseCase{}

	tests := []struct {
		name           string
		options        ListReviewsOptions
		expectedQuery  string
		expectedParams []any
	}{
		{
			name: "Only Match",
			options: ListReviewsOptions{
				Match: "123",
			},
			expectedQuery:  "WHERE match_id = $1",
			expectedParams: []any{"123"},
		},
		{
			name: "Match and Team",
			options: ListReviewsOptions{
				Match: "123",
				Team:  "Team A",
			},
			expectedQuery:  "WHERE match_id = $1 AND (home_team_id = $2 OR away_team_id = $2)",
			expectedParams: []any{"123", "Team A"},
		},
		{
			name: "Match, Team, and Username",
			options: ListReviewsOptions{
				Match:    "123",
				Team:     "Team A",
				Username: "john_doe",
			},
			expectedQuery:  "WHERE match_id = $1 AND (home_team_id = $2 OR away_team_id = $2) AND username = $3",
			expectedParams: []any{"123", "Team A", "john_doe"},
		},
		{
			name: "Match, Team, Username, and SearchTerm",
			options: ListReviewsOptions{
				Match:      "123",
				Team:       "Team A",
				Username:   "john_doe",
				SearchTerm: "good match",
			},
			expectedQuery:  "WHERE match_id = $1 AND (home_team_id = $2 OR away_team_id = $2) AND username = $3 AND r.search_vector @@ websearch_to_tsquery('portuguese', $4)",
			expectedParams: []any{"123", "Team A", "john_doe", "good match"},
		},
		{
			name:           "No Filters",
			options:        ListReviewsOptions{},
			expectedQuery:  "",
			expectedParams: []any{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, params := uc.buildWhereClause(tt.options)

			assert.Equal(t, tt.expectedQuery, query)
			assert.ElementsMatch(t, tt.expectedParams, params)
		})
	}
}
