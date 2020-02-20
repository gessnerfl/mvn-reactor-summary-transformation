package transformation

import (
	"fmt"
	"time"
)

//ModuleBuildTime type definition representing the build time of a Module
type ModuleBuildTime struct {
	ModuleName string
	BuildTime  time.Duration
}

//ReactorSummary type definition for the parsed Maven reactor summary of a build
type ReactorSummary []ModuleBuildTime

//Release type definition of a maven release number
type Release string

const unknownRelease = Release("unknown")

//Report type definition of a Report containing all Reactor summaries
type Report map[Release][]ReactorSummary

//CliParameters type for input parameters from cli
type CliParameters struct {
	SourceFolder *string
	OutputFormat OutputFormat
	OutputFile   *string
	EsAddress    *string
	Username     *string
	Password     *string
}

//OutputFormat type definition for output formats
type OutputFormat string

const (
	//NoneFormat dummy output format when no format was provided
	NoneFormat = OutputFormat("NONE")
	//CSVFormat constant value for CSV output format
	CSVFormat = OutputFormat("CSV")
	//ElasticsearchFormat constant value for ElasticSearch output format
	ElasticsearchFormat = OutputFormat("ES")
)

//Apply executs the transformation using the provided input parameters
func Apply(params *CliParameters) error {
	r := createReader()
	err := r.Validate(params)
	if err != nil {
		return err
	}

	w, err := createWriter(params)
	if err != nil {
		return err
	}
	err = w.Validate(params)
	if err != nil {
		return err
	}

	data, err := r.Read(params)
	if err != nil {
		return err
	}

	return w.Write(params, data)
}

func createWriter(params *CliParameters) (writer, error) {
	switch params.OutputFormat {
	case CSVFormat:
		return &csvWriter{}, nil
	case ElasticsearchFormat:
		return &es7Writer{}, nil
	}
	return nil, fmt.Errorf("Unsupported output format %s", params.OutputFormat)
}

type writer interface {
	Write(parmas *CliParameters, summaries Report) error
	Validate(params *CliParameters) error
}

func createReader() reader {
	return &sourceFolderReader{}
}

type reader interface {
	Read(params *CliParameters) (Report, error)
	Validate(params *CliParameters) error
}
