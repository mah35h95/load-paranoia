package utils

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"load_paranoia/model"
)

func WriteToFile(name string, data []byte) error {
	return os.WriteFile(name, data, 0644)
}

func GetQueryLogs(entries []model.Entry) []model.QueryLog {
	queryLogs := []model.QueryLog{}

	for _, entry := range entries {
		fromTime, toTime := extractFromAndToTimestamp(entry.ProtoPayload.ServiceData.JobGetQueryResultsResponse.Job.JobConfiguration.Query.Query)

		queryLog := model.QueryLog{
			JobID:          entry.ProtoPayload.ServiceData.JobGetQueryResultsResponse.Job.JobName.JobID,
			OutputRowCount: entry.ProtoPayload.ServiceData.JobGetQueryResultsResponse.Job.JobStatistics.QueryOutputRowCount,
			From:           fromTime,
			To:             toTime,
			TimestampFrom:  fromTime,
			TimestampTo:    toTime,
			StartTime:      entry.ProtoPayload.ServiceData.JobGetQueryResultsResponse.Job.JobStatistics.StartTime,
			EndTime:        entry.ProtoPayload.ServiceData.JobGetQueryResultsResponse.Job.JobStatistics.EndTime,
		}
		queryLogs = append(queryLogs, queryLog)
	}

	for i := 0; i < len(queryLogs)-1; i++ {
		if queryLogs[i].From.UnixMicro() < queryLogs[i+1].To.UnixMicro() {
			queryLogs[i+1].To = queryLogs[i].From
		}
	}

	return queryLogs
}

func extractFromAndToTimestamp(query string) (time.Time, time.Time) {
	fromRe := regexp.MustCompile(`where +recordstamp +> +'.*' +and`)
	fromMatch := fromRe.FindString(query)
	fromTime := extractTimestamp(fromMatch)

	toRe := regexp.MustCompile(`and +recordstamp +<= +'.*'`)
	toMatch := toRe.FindString(query)
	toTime := extractTimestamp(toMatch)

	return fromTime, toTime
}

func extractTimestamp(match string) time.Time {
	timestamp := time.Time{}

	timestampRe := regexp.MustCompile(`'.*'`)
	timestampMatch := strings.ReplaceAll(timestampRe.FindString(match), "'", "")

	customLayout := "2006-01-02 15:04:05.999999-07:00"
	timestamp, err := time.Parse(customLayout, timestampMatch)
	if err != nil {
		fmt.Printf("Failed Parsing Time: %+v\n", err)
		return timestamp
	}

	return timestamp
}

func CombineRowCount(queryLogs []model.QueryLog, intervalCounts []model.IntervalRowCountResult) string {
	combined := "SLTLoadedRowCount,EffectedOutputRowCount,FromEpochMicro,ToEpochMicro,JobID,QueryStartEpochMicro,QueryEndEpochMicro,QueryFromEpochMicro,QueryToEpochMicro\n"

	for i := range queryLogs {
		combined += fmt.Sprintf(
			"%d,%s,%d,%d,%s,%d,%d,%d,%d\n",
			intervalCounts[i].EffectedRowCount.Int64,
			queryLogs[i].OutputRowCount,
			queryLogs[i].From.UnixMicro(),
			queryLogs[i].To.UnixMicro(),
			queryLogs[i].JobID,
			queryLogs[i].StartTime.UnixMicro(),
			queryLogs[i].EndTime.UnixMicro(),
			queryLogs[i].TimestampFrom.UnixMicro(),
			queryLogs[i].TimestampTo.UnixMicro(),
		)
	}

	return combined
}
