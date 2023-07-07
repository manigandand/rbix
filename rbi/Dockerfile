FROM golang:1.18-alpine AS builder

LABEL app="rbi"
LABEL maintainer="manigandan.jeff@gmail.com"
LABEL version="1.0.0"
LABEL description="Remote Browser Isolation - A simple container isolate your browser from the internet"

RUN mkdir -p /app

WORKDIR /app

COPY sqrx-rbi sqrx-rbi

FROM ubuntu:latest

ENV USER=root
ENV DEBIAN_FRONTEND=noninteractive
ENV DEBCONF_NONINTERACTIVE_SEEN=true

WORKDIR /app
RUN apt-get update
RUN apt-get install -y xfce4 xfce4-goodies abiword gnupg apt-transport-https software-properties-common
RUN apt-get install -y ratpoison novnc websockify libxv1 libglu1-mesa xauth x11-utils xorg
RUN apt-get install -y tightvncserver

# Set VNC password
RUN mkdir -p ~/.vnc
RUN echo "password" | vncpasswd -f > ~/.vnc/passwd
RUN chmod 600 ~/.vnc/passwd

COPY --from=builder /app/sqrx-rbi /app/sqrx-rbi

# Load ENV
# Set Locales to support UTF-8 encoding to support non ASCII characters i.e. emojis
ENV LANG C.UTF-8
ENV LC_ALL C.UTF-8

EXPOSE 8888

# Start VNC server
CMD tightvncserver :1 -rfbauth ~/.vnc/passwd -auth ~/.vnc/passwd -geometry 1280x800 -depth 24 && \
    /app/sqrx-rbi && tail -f /dev/null