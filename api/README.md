# SquareX - API Server

API Server - responsible for managing the containers and provisioning new sqrx-rbi containers.

- /try -> endpoint to spin up a new container
- /status/{container_id} -> endpoint to check status of a container (running or not)
- /stop/{container_id} -> endpoint to stop a container
