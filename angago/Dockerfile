FROM golang:1.18-alpine AS builder

LABEL app="angago"
LABEL maintainer="manigandan.jeff@gmail.com"
LABEL version="1.0.0"
LABEL description="A simple websocket reverse proxy server"

RUN mkdir -p /app

WORKDIR /app

COPY sqrx-angago sqrx-angago

FROM alpine:latest

WORKDIR /app
RUN apk update && apk add --no-cache ca-certificates

COPY --from=builder /app/sqrx-angago /app/sqrx-angago

# Load ENV
# Set Locales to support UTF-8 encoding to support non ASCII characters i.e. emojis
ENV LANG C.UTF-8
ENV LC_ALL C.UTF-8

EXPOSE 8081

RUN ls

CMD ["/app/sqrx-angago"]