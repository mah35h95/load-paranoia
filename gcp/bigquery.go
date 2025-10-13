package gcp

import (
	"context"
	"fmt"
	"strings"

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
func (bq *BqClient) RunIntervalRowCountQuery(intervalQuery string) []model.IntervalRowCountResult {
	rowCountIntervals := []model.IntervalRowCountResult{}

	query := bq.client.Query(intervalQuery)
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

func GetChunkedQueries(projectID, datasetID string, table model.TableDetails, queryLogs []model.QueryLog) []string {
	chunkQueries := []string{}

	cteQueries := []string{}
	selectQueries := []string{}

	for index, queryLog := range queryLogs {
		cteQuery, selectQuery := getIntervalRowCountCteAndQuery(index, projectID, datasetID, table, queryLog)

		subQueries := len(cteQueries) + 1
		if subQueries > model.MaxSubQueries {
			chunkQueries = append(chunkQueries, getCombinedIntervalRowCountQuery(cteQueries, selectQueries))

			cteQueries = []string{cteQuery}
			selectQueries = []string{selectQuery}
			continue
		}

		cteQueries = append(cteQueries, cteQuery)
		selectQueries = append(selectQueries, selectQuery)
	}

	chunkQueries = append(chunkQueries, getCombinedIntervalRowCountQuery(cteQueries, selectQueries))
	return chunkQueries
}

func getIntervalRowCountCteAndQuery(
	index int,
	projectID,
	datasetID string,
	table model.TableDetails,
	queryLogs model.QueryLog,
) (string, string) {
	return fmt.Sprintf(
			`latest_records_%d AS (
SELECT
CASE
 WHEN operation_flag = 'D' THEN 1
 WHEN operation_flag = 'U' THEN 2
 WHEN operation_flag = 'I' THEN 3
ELSE 4
END AS operation_rank FROM %s a
WHERE recordstamp > TIMESTAMP_MICROS(%d) AND recordstamp <= TIMESTAMP_MICROS(%d)
QUALIFY ROW_NUMBER() OVER (PARTITION BY %s ORDER BY recordstamp DESC, operation_rank ASC) = 1 )`,
			index,
			fmt.Sprintf("`%s.%s.%s`", projectID, datasetID, table.TableID),
			queryLogs.TimestampFrom.UnixMicro(),
			queryLogs.TimestampTo.UnixMicro(),
			strings.Join(table.Columns, ","),
		), fmt.Sprintf(
			"SELECT COUNT(*) AS effectedRowCount, %d AS fromTimestamp, %d AS toTimestamp, FROM latest_records_%d\n",
			queryLogs.TimestampFrom.UnixMicro(),
			queryLogs.TimestampTo.UnixMicro(),
			index,
		)
}

func getCombinedIntervalRowCountQuery(cteQueries, selectQueries []string) string {
	return fmt.Sprintf(
		"WITH\n%s\n%sORDER BY fromTimestamp DESC",
		strings.Join(cteQueries, ","),
		strings.Join(selectQueries, "UNION ALL "),
	)
}
