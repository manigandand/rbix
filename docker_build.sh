#!/bin/sh

# build binary and docker images
# sqrx-api
echo "======= Building sqrx-api =======>"
cd api
CGO_ENABLED=0 GOOS=linux go build -o sqrx-api -ldflags "-w -s"
docker build -t manigandanjeff/sqrx-api:latest .
docker push manigandanjeff/sqrx-api:latest

# sqrx-angago
echo "======= Building sqrx-angago =======>"
cd ../angago

CGO_ENABLED=0 GOOS=linux go build -o sqrx-angago -ldflags "-w -s"
docker build -t manigandanjeff/sqrx-angago:latest .
docker push manigandanjeff/sqrx-angago:latest

# sqrx-rbi
echo "======= Building sqrx-rbi =======>"
cd ../rbi

CGO_ENABLED=0 GOOS=linux go build -o sqrx-rbi -ldflags "-w -s"
docker build -t manigandanjeff/sqrx-rbi:latest .
docker push manigandanjeff/sqrx-rbi:latest