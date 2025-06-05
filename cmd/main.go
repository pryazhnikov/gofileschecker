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
	paths          []string // Paths to directories for scanning
	debug          bool     // Enable debug logging
	fullFilePath   bool     // Show full file paths in output
	skipEmptyFiles bool     // Do not process empty files
}

func parseParameters() (*runParameters, error) {
	params := &runParameters{
		paths: make([]string, 0),
	}

	flag.BoolVar(&params.debug, "debug", false, "Enable debug logging")
	flag.BoolVar(&params.fullFilePath, "fullpath", false, "Show full file paths in output")
	flag.BoolVar(&params.skipEmptyFiles, "skipempty", false, "Skip empty files during scanning")
	flag.Func("path", "Path to directory for scanning", func(flagValue string) error {
		params.paths = append(params.paths, flagValue)
		return nil
	})

	flag.Parse()

	// Validate required parameters
	if len(params.paths) == 0 {
		return nil, fmt.Errorf("at least one path parameter is required")
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

	// Scanning all directories
	for _, path := range params.paths {
		logger.Info().Msgf("Path to process: %s", path)
		err = scanner.Scan(path)
		if err != nil {
			logger.Fatal().Err(err).Msgf("Cannot scan directory: %s", path)
		}
	}

	// Results combining
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
