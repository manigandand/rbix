# Rbix - API Server

API Server - responsible for managing the containers and provisioning new rbix-rbi containers.

```bash
# build the docker image
CGO_ENABLED=0 GOOS=linux go build -o rbix-api -ldflags "-w -s"

docker build -t manigandanjeff/rbix-api:latest .
docker push manigandanjeff/rbix-api:latest
# ---------------------------------------------

docker run -d -p 8080:8080 \
    --network rbix-network \
    --name rbix-api-1 \
    --hostname rbix-api \
    manigandanjeff/rbix-api:latest
```

- POST /try -> endpoint to spin up a new container

```json
{
	"session": "localhost:8081/dfv-4c1c34a5-4f1c-47ae-a812-e414f0fc41c9/ws",
	"termination_token": "728bd120-aeb8-4b88-bcde-941e386d0e39",
	"created_at": "2023-07-02T23:20:28.957607349+05:30",
	"started_at": "2023-07-02T23:20:29.154287117+05:30",
	"valid_till": "2023-07-02T23:30:29.257156717+05:30"
}
```

- /status/{container_id} -> endpoint to check status of a container (running or not)

```json
{
	"error": "container not found",
	"status": 404
}
```

- /stop/{container_id} -> endpoint to stop a container
