# Introduction #
This tool aims to provide a quick and easy way to set up a benchmark for federation engines. By following the provided step by step procedure one can create a dockec-compose.yml file that
is used to deploy a benchmark using Docker containers. The user can deploy one of the preconfigured datasets (OPS, Cross Domain and Life Science Domain) or use other datasets. Each dataset by default is loaded in a Virtuoso store and sets up a SPARQL endpoint at http://<name>:8890/sparql, where <name> the name provided by the user for each dataset. The datasets can then be federated using Semagrow or Fedx.  Both Semagrow and Fedx also set up a SPARQL endpoint at http://<name>:8080/sparql. This endpoint can be used to apply a query set in order to benchmark the federation.

# Requirements #
To build and run this you need Java 7+ and Maven 3+.

# Build #
To build this issue

```
mvn clean package
```
# Usage #
The federation-composer.sh script is used to create and modify a docker-compose.yml file. The script and the produced file are located under the root directory.

To see usage issue
```
./federation-composer.sh -h
```

To create a docker compose file issue
```
./federation-composer.sh -a
```
Then you are have to choose from one of the following options:

```
(1) add a data node in the compose file. You have to provide a name and either URL to download the dump or the directory that contains the dump.
(2) add Semagrow. You have to provide a name and the directory that contains the semagrow configuration file(s).
(3) add Fedx. You have to provide a name and the directory that contains the fedx configuration file.
(4) add OPF. Will add to the compose file data nodes with OPF dumps loaded.
(5) add Cross Domain. Will add to the compose file data nodes with Cross Domain dumps loaded.
(6) add Life Science Domain. Will add to the compose file data nodes with Life Science Domain dumps loaded.
```


To list all currently added container issue
```
./federation-composer.sh -l
```

To remove a container issue
```
./federation-composer.sh -r <container_name>
```

To remove all container issue
```
./federation-composer.sh -c
```

# Notes #
The produced docker-compose.yml file can be also modified without using the federation-composer.sh script but the script will overwrite all changes that where not created using it.

You can change the default images used in the above operations by editing the default_images file.
