### RBI - Remote Browser Isolation

- Pixel pushing/streaming
- Page scrubbing

```
CGO_ENABLED=0 GOOS=linux go build -o sqrx-rbi -ldflags "-w -s"

docker build -t manigandanjeff/sqrx-rbi:latest .

docker run -d --rm \
    --network sqrx-network \
    --name box-sqrx-rbi-1 \
    --hostname box-sqrx-rbi-1 \
    -e CONTAINER_ID=box-sqrx-rbi-1 \
    manigandanjeff/sqrx-rbi:latest

```

# TODO: build scripts
