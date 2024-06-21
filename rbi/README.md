### RBI - Remote Browser Isolation

- Pixel pushing/streaming
- Page scrubbing

```
CGO_ENABLED=0 GOOS=linux go build -o rbix-rbi -ldflags "-w -s"

docker build -t manigandanjeff/rbix-rbi:latest .

docker run -d --rm \
    --network rbix-network \
    --name box-rbix-rbi-1 \
    --hostname box-rbix-rbi-1 \
    -e CONTAINER_ID=box-rbix-rbi-1 \
    manigandanjeff/rbix-rbi:latest

```

# TODO: build scripts
