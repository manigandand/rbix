## Malware RIP

Web client (vnc client) to vncserver stream containers through websockets.

![Malware RIP](/static/img/app.sqrx.com.png)

---

```bash
# build the docker image
docker build -t manigandanjeff/sqrx-app:latest .
docker push manigandanjeff/sqrx-app:latest
# ---------------------------------------------

docker run -d -p 3000:3000 \
    --network sqrx-network \
    --name sqrx-app-1 \
    --hostname sqrx-app \
    manigandanjeff/sqrx-app:latest
```
