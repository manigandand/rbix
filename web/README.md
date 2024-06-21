## Malware Samathi RIP

Web client (vnc client) to vncserver stream containers through websockets.

![Malware RIP](/static/img/app.rbix.com.png)

---

```bash
# build the docker image
docker build -t manigandanjeff/rbix-app:latest .
docker push manigandanjeff/rbix-app:latest
# ---------------------------------------------

docker run -d -p 3000:3000 \
    --network rbix-network \
    --name rbix-app-1 \
    --hostname rbix-app \
    manigandanjeff/rbix-app:latest
```
