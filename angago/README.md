# angago

---

angago(anga po/அங்க போ) means Go there. Localhost Proxy Tunnel(reverse proxy).

```
CGO_ENABLED=0 GOOS=linux go build -o sqrx-angago -ldflags "-w -s"

docker build -t manigandanjeff/sqrx-angago:latest .

docker run -d -p 8081:8081 \
    --network sqrx-network \
    --name box-sqrx-angago-1 \
    --hostname box-sqrx-angago \
    manigandanjeff/sqrx-angago:latest

```
