#! /bin/bash

echo Running Tests
go test ./...
echo Tests Finished
echo Building
go build github.com/amazing0x41/tsm-tool/cmd/tsm-tool
echo Finished building