FROM ubuntu:latest

ENV USER=root
ENV DEBIAN_FRONTEND=noninteractive
ENV DEBCONF_NONINTERACTIVE_SEEN=true

# Install necessary packages
RUN apt-get update
RUN apt-get install -y xfce4 xfce4-goodies abiword gnupg apt-transport-https wget software-properties-common
RUN apt-get install -y ratpoison novnc websockify libxv1 libglu1-mesa xauth x11-utils xorg
RUN apt-get install -y tightvncserver

# Set VNC password
RUN mkdir -p ~/.vnc
RUN echo "password" | vncpasswd -f > ~/.vnc/passwd
RUN chmod 600 ~/.vnc/passwd

# Expose VNC port
EXPOSE 80

# Start VNC server
CMD tightvncserver :1 -rfbauth ~/.vnc/passwd -auth ~/.vnc/passwd -geometry 1280x800 -depth 24 && \
    websockify -D --web=/usr/share/novnc/ --cert=~/novnc.pem 80 localhost:5901 && tail -f /dev/null
