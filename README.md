# KOBE: Cloud-Native Open Benchmark Engine for SPARQL Query Processors

KOBE is a benchmarking system that leverages
[Docker](https://docker.io) and [Kubernetes](https://kubernetes.io) in
order to reproduce experiments of federated query processing in a
collections of data sources.

## Overview

In the SPARQL query processing community, as well as in the wider
databases community, benchmark reproducibility is based on releasing
datasets and query workloads. However, this paradigm breaks down for
federated query processors, as these systems do not manage the data
they serve to their clients but provide a data-integration abstraction
over the actual query processors that are in direct contact with the
data.

The KOBE benchmarking engine is a system that aims to provide a
generic platform to perform benchmarking and experimentation that can
be reproducible in different environments. It was designed with the
following objectives in mind:

1. to allow for benchmark and experiment specifications to be
   reproduced in different environments and be able to produce
   comparable and reliable results;
2. to ease the deployment of complex benchmarking experiments by
   automating the tedious tasks of initialization and execution.

## Installation

### Prerequisites

- `Kubernetes` >= 1.8.0
- `nfs-commons` installed in the nodes of the cluster. If in debian or
   Ubuntu you can install it using `apt-get install nfs-common`

### Installation of the Kubernetes operator

KOBE needs the Kubernetes operator that needs to be installed in the
Kubernetes cluster. To quickly install the KOBE operator in a
Kubernetes cluster, run the following

```
kubectl apply -f operator/deploy/crds
kubectl apply -f operator/deploy/service_account.yaml
kubectl apply -f operator/deploy/clusterrole.yaml
kubectl apply -f operator/deploy/clusterrole_binding.yaml
kubectl apply -f operator/deploy/role.yaml
kubectl apply -f operator/deploy/operator.yaml
```

You will get a confirmation message that each resource has
successfully been created.
This will set the operator running in your Kubernetes cluster and
needs to be done only once.

### Installation of Networking subsystem

KOBE uses istio to support network delays between the different 
deployments. To install istio you can consult the official 
[installation guide](https://istio.io/docs/setup/getting-started/) 
or you type the following commands.

```
curl -L https://istio.io/downloadIstio | sh -
export PATH=`pwd`/istio-1.6.0/bin:$PATH
istioctl manifest apply --set profile=default
```

### Installation of the Evaluation Metrics Extraction subsystem

To enable the evaluation metrics extraction subsystem, run the following
```
helm repo add elastic https://helm.elastic.co
helm repo add kiwigrid https://kiwigrid.github.io
helm install elastic/elasticsearch --name elasticsearch --set persistence.enabled=false --set replicas=1 --version 7.6.2
helm install elastic/kibana --name kibana --set service.type=NodePort --version 7.6.2
helm install --name fluentd -f operator/deploy/efk-config/fluentd-values.yaml kiwigrid/fluentd-elasticsearch --version 3.0.1
kubectl apply -f operator/deploy/efk-config/kobe-kibana-configuration.yaml
```

These result in the simplest setup of an one-node
[ElasticSearch](https://github.com/elastic/helm-charts/blob/master/elasticsearch)
that does not persist data across pod recreation, a
[Fluentd](https://github.com/kiwigrid/helm-charts/tree/master/charts/fluentd-elasticsearch)
`DaemonSet` and a
[Kibana](https://github.com/elastic/helm-charts/tree/master/kibana)
node that exposes a `NodePort`. 

After all pods are in Running state Kibana dashboards can be accessed
at 
```
http://<NODE-IP>:<NODEPORT>/app/kibana#/dashboard/
``` 
where `<NODE-IP>` the IP of any of the Kubernetes cluster nodes and
`<NODEPORT>` the result of `kubectl get -o
jsonpath="{.spec.ports[0].nodePort}" services kibana-kibana`.

The setup can be customized by changing the configuration parameters
of each helm chart. Please check the corresponding documentation of
each chart for more info.

## Removal

```
kubectl delete -f operator/deploy/operator.yaml
kubectl delete -f operator/deploy/role.yaml
kubectl delete -f operator/deploy/clusterrole_binding.yaml
kubectl delete -f operator/deploy/clusterrole.yaml
kubectl delete -f operator/deploy/service_account.yaml
kubectl delete -f operator/deploy/crds
```

To remove the evaluation metrics extraction subsystem run
```
helm delete --purge elasticsearch
helm delete --purge kibana
helm delete --purge fluentd
helm repo remove elastic
helm repo remove kiwigrid
```
and then in each Kubernetes node
```
rm -rf /var/log/fluentd-buffers/
rm /var/log/containers.log.pos
```

## Example

The typical workflow of defining a KOBE experiment is the following.
1. Create a set of datasets by defining new [Datasets](docs/api.md#kobedataset)
2. Define one or more [Benchmark](docs/api.md#kobebenchmark)
   and one or more [Federators](docs/api.md#kobefederator).
3. Define a [Experiment](docs/api.md#kobeexperiment).

A simple example can be found in the [examples](examples/) directory.
