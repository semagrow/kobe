# Experiment Walkthough

This walkthrough illustrates the steps required from the *experimenter* in order to create an [Experiment](../operator/docs/api.md#experiment) specification.

In KOBE, an experiment is defined over a benchmark and federator.
This resource provides the necessary parameters for instantiating a federation of querying systems.

## Prerequisites

In this walkthrough we assume that you already have already prepared the following:

* A [Benchmark](../operator/docs/api.md#benchmark) for the benchmark you want to use in your experiment.
* A [FederatorTemplate](../operator/docs/api.md#federatortemplate) for the federation engine you want to use in your experiment.
* A docker image of the evaluator, which is a piece of software that will pose the queries to the federator.

We have already prepared several benchmarks and federator templates to use.
If you want to create your own dataset server template, check out [this guide](./benchmarkWalkthrough.md).
Moreover, if you want to create your own federator template, check out [this guide](..).
Regarding the evaluator, we currently we offer the docker image `semagrow/kobe-sequential-evaluator`, which executes the queries of the benchmark in a sequential manner.



## Step 1. Prepare your YAML file

An experiment is characterized by its *name* and is parameterized with a *benchmark* to be executed and a *federator template* to be used.
A typical experiment specification should look like this:

```yaml
apiVersion: kobe.semagrow.org/v1alpha1
kind: Experiment
metadata:
  # Each experiment can be uniquely identified by its name.
  name: myexperiment
spec:
  # Specify the name of the benchmark to be executed.
  benchmark: mybench
  
  # Specify the name of the federation engine of the experiment.
  federatorName: myfederator
  
  # Specify the name of the federator template to be used.
  federatorTemplateRef: federatortemplate
  
  # Specify the docker image of the evaluator.
  evaluator:
    image: semagrow/kobe-sequential-evaluator
  
  # Specify the number of runs of the experiment, i.e. how many times each query 
  # of the benchmark should be executed.
  timesToRun: runs
  
  # If you set this parameter to true, KOBE will only build the federation 
  # and will not start the experiment.
  dryRun: false
  
  # If you set this parameter to false, KOBE will not build the federation
  # if it was already built in previous executions of this experiment.
  forceNewInit: true 
```

Check the following link in which we illustrate a simple example of the above specification:

* [experiment-toybench/toyexp-simple.yaml](../examples/experiment-toybench/toyexp-simple.yaml)

In this example, we define an experiment over the `toybench-simple` benchmark, and we use the Semagrow federation engine.
The queries of the benchmark are executed in a sequential manner, and each query of the benchmark is executed 3 times.
Since `toybench-simple` contains the queries `tq1`, `tq2`, `tq3`, in the example experiment the queries will be executed with the following order:
`tq1`, `tq2`, `tq3`, `tq1`, `tq2`, `tq3`, `tq1`, `tq2`, `tq3`.

## Examples

We have already prepared several experiment specifications to experiment with:

* [experiment-fedbench](../examples/experiment-fedbench)
* [experiment-geofedbench](../examples/experiment-geofedbench)
* [experiment-geographica](../examples/experiment-geographica)
* [experiment-toybench](../examples/experiment-toybench)

> Notice: We plan to define more experiment specifications in the future.
> We place all experiment specifications in the [examples/](../examples/) directory
> under a subdirectory with the prefix `experiment-*`. 
