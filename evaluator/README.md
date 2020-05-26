# KoBE QuerySet Evaluator

Query evaluator of the KoBE Benchmark Engine.

### Build
To build and run this you need Java 11+ and Maven 3+.

```
mvn clean package
mvn dependency:copy-dependencies
```

### Usage

In order to run an experiment issue the following command:

```
./runEval.sh [sparqlendpoint url] [properties file]
```

where [sparqlendpoint url] is the sparql endpoint of the federation engine, and [properties file] is the (optional)
configuration file for the experiment, where you can define the queryset, the number of the runs for each cycle of the
experiment, and the timeout for each query, etc. (see config.prop for a simple example). The default path of the
queryset is the /etc/querySet/ directory.

### Results

The results of the evaluation are placed in result/ directory in two csv files.
