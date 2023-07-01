### RBI - Remote Browser Isolation

- Pixel pushing/streaming
- Page scrubbing

```
CGO_ENABLED=0 GOOS=linux go build -o sqrx-rbi -ldflags "-w -s"

docker build -t sqrx/rbi:latest .

docker run -d --network sqrx-network --name box-sqrx-rbi-1 --hostname box-sqrx-rbi-1 sqrx/rbi:latest

```

# TODO: build scripts
