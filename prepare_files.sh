#!/usr/bin/env bash

# Execute this script to prepare superwatcher misc files,
# including copying config example

DB_PATHS=( "./tmp" "./pkg/tmp" "./pkg/reorgsim/tmp" "./pkg/datagateway/tmp" )

for d in "${DB_PATHS[@]}"; do
    echo "preparing $d";

    mkdir -pv $d;
    touch "$d/fakeredis.db";
done;

[ ! -f "./superwatcher-demo/config/config.yaml" ]\
 && echo "copying superwatcher-demo config file"\
 && cp -v ./superwatcher-demo/config/config.yaml.example ./superwatcher-demo/config/config.yaml;
