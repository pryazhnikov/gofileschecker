package parameters

import (
	"flag"
	"fmt"
)

type RunParameters struct {
	Paths          []string // Paths to directories for scanning
	Debug          bool     // Enable debug logging
	FullFilePath   bool     // Show full file paths in output
	SkipEmptyFiles bool     // Do not process empty files
}

type runParametersParser struct {
	parsedParams *RunParameters
}

func NewRunParametersParser() *runParametersParser {
	parser := &runParametersParser{
		parsedParams: nil, // It should be initialized at Parse() method
	}

	return parser
}

func (p *runParametersParser) initFlagSet(name string) (*flag.FlagSet, *RunParameters) {
	parsedParams := &RunParameters{
		Paths: make([]string, 0),
	}
	flagSet := flag.NewFlagSet(name, flag.ContinueOnError)
	flagSet.BoolVar(&parsedParams.Debug, "debug", false, "Enable debug logging")
	flagSet.BoolVar(&parsedParams.FullFilePath, "fullpath", false, "Show full file paths in output")
	flagSet.BoolVar(&parsedParams.SkipEmptyFiles, "skipempty", false, "Skip empty files during scanning")
	flagSet.Func("path", "Path to directory for scanning (multiple usage is allowed)", func(flagValue string) error {
		parsedParams.Paths = append(parsedParams.Paths, flagValue)
		return nil
	})

	return flagSet, parsedParams
}

func (p *runParametersParser) Parse(args []string) (*RunParameters, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("no arguments provided")
	}

	flagSet, parsedParams := p.initFlagSet(args[0])
	if err := flagSet.Parse(args[1:]); err != nil {
		return nil, err
	}

	// Validate required parameters
	p.parsedParams = parsedParams
	if len(parsedParams.Paths) == 0 {
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
