package utils

import (
	"testing"
)

func TestGenerateUniqueID(t *testing.T) {
	t.Parallel()

	// Create a map to store generated IDs and check for collisions
	uniqueIDMap := make(map[string]bool)

	// Generate a large number of IDs
	numIDs := 1000

	for range numIDs {
		uniqueID := GenerateUniqueID()

		// Check if the ID is of the correct length
		if len(uniqueID) != 36 {
			t.Errorf("Generated ID has incorrect length: %s", uniqueID)
		}

		// Check for collisions
		if uniqueIDMap[uniqueID] {
			t.Errorf("Collision detected for ID: %s", uniqueID)
		}

		// Add the generated ID to the map
		uniqueIDMap[uniqueID] = true
	}
}
