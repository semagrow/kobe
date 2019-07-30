# API Docs

This Document documents the types introduced by the KOBE Operator to be consumed by users.

## Table of Contents
* [KobeDataset](#kobedataset)
* [KobeBenchmark](#kobebenchmark)
* [KobeFederator](#kobefederator)
* [KobeExperiment](#kobeexperiment)

## KobeDataset

KobeDataset defines a dataset that could be used in an experiment.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| image | Image of the database system. Currently fixed to kostbabis/virtuoso | string | true |
| forceLoad | Forces to download and load from dump files | boolean | false |
| downloadFrom | the dump location. | url | true |
| count | how many instances of this database you want in your cluster (under same service). | integer | false |

[Back to TOC](#table-of-contents)

## KobeBenchmark

KobeBenchmark defines a benchmark in kobe.
A benchmark consists of a set of datasets that must be already defined with the [KobeDataset](#kobedataset). 
It also contains the definition of one or more [SPARQL](https://www.w3.org/TR/sparql11-query/) queries 
that are going to get tested against the datasets in the benchmark.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| datasets | Datasets to be used for the benchmark | [][KobeDataset](#kobedataset) | true |
| queries | Query set | []Query | true |

- Under `spec.datasets.name[*]` you must write down the name of the datasets your benchmark will include. 
  The names must be the same as the `metadata.name` of the KobeDataset custom resources defined above.
- Under `spec.queries[*]` you must write down the queries of your benchmark. Query name is the name of the query. 
  The field `language` for now should always be set to `sparql` and `queryString` should be the string that contains your query.

[Back to TOC](#table-of-contents)

## KobeFederator

KobeFederator defines a federator.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| image | Image of the federator. | string | true |
| imagePullPolicy | Image pull policy of pulling the image. One of Always, Never, IfNotPresent. | string | true |
| port | Number of port to expose on the host. | [ContainerPort](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.11/#containerport-v1-core) | true |
| sparqlEnding | The suffix of your federators sparql endpoint | string | true |
| fedConfDir | The directory your federator expects to find its metadata files in order to operate properly. | string | true |
| confFromFileImage | The image that configures the federator. | string | true |
| inputDumpDir | | string | true |
| outputDumpDir |  | string | true |
| confImage |  | string | true |
| inputDir |  | string | true |
| outputDir |  | string | true |

- Under `spec.confFromFileImage` you must provide the name of an image that does the following.
  It creates a container that reads from `/kobe/input_dump` files of a dataset and 
  writes at `/kobe/output_metadata` configuration files for that dataset.
  It can also instead query directly the database SPARQL endpoint to create 
  the metadata file since we provide the init container with an environment 
  variable called `END_POINT` which contains the full url of the SPARQL endpoint of the dataset
  The image should be oblivious of what dataset it makes the metadata for and incorporate 
  only the necessary logic to make that file. For example with semagrow we provide an 
  image that uses the `sevod-scraper` (check it under semagrow in github) to process 
  the dump files of a dataset (f.e dbpedia) and return a dbpedia.ttl file for this specific set.
  The read and write directories of your image can be changed from the following two 
  fields in the yaml `spec.inputDumpDir` and `spec.outputDumpDir` if its convenient.
  They automatically default to `/kobe/input` and `/kobe/output` respectively.
- Under `spec.ConfImage` you must provide the name of an image that does the following.
  It reads from `/kobe/input` a set of different metadata files and combines them 
  to one big configuration file of metadata for the benchmark. Your image should 
  not care about what datasets the files belong to and only do the union of them.
  For example, with semagrow we just need to turn each dataset metadata from `.ttl` to `.nt` 
  then concatenate them and turn them back to `.ttl`. Again if you want to change 
  the input and output directories your image expects to find the files and write to, 
  you can with the following fields `spec.inputDir` and `spec.outputDir`.

[Back to TOC](#table-of-contents)

## KobeExperiment

KobeExperiment defines the actual experiment. It consists of a [KobeFederator](#kobefederator)
that will get benchmarked. Also it requires the name of a [KobeBenchmark](#kobebenchmark) that will be used.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| benchmark | The benchmark name. It must be the same as the name of the [KobeBenchmark](#kobebenchmark). | [KobeBenchmark](#kobebenchmark) | true |
| federator | The federator name. It must be the same as the name of the [KobeFederator](#kobefederator). | [KobeFederator](#kobefederator) | false |
| timesToRun | The number of times you want the benchmark experiment to repeat | integer | false |
| dryRun | If set to true the federation will be created and the federator initialized. The health checks will also happen but the experiment will hang there and no evaluation job will run till this flag is changed. | boolean | false |
| forceNewInit | if set to true it will always try to run the init image that create a metadata file from a dataset for this federator. If set to false it will check and use pre-existing metadata files if they exist for a pair of dataset and federator. It can be used to save time since metadata extraction for a big dataset take a long time and makes sense to not repeat this process. This affects only the first init process with the image that makes a metadata file from a dataset dump or endpoint. The second init process that combines many init files to one will always run again before init complete. | boolean | false |
| evalImage | The image of the evaluator | string | false |

[Back to TOC](#table-of-contents)
