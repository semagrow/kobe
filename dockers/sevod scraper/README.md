# docker-sevod-scraper

This is the docker container for the sevod-scraper tool. For more info about the tool see https://github.com/semagrow/sevod-scraper

To build docker-sevod-scraper go into the clone directory and run

    docker build -t sevod-scraper .

To run it against a RDF dump issue

    docker run --rm -it -v </path/to/output>:/output -v </path/to/dump>:/dump sevod-scraper rdfdump /dump/<dump_name> <endpoint_url> <flags> /output/<output_name>

where:
* **/path/to/output** the directory to write the output
* **/path/to/dump** the directory that contains the dump
* **dump_name** the filename of the dump
* **endpoint_url** the endpoint URL where the dump is stored
* **flags** the flags to run sevod-scraper
* **output_name** the the filename of the output which will be located at **/path/to/output**/output_name

To run it against a Cassandra store issue

    docker run --rm -it -v </path/to/output>:/output sevod-scraper cassandra <cassandra_ip> <cassandra_port> <keyspace> <base_url> /output/<output_name>

where:
* **/path/to/output** the directory to write the output
* **cassandra_ip** the IP of the Cassandra store
* **cassandra_port** the port of the Cassandra store
* **keyspace** the Cassandra keyspace to scrap
* **base_url** the base url to use in the output
* **output_name** the the filename of the output which will be located at **/path/to/output**/output_name

Log properties can be changed in `log4j.properties` file.
