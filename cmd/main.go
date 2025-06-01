package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/pryazhnikov/gofileschecker/internal/scanner"
)

type runParameters struct {
	path string
	// Future parameters can be easily added here
}

func parseParameters() *runParameters {
	params := &runParameters{}

	flag.StringVar(&params.path, "path", "", "Path to directory for scanning")

	flag.Parse()

	// Validate required parameters
	if params.path == "" {
		fmt.Fprintln(os.Stderr, "Error: path parameter is required")
		flag.Usage()
		os.Exit(1)
	}

	return params
}

func main() {
	params := parseParameters()

	scanner := scanner.NewDirectoryScanner(params.path)
	err := scanner.Scan()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Now you can use params.path for directory scanning
	fmt.Printf("Starting directory scan at: %s\n", params.path)
}
