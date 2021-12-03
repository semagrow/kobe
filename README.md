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
- `kubectl` configured for the Kubernetes cluster
- `Helm` version 3 (for the Evaluation Metrics Extraction subsystem)
- `nfs-commons` installed in the nodes of the cluster. If in Debian or
   Ubuntu you can install it using `apt-get install nfs-common`

### Installation of the Kubernetes operator

KOBE needs the Kubernetes operator that needs to be installed in the
Kubernetes cluster. To quickly install the KOBE operator in a
Kubernetes cluster, you can use the `kobectl` script found in the
[bin](bin/) directory:

```
export PATH=`pwd`/bin:$PATH
kobectl install operator .
```
If you are using kubernetes version 1.15 and below you should instead use 
```
kobectl install operator-v1beta1 
```

Alternatively, you could run the following commands:

```
kubectl apply -f operator/deploy/crds
kubectl apply -f operator/deploy/service_account.yaml
kubectl apply -f operator/deploy/clusterrole.yaml
kubectl apply -f operator/deploy/clusterrole_binding.yaml
kubectl apply -f operator/deploy/operator.yaml
```
For Kubernetes version 1.15 and below  swap

```
kubectl apply -f operator/deploy/crds
```
with
```
kubectl apply -f operator/deploy/crds-v1beta1
```

You will get a confirmation message that each resource has
successfully been created.
This will set the operator running in your Kubernetes cluster and
needs to be done only once.

### Installation of Networking subsystem

KOBE uses [Istio](https://istio.io/) to support network delays between the different 
deployments. To install Istio first define the version (KOBE was tested with version 1.11.3)

```
export ISTIO_VERSION=1.11.3
```

then deploy Istio:
```
kobectl install istio .
```
Alternatively, you can consult the official 
[installation guide](https://istio.io/docs/setup/getting-started/) 
or you can type the following commands.

```
curl -L https://istio.io/downloadIstio | sh -
./istio-*/bin/istioctl manifest apply --set profile=default
```

### Installation of the Evaluation Metrics Extraction subsystem

To enable the evaluation metrics extraction subsystem, run
```
kobectl install efk .
```
or alternatively the following
```
helm repo add elastic https://helm.elastic.co
helm repo add kiwigrid https://kiwigrid.github.io
helm install elasticsearch elastic/elasticsearch --set persistence.enabled=false --set replicas=1 --version 7.6.2
helm install elasticsearch elastic/elasticsearch --set persistence.enabled=false --set replicas=1 --version 7.6.2
helm install fluentd kiwigrid/fluentd-elasticsearch -f operator/deploy/efk-config/fluentd-values.yaml --version 8.0.1
kubectl apply -f operator/deploy/efk-config/kobe-kibana-configuration.yaml
```

These result in the simplest setup of an single-node
[Elasticsearch](https://github.com/elastic/helm-charts/blob/master/elasticsearch)
that does not persist data across pod recreation, a
[Fluentd](https://github.com/kiwigrid/helm-charts/tree/master/charts/fluentd-elasticsearch)
`DaemonSet` and a
[Kibana](https://github.com/elastic/helm-charts/tree/master/kibana)
node that exposes a `NodePort`. 

After all pods are in `Running` state Kibana dashboards can be accessed
at 
```
http://<NODE-IP>:<NODEPORT>/app/kibana#/dashboard/
``` 
where `<NODE-IP>` the IP of any of the Kubernetes worker nodes and
`<NODEPORT>` the result of `kubectl get -o
jsonpath="{.spec.ports[0].nodePort}" services kibana-kibana`.

The setup can be customized by changing the configuration parameters
of each helm chart. Please check the corresponding documentation of
each chart for more info.

## Example

The typical workflow of defining a KOBE experiment is the following.
1. Create one [DatasetTemplate](operator/docs/api.md#datasettemplate)
   for each dataset server you want to use in your benchmark.
2. Define your [Benchmark](operator/docs/api.md#benchmark),
   which should contain a list of datasets and a list of queries.
2. Create one [FederatorTemplate](operator/docs/api.md#federatortemplate)
   for the federator engine you want to use in your experiment. 
3. Define an [Experiment](operator/docs/api.md#experiment) over your previously defined benchmark.

Several examples of the above specifications can be found in the [examples](examples/) directory.

In the following, we show the steps for deploying an experiment on a simple benchmark that comprises
three queries over a Semagrow federation of two Virtuoso endpoints.

You can use the `kobectl` script found in the [bin](bin/) directory for controlling your experiments:

```
export PATH=`pwd`/bin:$PATH
kobectl help
```

First, apply the templates for Virtuoso and Semagrow:

```
kobectl apply examples/dataset-virtuoso/virtuosotemplate.yaml
kobectl apply examples/federator-semagrow/semagrowtemplate.yaml
```
Then, apply the benchmark.

```
kobectl apply examples/benchmark-toybench/toybench-simple.yaml
```
Before running the experiment, you should verify that the datasets are loaded.
Use the following command:
```
kobectl show benchmark toybench-simple
```
When the datasets are loaded, you should get the following output:
```
NAME  STATUS
toy1  Running
toy2  Running
```
Proceed now with the execution of the experiment:
```
kobectl apply examples/experiment-toybench/toyexp-simple.yaml
```
As perviously, you can review the state of the experiment with the following command:
```
kobectl show experiment toyexp-simple
```
You can now view the evaluation metrics in the Kibana dashboards.

For removing all of the above, issue the following commands:
```
kobectl delete experiment toyexp-simple
kobectl delete benchmark toybench-simple
kobectl delete federatortemplate semagrowtemplate
kobectl delete datasettemplate virtuosotemplate
```
For more advanced control options for KOBE, use [kubectl](https://kubernetes.io/docs/reference/kubectl/overview/).

## Removal

To remove KOBE from your cluster, run the following command:
```
kobectl purge .
```
To remove KOBE operator manually, run
```
kubectl delete -f operator/deploy/operator.yaml
kubectl delete -f operator/deploy/role.yaml
kubectl delete -f operator/deploy/clusterrole_binding.yaml
kubectl delete -f operator/deploy/clusterrole.yaml
kubectl delete -f operator/deploy/service_account.yaml
kubectl delete -f operator/deploy/crds
```
To remove Istio manually, run
```
./istio-*/bin/istioctl manifest generate --set profile=default | kubectl delete -f -
kubectl delete namespace istio-system
```
To remove the evaluation metrics extraction subsystem manually, run
```
helm delete --purge elasticsearch
helm delete --purge kibana
helm delete --purge fluentd
helm repo remove elastic
helm repo remove kiwigrid
kubectl delete jobs.batch kobe-kibana-configuration
kubectl delete configmaps kobe-kibana-config
```
and then in each Kubernetes node
```
rm -rf /var/log/fluentd-buffers/
rm /var/log/containers.log.pos
```

