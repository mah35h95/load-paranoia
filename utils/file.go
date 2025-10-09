package utils

import (
	"fmt"
	"os"
	"time"

	"load_paranoia/model"
)

func WriteToFile(name string, data []byte) error {
	return os.WriteFile(name, data, 0644)
}

func CombineQueryOutputRowCount(entries []model.Entry) string {
	combined := "epochMillis,timestamp,queryOutputRowCount\n"

	for _, entry := range entries {
		combined += fmt.Sprintf(
			"%d,%s,%s\n",
			entry.Timestamp.UnixMilli(),
			time.UnixMilli(entry.Timestamp.UnixMilli()).UTC().Format(time.DateTime),
			entry.ProtoPayload.ServiceData.JobGetQueryResultsResponse.Job.JobStatistics.QueryOutputRowCount,
		)
	}

	return combined
}

func CombineRowIntervalCount(intervalCounts []model.IntervalRowCountResult) string {
	combined := "epochMillis,timestamp,effectedRowCount\n"

	for _, intervalCount := range intervalCounts {
		combined += fmt.Sprintf(
			"%d,%s,%d\n",
			intervalCount.Timestamp.Int64,
			time.UnixMilli(intervalCount.Timestamp.Int64).UTC().Format(time.DateTime),
			intervalCount.EffectedRowCount.Int64,
		)
	}

	return combined
}
