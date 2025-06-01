package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/pryazhnikov/gofileschecker/internal/checkers"
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

	log.Printf("Starting directory scan at: %s\n", params.path)

	fileChecker := checkers.NewFileChecker()

	scanner := scanner.NewDirectoryScanner(params.path, fileChecker)
	err := scanner.Scan()
	if err != nil {
		log.Fatalf("Cannot scan directory: %v", err)
	}

	log.Println("Directory scan completed, getting the results...")

	fcg := fileChecker.GetDuplicatedFileGroups()
	if len(fcg) == 0 {
		log.Println("No duplicated files found")
		return
	}

	fmt.Printf("Found %d duplicated files groups\n", len(fcg))
	for _, fcg := range fcg {
		fmt.Printf("Duplicated files group: %s\n", fcg.Hash())
		for _, file := range fcg.Files() {
			fmt.Printf("- %s\n", file)
		}

		fmt.Println()
	}

	log.Println("Done")
}
