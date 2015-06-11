@echo off
cd ..\..
echo Starting local SPARQL endpoints

REM DBPedia351 => 10000
echo Starting DBPedia Endpoint on port 10000
start startSparqlEndpoint.bat data\repositories\native-storage.dbpedia351 10000

REM Nytimes => 10001
echo Starting Nytimes Endpoint on port 10001
start startSparqlEndpoint.bat data\repositories\native-storage.nytimes 10001

REM LinkedMDB => 10002
echo Starting LinkedMDB Endpoint on port 10002
start startSparqlEndpoint.bat data\repositories\native-storage.linkedmdb 10002

REM Jamendo => 10003
echo Starting Jamendo Endpoint on port 10003
start startSparqlEndpoint.bat data\repositories\native-storage.jamendo 10003

REM Geonames => 10004
echo Starting DBPedia Geonames on port 10004
start startSparqlEndpoint.bat data\repositories\native-storage.geonames 10004
