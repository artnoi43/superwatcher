#!/usr/bin/env bash

# Execute this script to prepare superwatcher misc files,
# including copying config example

DB_PATHS=( "./tmp" "./pkg/tmp" "./pkg/reorgsim/tmp" "./pkg/datagateway/tmp" )

for d in "${DB_PATHS[@]}"; do
    echo "preparing $d";

    mkdir -pv $d;
    touch "$d/fakeredis.db";
done;

SERVICE_CONF_PATH="./examples/demoservice/config"

[ ! -f "${SERVICE_CONF_PATH}/config.yaml" ]\
 && echo "copying config file in ${SERVICE_CONF_PATH}"\
 && cp -v "${SERVICE_CONF_PATH}/config.yaml.example" "${SERVICE_CONF_PATH}/config.yaml";
