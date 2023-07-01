FROM golang:1.18-alpine AS builder

LABEL app="rbi"
LABEL maintainer="manigandan.jeff@gmail.com"
LABEL version="1.0.0"
LABEL description="Remote Browser Isolation - A simple container isolate your browser from the internet"

RUN mkdir -p /app

WORKDIR /app

COPY sqrx-rbi sqrx-rbi
# COPY . /app
# RUN CGO_ENABLED=0 GOOS=linux go build -o sqrx-rbi -ldflags "-w -s"

FROM alpine:latest

WORKDIR /app
RUN apk update && apk add --no-cache ca-certificates

COPY --from=builder /app/sqrx-rbi /app/sqrx-rbi

# Load ENV
# Set Locales to support UTF-8 encoding to support non ASCII characters i.e. emojis
ENV LANG C.UTF-8
ENV LC_ALL C.UTF-8
# ENV ENV prod
# ENV PORT 80
# ENV RBI_CONFIG_PATH "/root/.rbi/config.yaml"

EXPOSE 8080

RUN ls

CMD ["/app/sqrx-rbi"]