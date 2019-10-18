# kobe-operator

kobe-operator is a Kubernetes operator to automate the benchmarking of federated 
query processors.

## Installing the Chart

TBA.

To install the chart with the release name `my-release`:
```
helm repo add kobe https://charts.kobe.semagrow.github.io

helm install --name my-release --namespace kobe-operator kobe/kobe-operator
```

## Uninstalling the Chart

To uninstall/delete the `my-release` deployment:
```
helm delete my-release
```
The command removes all the Kubernetes components associated with the chart and deletes the release

## Configuration

The following table lists the configuration parameters of the kobe-operator chart and their default values

| Parameter | Description | Default |
| --------- | ----------- | ------- |
| `global.imagePullSecrets` | Reference to one or more secrets to be used when pulling images | `[]` |
| `global.rbac.create` | If `true`, create and use RBAC resources | `true` |
| `image.registry` | Image registry | `docker.io` |
| `image.repository` | Image repository | `kostbabis/kobe-operator` |
| `image.tag` | Image tag | `0.1.0` |
| `image.pullPolicy` | Image pull policy | `IfNotPresent` |
| `replicaCount`  | Number of cert-manager replicas  | `1` |
| `extraEnv` | Optional environment variables for kobe-manager | `[]` |
| `serviceAccount.create` | If `true`, create a new service account | `true` |
| `serviceAccount.name` | Service account to be used. If not set and `serviceAccount.create` is `true`, a name is generated using the fullname template |  |
| `resources` | CPU/memory resource requests/limits | |
| `nodeSelector` | Node labels for pod assignment | `{}` |
| `affinity` | Node affinity for pod assignment | `{}` |
| `tolerations` | Node tolerations for pod assignment | `[]` |
| `podAnnotations` | Annotations to add to the kobe-operator pod | `{}` |
| `podLabels` | Labels to add to the kobe-operator pod | `{}` |

Specify each parameter using the `--set key=value,[key=value]` argument to `helm install`

Alternatively, a YAML file that specifies the values for the above parameters can be provided while installing the chart.
For example,
```
helm install --name my-release -f values.yaml .
```
You can use the default [values.yaml](values.yaml).
