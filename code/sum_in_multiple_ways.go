package code

import (
	"runtime"
	"sync"
)

// 典型的顺序累加求和
func sumSequential(nums []int) int64 {
	var total int64 = 0 // 累加器 total
	for _, n := range nums {
		total += int64(n)
	}
	return total
}

// 分块并行求和
func sumParallelChunks(nums []int, numChunks int) int64 {
	if len(nums) == 0 {
		return 0
	}
	if numChunks <= 0 {
		numChunks = runtime.NumCPU()
	} // 默认使用CPU核心数作为块数
	if len(nums) < numChunks {
		numChunks = len(nums)
	}

	results := make(chan int64, numChunks)
	chunkSize := (len(nums) + numChunks - 1) / numChunks

	for i := 0; i < numChunks; i++ {
		start := i * chunkSize
		end := (i + 1) * chunkSize
		if end > len(nums) {
			end = len(nums)
		}

		// 每个goroutine处理一个独立的块
		go func(chunk []int) {
			var localSum int64 = 0
			for _, n := range chunk { // 块内部仍然是顺序累加，但这是局部行为
				localSum += int64(n)
			}
			results <- localSum // 将局部结果发送到channel
		}(nums[start:end])
	}

	var total int64 = 0
	for i := 0; i < numChunks; i++ {
		total += <-results // 合并结果，加法是结合的，顺序不重要
	}
	return total
}

// 辅助函数
func sumRecursiveParallelEntry(nums []int) int64 {
	// 设定一个阈值，小于此阈值则顺序计算，避免过多goroutine开销
	const threshold = 1024
	return sumRecursiveParallel(nums, threshold)
}

// 递归分治的并行求和
func sumRecursiveParallel(nums []int, threshold int) int64 {
	if len(nums) == 0 {
		return 0
	}
	if len(nums) < threshold {
		return sumSequential(nums) // 小任务直接顺序计算
	}

	mid := len(nums) / 2

	var sumLeft int64
	var wg sync.WaitGroup
	wg.Add(1) // 我们需要等待左半部分的计算结果
	go func() {
		defer wg.Done()
		sumLeft = sumRecursiveParallel(nums[:mid], threshold)
	}()

	// 右半部分可以在当前goroutine计算，也可以再开一个goroutine
	sumRight := sumRecursiveParallel(nums[mid:], threshold)

	wg.Wait() // 等待左半部分完成

	return sumLeft + sumRight // 合并，加法是结合的
}
