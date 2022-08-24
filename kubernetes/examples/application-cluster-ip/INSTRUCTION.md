# Instruction

## Apply deployment
`kubectl apply -f deployment.yaml`

## Apply service
`kubectl apply -f clusterip-service.yaml`

## Apply ingress
`kubectl apply -f ingress.yaml `  
**Do not forget to install ingress addon**

## Add ingress host to /etc/hosts file
```yaml
spec:
  ingressClassName: nginx
  rules:
    - host: ingress.local
```
You host file should look like:
```
127.0.0.1      localhost
127.0.0.1      ingress.local
```
In case you are using minikube:
```
$minikube_ip   ingress.local
```

## Test your application
* Open terminal -> `curl -v ingress.local`   
* Go to browser and try `http://ingress.local`, you should see Nginx main page