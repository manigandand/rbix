FROM golang:1.18-alpine AS builder

LABEL app="angago"
LABEL maintainer="manigandan.jeff@gmail.com"
LABEL version="1.0.0"
LABEL description="A simple websocket reverse proxy server"

RUN mkdir -p /app

WORKDIR /app

COPY rbix-angago rbix-angago

FROM alpine:latest

WORKDIR /app
RUN apk update && apk add --no-cache ca-certificates

COPY --from=builder /app/rbix-angago /app/rbix-angago

# Load ENV
# Set Locales to support UTF-8 encoding to support non ASCII characters i.e. emojis
ENV LANG C.UTF-8
ENV LC_ALL C.UTF-8

EXPOSE 8081

CMD /app/rbix-angago /config/config.angago.json