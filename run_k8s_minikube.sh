#!/bin/sh

minikube addons enable ingress
# 
minikube ip 

echo "192.168.49.2 api.sqrx.com in.malwareriplabs.sqrx.com" | sudo tee -a /etc/hosts