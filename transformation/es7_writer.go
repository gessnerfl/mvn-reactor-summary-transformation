package transformation

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/olivere/elastic/v7"
	"github.com/rs/xid"
)

type moduleBuildReport struct {
	Release        string `json:"release"`
	BuildIteration string `json:"build_iteration"`
	ModuleName     string `json:"module_name"`
	Duration       int64  `json:"duration"`
}

const mapping = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings":{
		"properties":{
			"release": { "type":"keyword" },
			"build_iteration": { "type":"keyword" },
			"module_name": { "type":"keyword" },
			"duration": { "type":"integer" }
		}
	}
}`

type es7Writer struct{}

func (p *es7Writer) Validate(params *CliParameters) error {
	if params.EsAddress == nil || strings.TrimSpace(*params.EsAddress) == "" {
		return errors.New("Elasticsearch address not provided")
	}
	return nil
}

func (p *es7Writer) Write(params *CliParameters, data Report) error {
	ctx := context.Background()

	options := []elastic.ClientOptionFunc{elastic.SetSniff(false), elastic.SetURL(*params.EsAddress)}
	if params.Username != nil && params.Password != nil {
		options = append(options, elastic.SetBasicAuth(*params.Username, *params.Password))
	}

	client, err := elastic.NewClient(options...)
	if err != nil {
		return fmt.Errorf("Failed to connect to elasticsearch; %s", err.Error())
	}

	indexName := "mvn-reactor-analysis-" + time.Now().Format("20060102_150405")
	_, err = client.CreateIndex(indexName).BodyString(mapping).Do(ctx)
	if err != nil {
		return fmt.Errorf("Failed to create Index, %s", err.Error())
	}

	for rel, reactorSummaries := range data {
		for iter, moduleBuildTimes := range reactorSummaries {
			for _, moduleBuildTime := range moduleBuildTimes {
				doc := p.createElasticSearchDocument(rel, iter, moduleBuildTime)
				id := xid.New()
				_, err = client.Index().Index(indexName).Id(id.String()).BodyJson(doc).Do(ctx)
				if err != nil {
					return fmt.Errorf("Failed to index build report document; %s", err.Error())
				}
			}

		}
	}
	return nil
}

func (p *es7Writer) createElasticSearchDocument(release Release, buildIteration int, moduleBuildTime ModuleBuildTime) *moduleBuildReport {
	return &moduleBuildReport{
		Release:        string(release),
		BuildIteration: fmt.Sprintf("%s-%d", string(release), buildIteration),
		ModuleName:     moduleBuildTime.ModuleName,
		Duration:       moduleBuildTime.BuildTime.Milliseconds(),
	}
}
