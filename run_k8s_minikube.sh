#!/bin/sh

minikube addons enable ingress
minikube addons enable ingress-dns

kubectl create clusterrolebinding default-view --clusterrole=view --serviceaccount=default:default

# update local dns /etc/hosts
minikube ip 
echo "$(minikube ip) api.sqrx.com in.malwareriplabs.sqrx.com" | sudo tee -a /etc/hosts

# deploy sqrx-api

# after deployment, create a role binding for the service account
kubectl create clusterrolebinding ns:platform-u:default-r:cluster-admin --clusterrole=cluster-admin --serviceaccount=platform:default