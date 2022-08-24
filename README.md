# Kubernetes main abstractions

## Preparation

If you want to test application inside repository - do not forget to:   
* Span all preparation containers - kubectl apply -f kubernetes/preparation
* Build application - docker build -t $registry/$image_name:$tag $directory
* Push application - docker push -t (tag from previous step)
* Change image section in deployment.yaml file to new $registry/$image_name:$tag

To access docker registry:
* Register on https://hub.docker.com/
* Now your image tag will look like $your_login/$image_name:$tag   

To create your own registry:
* See https://docs.docker.com/registry/deploying/

To use minikube registry:
* See https://minikube.sigs.k8s.io/docs/handbook/registry/

To install image directly to minikube:
* minikube image build -t $image:$tag -f $dockerfile_path $dockerfile_directory
  * minikube image build -t books:minikube -f Dockerfile . (Example)

---

## Application

### Application with cluster IP and ingress
* kubectl apply -f kubernetes/examples/application-cluster-ip   
  * **Do not forget to turn Ingress addon on if you are using minikube (minikube addons enable ingress)**   
* Add IP of your cluster to /etc/hosts file (Windows C:\Windows\System32\drivers\etc\hosts):   
  * 192.168.49.2 ingress.local
* Check if it is working: curl -v ingress.local   

### Application with nodePort
* kubectl apply -f kubernetes/examples/application-nodeport   
* Check if it is working: curl $(minikube ip):30007 | curl localhost:30007
  * Be aware, on OS (Windows, Mac with Minikube) abstraction won't be reachable. Don't be afraid, it is working.   
  * Access NodePort on minikube (https://minikube.sigs.k8s.io/docs/handbook/accessing/)
---

## Folder structure
* kubernetes/application - application abstractions   
* kubernetes/configuration - configuration abstractions   
* kubernetes/network - network abstractions   
* kubernetes/preparation - abstractions for presentation   
* kubernetes/job - abstractions for job 
* kubernetes/examples - ready to deploy applications based on Nginx image

---

## Kubernetes network
* https://www.youtube.com/c/%D0%90%D1%80%D1%82%D1%83%D1%80%D0%9A%D1%80%D1%8E%D0%BA%D0%BE%D0%B2

---

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
* kubect config use-context <-context_name-> – use the context for kubectl

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

---

## Minikube commands
* minikube node add - add another node for minikube
* minikube addons list - all addons for minikube
* minikube addons enable <addon_name> -  enable any addon for minikube
* minikube image build -t <tag> <directory> - build image directly to minikube 
* minikube mount <from>:<to> - mount directory to minikube
* minikube ip - show ip address of minikube
* minikube —help - help description of minikube commands