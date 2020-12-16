
## Prerequisites

- `Kubernetes` >= 1.8.0
- `nfs-commons` installed in the nodes of the cluster. If in debian or
   Ubuntu you can install it using `apt-get install nfs-common`

## Installation of the Kubernetes operator

KOBE needs the Kubernetes operator to be installed in the Kubernetes cluster. To
quickly install the KOBE operator in a Kubernetes cluster. 

You can use the `kobectl` script found in the [bin](bin/) directory:

```
export PATH=`pwd`/bin:$PATH
kobectl install operator .
```

Alternatively, you could run the following commands:
```
kubectl apply -f operator/deploy/crds
kubectl apply -f operator/deploy/service_account.yaml
kubectl apply -f operator/deploy/clusterrole.yaml
kubectl apply -f operator/deploy/clusterrole_binding.yaml
kubectl apply -f operator/deploy/role.yaml
kubectl apply -f operator/deploy/operator.yaml
```

You will get a confirmation message that each resource has successfully been
created. This will set the operator running in your Kubernetes cluster and needs
to be done only once.

## Installation of Networking subsystem

KOBE uses istio to support network delays between the different deployments. To
install istio you can run the following:
```
kobectl install istio .
```
Alternatively, you can consult the official [installation
guide](https://istio.io/docs/setup/getting-started/) or you can type the
following commands.

```
curl -L https://istio.io/downloadIstio | sh -
export PATH=`pwd`/istio-1.6.0/bin:$PATH
istioctl manifest apply --set profile=default
```

## Installation of the Evaluation Metrics Extraction subsystem

To enable the evaluation metrics extraction subsystem, run
```
kobectl install efk .
```
or alternatively the following
```
helm repo add elastic https://helm.elastic.co
helm repo add kiwigrid https://kiwigrid.github.io
helm install elastic/elasticsearch --name elasticsearch --set persistence.enabled=false --set replicas=1 --version 7.6.2
helm install elastic/kibana --name kibana --set service.type=NodePort --version 7.6.2
helm install --name fluentd -f operator/deploy/efk-config/fluentd-values.yaml kiwigrid/fluentd-elasticsearch --version 8.0.1
kubectl apply -f operator/deploy/efk-config/kobe-kibana-configuration.yaml
```

These result in the simplest setup of an one-node
[Elasticsearch](https://github.com/elastic/helm-charts/blob/master/elasticsearch)
that does not persist data across pod recreation, a
[Fluentd](https://github.com/kiwigrid/helm-charts/tree/master/charts/fluentd-elasticsearch)
`DaemonSet` and a
[Kibana](https://github.com/elastic/helm-charts/tree/master/kibana)
node that exposes a `NodePort`. 

After all pods are in Running state Kibana dashboards can be accessed at 
```
http://<NODE-IP>:<NODEPORT>/app/kibana#/dashboard/
``` 
where `<NODE-IP>` the IP of any of the Kubernetes cluster nodes and `<NODEPORT>`
the result of `kubectl get -o jsonpath="{.spec.ports[0].nodePort}" services
kibana-kibana`.

The setup can be customized by changing the configuration parameters of each
helm chart. Please check the corresponding documentation of each chart for more
info.
