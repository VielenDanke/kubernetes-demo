# Instruction

## Apply deployment
`kubectl apply -f deployment.yaml`

## Apply service
`kubectl apply -f clusterip-service.yaml`

## Test your application
```yaml
ports:
  - port: 80
    targetPort: http
    nodePort: 30007
```
* Docker desktop:
  * Open terminal -> `curl -v localhost:$nodePort`
  * Go to browser and try `http://localhost:$nodePort`, you should see Nginx main page
* Minikube
  * Open terminal -> `curl -v $minikube_ip:$nodePort`
  * Go to browser and try `http://minikube_ip:$nodePort`, you should see Nginx main page
  * Minikube IP - tap in terminal `minikube ip`