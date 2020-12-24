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
