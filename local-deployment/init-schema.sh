#! /bin/bash

LBT_PORT=8086

echo "initializing bigtable schema"
PROJECT="explorer"
INSTANCE="explorer"
HOST="127.0.0.1:$LBT_PORT"
cd ..
go run ./cmd/misc/main.go -config local-deployment/config.yml -command initBigtableSchema

echo "bigtable schema initialization completed"

echo "provisioning postgres db schema"
go run ./cmd/misc/main.go -config local-deployment/config.yml -command applyDbSchema
echo "postgres db schema initialization completed"