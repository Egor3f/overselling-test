package ram

import (
	"fmt"
	"github.com/pbnjay/memory"
	"log"
	"os"
	"runtime"
	"time"
)

type allocateType uint8

const Gb = 1024 * 1024 * 2014
const Mb = 1024 * 2014
const Kb = 1024

var logger *log.Logger
var slices [][]allocateType

func init() {
	logger = log.New(os.Stdout, "[RAM] ", 0)
}

func RunTest(remainFreeMb int64, allocMb int64) {
	const (
		defaultRemainFreeMb = 50
	)
	logger.Println("Started")

	pageSize := int64(os.Getpagesize())
	logger.Printf("Page size: %v", pageSize)

	freeMemoryBefore := int64(memory.FreeMemory())
	logger.Printf("Free memory before test: %v", freeMemoryBefore)

	var totalEstimate, toAllocate int64
	if allocMb > 0 {
		toAllocate = allocMb * Mb
	} else {
		if remainFreeMb <= 0 {
			remainFreeMb = defaultRemainFreeMb
		}
		toAllocate = absInt(freeMemoryBefore - remainFreeMb)
	}

	for totalEstimate = 0; totalEstimate < toAllocate; totalEstimate += int64(pageSize) {
		timeElapsed := allocateRam(pageSize)
		totalEstimate += pageSize
		logger.Printf("Allocated about %v, took %v microseconds", formatMemory(totalEstimate), formatMemory(timeElapsed))
	}

	freeMemoryAfterTest := int64(memory.FreeMemory())
	calculationError := absInt(totalEstimate - (freeMemoryBefore - freeMemoryAfterTest))
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
	newSlice := make([]allocateType, byteCount, byteCount)
	slices = append(slices, newSlice)
	return time.Since(startTime).Microseconds()
}

func releaseRam() {
	slices = nil
	runtime.GC()
}

func formatMemory(ms int64) string {
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
