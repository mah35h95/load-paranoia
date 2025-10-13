package model

import "cloud.google.com/go/bigquery"

const (
	// QueryMaxLength - The maximum standard SQL query length is 1024.00K characters,
	// including comments and white space characters.
	QueryMaxLength int = 1000000
	// MaxSubQueries - The maximum number of subqueries that can be combined in a single query
	// is 300. This is a limitation of BigQuery on query complexity and can be adjusted as needed.
	MaxSubQueries int = 300
)

type IntervalRowCountResult struct {
	EffectedRowCount bigquery.NullInt64
	FromTimestamp    bigquery.NullInt64
	ToTimestamp      bigquery.NullInt64
}

type TableDetails struct {
	TableID string
	Columns []string
}
