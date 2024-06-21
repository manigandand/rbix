#!/bin/sh

# build binary and docker images
# rbix-api
echo "======= Building rbix-api =======>"
cd api
CGO_ENABLED=0 GOOS=linux go build -o rbix-api -ldflags "-w -s"

# rbix-angago
echo "======= Building rbix-angago =======>"
cd ../angago

CGO_ENABLED=0 GOOS=linux go build -o rbix-angago -ldflags "-w -s"

docker build -t manigandanjeff/rbix-angago:latest .

# rbix-rbi
echo "======= Building rbix-rbi =======>"
cd ../rbi

CGO_ENABLED=0 GOOS=linux go build -o rbix-rbi -ldflags "-w -s"

docker build -t manigandanjeff/rbix-rbi:latest .

# create a docker network
echo "======= Creating docker network =======>"
docker network create rbix-network

cd ..

# run rbix-angago
echo "======= Running rbix-angago docker container =======>"
docker run --rm -d -p 8081:8081 \
    -v /angago/config.angago.json:/mnt/config/config.angago.json \
    --network rbix-network \
    --name box-rbix-angago-1 \
    --hostname box-rbix-angago \
    manigandanjeff/rbix-angago:latest

# run rbix-api
echo "======= Running rbix-api docker container =======>"
cd api
./rbix-api config.api.json