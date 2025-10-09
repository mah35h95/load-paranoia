package utils

// ChunkJobs - chunks jobs to a set chunk size
func ChunkJobs(jobIDs []string, chunkSize int) [][]string {
	chunkJobIDs := [][]string{}

	for i := 0; i < len(jobIDs); i += chunkSize {
		if i+chunkSize <= len(jobIDs) {
			chunkJobIDs = append(chunkJobIDs, jobIDs[i:i+chunkSize])
			continue
		}
		chunkJobIDs = append(chunkJobIDs, jobIDs[i:])
	}

	return chunkJobIDs
}
