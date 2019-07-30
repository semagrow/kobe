# Design

This document describes the design and interactions between
the custom resource descriptions that the KOBE operator introduces.

The custom resources that the KOBE operator introduces are:
* `KobeDataset`
* `KobeFederator`
* `KobeBenchmark`
* `KobeExperiment`


## KobeDataset

The `KobeDataset` custom resource definition (CRD) declaratively
defines a dataset that can be used in an experiment. The operator
will create and mantain a pod that runs a `virtuoso` instance with that dataset. It will also cache the `db` file and dump files for future retrieval if the pod dies and restarts or if the user deletes the kobedataset and want to redefine it.

One dataset can be used in many experiments and needs to only be defined once.

## KobeFederator

The `KobeFederator` custom resource definition (CRD) defines a 
federator that is able to federate `KobeDataset`s. The resource 
has the option to define the image of the federator to use and 
two initialization images to configure the federator given the 
datasets.

## KobeBenchmark

The `KobeBenchmark` custom resource definition (CRD) defines a 
benchmark that comprises a set of datasets and a set of queries.
The `KobeBenchmark` is independent of a specific federator and 
can be used to benchmark different federators.

## KobeExperiment

The `KobeExperiment` custom resource definition (CRD) defines an 
incarnation of a benchmark. It pairs a `KobeBenchmark` with a `KobeFederator`. Moreover, it specifies an evaluator to use for issuing the queries of the `KobeBenchmark` to the `KobeFederator`.