# angago

---

angago(anga po/அங்க போ) means Go there. Localhost Proxy Tunnel(reverse proxy).

```
CGO_ENABLED=0 GOOS=linux go build -o rbix-angago -ldflags "-w -s"

docker build -t manigandanjeff/rbix-angago:latest .

docker run -d -p 8081:8081 \
    --network rbix-network \
    --name box-rbix-angago-1 \
    --hostname box-rbix-angago \
    manigandanjeff/rbix-angago:latest

```
