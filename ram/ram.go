package ram

import (
	"fmt"
	"github.com/pbnjay/memory"
	"log"
	"os"
	"runtime"
	"time"
)

type allocateType uint64

const allocateTypeSize = 8

var logger *log.Logger
var slices [][]allocateType

func init() {
	logger = log.New(os.Stdout, "[RAM] ", 0)
}

func RunTest(remainFree int64, alloc int64) {
	const (
		defaultRemainFree = 50 * 1024 * 1024
	)
	logger.Println("Started")

	pageSize := int64(os.Getpagesize())
	chunkSize := pageSize * 256
	logger.Printf("Page size: %v, chunk size: %v", pageSize, chunkSize)

	freeMemoryBefore := int64(memory.FreeMemory())
	logger.Printf("Free memory before test: %v", freeMemoryBefore)

	var toAllocate, totalAllocated int64
	if alloc > 0 {
		toAllocate = alloc
	} else {
		if remainFree <= 0 {
			remainFree = defaultRemainFree
		}
		toAllocate = absInt(freeMemoryBefore - remainFree)
	}

	for totalAllocated = 0; totalAllocated < toAllocate; totalAllocated += chunkSize {
		timeElapsed := allocateRam(chunkSize)
		totalAllocated += chunkSize
		logger.Printf("Allocated about %v, took %v microseconds", formatMemory(totalAllocated), formatMemory(timeElapsed))
	}

	freeMemoryAfterTest := int64(memory.FreeMemory())
	calculationError := absInt(totalAllocated - (freeMemoryBefore - freeMemoryAfterTest))
	logger.Printf(
		"Free memory before test: %v, after test: %v, allocated totally %v, calculation error %v",
		formatMemory(freeMemoryBefore),
		formatMemory(freeMemoryAfterTest),
		formatMemory(freeMemoryBefore-freeMemoryAfterTest),
		formatMemory(calculationError),
	)

	releaseRam()
	freeMemoryAfterRelease := int64(memory.FreeMemory())
	logger.Printf(
		"Released RAM, free memory: %v, diff: %v",
		formatMemory(freeMemoryAfterRelease),
		formatMemory(absInt(freeMemoryAfterRelease-freeMemoryAfterTest)),
	)
}

func allocateRam(byteCount int64) (microsecs int64) {
	startTime := time.Now()
	elemCount := byteCount / allocateTypeSize
	newSlice := make([]allocateType, elemCount, elemCount)
	slices = append(slices, newSlice)
	return time.Since(startTime).Microseconds()
}

func releaseRam() {
	slices = nil
	runtime.GC()
}

func formatMemory(ms int64) string {
	const Gb = 1024 * 1024 * 2014
	const Mb = 1024 * 2014
	const Kb = 1024

	switch {
	case ms > Gb:
		return fmt.Sprintf("%.3f Gb", float64(ms)/Gb)
	case ms > Mb:
		return fmt.Sprintf("%.3f Mb", float64(ms)/Mb)
	case ms > Kb:
		return fmt.Sprintf("%.3f Kb", float64(ms)/Kb)
	default:
		return fmt.Sprintf("%d", uint64(ms))
	}
}

func absInt(val int64) int64 {
	if val >= 0 {
		return val
	}
	return -val
}
