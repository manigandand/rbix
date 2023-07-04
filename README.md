# squarex

Container Isolations- Websocket reverse proxy - Remote Browser Isolations

---

![squarex](https://private-user-images.githubusercontent.com/9547223/250960270-0cc3ad6f-e429-49fe-8bc8-38292f4296cb.png?jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJrZXkiOiJrZXkxIiwiZXhwIjoxNjg4NDkwNDczLCJuYmYiOjE2ODg0OTAxNzMsInBhdGgiOiIvOTU0NzIyMy8yNTA5NjAyNzAtMGNjM2FkNmYtZTQyOS00OWZlLThiYzgtMzgyOTJmNDI5NmNiLnBuZz9YLUFtei1BbGdvcml0aG09QVdTNC1ITUFDLVNIQTI1NiZYLUFtei1DcmVkZW50aWFsPUFLSUFJV05KWUFYNENTVkVINTNBJTJGMjAyMzA3MDQlMkZ1cy1lYXN0LTElMkZzMyUyRmF3czRfcmVxdWVzdCZYLUFtei1EYXRlPTIwMjMwNzA0VDE3MDI1M1omWC1BbXotRXhwaXJlcz0zMDAmWC1BbXotU2lnbmF0dXJlPTMyMjc4YjRmZTE2ZjFkNTE5MDRjYzA5MDRlNTEyOGMxYzI2MjcyZTdmNGI0M2I4OGE2ZDNmMjFjNzZhZWQ1NzkmWC1BbXotU2lnbmVkSGVhZGVycz1ob3N0JmFjdG9yX2lkPTAma2V5X2lkPTAmcmVwb19pZD0wIn0.HcghCskK6B-MkIDJF0r2U9tm64BWxm5y8rYJR_VJ9Kk)

---

> Demo video

### How to run

```bash
./run.sh
```

---

- Will create a docker network `sqrx-network`

### sqrx-api-server

- responsible for managing the containers and provisioning new sqrx-rbi containers.
- expose few endpoints

  - /try endpoint to spin up a new container
  - /status/{container_id} endpoint to check status of a container (running or not)
  - /stop/{termination_token} endpoint to stop a container

- Client, calls `POST /try` -> sqrx-api-server now spins up a new container and returns the session_id and termination_token, also schedules background job to terminate the container after 10 mins.

```json
{
	"session": "localhost:8081/dfv-4c1c34a5-4f1c-47ae-a812-e414f0fc41c9/ws",
	"termination_token": "728bd120-aeb8-4b88-bcde-941e386d0e39",
	"created_at": "2023-07-02T23:20:28.957607349+05:30",
	"started_at": "2023-07-02T23:20:29.154287117+05:30",
	"valid_till": "2023-07-02T23:30:29.257156717+05:30"
}
```

### sqrx-angago

- it's a simple reverseproxy server, which will proxy the websocket connection to the specific container.
- when client connects to reverseproxy `ws://localhost:8081/dfv-4c1c34a5-4f1c-47ae-a812-e414f0fc41c9/ws`,
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

- /try endpoint to spin up a new container
- /status endpoint to check status of a container (running or not)
- /stop endpoint to stop a container
