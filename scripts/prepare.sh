#!/bin/bash

go install github.com/pressly/goose/v3/cmd/goose@latest
go build -C ./src/main -o app