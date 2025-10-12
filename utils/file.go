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
			OutputRowCount: entry.ProtoPayload.ServiceData.JobGetQueryResultsResponse.Job.JobStatistics.QueryOutputRowCount,
			TimestampFrom:  fromTime,
			TimestampTo:    toTime,
			StartTime:      entry.ProtoPayload.ServiceData.JobGetQueryResultsResponse.Job.JobStatistics.StartTime,
			EndTime:        entry.ProtoPayload.ServiceData.JobGetQueryResultsResponse.Job.JobStatistics.EndTime,
		}
		queryLogs = append(queryLogs, queryLog)
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

func CombineQueryOutputRowCount(queryLogs []model.QueryLog) string {
	combined := "OutputRowCount,EpochTimestampFrom,EpochTimestampTo,EpochStartTime,EpochEndTime,TimestampFrom,TimestampTo,StartTime,EndTime\n"

	for _, queryLog := range queryLogs {
		combined += fmt.Sprintf(
			"%s,%d,%d,%d,%d,%s,%s,%s,%s\n",
			queryLog.OutputRowCount,
			queryLog.TimestampFrom.UnixMilli(),
			queryLog.TimestampTo.UnixMilli(),
			queryLog.StartTime.UnixMilli(),
			queryLog.EndTime.UnixMilli(),
			time.UnixMilli(queryLog.TimestampFrom.UnixMilli()).UTC().Format(time.DateTime),
			time.UnixMilli(queryLog.TimestampTo.UnixMilli()).UTC().Format(time.DateTime),
			time.UnixMilli(queryLog.StartTime.UnixMilli()).UTC().Format(time.DateTime),
			time.UnixMilli(queryLog.EndTime.UnixMilli()).UTC().Format(time.DateTime),
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
