
# API Docs

This Document documents the types introduced by the Kobe Operator to be consumed by users.

> Note this document is generated from code comments. When contributing a change to this document please do so by changing the code comments.

## Table of Contents
* [Benchmark](#benchmark)
* [BenchmarkList](#benchmarklist)
* [BenchmarkSpec](#benchmarkspec)
* [Dataset](#dataset)
* [DatasetList](#datasetlist)
* [DatasetSpec](#datasetspec)
* [DatasetStatus](#datasetstatus)
* [Evaluator](#evaluator)
* [Experiment](#experiment)
* [ExperimentList](#experimentlist)
* [ExperimentSpec](#experimentspec)
* [ExperimentStatus](#experimentstatus)
* [Federation](#federation)
* [FederationList](#federationlist)
* [FederationSpec](#federationspec)
* [FederationStatus](#federationstatus)
* [Federator](#federator)
* [FederatorList](#federatorlist)
* [FederatorTemplate](#federatortemplate)
* [KobeUtil](#kobeutil)
* [KobeUtilList](#kobeutillist)
* [Query](#query)

## Benchmark

Benchmark is the Schema for the benchmarks API

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#objectmeta-v1-meta) | false |
| spec |  | [BenchmarkSpec](#benchmarkspec) | false |
| status |  | [BenchmarkStatus](#benchmarkstatus) | false |

[Back to TOC](#table-of-contents)

## BenchmarkList

BenchmarkList contains a list of Benchmark

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#listmeta-v1-meta) | false |
| items |  | [][Benchmark](#benchmark) | true |

[Back to TOC](#table-of-contents)

## BenchmarkSpec

BenchmarkSpec defines the components for this benchmark setup

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| datasets |  | []string | true |
| queries |  | [][Query](#query) | true |

[Back to TOC](#table-of-contents)

## Dataset

Dataset is the Schema for the kobedatasets API

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#objectmeta-v1-meta) | false |
| spec |  | [DatasetSpec](#datasetspec) | false |
| status |  | [DatasetStatus](#datasetstatus) | false |

[Back to TOC](#table-of-contents)

## DatasetList

DatasetList contains a list of KobeDataset

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#listmeta-v1-meta) | false |
| items |  | [][Dataset](#dataset) | true |

[Back to TOC](#table-of-contents)

## DatasetSpec

DatasetSpec defines the desired state of Dataset

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| image | Docker image name. More info: https://kubernetes.io/docs/concepts/containers/images | string | true |
| imagePullPolicy | Image pull policy. One of Always, Never, IfNotPresent. Defaults to Always if :latest tag is specified, or IfNotPresent otherwise. Cannot be updated. More info: https://kubernetes.io/docs/concepts/containers/images#updating-images | v1.PullPolicy | true |
| replicas | Replicas is the number of desired replicas. This is a pointer to distinguish between explicit zero and unspecified. Defaults to 1. More info: https://kubernetes.io/docs/concepts/workloads/controllers/replicationcontroller#what-is-a-replicationcontroller | *int32 | true |
| forceLoad | Forces to download and load from dataset file | bool | true |
| downloadFrom | A URL that points to the compressed dataset file | string | true |
| port | Number of port to expose on the host. If specified, this must be a valid port number, 0 < x < 65536. | int32 | true |
| path | Path that the container will listen for queries | string | true |
| env | List of environment variables to set in the container. Cannot be updated. | [][v1.EnvVar](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#envvar-v1-core) | false |
| affinity | If specified, the pod's scheduling constraints | *[v1.Affinity](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#affinity-v1-core) | false |
| resources | Resources are not allowed for ephemeral containers. Ephemeral containers use spare resources already allocated to the pod. | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#resourcerequirements-v1-core) | false |

[Back to TOC](#table-of-contents)

## DatasetStatus

DatasetStatus defines the observed state of Dataset

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| podNames |  | []string | true |
| phase |  | string | true |

[Back to TOC](#table-of-contents)

## Evaluator

Evaluator defines the

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| image |  | string | true |
| imagePullPolicy |  | v1.PullPolicy | true |
| command |  | []string | true |
| parallelism |  | int32 | true |

[Back to TOC](#table-of-contents)

## Experiment

Experiment is the Schema for the experiments API

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#objectmeta-v1-meta) | false |
| spec |  | [ExperimentSpec](#experimentspec) | false |
| status |  | [ExperimentStatus](#experimentstatus) | false |

[Back to TOC](#table-of-contents)

## ExperimentList

ExperimentList contains a list of Experiment

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#listmeta-v1-meta) | false |
| items |  | [][Experiment](#experiment) | true |

[Back to TOC](#table-of-contents)

## ExperimentSpec

ExperimentSpec defines the desired state of Experiment

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| benchmark |  | string | true |
| federator |  | string | true |
| evaluator |  | [Evaluator](#evaluator) | true |
| timesToRun |  | int | true |
| restartPolicy |  | RestartPolicy | false |
| dryRun |  | bool | true |
| forceNewInit |  | bool | true |

[Back to TOC](#table-of-contents)

## ExperimentStatus

ExperimentStatus defines the observed state of Experiment

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| startTime | Time at which this workflow started | [metav1.Time](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#time-v1-meta) | false |
| completionTime | Time at which this workflow completed | [metav1.Time](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#time-v1-meta) | false |
| run | The current iteration of the experiment It should be zero if not started yet | int | false |
| phase | The phase of the experiment | ExperimentPhase | true |

[Back to TOC](#table-of-contents)

## Federation

Federation is the Schema for the federations API

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#objectmeta-v1-meta) | false |
| spec |  | [FederationSpec](#federationspec) | false |
| status |  | [FederationStatus](#federationstatus) | false |

[Back to TOC](#table-of-contents)

## FederationList

FederationList contains a list of Federation

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#listmeta-v1-meta) | false |
| items |  | [][Federation](#federation) | true |

[Back to TOC](#table-of-contents)

## FederationSpec

FederationSpec defines the desired state of Federation

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| federatorName |  | string | true |
| template |  | [FederatorTemplate](#federatortemplate) | true |
| endpoints |  | []string | true |
| datasets |  | []string | true |
| forceNewInit |  | bool | true |
| init |  | bool | true |

[Back to TOC](#table-of-contents)

## FederationStatus

FederationStatus defines the observed state of KobeFederation

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| podNames | INSERT ADDITIONAL STATUS FIELD - define observed state of cluster Important: Run \"operator-sdk generate k8s\" to regenerate code after modifying this file Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html | []string | true |
| phase |  | string | true |

[Back to TOC](#table-of-contents)

## Federator

Federator is the Schema for the federators API

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#objectmeta-v1-meta) | false |
| spec |  | [FederatorSpec](#federatorspec) | false |

[Back to TOC](#table-of-contents)

## FederatorList

FederatorList contains a list of Federator

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#listmeta-v1-meta) | false |
| items |  | [][Federator](#federator) | true |

[Back to TOC](#table-of-contents)

## FederatorTemplate

FederatorTemplate defines

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| image | Docker image name. More info: https://kubernetes.io/docs/concepts/containers/images | string | true |
| imagePullPolicy | Image pull policy. One of Always, Never, IfNotPresent. Defaults to Always if :latest tag is specified, or IfNotPresent otherwise. Cannot be updated. More info: https://kubernetes.io/docs/concepts/containers/images#updating-images | v1.PullPolicy | true |
| port | Number of port to expose on the host. If specified, this must be a valid port number, 0 < x < 65536. | int32 | true |
| path | suffix to be added to endpoint of federator f.e ../SemaGrow/sparql | string | true |
| affinity | If specified, the pod's scheduling constraints | *[v1.Affinity](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#affinity-v1-core) | false |
| resources | Resources are not allowed for ephemeral containers. Ephemeral containers use spare resources already allocated to the pod. | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#resourcerequirements-v1-core) | false |
| confFromFileImage | The Docker image that receives a compressed dataset and may produce configuration needed from the federator to federate this specific dataset This container will run one time for each dataset in the federation | string | true |
| inputDumpDir | where the above image expects the dump to be (if from dump) | string | true |
| outputDumpDir | where the above image will place its result config file | string | true |
| confImage | The Docker image that initializes the federator (equivalent to initContainers) | string | true |
| inputDir |  | string | true |
| outputDir |  | string | true |
| fedConfDir | which directory the federator needs the metadata config files in order to find them | string | true |

[Back to TOC](#table-of-contents)

## KobeUtil

KobeUtil is the Schema for the kobeutils API

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#objectmeta-v1-meta) | false |

[Back to TOC](#table-of-contents)

## KobeUtilList

KobeUtilList contains a list of KobeUtil

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#listmeta-v1-meta) | false |
| items |  | [][KobeUtil](#kobeutil) | true |

[Back to TOC](#table-of-contents)

## Query

Query contains the query info

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| name |  | string | true |
| language |  | string | true |
| queryString |  | string | true |

[Back to TOC](#table-of-contents)
