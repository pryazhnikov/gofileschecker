package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/pryazhnikov/gofileschecker/internal/checkers"
	"github.com/pryazhnikov/gofileschecker/internal/scanner"
	"github.com/rs/zerolog"
)

type runParameters struct {
	path            string // Path to directory for scanning
	debug           bool   // Enable debug logging
	showGroupPrefix bool   // Show common path prefix for the found groups
}

func parseParameters() (*runParameters, error) {
	params := &runParameters{}

	flag.StringVar(&params.path, "path", "", "Path to directory for scanning")
	flag.BoolVar(&params.debug, "debug", false, "Enable debug logging")
	flag.BoolVar(&params.showGroupPrefix, "prefix", false, "Show common path prefix for file groups")

	flag.Parse()

	// Validate required parameters
	if params.path == "" {
		return nil, fmt.Errorf("path parameter is required")
	}

	return params, nil
}

func newLogger(debug bool) zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	level := zerolog.InfoLevel
	if debug {
		level = zerolog.DebugLevel
	}

	return logger.Level(level)
}

func main() {
	params, err := parseParameters()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
		flag.Usage()
		os.Exit(1)
	}

	logger := newLogger(params.debug)
	logger.Info().Msgf("Logger level set: %s", logger.GetLevel().String())

	fileChecker := checkers.NewFileChecker()
	scanner := scanner.NewDirectoryScanner(logger, fileChecker)

	// todo: add an ability to scan multiple directories
	err = scanner.Scan(params.path)
	if err != nil {
		logger.Fatal().Err(err).Msg("Cannot scan directory")
	}

	logger.Info().Msg("Directory scan completed, getting the results...")

	fcg := fileChecker.GetDuplicatedFileGroups()
	if len(fcg) == 0 {
		logger.Info().Msg("No duplicated files found")
		return
	}

	fmt.Printf("Found %d duplicated files groups\n", len(fcg))
	for _, fcg := range fcg {
		fmt.Printf(
			"Duplicated files group: %s\n",
			fcg.Hash(),
		)

		pathPrefix := ""
		if params.showGroupPrefix {
			pathPrefix = fcg.CommonPathPrefix()
			if pathPrefix != "" {
				fmt.Println(pathPrefix)
			}
		}

		for _, file := range fcg.Files() {
			relativePath := strings.TrimPrefix(file, pathPrefix)
			fmt.Printf("- %s\n", relativePath)
		}

		fmt.Println()
	}

	logger.Info().Msg("Done")
}
