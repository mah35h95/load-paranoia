package gcp

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"

	"load_paranoia/model"
)

func GetTableResultLogs(logProjectID, projectID, datasetID, tableID, bearer, from, to string) []model.Entry {
	allLogEntries := []model.Entry{}

	filter := fmt.Sprintf(
		"protoPayload.authenticationInfo.principalEmail=\"svc-dbt-worker@prod-2763-entdatawh-bb5597.iam.gserviceaccount.com\" AND protoPayload.methodName=\"jobservice.getqueryresults\" AND protoPayload.serviceData.jobGetQueryResultsResponse.job.jobConfiguration.query.destinationTable.projectId=\"%s\" AND protoPayload.serviceData.jobGetQueryResultsResponse.job.jobConfiguration.query.destinationTable.datasetId=\"%s\" AND protoPayload.serviceData.jobGetQueryResultsResponse.job.jobConfiguration.query.destinationTable.tableId=\"%s\" AND protoPayload.serviceData.jobGetQueryResultsResponse.job.jobStatistics.queryOutputRowCount:*  AND timestamp > \"%s\" AND timestamp < \"%s\"",
		projectID,
		datasetID,
		tableID,
		from,
		to,
	)

	pageToken := ""
	count := 1
	for {
		logEntries, nextPageToken := getLogEntries(logProjectID, filter, pageToken, bearer)
		allLogEntries = append(allLogEntries, logEntries...)

		fmt.Printf("%s: Fetched logs %d times\n", tableID, count)
		count++

		pageToken = nextPageToken
		if pageToken == "" {
			break
		}
	}

	return allLogEntries
}

func getLogEntries(logProjectID, filter, pageToken, bearer string) ([]model.Entry, string) {
	loggingResponce := model.LoggingResponce{}

	path := "https://logging.googleapis.com/v2/entries:list"

	loggingRequest := model.LoggingRequest{
		ProjectIDS:    []string{logProjectID},
		ResourceNames: []string{"projects/" + logProjectID},
		Filter:        filter,
		OrderBy:       "timestamp desc",
		PageSize:      math.MaxInt32,
		PageToken:     pageToken,
	}

	byteBody, err := json.Marshal(loggingRequest)
	if err != nil {
		fmt.Printf("JSON Marshal: %+v\n", err)
		return loggingResponce.Entries, ""
	}

	req, err := http.NewRequest(http.MethodPost, path, strings.NewReader(string(byteBody)))
	if err != nil {
		fmt.Printf("New Request Create: %+v\n", err)
		return loggingResponce.Entries, ""
	}

	req.Header = http.Header{
		"Authorization": {bearer},
		"Accept":        {"application/json"},
		"Content-Type":  {"application/json"},
	}

	// Send req using http Client
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("Http Do: %+v\n", err)
		return loggingResponce.Entries, ""
	}
	if res.StatusCode == 403 {
		fmt.Printf("Un-Authorized: %+v\n", err)
		return loggingResponce.Entries, ""
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Read Body: %+v\n", err)
		return loggingResponce.Entries, ""
	}

	err = json.Unmarshal(resBody, &loggingResponce)
	if err != nil {
		fmt.Printf("JSON Unmarshaling: %+v\n", err)
		fmt.Printf("%s\n", string(resBody))
		return loggingResponce.Entries, ""
	}

	return loggingResponce.Entries, loggingResponce.NextPageToken
}
