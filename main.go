package main

import (
	"fmt"
	"sync"
	"time"

	"load_paranoia/auth"
	"load_paranoia/gcp"
	"load_paranoia/utils"
)

func main() {
	fmt.Println("Paranoia eradicator starting...")

	chunkSize := 5

	dbtProjectID := "prod-2134-entdatalake-5938ee"
	dbtDatasetID := "sap_s4_p41_lake"

	stageProjectID := "prod-2434-entdataingest-05104f"
	stageDatasetID := "S4HANA"

	dbtTableIDs := []string{"afko"}

	logProjectID := "prod-2763-entdatawh-bb5597"

	from := time.Now().AddDate(0, 0, -2).Format(time.RFC3339)
	to := time.Now().AddDate(0, 0, 1).Format(time.RFC3339)

	chunkTableIDs := utils.ChunkJobs(dbtTableIDs, chunkSize)

	for i := range chunkTableIDs {
		tableIDs := chunkTableIDs[i]

		fmt.Println("Fetching Access Token...")
		assesBearer := auth.GetAccessToken()

		bqClient, err := gcp.NewBigQueryClient(stageProjectID)
		if err != nil {
			fmt.Println("BQ Client Failed:", err)
			return
		}
		defer bqClient.CloseBigQueryClient()

		wg := sync.WaitGroup{}
		for j := range tableIDs {
			tableID := tableIDs[j]
			fmt.Printf("(%d/%d): %s - Start Log & BQ\n", (chunkSize*i)+j+1, len(dbtTableIDs), tableID)

			wg.Go(func() {
				dbtTableID := fmt.Sprintf("%s_current_v1", tableID)
				tableLogs := gcp.GetTableResultLogs(
					logProjectID,
					dbtProjectID,
					dbtDatasetID,
					dbtTableID,
					assesBearer,
					from,
					to,
				)
				data := utils.CombineQueryOutputRowCount(tableLogs)

				err := utils.WriteToFile(
					fmt.Sprintf("./output/%s_log.csv", tableID),
					[]byte(data),
				)
				if err != nil {
					fmt.Println("Error writing file:", err)
					return
				}
				fmt.Printf("Data written successfully for log-%s\n", tableID)

				fmt.Printf("(%d/%d): %s - Complete Log\n", (chunkSize*i)+j+1, len(dbtTableIDs), tableID)
			})

			wg.Go(func() {
				tableInterval := bqClient.RunIntervalRowCountQuery(stageProjectID, stageDatasetID, tableID, from, to)
				data := utils.CombineRowIntervalCount(tableInterval)

				err := utils.WriteToFile(
					fmt.Sprintf("./output/%s_bq.csv", tableID),
					[]byte(data),
				)
				if err != nil {
					fmt.Println("Error writing file:", err)
					return
				}
				fmt.Printf("Data written successfully for bq-%s\n", tableID)

				fmt.Printf("(%d/%d): %s - Complete BQ\n", (chunkSize*i)+j+1, len(dbtTableIDs), tableID)
			})
		}

		wg.Wait()
	}

	utils.PrintStuff()
}
