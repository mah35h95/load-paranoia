package utils

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"load_paranoia/model"

	excelize "github.com/xuri/excelize/v2"
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

	cleanQueryLogs := []model.QueryLog{}
	for _, queryLog := range queryLogs {
		if queryLog.From != queryLog.To {
			cleanQueryLogs = append(cleanQueryLogs, queryLog)
		}
	}
	return cleanQueryLogs
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
	combined := "FromEpochMicro,ToEpochMicro,SLTLoadedRowCount,NewLakeLoadedRowCount,JobID,JobStartEpochMicro,JobEndEpochMicro\n"

	for i := range queryLogs {
		combined += fmt.Sprintf(
			"%s,%s,%d,%s,%s,%s,%s\n",
			time.UnixMicro(queryLogs[i].From.UnixMicro()).UTC().Format(time.DateTime),
			time.UnixMicro(queryLogs[i].To.UnixMicro()).UTC().Format(time.DateTime),
			intervalCounts[i].EffectedRowCount.Int64,
			queryLogs[i].OutputRowCount,
			queryLogs[i].JobID,
			time.UnixMicro(queryLogs[i].StartTime.UnixMicro()).UTC().Format(time.DateTime),
			time.UnixMicro(queryLogs[i].EndTime.UnixMicro()).UTC().Format(time.DateTime),
		)

		// combined += fmt.Sprintf(
		// 	"%d,%d,%d,%s,%s,%d,%d\n",
		// 	queryLogs[i].From.UnixMicro(),
		// 	queryLogs[i].To.UnixMicro(),
		// 	intervalCounts[i].EffectedRowCount.Int64,
		// 	queryLogs[i].OutputRowCount,
		// 	queryLogs[i].JobID,
		// 	queryLogs[i].TimestampFrom.UnixMicro(),
		// 	queryLogs[i].TimestampTo.UnixMicro(),
		// )
	}

	return combined
}

func CombineAllCSVIntoExcel(csvDir string) {
	// Create a new Excel file
	f := excelize.NewFile()

	// Get all CSV files in the directory
	files, err := os.ReadDir(csvDir)
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		return
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".csv") {
			csvFilePath := filepath.Join(csvDir, file.Name())
			sheetName := strings.TrimSuffix(file.Name(), ".csv") // Use CSV filename as sheet name

			// Create a new sheet in the Excel file
			index, err := f.NewSheet(sheetName)
			if err != nil {
				fmt.Printf("Error creating new sheet %s: %v\n", sheetName, err)
				continue
			}
			f.SetActiveSheet(index)

			// Open and read the CSV file
			csvFile, err := os.Open(csvFilePath)
			if err != nil {
				fmt.Printf("Error opening CSV file %s: %v\n", csvFilePath, err)
				continue
			}
			defer csvFile.Close()

			reader := csv.NewReader(csvFile)
			reader.FieldsPerRecord = -1 // Allow variable number of fields per record

			// Write CSV data to the Excel sheet
			rowNum := 1
			for {
				record, err := reader.Read()
				if err == io.EOF {
					break
				}
				if err != nil {
					fmt.Printf("Error reading CSV record from %s: %v\n", csvFilePath, err)
					break
				}

				// Write each cell in the row
				for colNum, cellValue := range record {
					cellRef, err := excelize.CoordinatesToCellName(colNum+1, rowNum)
					if err != nil {
						fmt.Printf("Error converting coordinates to cell name: %v\n", err)
						continue
					}
					f.SetCellValue(sheetName, cellRef, cellValue)
				}
				rowNum++
			}
			fmt.Printf("Successfully imported %s to sheet %s\n", file.Name(), sheetName)
		}
	}

	// Delete the default "Sheet1" if it's empty and not needed
	if f.GetSheetName(0) == "Sheet1" && len(f.GetSheetList()) > 1 {
		f.DeleteSheet("Sheet1")
	}

	// Save the Excel file
	outputFileName := fmt.Sprintf("%s/output_csvs.xlsx", csvDir)
	if err := f.SaveAs(outputFileName); err != nil {
		fmt.Printf("Error saving Excel file: %v\n", err)
		return
	}
	fmt.Printf("Excel file '%s' created successfully.\n", outputFileName)
}
