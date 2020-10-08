# Benchmark Walkthough

This walkthrough illustrates the steps required from the *benchmark designer* in order to create a [Benchmark](../operator/docs/api.md#benchmark) specification.

In KOBE, a benchmark comprises a collection of data sources, the latency and throughput of these data sources, and a list of query strings.
Benchmarks are defined independently of the federator that is being benchmarked.

## Prerequisites

In this walkthrough we assume that you already have already prepared the following:

* The dump of each RDF dataset of the benchmark.
* A list of query strings of the benchmark.
* A [DatasetTemplate](../operator/docs/api.md#datasettemplate) for each dataset server you want to use in your benchmark.

Regarding the third perquisite, we have already prepared several dataset template to use.
If you want to create your own dataset server template, check out [this guide](..).

## Step 1. Prepare your dataset dumps

Create a .tar.gz file for each dataset, and upload it on a known location.

Place all files of the dataset into a directory, put this directory into a tar file and compress it with gzip.
Even though most dataset engines support the import of several RDF formats (such as RDF/XML, turtle, etc), the most simple format is N-TRIPLES.
Therefore, we suggest to store your dataset in a single .nt file.
To prepare a dump.nt file, do the following:
```
mkdir dataset/
mv dump.nt dataset/
tar czvf dataset.tar.gz dataset/
```
Finally, upload the .tar.gz file on a known location.
As an example, we have uploaded the datasets for the FedBench experiment in the following [location](https://users.iit.demokritos.gr/~gmouchakis/dumps/).

## Step 2. Prepare your YAML file

A benchmark is characterized by its *name* and it defined using a list of *datasets* and a set of *queries*.
A typical benchmark specification should look like this:

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
    # ...
  
  # Each benchmark consists of a set of queries.
  queries:
    # Each query can be uniquely identified by its name,
    # and is defined with
    #    The language in which the query is written (e.g., SPARQL).
    #    The actual query string to be posed to the federator.
    - name: query1
      language: sparql
      queryString: "SELECT * WHERE ... "
    # ...
```

In the following we illustrate an example using the `toybench-simple` benchmark.
This benchmark contains three SPARQL queries (namely `tq1`,`tq2`, and `tq3`), and two datasets (namely `toy1` and `toy2`), both of them served by Virtuoso.

```yaml
apiVersion: kobe.semagrow.org/v1alpha1
kind: Benchmark
metadata:
  name: toybench-simple
spec:
  datasets:
    - name: toy1
      files:
        - url: https://users.iit.demokritos.gr/~antru/dumps/toy1.tar.gz
      templateRef: virtuosotemplate
    - name: toy2
      files:
        - url: https://users.iit.demokritos.gr/~antru/dumps/toy2.tar.gz
      templateRef: virtuosotemplate
  queries:
    - name: tq1
      language: sparql
      queryString: "SELECT * WHERE {
        <http://example.org/33> <http://www.w3.org/2002/07/owl#sameAs> ?o .
        ?o <http://purl.org/dc/terms/creator> ?c .
        ?o <http://example.org/value> ?v .
      }"
    - name: tq2
      language: sparql
      queryString: "SELECT * WHERE {
        ?o <http://purl.org/dc/terms/creator> <http://semagrow.org/antru> .
        ?s <http://www.w3.org/2002/07/owl#sameAs> ?o .
      }"
    - name: tq3
      language: sparql
      queryString: "SELECT * WHERE {
        ?s <http://www.w3.org/2002/07/owl#sameAs> ?o .
        ?o <http://example.org/value> ?v .
        FILTER (?v > 33)
      }"
```

## Examples

* [benchmark-fedbench](../examples/benchmark-fedbench)
* [benchmark-geofedbench](../examples/benchmark-geofedbench)
* [benchmark-geographica](../examples/benchmark-geographica)
* [benchmark-toybench](../examples/benchmark-toybench)
