# Experiment Walkthough

This walkthrough illustrates the steps required from the *experimenter* in order to create an [Experiment](../operator/docs/api.md#experiment) specification.

In KOBE, an experiment is defined over a benchmark and federator.
This resource provides the necessary parameters for instantiating a federation of querying systems.
The experiment also provides an evaluator, that is a piece of software that will pose the queries to the federator.

## Prerequisites

In this walkthrough we assume that you already have already prepared the following:

* A [Benchmark](../operator/docs/api.md#benchmark) for the benchmark you want to use in your experiment.
* A [FederatorTemplate](../operator/docs/api.md#federatortemplate) for the federation engine you want to use in your experiment.

We have already prepared several benchmarks and federator templates to use.
If you want to create your own dataset server template, check out [this guide](./BenchmarkWalkthrough.md).
Moreover, if you want to create your own federator template, check out [this guide](..).

## Step 1. Prepare your YAML file

## Step 2 - Optional. Define your own evaluator.

## Examples

We have already prepared several experiment specifications to experiment with:

* [experiment-fedbench](../examples/experiment-fedbench)
* [experiment-geofedbench](../examples/experiment-geofedbench)
* [experiment-geographica](../examples/experiment-geographica)
* [experiment-toybench](../examples/experiment-toybench)

> Notice: We plan to define more experiment specifications in the future.
> We place all experiment specifications in the [examples/](../examples/) directory
> under a subdirectory with the prefix `experiment-*`. 
