To build this go to the root directory of this repository and issue:

    docker build -t kove-evaluator .

To run the container issue:

    docker run -e ENDPOINT=<ENDPOINT_TO_USE> kobe-evaluator

where `<ENDPOINT_TO_USE>` the SPARQL endpoint that will be evaluated. After the execution finishes the evaluation results will be printed.

To define query sets to run set as true the following available variables:

 - **CROSS_DOMAIN**, for the Cross Domain dataset
 - **LIFE_SCIENCE**, fot the Life Science dataset
 - **OPFBENCH**, for the OPFbench dataset
 - **SGPILOTS**, for the Semagrow pilots dataset

for example to run the evaluator on endpoint `http://endpoint.org` with query sets Cross Domain and Life Science issue:

    docker run -e ENDPOINT=http://endpoint.org -e CROSS_DOMAIN=true -e LIFE_SCIENCE=true kobe-evaluator

You can also use your own query sets by mounting a directory (e.g. `/path/to/queries`) to `/etc/querySet`. `/path/to/queries` is expected to contain text files that each of them contains a query. In this case issue:

    docker run -e ENDPOINT=http://endpoint.org -v /path/to/queries:/etc/querySet kobe-evaluator
