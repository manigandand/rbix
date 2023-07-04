FROM golang:1.18-alpine AS builder

LABEL app="api"
LABEL maintainer="manigandan.jeff@gmail.com"
LABEL version="1.0.0"
LABEL description="Sqrx API server - Orchestrator & Management API for Sqrx"

RUN mkdir -p /app

WORKDIR /app

COPY sqrx-api sqrx-api

FROM alpine:latest

WORKDIR /app
RUN apk update && apk add --no-cache ca-certificates

COPY --from=builder /app/sqrx-api /app/sqrx-api

# Load ENV
# Set Locales to support UTF-8 encoding to support non ASCII characters i.e. emojis
ENV LANG C.UTF-8
ENV LC_ALL C.UTF-8

EXPOSE 8080

CMD /app/sqrx-api /config/config.api.json