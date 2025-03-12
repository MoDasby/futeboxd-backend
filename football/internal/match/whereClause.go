package match

import (
	"fmt"
	"strings"
)

type ListMatchOptions struct {
	Team   int64
	Year   int
	Status string
}

func buildWhereClause(options ListMatchOptions) (string, []any) {
	var whereClause strings.Builder
	var params []any
	counter := 1
	isFirstCondition := true

	appendCondition := func(condition, operator string, param any) {
		if !isFirstCondition {
			whereClause.WriteString(" " + operator + " ")
		}
		whereClause.WriteString(condition)
		params = append(params, param)
		isFirstCondition = false
		counter++
	}

	if options.Team > 0 {
		appendCondition(
			fmt.Sprintf("(home_team_id = $%d OR away_team_id = $%d)", counter, counter),
			"AND",
			options.Team,
		)
	}

	if options.Year > 0 {
		appendCondition(
			fmt.Sprintf("EXTRACT(YEAR FROM m.match_date) = $%d", counter),
			"AND",
			options.Year,
		)
	}

	if options.Status == "completed" {
		appendCondition(
			fmt.Sprintf("completed = $%d", counter),
			"AND",
			"TRUE",
		)
	}

	if options.Status == "live" {
		appendCondition(
			fmt.Sprintf("completed = $%d", counter),
			"AND",
			"FALSE",
		)

		appendCondition(
			fmt.Sprintf("status_name <> $%d", counter),
			"AND",
			"STATUS_SCHEDULED",
		)

		appendCondition(
			fmt.Sprintf("status_name <> $%d", counter),
			"AND",
			"STATUS_POSTPONED",
		)
	}

	if isFirstCondition {
		return "", make([]any, 0)
	}

	return "WHERE " + whereClause.String(), params
}
