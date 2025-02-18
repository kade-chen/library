package shufflealgorithm

import (
	"math/rand"
	"time"
)

// Fisher-Yates 洗牌算法，用于随机打乱数组,保证每个元素出现的概率相等
// This function takes an array of strings as an argument and shuffles the elements in the array using the Fisher-Yates algorithm
func ShuffleAlgorithm(arr []string) {
	// 创建局部随机数生成器
	// Create a new random number generator using the current time in nanoseconds as the seed
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	n := len(arr)
	for i := n - 1; i > 0; i-- {
		j := r.Intn(i + 1)
		arr[i], arr[j] = arr[j], arr[i]
	}
}
