#!/bin/sh

# build binary and docker images
# rbix-api
echo "======= Building rbix-api =======>"
cd api
CGO_ENABLED=0 GOOS=linux go build -o rbix-api -ldflags "-w -s"
docker build -t manigandanjeff/rbix-api:latest .
docker push manigandanjeff/rbix-api:latest

# rbix-angago
echo "======= Building rbix-angago =======>"
cd ../angago

CGO_ENABLED=0 GOOS=linux go build -o rbix-angago -ldflags "-w -s"
docker build -t manigandanjeff/rbix-angago:latest .
docker push manigandanjeff/rbix-angago:latest

# rbix-rbi
echo "======= Building rbix-rbi =======>"
cd ../rbi

CGO_ENABLED=0 GOOS=linux go build -o rbix-rbi -ldflags "-w -s"
docker build -t manigandanjeff/rbix-rbi:latest .
docker push manigandanjeff/rbix-rbi:latest