To build this go to the root directory of this repository and issue:

    docker build -t kobe-evaluator .

To run the container issue:

    docker run -e ENDPOINT=<ENDPOINT_TO_USE> kobe-evaluator

where `<ENDPOINT_TO_USE>` the SPARQL endpoint that will be evaluated. After the execution finishes the evaluation results will be printed.

To define your query set, mount the directory that contains the queries (e.g. `/path/to/queries`) to `/etc/queries`. `/path/to/queries` is expected to contain text files that each of them contains a query. In this case issue:

    docker run -e ENDPOINT=http://endpoint.org -v /path/to/queries:/etc/queries kobe-evaluator
