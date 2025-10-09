package model

import "cloud.google.com/go/bigquery"

type IntervalRowCountResult struct {
	Timestamp        bigquery.NullInt64
	EffectedRowCount bigquery.NullInt64
}
