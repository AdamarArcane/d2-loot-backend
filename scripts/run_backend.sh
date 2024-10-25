#!/bin/bash

cd cmd/generate_constants

go run main.go

cd ../..

go run .