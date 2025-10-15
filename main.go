package main

import (
	"fmt"
	"sync"
	"time"

	"load_paranoia/auth"
	"load_paranoia/gcp"
	"load_paranoia/model"
	"load_paranoia/utils"
)

func main() {
	fmt.Println("Paranoia eradicator starting...")
	dbtProject := "dev-2763-entdatawh-591612"
	keyFilePath := "./.vscode/dbt-prod.json"

	chunkSize := 5

	dbtProjectID := "prod-2134-entdatalake-5938ee"
	dbtDatasetID := "sap_s4_p41_lake"

	stageProjectID := "prod-2434-entdataingest-05104f"
	stageDatasetID := "S4HANA"

	tableDetails := []model.TableDetails{
		{
			TableID: "bseg",
			Columns: []string{"mandt", "bukrs", "belnr", "gjahr", "buzei"},
		},
		{
			TableID: "matdoc",
			Columns: []string{"mandt", "key1", "key2", "key3", "key4", "key5", "key6"},
		},
		{
			TableID: "mldoc",
			Columns: []string{"mandt", "docref", "curtp"},
		},
	}

	logProjectID := "prod-2763-entdatawh-bb5597"

	// from := time.Now().AddDate(0, 0, -11).Format(time.RFC3339)
	// to := time.Now().AddDate(0, 0, -7).Format(time.RFC3339)

	fromParsed, err := time.Parse(time.RFC3339, "2025-09-30T21:00:00Z")
	if err != nil {
		fmt.Println("Error parsing from timestamp:", err)
		return
	}

	toParsed, err := time.Parse(time.RFC3339, "2025-10-07T21:00:00Z")
	if err != nil {
		fmt.Println("Error parsing to timestamp:", err)
		return
	}

	from := fromParsed.UTC().Format(time.RFC3339)
	to := toParsed.UTC().Format(time.RFC3339)

	chunkTableDetails := utils.ChunkJobs(tableDetails, chunkSize)

	for i := range chunkTableDetails {
		tableDetail := chunkTableDetails[i]

		fmt.Println("Fetching Access Token...")
		assesBearer := auth.GetAccessToken()

		bqClient, err := gcp.NewBigQueryClient(dbtProject, keyFilePath)
		if err != nil {
			fmt.Println("BQ Client Failed:", err)
			return
		}
		defer bqClient.CloseBigQueryClient()

		wg := sync.WaitGroup{}
		for j := range tableDetail {
			table := tableDetail[j]
			fmt.Printf("(%d/%d): %s - Start Log & BQ\n", (chunkSize*i)+j+1, len(tableDetails), table.TableID)

			wg.Go(func() {
				// GCP Logs Data
				dbtTableID := fmt.Sprintf("%s_current_v1__dbt_tmp", table.TableID)
				tableLogs := gcp.GetTableResultLogs(
					logProjectID,
					dbtProjectID,
					dbtDatasetID,
					dbtTableID,
					assesBearer,
					from,
					to,
				)
				tableQueryLogs := utils.GetQueryLogs(tableLogs)

				// BQ Table Data
				chunkQueries := gcp.GetChunkedQueries(stageProjectID, stageDatasetID, table, tableQueryLogs)
				tableIntervals := []model.IntervalRowCountResult{}
				for index, chunkQuery := range chunkQueries {
					tableIntervals = append(tableIntervals, bqClient.RunIntervalRowCountQuery(chunkQuery)...)
					fmt.Printf("%s: Fetched query result %d times\n", table.TableID, index+1)
				}

				// Write to csv
				data := utils.CombineRowCount(tableQueryLogs, tableIntervals)
				err := utils.WriteToFile(
					fmt.Sprintf("./output/%s.csv", table.TableID),
					[]byte(data),
				)
				if err != nil {
					fmt.Println("Error writing file:", err)
					return
				}
				fmt.Printf("Data written successfully for %s\n", table.TableID)

				fmt.Printf("(%d/%d): %s - Completed Log & BQ\n", (chunkSize*i)+j+1, len(tableDetails), table.TableID)
			})
		}

		wg.Wait()
	}

	utils.CombineAllCSVIntoExcel("./output")
	utils.PrintStuff()
}
