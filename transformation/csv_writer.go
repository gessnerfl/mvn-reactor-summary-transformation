package transformation

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
)

type csvWriter struct{}

var csvHeaders = []string{"Release", "Build Iteration", "Module", "Duration (ms)"}

func (p *csvWriter) Validate(params *CliParameters) error {
	if params.OutputFile == nil {
		return errors.New("No target output file specified")
	}
	_, err := os.Stat(*params.OutputFile)
	if !os.IsNotExist(err) {
		return errors.New("Target output file already exists")
	}
	return nil
}

func (p *csvWriter) Write(params *CliParameters, data Report) error {
	file, err := os.Create(*params.OutputFile)
	if err != nil {
		return fmt.Errorf("Failed to create output file %s; %s", *params.OutputFile, err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write(csvHeaders)
	if err != nil {
		return fmt.Errorf("Failed to write header to CSV file; %s", err)
	}
	for rel, reactorSummaries := range data {
		for iter, moduleBuildTimes := range reactorSummaries {
			for _, moduleBuildTime := range moduleBuildTimes {
				durationInMs := moduleBuildTime.BuildTime.Milliseconds()
				line := []string{string(rel), strconv.Itoa(iter), moduleBuildTime.ModuleName, strconv.FormatInt(durationInMs, 10)}
				err = writer.Write(line)
				if err != nil {
					return fmt.Errorf("Failed to write header to CSV file; %s", err)
				}
			}
		}
	}
	return nil
}
