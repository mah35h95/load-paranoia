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

	projectID := "prod-2134-entdatalake-5938ee"
	datasetID := "sap_s4_p41_lake"
	allTableIDs := []string{"afko_current_v1"}

	logProjectID := "prod-2763-entdatawh-bb5597"

	from := time.Now().AddDate(0, 0, -2).Format(time.RFC3339)
	to := time.Now().AddDate(0, 0, 1).Format(time.RFC3339)

	chunkTableIDs := utils.ChunkJobs(allTableIDs, chunkSize)

	for i := range chunkTableIDs {
		tableIDs := chunkTableIDs[i]

		fmt.Println("Fetching Access Token...")
		assesBearer := auth.GetAccessToken()

		wg := sync.WaitGroup{}
		wg.Add(len(tableIDs))

		for j := range tableIDs {
			tableID := tableIDs[j]
			fmt.Printf("(%d/%d): %s - Start\n", (chunkSize*i)+j+1, len(allTableIDs), tableID)

			go func() {
				defer wg.Done()
				tableLogs := gcp.GetTableResultLogs(logProjectID, projectID, datasetID, tableID, assesBearer, from, to)
				data := utils.CombineQueryOutputRowCount(tableLogs)

				err := utils.WriteToFile(
					fmt.Sprintf("./output/logs_%s.csv", tableID),
					[]byte(data),
				)
				if err != nil {
					fmt.Println("Error writing file:", err)
					return
				}
				fmt.Printf("Data written successfully for %s\n", tableID)

				fmt.Printf("(%d/%d): %s - Complete\n", (chunkSize*i)+j+1, len(allTableIDs), tableID)
			}()
		}

		wg.Wait()
	}

	utils.PrintStuff()
}
