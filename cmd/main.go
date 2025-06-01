package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/pryazhnikov/gofileschecker/internal/checkers"
	"github.com/pryazhnikov/gofileschecker/internal/scanner"
	"github.com/rs/zerolog"
)

type runParameters struct {
	path  string
	debug bool
	// Future parameters can be easily added here
}

func parseParameters() *runParameters {
	params := &runParameters{}

	flag.StringVar(&params.path, "path", "", "Path to directory for scanning")
	flag.BoolVar(&params.debug, "debug", false, "Enable debug logging")

	flag.Parse()

	// Validate required parameters
	if params.path == "" {
		fmt.Fprintln(os.Stderr, "Error: path parameter is required")
		flag.Usage()
		os.Exit(1)
	}

	return params
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
	params := parseParameters()
	logger := newLogger(params.debug)

	logger.Info().Str("path", params.path).Msg("Starting directory scan")

	fileChecker := checkers.NewFileChecker()

	scanner := scanner.NewDirectoryScanner(logger, params.path, fileChecker)
	err := scanner.Scan()
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
		fmt.Printf("Duplicated files group: %s\n", fcg.Hash())
		for _, file := range fcg.Files() {
			fmt.Printf("- %s\n", file)
		}

		fmt.Println()
	}

	logger.Info().Msg("Done")
}
