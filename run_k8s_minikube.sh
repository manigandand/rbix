#!/bin/sh

# expect minikube to be installed and running(started with `minikube start`) 

minikube addons enable ingress
minikube addons enable ingress-dns

kubectl create clusterrolebinding default-view --clusterrole=view --serviceaccount=default:default

# create a namespace `platform`
kubectl apply -f k8s/namespace/sqrx-ns-platform.yaml

# mount the configmaps
kubectl apply -f k8s/configmap/config.yaml

# update local dns /etc/hosts
minikube ip 
echo "$(minikube ip) api.sqrx.com in.malwareriplabs.sqrx.com" | sudo tee -a /etc/hosts

# deploy sqrx-api
kubectl apply -f k8s/sqrx-api/deployment.yaml
kubectl apply -f k8s/sqrx-api/service.yaml

# deploy sqrx-angago
kubectl apply -f k8s/sqrx-angago/deployment.yaml
kubectl apply -f k8s/sqrx-angago/service.yaml

# after deployment, create a role binding for the service account
kubectl create clusterrolebinding ns:platform-u:default-r:cluster-admin --clusterrole=cluster-admin --serviceaccount=platform:default