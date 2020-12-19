# Provide metric support

This walkthrough illustrates the steps required from the *implementor of a
federation engine* in order to provide full support for the available
benchmark metrics.

> NOTICE: This step is optional, in a sense that it is only needed in order to
> support all evaluation metrics of KOBE. KOBE will be able to visualize 1) the
> number of retured results and 2) the total time to reciece the full result set
> to experimenter "for free" if you choose to not follow this step.

## Logging subsystem concepts and available metrcs

One important feature of Kobe is that the experimenter can have easy access on
a set of several statistics and key performance indicators for each conducted
experiment. The metrics currently supported are the following:

* Number of returned results
* Total time to recieve the full result set
* Source selection time
* Query planning time
* Query execution time
* Number of sources accessed

Of these evaluation metrics, only the first two can can be computed by the
client side. Thus, the remaining metrics should be calculated by the federation
engine itself and can be presented via a log message. However, in order for
Kobe to be able to link the log message with its corresponding experiment
execution and with its specific query run, the log message should contain also
the following information:

* Experiment name
* Start time of the experiment
* Query name
* Run

These parameters are passed from the evaluator of Kobe to the federation engine
via a SPARQL comment that is attached in the query string.

## Step 1. Provide support for all evaluation metrics

In order to provide full support for evaluation metrics, you should do the
following:

1. Extend your query string parser to parse the first line of the query string
   which follows the according regex pattern:
   
   ```
   ^\#kobeQueryDesc Experiment: (?<experiment>[^ ]+) - Date: (?<date>[^ ]+ [^ ]+) - Query: (?<query>[^ ]+) - Run: (?<run>[0-9]+)$
   ```
   
2. Calculate some or all of the metrics discussed previously.
3. Provide a log message to output the metrics according to the following regex
   pattern:

   ```
   ^I - [^ ]+ [^ ]+ - .{12} - .{20} - [^ ]+ -  - Experiment: (?<experiment>[^ ]+) - Date: (?<date>[^ ]+ [^ ]+) - Query: (?<query>[^ ]+) - Run: (?<run>[0-9]+) - Source Selection Time: (?<source_selection_time>[0-9]+) - Compile Time: (?<compile_time>[0-9]+) - Sources: (?<sources>[0-9]+) - Execution time: (?<execution_time>[0-9]+)$
   ```
   
   where
  * `<experiment>`, `<date>`, `<query>`, `<run>` are obtained by the SPARQL
     comment of the query string.
  * `<source_selection_time>` is the time to perform source selection.
  * `<compile_time>` is the time to provide a query execution plan.
  * `<execution_time>` is the time to execute the plan.
  * `<sources>` is the number of sources that appear in the query plan.


## Example

As an example, consider this [pull request](https://github.com/semagrow/semagrow/pull/52)
which contains the integration needed for KOBE in the case for the Semagrow
federation engine. 
