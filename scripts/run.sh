#!/bin/bash

docker compose up -d
goose -dir ./migrations up