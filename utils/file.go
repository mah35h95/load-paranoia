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
			entry.Timestamp.Unix(),
			entry.ProtoPayload.ServiceData.JobGetQueryResultsResponse.Job.JobStatistics.QueryOutputRowCount,
		)
	}

	return combined
}
