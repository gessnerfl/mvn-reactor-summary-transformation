package transformation

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type sourceFolderReader struct{}

func (p *sourceFolderReader) Validate(params *CliParameters) error {
	if params.SourceFolder == nil {
		return errors.New("Source folder is missing")
	}
	info, err := os.Stat(*params.SourceFolder)
	if os.IsNotExist(err) {
		return errors.New("Source folder does not exist")
	}
	if !info.IsDir() {
		return errors.New("Source folder is not a directory")
	}
	return nil
}

func (p *sourceFolderReader) Read(params *CliParameters) (Report, error) {
	result := make(map[Release][]ReactorSummary)

	err := filepath.Walk(*params.SourceFolder, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			rel, summary, err := p.processFile(path)
			if err != nil {
				return err
			}
			result[rel] = append(result[rel], summary)
		}
		return nil
	})

	return result, err
}

func (p *sourceFolderReader) processFile(path string) (Release, ReactorSummary, error) {
	file, err := os.Open(path)
	defer file.Close()

	release := unknownRelease
	summary := ReactorSummary{}

	if err != nil {
		return release, summary, err
	}

	// Start reading from the file with a reader.
	reader := bufio.NewReader(file)

	regexVersion := regexp.MustCompile("(\\[INFO\\] Reactor Summary for )(\\S*)( )([\\d\\.]+(-SNAPSHOT))(:)")           //group 4 = release number
	regexModuleTime := regexp.MustCompile("(\\[INFO\\] )(\\S+)( \\.+ SUCCESS \\[\\s*)(\\d+[\\.\\:]\\d+)( )(\\w+)(\\])") // group 2 = module, group 4 = duration, group 6 = scale
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return release, summary, fmt.Errorf("Failed to process file %s; %s", path, err.Error())
		}

		if regexVersion.MatchString(line) {
			versionString := regexVersion.FindStringSubmatch(line)[4]
			release = Release(versionString)
		}

		if regexModuleTime.MatchString(line) {
			matches := regexModuleTime.FindStringSubmatch(line)
			module := matches[2]
			durationString := p.normalizeDuration(matches[4], matches[6])
			duration, err := time.ParseDuration(durationString)
			if err != nil {
				return release, summary, fmt.Errorf("Failed to get duration from %s; %s", line, err.Error())
			}
			summary = append(summary, ModuleBuildTime{
				ModuleName: module,
				BuildTime:  duration,
			})
		}
	}
	return release, summary, nil
}

func (p *sourceFolderReader) normalizeDuration(time string, scale string) string {
	if scale == "min" {
		return strings.ReplaceAll(time, ":", "m") + "s"
	}
	return time + scale
}
