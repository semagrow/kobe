# Welcome to KOBE

KOBE is a benchmarking system that leverages
[Docker](https://docker.io) and [Kubernetes](https://kubernetes.io) in
order to reproduce experiments of federated query processing in a
collections of data sources.

## Overview

In the SPARQL query processing community, as well as in the wider
databases community, benchmark reproducibility is based on releasing
datasets and query workloads. However, this paradigm breaks down for
federated query processors, as these systems do not manage the data
they serve to their clients but provide a data-integration abstraction
over the actual query processors that are in direct contact with the
data.

The KOBE benchmarking engine is a system that aims to provide a
generic platform to perform benchmarking and experimentation that can
be reproducible in different environments. It was designed with the
following objectives in mind:

1. to allow for benchmark and experiment specifications to be
   reproduced in different environments and be able to produce
   comparable and reliable results;
2. to ease the deployment of complex benchmarking experiments by
   automating the tedious tasks of initialization and execution.

