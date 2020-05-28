# KOBE operator

The KOBE operator acts as the orchestrator of the different components 
needed for a KOBE experiment to be deployed in a cluster of machines.

The KOBE operator implements the custom logic needed to react in those
KOBE-specific resources changes and maintain the necessary services in Kubernetes.

Once installed, the KOBE operator provides the following features:

- Monitors for new benchmark and experiment submissions to the
  Kubernetes cluster and triggers the deployment of the appropriate
  pods.
- Handles the initialization of the data sources and federators
  including the data import and metadata generation.
- Configures the network connections between the pods to realize 
  simulated network delays.

## API

To learn more about the CRDs have a look at the [API doc](docs/api.md).


## Installation of the operator

---
**NOTE**

This is a guide for installing the KOBE operator *only*
in a Kubernetes cluster. 
If you are looking for the installation of the entire KOBE system
please consult the top-level [README](../README.md).

---

## Prerequisites

- `Kubernetes` >= 1.8.0
- `nfs-commons` installed in the nodes of the cluster. If in debian or
   Ubuntu you can install it using `apt-get install nfs-common`

## Deployment

To quickly install the KOBE operator in a Kubernetes cluster, run the
following commands:
```
kubectl apply -f deploy/crds
kubectl apply -f deploy/service_account.yaml
kubectl apply -f deploy/clusterrole.yaml
kubectl apply -f deploy/clusterrole_binding.yaml
kubectl apply -f deploy/role.yaml
kubectl apply -f deploy/operator.yaml
```

## Removal of the operator

```
kubectl delete -f deploy/operator.yaml
kubectl delete -f deploy/role.yaml
kubectl delete -f deploy/clusterrole_binding.yaml
kubectl delete -f deploy/clusterrole.yaml
kubectl delete -f deploy/service_account.yaml
kubectl delete -f deploy/crds
```
