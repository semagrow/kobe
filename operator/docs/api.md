
# API Docs

This Document documents the types introduced by the Kobe Operator to be consumed by users.

> Note this document is generated from code comments. When contributing a change to this document please do so by changing the code comments.

## Table of Contents
* [Benchmark](#benchmark)
* [BenchmarkList](#benchmarklist)
* [BenchmarkSpec](#benchmarkspec)
* [Dataset](#dataset)
* [DatasetEndpoint](#datasetendpoint)
* [DatasetFile](#datasetfile)
* [DatasetTemplate](#datasettemplate)
* [DatasetTemplateList](#datasettemplatelist)
* [Delay](#delay)
* [EphemeralDataset](#ephemeraldataset)
* [EphemeralDatasetList](#ephemeraldatasetlist)
* [EphemeralDatasetStatus](#ephemeraldatasetstatus)
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
* [FederatorSpec](#federatorspec)
* [FederatorTemplate](#federatortemplate)
* [FederatorTemplateList](#federatortemplatelist)
* [KobeUtil](#kobeutil)
* [KobeUtilList](#kobeutillist)
* [NetworkConnection](#networkconnection)
* [Query](#query)
* [SystemDatasetSpec](#systemdatasetspec)

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
| datasets |  | [][Dataset](#dataset) | true |
| queries |  | [][Query](#query) | true |

[Back to TOC](#table-of-contents)

## Dataset



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| name |  | string | true |
| files |  | [][DatasetFile](#datasetfile) | true |
| systemspec |  | *[SystemDatasetSpec](#systemdatasetspec) | false |
| templateRef |  | string | false |
| affinity | If specified, the pod's scheduling constraints | *[v1.Affinity](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#affinity-v1-core) | false |
| resources | Resources are not allowed for ephemeral containers. Ephemeral containers use spare resources already allocated to the pod. | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#resourcerequirements-v1-core) | false |
| networkTopology | network delays | [][NetworkConnection](#networkconnection) | false |
| federatorConnection |  | *[NetworkConnection](#networkconnection) | false |

[Back to TOC](#table-of-contents)

## DatasetEndpoint



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| host |  | string | true |
| namespace |  | string | true |
| port |  | uint32 | true |
| path |  | string | true |

[Back to TOC](#table-of-contents)

## DatasetFile



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| url |  | string | true |
| checksum |  | string | false |

[Back to TOC](#table-of-contents)

## DatasetTemplate



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#objectmeta-v1-meta) | false |
| spec |  | [SystemDatasetSpec](#systemdatasetspec) | false |

[Back to TOC](#table-of-contents)

## DatasetTemplateList

FederatorList contains a list of Federator

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#listmeta-v1-meta) | false |
| items |  | [][DatasetTemplate](#datasettemplate) | true |

[Back to TOC](#table-of-contents)

## Delay



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| fixedDelaySec | Add a fixed delay before forwarding the request. Format: 1h/1m/1s/1ms. MUST be >=1ms. | *uint32 | false |
| fixedDelayMSec | Add a fixed delay before forwarding the request. Format: 1h/1m/1s/1ms. MUST be >=1ms. | *uint32 | false |
| percentage |  | *uint32 | false |
| percent |  | *uint32 | false |

[Back to TOC](#table-of-contents)

## EphemeralDataset

EphemeralDataset is the Schema for the kobedatasets API

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#objectmeta-v1-meta) | false |
| spec |  | [Dataset](#dataset) | false |
| status |  | [EphemeralDatasetStatus](#ephemeraldatasetstatus) | false |

[Back to TOC](#table-of-contents)

## EphemeralDatasetList

EphemeralDatasetList contains a list of EphemeralDataset

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#listmeta-v1-meta) | false |
| items |  | [][EphemeralDataset](#ephemeraldataset) | true |

[Back to TOC](#table-of-contents)

## EphemeralDatasetStatus



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| podNames |  | []string | true |
| phase |  | string | true |
| forceLoad |  | bool | true |

[Back to TOC](#table-of-contents)

## Evaluator

Evaluator defines the

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| image |  | string | true |
| imagePullPolicy |  | v1.PullPolicy | false |
| command |  | []string | false |
| parallelism |  | int32 | false |
| env |  | [][v1.EnvVar](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#envvar-v1-core) | false |

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
| federatorName |  | string | true |
| federatorSpec |  | *[FederatorSpec](#federatorspec) | false |
| federatorTemplateRef |  | string | false |
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
| spec |  | [FederatorSpec](#federatorspec) | true |
| datasets |  | [][DatasetEndpoint](#datasetendpoint) | true |
| topology |  | [][NetworkConnection](#networkconnection) | false |
| initPolicy |  | InitializationPolicy | false |

[Back to TOC](#table-of-contents)

## FederationStatus

FederationStatus defines the observed state of KobeFederation

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| podNames |  | []string | true |
| phase |  | FederationPhase | true |

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

## FederatorSpec

FederatorSpec contains all necessary information for a federator

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| initContainers |  | []v1.Container | false |
| containers |  | []v1.Container | true |
| port | Number of port to expose on the host. If specified, this must be a valid port number, 0 < x < 65536. | int32 | true |
| path | suffix to be added to endpoint of federator f.e ../SemaGrow/sparql | string | true |
| confFromFileImage | The Docker image that receives a compressed dataset and may produce configuration needed from the federator to federate this specific dataset This container will run one time for each dataset in the federation | string | true |
| inputDumpDir | where the above image expects the dump to be (if from dump) | string | true |
| outputDumpDir | where the above image will place its result config file | string | true |
| confImage | The Docker image that initializes the federator (equivalent to initContainers) | string | true |
| inputDir |  | string | true |
| outputDir |  | string | true |
| fedConfDir | which directory the federator needs the metadata config files in order to find them | string | true |

[Back to TOC](#table-of-contents)

## FederatorTemplate

FederatorTemplate defines a federator and its components that it needs to be installed.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#objectmeta-v1-meta) | false |
| spec |  | [FederatorSpec](#federatorspec) | false |

[Back to TOC](#table-of-contents)

## FederatorTemplateList



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#listmeta-v1-meta) | false |
| items |  | [][FederatorTemplate](#federatortemplate) | true |

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

## NetworkConnection



| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| datasetSource |  | *string | false |
| delayInjection |  | [Delay](#delay) | false |

[Back to TOC](#table-of-contents)

## Query

Query contains the query info

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| name |  | string | true |
| language |  | string | true |
| queryString |  | string | true |

[Back to TOC](#table-of-contents)

## SystemDatasetSpec

DatasetSpec defines the desired state of Dataset

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| importContainers |  | []v1.Container | false |
| initContainers | List of initialization containers belonging to the pod. Init containers are executed in order prior to containers being started. If any init container fails, the pod is considered to have failed and is handled according to its restartPolicy. The name for an init container or normal container must be unique among all containers. Init containers may not have Lifecycle actions, Readiness probes, or Liveness probes. The resourceRequirements of an init container are taken into account during scheduling by finding the highest request/limit for each resource type, and then using the max of of that value or the sum of the normal containers. Limits are applied to init containers in a similar fashion. Init containers cannot currently be added or removed. Cannot be updated. More info: https://kubernetes.io/docs/concepts/workloads/pods/init-containers/ | []v1.Container | false |
| containers | List of containers belonging to the pod. Containers cannot currently be added or removed. There must be at least one container in a Pod. Cannot be updated. | []v1.Container | true |
| initPolicy | Forces to download and load from dataset file | InitializationPolicy | true |
| port | Number of port to expose on the host. If specified, this must be a valid port number, 0 < x < 65536. | uint32 | true |
| path | Path that the container will listen for queries | string | true |

[Back to TOC](#table-of-contents)
