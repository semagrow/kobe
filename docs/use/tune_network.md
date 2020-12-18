# Tune network settings

This walkthrough illustrates the steps required from the *benchmark designer* in
order to configure the latency of the data sources of a
[Benchmark](../references/api.md#benchmark) specification.

## Prerequisites

In this walkthrough we assume that you already have already prepared the
following:

* A [Benchmark](../references/api.md#benchmark) for the benchmark you want to
  use in your experiment.

We have already prepared several benchmarks and federator templates to use. If
you want to create a benchmark specification, check out [this
guide](./create_benchmark.md).

## Step 1 - Inject latency for each source endpoint

KOBE allows simulating network traffic for all sources of the benchmark. For
every source dataset of the benchmark, you can:

* inject delay in the connection between the given source endpoint and the
  *federation engine*.
* inject delay in the connection between the given source endpoint and *another
  source endpoint*.

> The reason for injecting delays between the federated sources is the fact that
> every SPARQL endpoint can issue a SPARQL query to every other endpoint using
> the SERVICE SPARQL keyword.

The latency of each source can be configured using the following [delay
parameters](../references/api.md#delay). The functionality of these
parameters is offered by Istio. Check this
[link](https://istio.io/latest/docs/reference/config/networking/virtual-service/#HTTPFaultInjection-Delay)
for more information.

* The `fixedDelaySec` and `fixedDelayMSec` are used to indicate the *amount of
  delay* in seconds and in milliseconds.
* The `percentage` field can be used to only delay a certain *percentage of
  requests*.

You can extend your benchmark specification can be extended in order to define
the latency of the sources as follows:

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

Check the following link in which we illustrate a simple working example with
delays:

* [benchmark-toybench/toybench-delays.yaml](https://github.com/semagrow/kobe/tree/devel/examples/benchmark-toybench/toybench-delays.yaml)

This benchmark contains three SPARQL queries and two datasets (namely `toy1` and
`toy2`). All responces from `toy1` to the federator are delayed by 2 seconds and
150 milliseconds, all responces from `toy2` to the federator are delayed by 2
seconds, and the 50% of the responces from `toy1` to `toy2` are delayed by 3
seconds.

## Examples

We have already prepared a benchmark specification with delays to experiment with:

* [benchmark-toybench](https://github.com/semagrow/kobe/tree/devel/examples/benchmark-toybench)

> Notice: We plan to define more benchmark specifications with delays in the future. We
> place all benchmark specifications in the [examples/](https://github.com/semagrow/kobe/tree/devel/examples/) directory
> under a subdirectory with the prefix `benchmark-*`. 

