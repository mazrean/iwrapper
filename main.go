package main

import (
	"flag"
	"fmt"
	"os"

	iwrapper "github.com/mazrean/iwrapper/internal"
)

var (
	version          = "Unknown"
	revision         = "Unknown"
	versionFlag      bool
	srcFlag, dstFlag string
)

func init() {
	flag.BoolVar(&versionFlag, "version", false, "show version")
	flag.StringVar(&srcFlag, "src", "", "source file path")
	flag.StringVar(&dstFlag, "dst", "", "destination file path")
}

func main() {
	flag.Parse()

	if versionFlag {
		fmt.Printf("iwrapper %s (revision: %s)\n", version, revision)
		return
	}

	if len(srcFlag) == 0 {
		panic("source file path is required")
	}
	if len(dstFlag) == 0 {
		panic("destination file path is required")
	}

	f, err := os.Open(srcFlag)
	if err != nil {
		panic(fmt.Errorf("failed to open source file: %w", err))
	}
	defer f.Close()

	pkgName, results, err := iwrapper.ParseTarget(f)
	if err != nil {
		panic(fmt.Errorf("failed to parse target: %w", err))
	}

	confs, err := iwrapper.Convert(results)
	if err != nil {
		panic(fmt.Errorf("failed to convert: %w", err))
	}

	f, err = os.Create(dstFlag)
	if err != nil {
		panic(fmt.Errorf("failed to create destination file: %w", err))
	}
	defer f.Close()

	if err := iwrapper.Generate(f, pkgName, confs); err != nil {
		panic(fmt.Errorf("failed to generate wrapper: %w", err))
	}
}
