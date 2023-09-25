package main

import (
	"flag"
	"fmt"
	"os"

	iwrapper "github.com/mazrean/iwrapper/internal"
)

var (
	src, dst string
)

func init() {
	flag.StringVar(&src, "src", "", "source file path")
	flag.StringVar(&dst, "dst", "", "destination file path")
}

func main() {
	flag.Parse()

	f, err := os.Open(src)
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

	f, err = os.Create(dst)
	if err != nil {
		panic(fmt.Errorf("failed to create destination file: %w", err))
	}
	defer f.Close()

	if err := iwrapper.Generate(f, pkgName, confs); err != nil {
		panic(fmt.Errorf("failed to generate wrapper: %w", err))
	}
}
