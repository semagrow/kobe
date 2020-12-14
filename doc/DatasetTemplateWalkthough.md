# Dataset Template Walkthough

This walkthrough illustrates the steps required from the *implementor of a dataset server* in order to create a [DatasetTemplate](../operator/docs/api.md#datasettemplate) specification.

In KOBE, a dataset template is defined using a set of Docker images.
Additional parameters include the port and the path that the container will listen for queries.
Dataset templates are used in the [Benchmark](../operator/docs/api.md#benchmark) specifications in order to define the dataset server of the federated endpoints.

## Prerequisites

In this walkthrough we assume that you already have already prepared a Docker image that provides the SPARQL endpoint of the dataset server (e.g., https://hub.docker.com/r/openlink/virtuoso-opensource-7).

## Step 1. Prepare your Docker images

The first step is to provide a set of one or more Docker images that downloads the dataset, loads the data, and starts the dataset server.
Even though all this functionality can be provided with a single image, we suggest to split the various tasks into three separate images.
More specifically:

* A docker image that downloads a RDF dump from a known URL (found in the variable `$DATASET_URL`) and extracts its contents in the directory `/kobe/dataset/$DATASET/dump`.
* A docker image that loads the downloaded dump (already present in the directory `/kobe/dataset/$DATASET/dump`) into the dataset server.
  Optionally, it can back-up the contents of the database in some directory inside `/kobe/dataset/$DATASET/` such that the loading process to be executed only once. 
* A docker image that starts the dataset server and exposes its SPARQL endpoint.

The environment variables are initialized by the Kobe operator according to the specification of the benchmark to be executed.
Moreover, the shared volumes are managed through the Kobe operator too (ref. [here](../operator/docs/storage.md) for details about the shared storage of Kobe).

## Step 2. Prepare your YAML file


## Examples

We have already prepared several dataset template specifications to experiment with:

* [dataset-virtuoso](../examples/dataset-virtuoso)
* [dataset-strabon](../examples/dataset-strabon)

> Notice: We plan to define more dataset template specifications in the future.
> We place all benchmark specifications in the [examples/](../examples/) directory
> under a subdirectory with the prefix `dataset-*`. 

