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
			fmt.Sprintf("status_name = $%d", counter),
			"AND",
			"STATUS_FIRST_HALF",
		)

		appendCondition(
			fmt.Sprintf("status_name = $%d", counter),
			"AND",
			"STATUS_SCHEDULED",
		)

		appendCondition(
			fmt.Sprintf("status_name = $%d", counter),
			"AND",
			"STATUS_SECOND_HALF",
		)

		appendCondition(
			fmt.Sprintf("status_name = $%d", counter),
			"AND",
			"STATUS_HALFTIME",
		)

		appendCondition(
			fmt.Sprintf("status_name = $%d", counter),
			"AND",
			"STATUS_HALFTIME",
		)

		appendCondition(
			fmt.Sprintf("status_name = $%d", counter),
			"AND",
			"STATUS_OVERTIME",
		)

		appendCondition(
			fmt.Sprintf("status_name = $%d", counter),
			"AND",
			"STATUS_SHOOTOUT",
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

		appendCondition(
			fmt.Sprintf("status_name <> $%d", counter),
			"AND",
			"STATUS_ABANDONED",
		)
	}

	if isFirstCondition {
		return "", make([]any, 0)
	}

	return "WHERE " + whereClause.String(), params
}
