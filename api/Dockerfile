FROM golang:1.18-alpine AS builder

LABEL app="api"
LABEL maintainer="manigandan.jeff@gmail.com"
LABEL version="1.0.0"
LABEL description="RbiX API server - Orchestrator & Management API's"

RUN mkdir -p /app

WORKDIR /app

COPY rbix-api rbix-api

FROM alpine:latest

WORKDIR /app
RUN apk update && apk add --no-cache ca-certificates

COPY --from=builder /app/rbix-api /app/rbix-api

# Load ENV
# Set Locales to support UTF-8 encoding to support non ASCII characters i.e. emojis
ENV LANG C.UTF-8
ENV LC_ALL C.UTF-8

EXPOSE 8080

CMD /app/rbix-api /config/config.api.json