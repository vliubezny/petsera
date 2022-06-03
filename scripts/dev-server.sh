#!/bin/bash

set -e

PID_LIST=""

trap 'kill $PID_LIST' SIGINT

echo "Start postgres"
docker run -i --rm --name petseradb -p "127.0.0.1:5432:5432" \
    -e POSTGRES_USER=petsera -e POSTGRES_PASSWORD=root \
    -e POSTGRES_DB=petsera postgis/postgis:13-master & pid=$!
PID_LIST+=" $pid"

echo "Start backend"
air & pid=$!
PID_LIST+=" $pid"

# echo "Start frontend"
# cd ui && npm run dev & pid=$!
# PID_LIST+=" $pid"
# cd ..

wait $PID_LIST

echo "Dev server is stopped"