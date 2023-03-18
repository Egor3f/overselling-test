package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/egor3f/overselling-test/ram"
	"os"
)

var testRam bool
var remainFreeMb int64
var allocateMb int64

func main() {
	err := parseArgs()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error occured and program stopped executing: \n%v", err)
		return
	}
	if testRam {
		ram.RunTest(remainFreeMb, allocateMb)
	}
}

func parseArgs() (err error) {
	flag.BoolVar(&testRam, "ram", false, "Test RAM overselling")
	flag.Int64Var(&remainFreeMb, "remfree", 0, "How much memory to remain free (Megabytes)")
	flag.Int64Var(&allocateMb, "alloc", 0, "How much memory to allocate (Megabytes)")
	flag.Parse()

	if !testRam {
		err = errors.New("Please, choose at least one test")
	}

	if remainFreeMb > 0 && allocateMb > 0 {
		err = errors.New("Use only one of these args: remfree or alloc")
	}

	return
}
