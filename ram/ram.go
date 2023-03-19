package ram

import (
	"bytes"
	"fmt"
	"github.com/pbnjay/memory"
	"log"
	"os"
	"runtime/debug"
	"time"
)

var logger *log.Logger
var buffers []*bytes.Buffer

func init() {
	logger = log.New(os.Stdout, "[RAM] ", 0)
}

func RunTest(remainFree int64, alloc int64) {
	const (
		defaultRemainFree = 50 * 1024 * 1024
	)
	logger.Println("Started")

	pageSize := int64(os.Getpagesize()) // assuming 4096
	chunkSize := int64(1024 * 1024)
	freeMemoryBeforeTest := int64(memory.FreeMemory())
	var toAllocate, totalAllocated int64
	if alloc > 0 {
		toAllocate = alloc
	} else {
		if remainFree <= 0 {
			remainFree = defaultRemainFree
		}
		toAllocate = absInt(freeMemoryBeforeTest - remainFree)
	}
	logger.Printf(
		"Page size: %v, chunk size: %v, free memory: %v, to allocate: %v",
		formatMemory(pageSize),
		formatMemory(chunkSize),
		formatMemory(freeMemoryBeforeTest),
		formatMemory(toAllocate),
	)

	startTime := time.Now()
	for totalAllocated = 0; totalAllocated < toAllocate; totalAllocated += chunkSize {
		allocateRam(chunkSize)
	}
	testDuration := time.Since(startTime)

	freeMemoryAfterTest := int64(memory.FreeMemory())
	calculationError := totalAllocated - (freeMemoryBeforeTest - freeMemoryAfterTest)
	logger.Printf(
		"Free memory before test: %v, after test: %v, diff %v, calculation error %v",
		formatMemory(freeMemoryBeforeTest),
		formatMemory(freeMemoryAfterTest),
		formatMemory(freeMemoryBeforeTest-freeMemoryAfterTest),
		formatMemory(calculationError),
	)

	releaseRam()
	freeMemoryAfterRelease := int64(memory.FreeMemory())
	logger.Printf(
		"Released RAM, free memory: %v, diff: %v",
		formatMemory(freeMemoryAfterRelease),
		formatMemory(freeMemoryAfterRelease-freeMemoryAfterTest),
	)
	logger.Printf(
		"Total test duration: %v, microsecs per Gb: %v",
		testDuration,
		float64(testDuration.Microseconds())/(float64(totalAllocated)/1024/1024/1024),
	)
}

func allocateRam(byteCount int64) (microsecs int64) {
	startTime := time.Now()
	buf := bytes.NewBuffer([]byte{})
	buf.Grow(int(byteCount))
	buffers = append(buffers, buf)
	return time.Since(startTime).Microseconds()
}

func releaseRam() {
	buffers = nil
	debug.FreeOSMemory()
}

func formatMemory(ms int64) string {
	const Gb = 1024 * 1024 * 1024
	const Mb = 1024 * 1024
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
