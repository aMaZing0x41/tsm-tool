# tsm-tool
A tool for checking influx tsm files

## Using the tool
`./tsm-tool -file ./influx/000000000000001-000000001.tsm [-debug]`

## Building
1. Clone the repo into GOPATH
1. cd into project folder
1. `./build.sh`

## Future Condsiderations
- refactor out most TSM file handling and logic to package
- better error handling
- better logging
- more unit/integration tests with larger and more complex TSM files