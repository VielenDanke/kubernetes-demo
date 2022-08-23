# Kubernetes main abstractions

## Folder structure
* kubernetes/application - application abstractions   
* kubernetes/configuration - configuration abstractions   
* kubernetes/network - network abstractions   
* kubernetes/preparation - abstractions for presentation   
* kubernetes/job - abstractions for job   

## Kubernetes network
* https://www.youtube.com/c/%D0%90%D1%80%D1%82%D1%83%D1%80%D0%9A%D1%80%D1%8E%D0%BA%D0%BE%D0%B2

## Kubernetes commands
### Main commands:
Short names:   
* configmap – cm   
* pod – po   
* services – svc   
* ingress – ing

### Work with context:
All contexts are placed in file: ~/.kube/config
* kubectl config get-contexts – find all available contexts
* kubect config use-context <context_name> – use the context for kubectl

### To edit resource:

* kubectl edit <resource> -n <namespace> <resource_name>
  * Example: kubectl edit secrets -n trends plutus-tls

### To describe resource:
* kubectl describe <resource> -n <namespace> <resource_name>
  * Example: kubectl describe pod -n trends plutus-integ-asdasd-123we

### To download resource:
* kubectl get <resource> -n <namespace> <resource_name> -o <format> > some.<format> – download resource as file with specific format (json, yaml). 
  * Example: kubectl get cm -n trends coolmap -o yaml > some.yaml

### To get resource:
* kubectl get <resource> -n <namespace> – all resources in the namespace
* kubectl get <resource> --all-namespaces – all resource in all namespaces

### Pod:
* kubectl exec -n <namespace> -it <pod_name> - - /bin/sh – working inside a container
* kubectl logs -n <namespace> <pod_name> – get pod logs