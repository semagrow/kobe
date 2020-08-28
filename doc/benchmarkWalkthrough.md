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

## Step 2. Prepare your YAML file

