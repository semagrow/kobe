# Add a new dataset server

This walkthrough illustrates the steps required from the *implementor of a
dataset server* in order to create a [DatasetTemplate] specification.

In KOBE, a dataset template is defined using a set of Docker images. Additional
parameters include the port and the path that the container will listen for
queries. Dataset templates are used in the [Benchmark] specifications in order
to define the dataset server of the federated endpoints.

## Prerequisites

In this walkthrough we assume that you already have already prepared a Docker
image that provides the SPARQL endpoint of the dataset server (e.g.,
https://hub.docker.com/r/openlink/virtuoso-opensource-7).

## Step 1. Prepare your Docker images

The first step is to provide a set of one or more Docker images that downloads
the dataset, loads the data, and starts the dataset server. Even though all this
functionality can be provided with a single image, we suggest to split the
various tasks into three separate images. More specifically:

* A docker image that downloads a RDF dump from a known URL (found in the
  variable `$DATASET_URL`) and extracts its contents in the directory
  `/kobe/dataset/$DATASET_NAME/dump`.
* A docker image that loads the downloaded dump (already present in the
  directory `/kobe/dataset/$DATASET_NAME/dump`) into the dataset server.
  Optionally, it can back-up the contents of the database in some directory
  inside `/kobe/dataset/$DATASET_NAME` such that the loading process to be
  executed only once. 
* A docker image that starts the dataset server and exposes its SPARQL endpoint.

The environment variables are initialized by the KOBE operator according to the
specification of the benchmark to be executed. Moreover, the shared volumes are
managed through the KOBE operator too (ref.
[here](https://github.com/semagrow/kobe/tree/devel/operator/docs/storage.md) for
details about the shared storage of KOBE).

> In the benchmark walkthrough, we
> [suggest](../use/create_benchmark.md#step-1-prepare-your-dataset-dumps) that
> the dataset dumps should follow a specific format. Therefore, feel free to use
> `semagrow/url-donwnloader` (source code
> [here](https://github.com/semagrow/kobe/tree/devel/dockers/url-downloader)) as
> your first image. However, if you optionally want your template to support
> more dataset dump formats, you can implement your own url downloader. 

As an example, we present the images for two dataset servers (namely Virtuoso
and Strabon).

* Regarding the [Virtuoso] RDF store, we use the images
  `semagrow/url-donwnloader`, `semagrow/virtuoso-init` (source code
  [here](https://github.com/semagrow/kobe/tree/devel//examples/dataset-virtuoso/virtuoso-init)),
  and `semagrow/virtuoso-main` (source code
  [here](https://github.com/semagrow/kobe/tree/devel//examples/dataset-virtuoso/virtuoso-main)).
  We use the shared storage of KOBE, to keep a backup of the `/database`
  directory of Virtuoso, which is used to keep all the files used by the
  database. The last two images are built upon `openlink/virtuoso-opensource-7`.

* Regarding the [Strabon] geospatial RDF store, we use the images
  `semagrow/url-downloader`, `semagrow/strabon-init` (source code
  [here](https://github.com/semagrow/kobe/tree/devel//examples/dataset-strabon/strabon-init)),
  and `semagrow/strabon-main` (source code
  [here](https://github.com/semagrow/kobe/tree/devel//examples/dataset-strabon/strabon-main)).
  We use the shared storage of KOBE, to keep a backup of the PostGIS database
  (directory `/var/lib/postgresql/9.4/main`) which is where the data are kept
  inside [Strabon]. The last two images are built using the docker file of
  KR-suite (see http://github.com/GiorgosMandi/KR-Suite-docker)`.

## Step 2. Prepare your YAML file

Once you have prepared the docker images, creating the dataset template
specification for your dataset server is a straightforward task. It should look
like this (we use as an example the template for Virtuoso):

```yaml
apiVersion: kobe.semagrow.org/v1alpha1
kind: DatasetTemplate
metadata:
  # Each dataset template can be uniquely identified by its name.
  name: virtuosotemplate
spec:
  initContainers:
    # here you put the first two images (that is the images for initializing
    # your server in the order you want to be executed).
    - name: initcontainer0
      image: semagrow/url-downloader
    - name: initcontainer1
      image: semagrow/virtuoso-init
  containers:
    # here you put the last image (that is the image for serving the data)
    - name: maincontainer
      image: semagrow/virtuoso-main
      ports:
        - containerPort: 8890  # port to listen for queries
  port: 8890     # port to listen for queries
  path: /sparql  # path to listen for queries
```

The default URL for the SPARQL endpoint for [Virtuoso] is
`http://localhost:8890/sparql`, hence the port and the path to listen for
queries are `8890` and `/sparql` respectively.

## Examples

We have already prepared several dataset template specifications to experiment
with:

* [dataset-virtuoso](https://github.com/semagrow/kobe/tree/devel/examples/dataset-virtuoso)
* [dataset-strabon](https://github.com/semagrow/kobe/tree/devel/examples/dataset-strabon)

!!! note
    We plan to define more dataset template specifications in the future. We
    place all dataset template specifications in the
    [examples/](https://github.com/semagrow/kobe/tree/devel/examples/) directory
    under a subdirectory with the prefix `dataset-*`. 

[Benchmark]: ../references/api.md#benchmark
[DatasetTemplate]: ../references/api.md#datasettemplate

[Virtuoso]: https://virtuoso.openlinksw.com/
[Strabon]: http://strabon.di.uoa.gr/