# Design

This document describes the different resources, their purpose and the control
flow between the [Kobe operator](https://github.com/semagrow/kobe) and
[Kubernetes](https://kubernetes.io).

## Overview

<!--
[Kobe](https://github.com/semagrow/kobe) is a system that helps with the
benchmarking of federating query engines (or, for short, federators). In a
nutshell, a federator is an engine that sits between the user and a set of
databases and facilitates the query answering by delegating the work to the
appropriate databases of the federation. Benchmarking the performance of a
federator can become a tedious and possibly unreliable task because it requires
the deployment and configuration of multiple systems and a careful consideration
of their allocated resources. 
-->

[Kobe](https://github.com/semagrow/kobe) is a system that can easily setup the
ensemble of systems and orchestrate experiments for benchmarking federated query
engines. Kobe uses containers to abstract the installation of the systems and
their dependencies. [Kubernetes](https://kubernetes.io) is also employed for the
deployment of containers in a cluster of computers.

The [Kobe operator](https://github.com/semagrow/kobe/operator) is a custom
[Kubernetes
operator](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/) that
manages the resources deployed in the Kubernetes cluster 

## Actors

Kobe distinguishes several distinct roles for performing an experiment. Those
are 
* The *benchmark designer* that provides the queries to be executed in a set of
  datasets. The benchmark designer should also define the datasets that need to
  be loaded for that benchmark.
* The *implementor of a federation engine* that provides the instrumentation
  needed for the federator in order to be deployed, configured and monitored
  properly. 
* The *experimenter* that performs an experiment. The experimenter provides the
  additional details pertaining the experiment, such as, the systems to be used
  for each dataset of the experiment, the computational resources as well as the
  bandwidth that should be allocated in each system.

Each actor is responsible for providing part of the information needed for an
experiment. Information is expressed Kubernetes manifests and can be serialized
and redistributed as YAML files. 

## Model

A benchmarking experiment is described in Kobe in a set of resources. The main
resources that constitute the public API of Kobe are:
* The *Benchmark* that describe the query set (a ordered list of query strings)
  and provide information about the datasets needed to answer those queries.
  Each benchmark is properly designed to test certain characteristics and
  functionality of a federator.
* The *Federator* template that describes the software that should be used in
  order to initialize and run a federation engine. The software needed is
  expected to be provided in container images. This software include the
  initializer and the service that listen to queries.
* The *DatasetServer* template that describes the software that should be used 
  to load and serve a dataset. This software include the initializer and the 
  service that listen to queries. 
* The *Experiment* that describes a specific experiment of a given Benchmark and
  Federator. This resource provide the necessary parameters for instantiating a
  federation of querying systems. The experiment also provides an evaluator,
  that is a piece of software that will pose the queries to the federator.

The detailed schemas of each resource are presented in the [API](api.md).

## Flow

The aforementioned resources contain all the necessary information to setup and 
perform a benchmarking experiment. These resources are committed to a Kubernetes
cluster and monitored by (an already deployed) kobe operator.

The sequence of events of a typical flow are:
1. The `kobe-operator` is deployed in a Kubernetes cluster.
2. A set of `DatasetServer`s and `Federator`s are committed.
3. A `Benchmark` resource is applied in Kubernetes.
4. An `Experiment` is applied that must reference one of the `Federator`s and
   `Benchmark`s

Upon submission of the `Experiment` a series of actions occur from the
kobe-operator. 

* The operator checks whether the federator and the dataset servers are already
  defined.
* If datasets are not served by the specific `DatasetServer` required by the
  `Experiment`, then they are initialized.
* If there is no federation the controller creates a `Federation`, an internal
  resource that instantiates a `Federator` template with a specific set of
  datasets. 

When the `Federation` status changed to `Running` the operator starts an 
evaluator that receives as input the query-set defined in the `Benchmark`
and the endpoint of the federator. The evaluator completes after a predefined
number of runs. In each run it sends the queries to the federator and waits
for the responses.
