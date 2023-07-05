# squarex

Container Isolations- Websocket reverse proxy - Remote Browser Isolations

---

## Kubernetes deployment architecture

![squarex-k8s](/asset/sqrx-archi-k8s.jpg)

---

## Docker deployment architecture

![squarex](/asset/sqrx-archi-docker.jpg)

---

> Demo video

### How to run: via k8s minikube

```bash
# pre-requisite - install minikube
# curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
# sudo install minikube-linux-amd64 /usr/local/bin/minikube

# minikube start

./run_k8s_minikube.sh
```

### How to run: via docker containers

```bash
- Will create a docker network `sqrx-network`

./run.sh

```

---

### sqrx-api-server

- responsible for managing the containers and provisioning new sqrx-rbi containers.
- expose few endpoints

  - /try endpoint to spin up a new container
  - /status/{container_id} endpoint to check status of a container (running or not)
  - /stop/{termination_token} endpoint to stop a container

- Client, calls `POST http://api.sqrx.com/v1/try` -> sqrx-api-server now spins up a new container and returns the session_id and termination_token, also schedules background job to terminate the container after 10 mins.

```json
{
	"session": "in.malwareriplabs.sqrx.com/dfv-4c1c34a5-4f1c-47ae-a812-e414f0fc41c9/ws",
	"termination_token": "728bd120-aeb8-4b88-bcde-941e386d0e39",
	"created_at": "2023-07-02T23:20:28.957607349+05:30",
	"started_at": "2023-07-02T23:20:29.154287117+05:30",
	"valid_till": "2023-07-02T23:30:29.257156717+05:30"
}
```

### sqrx-angago

- it's a simple reverseproxy server, which will proxy the websocket connection to the specific container.
- when client connects to reverseproxy `ws://in.malwareriplabs.sqrx.com/dfv-4c1c34a5-4f1c-47ae-a812-e414f0fc41c9/ws`,
  it makes a downstream connection with the `sqrx-rbi` container and upgrade the client connection to websocket.
- upon successful connection, it will start streaming the message from the container to the client.

### sqrx-rbi

- simple websocket server, which will stream the message from the container to the client.
- just expose the container information to the client.
- doesn't implemented the remote browser isolation yet.

---

### Remote Browser Isolation

- Pixel pushing/streaming
- Page scrubbing

### limitations

- max 65536 socket connections per server
- each socket connection last for min 3 min to max 10 mins
  - in 1 hour, 6\*65536 = 3,93,216 connections
  - in 24 hours, (6*24)*65536 = 94,60,224 connections

```
# create network for sqrx
docker network create sqrx-network
```

## sqrx-api-server

- /v1/try endpoint to spin up a new container
- /v1/status endpoint to check status of a container (running or not)
- /v1/stop endpoint to stop a container
