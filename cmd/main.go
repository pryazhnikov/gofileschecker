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
	path           string // Path to directory for scanning
	debug          bool   // Enable debug logging
	fullFilePath   bool   // Show full file paths in output
	skipEmptyFiles bool   // Do not process empty files
}

func parseParameters() (*runParameters, error) {
	params := &runParameters{}

	flag.StringVar(&params.path, "path", "", "Path to directory for scanning")
	flag.BoolVar(&params.debug, "debug", false, "Enable debug logging")
	flag.BoolVar(&params.fullFilePath, "fullpath", false, "Show full file paths in output")
	flag.BoolVar(&params.skipEmptyFiles, "skipempty", false, "Skip empty files during scanning")

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

	fileChecker := checkers.NewFileChecker(params.skipEmptyFiles)
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

		pathPrefix := fcg.CommonPathPrefix()
		fmt.Printf("Location: %s\n", pathPrefix)

		for _, file := range fcg.Files() {
			fileView := file
			if !params.fullFilePath {
				fileView = strings.TrimPrefix(file, pathPrefix)
			}

			fmt.Printf("- %s\n", fileView)
		}

		fmt.Println()
	}

	logger.Info().Msg("Done")
}
