
## Typical workflow

The typical workflow of defining a KOBE experiment is the following.

1. Create one [DatasetTemplate](operator/docs/api.md#datasettemplate) for each
   dataset server you want to use in your benchmark.

2. Define your [Benchmark](operator/docs/api.md#benchmark), which should contain
   a list of datasets and a list of queries.

3. Create one [FederatorTemplate](operator/docs/api.md#federatortemplate) for
   the federator engine you want to use in your experiment. 

4. Define an [Experiment](operator/docs/api.md#experiment) over your previously
   defined benchmark.

Several examples of the above specifications can be found in the
[examples](examples/) directory.

In the following, we show the steps for deploying an experiment on a simple
benchmark that comprises three queries over a Semagrow federation of two
Virtuoso endpoints.

You can use the `kobectl` script found in the [bin](bin/) directory for
cotrolling your experiments:

```
export PATH=`pwd`/bin:$PATH
kobectl help
```

First, apply the templates for Virtuoso and Semagrow:

```
kobectl apply examples/dataset-virtuoso/virtuosotemplate.yaml
kobectl apply examples/federator-semagrow/semagrowtemplate.yaml
```
Then, apply the benchmark.

```
kobectl apply examples/benchmark-toybench/toybench-simple.yaml
```

Before running the experiment, you should verify that the datasets are loaded.
Use the following command:

```
kobectl show benchmark toybench-simple
```

When the datasets are loaded, you should get the following output:

```
NAME  STATUS
toy1  Running
toy2  Running
```

Proceed now with the execution of the experiment:

```
kobectl apply examples/experiment-toybench/toyexp-simple.yaml
```

As previously, you can review the state of the experiment with the following
command:

```
kobectl show experiment toyexp-simple
```
You can now view the evaluation metrics in the Kibana dashboards.

For removing all of the above, issue the following commands:
```
kobectl delete experiment toyexp-simple
kobectl delete benchmark toybench-simple
kobectl delete federatortemplate semagrowtemplate
kobectl delete datasettemplate virtuosotemplate
```
For more advanced control options for KOBE, use [kubectl].

[kubectl]: https://kubernetes.io/docs/reference/kubectl/overview/
