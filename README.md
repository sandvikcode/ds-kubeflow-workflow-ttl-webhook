# Mutating admission webhook 

This README will go through how to build and deploy the mutating admission webhook for argo. The goal with this project is to add the timetolive for the pods that are part of a argo workflow. The argo workflow is part of kubeflow deployment. 

```
k8s-kubeflow-mutate-webhook
│   Dockerfile
│   README.md  
│
└───pkg
│   │  
│   │
│   └───mutate
│       │   mutate_test.go
│       │   mutate.go
│ 
└───cmd
│    │   main.go
│    │   main_test.go
│ 
└───ssl
│    │  ssl_setup.sh
│ 
└───deployment
    │   csr.yaml
    │   deployment.yaml
    │   mutatingwebhookconfiguration.yaml
    │   secret.yaml
```

To deploy to k8s:

0) Get the CA bundle from the k8s API server. 
```bash
kubectl config view --raw --minify --flatten -o jsonpath='{.clusters[].cluster.certificate-authority-data}'
```
Update the caBundle value in the mutatingwebhookconfiguration.yaml file with the output. 

1)  Generate the needed SSL/TLS certificates using the bash script 
```bash 
sh ssl_setup.sh YOUR_APP_NAME YOUR_NAMESPACE
```

The bash script has the following requirements: 
- openssl is installed. 
- kubectl is available and the default namespace is set to YOUR_NAMESPACE. 

2) Add the ssl certificates as secrets to the k8s cluster in YOUR_NAMESPACE. 


``` bash 
export NAMESPACE=YOUR_NAMESPACE
export APP=YOUR_APP_NAME
kubectl create ns ${NAMESPACE}
kubectl create secret -n ${NAMESPACE} tls tls-secret --cert=${APP}.pem --key=${APP}.key
echo "Don't foreget to store the secrets in a safe place and NOT add them to git. " 
```

3) Build the Docker image. This can be build locally but I suggest building it using ``` gcloud builds submit ```. From the root of the project use the following command: 

```bash 
export PROJECT_ID=YOUR_PROJECT_ID
gcloud builds submit --tag gcr.io/${PROJECT_ID}/mutatingadmissiongwebhook .
```

4) Add the service, deployment, csr and mutatingwebhookconfiguration to k8s. 

```bash
kubectl apply -f deployment/service.yaml
kubectl apply -f deployment/deployment.yaml
kubectl apply -f deployment/mutatingwebhookconfiguration.yaml
```

5) Add the neede label to the namespace of interest in you k8s cluster. 

```bash 
kubectl label ns kubeflow mutateme=enabled
```

6) As default kubeflow currently havent given the access for the role "argo" to delete workflows. This is needed in order for argo to take action on the time to live. Thus this has to be updated: 

```bash 
kubectl edit clusterrole/argo
```
 
 If you have problems with that the workflows dont get deleted start the debugging with the following: 
 
 ```bash
 kubectl logs -n kubeflow deploy/workflow-controller
 ```

 ```bash 
 kubectl -n kubeflow get sa
 ```