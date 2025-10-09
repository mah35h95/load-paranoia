package gcp

import (
	"context"
	"fmt"

	"load_paranoia/model"

	"cloud.google.com/go/bigquery"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
)

// BqClient - Holds BigQuery client and context
type BqClient struct {
	ctx    context.Context
	client *bigquery.Client
}

// NewBigQueryClient - Creates and returns a new BigQuery client with context
func NewBigQueryClient(project string) (*BqClient, error) {
	ctx := context.Background()

	if project == "" {
		return nil, errors.New("project ID is empty")
	}
	client, err := bigquery.NewClient(ctx, project)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to create new bigquery client for %s", project)
	}

	return &BqClient{
		ctx:    ctx,
		client: client,
	}, nil
}

// CloseBigQueryClient - Closes the BigQuery client
func (bq *BqClient) CloseBigQueryClient() {
	bq.client.Close()
}

// RunQuery - Run the query provided in Big Query
func (bq *BqClient) RunIntervalRowCountQuery(projectID, datasetID, tableID, from, to string) []model.IntervalRowCountResult {
	rowCountIntervals := []model.IntervalRowCountResult{}

	queryString := fmt.Sprintf(`SELECT
  UNIX_MILLIS(TIMESTAMP_SECONDS(DIV(UNIX_SECONDS(recordstamp), 900) * 900 )) AS timestamp,
  COUNT(*) AS effectedRowCount
FROM
  %s
WHERE
  recordstamp >= "%s"
  AND recordstamp < "%s"
GROUP BY
  timestamp
ORDER BY
  timestamp desc;`,
		fmt.Sprintf("`%s.%s.%s`", projectID, datasetID, tableID),
		from,
		to,
	)
	query := bq.client.Query(queryString)

	it, err := query.Read(bq.ctx)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return rowCountIntervals
	}

	for {
		rowCountInterval := model.IntervalRowCountResult{}
		err := it.Next(&rowCountInterval)
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Println("Error iterating result:", err)
		}

		rowCountIntervals = append(rowCountIntervals, rowCountInterval)
	}

	return rowCountIntervals
}
