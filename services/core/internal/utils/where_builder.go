package utils

import (
	"fmt"
	"strings"
)

type filter struct {
	ftype string // AND ou OR se for uma unica consulta não há necessidade de usar
	name  string
	value string
}

type groupFilter struct {
	filters []filter
}

type WhereBuilder struct {
	filters []filter
	groups  []groupFilter
}

func NewWhereBuilder() *WhereBuilder {
	return &WhereBuilder{filters: make([]filter, 0), groups: make([]groupFilter, 0)}
}

func (wb *WhereBuilder) AddFilter(filterName string, value string, filterType string) {
	wb.filters = append(wb.filters, filter{ftype: filterType, value: value, name: filterName})
}

func (wb *WhereBuilder) NewFilter(filterName string, value string, filterType string) filter {
	return filter{ftype: filterType, value: value, name: filterName}
}

func (wb *WhereBuilder) AddGroupFilter(filterType string, filters ...filter) {
	wb.groups = append(wb.groups, groupFilter{filters: filters})
}

func (wb *WhereBuilder) filterValues(filters []filter) []filter {
	if len(filters) == 0 {
		return filters
	}

	filteredValues := make([]filter, 0)

	for _, filter := range filters {
		if filter.value != "" {
			filteredValues = append(filteredValues, filter)
		}
	}

	return filteredValues
}

func (wb *WhereBuilder) Build() (string, []any) {
	var builder strings.Builder
	params := make([]any, 0)
	counter := 1

	filteredValues := wb.filterValues(wb.filters)

	filteredGroups := make([][]filter, 0)
	for _, group := range wb.groups {
		filteredGroup := wb.filterValues(group.filters)
		if len(filteredGroup) > 0 {
			filteredGroups = append(filteredGroups, filteredGroup)
		}
	}

	for i, filter := range filteredValues {
		if builder.Len() == 0 {
			builder.WriteString(fmt.Sprintf("WHERE %s = $%d ", filter.name, counter))
		} else {
			builder.WriteString(fmt.Sprintf("%s = $%d ", filter.name, counter))
		}

		if i+1 < len(filteredValues) {
			builder.WriteString(fmt.Sprintf(" %s ", filter.ftype))
		}

		counter++
		params = append(params, filter.value)
	}

	for _, groupFilters := range filteredGroups {
		if builder.Len() == 0 {
			builder.WriteString("WHERE (")
		} else {
			builder.WriteString(" AND (")
		}

		for i, filter := range groupFilters {
			if i > 0 {
				builder.WriteString(fmt.Sprintf(" %s ", filter.ftype))
			}
			builder.WriteString(fmt.Sprintf("%s = $%d", filter.name, counter))
			counter++
			params = append(params, filter.value)
		}

		builder.WriteString(")")
	}

	return builder.String(), params
}
