package utils

import (
	"fmt"
	"os"

	"load_paranoia/model"
)

func WriteToFile(name string, data []byte) error {
	return os.WriteFile(name, data, 0644)
}

func CombineQueryOutputRowCount(entries []model.Entry) string {
	combined := "timestamp,queryOutputRowCount\n"

	for _, entry := range entries {
		combined += fmt.Sprintf(
			"%d,%s\n",
			entry.Timestamp.UnixMilli(),
			entry.ProtoPayload.ServiceData.JobGetQueryResultsResponse.Job.JobStatistics.QueryOutputRowCount,
		)
	}

	return combined
}

func CombineRowIntervalCount(intervalCounts []model.IntervalRowCountResult) string {
	combined := "timestamp,effectedRowCount\n"

	for _, intervalCount := range intervalCounts {
		combined += fmt.Sprintf(
			"%d,%d\n",
			intervalCount.Timestamp.Int64,
			intervalCount.EffectedRowCount.Int64,
		)
	}

	return combined
}
