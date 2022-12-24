#!/usr/bin/env bash

# Execute this script to prepare superwatcher misc files,
# including copying config example

mkdir -pv ./tmp ./pkg/reorgsim/tmp;
touch ./tmp/fakeredis.db ./pkg/reorgsim/tmp/fakeredis.db;
cp -v ./superwatcher-demo/config/config.yaml.example ./superwatcher-demo/config/config.yaml;
