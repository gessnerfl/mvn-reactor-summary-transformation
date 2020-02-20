# Maven Ractor Summary Transformation to CSV or Elasticsearch

This is a simple to read the build summaries of Maven Multimodule builds and convert/transform them into
CSV or Elasticsearch. It is expected that the summary files belong to the same project.

This is useful to analyze build times within a project over time. You can store multiple build summaries of
the same release and the tool will just add all of them with an iteration counter.

## How to Use

### CSV

1. Clone this repository
2. Run the program `go run main.go -src=PATH_TO_BUILD_SUMMARIES -fmt=csv -out=PATH_TO_CSV_FILE`

The CSV file will contain the following data:
`"Release", "Build Iteration", "Module", "Duration (ms)"`

### ElasticSearch

1. Clone this repository
2. If you want to use a local elasticsearch running in a docker container you can just use the provided docker-compose.yml file
3. Run teh program `go run main.go -src=PATH_TO_BUILD_SUMMARIES -fmt=es -es-address=ELASTIC_SEARCH_URL -username=ELASTICSEARCH_USER -password=ELASTICSEARCH_PASSWORD`

Username and Password are optional and depend on your installation. 

When you use the provided docker setup (docker-compose.yml) you can run the tool with the command: `go run main.go -src=PATH_TO_BUILD_SUMMARIES -fmt=csv -es-address=http://localhost:9200`

The data in Elasticsearch contains:

- release = the release number
- build_iteration = release number concatenated with a counter to uniquely identify a build iteration
- module_name = the maven module name
- duration = the build duration in milliseconds

The data is store in an index with the name `mvn-reactor-analysis-YYYYMMDD_HHMMSS` whereas `YYYYMMDD` is the date of execution `HHMMSS` is the time of execution