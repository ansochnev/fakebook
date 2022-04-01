package backend

import "strings"

func formatSQLQuery(query string) string {
	return strings.ReplaceAll(query, "\t", " ")
}
