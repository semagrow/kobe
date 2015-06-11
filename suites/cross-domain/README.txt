Cross Domain Benchmark

This package contains the configuration for the cross domain benchmark.

 - Execute the setup from /suites/setup to prepare all native stores
 - Run the three scenarios
 	 1) Single Native Repository (runSingleNativeBenchmark)
 	 2) Federated Local (runFederatedBenchmark)
 	 3) Federated SPARQL (runFederatedSparqlBenchmark)

 - Configuration
 	* each query is executed three times for each scenario
 	* timeout is set to 10min = 600000ms
 	* results are copied to suites/cross-domain/result