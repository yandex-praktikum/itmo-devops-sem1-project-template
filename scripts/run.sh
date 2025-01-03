#!/bin/bash

goose -dir ./migrations up
./src/main/app