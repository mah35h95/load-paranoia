package utils

import "load_paranoia/model"

// ChunkJobs - chunks jobs to a set chunk size
func ChunkJobs(jobIDs []model.TableDetails, chunkSize int) [][]model.TableDetails {
	chunkJobIDs := [][]model.TableDetails{}

	for i := 0; i < len(jobIDs); i += chunkSize {
		if i+chunkSize <= len(jobIDs) {
			chunkJobIDs = append(chunkJobIDs, jobIDs[i:i+chunkSize])
			continue
		}
		chunkJobIDs = append(chunkJobIDs, jobIDs[i:])
	}

	return chunkJobIDs
}
