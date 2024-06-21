#!/bin/sh

# expect minikube to be installed and running(started with `minikube start`) 

minikube addons enable ingress
minikube addons enable ingress-dns

kubectl create clusterrolebinding default-view --clusterrole=view --serviceaccount=default:default

# create a namespace `platform`
kubectl apply -f k8s/namespace/rbix-ns-platform.yaml

# mount the configmaps
kubectl apply -f k8s/configmap/config.yaml

# update local dns /etc/hosts
minikube ip
echo "$(minikube ip) app.rbixlabs.com api.rbixlabs.com in.malwaresamathi.rbixlabs.com" | sudo tee -a /etc/hosts

# deploy rbix-app
kubectl apply -f k8s/rbix-app/deployment.yaml
kubectl apply -f k8s/rbix-app/service.yaml

# deploy rbix-api
kubectl apply -f k8s/rbix-api/deployment.yaml
kubectl apply -f k8s/rbix-api/service.yaml

# deploy rbix-angago
kubectl apply -f k8s/rbix-angago/deployment.yaml
kubectl apply -f k8s/rbix-angago/service.yaml

# after deployment, create a role binding for the service account
kubectl create clusterrolebinding ns:platform-u:default-r:cluster-admin --clusterrole=cluster-admin --serviceaccount=platform:default