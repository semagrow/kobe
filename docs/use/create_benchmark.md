# Create a new benchmark

This walkthrough illustrates the steps required from the *benchmark designer* in
order to create a [Benchmark](../references/api.md#benchmark) specification.

In KOBE, a benchmark comprises a collection of data sources, the latency of
these data sources, and a list of query strings. Benchmarks are defined
independently of the federator that is being benchmarked.

## Prerequisites

In this walkthrough we assume that you already have already prepared the
following:

* The dump of each RDF dataset of the benchmark.
* A list of query strings of the benchmark.
* A [DatasetTemplate](../references/api.md#datasettemplate) for each dataset
  server you want to use in your benchmark.

Regarding the third prerequisite, we have already prepared several dataset
templates to use. If you want to create your own dataset server template, check
out [this guide](../extend/add_dataset_server.md).

## Step 1. Prepare your dataset dumps

Create a .tar.gz file for each dataset, and upload it on a known location.

Place all files of the dataset into a directory, put this directory into a tar
file and compress it with gzip. Even though most dataset engines support the
import of several RDF formats (such as RDF/XML, turtle, etc), the most simple
format is N-TRIPLES. Therefore, we suggest to store your dataset in a single .nt
file. If you choose to to prepare a dump.nt file, just do the following:

```
mkdir dataset/
mv dump.nt dataset/
tar czvf dataset.tar.gz dataset/
```

Finally, upload the .tar.gz file on a known location. As an example, we have
uploaded the datasets for the FedBench experiment in the following
[location](https://users.iit.demokritos.gr/~gmouchakis/dumps/).

## Step 2. Prepare your YAML file

A benchmark is characterized by its *name* and is parameterized using a list of
*datasets* and a set of *queries*. A typical benchmark specification should look
like this:

```yaml
apiVersion: kobe.semagrow.org/v1alpha1
kind: Benchmark
metadata:
  # Each benchmark can be uniquely identified by its name.
  name: mybench
spec:
  # Each benchmark consists of a set of dataset specifications.
  datasets:
    # Each dataset can be uniquely identified by its name,
    # and is defined with
    #   A list of URLs that contain the dump of the dataset to download.
    #   A specification of the dataset server to use (dataset template).
    - name: dataset1
      files:
        - url: https://path/to/download/the/dataset1.tar.gz
      templateRef: datasettemplate
    # ... add more datasets ...
  
  # Each benchmark consists of a set of queries.
  queries:
    # Each query can be uniquely identified by its name,
    # and is defined with
    #    The language in which the query is written (e.g., SPARQL).
    #    The actual query string to be posed to the federator.
    - name: query1
      language: sparql
      queryString: "SELECT * WHERE ... "
    # ... add more queries ...
```

Check the following link in which we illustrate a simple example of the above
specification:

* [benchmark-toybench/toybench-simple.yaml](https://github.com/semagrow/kobe/tree/devel/examples/benchmark-toybench/toybench-simple.yaml)

This benchmark contains three SPARQL queries (namely `tq1`,`tq2`, and `tq3`),
and two datasets (namely `toy1` and `toy2`), both of them served by Virtuoso.

## Examples

We have already prepared several benchmark specifications to experiment with:

* [benchmark-fedbench](https://github.com/semagrow/kobe/tree/devel/examples/benchmark-fedbench)
* [benchmark-geofedbench](https://github.com/semagrow/kobe/tree/devel/examples/benchmark-geofedbench)
* [benchmark-geographica](https://github.com/semagrow/kobe/tree/devel/examples/benchmark-geographica)
* [benchmark-toybench](https://github.com/semagrow/kobe/tree/devel/examples/benchmark-toybench)

> Notice: We plan to define more benchmark specifications in the future. We
> place all benchmark specifications in the [examples/](https://github.com/semagrow/kobe/tree/devel/examples/) directory
> under a subdirectory with the prefix `benchmark-*`. 
