# Ingress controller

## Information
Ingress - object in Kubernetes API, description how ingress should route traffic.   
To enable ingress we have to install Ingress controller.   
Ingress controllers - https://kubernetes.io/docs/concepts/services-networking/ingress-controllers/   
Deployment - https://kubernetes.github.io/ingress-nginx/deploy

## Ingress using NodePort
Open concrete port on each node (Configurable) 

## Ingress using HostNetwork
Using node network (Port is accessible on node network interface)