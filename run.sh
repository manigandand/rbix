#!/bin/sh

# build binary and docker images
# sqrx-api
echo "======= Building sqrx-api =======>"
cd api
CGO_ENABLED=0 GOOS=linux go build -o sqrx-api -ldflags "-w -s"

# sqrx-angago
echo "======= Building sqrx-angago =======>"
cd ../angago

CGO_ENABLED=0 GOOS=linux go build -o sqrx-angago -ldflags "-w -s"

docker build -t manigandanjeff/sqrx-angago:latest .

# sqrx-rbi
echo "======= Building sqrx-rbi =======>"
cd ../rbi

CGO_ENABLED=0 GOOS=linux go build -o sqrx-rbi -ldflags "-w -s"

docker build -t manigandanjeff/sqrx-rbi:latest .

# create a docker network
echo "======= Creating docker network =======>"
docker network create sqrx-network

cd ..

# run sqrx-angago
echo "======= Running sqrx-angago docker container =======>"
docker run --rm -d -p 8081:8081 \
    -v /angago/config.angago.json:/mnt/config/config.angago.json \
    --network sqrx-network \
    --name box-sqrx-angago-1 \
    --hostname box-sqrx-angago \
    manigandanjeff/sqrx-angago:latest

# run sqrx-api
echo "======= Running sqrx-api docker container =======>"
cd api
./sqrx-api config.api.json