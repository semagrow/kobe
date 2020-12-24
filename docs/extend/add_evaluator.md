# Add a new evaluator

This walkthrough illustrates the steps required from the *experimenter* in order
to implement a custom evaluator for KOBE.

In KOBE, evaluators is defined using a single Docker image, and are used in the
[Experiment] specifications in order to apply a query load to the SPARQL
endpoint of a federation engine.

## Step 1. Prepare your Docker image

Given a set of query strings and the SPARQL endpoint of the federator, the
evaluator should issue each query of the queryset one or more times to the
endpoint. For each query, it should calculate the number of the returned results
and the time passed to receive the results, and export this information via a
log message.

KOBE provides the following input to the evaluator:

* All queries of the experiment are stored in the `/queries` directory. For each
  query, there exists a text file that contains the SPARQL query. The name of
  the file is the is used to uniquely identify each query of the experiment.
* The SPARQL endpoint of the federation engine appears in the `$ENDPOINT`
  environment variable.
* The environment variable `$EVAL_RUNS` specify how many times each query of the
  experiment should be executed.
* The name of the experiment to be executed appears in the `$EXPERIMENT`
  environment variable. This information is mainly used for logging purposes.
 
Given this information, the evaluator can proceed in the evaluation of the
queryset with its strategy of choice. However, the evaluator for each query
should do the following:

* Before executing each query, attach a SPARQL comment in the first line of the
  query. The comment should look like this:
  
  ```
  ^\#kobeQueryDesc Experiment: (?<experiment>[^ ]+) - Date: (?<date>[^ ]+ [^ ]+) - Query: (?<query>[^ ]+) - Run: (?<run>[0-9]+)$
  ```
  
  where
    * `<experiment>` is the name of the experiment, obtained by `$EXPERIMENT`.
    * `<date>` is the date and time that this specific experiment started - same
      for all queries of this experiment.
    * `<query>` is the filename of the query that is going to be executed.
    * `<run>`  which is a number between 1 and `$EVAL_RUNS` and is used to
      identify the run of this query.
 
* After executing each query, output a log message of the following form:
  
  ```
  ^I - [^ ]+ [^ ]+ - .{12} - .{20} - [^ ]+ - Experiment: (?<experiment>[^ ]+) - Date: (?<date>[^ ]+ [^ ]+) - Query: (?<query>[^ ]+) - Run: (?<run>[0-9]+) - Query Evaluation Time: (?<evaluation_time>[0-9]+) - Results: (?<results>[0-9]+)$
  ```
  
  where
    * `<experiment>`, `<date>`, `<query>`, `<run>` are defined as previously.
    * `<evaluation_time>` is the time passed in ms to get the full
       result set of the query.
    * `<results>` is number of returned results of the query.
  
!!! note
    This guide describes the steps for the standard evaluator metrics in KOBE.
    If you want to add more metrics from the side of the evaluator, please refer
    to this [guide](add_metrics.md).

## Example

As an example, we have prepared `semagrow/kobe-sequential-evaluator`, which
evaluates the queries in a sequential order, which is a slightly modified
version of that of [FedBench].

* [./evaluator](https://github.com/semagrow/kobe/tree/devel/evaluator)

We plan to implement more evaluators in the future.

[Experiment]: ../use/create_experiment.md
[FedBench]: https://code.google.com/archive/p/fbench/