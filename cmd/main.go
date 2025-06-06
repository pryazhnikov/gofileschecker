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

type runParametersParser struct {
	parsedParams *runParameters
}

func newRunParametersParser() *runParametersParser {
	parser := &runParametersParser{
		parsedParams: nil, // It should be initialized at Parse() method
	}

	return parser
}

func (p *runParametersParser) initFlagSet(name string) (*flag.FlagSet, *runParameters) {
	parsedParams := &runParameters{
		paths: make([]string, 0),
	}
	flagSet := flag.NewFlagSet(name, flag.ContinueOnError)
	flagSet.BoolVar(&parsedParams.debug, "debug", false, "Enable debug logging")
	flagSet.BoolVar(&parsedParams.fullFilePath, "fullpath", false, "Show full file paths in output")
	flagSet.BoolVar(&parsedParams.skipEmptyFiles, "skipempty", false, "Skip empty files during scanning")
	flagSet.Func("path", "Path to directory for scanning (multiple usage is allowed)", func(flagValue string) error {
		parsedParams.paths = append(parsedParams.paths, flagValue)
		return nil
	})

	return flagSet, parsedParams
}

func (p *runParametersParser) Parse(args []string) (*runParameters, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("no arguments provided")
	}

	flagSet, parsedParams := p.initFlagSet(args[0])
	if err := flagSet.Parse(args[1:]); err != nil {
		return nil, err
	}

	// Validate required parameters
	p.parsedParams = parsedParams
	if len(parsedParams.paths) == 0 {
		return nil, fmt.Errorf("at least one path parameter is required")
	}

	return parsedParams, nil
}

func (p *runParametersParser) IsParsed() bool {
	return p.parsedParams != nil
}

func (p *runParametersParser) Usage() {
	flagSet, _ := p.initFlagSet("gofilechecker")
	flagSet.Usage()
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
	paramsParser := newRunParametersParser()
	params, err := paramsParser.Parse(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
		paramsParser.Usage()
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
			if !params.fullFilePath {
				fileView = strings.TrimPrefix(file, pathPrefix)
			}

			fmt.Printf("- %s\n", fileView)
		}

		fmt.Println()
	}

	logger.Info().Msg("Done")
}
