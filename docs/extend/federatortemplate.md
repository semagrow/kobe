# FederatorTemplate Walkthough

This walkthrough illustrates the steps required from the *implementor of a
federation engine* in order to create a
[FederatorTemplate](../operator/docs/api.md#federatortemplate) specification.

In KOBE, a federator template is defined using a set of Docker images.
Additional parameters include the port and the path that the container will
listen for queries. Federator templates are used in the
[Experiment](../operator/docs/api.md#experiment) specifications in order to
define the federation engine to be benchmarked.

## Prerequisites

In this walkthrough we assume that you already have already prepared a Docker
image that provides the SPARQL endpoint of the federation engine (e.g.,
https://hub.docker.com/r/semagrow/semagrow/). Moreover, you should have a piece
of software that automatically constructs the configuration requred for your
federator to operate (e.g., https://github.com/semagrow/sevod-scraper). 

## Step 1. Optonal. Provide support for all evaluation metrics

> NOTICE: This step is optional, in a sense that it is only needed in order to
> support all evaluation metrics of KOBE. KOBE will be able to visualize 1) the
> number of retured results and 2) the total time to reciece the full result set
> to experimenter "for free" if you choose to not follow this step.

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

As an example, consider this [pull request](https://github.com/semagrow/semagrow/pull/52)
which contains the integration needed for KOBE in the case for the Semagrow
federation engine. 

## Step 2. Prepare your Docker images

Usually, a federation engine requires some configuration files that depend on
the federated endpoints (e.g., the URLs of the federated SPARQL endpoints).
Thus, apart from the Docker image with the SPARQL endpoint of the federation
engine, you should provide a docker image that constructs any desired
configuration for each of the source endpoints, and a docker image that
initializes the federator that possibly takes into account the configuration
files of the source endpoints. More specifically, prepare the following images:

* A docker image that constructs a configuration file for a source endpoint and
  places it in an *output directory* of your choice. Assume that the source
  endpoint and the dataset name are available in the environment variables
  `$DATASET_NAME` and `$DATASET_URL` respectively, and that the dump file of the
  dataset is present in an *input directory* of your choice.
* A docker image that constructs a configuration file for the federation engine
  and places it in an *output directory* of your choice. Assume that all the
  configuration files produced ine the previous step are present in an *input
  directory* of your choice.
* A docker image that starts the federation engine and exposes its SPARQL
  endpoint.

The environment variables are initialized by the Kobe operator according to the
specification of the benchmark to be executed.

As an example, we present the images for two federation engines (namely fedx and
Semagrow).

* Regarding the [Semagrow](http://semagrow.github.io/) federation engine, we use
  the images `semagrow/semagrow-init` (source code
  [here](../examples/federator-semagrow/semagrow-init)),
  `semagrow/semagrow-init-all`, (source code
  [here](../examples/federator-semagrow/semagrow-init-all)), and
  `semagrow/semagrow` (see [here](https://hub.docker.com/r/semagrow/semagrow/)).
  The first image uses the
  [sevod-scraper](https://github.com/semagrow/sevod-scraper) tool to create a
  ttl metadata file from the dump file, and the second image concatenates all
  metadata files of each of the source endpoints into a single metadata file.

* Regarding the
  [fedx](http://iswc2011.semanticweb.org/fileadmin/iswc/Papers/Research_Paper/05/70310592.pdf)
  federation engine, we use the images `semagrow/fedx-init` (source code
  [here](../examples/federator-fedx/fedx-init)), `semagrow/fedx-init-all`,
  (source code [here](../examples/federator-fedx/fedx-init-all)), and
  `semagrow/fedx-server` (source code
  [here](https://github.com/semagrow/docker-fedx-server)). Fedx is known for not
  using any dataset statistics, but it uses only a ttl file that contains only
  the SPARQL endpoints of the federation. The first image creates a ttl file
  that defines the SPARQL endpoint of each dataset and the second image
  concatenates all ttil files of each source endpoints into a single
  configuration file.


## Step 3. Prepare your YAML file

Once you have prepared the docker images, creating the federator template
specification for your dataset server is a straightforward task. It should look
like this (we use as an example the template for Semagrow):

```yaml
apiVersion: kobe.semagrow.org/v1alpha1
kind: FederatorTemplate
metadata:
  # Each federator template can be uniquely identified by its name.
  name: semagrowtemplate
spec:
  containers:
    # here you put the last image (that is the image for the
    # SPARQL endpoint of the federation engine)
    - name: maincontainer 
      image: semagrow/semagrow
      ports:
      - containerPort: 8080             # port to listen for queries
  port: 8080                            # port to listen for queries
  path: /SemaGrow/sparql                # path to listen for queries
  fedConfDir: /etc/default/semagrow     # where the federator expects to find its configuration
  
  # federator configuration step 1 (for each dataset):
  confFromFileImage: semagrow/semagrow-init  # first docker image
  inputDumpDir: /sevod-scraper/input         # where to find the dump file for the dataset
  outputDumpDir: /sevod-scraper/output       # where to place the configuration for the dataset
  
  # federator configuration step 2 (combination step):
  confImage: semagrow/semagrow-init-all      # second docker image
  inputDir: /kobe/input                      # where to find all dataset configurations
  outputDir: /kobe/output                    # where to place the final (combined) configuration

```

The default URL for the SPARQL endpoint for Virtuoso is
`http://localhost:8080/SemaGrow/sparql`, hence the port and the path to listen
for queries are `8080` and `/SemaGrow/sparql` respectively. The input and output
directories of the images mentioned previously are configured using the
parameters `inputDumpDir`,`outputDumpDir`,`inputDir`,`outputDir`.

## Examples

We have already prepared several federator template specifications to experiment
with:

* [federator-fedx](../examples/federator-fedx)
* [federator-semagrow](../examples/federator-semagrow)
* [federator-uno](../examples/federator-uno)

> Notice: We plan to define more federator template specifications in the
> future. We place all federator template specifications in the
> [examples/](../examples/) directory under a subdirectory with the prefix
> `federator-*`. 
