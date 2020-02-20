package main

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/gessnerfl/mvn-reactor-summary-transformation/transformation"
)

func main() {
	sourceFolder := flag.String("src", "", "The source folder where the build summaries are stored")
	outputFormatString := flag.String("fmt", "", "The output format [csv,es]")
	outFile := flag.String("out", "", "CSV: The target file path")
	esAddress := flag.String("es-address", "", "ES: The address of the elastic search host")
	esUsername := flag.String("username", "", "ES: the username to access elasticsearch - optional")
	esPassword := flag.String("password", "", "ES: the password of the user to access elasticsearch - optional")

	flag.Parse()

	outputFormat, err := parseOutputFormat(outputFormatString)

	if err != nil {
		fmt.Println(err)
		flag.Usage()
	}

	params := transformation.CliParameters{
		SourceFolder: sourceFolder,
		OutputFormat: outputFormat,
		OutputFile:   outFile,
		EsAddress:    esAddress,
		Username:     esUsername,
		Password:     esPassword,
	}

	err = transformation.Apply(&params)

	if err != nil {
		fmt.Println(err)
	}
}

func parseOutputFormat(str *string) (transformation.OutputFormat, error) {
	if str != nil {
		if strings.ToLower(*str) == "csv" {
			return transformation.CSVFormat, nil
		}
		if strings.ToLower(*str) == "es" {
			return transformation.ElasticsearchFormat, nil
		}
	}
	return transformation.NoneFormat, errors.New("Output format missing or not valid. Supported formats are [csv]")
}
