# Benchmark Walkthough

This walkthrough illustrates the steps required from the *benchmark designer* in order to create a [Benchmark](../operator/docs/api.md#benchmark) specification.

In KOBE, a benchmark comprises a collection of data sources, the latency of these data sources, and a list of query strings.
Benchmarks are defined independently of the federator that is being benchmarked.

## Prerequisites

In this walkthrough we assume that you already have already prepared the following:

* The dump of each RDF dataset of the benchmark.
* A list of query strings of the benchmark.
* A [DatasetTemplate](../operator/docs/api.md#datasettemplate) for each dataset server you want to use in your benchmark.

Regarding the third prerequisite, we have already prepared several dataset templates to use.
If you want to create your own dataset server template, check out [this guide](DatasetTemplateWalkthough.md).

## Step 1. Prepare your dataset dumps

Create a .tar.gz file for each dataset, and upload it on a known location.

Place all files of the dataset into a directory, put this directory into a tar file and compress it with gzip.
Even though most dataset engines support the import of several RDF formats (such as RDF/XML, turtle, etc), the most simple format is N-TRIPLES.
Therefore, we suggest to store your dataset in a single .nt file.
If you choose to to prepare a dump.nt file, just do the following:
```
mkdir dataset/
mv dump.nt dataset/
tar czvf dataset.tar.gz dataset/
```
Finally, upload the .tar.gz file on a known location.
As an example, we have uploaded the datasets for the FedBench experiment in the following [location](https://users.iit.demokritos.gr/~gmouchakis/dumps/).

## Step 2. Prepare your YAML file

A benchmark is characterized by its *name* and is parameterized using a list of *datasets* and a set of *queries*.
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

Check the following link in which we illustrate a simple example of the above specification:

* [benchmark-toybench/toybench-simple.yaml](../examples/benchmark-toybench/toybench-simple.yaml)

This benchmark contains three SPARQL queries (namely `tq1`,`tq2`, and `tq3`), and two datasets (namely `toy1` and `toy2`), both of them served by Virtuoso.

## Step 3 - Optional. Inject latency for each source endpoint

KOBE allows simulating network traffic for all sources of the benchmark.
For every source dataset of the benchmark, you can:

* inject delay in the connection between the given source endpoint and the *federation engine*.
* inject delay in the connection between the given source endpoint and *another source endpoint*.

> The reason for injecting delays between the federated sources is the fact that every SPARQL endpoint can issue a SPARQL query to every other endpoint using the SERVICE SPARQL keyword.

The latency of each source can be configured using the following [delay parameters](../operator/docs/api.md#delay).
The functionality of these parameters is offered by Istio.
Check this [link](https://istio.io/latest/docs/reference/config/networking/virtual-service/#HTTPFaultInjection-Delay) for more information.

* The `fixedDelaySec` and `fixedDelayMSec` are used to indicate the *amount of delay* in seconds and in milliseconds.
* The `percentage` field can be used to only delay a certain *percentage of requests*.

You can extend your benchmark specification can be extended in order to define the latency of the sources as follows:

```yaml
# In this example we will use two datasets, ds1 and ds2.
spec:
  datasets:
    - name: ds1
      # adds 1 second of delay before forwarding all responces to the federator
      federatorConnection:
         delayInjection:
           fixedDelaySec: 1
           percentage: 100
      networkTopology:
        # adds 2 sec of delay before forwarding the 50% of responces to the source ds1
        - datasetSource: ds2
          delayInjection:
            fixedDelaySec: 2
            percentage: 50
      # ... add remaining parameters for ds1
      
    - name: ds2
      # ... add remaining parameters for ds2
```

Check the following link in which we illustrate a simple working example with delays:

* [benchmark-toybench/toybench-delays.yaml](../examples/benchmark-toybench/toybench-delays.yaml)

This benchmark contains three SPARQL queries and two datasets (namely `toy1` and `toy2`).
All responces from `toy1` to the federator are delayed by 2 seconds and 150 milliseconds,
all responces from `toy2` to the federator are delayed by 2 seconds, and
the 50% of the responces from `toy1` to `toy2` are delayed by 3 seconds.


## Examples

We have already prepared several benchmark specifications to experiment with:

* [benchmark-fedbench](../examples/benchmark-fedbench)
* [benchmark-geofedbench](../examples/benchmark-geofedbench)
* [benchmark-geographica](../examples/benchmark-geographica)
* [benchmark-toybench](../examples/benchmark-toybench)

> Notice: We plan to define more benchmark specifications in the future.
> We place all benchmark specifications in the [examples/](../examples/) directory
> under a subdirectory with the prefix `benchmark-*`. 
