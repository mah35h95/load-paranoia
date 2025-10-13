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
	diceProject := "prod-2367-entdataingst-7010d5"

	chunkSize := 5

	dbtProjectID := "prod-2134-entdatalake-5938ee"
	dbtDatasetID := "sap_s4_p41_lake"

	stageProjectID := "prod-2434-entdataingest-05104f"
	stageDatasetID := "S4HANA"

	tableDetails := []model.TableDetails{
		{
			TableID: "afko",
			Columns: []string{"mandt", "aufnr"},
		},
	}

	logProjectID := "prod-2763-entdatawh-bb5597"

	from := time.Now().AddDate(0, 0, -1).Format(time.RFC3339)
	to := time.Now().AddDate(0, 0, 0).Format(time.RFC3339)

	chunkTableDetails := utils.ChunkJobs(tableDetails, chunkSize)

	for i := range chunkTableDetails {
		tableDetail := chunkTableDetails[i]

		fmt.Println("Fetching Access Token...")
		assesBearer := auth.GetAccessToken()

		bqClient, err := gcp.NewBigQueryClient(diceProject)
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
				logData := utils.CombineQueryOutputRowCount(tableQueryLogs)

				err := utils.WriteToFile(
					fmt.Sprintf("./output/%s_log.csv", table.TableID),
					[]byte(logData),
				)
				if err != nil {
					fmt.Println("Error writing file:", err)
					return
				}
				fmt.Printf("Data written successfully for log-%s\n", table.TableID)

				// BQ Table Data
				chunkQueries := gcp.GetChunkedQueries(stageProjectID, stageDatasetID, table, tableQueryLogs)
				tableIntervals := []model.IntervalRowCountResult{}
				for index, chunkQuery := range chunkQueries {
					tableIntervals = append(tableIntervals, bqClient.RunIntervalRowCountQuery(chunkQuery)...)
					fmt.Printf("%s: Fetched query result %d times\n", table.TableID, index+1)
				}
				queryData := utils.CombineRowIntervalCount(tableIntervals)

				err = utils.WriteToFile(
					fmt.Sprintf("./output/%s_bq.csv", table.TableID),
					[]byte(queryData),
				)
				if err != nil {
					fmt.Println("Error writing file:", err)
					return
				}
				fmt.Printf("Data written successfully for bq-%s\n", table.TableID)

				fmt.Printf("(%d/%d): %s - Completed Log & BQ\n", (chunkSize*i)+j+1, len(tableDetails), table.TableID)
			})
		}

		wg.Wait()
	}

	utils.PrintStuff()
}
