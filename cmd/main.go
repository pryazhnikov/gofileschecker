package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/pryazhnikov/gofileschecker/internal/checkers"
	"github.com/pryazhnikov/gofileschecker/internal/parameters"
	"github.com/pryazhnikov/gofileschecker/internal/scanner"
	"github.com/rs/zerolog"
)

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
	paramsParser := parameters.NewRunParametersParser()
	params, err := paramsParser.Parse(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
		paramsParser.Usage()
		os.Exit(1)
	}

	logger := newLogger(params.Debug)
	logger.Info().Msgf("Logger level set: %s", logger.GetLevel().String())

	fileChecker := checkers.NewFileChecker(params.SkipEmptyFiles)
	scanner := scanner.NewDirectoryScanner(logger, fileChecker)

	// Scanning all directories
	for _, path := range params.Paths {
		logger.Info().Msgf("Path to process: %s", path)
		err = scanner.Scan(path)
		if err != nil {
			logger.Fatal().Err(err).Msgf("Cannot scan directory: %s", path)
		}
	}

	scanRes := scanner.Summary()
	logger.Info().Msgf(
		"Directories: %d, files: %d, errors: %d, skipped: %d",
		scanRes.Directories(),
		scanRes.Files(),
		scanRes.Errors(),
		scanRes.Skipped(),
	)

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
			if !params.FullFilePath {
				fileView = strings.TrimPrefix(file, pathPrefix)
			}

			fmt.Printf("- %s\n", fileView)
		}

		fmt.Println()
	}

	logger.Info().Msg("Done")
}
