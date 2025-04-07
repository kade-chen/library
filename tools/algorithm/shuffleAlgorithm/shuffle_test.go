package shufflealgorithm_test

import (
	"fmt"
	"testing"

	shufflealgorithm "github.com/kade-chen/library/tools/algorithm/shuffleAlgorithm"
)


// This function tests the shuffle algorithm
func TestExampleShuffle(t *testing.T) {
	// Create an array of strings
	zones := []string{"us-east1-a", "us-east1-b", "us-east1-c", "us-west1-a", "us-west1-b"}

	// Print the original array
	fmt.Println("原始数组:", zones)
	// Call the shuffle algorithm
	shufflealgorithm.ShuffleAlgorithm(zones)

	// Print the shuffled array
	fmt.Println("洗牌后数组:", zones)
}
